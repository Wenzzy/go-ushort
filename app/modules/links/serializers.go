package links

import (
	"github.com/gin-gonic/gin"
	db_utils "go-ushort/app/common/db-utils"
	"go-ushort/app/models"
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
	C           *gin.Context
	Links       []models.LinkModel
	TotalCount  int64
	Took        int
	CurrentPage int
}

func (s *LinksSerializer) Response() db_utils.GetAllResponse {
	data := make([]LinkResponse, 0, len(s.Links))
	for _, lm := range s.Links {
		data = append(data, LinkResponse{
			lm.ID,
			*lm.Name,
			lm.RealUrl,
			lm.GeneratedAlias,
			lm.CreatedAt.String(),
		})
	}
	getAllSrz := db_utils.GetAllSerializer{
		Data:        data,
		TotalCount:  s.TotalCount,
		Took:        s.Took,
		CurrentPage: s.CurrentPage,
	}
	return getAllSrz.Response()
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
