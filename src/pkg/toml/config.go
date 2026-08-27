package toml

import "C"
import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml"
)

type Config struct {
	DBConnectionString string `toml:"db_connection_string"` // DB connection string (psql)
	DBDriver           string `toml:"db_driver"`            // DB driver
	// Domain   string `toml:"domain"`   // Domain name, https://site.com
	// Hostname string `toml:"hostname"` // Host name, site.com
	// SessionKey   string `toml:"session_key"` // User session cookie key
	// CSRFKey      string `toml:"csrf_key"`    // CSRF cookie key
	// TZLocation string `toml:"tz_location"` // Timezone location
}

// LoadConfig loads config from file
func LoadConfig(configFilename string) (*Config, error) {
	var Conf Config
	f, err := os.Open(configFilename)
	if err != nil {
		return nil, err
	}
	decoder := toml.NewDecoder(f).Strict(true)
	if decoder == nil {
		return nil, fmt.Errorf("couldn't create decoder")
	}
	if err := decoder.Decode(&Conf); err != nil {
		return nil, err

	}
	return &Conf, nil
}
