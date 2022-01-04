FROM golang:alpine AS build
RUN mkdir /go/src/report-downloader
ADD ./ /go/src/report-downloader
WORKDIR /go/src/report-downloader
RUN apk update && apk add --no-cache git && apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main cmd/main.go

FROM scratch AS runtime
COPY --from=build /main /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
CMD ["/main"]



