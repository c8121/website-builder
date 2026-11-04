package entity

import (
	"errors"
	"fmt"

	"github.com/c8121/go_db_entities"
)

type Page struct {
	Id          int64
	ParentId    int64  `dbe:"Index"`
	Name        string `dbe:"Size: 1024, Index"`
	ContentType string
	Content     string
	Sort        int64
}

func (p *Page) GetEntityName() string {
	return "Page"
}

func (p *Page) GetId() int64 {
	return p.Id
}

func (p *Page) SetId(id int64) {
	p.Id = id
}

func ListPages(parentId int64) ([]Page, error) {

	filter := go_db_entities.NewFilter("ParentId", go_db_entities.CompareEqual, parentId)
	filter.OrderBy = []string{"Sort"}

	pageSupplier := func() go_db_entities.Entity {
		return &Page{}
	}

	items, err := go_db_entities.Select(DbContext, pageSupplier, *filter)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return nil, err
	}

	pages := make([]Page, len(items))
	for i, item := range items {
		page := item.(*Page)
		pages[i] = *page
	}

	return pages, nil

}

func LoadAllChildren(parentId int64, into *[]Page) error {

	if into == nil {
		return fmt.Errorf("into-parameter is nil")
	}

	children, err := ListPages(parentId)
	if err != nil {
		return err
	}

	for _, child := range children {
		*into = append(*into, child)
		err = LoadAllChildren(child.Id, into)
		if err != nil {
			return err
		}
	}

	return nil
}
