package links

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/models"
	"net/http"
)

// GetAllLinks godoc
//
//	@Summary	Get all links for user
//	@Schemes
//	@Tags		Links
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Success	200	{object}	[]LinkResponse
//	@Router		/links [get]
func GetAllLinks(c *gin.Context) {

	links, err := models.FindAllLinks(&models.LinkModel{UserID: c.MustGet("user_id").(uint)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	sz := Serializer{c, links}
	res := sz.Response()

	c.JSON(http.StatusOK, res)
}
