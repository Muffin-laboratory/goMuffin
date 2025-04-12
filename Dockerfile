FROM golang:1.24.1

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go build -o build/goMuffin git.wh64.net/muffin/goMuffin

ENTRYPOINT [ "./build/goMuffin" ]