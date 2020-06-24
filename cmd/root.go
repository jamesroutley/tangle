package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/jamesroutley/tangle/tangle"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	watch      bool
	outfile    string
	configfile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tangle <file.md>",
	Short: "Extracts and concatenates code from code blocks in Markdown files",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
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
	rootCmd.Flags().StringVar(&configfile, "configfile", ".tangle.json", "Path to a config file")
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
	// The operation of tangle depends on whether a file to operate on has been
	// provided or not.
	// If one has, we'll extract code from it, using the default config. If
	// not, we extract code using the config in `configfile`
	// TODO: explicit mode?
	var config *Config
	if len(args) == 1 {
		config = defaultConfig(args[0], outfile)
	} else {
		var err error
		config, err = readConfig(configfile)
		if err != nil {
			return err
		}
	}

	group := errgroup.Group{}

	for _, target := range config.Targets {
		target := target
		group.Go(func() error {
			var options []tangle.TanglerOption
			if target.Filters.Language != "" {
				options = append(options, tangle.LanguageFilterOption(target.Filters.Language))
			}
			if target.Filters.Regex != "" {
				re, err := regexp.Compile(target.Filters.Regex)
				if err != nil {
					return err
				}
				options = append(options, tangle.RegexFilterOption(re))
			}

			tangler := tangle.NewTangler(options...)

			if watch {
				fw, err := newFileWatcher(target.Sources...)
				if err != nil {
					return err
				}

				for _ = range fw.events {
					log.Printf("Generating %s", target.Outfile)
					if err := runTangler(tangler, target.Sources, target.Outfile); err != nil {
						// In watch mode, we don't want to stop the binary, but
						// let the user fix the error
						fmt.Fprintln(os.Stderr, err)
					}
				}
				return nil
			}

			if err := runTangler(tangler, target.Sources, target.Outfile); err != nil {
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func runTangler(tangler *tangle.Tangler, sources []string, outputFile string) error {
	code, err := tangler.Tangle(sources...)
	if err != nil {
		return err
	}

	if outputFile != "" {
		if err := ioutil.WriteFile(outputFile, code, 0644); err != nil {
			return err
		}
	} else {
		fmt.Printf("%s", code)
	}

	return nil
}
