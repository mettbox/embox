package env

import (
	"errors"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func SanitizeEnv(envName string) (string, error) {
	if len(envName) == 0 {
		return "", errors.New("Environment Variable Name Should Not Empty")
	}

	retValue := strings.TrimSpace(os.Getenv(envName))
	if len(retValue) == 0 {
		return "", errors.New("Environment Variable '" + envName + "' Has an Empty Value")
	}

	return retValue, nil
}

func GetEnv(envName string, defaultValue string) string {
	envValue, err := SanitizeEnv(envName)
	if err != nil {
		return defaultValue
	}

	return envValue
}

func GetEnvAsInt(envName string, defaultValue int) int {
	envValue, err := SanitizeEnv(envName)
	if err != nil {
		return defaultValue
	}

	retValue, err := strconv.Atoi(envValue)
	if err != nil {
		return defaultValue
	}

	return retValue
}

func GetEnvAsBool(envName string, defaultValue bool) bool {
	envValue, err := SanitizeEnv(envName)
	if err != nil {
		return defaultValue
	}

	retValue, err := strconv.ParseBool(envValue)
	if err != nil {
		return defaultValue
	}

	return retValue
}

func GetEnvSlice(key string, defaultVal []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return strings.Split(val, ",")
}
