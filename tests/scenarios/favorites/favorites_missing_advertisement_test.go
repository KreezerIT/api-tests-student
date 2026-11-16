package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestFavoritesMissingAdvertisement(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	cases := []struct {
		name, message string
		request       func(*testing.T) board.APIError
	}{
		{"добавление", "Advertisement not found", func(t *testing.T) board.APIError {
			return board.AddFavorite(t, user.Token, missingID, http.StatusNotFound)
		}},
		{"удаление", "Favorite entry not found", func(t *testing.T) board.APIError {
			return board.DeleteFavorite(t, user.Token, missingID, http.StatusNotFound)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, board.APIError{Error: "not_found", Message: tc.message}, tc.request(t))
		})
	}
}
