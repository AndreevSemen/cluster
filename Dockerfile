FROM ubuntu

ENV GOPATH /src/go

RUN apt-get update                 && \
    apt-get -y install python3        \
                       python3-pip    \
                       golang         \
                       nginx          \
                       git         && \
    pip3 install requests          && \
    go get encoding/json              \
           net/http                   \
           log                        \
           time                       \
           github.com/gorilla/mux     

EXPOSE 80

COPY . .

ENTRYPOINT nginx -c /nginx.conf           & \
           go run REST-api/http_table.go  & \
           go run REST-api/http_time.go

