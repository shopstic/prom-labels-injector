FROM golang:1.17-buster as builder
WORKDIR /app
COPY . ./
RUN go build

FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=builder /app/prom-labels-injector ./
USER nonroot:nonroot
CMD ["/app/prom-labels-injector"]
