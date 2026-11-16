package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestFavoriteDuplicate(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	owner := fixtures.Login(t, "FAVORITES_OWNER_LOGIN", "FAVORITES_OWNER_PASSWORD")
	require.NotEqual(t, user.User.Id, owner.User.Id, "для избранного нужен другой владелец")

	cases := []struct {
		name               string
		previousDuplicates int
		expected           board.APIError
	}{
		{"повторное добавление", 0, board.APIError{Error: "already_exists", Message: "Advertisement is already in favorites"}},
		{"добавление после отказа", 1, board.APIError{Error: "already_exists", Message: "Advertisement is already in favorites"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fields := fixtures.Fields()
			ad := board.Create(t, owner.Token, fields, fixtures.PhotoPath())
			t.Cleanup(func() { board.CleanupFavorite(t, user.Token, ad.ID) })
			board.AddFavorite(t, user.Token, ad.ID, http.StatusCreated)
			for i := 0; i < tc.previousDuplicates; i++ {
				actual := board.AddFavorite(t, user.Token, ad.ID, http.StatusConflict)
				require.Equal(t, tc.expected, actual)
			}
			actual := board.AddFavorite(t, user.Token, ad.ID, http.StatusConflict)
			require.Equal(t, tc.expected, actual)
			assertFavorite(t, user.Token, fields["title"], ad.ID, true)
		})
	}
}
