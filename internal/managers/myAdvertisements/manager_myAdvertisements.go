package myAdvertisements

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/myAdvertisements"
)

// GetMyAdvertisements возвращает список объявлений авторизованного пользователя и проверяет ожидаемый статус-код
func GetMyAdvertisements(t *testing.T, token string, expectedStatusCode int, limit int) string {
	t.Helper()
	myAdvertisementsRawResponse := myAdvertisements.HttpGetMyAdvertisements(t, token, limit)
	defer myAdvertisementsRawResponse.Body.Close()

	require.Equalf(t, expectedStatusCode, myAdvertisementsRawResponse.StatusCode,
		"получение собственных объявлений: HTTP-статус GET /my/advertisements должен быть %d", expectedStatusCode)

	myAdvertisementsResponse, errResponse := io.ReadAll(myAdvertisementsRawResponse.Body)
	require.NoError(t, errResponse, "получение собственных объявлений: не удалось прочитать тело ответа GET /my/advertisements")

	return string(myAdvertisementsResponse)
}
