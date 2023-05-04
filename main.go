package main

import (
    "errors"
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

func modifiedLocalRepository() (bool, error) {
    gitStatus, err := exec.Command("git", "status", "--porcelain=v2").Output()
    if err != nil {
        return false, err
    }
    return gitStatus != nil, nil
}

func exportCommand() error {
    modified, err := modifiedLocalRepository()
    if err != nil {
        fmt.Println("cannot get git status")
        return err
    } else if modified {
        fmt.Print("Local Change detected! continue to delete .git folder? [y/N]:")
        var input string
        fmt.Scan(&input)
        if input != "y" {
            return errors.New("Aborted.")
        }
    }

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
        fmt.Println("cannot remove .git folder")
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
    var command string
    if len(arg) == 0 {
        //Auto mode
        if existTransporterFile() {
            fmt.Print("Auto Mode: execute `Import` ok? [y/N]: ")
            var input string
            fmt.Scan(&input)
            if input == "y" {
                command = "import"
            }
        } else {
            fmt.Print("Auto Mode: execute `Export` ok? [y/N]: ")
            var input string
            fmt.Scan(&input)
            if input == "y" {
                command = "export"
            }
        }
    } else {
        switch arg[0] {
        case "export":
            command = "export"
        case "import":
            command = "import"
        default:
        }
    }

    switch command {
    case "import":
        if err := importCommand(); err != nil {
            fmt.Println(err)
            fmt.Println("Process Abort")
            os.Exit(1)
        }
    case "export":
        if err := exportCommand(); err != nil {
            fmt.Println(err)
            fmt.Println("Process Abort")
            os.Exit(1)
        }
    default:
        fmt.Println("Invalid Command: `github-transporter {export/import}`")
        os.Exit(1)
    }
}
