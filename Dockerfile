FROM golang:1.24.2

RUN mkdir /app
WORKDIR /app

COPY . .

RUN make

ENTRYPOINT [ "./build/goMuffin" ]