package config

import "os"

func Get(key string) (value string) {
	return os.Getenv(key)
}
