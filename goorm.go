package GO_ORM

import (
	"./log"
	"./session"
	"database/sql"
)

type Engine struct {
	db *sql.DB
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
	log.Info("Connect database success")
	e=&Engine{db: db}
	return
}

func (e *Engine) Close()  {
	if err:=e.db.Close();err!=nil{
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}
