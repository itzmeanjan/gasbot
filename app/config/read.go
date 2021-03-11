package config

import "github.com/spf13/viper"

// Read - Reading .env file content
// @note Supposed to be invoked, during application start up
func Read(file string) error {
	viper.SetConfigFile(file)

	return viper.ReadInConfig()
}

// Get - Get config value ( as string ) by key
func Get(key string) string {
	return viper.GetString(key)
}

// GetUint - Get config value ( as unsigned integer ) by key
func GetUint(key string) uint {
	return viper.GetUint(key)
}
