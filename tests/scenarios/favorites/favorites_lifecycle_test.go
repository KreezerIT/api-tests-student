package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestFavoritesLifecycle(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	owner := fixtures.Login(t, "FAVORITES_OWNER_LOGIN", "FAVORITES_OWNER_PASSWORD")
	require.NotEqual(t, user.User.Id, owner.User.Id, "для избранного нужен другой владелец")

	cases := []struct {
		name            string
		remove          bool
		expectedPresent bool
	}{
		{"добавленное объявление появляется в избранном", false, true},
		{"удалённое объявление исчезает из избранного", true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fields := fixtures.Fields()
			ad := board.Create(t, owner.Token, fields, fixtures.PhotoPath())
			t.Cleanup(func() { board.CleanupFavorite(t, user.Token, ad.ID) })
			board.AddFavorite(t, user.Token, ad.ID, http.StatusCreated)
			if tc.remove {
				assertFavorite(t, user.Token, fields["title"], ad.ID, true)
				board.DeleteFavorite(t, user.Token, ad.ID, http.StatusNoContent)
			}
			assertFavorite(t, user.Token, fields["title"], ad.ID, tc.expectedPresent)
		})
	}
}
