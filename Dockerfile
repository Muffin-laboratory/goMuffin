FROM golang:1.24.5

RUN mkdir /app
WORKDIR /app

COPY . .

RUN make deps
RUN make

ENTRYPOINT [ "./build/goMuffin" ]
