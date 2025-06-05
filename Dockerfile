FROM golang:1.24.3 AS builder
ARG VERSION
RUN test -n "$VERSION"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags="-X main.Version=$VERSION" -o prometheus-whistleblower

FROM alpine:3.22.0
WORKDIR /app
COPY --from=builder /app /app
RUN adduser -D -H -h /app appuser
RUN chown -R appuser:appuser /app
EXPOSE 8080
USER appuser
ENV GIN_MODE="release"
ENTRYPOINT ["/app/prometheus-whistleblower"]