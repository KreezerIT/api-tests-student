package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	base "api-tests-template/tests"

	"github.com/stretchr/testify/require"
)

func TestFavoritesUnauthorized(t *testing.T) {
	base.SetupSuite()
	for _, tc := range []struct {
		name    string
		request func(t *testing.T) board.APIError
	}{
		{"list", func(t *testing.T) board.APIError { return board.FavoritesError(t, "", http.StatusUnauthorized) }},
		{"add", func(t *testing.T) board.APIError { return board.AddFavorite(t, "", missingID, http.StatusUnauthorized) }},
		{"delete", func(t *testing.T) board.APIError {
			return board.DeleteFavorite(t, "", missingID, http.StatusUnauthorized)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			apiError := tc.request(t)
			require.Equal(t, board.APIError{Error: "unauthorized", Message: "Authorization header required"}, apiError)
		})
	}
}
