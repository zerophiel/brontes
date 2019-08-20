FROM golang:latest AS Builder
ADD . /src
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN git config --global http.sslverify false
WORKDIR /src/main
RUN go build -a -installsuffix cgo -o /src/bin/app -mod vendor

# Run Stage
FROM alpine
ENV GIN_MODE=release
RUN apk update && apk add ca-certificates
EXPOSE 8080
WORKDIR /app
COPY --from=Builder /src/bin/app /app
ENTRYPOINT ./app