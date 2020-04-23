FROM golang:1.14-stretch

WORKDIR /go-prometheus

COPY . .

RUN go mod download
RUN go get github.com/cespare/reflex

COPY reflex.conf /

EXPOSE 3001

ENTRYPOINT ["reflex", "-c", "/reflex.conf"]
