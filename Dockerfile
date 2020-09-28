FROM --platform=$TARGETPLATFORM alpine

ARG TARGETARCH
ARG NAME
ARG VER
ENV NAME=$NAME VER=$VER

COPY config.json.example /config.json
COPY ./bin/$TARGETARCH /$NAME-$VER
RUN chmod +x /$NAME-$VER

RUN mkdir -p /data
RUN chmod -R 777 /data

ENTRYPOINT [ "sh", "-c", "/$NAME-$VER" ]