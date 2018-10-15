package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/rvadim/docker-registry-cleaner/registry"

	version "github.com/hashicorp/go-version"
	"github.com/urfave/cli"
)

type registryParams struct {
	URL      string
	username string
	password string
}

func main() {

	app := cli.NewApp()
	app.Name = "Docker Registry Cleaner"
	app.Version = "0.1.1"
	app.Compiled = time.Now()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "path",
			Usage:  "Regular docker path '<registry>/repo/image'",
			EnvVar: "PATH",
		},
		cli.StringFlag{
			Name:   "username, u",
			Usage:  "Registry username (optional)",
			EnvVar: "USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "Registry password (optional)",
			EnvVar: "PASSWORD",
		},
		cli.StringFlag{
			Name:   "imageversion, iv",
			Value:  ".*-SNAPSHOT.*",
			Usage:  "Image Version to delete, this can be a regex \".*-SNAPSHOT.*\"",
			EnvVar: "IMAGE_VERSION",
		},
		cli.IntFlag{
			Name:   "keep, k",
			Value:  3,
			Usage:  "The number of images you want to keep, usefully if you are deleting images by regex",
			EnvVar: "KEEP",
		},
		cli.BoolFlag{
			Name:   "dryrun, d",
			Usage:  "Do not actually delete anything",
			EnvVar: "DRYRUN",
		},
	}

	app.Action = func(c *cli.Context) error {
		rawURL := c.String("path")
		if !strings.HasPrefix(rawURL, "http") {
			rawURL = fmt.Sprintf("https://%s", rawURL)
		}
		url, err := url.Parse(rawURL)
		if err != nil {
			fmt.Printf("Unable to parse url: %s", err)
			return cli.ShowAppHelp(c)
		}
		path := url.Path
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}

		fmt.Printf("Prased hostname: %s\n", url.Hostname())
		fmt.Printf("Prased protocol: %s\n", url.Scheme)
		fmt.Printf("Prased image: %s\n", path)

		r := registryParams{
			URL:      fmt.Sprintf("%s://%s", url.Scheme, url.Hostname()),
			username: c.String("username"),
			password: c.String("password"),
		}

		imageName := path
		regxVersion := c.String("imageversion")
		numKeep := c.Int("keep")

		if r.URL == "" {
			return cli.ShowAppHelp(c)
		}

		if imageName == "" {
			return cli.ShowAppHelp(c)
		}

		hub, err := registry.New(r.URL, r.username, r.password)
		if err != nil {
			fmt.Printf("%s", err)
		}

		versionsRaw := []string{}

		tags, err := hub.Tags(imageName)

		// Get list of versions that match the image name
		for _, element := range tags {
			r := regexp.MustCompile(regxVersion)
			matches := r.FindString(element)
			if len(matches) != 0 {
				if matches == element {
					versionsRaw = append(versionsRaw, matches)
				}
			}
		}

		// Prep versions to be sorted
		versions := make([]*version.Version, len(versionsRaw))
		for i, raw := range versionsRaw {
			v, _ := version.NewVersion(raw)
			versions[i] = v
		}

		// After this, the versions are properly sorted from high to low
		sort.Sort(sort.Reverse(version.Collection(versions)))

		if c.Bool("dryrun") {
			fmt.Printf("\nDRY RUN - nothing will be deleted \n")
		}

		fmt.Printf("\nFound %d images that match, keeping the %d latest versions and deleting the rest \n", len(versions), numKeep)

		for i, v := range versions {
			//fmt.Printf("\nsorted version: %s\n", v)
			if i >= numKeep {
				fmt.Printf("\nDelete version: %s\n", v)

				if !c.Bool("dryrun") {
					// Get the manifest digest for the image
					digest, err := hub.ManifestDigest(imageName, v.String())
					fmt.Printf("Deleting Manifest Digest: %s\n", digest)
					// delete manifest
					err = hub.DeleteManifest(imageName, digest)
					if err != nil {
						fmt.Printf("%s", err)
					}
				}
			} else {
				fmt.Printf("\nKeep version: %s\n", v)
			}

		}
		return nil
	}

	app.Run(os.Args)
}
