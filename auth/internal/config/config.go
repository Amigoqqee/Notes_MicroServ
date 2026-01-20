package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                   string
	Host                   string
	Timeout                int
	DBDSN                  string
	DBSSL                  string
	DBTimeout              int
	JWTSecretKey           string
	AccessTokenExpiration  int
	RefreshTokenExpiration int
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s: %s", "не установлена переменная окружения", key)
	}
	return value, nil
}

func NewConfig() *Config {
	port, err := getEnv("PORT")
	if err != nil {
		fmt.Println("Не удалось получить PORT из переменной окружения, используется порт по умолчанию")
	}

	host, err := getEnv("HOST")
	if err != nil {
		fmt.Println("Не удалось получить HOST из переменной окружения, исользуется хост по умолчанию")
	}

	timeout := 10
	if envValue, err := getEnv("SERVER_TIMEOUT"); err == nil {
		if parsed, parseErr := strconv.Atoi(envValue); parseErr == nil {
			timeout = parsed
		}
	} else {
		fmt.Println("Не удалось получить SERVER_TIMEOUT из переменной окружения, используется 10 сек")
	}

	dbTimeout := 5
	if envValue, err := getEnv("DB_TIMEOUT"); err == nil {
		if parsed, parseErr := strconv.Atoi(envValue); parseErr == nil {
			dbTimeout = parsed
		}
	} else {
		fmt.Println("Не удалось получить DB_TIMEOUT из переменной окружения, используется 5 секунд")
	}

	dbHost, err := getEnv("POSTGRES_HOST")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_HOST из переменной окружения")
	}
	dbPort, err := getEnv("POSTGRES_PORT")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_PORT из переменной окружения")
	}
	dbUser, err := getEnv("POSTGRES_USER")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_USER из переменной окружения")
	}
	dbPassword, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_PASSWORD из переменной окружения")
	}
	dbName, err := getEnv("POSTGRES_DB")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_DB из переменной окружения")
	}
	dbSSL, err := getEnv("POSTGRES_USE_SSL")
	if err != nil {
		fmt.Println("Не удалось получить POSTGRES_USE_SSL из переменной окружения")
	}

	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSL)

	jwtSecretKey, err := getEnv("JWT_SECRET_KEY")
	if err != nil {
		fmt.Println("Не удалось получить JWT_SECRET_KEY из переменной окружения, используется значение по умолчанию")
	}

	accessTokenExpiration := 24
	if envValue, err := getEnv("JWT_ACCESS_TOKEN_EXPIRATION"); err == nil {
		if parsed, parseErr := strconv.Atoi(envValue); parseErr == nil {
			accessTokenExpiration = parsed
		}
	}

	refreshTokenExpiration := 168
	if envValue, err := getEnv("JWT_REFRESH_TOKEN_EXPIRATION"); err == nil {
		if parsed, parseErr := strconv.Atoi(envValue); parseErr == nil {
			refreshTokenExpiration = parsed
		}
	}

	return &Config{
		Port:                   port,
		Host:                   host,
		DBDSN:                  dbDSN,
		DBSSL:                  dbSSL,
		JWTSecretKey:           jwtSecretKey,
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
		Timeout:                timeout,
		DBTimeout:              dbTimeout,
	}
}
