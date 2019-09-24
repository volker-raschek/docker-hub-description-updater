FROM volkerraschek/build-image:1.4.0 AS build-env

ADD ./ /workspace

RUN make bin/linux/amd64/dhdu

FROM busybox:latest
COPY --from=build-env /workspace/bin/linux/amd64/dhdu /usr/bin/dhdu
ENTRYPOINT [ "/usr/bin/dhdu" ]