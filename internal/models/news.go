package models

import (
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type newsTableType struct {
	s parse.StructInfo
	z []interface{}
}

func (n newsTableType) Schema() string {
	return n.s.SQLSchema
}

func (n newsTableType) Name() string {
	//TODO implement me
	panic("implement me")
}

func (n newsTableType) Columns() []string {
	//TODO implement me
	panic("implement me")
}

func (n newsTableType) NewStruct() reform.Struct {
	//TODO implement me
	panic("implement me")
}

type FieldInfo struct {
	Name   string // field name as defined in source file, e.g. Name
	Type   string // field type as defined in source file, e.g. string; always present for primary key, may be absent otherwise
	Column string // SQL database column name from "reform:" struct field tag, e.g. name
}

// StructInfo represents information about struct.
type StructInfo struct {
	Type         string      // struct type as defined in source file, e.g. User
	SQLSchema    string      // SQL database schema name from magic "reform:" comment, e.g. public
	SQLName      string      // SQL database view or table name from magic "reform:" comment, e.g. users
	Fields       []FieldInfo // fields info
	PKFieldIndex int         // index of primary key field in Fields, -1 if none
}

//go:generate reform

type News struct {
	ID      int64  `reform:"Id,pk"`
	Title   string `reform:"Title"`
	Content string `reform:"Content"`
}

var NewsTable = &newsTableType{
	s: parse.StructInfo{
		Type:    "News",
		SQLName: "news",
		Fields: []parse.FieldInfo{
			{Name: "Id", Type: "int32", Column: "id"},
			{Name: "Title", Type: "*int32", Column: "title"},
			{Name: "Content", Type: "string", Column: "ontent"},
		},
		PKFieldIndex: 0,
	},
	z: new(News).Values(),
}

func (n News) Schema() string {
	//TODO implement me
	panic("implement me")
}

func (n News) Name() string {
	//TODO implement me
	panic("implement me")
}

func (n News) Columns() []string {
	//TODO implement me
	panic("implement me")
}

func (n News) NewStruct() reform.Struct {
	//TODO implement me
	panic("implement me")
}

func (n News) String() string {
	//TODO implement me
	panic("implement me")
}

func (n News) Values() []interface{} {
	return []interface{}{
		n.ID,
		n.Title,
		n.Content,
	}
}

func (n News) Pointers() []interface{} {
	//TODO implement me
	panic("implement me")
}

func (n News) View() reform.View {
	//TODO implement me
	panic("implement me")
}

func (n News) Table() reform.Table {
	//TODO implement me
	panic("implement me")
}

func (n News) PKValue() interface{} {
	//TODO implement me
	panic("implement me")
}

func (n News) PKPointer() interface{} {
	//TODO implement me
	panic("implement me")
}

func (n News) HasPK() bool {
	news := new(News)
	news.ID = 2
	return true
}

func (n News) SetPK(pk interface{}) {
	//TODO implement me
	panic("implement me")
}
