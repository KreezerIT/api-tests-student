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
	require.NotEqual(t, user.User.Id, owner.User.Id, "для избранного нужен другой владелец")

	current := board.GetFavorites(t, user.Token, "")
	require.NotNil(t, current.Total, "в списке избранного отсутствует total")
	require.LessOrEqual(t, *current.Total, maxFavorites, "количество избранных превышает лимит")
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
			require.NotNil(caseT, before.Total, "в списке избранного отсутствует total")
			expectedBefore := maxFavorites
			if !tc.fillLastSlot {
				expectedBefore--
			}
			require.Equal(caseT, expectedBefore, *before.Total, "неверное количество избранных перед запросом")

			fields := fixtures.Fields()
			ad := board.Create(caseT, owner.Token, fields, fixtures.PhotoPath())
			caseT.Cleanup(func() { board.CleanupFavorite(caseT, user.Token, ad.ID) })
			actual := board.AddFavorite(caseT, user.Token, ad.ID, tc.status)
			require.Equal(caseT, tc.expectedError, actual)

			after := board.GetFavorites(caseT, user.Token, "")
			require.NotNil(caseT, after.Total, "в списке избранного отсутствует total")
			require.Equal(caseT, maxFavorites, *after.Total, "неверное количество избранных после запроса")
			assertFavorite(caseT, user.Token, fields["title"], ad.ID, tc.expectedPresent)
		})
	}
}
