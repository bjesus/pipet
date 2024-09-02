package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/bjesus/pipet/common"
	"github.com/bjesus/pipet/internal/app"
	"github.com/bjesus/pipet/outputs"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	app := &cli.App{
		Name:  "pipet",
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
			&cli.IntFlag{
				Name:  "interval",
				Value: 0,
				Usage: "Maximum number of pages to scrape",
			},
			&cli.StringFlag{
				Name:  "on-change",
				Usage: "Path to template file for output",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("spec argument is required")
			}
			spec := c.Args().Get(0)
			return runPipet(c, spec)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runPipet(c *cli.Context, specFile string) error {
	jsonOutput := c.Bool("json")
	separators := c.StringSlice("separator")
	templateFile := c.String("template")
	onChange := c.String("on-change")
	maxPages := c.Int("max-pages")
	interval := c.Int("interval")

	pipet := &common.PipetApp{
		MaxPages:  maxPages,
		Separator: separators,
	}

	log.Println("Parsing spec file:", specFile)
	err := app.ParseSpecFile(pipet, specFile)
	if err != nil {
		return fmt.Errorf("error parsing spec file: %w", err)
	}

	iterate := true
	previousValue := ""

	for iterate {
		newValue := ""
		log.Println("Executing blocks")
		err = app.ExecuteBlocks(pipet)
		if err != nil {
			return fmt.Errorf("error executing blocks: %w", err)
		}

		log.Println("Generating output")
		if jsonOutput {
			newValue = outputs.OutputJSON(pipet)
		} else if templateFile != "" {
			newValue = outputs.OutputTemplate(pipet, templateFile)
		} else {
			newValue = outputs.OutputText(pipet)
		}

		fmt.Print(newValue)

		if interval > 0 {
			if onChange != "" && previousValue != newValue {
				command := strings.ReplaceAll(onChange, "{}", strconv.Quote(newValue))
				log.Println("Executing on change command: " + command)
				cmd := exec.Command("bash", "-c", command)
				cmd.Output()
				previousValue = newValue
			}
			pipet.Data = []interface{}{}
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			iterate = false
		}
	}
	return nil
}
