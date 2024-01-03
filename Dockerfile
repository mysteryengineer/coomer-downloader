### Build Image ###
FROM --platform=$BUILDPLATFORM alpine:edge AS BUILD_IMAGE
ARG TARGETARCH
ARG VERSION

# Dependencies
RUN apk add --no-cache wget --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community

# Download the binary
RUN mkdir /coomer
WORKDIR /coomer
RUN wget https://github.com/mysteryengineer/coomer-downloader/releases/download/$VERSION/coomer-dl_linux_$TARGETARCH.zip
RUN unzip coomer-dl_linux_$TARGETARCH.zip

### Main Image ###
FROM alpine:edge
LABEL org.opencontainers.image.source="https://github.com/mysteryengineer/coomer-downloader"

# Dependencies
RUN apk add --no-cache libavif-apps ffmpeg --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community

# Define the image version
ENV IMAGE_VERSION=$VERSION

COPY --from=BUILD_IMAGE /coomer/coomer-dl /usr/local/bin/

CMD ["coomer-dl", "-d", "/tmp/coomer"]