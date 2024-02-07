
FROM golang:1.21.5-alpine3.19 as builder

LABEL author="Joseph Akayesi <josephakayesi@gmail.com>"

WORKDIR /dist

COPY . .

RUN go mod download

RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/main ./cmd/main.go

FROM scratch

COPY --from=golang /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /etc/passwd /etc/passwd
COPY --from=golang /etc/group /etc/group

COPY --from=builder /bin/main /bin/main

EXPOSE 5000

ENTRYPOINT ["/bin/main"]
