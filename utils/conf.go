package utils

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
)

const (
	KEY_ORG_ID       string = "orgId"
	KEY_TOKEN        string = "authToken"
	KEY_URI          string = "uri"
	KEY_DEBUG        string = "debug_mode"
	KEY_FORMAT       string = "format"
	CONFIG_FILE_NAME string = ".anypoint-cli"
)

func WriteConfig() {
	Debug(func() {
		maps := viper.AllSettings()
		fmt.Printf("All settings %s\n", maps)
	})

	fileName := viper.ConfigFileUsed()

	//	fmt.Printf("Config fileName used: %s\n", fileName)
	if fileName == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if home[len(home)-1] != os.PathSeparator {
			home += string(os.PathSeparator)
		}

		fileName = home + CONFIG_FILE_NAME + ".json"
	}

	settings := viper.AllSettings()
	//DO NOT PERSIST THE DEBUG SETTING
	delete(settings, KEY_DEBUG)

	fileContent, err := json.MarshalIndent(settings, " ", "\t")
	if err != nil {
		log.Fatalf("Error while saving configuration file %s : %s", fileName, err)
	}

	ioutil.WriteFile(fileName, fileContent, 755)

}
