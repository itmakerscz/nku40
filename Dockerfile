# The base go-image
FROM golang:1.17-alpine as builder

RUN go get github.com/go-sql-driver/mysql
 
# Create a directory for the app
RUN mkdir /hackujstat
 
# Copy all files from the current directory to the app directory
COPY . /hackujstat
 
# Set working directory
WORKDIR /hackujstat
 
# Run command as described:
# go build will build an executable file named server in the current directory
# RUN go mod vendor
# RUN go help modules


RUN go mod init itmakers.cz/hackujstat

RUN go mod download

RUN go mod vendor

# Run the Go build and output binary
RUN go build -o main .

# Make sure to expose the port the HTTP server is using
#EXPOSE 8080
# Run the app binary when we run the container
#ENTRYPOINT ["/hackujstat"]

# Run the server executable
CMD [ "/hackujstat/main" ]