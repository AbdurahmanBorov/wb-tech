# Сервис предоставляющий информацию о заказе

Демонстрационный сервис с простейшим интерфейсом для получения данных по некоторому uid-заказа.

## Требования ⚙️
Для запуска этого сервиса вас потребуется следующее:
- Docker
- Создать в корневой директории файл .env и заполнить по шаблону:
```
#POSTGRES ENVIRONMENTS
PG_USER=user
PG_PASSWORD=password
PGHOST=postgres
PGPORT=5432
PGDATABASE=db
PGSSLMODE=disable

HTTP_PORT=8080

NATSCLUSTERID=test-cluster
NATSCLIENTID=test-client
NATSCHANNEL=orders
NATSURL=http://nats-streaming:4222/
```

#### Миграции для базы данных находятся в директории [/migration](./migrations)

## Запуск 🔧

Для запуска пропишите команду ```make run```



## Интерфейс 🌐
После успешного запуска и перехода по пути ```http://localhost:8080/delivery``` (если указан тот же порт, что и в шаблоне)
у вас откроется страница, где вы можете ввести номер заказа и получить информацию по заказу.

## Публикация в канал nats-streaming

Есть скрипт publisher, после запуска считывает json из файла model.json и загружает в базу. Для запуска нужен запущенный сервис и прописать команду ```go run cmd/publisher/main.go```