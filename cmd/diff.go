package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hsnodgrass/pufctl/internal/validators"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hsnodgrass/pufctl/internal/helpers"
	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/internal/uitext"
	"github.com/hsnodgrass/pufctl/pkg/puppetfileparser/ast"
)

var (
	diffMeta  bool
	branchOne string
	branchTwo string

	diffCmd = &cobra.Command{
		Use:   uitext.DiffUse,
		Short: uitext.DiffShort,
		Long:  uitext.DiffLong,
		Args:  cobra.MaximumNArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				logging.Errorln("At least one Puppetfile path is required")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var pOnePath string
			var pTwoPath string
			var pOne *ast.Puppetfile
			var pTwo *ast.Puppetfile
			var err error
			var outString string
			var exitCode int
			var showPathOne string
			var showPathTwo string
			diffParseOptsOne := helpers.ParseOptions{
				Verbose:  helpers.MaxBools(viper.GetBool("always.verbose"), verbose),
				GitRef:   branchOne,
				SSHKey:   viper.GetString("auth.ssh_key"),
				Username: viper.GetString("auth.username"),
				Password: viper.GetString("auth.password"),
				Token:    viper.GetString("auth.token"),
			}
			diffParseOptsTwo := helpers.ParseOptions{
				Verbose:  helpers.MaxBools(viper.GetBool("always.verbose"), verbose),
				GitRef:   branchTwo,
				SSHKey:   viper.GetString("auth.ssh_key"),
				Username: viper.GetString("auth.username"),
				Password: viper.GetString("auth.password"),
				Token:    viper.GetString("auth.token"),
			}
			_show := helpers.MaxBools(show, viper.GetBool("always.show"))
			if len(args) == 1 {
				pOnePath = viper.GetString("puppetfile")
				pTwoPath = args[0]
				if validators.IsGitURL(pOnePath) || validators.IsGitURLNoSuffix(pOnePath) {
					showPathOne = fmt.Sprintf("%s#%s", pOnePath, diffParseOptsOne.GitRef)
				} else {
					showPathOne = pOnePath
				}
				if validators.IsGitURL(pTwoPath) || validators.IsGitURLNoSuffix(pTwoPath) {
					showPathTwo = fmt.Sprintf("%s#%s", pTwoPath, diffParseOptsTwo.GitRef)
				} else {
					showPathTwo = pTwoPath
				}
				pOne, err = helpers.Parse(pOnePath, diffParseOptsOne)
				if err != nil {
					logging.Errorln(fmt.Errorf("Failed to parse first Puppetfile %s with error: %w", showPathOne, err))
				}
				pTwo, err = helpers.Parse(pTwoPath, diffParseOptsTwo)
				if err != nil {
					logging.Errorf("Failed to parse second Puppetfile %s with error: %v", showPathTwo, err)
				}
			} else {
				pOnePath = args[0]
				pTwoPath = args[1]
				if validators.IsGitURL(pOnePath) || validators.IsGitURLNoSuffix(pOnePath) {
					showPathOne = fmt.Sprintf("%s#%s", pOnePath, diffParseOptsOne.GitRef)
				} else {
					showPathOne = pOnePath
				}
				if validators.IsGitURL(pTwoPath) || validators.IsGitURLNoSuffix(pTwoPath) {
					showPathTwo = fmt.Sprintf("%s#%s", pTwoPath, diffParseOptsTwo.GitRef)
				} else {
					showPathTwo = pTwoPath
				}
				pOne, err = helpers.Parse(pOnePath, diffParseOptsOne)
				if err != nil {
					logging.Errorln(fmt.Errorf("Failed to parse first Puppetfile %s with error: %w", showPathOne, err))
				}
				pTwo, err = helpers.Parse(pTwoPath, diffParseOptsTwo)
				if err != nil {
					logging.Errorln(fmt.Errorf("Failed to parse second Puppetfile %s with error: %w", showPathTwo, err))
				}
			}
			exOne, exTwo, diffResult := diffPuppetfiles(pOne, pTwo, diffMeta)
			if len(exOne) == 0 && len(exTwo) == 0 && len(diffResult) == 0 {
				logging.Debugln("No difference between Puppetfiles")
				outString = "No difference between Puppetfiles"
				exitCode = 0
			}
			exStrParts := []string{fmt.Sprintf("\nExclusive Statements\n%s\n", uitext.DashSep)}
			diffStrParts := []string{fmt.Sprintf("\n%s\nDiff Statements\n", uitext.DashSep)}
			if len(exOne) > 0 {
				exStrParts = append(exStrParts, fmt.Sprintf("Only found in %s\n%s", showPathOne, uitext.DashSep))
				logging.Debugln("Found exclusives in first Puppetfile")
				for _, d := range exOne {
					exStrParts = append(exStrParts, fmt.Sprintf("Module: %s", d.Module.Name))
				}
				if len(exTwo) == 0 {
					exitCode = 3
				}
			}
			if len(exTwo) > 0 {
				exStrParts = append(exStrParts, fmt.Sprintf("\n%s\nOnly found in %s\n%s", uitext.DashSep, showPathTwo, uitext.DashSep))
				logging.Debugln("Found exclusives in second Puppetfile")
				for _, d := range exTwo {
					exStrParts = append(exStrParts, fmt.Sprintf("Module: %s", d.Module.Name))
				}
				if len(exOne) == 0 {
					exitCode = 4
				}
			}
			if len(exOne) > 0 && len(exTwo) > 0 {
				exitCode = 5
			}
			if len(diffResult) > 0 {
				for _, d := range diffResult {
					diffStrParts = append(diffStrParts, fmt.Sprintf("%s\nMODULE: %s\n%s", uitext.DashSep, d.Name, uitext.DashSep))
					diffStrParts = append(diffStrParts, fmt.Sprintf("SOURCE: %s\n%s", showPathOne, d.One))
					diffStrParts = append(diffStrParts, fmt.Sprintf("SOURCE: %s\n%s", showPathTwo, d.Two))
				}
				exitCode = 5
			}
			if len(exStrParts) > 0 {
				outString = strings.Join(exStrParts, "\n")
			}
			if len(diffStrParts) > 0 {
				diffout := strings.Join(diffStrParts, "\n")
				outString = fmt.Sprintf("%s\n%s", outString, diffout)
			}
			outString = fmt.Sprintf("%s\n%s\n", outString, uitext.DashSep)
			if outFile != "" {
				err = helpers.PromptConfirmFile(outFile, outString, confirm)
				if err != nil {
					logging.Errorln("Failed to write output to file with error:", err)
				}
			}
			if _show {
				fmt.Println(outString)
			} else {
				logging.Warnln("If you would like to see the diff, please use the flag --show (-s)")
			}
			os.Exit(exitCode)
		},
	}
)

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolVar(&diffMeta, "diff-meta", false, "Include meta statements (comments, etc.) in diff")
	diffCmd.Flags().StringVar(&branchOne, "branch-one", "production", "The branch for the first Puppetfile (if using a Git source)")
	diffCmd.Flags().StringVar(&branchTwo, "branch-two", "development", "The branch for the second Puppetfile (if using a Git source)")
}

type modDiff struct {
	Name string
	One  string
	Two  string
}

func diffPuppetfiles(pOne, pTwo *ast.Puppetfile, diffMeta bool) ([]*ast.Statement, []*ast.Statement, []modDiff) {
	// diffPuppetfiles returns the difference between two Puppetfiles
	// If diffMeta is true, the diff will contain comments in the diff.
	// The diff does not account for the ordering of the Statements in
	// the Puppethelpers.
	stmtsOne := []*ast.Statement{}
	stmtsTwo := []*ast.Statement{}
	if diffMeta {
		logging.Debugln("Including metadata in diff")
		stmtsOne = append(stmtsOne, pOne.Statements...)
		stmtsTwo = append(stmtsTwo, pTwo.Statements...)
	} else {
		logging.Debugln("Only diffing Module statements")
		for _, s := range pOne.Statements {
			if s.Module != nil {
				stmtsOne = append(stmtsOne, s)
			}
		}
		for _, s := range pTwo.Statements {
			if s.Module != nil {
				stmtsTwo = append(stmtsTwo, s)
			}
		}
	}
	targetOne := map[string]string{}
	targetTwo := map[string]string{}
	for _, s := range stmtsTwo {
		// normalize the name before adding it
		var name string
		var err error
		if s.Module != nil && validators.IsModuleRef(s.Module.Name) {
			name, err = helpers.ModStringToSlug(s.Module.Name)
			if err != nil {
				logging.Errorf("Failed to normalized module name %s: %w", s.Module.Name, err)
			}
		} else if s.Module != nil {
			name = s.Module.Name
		} else {
			name = "metadata"
		}
		logging.Debugf("Adding %s to target two map\n", name)
		targetTwo[name] = s.Sprint()
	}
	for _, s := range stmtsOne {
		// normalize the name before adding it
		var name string
		var err error
		if s.Module != nil && validators.IsModuleRef(s.Module.Name) {
			name, err = helpers.ModStringToSlug(s.Module.Name)
			if err != nil {
				logging.Errorf("Failed to normalized module name %s: %w", s.Module.Name, err)
			}
		} else if s.Module != nil {
			name = s.Module.Name
		} else {
			name = "metadata"
		}
		logging.Debugf("Adding %s to target one map\n", name)
		targetOne[name] = s.Sprint()
	}
	diffResult := []modDiff{}
	exclusiveOne := []*ast.Statement{}
	exclusiveTwo := []*ast.Statement{}
	for _, s := range stmtsOne {
		// normalize the name
		var name string
		var err error
		if s.Module != nil && validators.IsModuleRef(s.Module.Name) {
			name, err = helpers.ModStringToSlug(s.Module.Name)
			if err != nil {
				logging.Errorf("Failed to normalized module name %s: %w", s.Module.Name, err)
			}
		} else if s.Module != nil {
			name = s.Module.Name
		} else {
			name = "metadata"
		}
		val, ok := targetTwo[name]
		if !ok {
			exclusiveOne = append(exclusiveOne, s)
		} else {
			if val != targetOne[name] {
				_diff := modDiff{
					Name: name,
					One:  targetOne[name],
					Two:  val,
				}
				diffResult = append(diffResult, _diff)
			}
		}
	}
	for _, s := range stmtsTwo {
		// normalize the name
		var name string
		var err error
		if s.Module != nil && validators.IsModuleRef(s.Module.Name) {
			name, err = helpers.ModStringToSlug(s.Module.Name)
			if err != nil {
				logging.Errorf("Failed to normalized module name %s: %w", s.Module.Name, err)
			}
		} else if s.Module != nil {
			name = s.Module.Name
		} else {
			name = "metadata"
		}
		if _, ok := targetOne[name]; !ok {
			exclusiveTwo = append(exclusiveTwo, s)
		}
	}
	return exclusiveOne, exclusiveTwo, diffResult
}

// ModuleDiff is used to hold the relevant diff properties of a Puppetfile
// TODO: Actually implement this (more detailed diff reporting)
type ModuleDiff struct {
	Name       string
	Properties []string
	Puppetfile string
}
