package links

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/models"
)

type Serializer struct {
	C     *gin.Context
	Links []models.LinkModel
}

type LinkSerializer struct {
	C *gin.Context
	models.LinkModel
}

type LinkResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	RealUrl   string `json:"real_url"`
	Alias     string `json:"alias"`
	CreatedAt string `json:"created_at"`
}

func (s *LinkSerializer) Response() LinkResponse {
	return LinkResponse{
		s.ID,
		*s.Name,
		s.RealUrl,
		s.GeneratedAlias,
		s.CreatedAt.String(),
	}
}

func (s *Serializer) Response() []LinkResponse {
	res := make([]LinkResponse, 0, len(s.Links))
	for _, lm := range s.Links {
		res = append(res, LinkResponse{
			lm.ID,
			*lm.Name,
			lm.RealUrl,
			lm.GeneratedAlias,
			lm.CreatedAt.String(),
		})
	}
	return res
}
