# build stage
FROM golang:1.10-alpine
COPY . ./src/fragcenter
WORKDIR /go/src/fragcenter
RUN go build

# final stage
FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /app
COPY --from=build-env /go/src/fragcenter/fragcenter /app/
CMD ["/app/fragcenter"]