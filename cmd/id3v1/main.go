package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/xonyagar/id3/v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "ID3v1"
	app.Usage = "reads id3v1 tags"
	app.Description = "an id3v1 tag reader"
	app.Version = "0.1.0"
	app.Commands = commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func commands() []cli.Command {

	return []cli.Command{
		{
			Name:   "version",
			Usage:  "Return id3v1 version",
			Action: commandVersion,
		},
		{
			Name:   "show",
			Usage:  "Show id3v1 tag",
			Action: commandShow,
		},
		{
			Name:   "clear",
			Usage:  "Clear id3v1 tag",
			Action: commandClear,
		},
		{
			Name:   "remove",
			Usage:  "Remove id3v1 tag",
			Action: commandRemove,
		},
		{
			Name:   "set",
			Usage:  "Set id3v1 tag",
			Action: commandSet,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title",
					Usage: "Set title field (Maximum 30 character)",
				},
				cli.StringFlag{
					Name:  "artist",
					Usage: "Set artist field (Maximum 30 character)",
				},
				cli.StringFlag{
					Name:  "album",
					Usage: "Set album field (Maximum 30 character)",
				},
				cli.StringFlag{
					Name:  "year",
					Usage: "Set year field (Maximum 4 character)",
				},
				cli.StringFlag{
					Name:  "comment",
					Usage: "Set comment field (Maximum 30 character in v1, 28 character in v1.1)",
				},
				cli.IntFlag{
					Name:  "album_track",
					Usage: "Set album track field (Only in v1.1, Range from 1 to 255)",
				},
				cli.IntFlag{
					Name:  "genre_index",
					Usage: "Set genre index field (Range from 0 to 255, Also set version to 1.1)",
				},
				cli.StringFlag{
					Name:  "genre",
					Usage: "Set genre field by title",
				},
				cli.BoolFlag{
					Name:  "tov10",
					Usage: "Set id3 version to 1",
				},
				cli.BoolFlag{
					Name:  "tov11",
					Usage: "Set id3 version to 1.1",
				},
			},
		},
	}
}

func commandVersion(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := v1.New(f)
	if err != nil {
		return err
	}

	fmt.Println(tag.Version())

	return nil
}

func commandSet(c *cli.Context) error {
	f, err := os.OpenFile(c.Args().First(), os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := v1.New(f)
	if err != nil && err != v1.ErrTagNotFound {
		return err
	}

	if title := c.String("title"); title != "" {
		tag.SetTitle(title)
	}

	if artist := c.String("artist"); artist != "" {
		tag.SetArtist(artist)
	}

	if album := c.String("album"); album != "" {
		tag.SetAlbum(album)
	}

	if year := c.String("year"); year != "" {
		tag.SetYear(year)
	}

	if comment := c.String("comment"); comment != "" {
		tag.SetComment(comment)
	}

	if genre := c.String("genre"); genre != "" {
		tag.SetGenre(genre)
	}

	if albumTrack := c.Int("album_track"); albumTrack != 0 {
		tag.SetAlbumTrack(albumTrack)
	}

	if genreIndex := c.Int("genre_index"); genreIndex != 0 {
		tag.SetGenreIndex(genreIndex)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	err = tag.Write(f)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

func commandClear(c *cli.Context) error {
	f, err := os.OpenFile(c.Args().First(), os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := v1.New(f)
	if err != nil {
		return err
	}

	tag.Clear()

	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	err = tag.Write(f)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

func commandRemove(c *cli.Context) error {
	f, err := os.OpenFile(c.Args().First(), os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	err = v1.Remove(f)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

func commandShow(c *cli.Context) error {
	f, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer f.Close()

	tag, err := v1.New(f)
	if err != nil {
		return err
	}

	fmt.Println(c.Args().First())
	fmt.Printf("Version: %s\n", tag.Version())
	fmt.Printf("Title: %s\n", tag.Title())
	fmt.Printf("Artist: %s\n", tag.Artist())
	fmt.Printf("Album: %s\n", tag.Album())
	fmt.Printf("Year: %s\n", tag.Year())
	fmt.Printf("Comment: %s\n", tag.Comment())
	if tag.Version() == v1.Version11 {
		fmt.Printf("Track: %d\n", tag.AlbumTrack())
	}
	fmt.Printf("Genre: %s (%d)\n", tag.Genre(), tag.GenreIndex())

	return nil
}
