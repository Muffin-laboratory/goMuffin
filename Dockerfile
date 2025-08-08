FROM golang:1.24.5

RUN mkdir /app
WORKDIR /app

COPY . .

RUN make

ENTRYPOINT [ "./build/goMuffin" ]
