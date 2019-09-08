package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kmwenja/rkyv"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rkyv"
	app.Commands = []cli.Command{
		createCmd(),
		updateCmd(),
		extractCmd(),
		infoCmd(),
		listCmd(),
		scanCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error while parsing arguments: %v\n", err)
		os.Exit(1)
	}
}

func createCmd() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create .rkyv file",
		Action: func(c *cli.Context) error {
			if len(c.Args()) == 0 {
				return fmt.Errorf("no arguments specified")
			}

			r := rkyv.Rkyv{}

			for _, filename := range c.Args() {
				stat, err := os.Stat(filename)
				if err != nil {
					return fmt.Errorf("could not get file stats `%s`: %v", filename, err)
				}

				if stat.IsDir() {
					return fmt.Errorf("directories are not supported: `%s`", filename)
				}

				data, err := ioutil.ReadFile(filename)
				if err != nil {
					return fmt.Errorf("could not read file `%s`: %v", filename, err)
				}

				r.AddFile(stat.Name(), data)
			}

			err := r.Flush()
			if err != nil {
				return fmt.Errorf("could not create .rkyv file: %v", err)
			}

			fmt.Printf("Created rkyv file: %s\n", r.Filename())
			return nil
		},
	}
}

func updateCmd() cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "update .rkyv file",
		Action: func(c *cli.Context) error {
			fmt.Println("update")
			return nil
		},
	}
}

func extractCmd() cli.Command {
	return cli.Command{
		Name:  "extract",
		Usage: "extract files from .rkyv file",
		Action: func(c *cli.Context) error {
			fmt.Println("extract")
			return nil
		},
	}
}

func infoCmd() cli.Command {
	return cli.Command{
		Name:  "info",
		Usage: "show details about .rkyv file",
		Action: func(c *cli.Context) error {
			if len(c.Args()) == 0 {
				return fmt.Errorf("no arguments specified")
			}

			filename := c.Args().Get(0)
			r, err := rkyv.OpenFile(filename)
			if err != nil {
				return err
			}

			data, err := json.Marshal(r)
			if err != nil {
				return fmt.Errorf("could not marshal json: %v", err)
			}

			fmt.Printf("%s\n", data)
			return nil
		},
	}
}

func listCmd() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list all indexed .rkyv files",
		Action: func(c *cli.Context) error {
			fmt.Println("list")
			return nil
		},
	}
}

func scanCmd() cli.Command {
	return cli.Command{
		Name:  "scan",
		Usage: "scan directories for .rkyv files then index them",
		Action: func(c *cli.Context) error {
			fmt.Println("scan")
			return nil
		},
	}
}
