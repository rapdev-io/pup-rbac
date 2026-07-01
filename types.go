package main

// ── Datadog API response types ───────────────────────────────────────────────

type permAttr struct {
	Name        string `json:"name"`
	GroupName   string `json:"group_name"`
	DisplayName string `json:"display_name"`
}

type permEntry struct {
	ID         string   `json:"id"`
	Attributes permAttr `json:"attributes"`
}

type permsResp struct {
	Data []permEntry `json:"data"`
}

type roleAttr struct {
	Name      string `json:"name"`
	BuiltIn   bool   `json:"built_in"`
	UserCount int    `json:"user_count"`
}

type roleEntry struct {
	ID         string   `json:"id"`
	Attributes roleAttr `json:"attributes"`
}

// ── Output types ─────────────────────────────────────────────────────────────

type Permission struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	DisplayName string `json:"display_name"`
}

type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	BuiltIn     bool         `json:"built_in"`
	UserCount   int          `json:"user_count"`
	Permissions []Permission `json:"permissions"`
}

type DumpOutput struct {
	Org         string `json:"org"`
	TotalRoles  int    `json:"total_roles"`
	CustomRoles int    `json:"custom_roles"`
	Roles       []Role `json:"roles"`
}
