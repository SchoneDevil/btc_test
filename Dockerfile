FROM golang:1.18 as builder
WORKDIR /usr/local/go/src/
ADD app/ /usr/local/go/src/
RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/main/app.go


FROM scratch
COPY --from=builder /usr/local/go/src/app /
EXPOSE 8080
CMD ["/app"]
