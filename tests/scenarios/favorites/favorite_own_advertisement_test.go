package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

func TestFavoriteOwnAdvertisement(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	cases := []struct {
		name             string
		previousAttempts int
	}{
		{"первая попытка", 0},
		{"повторная попытка после отказа", 1},
	}
	expected := board.APIError{Error: "forbidden", Message: "Cannot add own advertisement to favorites"}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fields := fixtures.Fields()
			ad := board.Create(t, user.Token, fields, fixtures.PhotoPath())
			t.Cleanup(func() { board.CleanupFavorite(t, user.Token, ad.ID) })
			for i := 0; i < tc.previousAttempts; i++ {
				require.Equal(t, expected, board.AddFavorite(t, user.Token, ad.ID, http.StatusForbidden),
					"добавление собственного объявления %s: поля error и message должны сообщать о запрете", ad.ID)
			}
			require.Equal(t, expected, board.AddFavorite(t, user.Token, ad.ID, http.StatusForbidden),
				"добавление собственного объявления %s: поля error и message должны сообщать о запрете", ad.ID)
			assertFavorite(t, user.Token, fields["title"], ad.ID, false)
		})
	}
}
