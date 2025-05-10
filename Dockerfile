FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o archviz ./cmd/archviz

FROM scratch

COPY --from=builder /build/archviz /
COPY --from=builder /build/cmd/archviz/static /static

ENTRYPOINT ["/archviz"]
