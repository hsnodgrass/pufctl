package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hsnodgrass/pufctl/internal/util"
	"github.com/hsnodgrass/pufctl/pkg/ast"
)

var (
	tag              bool
	data             bool
	module           bool
	metaWriteInPlace bool
	metaOutFile      string

	metaCmd = &cobra.Command{
		Use:   "meta [subcommand]",
		Short: "Commands for working with Puppetfile metadata",
		Long: `meta gives you several subcommands to work with
Puppetfile metadata. See the help of each command to see what
they do.`,
	}

	searchCmd = &cobra.Command{
		Use:   "search [meta (tag or data)]",
		Short: "Search for modules by their metadata (tags or data itself).",
		Long: `search allows you to see modules that have
been tagged with a specific meta data, whether it's a tag
or the actual string value of the data.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			trueCount := 0
			for _, b := range []bool{tag, data, module} {
				if b {
					trueCount++
				}
			}
			if trueCount > 1 {
				util.ExitErr("You may only use one search option at a time")
			}
			if trueCount < 1 {
				util.ExitErr("You must use a search option. Valid options are --tag (-t), --data (-d), and --module (-m)")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, a := range args {
				_ = metaSearch(a)
			}
		},
	}

	addCmd = &cobra.Command{
		Use:   "add [module name] [tag] [(optional) data]",
		Short: "Add metadata to a module",
		Long:  `add allows you to add metadata to a Puppetfile programatically.`,
		Run: func(cmd *cobra.Command, args []string) {
			puppetfile, err := util.ParseFile(viper.GetString("puppetfile"))
			if err != nil {
				util.ExitErr(err)
			}
			metaAdd(puppetfile, args[0], args[1], args[2])
			if metaWriteInPlace {
				if err = metaInPlace(viper.GetString("puppetfile"), puppetfile.Sprint()); err != nil {
					util.ExitErr(err)
				}
			}
			if metaOutFile != "" {
				fmt.Printf("Creating new file at %s\n", metaOutFile)
				err := metaWriteFile(&metaOutFile, puppetfile.Sprint())
				if err != nil {
					util.ExitErr(err)
				}
				// if err = ioutil.WriteFile(sortOutFile, []byte(puppetfile.Sprint()), 0755); err != nil {
				// 	util.ExitErr(err)
				// }
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(metaCmd)
	metaCmd.AddCommand(addCmd)
	metaCmd.AddCommand(searchCmd)
	addCmd.Flags().BoolVarP(&metaWriteInPlace, "write-in-place", "w", false, "Overwrite the Puppetfile with the changed version")
	addCmd.Flags().StringVarP(&metaOutFile, "out-file", "o", "", "Write the changed Puppetfile to a new file")
	searchCmd.Flags().BoolVarP(&tag, "tag", "t", false, "Search by metadata tag ( # @<tag>: ...)")
	searchCmd.Flags().BoolVarP(&data, "data", "d", false, "Search by metadata data ( # @tag: <data>)")
	searchCmd.Flags().BoolVarP(&module, "module", "m", false, "Search for metadata associated with a module")
}

func metaSearch(input string) error {
	puppetfile, err := util.ParseFile(viper.GetString("puppetfile"))
	if err != nil {
		util.ExitErr(err)
	}
	if tag {
		tagSearch(puppetfile, input)
	}
	if data {
		dataSearch(puppetfile, input)
	}
	if module {
		moduleSearch(puppetfile, input)
	}
	return nil
}

func tagSearch(puppetfile *ast.Puppetfile, input string) error {
	names := puppetfile.SearchModulesByMetaTag(input)
	if len(names) > 0 {
		for _, p := range names {
			fmt.Printf("Module \"%s\" tagged with @%s\n", p, input)
		}
	} else {
		fmt.Println("No results found")
	}
	return nil
}

func dataSearch(puppetfile *ast.Puppetfile, input string) error {
	names := puppetfile.SearchModulesByMetaData(input)
	if len(names) > 0 {
		for _, p := range names {
			fmt.Printf("Module \"%s\" has metadata \"%s\"\n", p, input)
		}
	} else {
		fmt.Println("No results found")
	}
	return nil
}

func moduleSearch(puppetfile *ast.Puppetfile, input string) error {
	out := puppetfile.SearchMetaByModuleName(input)
	if out != "" {
		fmt.Println(out)
	} else {
		fmt.Printf("Module \"%s\" has no metadata associated with it\n", input)
	}
	return nil
}

func metaAdd(puppetfile *ast.Puppetfile, name string, tag string, data string) error {
	_ = puppetfile.AddModuleMetadata(name, tag, data)
	fmt.Println(puppetfile.Sprint())
	return nil
}

func metaInPlace(path string, data string) error {
	original, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, []byte(data), 0755)
	if err != nil {
		e := ioutil.WriteFile(path, []byte(original), 0755)
		if e != nil {
			fmt.Println("FATAL: Error occured and original file could not be restored!")
			panic(e)
		}
		fmt.Println("ERROR: Error occured, but original file was restored!")
		return err
	}
	return nil
}

func metaWriteFile(path *string, content string) error {
	fmt.Printf("Creating new file at %s\n", *path)
	f, err := os.Create(*path)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Printf("Writing sorted Puppetfile to new file\n")
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	fmt.Printf("New Puppetfile successfully created at %s", *path)
	return nil
}
