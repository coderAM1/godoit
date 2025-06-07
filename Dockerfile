FROM golang:1.24.4-alpine3.22 as builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o bin/intserver integration/intserver/*.go
RUN mv bin/* /

FROM alpine:3.22

USER 1001:1001

COPY --from=builder --chown=1001:1001 --chmod=750 /intserver /server

WORKDIR /

CMD ["/bin/sh", "-c","/server"]