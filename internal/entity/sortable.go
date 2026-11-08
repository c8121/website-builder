package entity

import (
	"errors"

	"github.com/c8121/go_db_entities"
)

var (
	// ErrSortMismatch is returned if a list of ids passed for sorting does not
	// match the existing siblings (unknown, missing or duplicate ids).
	ErrSortMismatch = errors.New("ids do not match the existing siblings")

	ErrNoSortableEntity = errors.New("not an entity")
)

type Sortable interface {
	GetId() int64
	GetSort() int64
	SetSort(int64)
}

// Sort sets the order of all items to the order of ids.
// ids must contain exactly the ids of the existing children, each once.
func Sort[T any, PT interface {
	*T
	go_db_entities.Entity
}](items []T, ids []int64) error {

	if len(ids) != len(items) {
		return ErrSortMismatch
	}

	byId := make(map[int64]go_db_entities.Entity, len(items))
	for i := range items {
		var ptr PT = &items[i]
		byId[ptr.GetId()] = ptr
	}

	for i, id := range ids {
		item, ok := byId[id]
		if !ok { // unknown or duplicate
			return ErrSortMismatch
		}
		delete(byId, id)

		sortable, ok := item.(Sortable)
		if !ok {
			return ErrNoSortableEntity
		}
		sortable.SetSort(int64(i + 1))

		err := go_db_entities.Save(DbContext, item)
		if err != nil {
			return err
		}
	}

	return nil
}
