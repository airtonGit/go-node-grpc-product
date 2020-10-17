FROM golang as build-env
WORKDIR /app
ADD . /app

RUN cd /app && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app

# FROM golang:alpine
FROM scratch
COPY --from=build-env /app/app /app/app
WORKDIR /app
ENTRYPOINT [ "./app"]
