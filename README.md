# pup-rbac

RapDev pup extension — fetches all Datadog roles with fully resolved permission sets in one call.

## Install

```bash
pup extension install rapdev-io/pup-rbac
```

> **Private repo**: prefix with `GH_TOKEN=$(gh auth token)` if pup returns "no releases found".

## Usage

```bash
pup rbac dump [--org <name>] [--include-builtin]
```

Output (JSON):
```json
{
  "org": "acme",
  "total_roles": 12,
  "custom_roles": 9,
  "roles": [
    {
      "id": "...",
      "name": "RD-Advanced",
      "built_in": false,
      "user_count": 5,
      "permissions": [
        { "name": "dashboards_read", "group": "Dashboards", "display_name": "Dashboards Read" }
      ]
    }
  ]
}
```

By default only custom roles are returned. Pass `--include-builtin` to include Datadog's built-in roles.

## Auth

Credentials are forwarded automatically by pup via environment variables (`DD_ACCESS_TOKEN`, `DD_API_KEY`, `DD_APP_KEY`, `DD_SITE`, `DD_ORG`). No configuration needed beyond a valid `pup auth` session.

## Release

Push a `v*` tag to trigger a GoReleaser build:

```bash
git tag v0.1.0 && git push origin v0.1.0
```

Release assets are standalone binaries: `pup-rbac-<os>-<arch>`.
