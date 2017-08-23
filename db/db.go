package db

import (
	"database/sql"
	"fmt"

	"github.com/domego/gorp"
)

// DBConfig db config
type DBConfig struct {
	Database     string `yaml:"database"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// MongoDBConfig mongo db config
type MongoDBConfig struct {
	Address   []string `yaml:"address"`
	Database  string   `yaml:"database"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	PoolLimit int      `yaml:"pool_limit"`
}

func (cfg *DBConfig) Source() string {
	pwd := cfg.Password
	if pwd != "" {
		pwd = ":" + pwd
	}
	port := cfg.Port
	if port == 0 {
		port = 3306
	}
	return fmt.Sprintf("%s%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&interpolateParams=true", cfg.Username, pwd, cfg.Host, port, cfg.Database)
}

// SourceUTF8 备注旧版本数据库使用UTF8编码,以及时区用UTC 0时区,跟本地时区相差8小时,这里要特别注意
func (cfg *DBConfig) SourceUTF8() string {
	pwd := cfg.Password
	if pwd != "" {
		pwd = ":" + pwd
	}
	port := cfg.Port
	if port == 0 {
		port = 3306
	}
	return fmt.Sprintf("%s%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&interpolateParams=true", cfg.Username, pwd, cfg.Host, port, cfg.Database)
}

func (cfg *DBConfig) SerialSource() string {
	return cfg.Source() + "&tx_isolation=SERIALIZABLE"
}

type DBInfo struct {
	*gorp.DbMap

	DB *sql.DB
}

func (dbInfo *DBInfo) Begin(trans *Transaction, objects ...TransactionObject) (*Transaction, error) {
	if trans == nil {
		_trans, err := dbInfo.DbMap.Begin()
		if err != nil {
			return nil, err
		}
		trans = &Transaction{Transaction: _trans, objects: []TransactionObject{}}
	}
	for _, o := range objects {
		trans.AddObject(o)
	}
	return trans, nil
}

type TransactionObject interface {
	SetTransaction(trans *Transaction)
}

type Transaction struct {
	*gorp.Transaction
	objects []TransactionObject
}

func (trans *Transaction) AddObject(o TransactionObject) {
	trans.objects = append(trans.objects, o)
	o.SetTransaction(trans)
}

func (trans *Transaction) Commit() error {
	err := trans.Transaction.Commit()
	for _, o := range trans.objects {
		o.SetTransaction(nil)
	}
	return err
}

func (trans *Transaction) Rollback() error {
	err := trans.Transaction.Rollback()
	for _, o := range trans.objects {
		o.SetTransaction(nil)
	}
	return err
}
