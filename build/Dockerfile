FROM alpine:latest

LABEL maintainer Fred <fred@gcreativeprojects.tech>

COPY catalogue /usr/bin/

VOLUME /catalogue
WORKDIR /catalogue

ENTRYPOINT ["catalogue"]
CMD ["--help"]
