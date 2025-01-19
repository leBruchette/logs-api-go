FROM golang:1.23.5-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /lapi-server ./main/server.go
EXPOSE 8080
CMD ["/lapi-server"]

