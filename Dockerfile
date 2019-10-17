FROM golang:1.13.1 AS builder

WORKDIR /go/src/github.com/bhb603/lag

RUN go get -u github.com/golang/dep/cmd/dep
COPY Gopkg.* ./
RUN dep ensure -vendor-only

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/lag main.go

FROM alpine:latest
COPY --from=builder /bin/lag /bin/lag
CMD ["/bin/lag"]
