package advertisement

import (
	"net/http"
	"strings"
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestCreateAdvertisementInvalid(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	cases := []struct {
		name, field, value, message string
		noPhoto                     bool
	}{
		{name: "empty_title", field: "title", value: "", message: "title is required"},
		{name: "negative_price", field: "price", value: "-1", message: "price must be between 0 and 1000000 inclusive"},
		{name: "long_title", field: "title", value: strings.Repeat("a", 256), message: "title must not exceed 50 characters provided 256"},
		{name: "missing_photo", noPhoto: true, message: "at least 1 photo is required"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fields := fixtures.Fields()
			if tc.field != "" {
				fields[tc.field] = tc.value
			}
			photo := fixtures.PhotoPath()
			if tc.noPhoto {
				photo = ""
			}
			t.Log("Отправляем некорректные данные и проверяем ошибку")
			apiError := board.CreateError(t, user.Token, fields, photo, http.StatusBadRequest)
			require.Equal(t, board.APIError{Error: "invalid_request", Message: tc.message}, apiError)
		})
	}
}
