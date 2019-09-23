package discord

type user struct {
	id   string
	name string
}

// ID implements the plugin.User interface.
func (u *user) ID() string {
	return u.id
}

// Name implements the plugin.User interface.
func (u *user) Name() string {
	return u.name
}
