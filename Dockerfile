### Build Image ###
FROM --platform=$BUILDPLATFORM alpine:edge AS BUILD_IMAGE
ARG TARGETARCH
ARG VERSION

# Installing UnZip
RUN apk add --no-cache go --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community

# Build project
COPY . /coomer
WORKDIR /coomer
RUN sed -i "s/<version>/$VERSION/g" main.go
RUN GOARCH=$TARGETARCH go build -o coomer-dl

### Main Image ###
FROM alpine:edge
LABEL org.opencontainers.image.source="https://github.com/mysteryengineer/coomer-downloader"

# Dependencies
RUN apk add --no-cache libwebp-tools ffmpeg --repository=https://dl-cdn.alpinelinux.org/alpine/edge/community

# Define the image version
ARG VERSION
ENV IMAGE_VERSION=$VERSION

COPY --from=BUILD_IMAGE /coomer/coomer-dl /usr/local/bin/

CMD ["coomer-dl", "-d", "/tmp/coomer"]