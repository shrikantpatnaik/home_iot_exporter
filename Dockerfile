
# docker build --force-rm=true -t home_iot_exporter .

# RUN:
# docker run -it -p 8080:8080 home_iot_exporter

FROM golang as builder

RUN mkdir /app
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app

# Go dep
RUN go get -u github.com/prometheus/client_golang/prometheus && \
    go get -u github.com/prometheus/client_golang/prometheus/promhttp && \
    go get -u github.com/gorilla/mux


# Build
RUN set -ex && \
  CGO_ENABLED=0 go build \
        -tags netgo \
        -o /app/home_iot_exporter \
        -v -a \
        -ldflags '-extldflags "-static"' && \
  ls

# Create the second stage with a basic image.
# this will drop any previous
# stages (defined as `FROM <some_image> as <some_name>`)
# allowing us to start with a fat build image and end up with
# a very small runtime image.

FROM busybox

# add compiled binary
COPY --from=builder /app/home_iot_exporter /home_iot_exporter

# run
EXPOSE 8080
CMD ["/home_iot_exporter"]
