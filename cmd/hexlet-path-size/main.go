package main

import (
	"code"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "hexlet-path-size",
		Usage:     "print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		UsageText: "hexlet-path-size [global options] <path>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "recursive",
				DefaultText: "false",
				Usage:       "recursive size of directories",
				Aliases:     []string{"r"},
			},
			&cli.BoolFlag{
				Name:        "human",
				DefaultText: "false",
				Usage:       "human-readable sizes (auto-select unit)",
				Aliases:     []string{"H"},
			},
			&cli.BoolFlag{
				Name:        "all",
				DefaultText: "false",
				Usage:       "include hidden files and directories",
				Aliases:     []string{"a"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			path := cmd.Args().Get(0)
			if path == "" {
				return errors.New("path is required")
			}

			recursive := cmd.Bool("recursive")
			human := cmd.Bool("human")
			all := cmd.Bool("all")
			size, err := code.GetPathSize(path, recursive, human, all)
			if err != nil {
				return err
			}

			fmt.Printf("%s\t%s\n", size, path)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
