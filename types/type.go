package types

type OrderedImages struct {
	Images []FormattedImage
}

type FormattedImage struct {
	Id           string
	Title        string
	Date         string
	Url          string
	UploaderName string
	UserId       string
}
