package service

import (
	"company/system/microservices/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TITLE string = "Titutlo del Video 01"
const DESCRIPTION string = "Descripción del Video 01"
const URL string = "Url del Video 01"

func getVideo() entity.Video {
	return entity.Video{
		Title:       TITLE,
		Description: DESCRIPTION,
		Url:         URL,
	}
}

func TestFindAll(t *testing.T) {
	service := NewVideoService()

	service.Save(getVideo())

	videos := service.FindAll()

	firstVideo := videos[0]
	assert.NotNil(t, videos)
	assert.Equal(t, TITLE, firstVideo.Title)
	assert.Equal(t, DESCRIPTION, firstVideo.Description)
	assert.Equal(t, URL, firstVideo.Url)
}
