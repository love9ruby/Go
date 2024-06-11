package url

type ShortenRequest struct {
	URL      string `json:"url"`
	Password string `json:"password"`
}
