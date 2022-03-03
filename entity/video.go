package entity

type Video struct {
	Title       string `json:"title" binding:"min=2,max=50"`
	Description string `json:"description" binding:"max=100"`
	Url         string `json:"url" binding:"required,url"`
	Author      Person `json:"author" binding:"required"`
}
