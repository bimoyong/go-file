FROM --platform=$TARGETPLATFORM alpine:latest AS builder
ARG GIT_AUTH
ARG GOPRIVATE
ENV GIT_AUTH=$GIT_AUTH
ENV GOPRIVATE=$GOPRIVATE
RUN apk --no-cache add make git go gcc libtool musl-dev protoc
RUN git config --global url."https://$GIT_AUTH@github.com/bimoyong".insteadOf "https://github.com/bimoyong"

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin && \
    go get github.com/golang/protobuf/protoc-gen-go && \
    cp /go/bin/protoc-gen-go /usr/bin/

# Compile Go source
COPY . /src
WORKDIR /src
RUN cd /src && \
    make && \
    go clean -cache -modcache -i -r

FROM --platform=$TARGETPLATFORM alpine:latest
ARG NAME
ARG VER
ENV NAME=$NAME
ENV VER=$VER

COPY ./config.json.example /config.json
COPY --from=builder /src/bin/app /$NAME-$VER

ENTRYPOINT [ "sh", "-c", "/$NAME-$VER" ]