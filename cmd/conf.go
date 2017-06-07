package cmd

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"io/ioutil"
	"encoding/json"
	"log"
)
const (
	KEY_ORG_ID       string = "orgId"
	KEY_TOKEN        string = "authToken"
	KEY_URI          string = "uri"
	CONFIG_FILE_NAME string = ".anypoint-cli"
)

func WriteConfig() {
	maps := viper.AllSettings()
	fmt.Printf("All settings %s\n", maps)
	fileName := viper.ConfigFileUsed()
	fmt.Printf("Config fileName used: %s\n", fileName)
	if fileName == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if home[len(home)-1] != os.PathSeparator {
			home += string(os.PathSeparator)
		}

		fileName := home + CONFIG_FILE_NAME + ".json"
		fileContent, err := json.MarshalIndent(viper.AllSettings(), " ", "\t")
		if err != nil {
			log.Printf("Error while saving configuration file %s : %s", fileName, err)
		}

		ioutil.WriteFile(fileName, fileContent, 755)
	}

}