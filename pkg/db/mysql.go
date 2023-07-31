package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/harley9293/blotlog"
	"github.com/jmoiron/sqlx"
	"time"
)

type MysqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type MysqlClient struct {
	sourceName string
	conn       *sqlx.DB
}

func NewMysqlClient(config *MysqlConfig) *MysqlClient {
	client := &MysqlClient{}
	client.sourceName = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?readTimeout=5s&writeTimeout=5s", config.User, config.Password, config.Host, config.Port, config.DBName)
	client.connect()
	return client
}

func (c *MysqlClient) connect() {
	conn, err := sqlx.Open("mysql", c.sourceName)
	if err != nil {
		log.Error("DB Connect Error, err:%s", err)
	}
	c.conn = conn
	err = c.conn.Ping()
	if err != nil {
		log.Error("DB Ping Error, err:%s", err)
	}

	c.conn.SetConnMaxLifetime(time.Minute * 3)
}

func (c *MysqlClient) Get(dest any, query string, args ...any) {
	err := c.conn.Get(dest, query, args...)
	if err != nil && err != sql.ErrNoRows {
		log.Error("DB Get Error, err:%s, query:%s, args:%+v", err, query, args)
	}
}

func (c *MysqlClient) Insert(query string, args ...any) int64 {
	result, err := c.conn.Exec(query, args...)
	if err != nil {
		log.Error("DB Insert Error, err:%s, query:%s, args:%+v", err, query, args)
		return 0
	}

	id, _ := result.LastInsertId()
	return id
}

func (c *MysqlClient) Query(query string, args ...any) *sql.Rows {
	rows, err := c.conn.Query(query, args...)
	if err != nil {
		log.Error("DB Query Error, err:%s, query:%s, args:%+v", err, query, args)
		return nil
	}
	return rows
}

func (c *MysqlClient) Update(query string, args ...any) {
	_, err := c.conn.Exec(query, args...)
	if err != nil {
		log.Error("DB Update Error, err:%s, query:%s, args:%+v", err, query, args)
	}
}

func (c *MysqlClient) Delete(query string, args ...any) {
	_, err := c.conn.Exec(query, args...)
	if err != nil {
		log.Error("DB Delete Error, err:%s, query:%s, args:%+v", err, query, args)
	}
}

func (c *MysqlClient) KeepAlive() {
	err := c.conn.Ping()
	if err != nil {
		log.Error("DB Ping Error err:%s, Try Reconnect", err)
		c.Shutdown()
		c.connect()
	}
}

func (c *MysqlClient) Shutdown() {
	err := c.conn.Close()
	if err != nil {
		log.Error("mysql conn close error, err:%s", err.Error())
	}
}
