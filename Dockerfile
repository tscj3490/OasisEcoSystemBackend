FROM golang:1.9-alpine AS BUILDER
RUN apk add -U git
COPY . /go
WORKDIR /go/src/com/merkinsio/oasis-api/
RUN go get github.com/golang/dep/cmd/dep
RUN cd /go/src/com/merkinsio/oasis-api/ && dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.6
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=BUILDER /go/src/com/merkinsio/oasis-api /
RUN chmod +x /oasis-api
EXPOSE 4040
ENTRYPOINT ["/oasis-api"]
