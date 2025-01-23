package models

import (
	"gopkg.in/reform.v1"
	"strings"
)

//go:generate reform

type NewsCategories struct {
	NewsID     int64 `reform:"News_Id,pk"`
	CategoryID int64 `reform:"Category_Id,pk"`
}

func (n NewsCategories) String() string {
	res := make([]string, 6)
	res[0] = "NewsID: " + reform.Inspect(n.NewsID, true)
	res[1] = "CategoryID: " + reform.Inspect(n.CategoryID, true)
	return strings.Join(res, ", ")
}

func (n NewsCategories) Values() []interface{} {
	//TODO implement me
	panic("implement me")
}

func (n NewsCategories) Pointers() []interface{} {
	//TODO implement me
	panic("implement me")
}

func (n NewsCategories) View() reform.View {
	//TODO implement me
	panic("implement me")
}
