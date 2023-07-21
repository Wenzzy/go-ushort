package links

import (
	"github.com/gin-gonic/gin"
	"go-ushort/app/common/constants/emsgs"
	"go-ushort/app/common/database"
	db_utils "go-ushort/app/common/db-utils"
	"go-ushort/app/common/utils"
	"go-ushort/app/models"
	"net/http"
)

// GetAll godoc
//
//	@Summary	Get all links for user
//	@Schemes
//	@Tags		Links
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		take	query		int	false	"Limit links per page"
//	@Param		page	query		int	false	"Page number"
//	@Success	200		{object}	db_utils.GetAllResponse{data=[]LinkResponse}
//	@Router		/links [get]
func GetAll(c *gin.Context) {

	pag := db_utils.GenPagination(c)
	db := database.GetDB()

	var links []models.LinkModel
	var linksCount int64

	qb := db.
		Model(&models.LinkModel{}).
		Where(&models.LinkModel{UserID: c.MustGet("userId").(uint)})

	err := qb.Count(&linksCount).Error
	err = qb.
		Limit(pag.Take).
		Offset(pag.Offset).
		Order("id asc").
		Find(&links).
		Error

	if err != nil {
		c.JSON(utils.NewError(http.StatusInternalServerError, emsgs.Internal).H())
		return
	}

	sz := LinksSerializer{
		C:           c,
		Links:       links,
		Took:        pag.Take,
		TotalCount:  linksCount,
		CurrentPage: pag.Page,
	}
	res := sz.Response()

	c.JSON(http.StatusOK, res)
}

// Create godoc
//
//	@Summary	Create new link
//	@Schemes
//	@Tags		Links
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		request	body		LinkCreateValidator	true	"Request Body"
//	@Success	200		{object}	CreatedLinkResponse
//	@Failure	422		{object}	utils.CommonValidationError
//	@Failure	400		{object}	utils.CommonError
//	@Router		/links [post]
func Create(c *gin.Context) {
	validator := NewLinkCreateValidator()
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
		return
	}

	var foundLink *models.LinkModel

	db := database.GetDB()

	if err := db.Where(&models.LinkModel{RealUrl: validator.RealUrl}).First(&foundLink).Error; err == nil {
		sz := CreatedLinkSerializer{c, *foundLink}
		res := sz.Response()
		c.JSON(http.StatusOK, res)
		return
	}

	alias := db_utils.GenUniqueLinkAlias()

	model := &models.LinkModel{
		Name:           &validator.Name,
		RealUrl:        validator.RealUrl,
		UserID:         c.MustGet("userId").(uint),
		GeneratedAlias: alias,
	}

	if err := db_utils.SaveOne(model, "link"); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	sz := CreatedLinkSerializer{c, *model}
	res := sz.Response()
	c.JSON(http.StatusCreated, res)
}

// Update godoc
//
//	@Summary	Update link
//	@Schemes
//	@Tags		Links
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		linkId	path	string				true	"Link ID"
//	@Param		request	body	LinkUpdateValidator	true	"Request Body"
//	@Success	204		"Empty response"
//	@Failure	422		{object}	utils.CommonValidationError
//	@Failure	404		{object}	utils.CommonError
//	@Router		/links/{linkId} [patch]
func Update(c *gin.Context) {
	uri := LinkUpdateUriValidator{}
	if err := c.BindUri(&uri); err != nil {
		c.JSON(utils.NewError(http.StatusNotFound, emsgs.ObjectNotFound, "link").H())
		return
	}
	validator := NewLinkUpdateValidator()
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
		return
	}

	searchLink := &models.LinkModel{ID: uri.ID, UserID: c.MustGet("userId").(uint)}

	var foundLink *models.LinkModel

	db := database.GetDB()

	if err := db.Where(searchLink).First(&foundLink).Error; err != nil {
		c.JSON(utils.NewError(http.StatusNotFound, emsgs.ObjectNotFound, "link").H())
		return
	}

	if err := db.Model(foundLink).Where(searchLink).Update("name", validator.Name).Error; err != nil {
		c.JSON(utils.NewError(http.StatusInternalServerError, emsgs.Internal, "Can't update link").H())
		return
	}
	c.Status(http.StatusNoContent)
}

// Redirect godoc
//
//	@Summary	Redirect from alias to realLink
//	@Schemes
//	@Tags		Links
//	@Accept		json
//	@Produce	json
//	@Param		alias	path	string	true	"Short URL Alias"
//	@Success	302		"Redirect"
//	@Failure	404		{object}	utils.CommonError
//	@BasePath	/
//	@Router		/{alias} [get]
func Redirect(c *gin.Context) {
	uri := RedirectUriValidator{}
	if err := c.BindUri(&uri); err != nil {
		c.JSON(utils.NewError(http.StatusNotFound, emsgs.ObjectNotFound, "alias").H())
		return
	}
	var link *models.LinkModel

	db := database.GetDB()

	if err := db.Where(&models.LinkModel{GeneratedAlias: uri.Alias}).First(&link).Error; err != nil {
		c.JSON(utils.NewError(http.StatusNotFound, emsgs.ObjectNotFound, "alias").H())
		return
	}
	c.Redirect(http.StatusFound, link.RealUrl)
}
