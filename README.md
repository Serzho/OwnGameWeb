# OwnGameWeb
[![Go Tests and Linter](https://github.com/Serzho/OwnGameWeb/actions/workflows/go.yml/badge.svg)](https://github.com/Serzho/OwnGameWeb/actions/workflows/go.yml)  [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## Стек:

**Gin, PostgreSQL, vue.js, Docker, slog, jwt**

## Описание:
"Своя игра" — популярная телевизионная интеллектуальная викторина, где игроки выбирают вопросы из разных категорий. Веб реализация позволяет подключать пакеты вопросы, создавать комнаты и проводить игры.

## Запуск:

1. С помощью Makefile:
```
make build
```
2. Docker-compose:
```
go mod vendor
docker-compose up --build
```
3. Go build:
```
go mod download
go run ./cmd/app/
```

## Особенности реализации:

- Авторизация с помощью JWT-токенов
- Конфигурирование через файл .env в корне проекта
- Github CI: запуск линтера golangcilint (конфиг в файле .golangci.yml) и тестов
- Логирование в JSON формате
- Статические ошибки для каждого модуля
- Тестирование endpoints с помощью mock объектов
- Чистая архитектура с разделением на слои
- Vue.js интегрирован с помощью SDN
