FROM golang:1.24.3

ENV DATABASE_NAME=muffin_ai

RUN mkdir /app
WORKDIR /app

COPY . .

RUN make

ENTRYPOINT [ "./build/goMuffin" ]