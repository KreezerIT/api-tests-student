package board

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// readResponse проверяет статус и разбирает JSON-ответ
func readResponse(t *testing.T, response *http.Response, status int, out any) {
	t.Helper()
	require.NotNil(t, response, "сервер не вернул HTTP-ответ")
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать ответ")
	require.Equalf(t, status, response.StatusCode, "неожиданный HTTP-статус; ответ: %s", body)
	if out != nil {
		require.NoErrorf(t, json.Unmarshal(body, out), "некорректный JSON: %s", body)
	}
}

// assertError проверяет статус и поля ошибки API
func assertError(t *testing.T, response *http.Response, status int) APIError {
	t.Helper()
	var apiError APIError
	readResponse(t, response, status, &apiError)
	require.NotEmpty(t, apiError.Error, "в ответе отсутствует код ошибки")
	require.NotEmpty(t, apiError.Message, "в ответе отсутствует сообщение об ошибке")
	t.Logf("Ошибка API: %s: %s", apiError.Error, apiError.Message)
	return apiError
}

// assertCreateError проверяет ошибку создания и удаляет неожиданно созданное объявление
func assertCreateError(t *testing.T, response *http.Response, ownerToken string, status int) APIError {
	t.Helper()
	if response != nil && response.StatusCode == http.StatusCreated {
		var ad Advertisement
		readResponse(t, response, http.StatusCreated, &ad)
		require.NotEmpty(t, ad.ID, "ошибочно созданное объявление не содержит ID")
		t.Cleanup(func() { CleanupAdvertisement(t, ownerToken, ad.ID) })
		t.Fatalf("ошибочно создано объявление %s вместо отказа с HTTP %d", ad.ID, status)
	}
	return assertError(t, response, status)
}
