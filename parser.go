package main

import (
	"flag"
	"io"
	"os"
	"strings"
)

type Data struct {
	Rename string `json:"rename"`
	Link   string `json:"string"`
	Id     string `json:"id"`
}

var (
	outputPtr   = flag.String("o", "./", "Output Dir")
	inputPtr    = flag.String("i", "", "Input File")
	playlistPtr = flag.String("p", "", "Playlist ID")
	convertPtr  = flag.String("c", "", "Convert To MP3 Folders")
)

func parseArgs() ([]Data, string, bool, error) {
	flag.Parse()

	if *convertPtr != "" {
		return nil, "", false, nil
	}

	if *playlistPtr != "" {
		return []Data{{
			Link: *playlistPtr,
		}}, *outputPtr, true, nil
	}

	inFile, err := os.Open(*inputPtr)
	if err != nil {
		return nil, *outputPtr, false, err
	}
	defer inFile.Close()

	bytes, err := io.ReadAll(inFile)
	if err != nil {
		return nil, *outputPtr, false, err
	}

	inputs := strings.Split(string(bytes), "\n")
	var data []Data

	for _, i := range inputs {
		fs := strings.Split(strings.TrimSpace(i), ",")
		if i == "" || strings.TrimSpace(i) == "" {
			continue
		}
		rename := ""
		if len(fs) > 1 {
			rename = strings.TrimSpace(fs[1])
		}
		data = append(data, Data{
			Link:   fs[0],
			Rename: rename,
		})
	}

	if *playlistPtr != "" {
		return []Data{{
			Link: *playlistPtr,
		}}, *outputPtr, true, nil
	}

	return data, *outputPtr, false, nil
}
