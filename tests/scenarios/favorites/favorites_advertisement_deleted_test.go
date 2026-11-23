package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestFavoritesAdvertisementDeleted(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	owner := fixtures.Login(t, "FAVORITES_OWNER_LOGIN", "FAVORITES_OWNER_PASSWORD")
	require.NotEqual(t, user.User.Id, owner.User.Id, "подготовка избранного: user.id покупателя и владельца объявления должны различаться")

	cases := []struct {
		name                string
		deleteAdvertisement bool
		expectedPresent     bool
	}{
		{"существующее объявление остаётся в избранном", false, true},
		{"удалённое объявление исчезает из избранного", true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fields := fixtures.Fields()
			ad := board.CreateWithoutCleanup(t, owner.Token, fields, fixtures.PhotoPath())
			deleted := false
			t.Cleanup(func() {
				if !deleted {
					board.CleanupAdvertisement(t, owner.Token, ad.ID)
				}
			})
			t.Cleanup(func() { board.CleanupFavorite(t, user.Token, ad.ID) })
			board.AddFavorite(t, user.Token, ad.ID, http.StatusCreated)
			assertFavorite(t, user.Token, fields["title"], ad.ID, true)
			if tc.deleteAdvertisement {
				board.DeleteAdvertisement(t, owner.Token, ad.ID)
				deleted = true
			}
			assertFavorite(t, user.Token, fields["title"], ad.ID, tc.expectedPresent)
		})
	}
}
