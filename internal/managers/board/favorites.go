package board

import (
	"net/http"
	"testing"

	client "api-tests-template/internal/client/http/board"
)

// GetFavorites возвращает избранное с фильтром по названию
func GetFavorites(t *testing.T, token, title string) FavoriteList {
	t.Helper()
	var favorites FavoriteList
	readResponse(t, client.Favorites(t, token, title), http.StatusOK, &favorites)
	return favorites
}

// AddFavorite добавляет объявление в избранное и проверяет ответ
func AddFavorite(t *testing.T, token, id string, status int) APIError {
	t.Helper()
	response := client.AddFavorite(t, token, id)
	if status == http.StatusCreated {
		readResponse(t, response, status, nil)
		return APIError{}
	}
	return assertError(t, response, status)
}

// DeleteFavorite удаляет объявление из избранного и проверяет ответ
func DeleteFavorite(t *testing.T, token, id string, status int) APIError {
	t.Helper()
	response := client.DeleteFavorite(t, token, id)
	if status == http.StatusNoContent {
		readResponse(t, response, status, nil)
		return APIError{}
	}
	return assertError(t, response, status)
}

// FavoritesError проверяет ошибку при запросе избранного
func FavoritesError(t *testing.T, token string, status int) APIError {
	t.Helper()
	return assertError(t, client.Favorites(t, token, ""), status)
}
