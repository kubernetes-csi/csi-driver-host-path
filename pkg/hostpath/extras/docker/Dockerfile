FROM alpine
LABEL maintainers="Kubernetes Authors"
LABEL description="HostPath CSI Plugin"

COPY ./_output/hostpathplugin /hostpathplugin
COPY ./_output/livenessprobe /livenessprobe
ENTRYPOINT ["/hostpathplugin"]
