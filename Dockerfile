FROM alpine:latest

RUN apk add --no-cache go imagemagick

WORKDIR testdir
COPY . .

ENTRYPOINT ["go", "test", "-v", "./...", "-coverprofile", "cover.out"]
