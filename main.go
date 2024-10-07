package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

func getAllRetendFiles(directoryString string) []fs.DirEntry {
    directoryContents, err := os.ReadDir(directoryString)
    output := []fs.DirEntry{}
    if err != nil {
        log.Fatal(err)
    }
    for _, file := range directoryContents {
        if file.IsDir() {
            continue
        }
        if strings.Contains(file.Name(), ".retend") || strings.Contains(file.Name(), ".schedule") {
            output = append(output, file)
        }
    }
    
    if err != nil {
        fmt.Println(err)
    }
    return output
}

func main() {
    retendDirectoryPath := "/home/tjmisko/Tools/Retend/RetendExport"
    retend_files := getAllRetendFiles(retendDirectoryPath)
    for _, file := range retend_files {
        fmt.Println("Reading file:", file.Name())
        fileContents, err := os.Open(retendDirectoryPath + "/" + file.Name())
        if err != nil {
            log.Fatal(err)
        }
        defer fileContents.Close()

    }
}
