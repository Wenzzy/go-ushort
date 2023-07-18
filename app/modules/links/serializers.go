package links

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/models"
)

type LinkSerializer struct {
	C *gin.Context
	models.LinkModel
}

type LinkResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	RealUrl   string `json:"realUrl"`
	Alias     string `json:"alias"`
	CreatedAt string `json:"createdAt"`
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

type LinksSerializer struct {
	C     *gin.Context
	Links []models.LinkModel
}

func (s *LinksSerializer) Response() []LinkResponse {
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

type CreatedLinkSerializer struct {
	C *gin.Context
	models.LinkModel
}

type CreatedLinkResponse struct {
	ID    uint   `json:"id"`
	Alias string `json:"alias"`
}

func (s *CreatedLinkSerializer) Response() CreatedLinkResponse {
	return CreatedLinkResponse{
		s.ID,
		s.GeneratedAlias,
	}
}
