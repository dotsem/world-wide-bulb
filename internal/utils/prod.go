package utils

func IsProd() bool {
	return GetEnv("APP_ENV", "development") == "production"
}
