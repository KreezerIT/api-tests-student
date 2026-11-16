package board

import (
	"net/http"
	"testing"

	client "api-tests-template/internal/client/http/board"

	"github.com/stretchr/testify/require"
)

// Create создаёт объявление и регистрирует его удаление после теста
func Create(t *testing.T, token string, fields map[string]string, photo string) Advertisement {
	t.Helper()
	ad := CreateWithoutCleanup(t, token, fields, photo)
	t.Cleanup(func() {
		t.Logf("Очистка объявления %s", ad.ID)
		CleanupAdvertisement(t, token, ad.ID)
	})
	return ad
}

// GetAdvertisement возвращает объявление по ID
func GetAdvertisement(t *testing.T, id string) Advertisement {
	t.Helper()
	var ad Advertisement
	readResponse(t, client.Advertisement(t, id), http.StatusOK, &ad)
	return ad
}

// SearchAdvertisements возвращает результаты поиска по названию
func SearchAdvertisements(t *testing.T, title string) Advertisements {
	t.Helper()
	var ads Advertisements
	readResponse(t, client.Search(t, title), http.StatusOK, &ads)
	return ads
}

// GetPhotos возвращает фотографии объявления
func GetPhotos(t *testing.T, token, id string) []Photo {
	t.Helper()
	var photos []Photo
	readResponse(t, client.Photos(t, token, id), http.StatusOK, &photos)
	return photos
}

// CreateError проверяет ошибку при создании объявления
func CreateError(t *testing.T, token string, fields map[string]string, photo string, status int) APIError {
	t.Helper()
	return assertCreateError(t, client.CreateAdvertisement(t, token, fields, photo), token, status)
}

// DeleteAdvertisement удаляет объявление и проверяет ответ
func DeleteAdvertisement(t *testing.T, token, id string) {
	t.Helper()
	readResponse(t, client.DeleteAdvertisement(t, token, id), http.StatusNoContent, nil)
}

// CreateWithoutCleanup создаёт объявление без регистрации очистки
func CreateWithoutCleanup(t *testing.T, token string, fields map[string]string, photo string) Advertisement {
	t.Helper()
	var ad Advertisement
	readResponse(t, client.CreateAdvertisement(t, token, fields, photo), http.StatusCreated, &ad)
	require.NotEmpty(t, ad.ID, "сервер не вернул ID объявления")
	return ad
}
