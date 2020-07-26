package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hsnodgrass/pufctl/internal/helpers"
	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/internal/uitext"
	"github.com/hsnodgrass/pufctl/pkg/puppetfileparser/ast"
	"github.com/hsnodgrass/pufctl/pkg/semver"
)

var (
	bMajor bool
	bMinor bool

	bumpCmd = &cobra.Command{
		Use:   uitext.BumpUse,
		Short: uitext.BumpShort,
		Long:  uitext.BumpLong,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pfilePath := viper.GetString("puppetfile")
			_show := helpers.MaxBools(show, viper.GetBool("always.show"))
			_writeInPlace := helpers.MaxBools(writeInPlace, viper.GetBool("always.write_in_place"))
			puppetfile, err := helpers.Parse(viper.GetString("puppetfile"), parseOpts)
			changes := false
			if err != nil {
				logging.Errorln(err)
			}
			for _, m := range args {
				var sv semver.SemVer
				var prop *ast.Property
				mod := puppetfile.GetModule(m)
				if mod == nil {
					logging.Warnf("Could not find module \"%s\" in Puppetfile\n", m)
					continue
				}
				prop = mod.GetProperty(":tag")
				if prop == nil {
					prop = mod.GetProperty(":ref")
					if prop == nil {
						logging.Warnf("Module \"%s\" has no :tag or :ref\n", mod.Name)
						continue
					}
				}
				sv, err := semver.Make(prop.Value.String)
				if err != nil {
					logging.Warnln("Could not parse semver for module", mod.Name)
				}
				if bMajor {
					sv.BumpMajor()
				} else if bMinor {
					sv.BumpMinor()
				} else {
					sv.BumpPatch()
				}
				if sv.ParsedAsV {
					prop.Value.String = sv.VString()
				} else {
					prop.Value.String = fmt.Sprintf("%s", sv)
				}
				changes = true
				if _show {
					fmt.Println(uitext.DashSep)
					fmt.Println(mod.Sprint())
				}
			}
			bumpOutput(_writeInPlace, confirm, changes, pfilePath, outFile, puppetfile)
		},
	}
)

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.Flags().BoolVarP(&bMajor, "major", "X", false, "bump Major version (X.y.z)")
	bumpCmd.Flags().BoolVarP(&bMinor, "minor", "Y", false, "bump Minor version (x.Y.z)")
}

func bumpOutput(_writeInPlace, _confirm, _changes bool, _pfilePath, _outFile string, _puppetfile *ast.Puppetfile) {
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
