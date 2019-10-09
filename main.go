package main

import (
	"flag"
	"fmt"
	"github.com/javscrape/go-scrape"

	"github.com/javscrape/go-scrape/net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	proxy := flag.String("proxy", "", "set proxy")
	path := flag.String("video", "./", "set the video path")
	output := flag.String("output", "./video", "set the info output path")
	flag.Parse()
	fmt.Println("jav movie running")
	fmt.Println("read path:", *path)
	if *proxy != "" {
		e := net.RegisterProxy(*proxy)
		if e != nil {
			panic(e)
		}
	}

	list := getFileNames(*path)
	for _, n := range list {
		fmt.Println("name:", n)
		grab2 := scrape.NewGrabJavbus()
		grab3 := scrape.NewGrabJavdb()
		s := scrape.NewScrape(grab2, grab3)
		s.Output(*output)
		s.GrabSample(true)
		s.ImageCache("")
		msg, e := s.Find(n)
		if e != nil {
			panic(e)
		}
		for _, m := range *msg {
			fmt.Printf("message: %+v\n", m)
		}
	}
}

func getFileNames(path string) (files []string) {
	info, e := os.Stat(path)
	if e != nil {
		return nil
	}

	if info.IsDir() {
		file, e := os.Open(path)
		if e != nil {
			return nil
		}
		defer file.Close()
		names, e := file.Readdirnames(-1)
		if e != nil {
			return nil
		}
		var fullPath string
		for _, name := range names {
			fullPath = filepath.Join(path, name)
			files = append(files, getName(fullPath))
		}
	} else {
		files = append(files, getName(path))
	}

	return files
}

func getName(file string) string {
	info, e := os.Stat(file)
	if e != nil {
		return file
	}
	if info.IsDir() {
		return filepath.Base(file)
	}
	ext := filepath.Ext(file)
	return strings.TrimSuffix(filepath.Base(file), ext)
}
