# Subscription Service

REST-сервис для управления онлайн-подписками пользователей.

## Функционал

* CRUD операции с подписками (создание, чтение, обновление, удаление)
* Подсчет общей стоимости подписок за период с фильтрацией по пользователю и названию сервиса
* Swagger документация API

## Технологии

* Go 1.24
* Gin (HTTP фреймворк)
* GORM (ORM для PostgreSQL)
* PostgreSQL
* Docker и Docker Compose
* Swagger (документация API)
* Логирование с использованием `logrus`

## Как запустить

1. Создать файл `.env` по примеру `.env.example` с параметрами БД.

2. Запустить сервис через Docker Compose:

```bash
docker-compose up --build
```

3. Доступные сервисы:

* API: `http://localhost:8080`
* Swagger UI: `http://localhost:8080/swagger/index.html`

## Конфигурация

Настройки БД и другие параметры хранятся в `.env` файле:

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=subscriptions
```

## API

Полная документация доступна через Swagger UI.

## Особенности

* Проверка существования пользователя не реализована (вне зоны ответственности)
* Цена подписки — целое число (рубли, без копеек)
* Формат даты: месяц-год (`MM-YYYY`)
