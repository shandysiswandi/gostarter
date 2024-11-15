FROM golang:1.19.13-alpine AS builder
RUN apk add --no-cache upx
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o server . && upx --best --lzma /app/server

FROM gcr.io/distroless/static-debian12:nonroot
ARG TZ="UTC"
ENV TZ=${TZ}
COPY --from=builder /app/server /server
EXPOSE 8081 50001
ENTRYPOINT ["/server"]