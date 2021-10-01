package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func Print(msg string) {
	fmt.Println(SuccessColor, "✓", ColorReset, msg)
}

func Printl(msg string) {
	fmt.Println(InfoColor, "⌛", ColorReset, msg)
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

func addToList(pkg, list string) {
	file, err := os.OpenFile(list, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(pkg + "\n"); err != nil {
		log.Fatal(err)
	}
}
func gconf(gconfs, keyword string) string {

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

func src(pkg string) {
	pkgName := pkg
	data, err := ioutil.ReadFile("/etc/gpac.gconf")
	if err != nil {
		fmt.Println("File reading error. Please don´t try it again", err)
		return
	}
	var gpacGConfText string = string(data)

	repoUrl := ""
	if strings.Contains(pkg, "/") {

		trimmed := strings.Split(pkg, "/")
		//fmt.Println(trimmed[0])
		for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
			if gconf(string(line), trimmed[0]) != "" {
				repoUrl = gconf(string(line), trimmed[0])
				pkgName = trimmed[1]

				fmt.Println("using 3rd party reop: " + repoUrl)
			}
		}
	} else {

		for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
			if gconf(string(line), "repoUrl") != "" {
				repoUrl = gconf(string(line), "repoUrl")
			}
		}
		fmt.Println("using reop: " + repoUrl)

	}
	srcDir := ""
	for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
		if gconf(string(line), "srcDir") != "" {
			srcDir = gconf(string(line), "srcDir")
		}
	}
	buildOutput := ""
	for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
		if gconf(string(line), "buildOutput") != "" {
			buildOutput = gconf(string(line), "buildOutput")
		}
	}
	pkgUrl := repoUrl + pkgName + ".tar.gz"
	Printl("Downloading Source Code...")
	err = download(srcDir+pkgName+".tar.gz", pkgUrl)
	if err != nil {
		panic(err)
	}
	os.MkdirAll(srcDir, os.ModePerm)
	Print("Source Code downloaded")
	Printl("extracting Source Code...")
	cmd := exec.Command("tar", "--extracr", "-f", srcDir+pkgName+"tar.gz")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if buildOutput == "1" {
		cmd.Stdout = os.Stdout
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	Print("Souce Code extracted")
}

// build function
func build(pkg string) {
	// get repo
	pkgName := pkg
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
	if strings.Contains(pkg, "/") {

		trimmed := strings.Split(pkg, "/")
		//fmt.Println(trimmed[0])
		for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
			if gconf(string(line), trimmed[0]) != "" {
				repoUrl = gconf(string(line), trimmed[0])
				pkgName = trimmed[1]

				fmt.Println("using 3rd party reop: " + repoUrl)
			}
		}
	} else {

		for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
			if gconf(string(line), "repoUrl") != "" {
				repoUrl = gconf(string(line), "repoUrl")
			}
		}
		fmt.Println("using reop: " + repoUrl)
	}

	pkgList := ""
	for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
		if gconf(string(line), "pkgList") != "" {
			pkgList = gconf(string(line), "pkgList")
		}
	}

	tmpDir := ""
	for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
		if gconf(string(line), "tmpDir") != "" {
			tmpDir = gconf(string(line), "tmpDir") + pkgName + "/"
		}
	}
	buildOutput := ""
	for _, line := range strings.Split(strings.TrimRight(gpacGConfText, "\n"), "\n") {
		if gconf(string(line), "buildOutput") != "" {
			buildOutput = gconf(string(line), "buildOutput")
		}
	}
	// create tmp dir
	os.MkdirAll(tmpDir, os.ModePerm)
	// get gconf file for the build script
	gconfFile := repo + pkg + ".gconf"
	// fmt.Println(gconfFile)
	data, err = ioutil.ReadFile(gconfFile)
	if err != nil {
		fmt.Println("File reading error. Please don´t try it again", err)
		return
	}
	buildCommand := ""
	var gconfText string = string(data)
	for _, line := range strings.Split(strings.TrimRight(gconfText, "\n"), "\n") {
		if gconf(string(line), "build") != "" {
			buildCommand = gconf(string(line), "build")
		}
	}

	pkgMsg := ""
	for _, line := range strings.Split(strings.TrimRight(gconfText, "\n"), "\n") {
		if gconf(string(line), "msg") != "" {
			pkgMsg = gconf(string(line), "msg")
		}
	}
	// write the build script
	buildCommand = "cd " + tmpDir + "; " + "tar --extract -f " + pkgName + ".tar.gz" + ";" + "cd " + pkgName + ";" + buildCommand
	// fmt.Println(buildCommand)
	f, err := os.Create(tmpDir + "build")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	//fmt.Println(buildCommand)
	_, err2 := f.WriteString(buildCommand)
	if err2 != nil {
		log.Fatal(err2)
	}

	// fmt.Println("done")
	// fmt.Println(gconfFile)
	//fmt.Println(tmpDir + pkgName + ".tar.gz")
	pkgUrl := repoUrl + pkgName + ".tar.gz"

	err = download(tmpDir+pkgName+".tar.gz", pkgUrl)
	if err != nil {
		panic(err)
	}
	Print("Downloaded tar ball from: " + pkgUrl)
	//fmt.Println(pkgUrl)
	//fmt.Println(repoUrl)
	//fmt.Println(repo)
	//fmt.Println(pkg)
	Printl("running build script this can take a while")

	cmd := exec.Command("sh", tmpDir+"buil") // this breaks gpac use build instead of buil
	// i need it for tests
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if buildOutput == "1" {
		cmd.Stdout = os.Stdout
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	if pkgMsg != "" {
		fmt.Println("pkgMsg: " + pkgMsg)
	} else {
		fmt.Println("no pkgMsg was provided")
	}
	Print("package installed")
	//adding package to the package list
	addToList(pkg, pkgList)

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

			if os.Args[1] == "build" || os.Args[1] == "b" || os.Args[1] == "install" {

				build(arg)

			} else if os.Args[1] == "source" || os.Args[1] == "src" {
				src(arg)
			}
		}
	}
}
