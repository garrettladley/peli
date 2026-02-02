package tui

// Config holds all configurable values for the TUI.
type Config struct {
	NameCharLimit        int
	CuisineCharLimit     int
	InputWidth           int
	InitialElo           float64
	InitialScore         float64
	ViewportPadding      int
	ProgressBarWidth     int
	CardWidth            int
	NameTruncateWidth    int
	CuisineTruncateWidth int
	CompactWindowContext int
	ResultWindowContext  int
}

// Option is a functional option for configuring the TUI.
type Option func(*Config)

// DefaultConfig returns the default TUI configuration.
func DefaultConfig() Config {
	return Config{
		NameCharLimit:        50,
		CuisineCharLimit:     30,
		InputWidth:           38,
		InitialElo:           1500.0,
		InitialScore:         5.0,
		ViewportPadding:      8,
		ProgressBarWidth:     30,
		CardWidth:            24,
		NameTruncateWidth:    22,
		CuisineTruncateWidth: 12,
		CompactWindowContext: 2,
		ResultWindowContext:  5,
	}
}

// WithNameCharLimit sets the character limit for restaurant names.
func WithNameCharLimit(n int) Option {
	return func(c *Config) {
		c.NameCharLimit = n
	}
}

// WithCuisineCharLimit sets the character limit for cuisine types.
func WithCuisineCharLimit(n int) Option {
	return func(c *Config) {
		c.CuisineCharLimit = n
	}
}

// WithInputWidth sets the width of input fields.
func WithInputWidth(w int) Option {
	return func(c *Config) {
		c.InputWidth = w
	}
}

// WithInitialElo sets the starting Elo rating for new restaurants.
func WithInitialElo(elo float64) Option {
	return func(c *Config) {
		c.InitialElo = elo
	}
}

// WithInitialScore sets the starting score for new restaurants.
func WithInitialScore(score float64) Option {
	return func(c *Config) {
		c.InitialScore = score
	}
}

// WithViewportPadding sets the padding for the viewport.
func WithViewportPadding(p int) Option {
	return func(c *Config) {
		c.ViewportPadding = p
	}
}

// WithProgressBarWidth sets the width of progress bars.
func WithProgressBarWidth(w int) Option {
	return func(c *Config) {
		c.ProgressBarWidth = w
	}
}

// WithCardWidth sets the width of comparison cards.
func WithCardWidth(w int) Option {
	return func(c *Config) {
		c.CardWidth = w
	}
}

// WithNameTruncateWidth sets the truncation width for names in lists.
func WithNameTruncateWidth(w int) Option {
	return func(c *Config) {
		c.NameTruncateWidth = w
	}
}

// WithCuisineTruncateWidth sets the truncation width for cuisines in lists.
func WithCuisineTruncateWidth(w int) Option {
	return func(c *Config) {
		c.CuisineTruncateWidth = w
	}
}

// WithCompactWindowContext sets the context size for compact search view.
func WithCompactWindowContext(n int) Option {
	return func(c *Config) {
		c.CompactWindowContext = n
	}
}

// WithResultWindowContext sets the context size for result view.
func WithResultWindowContext(n int) Option {
	return func(c *Config) {
		c.ResultWindowContext = n
	}
}

// Apply applies functional options to the config.
func (c *Config) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}
