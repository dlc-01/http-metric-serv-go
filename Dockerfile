FROM golang:1.20 as builder
LABEL authors="dlc"
WORKDIR /App
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /App/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /server \
from alphine \
copy --builder /server  /server/App
EXPOSE 8080
CMD ["/servMetrics-go"]