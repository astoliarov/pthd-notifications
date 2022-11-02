# Builder image, where we build the example.
FROM golang:1.19-alpine3.16 AS builder

# Dependencies for build
RUN apk --no-cache add make build-base

WORKDIR /go/src/pthd-notifications
COPY . .
RUN make build


# Final image.
FROM alpine:3.16.1

RUN apk --no-cache add bash ca-certificates tzdata build-base
ENV TZ Europe/Minsk

WORKDIR "/pthd-notifications"
COPY --from=builder /go/src/pthd-notifications/bin/api .
CMD ["/pthd-notifications/api"]
