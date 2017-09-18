FROM golang:1.8
WORKDIR /go/src/github.com/rightstuff/inv-api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inv-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/rightstuff/inv-api/inv-api .
CMD ["./inv-api"]
