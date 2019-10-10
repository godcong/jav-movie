package main

import (
	"flag"
	"fmt"
	"github.com/javscrape/go-scrape"
	"path"

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
		if n == "" {
			continue
		}
		grab2 := scrape.NewGrabJavbus()
		grab3 := scrape.NewGrabJavdb()
		s := scrape.NewScrape(grab3, grab2)
		s.Output(*output)
		s.GrabSample(true)
		s.ImageCache("")
		msg, e := s.Find(getName(n))
		if e != nil {
			panic(e)
		}
		if len(*msg) == 0 {
			fmt.Println("no data:", n)
			continue
		}

		info, _ := os.Stat(n)
		if info.IsDir() {
			e := os.Rename(n, filepath.Join(*output, strings.ToUpper(filepath.Base(n))))
			if e != nil {
				fmt.Println("dir error:", e)
			}
		} else {
			_ = os.MkdirAll(filepath.Join(*output, strings.ToUpper(getName(n))), os.ModePerm)
			ext := filepath.Ext(n)
			name := strings.ToUpper(getName(n))
			e := os.Rename(n, filepath.Join(*output, name, name+ext))
			if e != nil {
				fmt.Println("file error:", e)
			}
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
			fileInfo, e := os.Stat(fullPath)
			if e != nil {
				continue
			}
			if !fileInfo.IsDir() && !IsVideo(fullPath) {
				continue
			}
			files = append(files, fullPath)
		}
	} else {
		files = append(files, path)
	}

	return files
}

func getName(file string) string {
	info, e := os.Stat(file)
	if e != nil {
		return ""
	}
	if info.IsDir() {
		return filepath.Base(file)
	}
	ext := filepath.Ext(file)
	return strings.TrimSuffix(filepath.Base(file), ext)
}

// IsVideo ...
func IsVideo(filename string) bool {
	video := `.swf,.flv,.3gp,.ogm,.vob,.m4v,.mkv,.mp4,.mpg,.mpeg,.avi,.rm,.rmvb,.mov,.wmv,.asf,.dat,.asx,.wvx,.mpe,.mpa`
	ext := strings.ToLower(path.Ext(filename))
	return strings.Index(video, ext) != -1
}
