FROM golang:alpine

# install git, ffmpeg, etc
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base && apk add --no-cache ffmpeg && apk add --no-cache tzdata

# ref https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o /build app/api/main.go

# Expose port 8888 to the outside world
EXPOSE 8888

# Run the executable
CMD [ "/build" ]