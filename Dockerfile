FROM golang:1.10.0-alpine AS go
RUN apk --no-cache add git
WORKDIR /go/src/github.com/sudarshan-reddy/globule

COPY . /go/src/github.com/sudarshan-reddy/globule

RUN go build 

FROM alpine:3.7
WORKDIR /root/

COPY --from=go /go/src/github.com/sudarshan-reddy/globule .

ENTRYPOINT ["./globule"]
