FROM golang:1.24.1

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /my_app 

CMD ["/my_app"]