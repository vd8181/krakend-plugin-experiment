FROM --platform=linux/amd64 devopsfaith/krakend:2.9.3 as builder

RUN apk add --no-cache build-base git go
RUN apk add --no-cache gcc musl-dev git go

RUN apk update && apk upgrade musl

WORKDIR /app

COPY go.mod go.sum ./

# Pin dependencies as needed for KrakenD compatibility
RUN go mod tidy && go mod download

COPY validator.go ./

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

RUN go get github.com/luraproject/lura/v2@v2.9.0 \
     github.com/luraproject/lura/v2/proxy@v2.9.0 \
     github.com/luraproject/lura/v2/config@v2.9.0 \
     github.com/luraproject/lura/v2/sd@v2.9.0
RUN go build -buildmode=plugin -o validator.so validator.go

FROM devopsfaith/krakend:2.9.3

WORKDIR /app

COPY --from=builder /app/validator.so .
COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .

EXPOSE 3000

CMD ["check-plugin", "/app/validator.so"]

#
# ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
#
# # Optional: Only needed if your go.mod doesn't already include these
# RUN go get github.com/luraproject/lura/v2@v2.9.0 \
#     github.com/luraproject/lura/v2/proxy@v2.9.0 \
#     github.com/luraproject/lura/v2/config@v2.9.0 \
#     github.com/luraproject/lura/v2/sd@v2.9.0
#
# RUN go build -x -buildmode=plugin -o validator.so validator.go
#
# # Runtime stage
# FROM devopsfaith/krakend:2.9.3
#
# WORKDIR /app
#
# COPY --from=builder /app/validator.so .
COPY --from=builder /app/go.mod .
COPY --from=builder /app/go.sum .
#
# EXPOSE 3000
#
# CMD ["check-plugin", "/app/validator.so"]