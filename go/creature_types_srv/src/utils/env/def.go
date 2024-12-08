package env

import "os"

func Get(env_name, def_value string) *string {
	value := os.Getenv(env_name)
	if len(value) == 0 {
		value = def_value
	}
	return &value
}
