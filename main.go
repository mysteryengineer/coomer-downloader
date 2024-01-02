package main

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var service string
	var user string
	var directory string
	var parallel int
	var limit int
	var noTelemetry bool
	var convertImages bool
	var convertVideos bool

	currentPath, _ := os.Getwd()
	extensions := make([]string, 0)

	app := &cli.App{
		Name:            "coomer-dl",
		Usage:           "a CLI tool to download files from https://coomer.su",
		UsageText:       "coomer-dl -s [onlyfans/fansly] -u [user] [global options]",
		Version:         "<version>",
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "service",
				Aliases:     []string{"s"},
				Usage:       "service where the files are hosted; 'onlyfans' or 'fansly'",
				Destination: &service,
				Category:    "Required:",
				EnvVars:     []string{"COOMER_SERVICE"},
				Action: func(context *cli.Context, s string) error {
					if s != "onlyfans" && s != "fansly" {
						return fmt.Errorf("Invalid service '%s'", service)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "user",
				Aliases:     []string{"u"},
				Usage:       "user that you want to download files from",
				Destination: &user,
				Category:    "Required:",
				EnvVars:     []string{"COOMER_USER"},
			},
			&cli.StringFlag{
				Name:        "dir",
				Aliases:     []string{"d"},
				Value:       currentPath,
				Usage:       "directory where the files will be saved",
				Destination: &directory,
				Category:    "Optional:",
				DefaultText: "current directory",
			},
			&cli.IntFlag{
				Name:        "parallel",
				Value:       3,
				Usage:       "the number of downloads to be done in parallel",
				Destination: &parallel,
				Category:    "Optional:",
				DefaultText: "3",
				EnvVars:     []string{"COOMER_PARALLEL"},
				Action: func(context *cli.Context, i int) error {
					if i < 1 || i > 5 {
						return fmt.Errorf("The number of parallel downloads should be between 1-5")
					}
					return nil
				},
			},
			&cli.IntFlag{
				Name:        "limit",
				Value:       1_000_000,
				Usage:       "the maximum number of files to be downloaded",
				Destination: &limit,
				Category:    "Optional:",
				EnvVars:     []string{"COOMER_LIMIT"},
				DefaultText: "all files",
				Action: func(context *cli.Context, i int) error {
					if i < 1 {
						return fmt.Errorf("The number of max downloads should be at least 1")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "extensions",
				Usage:       "filter the downloads to only certain file extensions, separated by comma",
				Category:    "Optional:",
				EnvVars:     []string{"COOMER_EXTENSIONS"},
				DefaultText: "all extensions",
				Action: func(context *cli.Context, s string) error {
					split := strings.Split(s, ",")
					split = funk.Map(split, func(ext string) string {
						return "." + strings.ToLower(strings.TrimSpace(ext))
					}).([]string)
					split = funk.Filter(split, func(ext string) bool {
						return len(ext) > 1
					}).([]string)

					extensions = append(extensions, split...)
					return nil
				},
			},
			&cli.BoolFlag{
				Name:               "no-telemetry",
				Value:              false,
				Usage:              "if you want to disable the telemetry",
				Destination:        &noTelemetry,
				Category:           "Optional:",
				DisableDefaultText: true,
				EnvVars:            []string{"COOMER_TELEMETRY"},
			},
			&cli.BoolFlag{
				Name:               "convert-images",
				Value:              false,
				Usage:              "enable the conversion of images to AVIF",
				Destination:        &convertImages,
				Category:           "Optional:",
				DisableDefaultText: true,
				EnvVars:            []string{"COOMER_CONVERT_IMAGES"},
			},
			&cli.BoolFlag{
				Name:               "convert-videos",
				Value:              false,
				Usage:              "enable the conversion of videos to AV1",
				Destination:        &convertVideos,
				Category:           "Optional:",
				DisableDefaultText: true,
				EnvVars:            []string{"COOMER_CONVERT_VIDEOS"},
			},
		},
		Action: func(cCtx *cli.Context) error {
			if service == "" {
				return fmt.Errorf("Required flag '--service', '-s' is missing")
			}

			if user == "" {
				return fmt.Errorf("Required flag '--user', '-u' is missing")
			}

			expandedDir, err := ExpandPath(directory)
			if err != nil {
				return fmt.Errorf("Directory path %s is invalid", directory)
			}

			err = startJob(
				cCtx.App.Version,
				service,
				user,
				expandedDir,
				parallel,
				limit,
				extensions,
				noTelemetry,
				convertImages,
				convertVideos,
			)

			return err
		},
		Commands: []*cli.Command{
			{
				Name:  "check-deps",
				Usage: "check if you have all the dependencies installed in your computer",
				Action: func(cCtx *cli.Context) error {
					CheckDeps()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		PrintError(err.Error())
	}
}

// region - Private functions

func startJob(
	version string,
	service string,
	user string,
	directory string,
	parallel int,
	limit int,
	extensions []string,
	noTelemetry bool,
	convertImages bool,
	convertVideos bool,
) error {
	if IsOutdated(version, "mysteryengineer/coomer-downloader") {
		pterm.Println(pterm.Yellow("\n✨ There's a new version of Coomer Downloader available for download!"))
	}

	name, err := CheckUser(service, user)
	if err != nil {
		return err
	}

	if !noTelemetry {
		TrackDownloadStart(version, service, name, parallel, limit, false, false)
	}

	fullDir := filepath.Join(directory, name)
	medias := GetMedias(service, user, fullDir, limit)
	numMedias := len(medias)

	if len(extensions) > 0 {
		medias = FilterExtensions(medias, extensions)
	}

	downloads := DownloadMedias(medias, parallel)
	successes := funk.Filter(downloads, func(download Download) bool { return download.IsSuccess }).([]Download)
	failures := funk.Filter(downloads, func(download Download) bool { return !download.IsSuccess }).([]Download)

	duplicated, successes := RemoveDuplicates(successes)

	if !noTelemetry {
		TrackDownloadEnd(version, service, name, numMedias, len(failures), duplicated)
	}

	CreateReport(fullDir, downloads)

	ConvertMedia(successes, convertImages, convertVideos)

	pterm.Println("\n🌟 Done!")
	return err
}

// region
