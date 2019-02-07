FROM alpine
LABEL maintainers="Kubernetes Authors"
LABEL description="HostPath Driver"

COPY ./bin/hostpathplugin /hostpathplugin
ENTRYPOINT ["/hostpathplugin"]
