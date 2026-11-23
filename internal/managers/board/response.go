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
	require.NotNil(t, response, "получение HTTP-ответа: сервер не вернул ответ для проверки статуса %d", status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "чтение HTTP-ответа: не удалось прочитать тело ответа со статусом %d", response.StatusCode)
	require.Equalf(t, status, response.StatusCode, "проверка HTTP-ответа: неожиданный статус; тело ответа: %s", body)
	if out != nil {
		require.NoErrorf(t, json.Unmarshal(body, out), "разбор HTTP-ответа со статусом %d: тело должно быть корректным JSON; ответ: %s", status, body)
	}
}

// assertError проверяет статус и поля ошибки API
func assertError(t *testing.T, response *http.Response, status int) APIError {
	t.Helper()
	var apiError APIError
	readResponse(t, response, status, &apiError)
	require.NotEmpty(t, apiError.Error, "проверка ошибки API с HTTP %d: поле error должно содержать код ошибки", status)
	require.NotEmpty(t, apiError.Message, "проверка ошибки API с HTTP %d: поле message должно содержать описание ошибки", status)
	t.Logf("Ошибка API: %s: %s", apiError.Error, apiError.Message)
	return apiError
}

// assertCreateError проверяет ошибку создания и удаляет неожиданно созданное объявление
func assertCreateError(t *testing.T, response *http.Response, ownerToken string, status int) APIError {
	t.Helper()
	if response != nil && response.StatusCode == http.StatusCreated {
		var ad Advertisement
		readResponse(t, response, http.StatusCreated, &ad)
		require.NotEmpty(t, ad.ID, "подготовка очистки после ошибочного создания: в ответе POST /advertisement отсутствует id")
		t.Cleanup(func() { CleanupAdvertisement(t, ownerToken, ad.ID) })
		t.Fatalf("ошибочно создано объявление %s вместо отказа с HTTP %d", ad.ID, status)
	}
	return assertError(t, response, status)
}
