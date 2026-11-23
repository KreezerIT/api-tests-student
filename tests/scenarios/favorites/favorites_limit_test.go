package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

const maxFavorites = 50

func TestFavoritesLimit(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	owner := fixtures.Login(t, "FAVORITES_OWNER_LOGIN", "FAVORITES_OWNER_PASSWORD")
	require.NotEqual(t, user.User.Id, owner.User.Id, "подготовка избранного: user.id покупателя и владельца объявления должны различаться")

	current := board.GetFavorites(t, user.Token, "")
	require.NotNil(t, current.Total, "подготовка проверки лимита: в исходном списке избранного отсутствует поле total")
	require.LessOrEqual(t, *current.Total, maxFavorites, "подготовка проверки лимита: исходное поле total превышает допустимые 50 записей")
	initialCount := *current.Total
	needed := maxFavorites - 1 - initialCount
	t.Logf("Подготовка: в избранном %d объявлений, добавляем %d для проверки границы", initialCount, max(needed, 0))
	for i := 0; i < needed; i++ {
		ad := board.Create(t, owner.Token, fixtures.Fields(), fixtures.PhotoPath())
		t.Cleanup(func() { board.CleanupFavorite(t, user.Token, ad.ID) })
		board.AddFavorite(t, user.Token, ad.ID, http.StatusCreated)
	}

	cases := []struct {
		name            string
		fillLastSlot    bool
		status          int
		expectedError   board.APIError
		expectedPresent bool
	}{
		{name: "добавление пятидесятого объявления", status: http.StatusCreated, expectedPresent: true},
		{name: "добавление сверх лимита", fillLastSlot: true, status: http.StatusBadRequest,
			expectedError: board.APIError{Error: "limit_reached", Message: "Favorites limit of 50 reached"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(caseT *testing.T) {
			if !tc.fillLastSlot && initialCount == maxFavorites {
				caseT.Skip("избранное уже заполнено до начала теста")
			}
			if tc.fillLastSlot && initialCount < maxFavorites {
				last := board.Create(caseT, owner.Token, fixtures.Fields(), fixtures.PhotoPath())
				caseT.Cleanup(func() { board.CleanupFavorite(caseT, user.Token, last.ID) })
				board.AddFavorite(caseT, user.Token, last.ID, http.StatusCreated)
			}

			before := board.GetFavorites(caseT, user.Token, "")
			require.NotNil(caseT, before.Total, "перед добавлением на границе лимита: в списке избранного отсутствует поле total")
			expectedBefore := maxFavorites
			if !tc.fillLastSlot {
				expectedBefore--
			}
			require.Equal(caseT, expectedBefore, *before.Total, "перед добавлением на границе лимита: поле total не соответствует подготовленному количеству записей")

			fields := fixtures.Fields()
			ad := board.Create(caseT, owner.Token, fields, fixtures.PhotoPath())
			caseT.Cleanup(func() { board.CleanupFavorite(caseT, user.Token, ad.ID) })
			actual := board.AddFavorite(caseT, user.Token, ad.ID, tc.status)
			require.Equal(caseT, tc.expectedError, actual,
				"добавление объявления %s на границе лимита (%s): поля error и message не соответствуют ожидаемому результату", ad.ID, tc.name)

			after := board.GetFavorites(caseT, user.Token, "")
			require.NotNil(caseT, after.Total, "после добавления на границе лимита: в списке избранного отсутствует поле total")
			require.Equal(caseT, maxFavorites, *after.Total, "после добавления на границе лимита: поле total должно оставаться равным 50")
			assertFavorite(caseT, user.Token, fields["title"], ad.ID, tc.expectedPresent)
		})
	}
}
