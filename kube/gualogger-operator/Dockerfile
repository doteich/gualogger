FROM rust:alpine AS builder

WORKDIR /app/src
RUN USER=root

RUN apk add pkgconfig openssl-dev libc-dev
COPY ./ ./
RUN cargo build --release

FROM alpine:latest

WORKDIR /app
RUN apk update \
    && apk add openssl ca-certificates

EXPOSE 80 443

COPY --from=builder /app/src/target/release/gualogger-operator /app/gualogger-operator

CMD ["/app/gualogger-operator"]
