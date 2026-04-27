module github.com/plexusone/omni-mixpanel/omnidxi

go 1.22

require (
	github.com/mixpanel/mixpanel-go/v2 v2.0.0
	github.com/plexusone/omnidxi-core v0.1.0
)

require (
	github.com/barkimedes/go-deepcopy v0.0.0-20220514131651-17c30cfc62df // indirect
	github.com/diegoholiveira/jsonlogic/v3 v3.9.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
)

// Local development - remove before release
replace github.com/plexusone/omnidxi-core => ../../omnidxi-core
