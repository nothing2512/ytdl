package youtube

import (
	"path"
	"sort"
	"strings"
)

type VideoFormat struct {
	Itag         int    `json:"itag"`
	Url          string `json:"url"`
	MimeType     string `json:"mimeType"`
	FPS          int    `json:"fps"`
	Quality      string `json:"quality"`
	AudioQuality string `json:"audioQuality"`
}

type VideoData struct {
	client     *Client     `json:"-"`
	downloader *Downloader `json:"-"`

	StreamingData struct {
		Formats         []VideoFormat `json:"formats"`
		AdaptiveFormats []VideoFormat `json:"adaptiveFormats"`
	} `json:"streamingData"`
	VideoDetails struct {
		Title string `json:"title"`
	} `json:"videoDetails"`
	Formats []VideoFormat
	Title   string
}

func (v *VideoData) reformat() {

	getAudioQuality := func(d string) int {
		if d == "AUDIO_QUALITY_HIGH" {
			return 3
		}
		if d == "AUDIO_QUALITY_MEDIUM" {
			return 2
		}
		if d == "AUDIO_QUALITY_LOW" {
			return 1
		}
		return 0
	}

	v.Title = v.VideoDetails.Title
	v.Formats = append(v.StreamingData.Formats, v.StreamingData.AdaptiveFormats...)

	sort.Slice(v.Formats, func(i, j int) bool {
		if getAudioQuality(v.Formats[i].AudioQuality) > getAudioQuality(v.Formats[j].AudioQuality) {
			return true
		}
		if getAudioQuality(v.Formats[i].AudioQuality) == getAudioQuality(v.Formats[j].AudioQuality) {
			if strings.Contains(v.Formats[i].MimeType, "audio/mp4") {
				return true
			}
			if strings.Contains(v.Formats[j].MimeType, "audio/mp4") {
				return false
			}
		}
		return v.Formats[i].FPS < v.Formats[j].FPS
	})
}

func (v *VideoData) Download(output string, from, to int) error {
	outputFile := path.Join(output, v.Title+pickIdealFileExtension(v.Formats[0].MimeType))
	return v.downloader.download(v.Formats[0].Url, outputFile, from, to)
}
