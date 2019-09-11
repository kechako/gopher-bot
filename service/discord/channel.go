package discord

type channel struct {
	id   string
	name string
}

// ID implements the plugin.Channel interface.
func (ch *channel) ID() string {
	return ch.id
}

// Name implements the plugin.Channel interface.
func (ch *channel) Name() string {
	return ch.name
}
