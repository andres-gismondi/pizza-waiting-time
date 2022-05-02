FROM golang:1.17.6-alpine

ENV GO111MODULE=on

WORKDIR /app
COPY . .

RUN go mod vendor
RUN go build -mod=vendor -ldflags "-X main.version=$(git describe HEAD || git log --format="%h" -1)" ./cmd/server

# final stage
FROM golang:1.17.6-alpine

COPY --from=builder /app /app
WORKDIR /app
EXPOSE 8080

ENTRYPOINT ["/app/server"]