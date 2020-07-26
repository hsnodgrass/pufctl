package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"

	pconf "github.com/hsnodgrass/pufctl/internal/config"
	"github.com/hsnodgrass/pufctl/internal/helpers"
	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/internal/uitext"
	pufctlver "github.com/hsnodgrass/pufctl/internal/version"
	"github.com/hsnodgrass/pufctl/pkg/puppetfileparser/ast"
)

var (
	cfgFile          string
	inFile           string
	verbose          bool
	confirm          bool
	forgeapiurl      string
	outFile          string
	show             bool
	versionsOnly     bool
	sshKeyPath       string
	docPath          string
	mdDocPath        string
	yamlDocPath      string
	manDocPath       string
	licenseFlag      bool
	globalConf       bool
	puppetfileBranch string
	user             string
	pass             string
	token            string

	parseOpts helpers.ParseOptions

	rootCmd = &cobra.Command{
		Use:     uitext.RootUse,
		Short:   uitext.RootShort,
		Long:    uitext.RootLong,
		Version: pufctlver.Version(),
		Run: func(cmd *cobra.Command, args []string) {
			if licenseFlag {
				fmt.Println(pufctlver.License())
			} else {
				cmd.Help()
			}
		},
	}

	docGenCmd = &cobra.Command{
		Use:   uitext.DocGenUse,
		Short: uitext.DocGenShort,
		Long:  uitext.DocGenLong,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				logging.Warnln("Command pufctl docgen does not accept args and will ignore them")
			}
			logging.Infof("Generating markdown docs at %s\n", viper.GetString("genopts.doc_gen_path"))
			err := doc.GenMarkdownTree(rootCmd, viper.GetString("genopts.doc_gen_path"))
			if err != nil {
				logging.Errorln("docgen failed to generate docs with error: ", err)
			}
		},
	}

	confGenCmd = &cobra.Command{
		Use:   uitext.ConfGenUse,
		Short: uitext.ConfGenShort,
		Long:  uitext.ConfGenLong,
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			_confirm := helpers.MaxBools(confirm, viper.GetBool("always.confirm"))
			var conf string
			if outFile != "" {
				conf = outFile
			} else {
				conf = pconf.FullConfPath
			}
			// if flag --global not specified, explicitly set config back to defaults
			if !globalConf {
				setAllConfigDefaults()
			}
			if _confirm {
				logging.Infoln("Overwriting config file ", conf)
				err := viper.WriteConfigAs(conf)
				if err != nil {
					logging.Errorln("genconfig failed to write config file with error: ", err)
				}
			} else {
				logging.Infoln("Attempting to create config file ", conf)
				err := viper.SafeWriteConfigAs(conf)
				if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
					logging.Errorf("%s! Please use -y to override", err)
				} else {
					logging.Errorln("genconfig failed to write config file with error: ", err)
				}
			}
		},
	}

	showCmd = &cobra.Command{
		Use:   uitext.ShowUse,
		Short: uitext.ShowShort,
		Long:  uitext.ShowLong,
		PreRun: func(cmd *cobra.Command, args []string) {
			parseOpts = helpers.ParseOptions{
				Verbose:  helpers.MaxBools(viper.GetBool("always.verbose"), verbose),
				GitRef:   viper.GetString("puppetfile_branch"),
				SSHKey:   viper.GetString("auth.ssh_key"),
				Username: viper.GetString("auth.username"),
				Password: viper.GetString("auth.password"),
				Token:    viper.GetString("auth.token"),
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Showing sorted and organized Puppetfile ", viper.GetString("puppetfile"))
			fmt.Println(uitext.StarSep)
			puppetfile, err := helpers.Parse(viper.GetString("puppetfile"), parseOpts)
			if err != nil {
				logging.Errorln("Failed to parse Puppetfile! Error: ", err)
			}
			if versionsOnly {
				for k, v := range puppetfile.ModuleVersionMap {
					fmt.Printf("%s: %s\n", k, v)
				}
			} else {
				fmt.Println(uitext.DashSep)
				fmt.Printf("%s", puppetfile.Sprint())
				fmt.Println(uitext.DashSep)
			}
		},
	}
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.Errorln("Pufctl failed to execute! Error: ", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolVarP(&licenseFlag, "license", "L", false, "show the license shorthand statement")
	rootCmd.PersistentFlags().StringVarP(&inFile, "puppetfile", "p", pconf.Puppetfile, "path to the Puppetfile to parse")
	rootCmd.PersistentFlags().StringVar(&puppetfileBranch, "puppetfile-branch", pconf.PuppetfileBranch, "The branch to use for a Puppetfile from Git")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", pconf.FullConfPath, "path to config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", pconf.AlwaysVerbose, "verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&confirm, "confirm", "y", false, "skip all confirmation checks")
	rootCmd.PersistentFlags().StringVar(&forgeapiurl, "forge-api", pconf.ForgeAPI, "Puppet Forge API URL")
	rootCmd.PersistentFlags().StringVarP(&outFile, "out-file", "o", "", "Write command output or changed Puppetfile to specified file")
	rootCmd.PersistentFlags().BoolVarP(&show, "show", "s", pconf.AlwaysShow, "Show Puppetfile after each command")
	rootCmd.PersistentFlags().StringVar(&sshKeyPath, "ssh-key", pconf.SSHKeyPath, "Path to your SSH key")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Forge / Git authentication token")
	rootCmd.PersistentFlags().StringVar(&user, "user", "", "Forge / Git authentication username")
	rootCmd.PersistentFlags().StringVar(&pass, "pass", "", "Forge / Git authentication password")

	viper.BindPFlag("auth.username", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("auth.password", rootCmd.PersistentFlags().Lookup("pass"))
	viper.BindPFlag("auth.token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("puppetfile", rootCmd.PersistentFlags().Lookup("puppetfile"))
	viper.BindPFlag("puppetfile_branch", rootCmd.PersistentFlags().Lookup("puppetfile-branch"))
	viper.BindPFlag("always.verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("always.confirm", rootCmd.PersistentFlags().Lookup("confirm"))
	viper.BindPFlag("forge.api_url", rootCmd.PersistentFlags().Lookup("forge-api"))
	viper.BindPFlag("always.show", rootCmd.PersistentFlags().Lookup("show"))
	viper.BindPFlag("auth.ssh_key", rootCmd.PersistentFlags().Lookup("ssh-key"))

	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVar(&versionsOnly, "versions-only", false, "Show a truncated for of Puppetfile with only <module name>: <version>")

	rootCmd.AddCommand(docGenCmd)
	docGenCmd.Flags().StringVar(&docPath, "doc-path", pconf.DocGenPath, "Path where docs will be generated")
	viper.BindPFlag("genopts.doc_gen_path", docGenCmd.Flags().Lookup("doc-path"))

	rootCmd.AddCommand(confGenCmd)
	docGenCmd.Flags().BoolVarP(&globalConf, "current", "g", false, "Include current global flag values in generated config file")

}

func initConfig() {
	initViperDefaults()

	viper.SetEnvPrefix("PUF")
	viper.AutomaticEnv()

	viper.SetConfigName(pconf.ConfigFile)
	viper.SetConfigType(pconf.ConfigFileExtension)
	viper.AddConfigPath(pconf.ConfPath)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err == nil {
		_verbose := helpers.MaxBools(verbose, viper.GetBool("always.verbose"))
		if _verbose {
			fmt.Println(uitext.DashSep)
			fmt.Println("Using config file:", viper.ConfigFileUsed())
			fmt.Println(uitext.DashSep)
		}
	}
}

func initViperDefaults() {
	viper.SetDefault("forge", pconf.ForgeStrDefaults)
	viper.SetDefault("forge.user_agent", pufctlver.UserAgent())
	viper.SetDefault("always", pconf.AlwaysBoolDefaults)
	viper.SetDefault("auth", pconf.AuthStrDefaults)
	viper.SetDefault("genopts", pconf.GenoptsStrDefaults)
	viper.SetDefault("puppetfile", pconf.Puppetfile)
	viper.SetDefault("puppetfile_branch", pconf.PuppetfileBranch)
}

func setAllConfigDefaults() {
	viper.Set("forge", pconf.ForgeStrDefaults)
	viper.Set("always", pconf.AlwaysBoolDefaults)
	viper.Set("auth", pconf.AuthStrDefaults)
	viper.Set("genopts", pconf.GenoptsStrDefaults)
	viper.Set("puppetfile", pconf.Puppetfile)
}

func checkShow(s bool, puppetfile *ast.Puppetfile) {
	if s {
		fmt.Println(uitext.DashSep)
		fmt.Println("Your Puppetfile starts below the dashed line")
		fmt.Println(uitext.DashSep)
		fmt.Printf("%s", puppetfile.Sprint())
	}
}

func checkWriteInPlace(wip, c bool, pfilePath string, puppetfile *ast.Puppetfile) error {
	if wip {
		err := helpers.PromptConfirmFile(pfilePath, puppetfile.Sprint(), c)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkOutFile(c bool, outFile, pfilePath string, puppetfile *ast.Puppetfile) error {
	if outFile != "" {
		err := helpers.PromptConfirmFile(outFile, puppetfile.Sprint(), c)
		if err != nil {
			return err
		}
	}
	return nil
}
