package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// clors
const (
	InfoColor   = "\033[1;34m"
	SucessColor = "\033[32m"
	ErrorColor  = "\033[1;31m"
	ColorReset  = "\033[0m"
)

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

func main() {
	if _, err := os.Stat("/etc/gpac.gconf"); os.IsNotExist(err) {
		fmt.Println(ErrorColor + "Error: /etc/gpac.gconf not found" + SucessColor)
		os.Exit(1)
	}
	if checkargs() {
		arguments()
	} else {
		help()
	}

	os.Exit(0)
}

// check if arguments are given
func checkargs() bool {
	return len(os.Args) > 1
}

// root-check func
func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func build(pkg string) {
	// root-check
	if !isRoot() {
		fmt.Println(ErrorColor+"Hey bitch! Run me as root!", ColorReset)
		os.Exit(127)
	}

	fmt.Println(InfoColor + "installing package: " + pkg)

	file, err := os.Open("/etc/gpac.gconf")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {

	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() { // internally, it advances token based on sperator
		if gconf(scanner.Text(), "repo") != "" {
			var repo string = gconf(scanner.Text(), "repo")

			var package_location string = repo + pkg

			fmt.Println(SucessColor, "✅", ColorReset, " Using repo at: "+repo)
			if _, err := os.Stat(package_location); os.IsNotExist(err) {
				fmt.Println(ErrorColor, "❌ Package "+pkg+" not found. Don't try it again!", ColorReset)
				os.Exit(1)
			}
			if _, err := os.Stat(package_location); !os.IsNotExist(err) {
				fmt.Println(SucessColor, "✅", ColorReset, " Package "+pkg+" found. But EGAL!")
			}
		}
	}

	var tmpdir string = "/tmp/"
	tmpdir = tmpdir + pkg
	fmt.Println(SucessColor, "✅ ", ColorReset, "Creating tmpdir: "+tmpdir)

	if tmpdir != "/" && strings.Contains(tmpdir, "/tmp") {
		tcmd2 := exec.Command("rm", "-rf", tmpdir)
		tcmd2.Stdout = os.Stdout
		tcmd2.Stderr = os.Stderr
		if err := tcmd2.Run(); err != nil {
			log.Fatal(err)
		}
	}

	// tmp-cmd
	tcmd := exec.Command("mkdir", tmpdir)
	tcmd.Stdout = os.Stdout
	tcmd.Stderr = os.Stderr
	if err := tcmd.Run(); err != nil {
		log.Fatal(err)
	}

	// get repo path
	file, err = os.Open("/etc/gpac.gconf")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner = bufio.NewScanner(file)

	for scanner.Scan() { // internally, it advances token based on sperator
		if gconf(scanner.Text(), "repo") != "" {
			bcmd := exec.Command("cp", "-r", gconf(scanner.Text(), "repo")+pkg+"/"+"build", tmpdir)
			bcmd.Stdout = os.Stdout
			bcmd.Stderr = os.Stderr
			if err := bcmd.Run(); err != nil {
				log.Fatal(err)
			}

			file, err := os.Open(gconf(scanner.Text(), "repo") + pkg + "/" + "url")
			if err != nil {
				log.Fatal(err)
			}

			defer func() {
				if err = file.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			scanner := bufio.NewScanner(file)

			for scanner.Scan() { // internally, it advances token based on sperator

				fmt.Println(scanner.Text())
				var url string = scanner.Text()
				// "curl",  url, ">",tmpdir + "/", os.Args[2]
				f, err := os.Create("/tmp/clurl.sh")
				if err != nil {
					fmt.Println(err)
					return
				}

				l, err := f.WriteString("curl " + "-LG " + url + " > " + tmpdir + "/" + pkg + ".tar.gz")
				if err != nil {

					f.Close()
					return
				}
				l = l
				err = f.Close()
				if err != nil {

					return
				}
				ccmd := exec.Command("sh", "/tmp/clurl.sh")
				ccmd.Stdout = os.Stdout
				ccmd.Stderr = os.Stderr
				if err := ccmd.Run(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	// build-cmd

	bcmd2 := exec.Command("sh", tmpdir+"/"+"build")
	bcmd2.Stdout = os.Stdout
	bcmd2.Stderr = os.Stderr
	if err := bcmd2.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(SucessColor, "✅ ", "Package "+pkg+" installed. So fuck off and have fun with your new bloat.", ColorReset)
}

func help() {
	fmt.Println("+-----------+\n" +
		"| gpac help |\n" +
		"+-----------+\n" +
		"Install: gpac b packagename")
	os.Exit(0)
}
func create(pkgname string) {
	file1, err := os.Open("/etc/gpac.gconf")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file1.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner1 := bufio.NewScanner(file1)
	gconf_file := ""
	rdir := ""
	for scanner1.Scan() {
		if gconf(scanner1.Text(), "grepo") != "" {
			gconf_file = gconf(scanner1.Text(), "grepo") + pkgname + ".gconf"
		}
		if gconf(scanner1.Text(), "repo") != "" {
			rdir = gconf(scanner1.Text(), "repo") + pkgname + "/"
		}
	}
	fmt.Println(rdir)
	rcmd := exec.Command("mkdir", "-p", rdir)
	rcmd.Stdout = os.Stdout
	rcmd.Stderr = os.Stderr
	if err := rcmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(gconf_file)
	file, err := os.Open(gconf_file)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		if gconf(scanner.Text(), "build") != "" {
			fmt.Println("package found")

			f, err := os.Create(rdir + "build")
			if err != nil {
				fmt.Println(err)
				return
			}

			l, err := f.WriteString(gconf(scanner.Text(), "build"))
			if err != nil {

				f.Close()
				return
			}
			l = l
		}
		if gconf(scanner.Text(), "url") != "" {

			f, err := os.Create(rdir + "url")
			if err != nil {
				fmt.Println(err)
				return
			}

			l, err := f.WriteString(gconf(scanner.Text(), "url"))
			if err != nil {

				f.Close()
				return
			}
			l = l
		}
	}
	build(pkgname)
}

func arguments() {

	if os.Args[1] == "help" || os.Args[1] == "h" {

		help()
	}

	for i, arg := range os.Args {

		if i == 0 || i == 1 {

		} else {

			if os.Args[1] == "build" || os.Args[1] == "b" {

				create(arg)

			}
		}
	}
}
