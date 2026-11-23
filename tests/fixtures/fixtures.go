package fixtures

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"api-tests-template/internal/managers/auth"
	authmodels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"

	"github.com/stretchr/testify/require"
)

// Login загружает настройки и авторизует тестового пользователя
func Login(t *testing.T, loginKey, passwordKey string) authmodels.LoginOkResponse {
	t.Helper()
	base.SetupSuite()
	result := auth.Login(t, os.Getenv(loginKey), os.Getenv(passwordKey))
	require.NotEmpty(t, result.Token, "подготовка пользователя из %s: в ответе авторизации отсутствует token", loginKey)
	return result
}

// Fields возвращает данные для объявления со случайным названием
func Fields() map[string]string {
	return map[string]string{
		"title":       "QA lab " + utils.RandomString(20),
		"description": "Тестовое объявление лабораторной работы",
		"price":       "1234.50",
		"quantity":    "3",
	}
}

// PhotoPath возвращает путь к тестовой фотографии
func PhotoPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "scenarios", "advertisement", "images", "Валерон.png")
}
