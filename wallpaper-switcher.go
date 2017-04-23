package main

import (
	"os"
	"strconv"
	"io/ioutil"
	"os/exec"
	"github.com/rjeczalik/notify"
	"time"
)

var currentWallpaperIndex = -1
var wallpapersPath string
var wallpapers []string
var watcherEvents = make(chan notify.EventInfo, 1)
var switchEvents = make(chan int, 1)

// Re-reads directory, that contains wallpapers,
// and saves file pathes to wallpapers variable
func updateWallpapersList() {
	files, err := ioutil.ReadDir(wallpapersPath)
	if err != nil {
		panic("Unable to read wallpapers directory")
	}
	wallpapers = []string{}
	for _, file := range files {
		if !file.IsDir() {
			wallpapers = append(wallpapers, wallpapersPath+file.Name())
		}
	}
}

// Calls `feh` to set wallpaper from given file
func setWallpaper(path string) {
	exec.Command("feh", "--bg-fill", path).Run()
}

// Listens fs events on directory with wallpapers
// for updating its list when something changes
func watchFs() {
	notify.Watch(wallpapersPath, watcherEvents, notify.Create, notify.Remove)
	go func() {
		for range watcherEvents {
			updateWallpapersList()
		}
	}()
}

// Switches to next wallpaper by timeout or
// when switch event comes
func runTimer(timeout int) {
	go func() {
		for {
			select {
			case _ = <-switchEvents:
				nextWallpaper()
			case <-time.After(time.Second * time.Duration(timeout)):
				nextWallpaper()
			}
		}
	}()
}

// Calculates index of next wallpaper and sets it
func nextWallpaper() {
	currentWallpaperIndex = (currentWallpaperIndex + 1) % len(wallpapers)
	setWallpaper(wallpapers[currentWallpaperIndex])
}

// Infinitely reads from given fifo and
// emits switch event when something are written to fifo
func runReader(file string) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		exec.Command("mkfifo", file).Run()
	}
	for {
		ioutil.ReadFile(file)
		switchEvents <- 1
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		panic("Invalid arguments count, must be 3: wallpapers path, input fifo and switch time in seconds")
	}
	wallpapersPath = args[0]
	inputFifo := args[1]
	timeout, err := strconv.Atoi(args[2])
	if err != nil {
		panic(err)
	}
	updateWallpapersList()
	watchFs()
	runTimer(timeout)
	switchEvents <- 1
	runReader(inputFifo)
}
