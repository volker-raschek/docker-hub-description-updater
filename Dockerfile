FROM scratch AS build

COPY docker-hub-description-updater-* /usr/bin/docker-hub-description-updater

ENTRYPOINT [ "/usr/bin/docker-hub-description-updater" ]