package entity

type Video struct {
	Title       string `json:"title" binding:"min=2,max=10"`
	Description string `json:"description" binding:"max=20"`
	Url         string `json:"url" binding:"required,url"`
	Author      Person `json:"author" binding:"required"`
}
