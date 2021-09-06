package main



import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"log"
	"strings"
	"bufio"
)



const (
	InfoColor    = "\033[1;34m"
		NormalColor = "\033[32m"
	//NoticeColor  = "\033[1;36m%s\033[0m"
	//WarningColor = "\033[1;33m%s\033[0m"
        ErrorColor   = "\033[1;31m"
    	ColorReset = "\033[0m"
	// DebugColor   = "\033[0;36m%s\033[0m"
)
func isRoot() bool {
    currentUser, err := user.Current()
    if err != nil {
        log.Fatalf("[isRoot] Unable to get current user: %s", err)
    }
    return currentUser.Username == "root"
}

func main() {
	if _, err := os.Stat("/etc/gpac.conf"); os.IsNotExist(err) {
	fmt.Println("/etc/gpac.conf does not exist")
	os.Exit(1)
}

	arguments()
	os.Exit(0)
}



func build(pkg string) {
	// root-check
	if isRoot() == false{
		fmt.Println("rn me as root")
		os.Exit(127)
	}

	fmt.Println(InfoColor+"installing package: " +pkg)





    file, err := os.Open("/etc/gpac.conf")
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err = file.Close(); err != nil {
            // log.Fatal(err)
        }
    }()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {             // internally, it advances token based on sperator

		var repo string = scanner.Text()
		var package_location string = repo + pkg

		fmt.Println(NormalColor,"✅", ColorReset, " Using repo at: " + repo)
	if _, err := os.Stat(package_location); os.IsNotExist(err) {
	fmt.Println(ErrorColor, "❌ Package" + pkg+ " does not exist", ColorReset)
		// os.Exit(1)
}
	if _, err := os.Stat(package_location); !os.IsNotExist(err) {
	fmt.Println(NormalColor,"✅", ColorReset, " Package " + pkg + " found")
}


    }


	var tmpdir string = "/tmp/"
	tmpdir = tmpdir + pkg
	fmt.Println(NormalColor,"✅ ", ColorReset, "Creating tmpdir: " + tmpdir)


	if tmpdir != "/" && strings.Contains(tmpdir, "/tmp") == true {
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


    file, err = os.Open("/etc/gpac.conf")
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err = file.Close(); err != nil {
            log.Fatal(err)
        }
    }()

    scanner = bufio.NewScanner(file)

    for scanner.Scan() {             // internally, it advances token based on sperator
        fmt.Println(scanner.Text())  // token in unicode-char

    bcmd := exec.Command("cp","-r",scanner.Text() + pkg + "/" + "build",tmpdir)
	bcmd.Stdout = os.Stdout
    bcmd.Stderr = os.Stderr
    if err := bcmd.Run(); err != nil {
        log.Fatal(err)
    }


    file, err := os.Open(scanner.Text() + pkg + "/" + "url")
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        if err = file.Close(); err != nil {
            log.Fatal(err)
        }
    }()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {             // internally, it advances token based on sperator
		fmt.Println(tmpdir)
		fmt.Println(scanner.Text())
		var url string = scanner.Text()
		// "curl",  url, ">",tmpdir + "/", os.Args[2]
		     f, err := os.Create("/tmp/clurl.sh")
    if err != nil {
        fmt.Println(err)
        return
    }

    l, err := f.WriteString("curl " + "-LG " +  url+ " > "+tmpdir + "/"+ pkg + ".tar.gz")
    if err != nil {
        // fmt.Println(err)
        f.Close()
        return
    }
	l = l
		err = f.Close()
    if err != nil {
        // fmt.Println(err)
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

	// build-cmd

    bcmd2 := exec.Command("sh", tmpdir + "/" + "build")
	bcmd2.Stdout = os.Stdout
    bcmd2.Stderr = os.Stderr
    if err := bcmd2.Run(); err != nil {
        log.Fatal(err)
    }
	fmt.Println(NormalColor,"✅ ", "Package " + pkg+ " installed", ColorReset)
}
func help(){
	fmt.Println("+-----------+\n" +
				"| gpac help |\n" +
				"+-----------+" )
	os.Exit(0)
}

func  arguments() {

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
