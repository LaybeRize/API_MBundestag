FROM golang:1.19-alpine

WORKDIR /app

COPY database_old ./database
COPY ./dataLogic ./dataLogic
COPY help ./helper
COPY ./htmlHandler ./htmlHandler
COPY ./htmlWrapper ./htmlWrapper
COPY ./main ./main
COPY ./public ./public
COPY ./templates ./templates
COPY go.mod .
RUN go mod tidy

RUN go build -o /mbundestag /app/main

EXPOSE 8080

CMD ["/mbundestag"]