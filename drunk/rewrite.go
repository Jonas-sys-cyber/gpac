package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

const (
	InfoColor    = "\033[1;34m"
	SuccessColor = "\033[32m"
	ErrorColor   = "\033[1;31m"
	ColorReset   = "\033[0m"
)

// error messages
func errorPrint(msg string, exitCode int) {
	fmt.Println(ErrorColor, msg, ColorReset)
	os.Exit(exitCode)

}
func main() {
	if _, err := os.Stat("/etc/gpac.gconf"); os.IsNotExist(err) {
		errorPrint("Error: /etc/gpac.gconf not found. Just install the shit right, man!", 127)
	}
	if checkargs() {
		arguments()
	} else {
		help()
	}

	os.Exit(0)
}

func gconf(gconfs string, keyword string) string {

	for _, line := range strings.Split(strings.TrimRight(gconfs, "\n"), "\n") {

		if string(line[0:len(keyword)+1]) == ":"+keyword {
			text := line[:len(line)-1]
			text = text[len(keyword)+2:]
			return text
		} else {
			return ""
		}
	}
	panic("should never happen")

}

// build function
func build(pkg string) {
	// get repo
	data, err := ioutil.ReadFile("/etc/gpac.gconf")
	if err != nil {
		fmt.Println("File reading error. Please don´t try it again", err)
		return
	}
	var gpacGConfText string = string(data)
	repo := ""
	for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
		if gconf(string(line), "repoPath") != "" {
			repo = gconf(string(line), "repoPath")
		}
	}
	repoUrl := ""
	for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
		if gconf(string(line), "repo") != "" {
			repoUrl = gconf(string(line), "repoUrl")
		}
	}
	pkgUrl := repoUrl + pkg + "tar.gz"
	fmt.Println(pkgUrl)
	fmt.Println(repoUrl)
	fmt.Println(repo)
	fmt.Println(pkg)
}

// check if arguments are given
func checkargs() bool {
	return len(os.Args) > 1
}

// root check
func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

// help
func help() {
	fmt.Println("Real programmers don´t need help!")

}
func arguments() {

	if os.Args[1] == "help" || os.Args[1] == "h" {

		help()
	}

	for i, arg := range os.Args {

		if i == 0 || i == 1 {

		} else {

			if os.Args[1] == "build" || os.Args[1] == "b" {

				build(arg)

			}
		}
	}
}
