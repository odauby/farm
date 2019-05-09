package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/xml"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var seed = map[string]string{} // map[rom crc32]"file path"

type Datfile struct {
	Header Header `xml:"header"`
	Sets   []Set  `xml:"game"`
}

type Header struct {
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Category    string `xml:"category"`
	Version     string `xml:"version"`
	Author      string `xml:"author"`
	Comment     string `xml:"comment"`
}

type Set struct {
	Name         string `xml:"name,attr"`
	CloneOf      string `xml:"cloneof,attr"`
	RomOf        string `xml:"romof,attr"`
	Description  string `xml:"description"`
	Year         string `xml:"year"`
	Manufacturer string `xml:"manufacturer"`
	Roms         []Rom  `xml:"rom"`
}

type Rom struct {
	Name  string `xml:"name,attr"`
	Merge string `xml:"merge,attr"`
	Size  string `xml:"size,attr"`
	Crc   string `xml:"crc,attr"`
	Sha1  string `xml:"sha1,attr"`
}

func getFileSha1(file string) string {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	h := sha1.New()
	h.Write(dat)
	return fmt.Sprintf("%x", h.Sum(nil))

}

func getFileCRC32(file string) string {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	h := crc32.New(crc32.IEEETable)
	h.Write(dat)
	return fmt.Sprintf("%x", h.Sum(nil))

}

func l(line ...interface{}) {

	log.Println(line...)

}

func walkSourceDir(root string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("error: ", err)
		}
		if !info.IsDir() {
			currentCRC32 := getFileCRC32(path)
			seed[currentCRC32] = path
			// fmt.Print("\r", path, " ", currentCRC32)

		}

		return nil
	})
	if err != nil {
		fmt.Println("error: ", err)
	}
	// fmt.Print("\r")
}

func zipSet(destDir string, romSet *Set, mode string) {
	zipFile := destDir + string(os.PathSeparator) + romSet.Name + ".zip"
	l("creating zip:", zipFile)

	zipf, _ := os.Create(zipFile)
	defer zipf.Close()

	zipw := zip.NewWriter(zipf)
	defer zipw.Close()

	for _, r := range romSet.Roms {
		currentRomFile := seed[r.Crc]
		switch mode {
		case "nonmerge":
			f, _ := os.Open(currentRomFile)
			defer f.Close()
			w, _ := zipw.Create(r.Name)
			io.Copy(w, f)

		case "split":
			if r.Merge == "" {
				f, _ := os.Open(currentRomFile)
				defer f.Close()
				w, _ := zipw.Create(r.Name)
				io.Copy(w, f)
			}

		}

	}

}

func main() {
	l("f.a.r.m. is starting")

	datPtr := flag.String("datfile", "", "DAT file location")
	sourcePtr := flag.String("source", "", "Location of the source files")
	destPtr := flag.String("dest", "", "Location of the destination files")
	modePtr := flag.String("mode", "nonmerge", "Desired output mode (nonmerge, split or merge")

	flag.Parse()

	canStart := true
	if *datPtr == "" {
		canStart = false
		l("[ERROR] no DAT file location provided, use \"-datfile\" flag to specify what reference DAT file you want to use")
	}

	if *sourcePtr == "" {
		canStart = false
		l("[ERROR] no source directory specified, use \"-source\" flag to specify where your source files (roms) are located")
	}

	if *destPtr == "" {
		canStart = false
		l("[ERROR] no destination directory specified, use \"-dest\" flag to specify in what directory you want f.a.r.m. to write your reorganized files")
	}

	if *sourcePtr == *destPtr && *sourcePtr != "" {
		canStart = false
		l("[ERROR] source and destination directories are the same ... this can't be good :-/")
	}

	if *modePtr == "merge" {
		canStart = false
		l("[ERROR] merge collection mode is not supported, yet. Sorry about that.")
	}

	if !canStart {
		os.Exit(1)
	}

	l("start reading file " + *datPtr)
	dat, err := ioutil.ReadFile(*datPtr)
	if err != nil {
		panic(err)
	}

	l("done reading file ")
	var datFile Datfile
	l("start marshalling XML")
	xml.Unmarshal(dat, &datFile)
	l("done marshalling XML")
	l("[name]       ", datFile.Header.Name)
	l("[description]", datFile.Header.Description)
	l("[category]   ", datFile.Header.Category)
	l("[version]    ", datFile.Header.Version)
	l("[author]     ", datFile.Header.Author)
	l("[comment]    ", datFile.Header.Comment)
	l("start parsing sets and roms")

	rom := map[string]Rom{} // map[rom crc32]{rom details}

	for _, s := range datFile.Sets {

		for _, r := range s.Roms {
			if r.Merge == "" {
				rom[r.Crc] = r
			}

		}

	}
	l("done parsing sets and roms")
	l("found", len(datFile.Sets), "sets and", len(rom), "roms")

	l("parsing source directory", *sourcePtr)
	walkSourceDir(*sourcePtr)
	l("identified", len(seed), "unique roms in directory", *sourcePtr)

	// Let's see what we have...
	for _, s := range datFile.Sets {
		complete := len(s.Roms) > 0
		for _, r := range s.Roms {

			_, gotIt := seed[r.Crc]
			complete = complete && gotIt
		}
		if complete {
			l("Set", s.Name, "is complete")
			zipSet(*destPtr, &s, *modePtr)
		}
	}

}
