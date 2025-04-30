FROM golang:1.23.5 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN go mod tidy 
RUN go build -o /app/main ./cmd 

FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 9998
CMD ["./main"]

