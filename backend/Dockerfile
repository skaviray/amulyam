FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

FROM gcr.io/distroless/base-nossl:debug-nonroot

WORKDIR /

# COPY ./app.env .

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["/main"]