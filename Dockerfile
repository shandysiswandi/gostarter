FROM golang:1.19.13-alpine AS builder
RUN apk add --no-cache upx tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o server .
RUN upx --best --lzma /app/server

FROM scratch
ARG TZ="UTC"
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
ENV TZ=${TZ}
COPY --from=builder /app/server /server
EXPOSE 8081 50001
ENTRYPOINT ["/server"]
