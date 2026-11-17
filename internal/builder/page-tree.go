package builder

import (
	"log/slog"

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
	Page              *entity.Page
	Config            *BuildConfig
	IncludeVisibility int
}

func NewPageTree(page *entity.Page, config *BuildConfig, visibility int) *PageTree {
	return &PageTree{
		Page:              page,
		Config:            config,
		IncludeVisibility: visibility,
	}
}

func (p *PageTree) CreatePageTree() (*TreeItem, error) {
	return p.LoadPageTree(&entity.Page{})
}

func (p *PageTree) LoadPageTree(start *entity.Page) (*TreeItem, error) {

	children, err := entity.ListPages(entity.NewSession(), start.Id)
	if err != nil {
		return nil, err
	}

	childTreeItems := make([]*TreeItem, 0)
	for _, child := range children {

		if child.Visibility != 0 && (child.Visibility&p.IncludeVisibility == 0) {
			slog.Info("Ignore page", "page", child.Id, "visibility", child.Visibility)
			continue
		}

		childTreeItem, err := p.LoadPageTree(&child)
		if err != nil {
			return nil, err
		}
		childTreeItems = append(childTreeItems, childTreeItem)
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
