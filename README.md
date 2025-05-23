# Github rep tracking bot
Telegram бот для отслеживания commit, issue, pull request в репозиториях на github.
## Функциональность:
 - Добавление, удаление и отображение отслеживаемых репозиториев.
 - Рассылка уведомлений об изменениях в репозиториях.
## Команды бота
 - `/help` - информация о доступных командах.
 - `/start` - запуск.
 - `/repos` - отслеживаемые репозитории.
 - `/add` - добавление репозитория в отслеживание.
 - `/del` - удаление репозитория из отслеживания.
 - `/cancel` - прекращение любой команды.
## Архитектура проекта
 - Проект состоит из 2 микросервисов - основного сервера, выполняющего CRUD операции для отслеживаемых репозиториев и отслеживающего изменения в них, и сервиса, осуществляющего взаимодействие с пользователем через Telegram API. 
 - Взаимодействие между сервисами для выполнения CRUD операций осуществляется по http, передача уведомлений об изменениях - через Kafka. 
 - Для обработки команд пользователя в Telegram реализована машина состояний.
## Технологический стек
 - Kafka.
 - Golang.
 - segmentio/kafka-go - драйвер для kafka.
#### Основной сервер
 - Gorilla/mux - маршрутизация.
 - PostgreSQL - хранение информации о пользователях и отслеживаемых репозиториях.
 - PGX - драйвер для PostgreSQL.
 - Gocron - планировщик задач.
#### Сервис для работы с Telegram API
 - go-telegram-bot-api - взаимодействия с пользователем через Telegram.
