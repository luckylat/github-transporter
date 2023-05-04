package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
    "regexp"
    "strings"
)

const filename = ".github-transporter"

//TODO: Restore function

func exportCommand() error {
	//get origin URL
	repositoryName, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		fmt.Println("cannot get git remote file")
		return err
	}

	//make `.github-transporter` file
	dotfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("cannot touch", filename, ".")
		return err
	}
	defer dotfile.Close()

	dotfile.WriteString("remote = " + string(repositoryName))

	if err := os.RemoveAll(".git"); err != nil {
		fmt.Println("cannot remove .git file")
		return err
	}
	return nil
}

func importCommand() error {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("cannot open", filename, ".")
		return err
	}
    formattedFileData := strings.ReplaceAll(string(fileData), "\n", "")

    reg := regexp.MustCompile(`remote = `)
    repositoryName := reg.ReplaceAllString(formattedFileData, "")

    commandStr := "git init; git remote add origin " + repositoryName + "; git fetch;"
    gitCommand := exec.Command("/bin/sh", "-c", commandStr)

    if err := gitCommand.Run(); err != nil {
        fmt.Println("cannot add remote URL:", repositoryName, ".")
        return err
    }

    if err := os.Remove(".github-transporter"); err != nil {
        fmt.Println("cannnot remote .github-transporter file")
        return err
    }

	return nil
}

func existTransporterFile() bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func main() {
	flag.Parse()
	arg := flag.Args()
	if len(arg) == 0 {
		//Auto mode
		if existTransporterFile() {
			fmt.Println("Auto Mode: execute `Import` ok? [y/N]: ")
		} else {
			fmt.Println("Auto Mode: execute `Export` ok? [y/N]: ")
		}
	} else {
		switch arg[0] {
		case "export":
            if err := exportCommand(); err != nil {
                fmt.Println(err)
                fmt.Println("Process Abort")
                os.Exit(1)
            }
		case "import":
            if err := importCommand(); err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
		default:
			fmt.Println("Invalid Command: `github-transporter {export/import}`")
			return
		}
	}
}
