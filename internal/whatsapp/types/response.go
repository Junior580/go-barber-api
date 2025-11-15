package types

type Option struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Response struct {
	Data struct {
		Options []Option `json:"options"`
	} `json:"data"`
}
