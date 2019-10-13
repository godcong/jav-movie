package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/javscrape/go-scrape"
	"path"
	"strconv"

	"github.com/javscrape/go-scrape/net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	proxy := flag.String("proxy", "", "set proxy")
	path := flag.String("video", "./", "set the video path")
	output := flag.String("output", "./video", "set the info output path")
	failed := flag.String("failed", "./failed", "set the failed output path")
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
	_ = os.MkdirAll(*failed, os.ModePerm)
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
		name := getName(n)
		if mulitiVideos(&name) {
			fmt.Println("multi:", n)
		}
		msg, e := s.Find(name)
		if e != nil {
			fmt.Println("find error:", e)
			continue
		}
		if len(*msg) == 0 {
			fmt.Println("no data:", n)
			e = moveTo(n, *failed, true)
			if e != nil {
				fmt.Println("move error:", e)
			}
			continue
		}
		e = moveTo(n, *output, false)
		if e != nil {
			fmt.Println("move error:", e)
		}

		for _, m := range *msg {
			fmt.Printf("message: %+v\n", m)
		}
	}
}

func mulitiVideos(name *string) bool {
	split := strings.Split(*name, "@")
	if len(split) == 2 {
		*name = split[0]
		return true
	}
	return false
}

func moveTo(file string, path string, namepath bool) (e error) {
	info, _ := os.Stat(file)
	if info.IsDir() {
		return moveDir(file, path, namepath)
	}
	ext := filepath.Ext(file)
	name := strings.ToUpper(getName(file))

	target := filepath.Join(path, name, name+ext)
	if namepath {
		target = filepath.Join(path, name+ext)
	} else {
		_ = os.MkdirAll(filepath.Join(path, strings.ToUpper(getName(file))), os.ModePerm)
	}
	_, e = os.Stat(target)
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if os.IsNotExist(e) {
		return os.Rename(file, target)
	}
	// exist create dir:name
	return moveBak(file, target)
}

func moveDir(file string, path string, namepath bool) (e error) {
	f, e := os.Open(file)
	if e != nil {
		return nil
	}
	defer f.Close()
	names, e := f.Readdirnames(-1)
	if e != nil {
		return e
	}
	var fullPath string
	for _, name := range names {
		fullPath = filepath.Join(file, name)
		e = moveTo(fullPath, path, namepath)
		if e != nil {
			fmt.Println("dir error:", file)
			continue
		}
	}
	return nil
}

func moveBak(file string, path string) (e error) {
	for count := 1; count < 10; count++ {
		dir := filepath.Dir(path)
		name := getName(file)
		ext := filepath.Ext(file)

		target := filepath.Join(dir, name+"_"+strconv.Itoa(count)+ext)
		_, e = os.Stat(target)
		if e != nil && !os.IsNotExist(e) {
			return e
		}
		if os.IsNotExist(e) {
			return os.Rename(file, target)
		}
		// exist create dir:name
	}
	return errors.New("unmoved")
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
