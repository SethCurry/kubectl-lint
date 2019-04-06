package cmd

import (
	"fmt"
	"os"

	"text/tabwriter"

	"github.com/SethCurry/kubectl-lint/internal/lint"
	"github.com/SethCurry/kubectl-lint/internal/linters"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
)

var cfgFile string

var builderFlags = genericclioptions.NewResourceBuilderFlags()
var configFlags = genericclioptions.NewConfigFlags(true)

var logger = zap.NewNop()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-lint",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if viper.GetBool("verbose") {
			var err error
			logger, err = zap.NewDevelopment()
			if err != nil {
				fmt.Println("Failed to create logger!")
				panic(err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		runner := lint.NewRunner(lint.WithLinters(linters.All()), lint.WithLogger(logger))

		disabledLinters := viper.GetStringSlice("disable")
		logger.Debug("disabling linters",
			zap.Strings("error_codes", disabledLinters),
		)
		for _, v := range viper.GetStringSlice("disable") {
			err := runner.DisableLinter(lint.ErrorCode(v))
			if err != nil {
				fmt.Println("Failed to disable linter.")
				panic(err)
			}
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		visitor := builderFlags.ToBuilder(configFlags, args).Do()
		visitErr := visitor.Visit(func(info *resource.Info, err error) error {
			err = info.Get()
			if err != nil {
				return err
			}

			lintErrors, err := runner.Lint(info)
			if err != nil {
				return err
			}

			for _, v := range lintErrors {
				fmt.Fprintf(w, "%s\t| %s.%s.%s\t| %s\t| %s\n", v.Severity.String(), v.SourceNamespace, v.SourceKind, v.SourceName, string(v.ErrorCode), string(v.ErrorMessage))
			}

			return nil
		})

		w.Flush()

		if visitErr != nil {
			panic(visitErr)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	builderFlags.WithLatest()

	configFlags.AddFlags(rootCmd.Flags())
	builderFlags.AddFlags(rootCmd.Flags())
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubectl-lint.yaml)")

	rootCmd.Flags().StringSliceP("disable", "d", []string{}, "Comma separated list of error codes to disable.")
	viper.BindPFlag("disable", rootCmd.Flags().Lookup("disable"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enables the built-in logger for debugging.")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".kubectl-lint" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".kubectl-lint")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
