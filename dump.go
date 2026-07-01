package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/rapdev-io/pup-rbac/internal/pupapi"
)

func runDump(org string, includeBuiltin bool) {
	client := pupapi.New()

	// Step 1: permission catalog for UUID resolution
	var permsData permsResp
	if err := client.Get("v2/permissions", &permsData); err != nil {
		fatalf("fetching permissions: %v", err)
	}
	permLookup := make(map[string]Permission, len(permsData.Data))
	for _, p := range permsData.Data {
		permLookup[p.ID] = Permission{
			Name:        p.Attributes.Name,
			Group:       p.Attributes.GroupName,
			DisplayName: p.Attributes.DisplayName,
		}
	}

	// Step 2: all roles (paginated)
	rawRoles, err := client.Paginate("v2/roles")
	if err != nil {
		fatalf("fetching roles: %v", err)
	}
	var allRoles []roleEntry
	for _, raw := range rawRoles {
		var r roleEntry
		if err := json.Unmarshal(raw, &r); err == nil {
			allRoles = append(allRoles, r)
		}
	}

	var targets []roleEntry
	for _, r := range allRoles {
		if !r.Attributes.BuiltIn || includeBuiltin {
			targets = append(targets, r)
		}
	}

	// Step 3: permissions per role (concurrent, max 10 in flight)
	type result struct {
		perms []Permission
		err   error
	}
	results := make([]result, len(targets))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for i, role := range targets {
		wg.Add(1)
		go func(idx int, roleID string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			var resp struct {
				Data []struct {
					ID string `json:"id"`
				} `json:"data"`
			}
			if err := client.Get(fmt.Sprintf("v2/roles/%s/permissions", roleID), &resp); err != nil {
				results[idx] = result{err: err}
				return
			}
			var perms []Permission
			for _, p := range resp.Data {
				if perm, ok := permLookup[p.ID]; ok {
					perms = append(perms, perm)
				}
			}
			results[idx] = result{perms: perms}
		}(i, role.ID)
	}
	wg.Wait()

	// Assemble
	out := DumpOutput{
		Org:        org,
		TotalRoles: len(allRoles),
		Roles:      make([]Role, 0, len(targets)),
	}
	for i, target := range targets {
		if !target.Attributes.BuiltIn {
			out.CustomRoles++
		}
		perms := results[i].perms
		if perms == nil {
			perms = []Permission{}
		}
		out.Roles = append(out.Roles, Role{
			ID:          target.ID,
			Name:        target.Attributes.Name,
			BuiltIn:     target.Attributes.BuiltIn,
			UserCount:   target.Attributes.UserCount,
			Permissions: perms,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fatalf("encoding output: %v", err)
	}
}
