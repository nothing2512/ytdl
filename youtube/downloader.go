package youtube

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Downloader struct {
	totalBytes      int
	downloadedBytes int
	filename        string
	f               *os.File
	from            int
	to              int
}

func (d *Downloader) download(url, filename string, from, to int) error {
	d.filename = filename
	d.from = from
	d.to = to

	_, err := os.Stat(d.filename)
	if err == nil {
		return nil
	}

	_, err = os.Stat(strings.Split(d.filename, ".")[0] + ".mp3")
	if err == nil {
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contentLength := resp.Header.Get("Content-Length")

	if contentLength != "" {
		d.totalBytes, err = strconv.Atoi(contentLength)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Downloading... (unknown size)")
	}

	d.f, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer d.f.Close()

	_, err = io.Copy(d, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println() // Newline after progress bar

	// d.convertAudio()

	return nil
}

func (d *Downloader) convertAudio() error {
	mp3 := strings.Split(d.filename, ".")[0] + ".mp3"
	cmd := exec.Command("ffmpeg", "-i", d.filename, mp3)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	os.Remove(d.filename)
	return nil
}

func (d *Downloader) printProgress() {
	if d.totalBytes == 0 {
		return
	}
	completed := int(float64(d.downloadedBytes) / float64(d.totalBytes) * 100)
	fmt.Printf("[%v/%v] %s [%s%s] %d%%\r", d.from, d.to, d.filename, strings.Repeat("=", completed/2), strings.Repeat(" ", (100-completed)/2), completed)
}

func (d *Downloader) Write(p []byte) (n int, err error) {
	n = len(p)
	d.downloadedBytes += n
	d.printProgress()
	if _, err := d.f.Write(p); err != nil {
		return 0, err
	}
	return
}

var canonicals = map[string]string{
	"video/quicktime":  ".mov",
	"video/x-msvideo":  ".avi",
	"video/x-matroska": ".mkv",
	"video/mpeg":       ".mpeg",
	"video/webm":       ".webm",
	"video/3gpp2":      ".3g2",
	"video/x-flv":      ".flv",
	"video/3gpp":       ".3gp",
	"video/mp4":        ".mp4",
	"video/ogg":        ".ogv",
	"video/mp2t":       ".ts",
	"audio/mp4":        ".m4a",
}

func pickIdealFileExtension(mediaType string) string {
	mediaType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return ".mov"
	}

	if extension, ok := canonicals[mediaType]; ok {
		return extension
	}

	// Our last resort is to ask the operating system, but these give multiple results and are rarely canonical.
	extensions, err := mime.ExtensionsByType(mediaType)
	if err != nil || extensions == nil {
		return ".mov"
	}

	return extensions[0]
}
