package entity

import (
	"errors"
	"fmt"
	"slices"

	"github.com/c8121/go_db_entities"
)

type Page struct {
	Id       int64
	ParentId int64  `dbe:"Index"`
	Name     string `dbe:"Size: 1024, Index"`
	Sort     int64
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

func (p *Page) GetSort() int64 {
	return p.Sort
}

func (p *Page) SetSort(sort int64) {
	p.Sort = sort
}

func ListPages(parentId int64) ([]Page, error) {

	filter := go_db_entities.NewFilter("ParentId", go_db_entities.CompareEqual, parentId)
	filter.OrderBy = []string{"Sort", "Id"}

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

func CheckForCycle(page *Page) error {
	var children []Page
	err := LoadAllChildren(page.Id, &children)
	if err != nil {
		return err
	}
	for _, child := range children {
		if child.Id == page.ParentId {
			return fmt.Errorf("cyclic reference: %d", page.Id)
		}
	}
	return nil
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

func LoadPagePath(pageId int64) ([]Page, error) {

	path := make([]Page, 0)

	var page *Page

	for page == nil || pageId > 0 {

		page = &Page{}
		err := go_db_entities.Load(DbContext, page, *go_db_entities.NewIdFilter(pageId))
		if err != nil {
			return nil, err
		}

		path = append(path, *page)
		pageId = page.ParentId
	}

	slices.Reverse(path)
	return path, nil

}

func WalkPages(handler func(page *Page)) error {
	return walkPages(0, handler)
}

func walkPages(parentId int64, handler func(page *Page)) error {
	children, err := ListPages(parentId)
	if err != nil {
		return err
	}

	for _, child := range children {
		handler(&child)
		walkPages(child.Id, handler)
	}

	return nil
}
