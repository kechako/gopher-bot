package slack

type user struct {
	id          string
	name        string
	fullName    string
	displayName string
}

// ID implements the plugin.User interface.
func (u *user) ID() string {
	return u.id
}

// Name implements the plugin.User interface.
func (u *user) Name() string {
	return u.name
}

// FullName implements the plugin.User interface.
func (u *user) FullName() string {
	if u.fullName == "" {
		return u.name
	}
	return u.fullName
}

// DisplayName implements the plugin.User interface.
func (u *user) DisplayName() string {
	if u.displayName == "" {
		return u.FullName()
	}
	return u.displayName
}
