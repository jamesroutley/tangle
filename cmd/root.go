package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jamesroutley/tangle/tangle"
	"github.com/spf13/cobra"
)

var (
	watch   bool
	outfile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tangle <file.md>",
	Short: "Extracts and concatenates code from code blocks in Markdown files",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(cmd, args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().BoolVar(&watch, "watch", false, "Watch the input file, and recompile when it changes")
	rootCmd.Flags().StringVar(&outfile, "outfile", "", "Name of a file to write the output to. Writes to stdout if none provided.")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	filename := args[0]

	tangler := tangle.NewTangler()
	code, err := tangler.Tangle(filename)
	if err != nil {
		return err
	}

	if outfile != "" {
		if err := ioutil.WriteFile(outfile, code, 0644); err != nil {
			return err
		}
	} else {
		fmt.Printf("%s", code)
	}

	return nil
}
