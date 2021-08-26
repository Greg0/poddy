package main

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	cmd := exec.Command("kubectl", "get", "pods", "--template", "{{range .items}}{{.metadata.name}}{{\";;\"}}{{end}}")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	options := strings.Split(out.String(), ";;")
	selectedPods := []string{}
	podPrompt := &survey.MultiSelect{
		Message: "Select pods to execute command:",
		Options: options,
	}
	err = survey.AskOne(podPrompt, &selectedPods)
	if err != nil {
		return
	}

	cmdType := ""
	cmdTypePrompt := &survey.Select{
		Message: "Choose command type:",
		Options: []string{"logs"},
	}
	err = survey.AskOne(cmdTypePrompt, &cmdType)

	if err != nil {
		return
	}

	if cmdType == "logs" {
		dirToSave := "logs"
		prompt := &survey.Input{
			Message: "Choose dir to save logs:",
			Suggest: ListDirectories,
		}
		survey.AskOne(prompt, &dirToSave)

		for _, podName := range selectedPods {
			logFile := dirToSave + "/" + podName
			cmd = exec.Command("kubectl", "logs", podName)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				return
			}
			err = ioutil.WriteFile(logFile, out.Bytes(), 0644)
			check(err)

			fmt.Println("[" + podName + "] - logs saved to " + logFile)
		}
	}
}

func ListDirectories(toComplete string) []string {
	var directories []string
	files, _ := filepath.Glob(toComplete + "*")

	for _, path := range files {
		fileInfo, _ := os.Stat(path)
		if fileInfo.IsDir() {
			directories = append(directories, path)
		}
	}

	return directories
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
