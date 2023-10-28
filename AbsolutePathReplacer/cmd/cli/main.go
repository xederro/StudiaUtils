package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	wg  *sync.WaitGroup = new(sync.WaitGroup)
	ex  *regexp.Regexp
	emd *regexp.Regexp
	abs *regexp.Regexp
	md  []string
)

func main() {
	pathFlag := flag.String("path", ".", "path to directory")
	gitFlag := flag.Bool("git", false, "looks at changed, added etc. files in git")
	flag.Parse()

	c, err := regexp.Compile("excalidraw")
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}
	ex = c

	c, err = regexp.Compile("\\.md(\")?$")
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}
	emd = c

	c, err = regexp.Compile("\\(Notatki")
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}
	abs = c

	wg.Add(1)
	if *gitFlag {
		go MDFilesFromGit(*pathFlag)
	} else {
		go MDFilesFromDir(*pathFlag)
	}

	wg.Wait()
	fmt.Println(md)
	for _, path := range md {
		fmt.Println(path)
		wg.Add(1)
		go Replace(path)
	}
	wg.Wait()
}

func MDFilesFromGit(dir string) {
	cmd := exec.Command("git", "-C", dir, "-c", "core.quotepath=false", "ls-files", "-o", "-d", "-c", "-m")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	s := strings.Split(out.String(), "\n")

	for _, f := range s {
		if emd.Match([]byte(f)) && !ex.Match([]byte(f)) {
			md = append(md, f)
		}
	}
	wg.Done()
}

func MDFilesFromDir(dir string) {
	open, err := os.Open(dir)

	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {

		}
	}(open)

	infos, err := open.Readdir(-1)
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	for _, info := range infos {
		n := info.Name()
		if info.IsDir() {
			wg.Add(1)
			go MDFilesFromDir(dir + "/" + n)
		} else if emd.Match([]byte(n)) && !ex.Match([]byte(n)) {
			md = append(md, dir+"/"+n)
		}
	}

	wg.Done()
}

func Replace(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("An Error Occurred:", err)
	}

	sp := abs.Split(string(data), -1)
	if len(sp) != 1 {
		tmp, err := os.CreateTemp(filepath.Dir(path), "replace-*")
		if err != nil {
			log.Fatal(err)
		}
		defer func(tmp *os.File) {
			err := tmp.Close()
			if err != nil {

			}
		}(tmp)

		for k, f := range sp {
			_, err := tmp.Write([]byte(f))
			if err != nil {
				log.Fatalln("An Error Occurred:", err)
			}
			if k != len(sp)-1 {
				_, err = tmp.Write([]byte("(/Notatki"))
				if err != nil {
					log.Fatalln("An Error Occurred:", err)
				}
			}
		}

		if err := tmp.Close(); err != nil {
			log.Fatalln("An Error Occurred:", err)
		}

		if err := os.Rename(tmp.Name(), path); err != nil {
			log.Fatalln("An Error Occurred:", err)
		}
	}
	wg.Done()
}
