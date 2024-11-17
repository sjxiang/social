package config

import (
	"os"
	"strconv"
)


func defaultEnvString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func defaultEnvNumeric(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInteger, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInteger
}


func defaultEnvBoolean(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}