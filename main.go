package main

import (
	"fmt"
	"os"
)

func main() {
	includeBuiltin := false
	org := os.Getenv("DD_ORG")

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--include-builtin":
			includeBuiltin = true
		case "--org":
			if i+1 < len(os.Args) {
				org = os.Args[i+1]
				i++
			}
		case "--help", "-h":
			fmt.Fprintln(os.Stderr, "Usage: pup rbac dump [--org <name>] [--include-builtin]")
			fmt.Fprintln(os.Stderr, "  Fetches all Datadog roles with enriched permission sets.")
			fmt.Fprintln(os.Stderr, "  By default, only custom roles are returned.")
			os.Exit(0)
		}
	}

	runDump(org, includeBuiltin)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "pup-rbac: "+format+"\n", args...)
	os.Exit(1)
}
