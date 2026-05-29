FROM golang:1.22 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/service

FROM gcr.io/distroless/base-debian12
COPY --from=builder /out/app /app
USER 1000
ENTRYPOINT ["/app"]
