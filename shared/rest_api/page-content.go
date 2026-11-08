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
		util.AbortWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	page := &entity.PageContent{}
	err = go_db_entities.Load(entity.DbContext, page, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "failed to load content", err)
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
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to load content", err)
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
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to parse request", err)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		util.AbortWithStatus(c, http.StatusBadRequest, "Empty name not allowed", nil)
		return
	}

	if req.PageId == 0 {
		util.AbortWithStatus(c, http.StatusBadRequest, "PageId == 0", nil)
		return
	}

	req.Sort, err = findValidPageContentSort(req)
	if errors.Is(err, go_db_entities.ErrNotFound) {
		util.AbortWithStatus(c, http.StatusNotFound, "content not found", err)
		return
	}
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

func findValidPageContentSort(content *entity.PageContent) (int64, error) {

	// New contents are appended to the end of the page's contents
	if content.Id == 0 {
		return nextPageContentSort(content.PageId)
	}

	// Sort is managed by the reorder-endpoint, do not trust the client's (possibly stale) value.
	// A content moved to another page is appended to the end of that page's contents.
	existing := &entity.PageContent{}
	err := go_db_entities.Load(entity.DbContext, existing, *go_db_entities.NewIdFilter(content.Id))
	if err != nil {
		return 0, err
	}
	if existing.PageId == content.PageId {
		return existing.Sort, nil
	}

	return nextPageContentSort(content.PageId)
}

// nextPageContentSort returns the Sort value for a content appended to the end of pageId's contents.
func nextPageContentSort(pageId int64) (int64, error) {

	contents, err := entity.ListPageContents(pageId)
	if err != nil {
		return 0, err
	}

	if len(contents) == 0 {
		return 1, nil
	}
	return contents[len(contents)-1].Sort + 1, nil
}

type reorderPageContentsRequest struct {
	PageId int64
	Ids    []int64
}

// ReorderPageContents sets the order of the contents of PageId to the order of Ids.
func ReorderPageContents(c *gin.Context) {

	req := &reorderPageContentsRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	if req.PageId == 0 {
		util.AbortWithStatus(c, http.StatusBadRequest, "PageId == 0", nil)
		return
	}

	contents, err := entity.ListPageContents(req.PageId)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to load contents", err)
		return
	}

	err = entity.Sort(contents, req.Ids)
	if errors.Is(err, entity.ErrSortMismatch) {
		util.AbortWithStatus(c, http.StatusBadRequest, "Ids do not match the existing contents", err)
		return
	}
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to sort contents", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &util.Response{
		Status:  "ok",
		Message: "Contents sorted",
	})
}

func DeletePageContent(c *gin.Context) {

	id, err := strconv.ParseInt(c.PostForm("id"), 10, 64)
	if err != nil {
		util.AbortWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	page := &entity.PageContent{}
	err = go_db_entities.Load(entity.DbContext, page, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "failed to load page", err)
		return
	}

	err = go_db_entities.Delete(entity.DbContext, page)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "failed to delete page", err)
	}

	c.IndentedJSON(http.StatusOK, &util.Response{
		Status:  "ok",
		Message: "Page deleted",
	})
}
