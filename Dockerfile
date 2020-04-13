FROM golang:alpine AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -tags netgo -ldflags '-extldflags "-static" -s -w' -o /cloud ./cmd/cloud

FROM scratch

COPY --from=builder /cloud /cloud

EXPOSE 8080

ENTRYPOINT ["/cloud"]
