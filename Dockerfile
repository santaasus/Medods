FROM golang:1.22.4 as builder
WORKDIR /srv/medods-app
COPY . .
RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o medods .

FROM golang:1.22.4
WORKDIR /srv/medods-app
COPY --from=builder /srv/medods-app/config.json .
COPY --from=builder /srv/medods-app/scheme.sql .
COPY --from=builder /srv/medods-app/medods .

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8080

CMD ["/go/bin/dlv", "--listen=:8080", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/srv/medods-app/medods"]