package config

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`

	Twitch struct {
		ClientID      string `yaml:"client_id"`
		ClientSecret  string `yaml:"client_secret"`
		BroadcasterID string `yaml:"broadcaster_id"`
	} `yaml:"twitch"`

	OBS struct {
		WSAddress string `yaml:"ws_address"`
	} `yaml:"obs"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Set default values
	cfg.Server.Address = ":8080"
	cfg.OBS.WSAddress = ":8081"

	// Here you would load from config file
	// or environment variables

	return cfg, nil
}
