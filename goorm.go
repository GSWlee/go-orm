package main

import (
	"./dialect"
	"./log"
	"./session"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(dirver string, source string) (e *Engine, err error) {
	db, err := sql.Open(dirver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(dirver)
	if !ok {
		log.Errorf("dialect %s Not Found", dirver)
		return
	}
	log.Info("Connect database success")
	e = &Engine{db: db, dialect: dial}
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(session *session.Session) (interface{}, error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}
