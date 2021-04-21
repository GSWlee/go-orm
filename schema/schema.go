package schema

import (
	"../dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model     interface{}
	Name      string
	Fields    []*Field
	FieldName []string
	fieldMap  map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType:=reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema:=&Schema{
		Model:     dest,
		Name:      modelType.Name(),
		fieldMap:  make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("goorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldName = append(schema.FieldName, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue:=reflect.Indirect(reflect.ValueOf(dest))
	var vars []interface{}
	for _,field:=range s.Fields{
		vars=append(vars,destValue.FieldByName(field.Name).Interface())
	}
	return vars
}