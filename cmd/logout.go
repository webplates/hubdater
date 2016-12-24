package cmd

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/mitchellh/go-homedir"
)

func init() {
	RootCmd.AddCommand(logoutCmd)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Docker Hub",

	Run: func(cmd *cobra.Command, args []string) {
		path, _ := homedir.Expand("~/.hubdater")
		if _, err := os.Stat(path); os.IsNotExist(err) {
		    os.Mkdir(path, 0755)
		}
		path = filepath.Join(path, "token")

		err := ioutil.WriteFile(path, []byte(""), 0644)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Successfully logged out!\n")
	},
}
