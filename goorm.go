package GO_ORM

import (
	"./log"
	"./session"
	"./dialect"
	"database/sql"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
}

func NewEngine(dirver string,source string) (e *Engine,err error) {
	db,err:=sql.Open(dirver,source)
	if err!=nil{
		log.Error(err)
		return
	}

	if err=db.Ping();err!=nil{
		log.Error(err)
		return
	}

	dial,ok:=dialect.GetDialect(dirver)
	if !ok{
		log.Errorf("dialect %s Not Found",dirver)
		return
	}
	log.Info("Connect database success")
	e=&Engine{db: db,dialect: dial}
	return
}

func (e *Engine) Close()  {
	if err:=e.db.Close();err!=nil{
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db,e.dialect)
}
