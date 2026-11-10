package entity

import (
	"errors"

	"github.com/c8121/go_db_entities"
)

type SiteConfig struct {
	Id    int64
	Name  string
	Value string
	Sort  int64
}

func (s *SiteConfig) GetEntityName() string {
	return "Site"
}

func (s *SiteConfig) GetId() int64 {
	return s.Id
}

func (s *SiteConfig) SetId(id int64) {
	s.Id = id
}

func (s *SiteConfig) GetSort() int64 {
	return s.Sort
}

func (s *SiteConfig) SetSort(sort int64) {
	s.Sort = sort
}

func ListSiteConfig() ([]SiteConfig, error) {

	filter := &go_db_entities.Filter{}
	filter.OrderBy = []string{"Sort", "Id"}

	siteConfigSupplier := func() go_db_entities.Entity {
		return &SiteConfig{}
	}

	items, err := go_db_entities.Select(DbContext, siteConfigSupplier, *filter)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return nil, err
	}

	siteConfigs := make([]SiteConfig, len(items))
	for i, item := range items {
		siteConfig := item.(*SiteConfig)
		siteConfigs[i] = *siteConfig
	}

	return siteConfigs, nil
}

func SiteConfigAsMap(siteConfigs []SiteConfig) map[string]string {
	config := map[string]string{}
	for _, siteConfig := range siteConfigs {
		config[siteConfig.Name] = siteConfig.Value
	}
	return config
}
