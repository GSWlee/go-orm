package clause

import (
	"fmt"
	"strings"
)

type generator func(value ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

func genBindVars(num int)string{
	var vars []string
	for i:=0;i<num;i++{
		vars=append(vars,"?")
	}
	return strings.Join(vars,", ")
}

func _insert(value ...interface{}) (string,[]interface{}) {
	tableName:=value[0]
	field :=strings.Join(value[1].([]string),",")
	return fmt.Sprintf("INSERT INTO %s (%v)",tableName,field), []interface{}{}
}

func _values(values ...interface{}) (string,[]interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i,value:=range values{
		v:=value.([]interface{})
		if bindStr==""{
			bindStr=genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)",bindStr))
		if i+1!=len(values){
			sql.WriteString(", ")
		}
		vars=append(vars,v...)
	}
	return sql.String(),vars
}

func _select(value ...interface{}) (string,[]interface{}) {
	tableName:=value[0]
	field:=strings.Join(value[1].([]string),",")
	return fmt.Sprintf("SELECT %v FROM %s",field,tableName),[]interface{}{}
}

func _limit(values ...interface{}) (string,[]interface{}) {
	return "LIMIT ?",values
}

func _where(values ...interface{}) (string,[]interface{}) {
	desc,vars:=values[0],values[1:]
	return fmt.Sprintf("WHERE %s",desc),vars
}

func _orderBy(values ...interface{}) (string,[]interface{}) {
	return fmt.Sprintf("ORDER BY %s",values[0]),[]interface{}{}
}