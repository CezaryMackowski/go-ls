FROM golang:1.24 AS build-stage

WORKDIR /usr/src/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-ls .

FROM build-stage AS run-test-stage

CMD ["go", "test", "-v", "./..."]

FROM gcr.io/distroless/base-debian12 AS run-stage

WORKDIR /

COPY --from=build-stage /usr/src/go-ls app

ENTRYPOINT ["./app"]