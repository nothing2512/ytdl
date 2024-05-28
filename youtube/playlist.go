package youtube

type PlaylistItem struct {
	ContentDetails struct {
		VideoId string `json:"videoId"`
	} `json:"contentDetails"`
}

type Playlist struct {
	client *Client `json:"-"`

	Items         []PlaylistItem `json:"items"`
	NextPageToken string         `json:"nextPageToken"`
}

func (p *Playlist) Download(output string) error {
	to := len(p.Items)
	for k, v := range p.Items {
		video, err := p.client.GetVideo(v.ContentDetails.VideoId)
		if err != nil {
			return err
		}
		err = video.Download(output, k+1, to)
		if err != nil {
			return err
		}
	}
	return nil
}
