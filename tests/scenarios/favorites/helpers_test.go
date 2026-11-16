package favorites

import (
	"testing"

	"api-tests-template/internal/managers/board"

	"github.com/stretchr/testify/require"
)

const missingID = "00000000-0000-4000-8000-000000000000"

// assertFavorite проверяет, что объявление встречается в избранном ровно один раз или отсутствует
func assertFavorite(t *testing.T, token, title, id string, present bool) {
	t.Helper()
	list := board.GetFavorites(t, token, title)
	count := 0
	for _, item := range list.Items {
		if item.Advertisement.ID == id {
			count++
		}
	}
	expectedCount := 0
	if present {
		expectedCount = 1
	}
	require.Equal(t, expectedCount, count, "неверное количество записей объявления %s в избранном", id)
}
