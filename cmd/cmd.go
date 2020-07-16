package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stilvoid/clap/pres"
	"gopkg.in/yaml.v2"
)

// Root represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "clap",
	Long: "Clap is a yaml-based presentation tool for the shell",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fn := args[0]

		f, err := os.Open(fn)
		if err != nil {
			panic(fmt.Errorf("Unable to open '%s': %s", fn, err))
		}

		p := pres.Pres{}
		d := yaml.NewDecoder(f)
		err = d.Decode(&p)
		if err != nil {
			panic(fmt.Errorf("Invalid yaml in '%s': %s", fn, err))
		}

		p.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
