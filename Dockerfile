# Используем официальный образ Go
FROM golang:1.24.1-alpine

RUN apk add --no-cache make
# Указываем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы Go модуля в контейнер
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем все файлы исходного кода в контейнер
COPY . .

RUN if [ "$APP_ENV" = "test" ]; then RUN rm -f .env; cp .env.test .env; fi

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN make tidy


# Собираем приложение
RUN make build

RUN make migration-up

# Указываем порт, на котором будет работать приложение
EXPOSE 9000

# Команда для запуска приложения
CMD ["sh", "-c", "./build/app APP_ENV=${APP_ENV}"]
