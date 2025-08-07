FROM golang:latest

WORKDIR /go/src/app/

COPY . .

RUN apt update

RUN go mod download -x
# For hot reloading
RUN go install github.com/air-verse/air@latest

EXPOSE 3000

CMD ["air", "-c", ".air.toml"]