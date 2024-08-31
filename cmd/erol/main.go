package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/bjesus/erol/common"
	"github.com/bjesus/erol/internal/app"
	"github.com/bjesus/erol/outputs"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	app := &cli.App{
		Name:  "erol",
		Usage: "Easy web scraping CLI tool",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output as JSON",
			},
			&cli.StringSliceFlag{
				Name:  "separator",
				Usage: "Separator for text output (can be used multiple times)",
			},
			&cli.StringFlag{
				Name:  "template",
				Usage: "Path to template file for output",
			},
			&cli.IntFlag{
				Name:  "max-pages",
				Value: 3,
				Usage: "Maximum number of pages to scrape",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("spec argument is required")
			}
			spec := c.Args().Get(0)
			return runErol(c, spec)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runErol(c *cli.Context, specFile string) error {
	jsonOutput := c.Bool("json")
	separators := c.StringSlice("separator")
	templateFile := c.String("template")
	maxPages := c.Int("max-pages")

	erol := &common.ErolApp{
		MaxPages:  maxPages,
		Separator: separators,
	}

	log.Println("Parsing spec file:", specFile)
	err := app.ParseSpecFile(erol, specFile)
	if err != nil {
		return fmt.Errorf("error parsing spec file: %w", err)
	}

	log.Println("Executing blocks")
	err = app.ExecuteBlocks(erol)
	if err != nil {
		return fmt.Errorf("error executing blocks: %w", err)
	}

	log.Println("Generating output")
	if jsonOutput {
		return outputs.OutputJSON(erol)
	} else if templateFile != "" {
		return outputs.OutputTemplate(erol, templateFile)
	} else {
		return outputs.OutputText(erol)
	}
}
