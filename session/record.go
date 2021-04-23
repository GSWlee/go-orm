package session

import (
	"../clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64,error) {
	recordValues:=make([]interface{},0)
	for _,value:=range values{
		s.CallMethod(BeforeInsert,value)
		table:=s.Model(value).refTable
		s.clause.Set(clause.INSERT,table.Name,table.FieldName)
		recordValues=append(recordValues,table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES,recordValues...)
	sql,vars:=s.clause.Build(clause.INSERT,clause.VALUES)

	result,err:=s.Raw(sql,vars...).Exec()
	if err!=nil{
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(value interface{}) error {
	destSlice:=reflect.Indirect(reflect.ValueOf(value))
	destType:=destSlice.Type().Elem()
	table:=s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT,table.Name,table.FieldName)
	sql,vars:=s.clause.Build(clause.SELECT,clause.WHERE,clause.ORDERBY,clause.LIMIT)
	rows,err:=s.Raw(sql,vars...).QueryRows()
	if err!=nil{
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldName {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}

		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery,dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) Update(kv ...interface{}) (int64, error) {
	m,ok:=kv[0].(map[string]interface{})
	if !ok{
		m=make(map[string]interface{})
		for i:=0;i<len(kv);i+=2{
			m[kv[i].(string)]=kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE,s.RefTable().Name,m)
	sql,vars:=s.clause.Build(clause.UPDATE,clause.WHERE)
	result,err:=s.Raw(sql,vars...).Exec()
	if err!=nil{
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64,error) {
	s.clause.Set(clause.DELETE,s.RefTable().Name)
	sql,vars:=s.clause.Build(clause.DELETE,clause.WHERE)
	result,err:=s.Raw(sql,vars...).Exec()
	if err!=nil{
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64,error) {
	s.clause.Set(clause.COUNT,s.RefTable().Name)
	sql,vars:=s.clause.Build(clause.COUNT,clause.WHERE)
	row:=s.Raw(sql,vars...).QueryRow()
	var tmp int64
	if err:=row.Scan(&tmp);err!=nil{
		return 0, err
	}
	return tmp,nil
}

func (s *Session) Where(dest string,values ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE,append(append(vars,dest),values...)...)
	return s
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT,num)
	return s
}

func (s *Session) Orderby(desc string) *Session {
	s.clause.Set(clause.ORDERBY,desc)
	return s
}

func (s *Session) First(value interface{}) error {
	dest:=reflect.Indirect(reflect.ValueOf(value))
	destslice:=reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err:=s.Limit(1).Find(destslice.Addr().Interface());err!=nil{
		return err
	}
	dest.Set(destslice.Index(0))
	return nil
}