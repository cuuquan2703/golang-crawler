package utils

type JSONFile struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
	ImgUrl     []string `json:"imgUrl"`
	RelatedUrl []string `json:"relatedUrl"`
}

type Option struct {
	MaxDepth int      `json:"maxDepth"`
	Tag      []string `json:"tag"`
	BoldText []string `json:"boldText"`
}

type Body struct {
	Url     []string `json:"url"`
	Options Option   `json:"options"`
}

type Response struct {
	Status  string `json:"status"`
	Message any    `json:"message"`
}

type TagHTML struct {
	Class string `json:"class"`
	Text  string `json:"text"`
}
