package advertisement

import (
	"testing"

	"api-tests-template/internal/managers/board"
	"api-tests-template/tests/fixtures"

	"github.com/stretchr/testify/require"
)

type expectedAdvertisement struct {
	ID, UserID, Title, Description string
	Price                          float64
	Quantity                       int
}

func TestCreateAdvertisement(t *testing.T) {
	user := fixtures.Login(t, "TEST_LOGIN", "TEST_PASSWORD")
	cases := []struct {
		name, price, quantity string
		expectedPrice         float64
		expectedQuantity      int
	}{
		{"со всеми параметрами", "1234.50", "3", 1234.5, 3},
		{"нулевая цена и одна единица", "0", "1", 0, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fields := fixtures.Fields()
			fields["price"], fields["quantity"] = tc.price, tc.quantity
			created := board.Create(t, user.Token, fields, fixtures.PhotoPath())

			var ad board.Advertisement
			require.True(t, t.Run("объявление доступно по ID и содержит переданные параметры", func(t *testing.T) {
				ad = board.GetAdvertisement(t, created.ID)
				require.Equal(t, expectedAdvertisement{
					ID: created.ID, UserID: user.User.Id, Title: fields["title"],
					Description: fields["description"], Price: tc.expectedPrice, Quantity: tc.expectedQuantity,
				}, expectedAdvertisement{
					ID: ad.ID, UserID: ad.UserID, Title: ad.Title,
					Description: ad.Description, Price: ad.Price, Quantity: ad.Quantity,
				})
			}))

			t.Run("объявление находится по полному названию", func(t *testing.T) {
				results := board.SearchAdvertisements(t, fields["title"])
				ids := make([]string, 0, len(results.Items))
				for _, item := range results.Items {
					ids = append(ids, item.ID)
				}
				require.Contains(t, ids, created.ID, "созданное объявление отсутствует в поиске")
			})

			t.Run("загруженная фотография доступна по API", func(t *testing.T) {
				photos := board.GetPhotos(t, user.Token, created.ID)
				require.Len(t, photos, 1)
				require.NotEmpty(t, photos[0].ID)
				require.NotEmpty(t, photos[0].URL)
				require.Len(t, ad.Photos, 1)
				require.Equal(t, ad.Photos[0].ID, photos[0].ID)
			})
		})
	}
}
