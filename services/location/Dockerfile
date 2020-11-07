FROM golang:1.15-alpine3.12 AS builder
RUN apk add --no-cache git make
WORKDIR /go/src/github.com/edebernis/social-life-manager/services/location
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN make build

FROM alpine:3.12
RUN apk add --no-cache ca-certificates make
COPY --from=builder /go/src/github.com/edebernis/social-life-manager/services/location/ /
WORKDIR /

EXPOSE 8080
EXPOSE 2112

# The entrypoint is the makefile, so calling "docker-compose run app <make target>"
# calls the given makefile target in the app environment
ENTRYPOINT ["make"]
CMD ["run"]
