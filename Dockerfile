FROM golang:1.26.0

# 폴더 생성
RUN mkdir /app
WORKDIR /app

# 코드 복사
COPY . .

# 종속성 설치 & 빌드
RUN make deps
RUN make

ENTRYPOINT [ "./build/goMuffin" ]
