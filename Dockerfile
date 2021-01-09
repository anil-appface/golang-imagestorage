# Start from golang base image
FROM golang:alpine as builder

# WORKDIR /go/src/github.com/anil/golang-imagestorage

RUN apk update && apk add --no-cache git

WORKDIR /app

# RUN go get -d -v golang.org/x/net/html
# RUN go get -d -v github.com/gorilla/mux
# RUN go get -d -v github.com/lib/pq 
# RUN go get -d -v github.com/sirupsen/logrus 
# RUN go get -d -v github.com/go-sql-driver/mysql

# COPY handlers /usr/src/app/handlers
# COPY store /usr/src/app/store
# COPY utils /usr/src/app/utils
# COPY main.go app.go go.mod go.sum  /usr/src/app/

# RUN go mod download 

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



# RUN CGO_ENABLED=0 GOOS=linux go build -o golang-imagestorage

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .       

# Expose port 8080 to the outside world
EXPOSE 8081

#Command to run the executable
CMD ["./main"]

# FROM alpine:latest

# RUN apk --no-cache add ca-certificates

# # WORKDIR /root/
# # COPY --from=builder /go/src/github.com/anil/golang-imagestorage .

# WORKDIR /app/golang-imagestorage
# COPY --from=build-env /go/src/github.com/anil/golang-imagestorage/golang-imagestorage /app/golang-imagestorage
# ENTRYPOINT ["/app/golang-imagestorage"]

# CMD ["./app/golang-imagestorage"]