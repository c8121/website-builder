package rest_api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_vue_spa/util"
	"github.com/c8121/website-builder/internal/entity"
	"github.com/gin-gonic/gin"
)

func GetSiteConfig(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.AbortWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	siteConfig := &entity.SiteConfig{}
	err = go_db_entities.Load(entity.DbContext, siteConfig, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "failed to load site", err)
		return
	}

	c.IndentedJSON(http.StatusOK, siteConfig)
}

func ListSiteConfig(c *gin.Context) {

	items, err := entity.ListSiteConfig()
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to load sites", err)
		return
	}

	if len(items) > 0 {
		c.IndentedJSON(http.StatusOK, items)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}

func AddSiteConfig(c *gin.Context) {

	req := &entity.SiteConfig{}
	err := c.BindJSON(req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to parse request", err)
		return
	}

	req.Sort, err = findValidSiteConfigSort(req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to sort", err)
		return
	}

	err = go_db_entities.Save(entity.DbContext, req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to save entity", err)
		return
	}

	c.IndentedJSON(http.StatusOK, req)
}

func findValidSiteConfigSort(siteConfig *entity.SiteConfig) (int64, error) {

	// New sites are appended to the end of their siblings
	if siteConfig.Id == 0 {
		return nextSiteConfigSort()
	}

	// Sort is managed by the reorder-endpoint, do not trust the client's (possibly stale) value.
	existing := &entity.SiteConfig{}
	err := go_db_entities.Load(entity.DbContext, existing, *go_db_entities.NewIdFilter(siteConfig.Id))
	if err != nil {
		return 0, err
	}
	return existing.Sort, nil
}

// nextSiteConfigSort returns the Sort value for a site appended to the end of sites liste.
func nextSiteConfigSort() (int64, error) {

	siteConfigs, err := entity.ListSiteConfig()
	if err != nil {
		return 0, err
	}

	if len(siteConfigs) == 0 {
		return 1, nil
	}
	return siteConfigs[len(siteConfigs)-1].Sort + 1, nil
}

type reorderSiteConfigsRequest struct {
	Ids []int64
}

// ReorderSiteConfigs sets the order of the sites to the order of Ids.
func ReorderSiteConfigs(c *gin.Context) {

	req := &reorderSiteConfigsRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	siteConfigs, err := entity.ListSiteConfig()
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to load configs", err)
		return
	}

	err = entity.Sort(siteConfigs, req.Ids)
	if errors.Is(err, entity.ErrSortMismatch) {
		util.AbortWithStatus(c, http.StatusBadRequest, "Ids do not match the existing configs", err)
		return
	}
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to sort configs", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &util.Response{
		Status:  "ok",
		Message: "Pages sorted",
	})
}

func DeleteSiteConfig(c *gin.Context) {

	id, err := strconv.ParseInt(c.PostForm("id"), 10, 64)
	if err != nil {
		util.AbortWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	siteConfig := &entity.SiteConfig{}
	err = go_db_entities.Load(entity.DbContext, siteConfig, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "failed to load site", err)
		return
	}

	err = go_db_entities.Delete(entity.DbContext, siteConfig)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "failed to delete site config", err)
	}

	c.IndentedJSON(http.StatusOK, &util.Response{
		Status:  "ok",
		Message: "Site deleted",
	})
}
