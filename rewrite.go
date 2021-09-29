package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
		errorPrint("Error: /etc/gpac.gconf not found", 127)
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

func download(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// build function
func build(pkg string) {
	// get repo
	data, err := ioutil.ReadFile("/etc/gpac.gconf")
	if err != nil {
		fmt.Println("File reading error", err)
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
	tmpDir := ""
	for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
		if gconf(string(line), "tmpDir") != "" {
			tmpDir = gconf(string(line), "tmpDir") + pkg + "/"
		}
	}
	// create tmp dir
	os.MkdirAll(tmpDir, os.ModePerm)
	// get gconf file
	fmt.Println(tmpDir + pkg + ".tar.gz")
	pkgUrl := repoUrl + pkg + ".tar.gz"
	err = download(tmpDir+pkg+".tar.gz", pkgUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded tar ball from: " + pkgUrl)
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
	fmt.Println("real programmers don't need help")

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
