# План разработки Go SDK для чат-ботов eXpress (BotX)

## 1. Анализ требований и источников
- Зафиксировать домены BotX API и их версии/пути: Bots, Notifications, Events, Chats, Users, Files, SmartApps, Stickers, OpenID, Metrics (и др.).
- Выделить поведенческие особенности BotX: асинхронные ответы с `sync_id`, callback-эндпоинт, sync SmartApp events, статус-эндпоинт.
- Снять ключевые паттерны из `pybotx`:
  - `Bot` как оркестратор API + runtime обработчиков.
  - `HandlerCollector` для команд/системных событий/SmartApp sync.
  - Middleware-цепочки и exception middleware.
  - CallbackManager (ожидание/таймауты) и CallbackRepo (in-memory интерфейс).
  - JWT-верификация входящих запросов по `aud/iss` и секрету бота.
  - Подход к optional-полям (`Missing`/`Undefined`) и сериализации.

## 2. Архитектура SDK и публичный API (Go-идиомы)
- Определить пакеты:
  - `botx` (публичный фасад),
  - `client` (HTTP/транспорт + доменные API),
  - `bot` (runtime обработчиков и вебхук-эндпоинты),
  - `models` (DTO/доменные типы),
  - `middleware`, `errors`, `auth`, `files`, `callbacks`.
- Зафиксировать API пользователя:
  - `Bot` + `HandlerCollector` (регистрация команд/системных событий/SmartApp sync),
  - функции-хелперы для status/callback/command endpoints (net/http-friendly),
  - унифицированные методы API (`SendMessage`, `EditMessage`, `CreateChat` и т.п.).
- Заложить Go-стиль: контексты, опциональные значения через указатели или `optional` типы, явные ошибки, компактные интерфейсы.

## 3. Базовый транспорт, авторизация, сериализация
- Реализовать базовый HTTP клиент:
  - таймауты, повторные попытки, логирование,
  - базовый `Do(ctx, req)` с хуками.
- Реализовать `AuthorizedMethod`:
  - Bearer token, ленивое получение/кэш,
  - единая обработка HTTP-статусов.
- Реализовать модель ошибок:
  - map status -> error,
  - обработка callback ошибок (reason -> error).
- Сериализация с поддержкой `Missing` полей (omitempty + explicit “absent”):
  - тип `Optional[T]` с флагом `Set`.

## 4. Callback-менеджер и асинхронные вызовы
- CallbackManager:
  - хранение pending callbacks,
  - ожидание по `sync_id`, таймауты,
  - интерфейс репозитория (in-memory + возможность внешнего).
- Единый API-метод для “wait callback or fire-and-forget”.

## 5. Runtime бот-процессора
- `HandlerCollector`:
  - команды (visible/hidden, описание, валидация `/cmd`),
  - `default` handler,
  - system events handlers,
  - sync SmartApp handler.
- Middleware цепочки (pre/post), exception middleware.
- Контекстные данные (bot_id/chat_id) через `context.Context`.
- JWT-валидация входящих запросов:
  - проверка `aud`, `iss`, HS256, `trusted_issuers`.
- Эндпоинты:
  - `/command` (async bot commands),
  - `/smartapps/request` (sync event),
  - `/status`,
  - `/notification/callback`.

## 6. Домены BotX API (v1 охват)
- Реализовать методы по доменам (Go-пакеты + типы запросов/ответов):
  - Bots, Notifications, Events, Chats, Users, Files, SmartApps, Stickers, OpenID, Metrics.
- Особые кейсы из `pybotx`:
  - SmartApps: manifest, events, notification, unread counter, custom notification, list.
  - Files/attachments: upload, download, async files.
  - Stickers: upload/download + валидации PNG/размеров.

## 7. Модели и конвертеры доменных типов
- DTO для входящих сообщений, системных событий, SmartApp событий.
- Модель статус-меню и команд (visible/hidden).
- Конвертеры “API -> domain” и “domain -> API” (как в `pybotx.models.*`).
- Поля `sync_id`, `bot_id`, `chat_id`, типы вложений и async-files.

## 8. Тестирование
- Табличные тесты сериализации/десериализации.
- Тесты runtime (Handlers, middleware, sync SmartApp, callback flow).
- Контрактные тесты API-клиентов (stub HTTP server).

## 9. Документация и примеры
- `README` с минимальным ботом на net/http.
- Примеры:
  - echo-bot,
  - smartapp-sync,
  - callback listener.
- Таблица соответствия методов BotX API и SDK.

## 10. Релиз и поддержка
- SemVer, changelog.
- План расширения на “менее используемые” домены и edge cases.

