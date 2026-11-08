package builder

import (
	"github.com/c8121/website-builder/internal/entity"
)

type TreeItem struct {
	Page       *entity.Page
	Name       string
	PathName   string
	Path       string
	Children   []*TreeItem
	IsHomePage bool
}

type PageTree struct {
	Config *BuildConfig
}

func NewPageTree(config *BuildConfig) *PageTree {
	return &PageTree{
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
		Page:       start,
		Name:       start.Name,
		PathName:   GetPageDirectoryName(*start),
		Path:       path,
		Children:   childTreeItems,
		IsHomePage: p.Config.IsHomePage(start),
	}

	return item, nil
}
