FROM golang:1.13-buster as build
WORKDIR /go/src/app
ADD ./mural /go/src/app
RUN go build ./cmd/mural/main.go -o /go/bin/app

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /
CMD ["/app"]