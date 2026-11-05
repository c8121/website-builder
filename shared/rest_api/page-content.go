package rest_api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_vue_spa/util"
	"github.com/c8121/website-builder/internal/entity"
	"github.com/gin-gonic/gin"
)

func GetPageContent(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.AboutWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	page := &entity.PageContent{}
	err = go_db_entities.Load(entity.DbContext, page, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AboutWithStatus(c, http.StatusNotFound, "failed to load content", err)
		return
	}

	c.IndentedJSON(http.StatusOK, page)
}

func ListPageContents(c *gin.Context) {

	pageId, err := strconv.ParseInt(c.Param("page"), 10, 64)
	if err != nil {
		pageId = 0
	}

	items, err := entity.ListPageContents(pageId)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		util.AboutWithStatus(c, http.StatusInternalServerError, "Failed to load content", err)
		return
	}

	if len(items) > 0 {
		c.IndentedJSON(http.StatusOK, items)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}

func AddPageContent(c *gin.Context) {

	req := &entity.PageContent{}
	err := c.BindJSON(req)
	if err != nil {
		util.AboutWithStatus(c, http.StatusInternalServerError, "Failed to parse request", err)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		util.AboutWithStatus(c, http.StatusBadRequest, "Empty name not allowed", nil)
		return
	}

	if req.PageId == 0 {
		util.AboutWithStatus(c, http.StatusBadRequest, "PageId == 0", nil)
		return
	}

	err = go_db_entities.Save(entity.DbContext, req)
	if err != nil {
		util.AboutWithStatus(c, http.StatusInternalServerError, "Failed to save entity", err)
		return
	}

	c.IndentedJSON(http.StatusOK, req)
}

func DeletePageContent(c *gin.Context) {

	id, err := strconv.ParseInt(c.PostForm("id"), 10, 64)
	if err != nil {
		util.AboutWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	page := &entity.PageContent{}
	err = go_db_entities.Load(entity.DbContext, page, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AboutWithStatus(c, http.StatusNotFound, "failed to load page", err)
		return
	}

	err = go_db_entities.Delete(entity.DbContext, page)
	if err != nil {
		util.AboutWithStatus(c, http.StatusInternalServerError, "failed to delete page", err)
	}

	c.IndentedJSON(http.StatusOK, &util.Response{
		Status:  "ok",
		Message: "Page deleted",
	})
}
