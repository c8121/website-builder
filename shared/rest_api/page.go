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
		util.AbortWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	page := &entity.Page{}
	err = go_db_entities.Load(entity.NewSession(), page, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "failed to load page", err)
		return
	}

	c.IndentedJSON(http.StatusOK, page)
}

func ListPages(c *gin.Context) {

	parentId, err := strconv.ParseInt(c.Param("parent"), 10, 64)
	if err != nil {
		parentId = 0
	}

	items, err := entity.ListPages(entity.NewSession(), parentId)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to load pages", err)
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
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to parse request", err)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		util.AbortWithStatus(c, http.StatusBadRequest, "Empty name not allowed", nil)
		return
	}

	if req.Id != 0 && req.Id == req.ParentId {
		util.AbortWithStatus(c, http.StatusBadRequest, "ParentId == Id", nil)
		return
	}

	session, tx, err := entity.NewSessionTx()
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to create transaction", err)
		return
	}
	defer tx.Rollback() // no effect after Commit (returns sql.ErrTxDone, which can be ignored)

	if req.Id != 0 {
		err := entity.CheckForCycle(session, req)
		if err != nil {
			util.AbortWithStatus(c, http.StatusBadRequest, "cyclic reference", err)
			return
		}
	}

	req.Sort, err = findValidPageSort(session, req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to sort", err)
		return
	}

	err = go_db_entities.Save(session, req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to save entity", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	c.IndentedJSON(http.StatusOK, req)
}

func findValidPageSort(session *go_db_entities.Session, page *entity.Page) (int64, error) {

	// New pages are appended to the end of their siblings
	if page.Id == 0 {
		return nextPageSort(session, page.ParentId)
	}

	// Sort is managed by the reorder-endpoint, do not trust the client's (possibly stale) value.
	// A page moved to another parent is appended to the end of the new siblings.
	existing := &entity.Page{}
	err := go_db_entities.Load(session, existing, *go_db_entities.NewIdFilter(page.Id))
	if err != nil {
		return 0, err
	}
	if existing.ParentId == page.ParentId {
		return existing.Sort, nil
	}

	return nextPageSort(session, page.ParentId)
}

// nextPageSort returns the Sort value for a page appended to the end of parentId's children.
func nextPageSort(session *go_db_entities.Session, parentId int64) (int64, error) {

	pages, err := entity.ListPages(session, parentId)
	if err != nil {
		return 0, err
	}

	if len(pages) == 0 {
		return 1, nil
	}
	return pages[len(pages)-1].Sort + 1, nil
}

type reorderPagesRequest struct {
	ParentId int64
	Ids      []int64
}

// ReorderPages sets the order of the children of ParentId to the order of Ids.
func ReorderPages(c *gin.Context) {

	req := &reorderPagesRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	session, tx, err := entity.NewSessionTx()
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to create transaction", err)
		return
	}
	defer tx.Rollback() // no effect after Commit (returns sql.ErrTxDone, which can be ignored)

	pages, err := entity.ListPages(session, req.ParentId)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to load pages", err)
		return
	}

	err = entity.Sort(session, pages, req.Ids)
	if errors.Is(err, entity.ErrSortMismatch) {
		util.AbortWithStatus(c, http.StatusBadRequest, "Ids do not match the existing pages", err)
		return
	}
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to sort pages", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &util.Response{
		Status:  "ok",
		Message: "Pages sorted",
	})
}

func DeletePage(c *gin.Context) {

	id, err := strconv.ParseInt(c.PostForm("id"), 10, 64)
	if err != nil {
		util.AbortWithStatus(c, http.StatusBadRequest, "missing id", err)
		return
	}

	session := entity.NewSession()

	page := &entity.Page{}
	err = go_db_entities.Load(session, page, *go_db_entities.NewIdFilter(id))
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "failed to load page", err)
		return
	}

	err = go_db_entities.Delete(session, page)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "failed to delete page", err)
	}

	c.IndentedJSON(http.StatusOK, &util.Response{
		Status:  "ok",
		Message: "Page deleted",
	})
}
