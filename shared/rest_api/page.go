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

func GetPage(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.AboutWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	page := &entity.Page{}
	err = go_db_entities.Load(entity.DbContext, page, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AboutWithStatus(c, http.StatusNotFound, "failed to load page", err)
		return
	}

	c.IndentedJSON(http.StatusOK, page)
}

func ListPages(c *gin.Context) {

	parentId, err := strconv.ParseInt(c.Param("parent"), 10, 64)
	if err != nil {
		parentId = 0
	}

	items, err := entity.ListPages(parentId)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		util.AboutWithStatus(c, http.StatusInternalServerError, "Failed to load pages", err)
		return
	}

	if len(items) > 0 {
		c.IndentedJSON(http.StatusOK, items)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}

func AddPage(c *gin.Context) {

	req := &entity.Page{}
	err := c.BindJSON(req)
	if err != nil {
		util.AboutWithStatus(c, http.StatusInternalServerError, "Failed to parse request", err)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		util.AboutWithStatus(c, http.StatusBadRequest, "Empty name not allowed", nil)
		return
	}

	if req.Id != 0 && req.Id == req.ParentId {
		util.AboutWithStatus(c, http.StatusBadRequest, "ParentId == Id", nil)
		return
	}

	if req.Id != 0 {
		children := make([]entity.Page, 0)
		err := entity.LoadAllChildren(req.Id, &children)
		if err != nil {
			util.AboutWithStatus(c, http.StatusInternalServerError, "Failed to check for exiting children", err)
			return
		}
		for _, child := range children {
			if child.Id == req.ParentId {
				util.AboutWithStatus(c, http.StatusBadRequest, "cyclic reference", nil)
				return
			}
		}
	}

	err = go_db_entities.Save(entity.DbContext, req)
	if err != nil {
		util.AboutWithStatus(c, http.StatusInternalServerError, "Failed to save entity", err)
		return
	}

	c.IndentedJSON(http.StatusOK, req)
}
