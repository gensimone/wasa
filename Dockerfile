FROM golang:1.25.1

WORKDIR /src/

COPY . .

RUN go build -o ./webapi ./cmd/webapi

EXPOSE 3000

VOLUME /data

CMD ["/src/webapi"]
