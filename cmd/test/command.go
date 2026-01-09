package test

import (
	"strings"

	"github.com/maximal/cloudflare-speedtest/internal/exit"
	"github.com/spf13/cobra"
)

var TestCmd = &cobra.Command{
	Use:   PROGRAM_NAME,
	Short: SHORT_DESCRIPTION,
	Long:  LONG_DESCRIPTION,
	Run: func(_ *cobra.Command, args []string) {
		if !validateFlags() {
			exit.Exit(exit.StatusInvalidFlag)
		}

		if flags.Version {
			print(PROGRAM_TITLE + " v" + VERSION)
			exit.Exit(exit.StatusOk)
		}

		if len(args) != 0 {
			stderrRed("No arguments allowed; got: %s", strings.Join(args, ", "))
			exit.Exit(exit.StatusInvalidArg)
		}

		exit.Exit(run())
	},
}

func init() {
	// cobra.OnInitialize(initConfig)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	TestCmd.PersistentFlags().
		BoolVarP(&flags.Insecure, "insecure", "i", false, "make HTTP requests instead of HTTPS/TLS ones")

	TestCmd.PersistentFlags().
		BoolVarP(&flags.Version, "version", "v", false, "print version information and exit")

	TestCmd.PersistentFlags().
		StringVarP(&flags.Format, "format", "f", "text", "output format: text, json, jsonl, influx, tsv")

	TestCmd.PersistentFlags().
		BoolVarP(&flags.NoProgress, "no-progress", "n", false, "do not print progress information to STDERR")

	// Not yet fully implemented
	// TestCmd.PersistentFlags().
	//	StringVarP(&flags.SoftLimitString, "soft-limit", "s", "0", "soft time limit: 0 (no limit), 66 (seconds), 45s, 1m30s, ...")

	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	// rootCmd.AddCommand(addCmd)
	// rootCmd.AddCommand(initCmd)
}
