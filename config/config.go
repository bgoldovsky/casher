package config

import "os"

// Port Получает порт для запуска приложения
// Или подставляет значение по умолчанию, если он не указан
func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// ConnectionString Получает строку подключения к БД
// Или подставляет значение по умолчанию, если она не указана
func ConnectionString() string {
	cs := os.Getenv("DATABASE_URL")
	if cs == "" {
		cs = "dbname=casher sslmode=disable"
	}
	return cs
}
