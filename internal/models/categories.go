package models

import (
	"gopkg.in/reform.v1"
	"strings"
)

//go:generate reform

type NewsCategories struct {
	news_Id     int64 `reform:"News_Id,pk"`
	category_id int64 `reform:"Category_Id,pk"`
}

func (n NewsCategories) String() string {
	res := make([]string, 6)
	res[0] = "news_Id: " + reform.Inspect(n.news_Id, true)
	res[1] = "category_id: " + reform.Inspect(n.category_id, true)
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
