FROM golang:1.24.4-alpine AS builder

ARG UPX_VERSION=4.2.4-r0

RUN apk update && apk add --no-cache ca-certificates upx=$UPX_VERSION tzdata && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 64000 \
  dnsupdater

COPY . .
RUN go build -buildvcs=false -tags netgo -trimpath -tags netgo -ldflags="-w -s" -o ./dnsupdater cmd/main.go

RUN upx --best --lzma dnsupdater

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=dnsupdater:dnsupdater /app/dnsupdater /app/dnsupdater

USER dnsupdater:dnsupdater

ENTRYPOINT ["/app/dnsupdater"]
