# Этап 1: Компиляция двоичного файла в контейнеризованном окружении Golang
#
FROM golang:1.17 as build
# Скопировать исходные файлы с хоста
COPY . /src
# Назначить рабочим каталог с исходным кодом
WORKDIR /src
# Собрать двоичный файл!
RUN CGO_ENABLED=0 GOOS=linux go build -o book-catalog
# Этап 2: Сборка образа со службой хранилища пар ключ/значение
#
# Использовать образ "scratch", не содержащий распространяемых файлов
FROM scratch
# Скопировать двоичный файл из контейнера build
COPY --from=build /src/book-catalog .
# Если предполагается использовать TLS, скопировать файлы .pem
#COPY --from=build /src/*.pem .
# Сообщить фреймворку Docker, что служба будет использовать порт 8080
EXPOSE 8080
# Команда, которая должна быть выполнена при запуске контейнера
CMD ["/book-catalog"]