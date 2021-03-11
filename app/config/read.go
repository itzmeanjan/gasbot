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

// GetGaszQueryURL - Returns full url of `gasz - Ethereum Gas Price Notifier`
// where GET request can be sent, for getting latest gas price recommendation
func GetGaszQueryURL() string {

	if url := Get("GASZ"); url != "" {
		return url
	}

	return "https://gasz.in/v1/latest"

}