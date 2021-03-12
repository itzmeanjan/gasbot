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

// GetGaszSubscribeURL - Subscribe to latest gas price feed of `gasz`
func GetGaszSubscribeURL() string {

	if url := Get("GASZ_Subscribe"); url != "" {
		return url
	}

	return "wss://gasz.in/v1/subscribe"

}

// GetPort - Service to run on this port number
func GetPort() uint {

	if port := GetUint("Port"); port > 1024 {
		return port
	}

	return 7000

}

// GetToken - Token for interacting with Telegram HTTP API
func GetToken() string {
	return Get("Token")
}

// GetURL - Returns public URL of this service
// to which telegram will talk to
func GetURL() string {
	return Get("Url")
}
