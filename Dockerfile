FROM golang:1.20
LABEL authors="dlc"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /servMetrics-go
EXPOSE 8080
CMD ["/servMetrics-go"]