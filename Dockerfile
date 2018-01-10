FROM golang:alpine
ARG release_version="0.3.2"

RUN apk add --update curl git

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v${release_version}/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/kochie
WORKDIR /go/src/github.com/kochie

COPY . .

RUN dep ensure

RUN go-wrapper install

CMD ["go-wrapper", "run"]
