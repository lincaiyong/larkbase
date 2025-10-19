package larkbase

import (
	"github.com/lincaiyong/larkbase/larkfield"
	"github.com/lincaiyong/larkbase/larksuite/bitable"
)

type Filter = bitable.FilterInfo
type ViewFilter = bitable.AppTableViewPropertyFilterInfo
type Condition = larkfield.Condition

func NewViewIdFindOption(viewId string) *FindOption {
	return &FindOption{viewId: viewId}
}

func NewFindOption(filter *bitable.FilterInfo, sorts ...*bitable.Sort) *FindOption {
	return &FindOption{filter: filter, sorts: sorts}
}

type FindOption struct {
	filter *bitable.FilterInfo
	sorts  []*bitable.Sort
	limit  int
	viewId string
}

func (o *FindOption) Limit(limit int) *FindOption {
	o.limit = limit
	return o
}

func (c *Connection[T]) Sort() *T {
	return c.condition
}

func (c *Connection[T]) Condition() *T {
	return c.condition
}

func (c *Connection[T]) FilterAnd(conditions ...*Condition) *Filter {
	s := make([]*bitable.Condition, len(conditions))
	for i, condition := range conditions {
		s[i] = condition.ToLarkCondition()
	}
	return bitable.NewFilterInfoBuilder().
		Conjunction(`and`).
		Conditions(s).Build()
}

func (c *Connection[T]) FilterOr(conditions ...*Condition) *Filter {
	s := make([]*bitable.Condition, len(conditions))
	for i, condition := range conditions {
		s[i] = condition.ToLarkCondition()
	}
	return bitable.NewFilterInfoBuilder().
		Conjunction(`or`).
		Conditions(s).Build()
}

func (c *Connection[T]) ViewFilterAnd(conditions ...*Condition) *ViewFilter {
	s := make([]*bitable.AppTableViewPropertyFilterInfoCondition, len(conditions))
	for i, condition := range conditions {
		s[i] = condition.ToLarkViewCondition()
	}
	return bitable.NewAppTableViewPropertyFilterInfoBuilder().
		Conjunction(`and`).
		Conditions(s).Build()
}

func (c *Connection[T]) ViewFilterOr(conditions ...*Condition) *ViewFilter {
	s := make([]*bitable.AppTableViewPropertyFilterInfoCondition, len(conditions))
	for i, condition := range conditions {
		s[i] = condition.ToLarkViewCondition()
	}
	return bitable.NewAppTableViewPropertyFilterInfoBuilder().
		Conjunction(`or`).
		Conditions(s).Build()
}
