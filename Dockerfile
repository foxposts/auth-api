FROM golang:1.19-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go


FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/main .
COPY api.env .
CMD ["./main"]