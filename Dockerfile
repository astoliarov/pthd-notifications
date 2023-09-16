FROM golang:1.21-alpine3.18 AS builder

# Dependencies for build
RUN apk --no-cache add make build-base

WORKDIR /go/src/pthd-notifications
COPY . .
RUN make build


# Final image.
FROM alpine:3.18

RUN apk --no-cache add bash ca-certificates tzdata build-base
ENV TZ Europe/Minsk

WORKDIR "/pthd-notifications"
COPY --from=builder /go/src/pthd-notifications/bin/api .
COPY --from=builder /go/src/pthd-notifications/bin/async-api .
CMD ["/pthd-notifications/api"]
