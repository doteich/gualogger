FROM golang:1.24-alpine3.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o geist-connector

# Stage 2
FROM alpine:3.22
WORKDIR /app
COPY --from=build /app/geist-connector .
# COPY --from=build /app/configs /app/configs
ENTRYPOINT [ "/app/geist-connector" ]
