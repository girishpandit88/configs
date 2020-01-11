package postgres

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type Object struct {
	key        string
	value      map[string]string
	createdAt  time.Time
	modifiedAt time.Time
}

func (o *Object) String() string {
	return fmt.Sprintf("\nKey: %s, "+
		"\nValue: [%v\n], "+
		"\nCreatedAt: %s, "+
		"\nModifiedAt: %s",
		o.Key(),
		valueAsString(o.value),
		o.CreatedAt(),
		o.ModifiedAt(),
	)
}

func valueAsString(value map[string]string) string {
	var s string
	for k, v := range value {
		s += fmt.Sprintf("\n%v->%v", k, v)
	}
	return s
}

type DBStore struct {
	DB *PostgresDb
}

func NewDBA(url string) (*DBStore, error) {
	d, err := NewPostgresDb(url, "configs")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &DBStore{DB: d}, nil
}

func (o *Object) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Object) SetCreatedAt(createdAt time.Time) {
	o.createdAt = createdAt
}

func (o *Object) ModifiedAt() time.Time {
	return o.modifiedAt
}

func (o *Object) SetModifiedAt(modifiedAt time.Time) {
	o.modifiedAt = modifiedAt
}

func (o *Object) Value() map[string]string {
	return o.value
}

func (o *Object) SetValue(value map[string]string) {
	o.value = value
}

func (o *Object) Key() string {
	return o.key
}

func (o *Object) SetKey(key string) {
	o.key = key
}

type Config interface {
	Save(*Object)
	GetConfig(key string) (*Object, error)
	GetConfigByProperty(key string, matchingValue string) (*[]Object, error)
}

func (db *DBStore) Save(key string, value map[string]string) error {
	object := Object{
		key:        key,
		value:      value,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
	}

	if err := db.DB.Save(&object); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (db *DBStore) GetConfig(key string) (*Object, error) {
	return db.DB.GetByKey(key)
}

func (db *DBStore) GetConfigByProperty(key string, matchingValue string) (*[]Object, error) {
	return db.DB.GetByProperty(key, matchingValue)
}

func (db *DBStore) CreateTableFromFile(fileName string) error {
	path := filepath.Join("github.com", "girishpandit88", "configs", fileName)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		// handle error
	}

	requests := strings.Split(string(file), ";")
	for _, stmt := range requests {
		err := db.DB.ExecStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
