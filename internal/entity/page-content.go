package entity

import (
	"errors"

	"github.com/c8121/go_db_entities"
)

type PageContent struct {
	Id          int64
	PageId      int64  `dbe:"Index"`
	Name        string `dbe:"Size: 1024, Index"`
	ContentType string
	Content     string
	Sort        int64
}

func (p *PageContent) GetEntityName() string {
	return "PageContent"
}

func (p *PageContent) GetId() int64 {
	return p.Id
}

func (p *PageContent) SetId(id int64) {
	p.Id = id
}

func ListPageContents(parentId int64) ([]PageContent, error) {

	filter := go_db_entities.NewFilter("PageId", go_db_entities.CompareEqual, parentId)
	filter.OrderBy = []string{"Sort"}

	pageSupplier := func() go_db_entities.Entity {
		return &PageContent{}
	}

	items, err := go_db_entities.Select(DbContext, pageSupplier, *filter)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return nil, err
	}

	contents := make([]PageContent, len(items))
	for i, item := range items {
		content := item.(*PageContent)
		contents[i] = *content
	}

	return contents, nil

}
