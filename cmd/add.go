package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hsnodgrass/pufctl/internal/auth"
	pconf "github.com/hsnodgrass/pufctl/internal/config"
	"github.com/hsnodgrass/pufctl/internal/helpers"
	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/internal/sources/forgesource"
	"github.com/hsnodgrass/pufctl/internal/sources/gitsource"
	"github.com/hsnodgrass/pufctl/internal/uitext"
	"github.com/hsnodgrass/pufctl/internal/validators"
	"github.com/hsnodgrass/pufctl/pkg/forgeapi"
	"github.com/hsnodgrass/pufctl/pkg/puppetfileparser/ast"
)

var (
	modProperties  []string
	modMetadata    []string
	modResolveDeps bool
	writeInPlace   bool

	addCmd = &cobra.Command{
		Use:   uitext.AddUse,
		Short: uitext.AddShort,
		Long:  uitext.AddLong,
	}

	addModuleCmd = &cobra.Command{
		Use:   uitext.AddModuleUse,
		Short: uitext.AddModuleShort,
		Long:  uitext.AddModuleLong,
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			addModulePreRun()
		},
		Run: func(cmd *cobra.Command, args []string) {
			var newModReq *newModuleRequest
			pfilePath := viper.GetString("puppetfile")
			_show := helpers.MaxBools(show, viper.GetBool("always.show"))
			_writeInPlace := helpers.MaxBools(writeInPlace, viper.GetBool("always.write_in_place"))
			puppetfile, err := helpers.Parse(pfilePath, parseOpts)
			if err != nil {
				logging.Errorln(err)
			}
			newModReq, err = createNewModuleRequest(args)
			if err != nil {
				logging.Errorln("Failed to create new module request with error:", err)
			}
			logging.Infoln("Adding module to Puppetfile. This may take a few seconds.")
			changes, err := newModReq.AddFunc(*newModReq, puppetfile, modResolveDeps)
			if err != nil {
				cmd.Help()
				logging.Errorln(err)
			}
			addOutput(_show, _writeInPlace, confirm, changes, pfilePath, outFile, puppetfile)
		},
	}
	addMetaCmd = &cobra.Command{
		Use:   uitext.AddMetaUse,
		Short: uitext.AddMetaShort,
		Long:  uitext.AddMetaLong,
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			addMetaPreRun(args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			pfilePath := viper.GetString("puppetfile")
			_confirm := helpers.MaxBools(confirm, viper.GetBool("always.confirm"))
			_show := helpers.MaxBools(show, viper.GetBool("always.show"))
			_writeInPlace := helpers.MaxBools(writeInPlace, viper.GetBool("always.write_in_place"))
			puppetfile, err := helpers.Parse(pfilePath, parseOpts)
			if err != nil {
				logging.Errorln("Failed to parse Puppetfile with error:", err)
			}
			switch args[0] {
			case "top":
				for _, m := range modMetadata {
					err := puppetfile.AddComment(args[0], m)
					if err != nil {
						logging.Errorln("Failed to add top-block comment with error:", err)
					}
				}
			case "bottom":
				for _, m := range modMetadata {
					err := puppetfile.AddComment(args[0], m)
					if err != nil {
						logging.Errorln("Failed to add bottom-block comment with error:", err)
					}
				}
			default:
				if f, _ := puppetfile.HasModule(args[0]); f {
					for _, m := range modMetadata {
						tag := ast.SingleMatch(ast.ReMetaTag, m)
						data := ast.SingleMatch(ast.ReMetaData, m)
						puppetfile.AddModuleMetadata(args[0], tag, data)
					}
				} else {
					logging.Errorf("Module %s does not exist in the Puppetfile", args[0])
				}
			}
			checkShow(_show, puppetfile)
			checkWriteInPlace(_writeInPlace, _confirm, pfilePath, puppetfile)
			checkOutFile(_confirm, outFile, pfilePath, puppetfile)
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.AddCommand(addModuleCmd)
	addModuleCmd.Flags().BoolVarP(&modResolveDeps, "resolve-deps", "D", false, "also adds dependecies of module to the Puppetfile")
	addModuleCmd.Flags().StringSliceVarP(&modProperties, "properties", "a", []string{}, "comma-separated list of properties (Ex. -p ':git=>https://fake.com,:ref=>production')")
	addModuleCmd.Flags().StringSliceVarP(&modMetadata, "metadata", "m", []string{}, "comma-separated list of metadata (Ex. -m '# @maintainer: team@fake.com','# @critical:')")
	addModuleCmd.Flags().BoolVarP(&writeInPlace, "write-in-place", "w", pconf.AlwaysWriteInPlace, "Overwrite the current Puppetfile with changes")
	viper.BindPFlag("always.write_in_place", addModuleCmd.Flags().Lookup("write-in-place"))

	addCmd.AddCommand(addMetaCmd)
	addMetaCmd.Flags().StringSliceVarP(&modMetadata, "metadata", "m", []string{}, "comma-separated list of metadata (Ex. -m '# @maintainer: team@fake.com','# @critical:')")
	addMetaCmd.Flags().BoolVarP(&writeInPlace, "write-in-place", "w", pconf.AlwaysWriteInPlace, "Overwrite the current Puppetfile with changes")
	viper.BindPFlag("always.write_in_place", addMetaCmd.Flags().Lookup("write-in-place"))
}

type newModuleRequest struct {
	Input        string
	Slug         string
	URL          string
	UserAgent    string
	AddFunc      func(newModuleRequest, *ast.Puppetfile, bool) (bool, error)
	ForgeModMeta forgeapi.ModuleMetadata
	auth.Auth
	Mod
}

func (m *newModuleRequest) SetAuth(auth auth.Auth) {
	m.Auth = auth
}

// Mod holds Module properties and Metadata
type Mod struct {
	Props []string
	Meta  []string
}

// GetMod returns a pointer to Mod
func (m *Mod) GetMod() *Mod {
	return m
}

func addModulePreRun() {
	if !validators.IsGitURL(viper.GetString("puppetfile")) && !validators.IsGitURLNoSuffix(viper.GetString("puppetfile")) {
		err := helpers.ValidatePath(viper.GetString("puppetfile"))
		if err != nil {
			logging.Errorln("Failed to validate Puppetfile with error:", err)
		}
	}
	_s := helpers.MaxBools(show, viper.GetBool("always.show"))
	_w := helpers.MaxBools(writeInPlace, viper.GetBool("always.write_in_place"))
	var _o bool
	if outFile == "" {
		_o = false
	} else {
		_o = true
	}
	if _s || _w || _o {
		logging.Debugln("Command prechecks complete")
	} else {
		logging.Errorln("You must specify --show (-s), --write-in-place (-w), or --out-file (-o)")
	}
	parseOpts = helpers.ParseOptions{
		Verbose:  helpers.MaxBools(viper.GetBool("always.verbose"), verbose),
		GitRef:   viper.GetString("puppetfile_branch"),
		SSHKey:   viper.GetString("auth.ssh_key"),
		Username: viper.GetString("auth.username"),
		Password: viper.GetString("auth.password"),
		Token:    viper.GetString("auth.token"),
	}
}

func addMetaPreRun(args []string) {
	err := helpers.ValidatePath(viper.GetString("puppetfile"))
	if err != nil {
		logging.Errorln("Failed to validate Puppetfile with error:", err)
	}
	if len(modMetadata) < 1 {
		logging.Errorln("The flag --metadata (-m) is required")
	}
	if len(args) > 1 {
		logging.Warnln("pufctl add meta only accepts one argument, all others will be ignored")
	}
	switch args[0] {
	case "top":
		logging.Debugln("Intending to add metadata to top-block comments")
	case "bottom":
		logging.Debugln("Intending to add metadata to bottom-block comments")
	default:
		logging.Debugln("Intending to add metadata to module:", args[0])
	}
	parseOpts = helpers.ParseOptions{
		Verbose:  helpers.MaxBools(viper.GetBool("always.verbose"), verbose),
		GitRef:   viper.GetString("puppetfile_branch"),
		SSHKey:   viper.GetString("auth.ssh_key"),
		Username: viper.GetString("auth.username"),
		Password: viper.GetString("auth.password"),
		Token:    viper.GetString("auth.token"),
	}
}

func createNewModuleRequest(args []string) (*newModuleRequest, error) {
	var newModReq *newModuleRequest
	var normalized string
	isGit := validators.IsGitURL(args[0])
	isGitNoSuf := validators.IsGitURLNoSuffix(args[0])
	isModSlug := validators.IsModuleSlug(args[0])
	isModRef := validators.IsModuleRef(args[0])
	switch {
	case isGit:
		normalized = args[0]
	case isGitNoSuf:
		normalized = fmt.Sprintf("%s.git", args[0])
	case isModSlug:
		normalized = args[0]
	case isModRef:
		converted, err := helpers.ModStringToSlug(args[0])
		if err != nil {
			return nil, fmt.Errorf("Failed to convert \"%s\" to mod slug: %w", args[0], err)
		}
		normalized = converted
	}
	if isGit || isGitNoSuf {
		auth, err := auth.GitAuth(normalized, viper.GetString("auth.username"), viper.GetString("auth.password"), viper.GetString("auth.token"), viper.GetString("auth.ssh_key"))
		if err != nil {
			return nil, fmt.Errorf("Failed to create new module request: %w", err)
		}
		newModReq = &newModuleRequest{
			Input:   args[0],
			URL:     normalized,
			Mod:     Mod{Props: modProperties, Meta: modMetadata},
			Auth:    auth,
			AddFunc: gitAddModule,
		}
		return newModReq, nil
	}
	if isModSlug || isModRef {
		newModReq = &newModuleRequest{
			Input:     args[0],
			Slug:      normalized,
			URL:       viper.GetString("forge.api_url"),
			UserAgent: viper.GetString("forge.user_agent"),
			Mod:       Mod{Props: modProperties, Meta: modMetadata},
			AddFunc:   forgeAddModule,
		}
		return newModReq, nil
	}
	return nil, fmt.Errorf("Input %s is invalid, please use a module slug (<namespace>-<modulename>) or a git repo URL", args[0])
}

func forgeAddModule(mod newModuleRequest, puppetfile *ast.Puppetfile, resolveDeps bool) (bool, error) {
	changes := false
	logging.Debugln("Resolving module from Puppet Forge")
	fm, err := forgesource.GetModule(mod.Slug, mod.URL, mod.UserAgent)
	if err != nil {
		logging.Errorln(err)
	}
	logging.Debugln("Successfully retrieved module from Puppet Forge")
	p := mod.GetMod().Props
	p = append(p, fm.CurrentRelease.Version)
	for _, prop := range p {
		logging.Debugln("Proceeding with property", prop)
	}
	err = puppetfile.AddModule(fm.Slug, p)
	if err != nil && !resolveDeps {
		return changes, fmt.Errorf("Module %s already present in Puppetfile", fm.Slug)
	}
	if err != nil && resolveDeps {
		logging.Warnf("Module %s is already in Puppetfile. Resolving dependencies anyways.", fm.Slug)
	} else {
		logging.Debugln("Successfully added module", fm.Slug)
		changes = true
	}
	if resolveDeps {
		logging.Infoln("Resolving dependencies")
		depChanges := forgesource.ResolveDeps(fm, mod.URL, mod.UserAgent, puppetfile)
		if depChanges {
			logging.Infoln("Successfully resolved dependencies")
			changes = depChanges
		}
	}
	return changes, nil
}

func gitAddModule(mod newModuleRequest, puppetfile *ast.Puppetfile, resolveDeps bool) (bool, error) {
	changes := false
	logging.Debugln("Resolving module from Git source")
	defBranch, meta := gitsource.GetModuleMeta(mod.URL, mod.Auth)
	source := meta.Source
	if source == "" {
		logging.Warnln("Cannot parse source from metadata, using input", mod.Input)
		source = mod.Input
	}
	logging.Debugln("Proceeding with source", source)
	if !strings.HasSuffix(source, ".git") {
		source = fmt.Sprintf("%s.git", source)
		logging.Debugln("Proceeding with modified source", source)
	}
	slug := strings.Replace(meta.Name, "/", "-", 1)
	logging.Debugln("Proceeding with slug", slug)
	p := mod.GetMod().Props
	if len(p) < 1 {
		gitProp := fmt.Sprintf(":git=>%s", source)
		logging.Debugln("Git source is", gitProp)
		dbProp := fmt.Sprintf(":default_branch=>%s", defBranch)
		logging.Debugln("Default branch is", dbProp)
		p = []string{gitProp, dbProp}
	}
	err := puppetfile.AddModule(slug, p)
	if err != nil && !resolveDeps {
		return changes, fmt.Errorf("Module %s already present in Puppetfile", slug)
	} else if err != nil && resolveDeps {
		logging.Warnf("Module %s is already in Puppetfile. Resolving dependencies anyways.", slug)
	} else {
		logging.Infoln("Successfully added module", slug)
		changes = true
	}
	if resolveDeps {
		depChanges := gitsource.ResolveDeps(meta, puppetfile)
		if depChanges {
			logging.Infoln("Successfully resolved dependencies")
			changes = depChanges
		}
	}
	return changes, nil
}

func addOutput(_show, _writeInPlace, _confirm, _changes bool, _pfilePath, _outFile string, _puppetfile *ast.Puppetfile) {
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
