package myAdvertisement

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/board"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/tests/fixtures"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Проверяем получение собственных объявлений в профиле
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	advertisementID := board.Create(s.T(), s.loginData.Token, fixtures.Fields(), fixtures.PhotoPath()).ID
	var advertisementsBody string
	if !s.Run("Получаем список собственных объявлений", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, 50)
	}) {
		return
	}

	var items gjson.Result
	if !s.Run("Проверяем, что список объявлений не пуст", func() {
		items = gjson.Get(advertisementsBody, "items")
		require.NotEmpty(s.T(), items.Array())
	}) {
		return
	}

	s.Run("Проверяем владельца и поля созданного объявления", func() {
		found := false
		for _, item := range items.Array() {
			if gjson.Get(item.String(), "id").String() != advertisementID {
				continue
			}
			found = true
			require.Equal(s.T(), s.loginData.User.Id, gjson.Get(item.String(), "user_id").String(), "user_id не совпадает")

			require.NotEmpty(s.T(), gjson.Get(item.String(), "id").String(), "id пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "created_at").String(), "created_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "updated_at").String(), "updated_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "title").String(), "Название пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "description").String(), "Описание пустое")
			require.True(s.T(), len(gjson.Get(item.String(), "photos").Array()) > 0, "Присутствуют пустые фотографии")
		}
		require.True(s.T(), found, "Созданное объявление не найдено в списке")
	})
}

func (s *TestSuite) TestGetMyAdvertisementsUnauthorized() {
	cases := []struct {
		name, token, message string
	}{
		{"без токена", "", "Authorization header required"},
		{"неверный токен", "incorrect_token", "Invalid or expired token"},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			body := myAdvertisements.GetMyAdvertisements(s.T(), tc.token, http.StatusUnauthorized, 50)
			actual := struct{ Error, Message string }{
				Error: gjson.Get(body, "error").String(), Message: gjson.Get(body, "message").String(),
			}
			require.Equal(s.T(), struct{ Error, Message string }{Error: "unauthorized", Message: tc.message}, actual)
		})
	}
}
