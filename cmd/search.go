package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hsnodgrass/pufctl/internal/logging"
	"github.com/hsnodgrass/pufctl/internal/sources/forgesource"
	"github.com/hsnodgrass/pufctl/internal/uitext"
	"github.com/hsnodgrass/pufctl/pkg/forgeapi"
)

var (
	fsLimit          int
	fsOffset         int
	fsSortBy         string
	fsTag            string
	fsOwner          string
	fsWithTasks      bool
	fsWithPlans      bool
	fsWithPDK        bool
	fsEndorsements   []string
	fsOS             string
	fsPEReq          string
	fsPupReq         string
	fsMinScore       int
	fsModuleGroups   []string
	fsShowDeleted    bool
	fsHideDeprecated bool
	fsOnlyLatest     bool
	fsSlugs          []string
	fsWithHTML       bool
	fsIncludeFields  []string
	fsExcludeFields  []string

	searchCmd = &cobra.Command{
		Use:   uitext.SearchUse,
		Short: uitext.SearchShort,
		Long:  uitext.SearchLong,
	}

	searchForgeCmd = &cobra.Command{
		Use:   uitext.ForgeSearchUse,
		Short: uitext.ForgeSearchShort,
		Long:  uitext.ForgeSearchLong,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_url := viper.GetString("forge.api_url")
			agent := viper.GetString("forge.user_agent")
			opts := newListModulesOpts(args)
			handleMissingFields(&opts)
			finds, err := forgesource.Search(_url, agent, opts)
			if err != nil {
				logging.Errorln("Search failed with error: ", err)
			}
			if len(finds) < 1 {
				logging.Errorln("No modules found, please try different search parameters")
			}
			for _, m := range finds {
				printModShort(m)
			}
			fmt.Println(uitext.DashSep)
			fmt.Printf("Returned %v results\n", len(finds))
			fmt.Println(uitext.DashSep)
		},
	}
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.AddCommand(searchForgeCmd)
	searchForgeCmd.Flags().IntVar(&fsLimit, "limit", 5, "limit the number of search hits returned")
	searchForgeCmd.Flags().IntVar(&fsOffset, "offset", 0, "return only search hits after the offset")
	searchForgeCmd.Flags().StringVarP(&fsSortBy, "sort-by", "b", "rank", "sort search hits by [rank|downloads|latest_release]")
	searchForgeCmd.Flags().StringVarP(&fsTag, "tag", "t", "", "return only modules with specified tag")
	searchForgeCmd.Flags().StringVar(&fsOwner, "owner", "", "return only modules from specified owner")
	searchForgeCmd.Flags().BoolVar(&fsWithTasks, "with-tasks", false, "return only modules with tasks")
	searchForgeCmd.Flags().BoolVar(&fsWithPlans, "with-plans", false, "return only modules with plans")
	searchForgeCmd.Flags().StringSliceVar(&fsEndorsements, "endorsements", []string{}, "return only modules with specified endorsements [supported|approved|partner]")
	searchForgeCmd.Flags().StringVar(&fsOS, "os", "", "return only modules that support specified operating system")
	searchForgeCmd.Flags().StringVar(&fsPEReq, "pe-requirement", "", "return only modules that list a PE version in the specified semver range")
	searchForgeCmd.Flags().StringVar(&fsPupReq, "puppet-requirement", "", "return only modules that list a Puppet version in the specified semver range")
	searchForgeCmd.Flags().IntVar(&fsMinScore, "minimum-score", 0, "return only modules above the specified quality score")
	searchForgeCmd.Flags().StringSliceVar(&fsModuleGroups, "module-groups", []string{}, "return only modules in the specified module groups: [base|pe_only]")
	searchForgeCmd.Flags().BoolVar(&fsShowDeleted, "show-deleted", false, "show deleted modules in search results")
	searchForgeCmd.Flags().BoolVarP(&fsHideDeprecated, "hide-deprecated", "H", false, "hide deprecated modules in search results")
	searchForgeCmd.Flags().BoolVar(&fsOnlyLatest, "only-latest", false, "only consider the latest module version in compat checks")
	searchForgeCmd.Flags().StringSliceVarP(&fsSlugs, "slugs", "S", []string{}, "search for specific modules by their name slugs")
	searchForgeCmd.Flags().BoolVar(&fsWithHTML, "with-html", false, "render markdown to HTML before returning results")
	searchForgeCmd.Flags().StringSliceVarP(&fsIncludeFields, "include-fields", "i", []string{}, "top level optional keys to include in response objects")
	searchForgeCmd.Flags().StringSliceVarP(&fsExcludeFields, "exclude-fields", "e", []string{}, "top level keys to exclude from response objects")
}

func newListModulesOpts(args []string) forgeapi.ListModulesOpts {
	parts := []string{}
	for _, a := range args {
		parts = append(parts, url.QueryEscape(a))
	}
	query := strings.Join(args, " ")
	opts := forgeapi.ListModulesOpts{
		Limit:          fsLimit,
		Offset:         fsOffset,
		SortBy:         fsSortBy,
		Query:          url.QueryEscape(query),
		WithTasks:      fsWithTasks,
		WithPlans:      fsWithPlans,
		WithPDK:        fsWithPDK,
		MinScore:       fsMinScore,
		ShowDeleted:    fsShowDeleted,
		HideDeprecated: fsHideDeprecated,
		OnlyLatest:     fsOnlyLatest,
		WithHTML:       fsWithHTML,
	}
	return opts
}

func handleMissingFields(opts *forgeapi.ListModulesOpts) {
	// Not sure if this function is even necessary with the struct having
	// omitempty field tags.
	noDefaultsStr := []string{fsTag, fsOwner, fsOS, fsPEReq, fsPupReq}
	noDefaultsSlice := [][]string{fsEndorsements, fsModuleGroups, fsSlugs, fsIncludeFields, fsExcludeFields}
	for i, a := range noDefaultsStr {
		if a != "" {
			switch i {
			case 0:
				opts.Tag = a
			case 1:
				opts.Owner = a
			case 2:
				opts.OperatingSystem = a
			case 3:
				opts.PERequirement = a
			case 4:
				opts.PupRequirement = a
			}
		}
	}
	for i, s := range noDefaultsSlice {
		if len(s) > 0 {
			switch i {
			case 0:
				opts.Endorsements = s
			case 1:
				opts.ModuleGroups = s
			case 2:
				opts.Slugs = s
			case 3:
				opts.IncludeFields = s
			case 4:
				opts.ExcludeFields = s
			}
		}
	}
}

func printModShort(m forgeapi.Module) {
	fmt.Println(uitext.StarSep)
	fmt.Println("Name: ", m.Slug)
	fmt.Println("Source: ", m.CurrentRelease.Metadata.Source)
	fmt.Println("Summary: ", m.CurrentRelease.Metadata.Summary)
	fmt.Println("Downloads: ", m.Downloads)
}
