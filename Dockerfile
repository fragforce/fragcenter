# build stage
FROM golang:1.10-alpine AS build-env
COPY . ./src/fragcenter
WORKDIR /go/src/fragcenter
RUN go build

# final stage
FROM alpine:latest
RUN apk add --no-cache bash \
 && mkdir -p /app/public
WORKDIR /app
COPY --from=build-env /go/src/fragcenter/fragcenter /app/
CMD ["/app/fragcenter"]