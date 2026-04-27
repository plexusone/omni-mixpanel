# Release Notes - v0.1.0

**Release Date:** 2026-04-27

## Overview

Initial release of omni-mixpanel, the Mixpanel adapter for the omnidxi DXI abstraction layer.

This package wraps the official [mixpanel/mixpanel-go](https://github.com/mixpanel/mixpanel-go) SDK (v2) and implements the omnidxi-core `Tracker` interface, enabling Mixpanel integration with the unified DXI API.

## Highlights

- Mixpanel DXI adapter wrapping mixpanel-go/v2
- Full implementation of omnidxi-core Tracker interface

## What's Included

### Tracker Implementation

Complete implementation of all Tracker interface methods:

| Method | Mixpanel Mapping |
|--------|------------------|
| `Track()` | `client.Track()` with NewEvent |
| `Identify()` | `client.PeopleSet()` with reserved properties |
| `Group()` | `client.GroupSet()` |
| `Alias()` | `client.Alias()` |
| `Flush()` | No-op (Mixpanel SDK sends immediately) |
| `Close()` | Cleanup resources |

### Features

- **Event Tracking** - NewEvent with timestamp support for accurate event timing
- **User Identification** - PeopleSet with reserved property mapping ($name, $email, etc.)
- **Group Analytics** - GroupSet for account-level tracking
- **Identity Linking** - Native Alias method for connecting anonymous to identified users
- **Insert ID Support** - Automatic deduplication via insert_id
- **Event Context Mapping** - Page, device, app, and UTM metadata mapped to Mixpanel properties

### Usage

```go
import (
    "context"
    mixpanel "github.com/plexusone/omni-mixpanel/omnidxi"
    "github.com/plexusone/omnidxi-core"
)

// Create tracker
tracker, err := mixpanel.New(
    mixpanel.WithToken("your-mixpanel-token"),
)
if err != nil {
    log.Fatal(err)
}
defer tracker.Close()

// Track events using canonical schema
err = tracker.Track(ctx, core.Event{
    Type:   core.EventTypeUIClick,
    Name:   "signup_button_clicked",
    UserID: "user-123",
    Properties: map[string]any{
        "button_location": "header",
    },
})
```

## Requirements

- Go 1.22 or later

## Dependencies

- `github.com/mixpanel/mixpanel-go/v2` v2.0.0
- `github.com/plexusone/omnidxi-core` v0.1.0

## Installation

```bash
go get github.com/plexusone/omni-mixpanel@v0.1.0
```

## Related Packages

- [omnidxi-core](https://github.com/plexusone/omnidxi-core) - Core interfaces and types
- [omnidxi](https://github.com/plexusone/omnidxi) - Batteries-included client with MultiTracker
- [omni-amplitude](https://github.com/plexusone/omni-amplitude) - Amplitude adapter
