FROM golang:1.11-alpine
WORKDIR /go/src/github.com/rvadim/docker-registry-cleaner
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docker-registry-cleaner .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/rvadim/docker-registry-cleaner/docker-registry-cleaner .
ENTRYPOINT ["/root/docker-registry-cleaner"]
CMD ["/root/docker-registry-cleaner"]
