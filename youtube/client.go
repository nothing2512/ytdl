package youtube

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	ApiKey string
}

func (c *Client) Initiate() {
	c.ApiKey = "AIzaSyACehpcKKyng24lMaMg1Djo4g1rHjwg5rk"
}

func (c *Client) GetVideo(id string) (*VideoData, error) {
	uri := "https://www.youtube.com/youtubei/v1/player?key=" + c.ApiKey
	body := map[string]interface{}{
		"videoId": id,
		"context": map[string]interface{}{
			"client": map[string]interface{}{
				"hl":                "en",
				"gl":                "US",
				"clientName":        "ANDROID",
				"clientVersion":     "18.11.34",
				"androidSDKVersion": 30,
				"userAgent":         "com.google.android.youtube/18.11.34 (Linux; U; Android 11) gzip",
				"timeZone":          "UTC",
				"utcOffsetMinutes":  0,
			},
		},
		"playbackContext": map[string]interface{}{
			"contentPlaybackContext": map[string]interface{}{
				"html5Preference": "HTML5_PREF_WANTS",
			},
		},
		"contentCheckOk": true,
		"racyCheckOk":    true,
		"params":         "CgIQBg==",
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Youtube-Client-Name", "3")
	req.Header.Set("X-Youtube-Client-Version", "18.11.34")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data VideoData
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	data.reformat()
	data.client = c
	data.downloader = &Downloader{}

	return &data, nil
}

func (c *Client) GetPlaylist(id string, pageToken string) (*Playlist, error) {

	queryParams := url.Values{}
	queryParams.Add("key", c.ApiKey)
	queryParams.Add("playlistId", id)
	queryParams.Add("part", "contentDetails")
	queryParams.Add("maxResults", "50")
	queryParams.Add("pageToken", pageToken)

	uri := "https://www.googleapis.com/youtube/v3/playlistItems?" + queryParams.Encode()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data Playlist
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	if data.NextPageToken != "" {
		_data, err := c.GetPlaylist(id, data.NextPageToken)
		if err != nil {
			return nil, err
		}
		data.Items = append(data.Items, _data.Items...)
	}

	data.client = c

	return &data, nil
}
