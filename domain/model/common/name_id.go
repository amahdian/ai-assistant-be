package common

import "github.com/samber/lo"

type NameId struct {
	Id   int64
	Name string
}

type NameIds []NameId

func (n NameIds) Names() []string {
	return lo.Map(n, func(item NameId, _ int) string {
		return item.Name
	})
}

func (n NameIds) Ids() []int64 {
	return lo.Map(n, func(item NameId, _ int) int64 {
		return item.Id
	})
}
