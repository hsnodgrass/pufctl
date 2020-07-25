package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hsnodgrass/pufctl/internal/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	show             bool
	sortWriteInPlace bool
	sortOutFile      string

	sortCmd = &cobra.Command{
		Use:   "sort [(optional) path to Puppetfile]",
		Short: "Sort modules in a Puppetfile",
		Long: `sort allows you to sort modules in a Puppetfile alphabetically by name.
A path to a Puppetfile can be passed as the first argument. If a path isn't passed, 
sort will default to the pufctl configs.`,
		Run: func(cmd *cobra.Command, args []string) {
			puppetfile, err := util.ParseFile(viper.GetString("puppetfile"))
			if err != nil {
				util.ExitErr(err)
			}
			if show {
				fmt.Printf("%s", puppetfile.Sprint())
			}
			if sortWriteInPlace {
				if err = sortInPlace(args[0], puppetfile.Sprint()); err != nil {
					util.ExitErr(err)
				}
			}
			if sortOutFile != "" {
				fmt.Printf("Creating new file at %s\n", sortOutFile)
				err := sortWriteFile(&sortOutFile, puppetfile.Sprint())
				if err != nil {
					util.ExitErr(err)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.Flags().BoolVarP(&show, "show", "s", false, "Print the sorted Puppetfile to the console and don't write to file")
	sortCmd.Flags().BoolVarP(&sortWriteInPlace, "write-in-place", "w", false, "Overwrite the Puppetfile with the sorted version")
	sortCmd.Flags().StringVarP(&sortOutFile, "out-file", "o", "", "Write the sorted Puppetfile to a new file")
}

func sortInPlace(path string, sorted string) error {
	original, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, []byte(sorted), 0755)
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

func sortWriteFile(path *string, content string) error {
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
	fmt.Printf("New Puppet file successfully created at %s", *path)
	return nil
}
