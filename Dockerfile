FROM alpine
LABEL maintainers="Kubernetes Authors"
LABEL description="HostPath Driver"
ARG binary=./bin/hostpathplugin

# Add util-linux to get a new version of losetup.
RUN apk add util-linux && apk update && apk upgrade
COPY ${binary} /hostpathplugin
ENTRYPOINT ["/hostpathplugin"]
