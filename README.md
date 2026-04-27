# omni-mixpanel

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
