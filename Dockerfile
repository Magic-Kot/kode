FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

COPY ./ ./

# build
RUN go mod download
RUN go build -v -o ./bin/app cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY internal/config/config.yml /config.yml

CMD ["/app"]