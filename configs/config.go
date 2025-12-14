package configs

import (
	"fmt"
	"os"
)

type conf struct {
	META_ACCESS_TOKEN string `mapstructure:"META_ACCESS_TOKEN"`
	NGROK_AUTH        string `mapstructure:"NGROK_AUTH"`
	PRIVATE_KEY_PATH  string `mapstructure:"PRIVATE_KEY_PATH"`
	PASSPHRASE        string `mapstructure:"PASSPHRASE"`
	// DBDriver          string `mapstructure:"DB_DRIVER"`
	// DBHost            string `mapstructure:"DB_HOST"`
	// DBPort            string `mapstructure:"DB_PORT"`
	// DBUser            string `mapstructure:"DB_USER"`
	// DBPassword        string `mapstructure:"DB_PASSWORD"`
	// DBName            string `mapstructure:"DB_NAME"`
	// WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	// JWTSecret         string `mapstructure:"JWT_SECRET"`
	// JWTExpiresIn      int    `mapstructure:"JWT_EXPIRESIN"`
	// TokenAuth         *jwtauth.JWTAuth
}

func LoadConfig() (*conf, error) {
	config := &conf{
		META_ACCESS_TOKEN: os.Getenv("META_ACCESS_TOKEN"),
		NGROK_AUTH:        os.Getenv("NGROK_AUTH"),
		PRIVATE_KEY_PATH:  os.Getenv("PRIVATE_KEY_PATH"),
		PASSPHRASE:        os.Getenv("PASSPHRASE"),
	}

	required := map[string]string{
		"META_ACCESS_TOKEN": config.META_ACCESS_TOKEN,
		"NGROK_AUTH":        config.NGROK_AUTH,
		"PRIVATE_KEY_PATH":  config.PRIVATE_KEY_PATH,
		"PASSPHRASE":        config.PASSPHRASE,
	}

	for name, value := range required {
		if value == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", name)
		}
	}

	return config, nil
}
