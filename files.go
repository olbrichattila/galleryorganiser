package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

const (
	paralellFileCount = 10
)

type filer interface {
	Split(string, string, bool) error
}

type files struct {
	sourceFolder      string
	destinationFolder string
	overwrite         bool
}

func (f *files) Split(sf, df string, overwrite bool) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, paralellFileCount)

	f.overwrite = overwrite
	f.sourceFolder = sf
	f.destinationFolder = df

	files, err := f.readDir()
	if err != nil {
		return err
	}

	for i, file := range *files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		folderName := f.extension(fileName) + "/" + file.ModTime().Format("2006/01")

		err := f.mkDir(folderName)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("%d File copy started\n", i)
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := f.copyFile(folderName, fileName)
			if err != nil {
				fmt.Println("Copy error", err)
			}

			fmt.Printf("%d File %s/%s copyed to %s/%s/%s\n", i, f.sourceFolder, fileName, f.destinationFolder, folderName, fileName)
		}(i)
	}

	wg.Wait()
	fmt.Println("Done.")

	return nil
}

func (f *files) readDir() (*[]fs.FileInfo, error) {
	d, err := os.Open(f.sourceFolder)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	return &files, nil
}

func (f *files) mkDir(folderName string) error {
	err := os.MkdirAll(f.destinationFolder+"/"+folderName, 0755)
	if err != nil {
		return err

	}

	return nil
}

func (f *files) copyFile(folderName, filaName string) error {
	destFileName := f.destinationFolder + "/" + folderName + "/" + filaName
	if !f.overwrite && f.fileExists(destFileName) {
		return nil
	}
	sourceFileName := f.sourceFolder + "/" + filaName

	sourceFile, err := os.Open(sourceFileName)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destFileName)
	if err != nil {
		return nil
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (f *files) extension(fp string) string {
	extension := filepath.Ext(fp)
	return extension[1:]
}

func (f *files) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
