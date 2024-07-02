FROM golang:1.22.3

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

EXPOSE 7540

ENV TODO_PORT=7540 TODO_PASSWORD=abc123

RUN GOOS=linux GOARCH=amd64 go build -o ./todo_app

CMD [ "./todo_app" ]