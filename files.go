package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

const (
	paralellFileCount = 10
)

type filer interface {
	Split(string, string, bool, bool) error
}

type files struct {
	sourceFolder      string
	destinationFolder string
	overwrite         bool
}

type fileInfo struct {
	path     string
	relPath  string
	fileInfo os.FileInfo
}

func (f *files) Split(sf, df string, overwrite, flat bool) error {
	semaphore := make(chan struct{}, paralellFileCount)
	safeCounter := &counter{}

	f.overwrite = overwrite
	f.sourceFolder = sf
	f.destinationFolder = df

	files, err := f.readDir()
	if err != nil {
		return err
	}

	fileCount := len(*files)
	fmt.Printf("Processing %d files\n", fileCount)

	for i, file := range *files {
		fileName := file.fileInfo.Name()
		var folderName string
		if flat {
			folderName = f.extension(fileName) + "/" + file.fileInfo.ModTime().Format("2006/01")
		} else {
			folderName = file.relPath + f.extension(fileName) + "/" + file.fileInfo.ModTime().Format("2006/01")
		}

		safeCounter.Increment()
		go func(c *counter, i, l int, file fileInfo) {
			defer c.Decrement()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := f.mkDir(folderName)
			if err == nil {
				err = f.copyFile(file.path, folderName, fileName)
				if err != nil {
					fmt.Println("Copy error", err)
				}
			}

		}(safeCounter, i, fileCount, file)

		for {
			if safeCounter.Value() <= paralellFileCount {
				fmt.Printf("\rProgress %.0f%%", math.Round(float64(i)/float64(fileCount)*100))
				break
			}
		}
	}

	for {
		if safeCounter.Value() == 0 {
			break
		}
	}

	fmt.Println("\nDone.")

	return nil
}

func (f *files) readDir() (*[]fileInfo, error) {
	var files []fileInfo

	err := filepath.Walk(f.sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(f.sourceFolder, path)
			if err != nil {
				return err
			}
			dir, _ := filepath.Split(relPath)

			files = append(files, fileInfo{fileInfo: info, path: path, relPath: dir})
		}

		return nil
	})

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

func (f *files) copyFile(sourceFileName, folderName, fileName string) error {
	destFileName := f.destinationFolder + "/" + folderName + "/" + fileName
	if !f.overwrite && f.fileExists(destFileName) {
		return nil
	}

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
