FROM gcr.io/distroless/static

ARG bin
ADD $bin /location

EXPOSE 8080 2112

ENTRYPOINT ["/location"]
