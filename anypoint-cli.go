package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"strings"
	"github.com/tech-nico/anypoint-cli/rest"
)

var hostname = flag.String("hostname", "", "MuleSoft Anypoint platform hostname")
var port = flag.Int("port", 443, "Port used to connect to the Anypoint platform. Default to 443 (HTTPS)")
var help = flag.Bool("help", false, "Shows how to use anypoint-cli")

const HOSTNAME_KEY string = "ANYPOINT_CLI_HOSTNAME"

func main() {
	flag.Parse()

	if *hostname == ""{
		fmt.Println("--hostname is a mandatory parameter")
		os.Exit(1)
	}

	login(*hostname)
}

func login (uri string){
	username, password := credentials()

	fmt.Printf("\nUsername: %s\n", username)
	fmt.Printf("Pwd: %s\n", password)

	//login(username, password)
	rest.NewAuth(uri, username, password)

}

func credentials() (string, string){
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username:")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter password:")
	bytepassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	password := string(bytepassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
