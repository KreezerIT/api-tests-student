package board

import (
	"net/http"
	"testing"

	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// CreateAdvertisement отправляет запрос на создание объявления с фотографией
func CreateAdvertisement(t *testing.T, token string, fields map[string]string, photo string) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().Post("/advertisement")
	for name, value := range fields {
		request.MultipartFormData(name, value)
	}
	if photo != "" {
		request.MultipartFile("photos", photo)
	}
	return request.Expect(t).End().Response
}

// Advertisement запрашивает объявление по ID
func Advertisement(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get("/advertisement").Query("id", id).Expect(t).End().Response
}

// DeleteAdvertisement отправляет запрос на удаление объявления
func DeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete("/advertisement").Query("id", id).Expect(t).End().Response
}

// Search ищет объявления по названию
func Search(t *testing.T, title string) *http.Response {
	return apiRunner.GetRunner().Create().Get("/advertisements").Query("search", title).Expect(t).End().Response
}

// Photos запрашивает фотографии объявления
func Photos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get("/advertisements/" + id + "/photos").Expect(t).End().Response
}

// Favorites запрашивает избранное с фильтром по названию
func Favorites(t *testing.T, token, title string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get("/favorites").Query("search", title).Expect(t).End().Response
}

// AddFavorite отправляет запрос на добавление в избранное
func AddFavorite(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post("/favorites").JSON(map[string]string{"advertisement_id": id}).Expect(t).End().Response
}

// DeleteFavorite отправляет запрос на удаление из избранного
func DeleteFavorite(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete("/favorites/" + id).Expect(t).End().Response
}
