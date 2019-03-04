FROM alpine
LABEL maintainers="Kubernetes Authors"
LABEL description="HostPath Driver"

# Add util-linux to get a new version of losetup.
RUN apk add util-linux
COPY ./bin/hostpathplugin /hostpathplugin
ENTRYPOINT ["/hostpathplugin"]
