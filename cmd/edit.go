package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hsnodgrass/pufctl/internal/helpers"
	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/internal/uitext"
	"github.com/hsnodgrass/pufctl/pkg/puppetfileparser/ast"
)

var (
	modName  string
	modProps []string

	editCmd = &cobra.Command{
		Use:   uitext.EditUse,
		Short: uitext.EditShort,
		Long:  uitext.EditLong,
	}

	editModuleCmd = &cobra.Command{
		Use:   uitext.EditModuleUse,
		Short: uitext.EditModuleShort,
		Long:  uitext.EditModuleLong,
		Args:  cobra.MaximumNArgs(1),
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
			pfilePath := viper.GetString("puppetfile")
			_show := helpers.MaxBools(show, viper.GetBool("always.show"))
			_writeInPlace := helpers.MaxBools(writeInPlace, viper.GetBool("always.write_in_place"))
			puppetfile, err := helpers.Parse(pfilePath, parseOpts)
			changes := false
			if err != nil {
				logging.Errorln(err)
			}
			if ok, _ := puppetfile.HasModule(args[0]); !ok {
				logging.Errorf("Module %s could not be found in Puppetfile\n", args[0])
			}
			mod := puppetfile.GetModule(args[0])
			if len(modProps) > 0 {
				for _, p := range modProps {
					parts := strings.Split(p, "=>")
					if len(parts) != 2 {
						logging.Errorf("Module property \"%s\" incorrectly formatted", p)
					}
					mod.EditProperty(parts[0], parts[1])
					changes = true
					if _show {
						fmt.Println(uitext.DashSep)
						mod.Sprint()
					}
				}
			}
			if modName != "" {
				err := puppetfile.RenameModule(args[0], modName)
				if err != nil {
					logging.Errorf("Renaming module \"%s\" failed with error: %w", args[0], err)
				}
				changes = true
			}
			editOutput(_show, _writeInPlace, confirm, changes, pfilePath, outFile, puppetfile)
		},
	}
)

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.AddCommand(editModuleCmd)
	editModuleCmd.Flags().StringVarP(&modName, "name", "n", "", "new module name")
	editModuleCmd.Flags().StringSliceVarP(&modProps, "key", "k", []string{}, "module property key=>value pairs")
}

func editOutput(_show, _writeInPlace, _confirm, _changes bool, _pfilePath, _outFile string, _puppetfile *ast.Puppetfile) {
	checkShow(_show, _puppetfile)
	if _changes {
		err := checkWriteInPlace(_writeInPlace, _confirm, _pfilePath, _puppetfile)
		if err != nil {
			logging.Errorln("Failed to write in place with error: ", err)
		}
	} else {
		logging.Infoln("No changes to write to Puppetfile ", _pfilePath)
	}
	err := checkOutFile(_confirm, _outFile, _pfilePath, _puppetfile)
	if err != nil {
		logging.Errorln("Failed to write output to file with error: ", err)
	}
}
