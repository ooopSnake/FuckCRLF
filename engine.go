package main

import (
	"flag"
	"os"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"strings"
)

const (
	kExitSuccess  = iota
	kExitNoExt
	kExitWTF
	kExitFileStat
)

func main() {
	flagParser := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cmdFileExt := flagParser.String("ext", "",
		`(required) set target file ext . eg lua go without dot.`)
	cmdWorkDir := flagParser.String("path", "",
		`(optional) path to handle, default is currentdir`)
	flagParser.Parse(os.Args[1:])
	targetFileExit := ""
	targetWorkDir := ""
	if len(*cmdFileExt) == 0 {
		flagParser.Usage()
		os.Exit(kExitNoExt)
	} else {
		targetFileExit = fmt.Sprintf(".%s", *cmdFileExt)
	}
	if len(*cmdWorkDir) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("WTF ? err:", err.Error())
			os.Exit(kExitWTF)
		}
		targetWorkDir = wd
	} else {
		fInfo, err := os.Stat(*cmdWorkDir)
		if err != nil {
			fmt.Println("file stat error ", err.Error())
			os.Exit(kExitFileStat)
		} else {
			targetWorkDir = fInfo.Name()
		}
	}
	relParent := fmt.Sprintf(".%c", filepath.Separator)
	filepath.Walk(targetWorkDir, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if info.IsDir() {
			return nil
		}
		if ext != targetFileExit {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("err , while reading ", path, " skip")
			return nil
		}
		newStr := strings.Replace(string(content), "\r\n", "\n", -1)
		err = ioutil.WriteFile(path, []byte(newStr), info.Mode())
		if err != nil {
			fmt.Println("err ,while writing ", path, " msg is ", err.Error())
			return nil
		}
		relStr, err := filepath.Rel(targetWorkDir, path)
		if err != nil {
			fmt.Println("can not get rel path msg is ", err.Error())
			return nil
		}
		fmt.Println(fmt.Sprintf("%s%s processed", relParent, relStr))
		return nil
	})
	os.Exit(kExitSuccess)
}
