package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/xonyagar/id3"
	"log"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "ID3"
	app.Usage = "reads id3 tags"
	app.Description = "an id3 tag reader"
	app.Version = "0.1.0"
	app.Commands = commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func commands() []cli.Command {
	versionFlags := []cli.Flag{
		cli.BoolFlag{
			Name:  "v1",
			Usage: "Show only version 1",
		},
		cli.BoolFlag{
			Name:  "v22",
			Usage: "Show only version 2.2.0",
		},
		cli.BoolFlag{
			Name:  "v23",
			Usage: "Show only version 2.3.0",
		},
		cli.BoolFlag{
			Name:  "v24",
			Usage: "Show only version 2.4.0",
		},
	}

	return []cli.Command{
		{
			Name:   "title",
			Usage:  "Return title",
			Action: commandTitle,
			Flags:  versionFlags,
		},
		{
			Name:   "artists",
			Usage:  "Return artist(s)",
			Action: commandArtists,
			Flags:  versionFlags,
		},
		{
			Name:   "album",
			Usage:  "Return album",
			Action: commandAlbum,
			Flags:  versionFlags,
		},
		{
			Name:   "album-artists",
			Usage:  "Return album artist(s)",
			Action: commandAlbumArtists,
			Flags:  versionFlags,
		},
		{
			Name:   "year",
			Usage:  "Return year",
			Action: commandYear,
			Flags:  versionFlags,
		},
		{
			Name:   "genres",
			Usage:  "Return genre(s)",
			Action: commandGenres,
		},
		{
			Name:   "track-number-and-position",
			Usage:  "Return track number and position",
			Action: commandTrackNumberAndPosition,
		},
	}
}

func commandTitle(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := id3.New(f)
	if err != nil {
		return err
	}

	if c.Bool("v1") {
		if tag.V1 != nil {
			fmt.Println(tag.V1.Title())
		}
	} else if c.Bool("v22") {
		if tag.V22 != nil {
			fmt.Println(tag.V22.Title())
		}
	} else if c.Bool("v23") {
		if tag.V23 != nil {
			fmt.Println(tag.V23.Title())
		}
	} else if c.Bool("v24") {
		if tag.V24 != nil {
			fmt.Println(tag.V24.Title())
		}
	} else {
		fmt.Println(tag.Title())
	}

	return nil
}

func commandAlbum(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := id3.New(f)
	if err != nil {
		return err
	}

	if c.Bool("v1") {
		if tag.V1 != nil {
			fmt.Println(tag.V1.Album())
		}
	} else if c.Bool("v22") {
		if tag.V22 != nil {
			fmt.Println(tag.V22.Album())
		}
	} else if c.Bool("v23") {
		if tag.V23 != nil {
			fmt.Println(tag.V23.Album())
		}
	} else if c.Bool("v24") {
		if tag.V24 != nil {
			fmt.Println(tag.V24.Album())
		}
	} else {
		fmt.Println(tag.Album())
	}

	return nil
}

func commandArtists(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := id3.New(f)
	if err != nil {
		return err
	}

	if c.Bool("v1") {
		if tag.V1 != nil {
			fmt.Println(tag.V1.Artist())
		}
	} else if c.Bool("v22") {
		if tag.V22 != nil {
			fmt.Println(strings.Join(tag.V22.Artists(), ", "))
		}
	} else if c.Bool("v23") {
		if tag.V23 != nil {
			fmt.Println(strings.Join(tag.V23.Artists(), ", "))
		}
	} else if c.Bool("v24") {
		if tag.V24 != nil {
			fmt.Println(strings.Join(tag.V24.Artists(), ", "))
		}
	} else {
		fmt.Println(strings.Join(tag.Artists(), ", "))
	}

	return nil
}

func commandAlbumArtists(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := id3.New(f)
	if err != nil {
		return err
	}

	if c.Bool("v1") {
		if tag.V1 != nil {
			fmt.Println(tag.V1.Artist())
		}
	} else if c.Bool("v22") {
		if tag.V22 != nil {
			fmt.Println(strings.Join(tag.V22.AlbumArtists(), ", "))
		}
	} else if c.Bool("v23") {
		if tag.V23 != nil {
			fmt.Println(strings.Join(tag.V23.AlbumArtists(), ", "))
		}
	} else if c.Bool("v24") {
		if tag.V24 != nil {
			fmt.Println(strings.Join(tag.V24.AlbumArtists(), ", "))
		}
	} else {
		fmt.Println(strings.Join(tag.Artists(), ", "))
	}
	return nil
}

func commandYear(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := id3.New(f)
	if err != nil {
		return err
	}

	if c.Bool("v1") {
		if tag.V1 != nil {
			fmt.Println(tag.V1.Year())
		}
	} else if c.Bool("v22") {
		if tag.V22 != nil {
			fmt.Println(tag.V22.Year())
		}
	} else if c.Bool("v23") {
		if tag.V23 != nil {
			fmt.Println(tag.V23.Year())
		}
	} else if c.Bool("v24") {
		if tag.V24 != nil {
			fmt.Println(tag.V24.Year())
		}
	} else {
		fmt.Println(tag.Album())
	}

	return nil
}

func commandTrackNumberAndPosition(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := id3.New(f)
	if err != nil {
		return err
	}

	a, b := tag.TrackNumberAndPosition()
	fmt.Printf("%d/%d\n", a, b)

	return nil
}

func commandGenres(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := id3.New(f)
	if err != nil {
		return err
	}

	fmt.Println(strings.Join(tag.Genres(), ", "))

	return nil
}
