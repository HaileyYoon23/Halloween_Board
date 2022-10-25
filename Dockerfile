FROM golang:latest

WORKDIR /opt/app

COPY . .

RUN go mod download

RUN go get \
    && go get github.com/getsentry/sentry-go \
    && go get github.com/gorilla/mux \
    && go get github.com/labstack/gommon \
    && go get github.com/go-sql-driver/mysql

RUN go build -o main .

#ENTRYPOINT ["/bin/bash", "-c"]
CMD ["./main"]
