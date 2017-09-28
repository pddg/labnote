package main

import (
	"github.com/pddg/labnote/md"
	"github.com/urfave/cli"
	"fmt"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "labnote"
	app.Version = "0.0.0"
	app.Description = "test description"
	app.Commands = []cli.Command{
		{
			Name: "parse",
			Usage: "Parse file.",
			Action: parse,
		},
	}
	app.Run(os.Args)
}

func parse(c *cli.Context) {
	if c.NArg() == 0 {
		fmt.Println("Please specify input file.")
		os.Exit(1)
	}
	args := c.Args()
	markdown := new(md.Markdown)
	if err := markdown.SetPath(args[0]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	markdown.SetTags()
	if markdown.SetDate() {
		if err := markdown.ModFirstLine(); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("Name:\t%s\n", markdown.Path)
	fmt.Printf("Head:\t%s\n", markdown.FirstLine)
	if markdown.HasDate {
		fmt.Printf("Date:\t%s\n", markdown.Date)
	}
	if markdown.HasTags {
		fmt.Printf("Tags:\t%s\n", markdown.Tags)
	}
}