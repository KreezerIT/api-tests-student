package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestFavoriteInvalidAdvertisementID(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	cases := []struct{ name, id, message string }{
		{"empty", "", "invalid UUID length: 0"},
		{"invalid", "not-a-uuid", "invalid UUID length: 10"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Проверяем некорректный ID объявления")
			apiError := board.AddFavorite(t, user.Token, tc.id, http.StatusBadRequest)
			require.Equal(t, board.APIError{Error: "invalid_request", Message: tc.message}, apiError,
				"добавление в избранное с advertisement_id=%q: поля error и message должны описывать неверный формат ID", tc.id)
		})
	}
}
