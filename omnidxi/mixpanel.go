// Package omnidxi provides a Mixpanel implementation of the omnidxi Tracker interface.
package omnidxi

import (
	"context"
	"sync"

	"github.com/mixpanel/mixpanel-go/v2"

	core "github.com/plexusone/omnidxi-core"
)

const (
	providerName = "mixpanel"
	version      = "0.1.0"
)

// Tracker wraps the Mixpanel mixpanel-go client with the omnidxi interface.
type Tracker struct {
	client *mixpanel.ApiClient
	config core.Config
	closed bool
	mu     sync.RWMutex
}

// New creates a new Mixpanel tracker with the given options.
func New(opts ...core.Option) *Tracker {
	cfg := core.NewConfig(opts...)

	client := mixpanel.NewApiClient(cfg.APIKey)

	return &Tracker{
		client: client,
		config: cfg,
	}
}

// Track sends an event to Mixpanel.
func (t *Tracker) Track(ctx context.Context, event core.Event) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.closed {
		return core.ErrClosed
	}
	if t.config.Disabled {
		return nil
	}

	props := t.buildProperties(event)
	mpEvent := t.client.NewEvent(t.eventName(event), t.distinctID(event), props)

	// Add timestamp if set
	if !event.Timestamp.IsZero() {
		mpEvent.AddTime(event.Timestamp)
	}

	// Add insert ID for deduplication
	if event.ID != "" {
		mpEvent.AddInsertID(event.ID)
	}

	if err := t.client.Track(ctx, []*mixpanel.Event{mpEvent}); err != nil {
		return core.NewProviderError(providerName, "track", err)
	}

	return nil
}

// Identify associates user properties with a user.
func (t *Tracker) Identify(ctx context.Context, user core.User) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.closed {
		return core.ErrClosed
	}
	if t.config.Disabled {
		return nil
	}

	props := t.buildUserProperties(user)
	peopleProps := mixpanel.NewPeopleProperties(user.ID, props)

	if err := t.client.PeopleSet(ctx, []*mixpanel.PeopleProperties{peopleProps}); err != nil {
		return core.NewProviderError(providerName, "identify", err)
	}

	return nil
}

// Group associates a user with a group.
func (t *Tracker) Group(ctx context.Context, group core.Group) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.closed {
		return core.ErrClosed
	}
	if t.config.Disabled {
		return nil
	}

	props := map[string]any{}
	if group.Name != "" {
		props["$name"] = group.Name
	}
	for k, v := range group.Traits {
		props[k] = v
	}

	if err := t.client.GroupSet(ctx, "company", group.ID, props); err != nil {
		return core.NewProviderError(providerName, "group", err)
	}

	// Also set the group on the user
	userProps := mixpanel.NewPeopleProperties(group.UserID, map[string]any{
		"$group_company": group.ID,
	})
	if err := t.client.PeopleSet(ctx, []*mixpanel.PeopleProperties{userProps}); err != nil {
		return core.NewProviderError(providerName, "group_user", err)
	}

	return nil
}

// Alias links two user identities.
func (t *Tracker) Alias(ctx context.Context, alias core.Alias) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.closed {
		return core.ErrClosed
	}
	if t.config.Disabled {
		return nil
	}

	// Mixpanel's Alias method
	if err := t.client.Alias(ctx, alias.PreviousID, alias.UserID); err != nil {
		return core.NewProviderError(providerName, "alias", err)
	}

	return nil
}

// Flush is a no-op for Mixpanel as it doesn't buffer client-side.
func (t *Tracker) Flush(ctx context.Context) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.closed {
		return core.ErrClosed
	}

	// Mixpanel API client doesn't buffer, so nothing to flush
	return nil
}

// Close shuts down the tracker.
func (t *Tracker) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true
	return nil
}

// Info returns metadata about this tracker.
func (t *Tracker) Info() core.TrackerInfo {
	return core.TrackerInfo{
		Name:       providerName,
		Version:    version,
		SDKVersion: "2.0.0",
	}
}

// eventName returns the event name, using Name if set or Type as fallback.
func (t *Tracker) eventName(e core.Event) string {
	if e.Name != "" {
		return e.Name
	}
	return string(e.Type)
}

// distinctID returns the user identifier for Mixpanel.
func (t *Tracker) distinctID(e core.Event) string {
	if e.UserID != "" {
		return e.UserID
	}
	return e.SessionID
}

// buildProperties builds Mixpanel event properties from an omnidxi Event.
func (t *Tracker) buildProperties(e core.Event) map[string]any {
	props := make(map[string]any)

	// Copy event properties
	for k, v := range e.Properties {
		props[k] = v
	}

	// Add standard fields
	if e.SessionID != "" {
		props["session_id"] = e.SessionID
	}
	if e.Type != "" {
		props["event_type"] = string(e.Type)
	}

	// Add context if present
	if e.Context != nil {
		t.addContextProperties(props, e.Context)
	}

	return props
}

// buildUserProperties builds Mixpanel people properties from an omnidxi User.
func (t *Tracker) buildUserProperties(u core.User) map[string]any {
	props := make(map[string]any)

	// Standard Mixpanel people properties use $ prefix
	if u.Traits.Email != "" {
		props["$email"] = u.Traits.Email
	}
	if u.Traits.Name != "" {
		props["$name"] = u.Traits.Name
	}
	if u.Traits.Phone != "" {
		props["$phone"] = u.Traits.Phone
	}
	if !u.Traits.CreatedAt.IsZero() {
		props["$created"] = u.Traits.CreatedAt.Format("2006-01-02T15:04:05")
	}
	if u.Traits.City != "" {
		props["$city"] = u.Traits.City
	}
	if u.Traits.Region != "" {
		props["$region"] = u.Traits.Region
	}
	if u.Traits.Country != "" {
		props["$country_code"] = u.Traits.Country
	}

	// Custom traits without $ prefix
	if u.Traits.Company != "" {
		props["company"] = u.Traits.Company
	}
	if u.Traits.Title != "" {
		props["title"] = u.Traits.Title
	}
	if u.Traits.Plan != "" {
		props["plan"] = u.Traits.Plan
	}
	if u.Traits.TotalSpend > 0 {
		props["total_spend"] = u.Traits.TotalSpend
	}

	// Extra traits
	for k, v := range u.Extra {
		props[k] = v
	}

	return props
}

// addContextProperties adds context fields to event properties.
func (t *Tracker) addContextProperties(props map[string]any, ctx *core.EventContext) {
	if ctx.PagePath != "" {
		props["page_path"] = ctx.PagePath
	}
	if ctx.PageTitle != "" {
		props["page_title"] = ctx.PageTitle
	}
	if ctx.PageURL != "" {
		props["$current_url"] = ctx.PageURL
	}
	if ctx.PageReferrer != "" {
		props["$referrer"] = ctx.PageReferrer
	}
	if ctx.DeviceType != "" {
		props["device_type"] = ctx.DeviceType
	}
	if ctx.Platform != "" {
		props["$os"] = ctx.Platform
	}
	if ctx.AppVersion != "" {
		props["app_version"] = ctx.AppVersion
	}

	// UTM parameters - Mixpanel standard names
	if ctx.UTMSource != "" {
		props["utm_source"] = ctx.UTMSource
	}
	if ctx.UTMCampaign != "" {
		props["utm_campaign"] = ctx.UTMCampaign
	}
	if ctx.UTMMedium != "" {
		props["utm_medium"] = ctx.UTMMedium
	}
	if ctx.UTMTerm != "" {
		props["utm_term"] = ctx.UTMTerm
	}
	if ctx.UTMContent != "" {
		props["utm_content"] = ctx.UTMContent
	}
}

// Ensure Tracker implements the interfaces.
var (
	_ core.Tracker         = (*Tracker)(nil)
	_ core.TrackerWithInfo = (*Tracker)(nil)
)
