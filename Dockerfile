FROM alpine:latest
ADD bin/changelog-generator.linux /usr/local/bin/changelog-generator
WORKDIR /usr/local/bin
ENTRYPOINT ["changelog-generator"]
