package goutils

import (
	"fmt"
	"strings"
)

func SqlMergeTable(src, desc string) string {
	srcTable := strmysqlparseTable(src)
	descTable := strmysqlparseTable(desc)
	fmt.Printf("--->\n%v\n\ndesc -->\n%v\n\n\n", srcTable, descTable)
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
			if item.OpType == "PRIMARY KEY" {
				fields = append(fields, fmt.Sprintf("ADD %s (`%s`)", item.OpType, strings.Join(item.Keys, "`,`")))
			} else {
				fields = append(fields, fmt.Sprintf("ADD %s `%s`(`%s`)", item.OpType, k, strings.Join(item.Keys, "`,`")))
			}

		}
	}
	for k, _ := range descTable.Indexs {
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
	fields := map[string]strmysqltablefield{}
	ret := strmysqltable{
		Fields: fields,
		Indexs: map[string]strmysqltableindex{},
	}
	index, tableName := StrGetBody(data, "`", "`")
	if index < 0 {
		return ret
	}
	ret.Name = tableName
	body := data[index:]
	index = strings.Index(body, "(")
	body = body[index+1:]
	// _, body := StrGetBody(data, "(", ")")
	lines := strings.Split(body, ",")
	Pre := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "`") { //字段
			ft := strlineParsestrmysqltablefield(line)
			if ft != nil {
				ft.Pre = Pre
				Pre = ft.Name
				ret.Fields[ft.Name] = *ft
			}
		} else { //索引
			ft := strlineParsesstrmysqltableindex(line)
			if ft != nil {
				ret.Indexs[ft.Name] = *ft
			}
		}
	}
	return ret
}
func strlineParsestrmysqltablefield(src string) *strmysqltablefield {
	ret := &strmysqltablefield{
		IsNull: false,
	}
	i, name := StrGetBody(src, "`", "`")
	if i < 0 {
		return nil
	}
	src = strings.TrimSpace(src[i:])
	data := strings.Split(src, " ")
	l := len(data)
	ret.Name = name
	for i := 0; i < l; i++ {
		tag := strings.ToUpper(strings.TrimSpace(data[i]))
		switch tag {
		case "":
			break
		case "UNSIGNED":
			ret.Type += " UNSIGNED"
		case "NOT":
			if strings.ToUpper(strings.TrimSpace(data[i+1])) == "NULL" {
				i += 1
				ret.IsNull = false
			}
		case "AUTO_INCREMENT":
			ret.Ext = "AUTO_INCREMENT"
		case "DEFAULT":
			val := strings.ToUpper(strings.TrimSpace(data[i+1]))
			if strings.HasPrefix(val, "'") {
				_, ret.Default = StrGetBody(val, "'", "'")
			} else {
				ret.Default = val
			}
			i += 1
		case "ON":
			op := strings.ToUpper(strings.TrimSpace(data[i+1]))
			val := strings.ToUpper(strings.TrimSpace(data[i+2]))
			ret.Ext = fmt.Sprintf("%s %s %s", tag, op, val)
			i += 2
		case "COMMENT":
			val := strings.ToUpper(strings.TrimSpace(data[i+1]))
			_, ret.Comment = StrGetBody(val, "'", "'")
			i += 1
		default:
			if ret.Type == "" {
				ret.Type = tag
			}
		}
	}
	return ret
}
func strlineParsesstrmysqltableindex(src string) *strmysqltableindex {
	src = strings.ToLower(strings.TrimSpace(src))
	ret := &strmysqltableindex{}
	if strings.HasPrefix(src, "primary") {
		i, val := StrGetBody(src, "(`", "`)")
		if i < 0 {
			return nil
		}
		ret.Name = ""
		ret.OpType = "PRIMARY KEY"
		ret.Keys = strings.Split(val, "`,`")
	} else if strings.HasPrefix(src, "key") || strings.HasPrefix(src, "index") {
		i, val := StrGetBody(src, "(`", "`)")
		if i < 0 {
			return nil
		}
		i, name := StrGetBody(src, "`", "`")
		if i < 0 {
			return nil
		}
		ret.Name = name
		ret.OpType = "INDEX"
		ret.Keys = strings.Split(val, "`,`")
	} else if strings.HasPrefix(src, "unique") {
		i, val := StrGetBody(src, "(`", "`)")
		if i < 0 {
			return nil
		}
		i, name := StrGetBody(src, "`", "`")
		if i < 0 {
			return nil
		}
		ret.Name = name
		ret.OpType = "UNIQUE"
		ret.Keys = strings.Split(val, "`,`")
	}
	if ret.OpType == "" || len(ret.Keys) == 0 {
		return nil
	}
	return ret

}

type strmysqltablefield struct {
	Name    string
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
