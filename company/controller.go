package company

import (
	"errors"
	"net/http"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/swaggo/swag/example/celler/httputil"
)

// Controller defines methods to a Controller
type Controller interface {
	Find(ctx *gin.Context)
	LoadWebsites(ctx *gin.Context)
	InitDatabase(string)
}

type companyController struct {
	service Service
}

// NewController return a new companyController
func NewController(service Service) Controller {
	return companyController{service}
}

func (c companyController) GetAll(ctx *gin.Context) {
	results, err := c.service.findAll()
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, results)
}

// Find godoc
// @Summary Show a company
// @Description get company by name and zipcode
// @ID get-company-by-name-and-zipcode
// @Produce json
// @Param name query string true "Name"
// @Param zipcode query string true "Zipcode"
// @Success 200 {array} company.Company
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /companies [get]
func (c companyController) Find(ctx *gin.Context) {
	name, hasName := ctx.GetQuery("name")
	zipcode, hasZip := ctx.GetQuery("zipcode")
	var result Company
	var err error

	if !hasName && !hasZip {
		c.GetAll(ctx)
		return
	}

	if hasName && hasZip {
		result, err = c.service.findByNameAndZipCode(name, zipcode)
	} else {
		err = errors.New("Missing parameters 'name' or 'zipcode'")
	}

	if err != nil {
		if err == mgo.ErrNotFound {
			log.WithError(err).Error("Company not found")
			ctx.Status(http.StatusNoContent)
			return
		}
		log.WithError(err).Error("fail")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// LoadWebsites godoc
// @Summary Load a csv file with websites to merge with companies data
// @Description post website file to merge with companies
// @ID post-load-websites
// @accept mpfd
// @Produce  plain
// @Success 200 {string} string "OK"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /companies/websites [get]
func (c companyController) LoadWebsites(ctx *gin.Context) {
	fileheader, err := ctx.FormFile("data")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	file, err := fileheader.Open()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = c.service.loadWebsites(file)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
	return
}

func (c companyController) InitDatabase(file string) {
	c.service.InitDatabase(file)
}
