FROM golang:1.15 as build
WORKDIR /db_course_project
COPY go.mod .
RUN go mod download
COPY . /db_course_project
RUN go build -o bin/server -v ./cmd/

FROM ubuntu:latest as server
RUN apt update && apt install ca-certificates -y && rm -rf /var/cache/apt/*
COPY --from=build /db_course_project/bin/server /
CMD ["./server"]

FROM ubuntu:latest
