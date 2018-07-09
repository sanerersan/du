package main

import (
	"flag"
	"fmt"
	"github.com/sanerersan/du/capslice"
	duUtils "github.com/sanerersan/du/utils"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
)

var sortbyName *bool = flag.Bool("sn", false, "sort as name")
var sortReverse *bool = flag.Bool("r", false, "sort reverse")
var rootDir *string = flag.String("dir", "/", "root directory")

var wg duUtils.WaitGroupWapper
var capInfos *capslice.CapSlice
var sema *duUtils.SemaphoreWarpper

func init() {
	capInfos = capslice.NewCapSlice()
	sema, _ = duUtils.NewSemaphore("du_sema", 20, 0)
}

func main() {
	flag.Parse()

	wg.Run(WalkDirList, *rootDir)
	wg.Wait()

	sortAndDisplay()
}

func sortAndDisplay() {
	if *sortbyName {
		capInfos.SetSortByName()
	}

	if *sortReverse {
		sort.Sort(sort.Reverse(capInfos))
	} else {
		sort.Sort(capInfos)
	}

	infos := capInfos.GetCaps()
	for _, info := range infos {
		fmt.Println(info.Path, duUtils.GetHumanReadableString(info.Cap))
	}
}

func WalkDirList(dirPath string) {
	sema.WaitSync()
	defer sema.Release()

	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Println(err)
		return
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if ("." == entryName) || (".." == entryName) {
			continue
		}
		fullName := filepath.Join(dirPath, entryName)
		if entry.IsDir() {
			wg.Run(WalkDirList, fullName)
		} else {
			fci := &capslice.FileCapacityInfo{
				Name: entry.Name(),
				Path: fullName,
				Cap:  entry.Size(),
			}
			capInfos.Append(fci)
		}
	}

	return
}
