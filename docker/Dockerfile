#FROM golang:1.16
#WORKDIR /app/gohttpserver
#ADD . /app/gohttpserver
#RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-X main.VERSION=docker' -o gohttpserver

FROM ubuntu:22.04
WORKDIR /app
RUN mkdir -p /app/public
RUN apt-get update && apt-get install -y ca-certificates
VOLUME /app/public
ADD assets ./assets
# COPY --from=0 /app/gohttpserver/gohttpserver .
COPY static-server .
# EXPOSE 8000
ENTRYPOINT [ "/app/static-server" ]
CMD []
