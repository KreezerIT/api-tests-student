package board

import (
	"net/http"
	"testing"
	"time"

	client "api-tests-template/internal/client/http/board"

	"github.com/stretchr/testify/require"
)

// CleanupAdvertisement удаляет тестовое объявление
func CleanupAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	cleanup(t, token, id, "объявления", false, client.DeleteAdvertisement)
}

// CleanupFavorite удаляет тестовую запись избранного
func CleanupFavorite(t *testing.T, token, id string) {
	t.Helper()
	cleanup(t, token, id, "избранного", true, client.DeleteFavorite)
}

// cleanup повторяет удаление при временной недоступности сервиса
func cleanup(t *testing.T, token, id, kind string, allowNotFound bool, remove func(*testing.T, string, string) *http.Response) {
	t.Helper()
	const attempts = 5
	for attempt := 1; attempt <= attempts; attempt++ {
		response := remove(t, token, id)
		require.NotNil(t, response, "при очистке %s %s сервер не вернул ответ", kind, id)
		if allowNotFound && response.StatusCode == http.StatusNotFound {
			response.Body.Close()
			return
		}
		if response.StatusCode != http.StatusServiceUnavailable || attempt == attempts {
			readResponse(t, response, http.StatusNoContent, nil)
			return
		}
		response.Body.Close()
		t.Logf("Очистка %s %s: HTTP 503, повтор %d из %d", kind, id, attempt, attempts)
		time.Sleep(time.Second)
	}
}
