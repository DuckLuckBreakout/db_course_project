FROM golang:1.15 as build
WORKDIR /db_course_project
COPY go.mod .
RUN go mod download
COPY . /db_course_project
RUN go build ./cmd/main.go

FROM ubuntu:20.04 as server
COPY . .
EXPOSE 5000
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get -y update && apt -y install postgresql-12
USER postgres
RUN  /etc/init.d/postgresql start &&\
    psql -f /init_db/init_db.sql -d postgres &&\
    /etc/init.d/postgresql stop
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/12/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/12/main/postgresql.conf
RUN echo "max_connections = 1200" >> /etc/postgresql/12/main/postgresql.conf
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root
COPY --from=build  /db_course_project/main /usr/bin
CMD /etc/init.d/postgresql start && main