package goutils

import (
	"fmt"
	"strings"
)

func SqlMergeTable(src, desc string) string {
	srcTable := strmysqlparseTable(src)
	descTable := strmysqlparseTable(desc)
	fields := []string{}
	for k, item := range srcTable.Fields {
		old, ok := descTable.Fields[k]
		if !ok { //新建
			fields = append(fields, buildsqlline("ADD", k, item, srcTable, descTable))
		} else {
			if old.Comment != item.Comment || old.Type != item.Type || old.Default != item.Default || old.IsNull != item.IsNull || old.Ext != item.Ext || old.Pre != item.Pre { //修改
				fields = append(fields, buildsqlline("MODIFY", k, item, srcTable, descTable))
			}
		}
	}
	for k, _ := range descTable.Fields {
		_, ok := srcTable.Fields[k]
		if !ok {
			fields = append(fields, fmt.Sprintf("DROP COLUMN `%s`", k))
		}
	}
	for k, item := range srcTable.Indexs {
		_, ok := descTable.Indexs[k]
		if !ok {
			fields = append(fields, fmt.Sprintf("ADD %s `%s`(`%s`)", item.OpType, k, strings.Join(item.Keys, "`,`")))
		}
	}
	for k, _ := range descTable.Fields {
		item, ok := srcTable.Indexs[k]
		if !ok {
			fields = append(fields, fmt.Sprintf("DROP %s `%s`", item.OpType, k))
		}
	}
	return fmt.Sprintf("ALTER TABLE `%s` %s", srcTable.Name, strings.Join(fields, ","))
}
func getalterAfter(k string, src, old strmysqltable) string {
	item, ok := old.Fields[k]
	if ok {
		if item.Pre != "" {
			return fmt.Sprintf("AFTER `%s`", item.Pre)
		}
	}
	return ""
}
func buildsqlline(op, key string, desc strmysqltablefield, src, old strmysqltable) string {
	nullDesc := "NOT NULL"
	defval := fmt.Sprintf("DEFAULT '%s'", desc.Default)
	comment := ""
	if desc.Comment != "" {
		comment = fmt.Sprintf("COMMENT '%s'", desc.Comment)
	}
	if desc.IsNull {
		nullDesc = ""
		defval = ""
	}

	return fmt.Sprintf("%s COLUMN `%s` %s %s %s %s %s %s", op, key, strings.ToUpper(desc.Type), nullDesc, defval, desc.Ext, comment, getalterAfter(key, src, old))
}

func strmysqlparseTable(data string) strmysqltable {
	tableName := ""
	fields := map[string]strmysqltablefield{}
	return strmysqltable{
		Name:   tableName,
		Fields: fields,
	}
}

type strmysqltablefield struct {
	Type    string
	IsNull  bool
	Default string
	Ext     string
	Comment string
	Pre     string
}
type strmysqltableindex struct {
	OpType string
	Name   string
	Keys   []string
}
type strmysqltable struct {
	Name   string
	Fields map[string]strmysqltablefield
	Indexs map[string]strmysqltableindex
}
