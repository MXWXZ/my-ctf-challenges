FROM golang:alpine as builder
RUN apk --no-cache add git ca-certificates
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w"

FROM golang:alpine as prod
WORKDIR /app
COPY --from=0 /app/darkforest .
COPY ./map.txt .
CMD ["./darkforest"]
