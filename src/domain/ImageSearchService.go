package domain

type ImageSearchService struct{}

func (E ImageSearchService) GetImagesByKeyword(keyword string) []Image {
	panic("implement me")
}

func (E ImageSearchService) GetImagesByRegex(regex string) []Image {
	panic("implement me")
}

func (E ImageSearchService) GetImagesByKeywords(keywords []string) []Image {
	panic("implement me")
}
