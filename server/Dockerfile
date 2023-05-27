
FROM golang as buildstage

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN env CGO_ENABLED=0 go build -o main .

FROM busybox
COPY --from=buildstage /app/main .
COPY --from=buildstage /app/data/users.json data/users.json

# CA certificates
COPY --from=buildstage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./main"]
