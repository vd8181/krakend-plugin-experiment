FROM --platform=linux/arm64/v8 alpine:3.21.3 as builder
RUN apk add --no-cache build-base git go musl musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY validator.go ./

ENV CGO_ENABLED=1 GOOS=linux GOARCH=arm64

RUN go get github.com/luraproject/lura/v2@v2.9.0 \
     github.com/luraproject/lura/v2/proxy@v2.9.0 \
     github.com/luraproject/lura/v2/config@v2.9.0 \
     github.com/luraproject/lura/v2/sd@v2.9.0

RUN go build -buildmode=plugin -o validator.so validator.go

FROM --platform=linux/arm64/v8 devopsfaith/krakend:2.9.3

WORKDIR /app


COPY --from=builder /app/validator.so .
COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .
COPY krakend.json /app/krakend.json

# Expose the default KrakenD port
EXPOSE 8080

# Update CMD to use the krakend executable
CMD ["krakend", "run", "-c", "/app/krakend.json"]