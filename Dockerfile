FROM golang:alpine3.12 AS build
ENV GOPROXY=https://goproxy.io,direct
RUN apk add --no-cache --update gcc musl-dev
WORKDIR /app
COPY . .
RUN go build

FROM alpine:3.12
ENV GIN_MODE=release
WORKDIR /app
EXPOSE 80
VOLUME [ "/app/data" ]
COPY --from=build /app/short-url .
CMD ["./short-url"]
