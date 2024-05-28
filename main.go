package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"main/youtube"
)

func handle(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func main() {
	data, out, isPlaylist, err := parseArgs()
	handle(err)

	if *convertPtr != "" {
		convertAudios()
		return
	}

	err = os.MkdirAll(out, os.ModePerm)
	handle(err)

	if isPlaylist {
		downloadPlaylist(data[0], out)
	} else {
		downloadFiles(data, out)
	}
}

func convertAudios() {
	files, err := os.ReadDir(*convertPtr)
	handle(err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), ".mp3") {
			continue
		}

		m4a := path.Join(*convertPtr, file.Name())
		mp3 := strings.Split(m4a, ".")[0] + ".mp3"
		cmd := exec.Command("ffmpeg", "-i", m4a, mp3)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		handle(err)
		os.Remove(m4a)
	}
}

func downloadFiles(datas []Data, output string) {
	for k, x := range datas {
		_, e := url.Parse(x.Link)
		if e != nil {
			datas[k].Id = strings.Split(x.Link, "?v=")[1]
		} else {
			datas[k].Id = x.Link
		}
	}

	c := youtube.Client{}
	c.Initiate()

	for _, x := range datas {
		video, err := c.GetVideo(x.Id)
		handle(err)
		if x.Rename != "" {
			video.Title = x.Rename
		}
		err = video.Download(output, 1, 1)
		handle(err)
	}
}

func downloadPlaylist(data Data, output string) {
	_, e := url.Parse(data.Link)
	if e != nil {
		data.Id = strings.Split(data.Link, "?v=")[1]
	} else {
		data.Id = data.Link
	}

	c := youtube.Client{}
	c.Initiate()

	playlist, err := c.GetPlaylist(data.Id, "")
	handle(err)
	err = playlist.Download(output)
	handle(err)
}

func exampleDownloadFile() {
	c := youtube.Client{}
	c.Initiate()
	video, err := c.GetVideo("k2qgadSvNyU")
	handle(err)
	err = video.Download("./", 1, 1)
	handle(err)
}

func exampleDownloadPlaylist() {
	c := youtube.Client{}
	c.Initiate()
	playlist, err := c.GetPlaylist("PLTxw13RLjMaiDCUsIljLCcNX3DcbyETKt", "")
	handle(err)
	err = playlist.Download("./")
	handle(err)
}
