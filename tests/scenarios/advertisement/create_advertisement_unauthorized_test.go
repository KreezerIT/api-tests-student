package advertisement

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	base "api-tests-template/tests"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestCreateAdvertisementUnauthorized(t *testing.T) {
	base.SetupSuite()
	cases := []struct{ name, token, message string }{
		{"без токена", "", "Authorization header required"},
		{"неверный токен", "invalid_token", "Invalid or expired token"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := board.CreateError(t, tc.token, fixtures.Fields(), fixtures.PhotoPath(), http.StatusUnauthorized)
			require.Equal(t, board.APIError{Error: "unauthorized", Message: tc.message}, actual,
				"создание объявления без авторизации (%s): поля error и message должны описывать ошибку токена", tc.name)
		})
	}
}
