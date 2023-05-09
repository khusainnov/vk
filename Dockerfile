FROM golang:1.20 as build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o vkbot main.go

FROM gcr.io/distroless/base-debian11

COPY --from=build app/vkbot .
COPY third_party ./third_party

EXPOSE 5000

CMD ["/vkbot"]
