package builder

import (
	"github.com/c8121/website-builder/internal/entity"
)

type TreeItem struct {
	CurrentPage *entity.Page //Current page the builder is currently working on
	Page        *entity.Page //Page this item is linked to
	Name        string
	PathName    string
	Path        string
	Children    []*TreeItem
	IsHomePage  bool
}

type PageTree struct {
	Page   *entity.Page
	Config *BuildConfig
}

func NewPageTree(page *entity.Page, config *BuildConfig) *PageTree {
	return &PageTree{
		Page:   page,
		Config: config,
	}
}

func (p *PageTree) CreatePageTree() (*TreeItem, error) {
	return p.LoadPageTree(&entity.Page{})
}

func (p *PageTree) LoadPageTree(start *entity.Page) (*TreeItem, error) {

	children, err := entity.ListPages(start.Id)
	if err != nil {
		return nil, err
	}

	childTreeItems := make([]*TreeItem, len(children))
	for i, child := range children {
		childTreeItem, err := p.LoadPageTree(&child)
		if err != nil {
			return nil, err
		}
		childTreeItems[i] = childTreeItem
	}

	path, err := GetPageDirectoryPath(start)
	if err != nil {
		path = ""
		//return nil, err
	}

	item := &TreeItem{
		CurrentPage: p.Page,
		Page:        start,
		Name:        start.Name,
		PathName:    GetPageDirectoryName(*start),
		Path:        path,
		Children:    childTreeItems,
		IsHomePage:  p.Config.IsHomePage(start),
	}

	return item, nil
}
