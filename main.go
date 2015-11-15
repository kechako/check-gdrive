package main

import (
	"fmt"
	"os"

	"github.com/kechako/check-gdrive/gdrive"
	"github.com/kechako/check-gdrive/local"

	"golang.org/x/net/context"
)

func makeGDriveNameMap(files []*gdrive.File) map[string]*gdrive.File {
	nameMap := make(map[string]*gdrive.File)

	for _, file := range files {
		nameMap[file.Title] = file
	}

	return nameMap
}

func makeLocalNameMap(files []*local.File) map[string]*local.File {
	nameMap := make(map[string]*local.File)

	for _, file := range files {
		nameMap[file.Name()] = file
	}

	return nameMap
}

func diffFile(gFile *gdrive.File, lFile *local.File) (bool, error) {
	if gFile.IsFolder() != lFile.IsDir() {
		return false, nil
	}

	if lFile.IsDir() {
		return true, nil
	}

	md5Checksum, err := lFile.Md5Checksum()
	if err != nil {
		return false, err
	}

	return md5Checksum == gFile.Md5Checksum, nil
}

func traverse(gFolder *gdrive.File, lFolder *local.File) error {
	gFiles, err := gFolder.GetFiles()
	if err != nil {
		return err
	}

	lFiles, err := lFolder.GetFiles()
	if err != nil {
		return err
	}

	gMap := makeGDriveNameMap(gFiles)
	lMap := makeLocalNameMap(lFiles)

	for _, lFile := range lFiles {
		gFile, ok := gMap[lFile.Name()]
		if ok {
			matched, err := diffFile(gFile, lFile)
			if err != nil {
				return err
			}
			if !matched {
				fmt.Printf("+-:%s:%s\n", gFile.Path, lFile.Path)
			}

			if lFile.IsDir() {
				err = traverse(gFile, lFile)
				if err != nil {
					return err
				}
			}
		} else {
			fmt.Printf("+ :%s:%s\n", gFolder.Join(lFile.Name()), lFile.Path)
		}
	}
	for _, gFile := range gFiles {
		_, ok := lMap[gFile.Title]
		if !ok {
			fmt.Printf(" -:%s:%s\n", gFile.Path, lFolder.Join(gFile.Title))
		}
	}

	return nil
}

func _main() (int, error) {
	folderId := "root"
	baseDir, err := os.Getwd()
	if err != nil {
		return 1, err
	}
	if len(os.Args) > 1 {
		folderId = os.Args[1]
	}
	if len(os.Args) > 2 {
		baseDir = os.Args[2]
	}

	ctx := context.Background()

	g, err := gdrive.New(ctx)
	if err != nil {
		return 2, err
	}

	gFolder, err := g.GetFile(folderId)
	if err != nil {
		return 3, err
	}
	lFolder, err := local.NewFile(baseDir)
	if err != nil {
		return 4, err
	}

	err = traverse(gFolder, lFolder)
	if err != nil {
		return 5, err
	}

	return 0, nil
}

func main() {
	if exitStatus, err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] %v\n", err)
		os.Exit(exitStatus)
	}
}
