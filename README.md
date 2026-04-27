# Omni-Mixpanel

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/omni-mixpanel/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/omni-mixpanel/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/omni-mixpanel/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/omni-mixpanel/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/omni-mixpanel/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/omni-mixpanel/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/omni-mixpanel
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/omni-mixpanel
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/omni-mixpanel
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/omni-mixpanel
 [docs-mkdoc-svg]: https://img.shields.io/badge/docs-guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.github.io/omni-mixpanel
 [viz-svg]: https://img.shields.io/badge/repo-visualization-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fomni-mixpanel
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omni-mixpanel/blob/main/LICENSE

Mixpanel integrations for the PlexusOne omni* family.

## Packages

| Package | Description |
|---------|-------------|
| `omnidxi` | DXI adapter for Mixpanel analytics |

## Installation

```bash
go get github.com/plexusone/omni-mixpanel/omnidxi
```

## Usage

```go
import (
    "context"

    "github.com/plexusone/omni-mixpanel/omnidxi"
    core "github.com/plexusone/omnidxi-core"
)

func main() {
    ctx := context.Background()

    // Create Mixpanel tracker
    tracker := omnidxi.New(core.WithAPIKey("your-mixpanel-token"))
    defer tracker.Close()

    // Track an event
    event := core.NewEvent(core.EventTypePageView, "Home Viewed").
        WithUserID("user_123").
        WithProperty("source", "direct")

    tracker.Track(ctx, event)

    // Identify user
    user := core.NewUser("user_123").
        WithTraits(core.UserTraits{
            Email: "user@example.com",
            Name:  "Jane Doe",
        })

    tracker.Identify(ctx, user)
}
```

## Related Projects

- [omnidxi-core](https://github.com/plexusone/omnidxi-core) - Core interfaces
- [omnidxi](https://github.com/plexusone/omnidxi) - Batteries-included package
- [omni-amplitude](https://github.com/plexusone/omni-amplitude) - Amplitude adapter

## License

MIT
