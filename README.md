# Coomer Downloader

Download all the pictures/videos posts from a particular user on Coomer. It also removes duplicates and convert files to better formats, with higher image quality and smaller file sizes.

## 🎥 Demo

[![asciicast](https://asciinema.org/a/BHbFOMd9pIaXRykcPYFfJAbgE.svg)](https://asciinema.org/a/BHbFOMd9pIaXRykcPYFfJAbgE)

## 🖼️ Usage

There are 2 ways to use this app: through the CLI tool or using Docker. Here are some points to consider to help you choose which solution is best for you:

1. **CLI tool**: if you just intend to use __coomer-dl__ to download the files, but have no intention to automatically convert them to better/smaller formats (AVIF/AV1), then stick with the CLI tool.

    - *Dependencies:* You can also use the CLI tool to [convert files](#convert-imagesvideos), however you must make sure that you have the dependencies [libavif](https://github.com/AOMediaCodec/libavif) and [FFmpeg](https://www.ffmpeg.org/download.html) before using the CLI tool. You can verify if all dependencies are properly installed by running the command `coomer-dl check-deps`.

2. **Docker**: if you want __coomer-dl__ to not only download the files, but also convert them to better formats, then Docker is probably a better option since it comes with all required dependencies installed.

### CLI Tool

Download the [latest version](https://github.com/mysteryengineer/coomer-downloader/releases) of __coomer-dl__ that matches your computer architecture and operating system. Extract the .zip file somewhere and then run the command below in the terminal:

```
$ coomer-dl -s onlyfans -u atomicbrunette18 -d /Downloads/Coomer
```

Where:

- `-s` (mandatory): the service where the files were originally posted; `onlyfans` or `fansly`.
- `-u` (mandatory): the user that you want to download images from.
- `-d` (optional): the directory where you want the files to be saved; default is the current directory.

For the full list of parameters, type `coomer-dl --help` in the terminal.

### Docker

Install [Docker](https://docs.docker.com/get-docker/) in your computer, then run the command below:

```
$ docker run --rm -t \
    -e COOMER_SERVICE=onlyfans \
    -e COOMER_USER=atomicbrunette18 \
    -v "/path/in/your/computer:/tmp/coomer" \
    ghcr.io/mysteryengineer/coomer-downloader
```

Where:

- `-e COOMER_SERVICE`: (mandatory): where the files were originally posted; `onlyfans` or `fansly`.
- `-e COOMER_USER`: (mandatory): the user that you want to download images from.

#### Volume

For those that are not familiar with Docker, the `-v` (volume) parameter defines where the media will be saved, and it's divided in two parts, separated by the colon sign `:`. You just need to worry about the first part, on the left side of the colon sign (**don't change anything on the right side**) and update it according to a path in your computer where you want the media to be downloaded.

For example, if you want to download the media in the directory `/Downloads/Coomer` then the volume parameter should look like this `-v "/Downloads/Coomer:/tmp/coomer"`.

## 💡 Features

### Remove duplicates

This application will automatically delete all files that are identical.

### Convert images/videos

You can convert the media downloaded to better formats (AVIF for images and AV1 for videos); this will make the files smaller but preserving the same quality. To do that you must:

- **CLI tool:** add the parameters `--convert-images` and/or `--convert-videos`, depending on what you need.

- **Docker:** add the environment variables `-e COOMER_CONVERT_IMAGES=true` and/or `-e COOMER_CONVERT_VIDEOS=true` when you run the container.

## 🛠️ Build

In the project's root folder run in the CLI:

Go:
```
$ go build -o coomer-dl
```

Docker:
```
$ docker build -t mysteryengineer/coomer-downloader . --build-arg="VERSION=23.12.17"
```

## 📈 Telemetry

This app collects information about the data that you're downloading to help me track bugs and improve the general stability of the software.

**No identifiable information about you or your computer is tracked.** But if you still want to stop the telemetry, you can do that by adding the flag `--no-telemetry` in the CLI tool or the environment variable `-e COOMER_TELEMETRY=false` when you run the Docker container.

## 📝 License

**coomer-dl** is released under the MIT License. See [LICENSE](LICENSE) for details.