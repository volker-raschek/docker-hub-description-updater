ARG BASE_IMAGE
ARG BUILD_IMAGE
ARG EXECUTABLE_TARGET
ARG VERSION

# BUILD
# ==============================================================================
FROM ${BUILD_IMAGE} AS build-env

ADD ./ /workspace

RUN make clean ${EXECUTABLE_TARGET} GOPROXY=${GOPROXY}

# TARGET
# ==============================================================================
FROM ${BASE_IMAGE}
COPY --from=build-env /workspace/${EXECUTABLE_TARGET} /usr/bin/dhdu
RUN chmod +x /usr/bin/dhdu
ENTRYPOINT [ "/usr/bin/dhdu" ]