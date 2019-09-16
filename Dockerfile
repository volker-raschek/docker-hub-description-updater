FROM volkerraschek/build-image:1.4.0 AS build-env

ADD ./ /workspace

RUN make bin/linux/amd64/dhd

FROM busybox:latest
COPY --from=build-env /workspace/bin/linux/amd64/dhd /usr/bin/dhd
ENTRYPOINT [ "/usr/bin/dhd" ]