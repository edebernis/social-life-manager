FROM gcr.io/distroless/static

ARG bin
ADD $bin /location
COPY ./config /config

EXPOSE 8080 2112

ENTRYPOINT ["/location"]
