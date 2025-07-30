package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the configuration for the application.
type Config struct {
	Port        string `mapstructure:"PORT"`
	Static_Path string `mapstructure:"STATIC_PATH"`
	Html_Path   string `mapstructure:"HTML_PATH"`
	DB_Path     string `mapstructure:"DB_PATH"`
	//JWTSecret string `mapstructure:"JWT_SECRET"`
}

func Load() *Config {

	// Set default values
	viper.SetDefault("PORT", ":33031")

	// Set default paths relative to the working directory
	viper.SetDefault("STATIC_PATH", filepath.Join("..", "frontend", "static"))
	viper.SetDefault("HTML_PATH", filepath.Join("..", "frontend", "templates"))
	viper.SetDefault("DB_PATH", filepath.Join("..", "data", "grom.db"))
	//viper.SetDefault("JWT_SECRET", "your-secret-key")

	//viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env") // Load environment variables from .env file
	viper.AddConfigPath(".")   // look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Go: No config file found, using defaults and environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Go: Unable to decode config:", err)
	}

	return &config
}
