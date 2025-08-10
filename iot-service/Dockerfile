FROM golang:1.23-alpine3.21 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gualogger

# Stage 2
FROM alpine:3.21
WORKDIR /app
COPY --from=build /app/gualogger .
COPY --from=build /app/configs /app/configs
ENTRYPOINT [ "/app/gualogger" ]
