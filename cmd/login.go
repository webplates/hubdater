package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/mitchellh/go-homedir"
)

func init() {
	RootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Docker Hub",

	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet("Username") == false {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter Username: ")
			username, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			viper.Set("Username", strings.TrimSpace(username))
		}

		if viper.IsSet("Password") == false {
			fmt.Print("Enter Password: ")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				panic(err)
			}

			fmt.Println("")

			viper.Set("Password", strings.TrimSpace(string(bytePassword)))
		}

		token, err := HubClient.Login(viper.GetString("Username"), viper.GetString("Password"))
		if err != nil {
			panic(err)
		}

		path, _ := homedir.Expand("~/.hubdater")
		if _, err := os.Stat(path); os.IsNotExist(err) {
		    os.Mkdir(path, 0755)
		}
		path = filepath.Join(path, "token")

		err = ioutil.WriteFile(path, []byte(token), 0644)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Logged in with username \"%s\"!\n", viper.GetString("Username"))
	},
}
