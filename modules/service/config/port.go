package config

func GetPort() string {
	return GetEnv("PORT", "8080")
}
