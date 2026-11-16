В [ПР'е](https://github.com/KreezerIT/api-tests-student/pull/1) описал TL;DR решения и приложил кусочек логов

---

# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
FAVORITES_OWNER_LOGIN=second_user_login
FAVORITES_OWNER_PASSWORD=second_user_password
```

Для сценариев избранного нужны два разных аккаунта TestBoard: второй создаёт объявления, которые первый добавляет в избранное. Тесты создают объявления с уникальными названиями и удаляют их после проверки. Не добавляйте `.env.override` в Git.

8. Для запуска всех тестов достаточно выполнить команду в корне проекта:

```bash
  go test -p 1 -v ./...
```

9. Запустите все тесты три раза без кэша, чтобы проверить повторяемость:

```bash
  go test -p 1 -count=3 ./...
```

10. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -v ./tests/scenarios/myAdvertisement/
```

HTTP runner выдерживает минимум 300 мс между запросами, включая авторизацию и очистку тестовых данных. Это снижает нагрузку на TestBoard, который при частых запросах возвращает 503. Ограничитель общий для клиентов одного процесса, поэтому запускайте пакеты с `-p 1` и не запускайте несколько прогонов одновременно.

## Логирование

Для включения логирования request/response в консоль выставьте `DEBUG=true` в `.env.override`. Логи содержат логины, пароли и токены: не публикуйте их. После отладки верните `DEBUG=false`.

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```
