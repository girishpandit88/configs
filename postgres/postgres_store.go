package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type PGError string

// implement error interface
func (o PGError) Error() string {
	return fmt.Sprintf("Error from postgres %s", o.Error())
}

type PostgresDb struct {
	db    *pgxpool.Pool
	table string
}

func NewPostgresDb(url string, table string) (*PostgresDb, error) {
	config, err := pgxpool.ParseConfig(url)
	dbConn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresDb{
		db:    dbConn,
		table: table,
	}, nil
}

func (p *PostgresDb) Save(o *Object) error {
	if _, err := p.db.Exec(context.Background(), generateUpsertQuery(p.table), o.Key(), o.Value(), o.CreatedAt(), o.ModifiedAt()); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDb) GetByKey(key string) (*Object, error) {
	row := p.db.QueryRow(context.Background(), generateLookupQuery(p.table), key+"%")
	var keyFromDB string
	var value map[string]string
	var cAt, mAt time.Time
	err := row.Scan(&keyFromDB, &value, &cAt, &mAt)
	if err != nil {
		log.Fatalf("error from db %s", err)
	}
	return &Object{createdAt: cAt, modifiedAt: mAt, key: keyFromDB, value: value}, nil
}

func (p *PostgresDb) GetByProperty(propKey string, propValueMatch string) (*[]Object, error) {
	query, err := p.db.Query(context.Background(), generatePropLookupQuery(p.table), propKey, "%"+propValueMatch+"%")
	defer query.Close()
	var objects []Object
	for query.Next() {
		var keyFromDB string
		var value map[string]string
		var cAt, mAt time.Time
		err := query.Scan(&keyFromDB, &value, &cAt, &mAt)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		o := Object{createdAt: cAt, modifiedAt: mAt, key: keyFromDB, value: value}
		objects = append(objects, o)
	}

	if err != nil {
		log.Fatalf("error from db %s", err)
	}
	return &objects, nil
}

func (p *PostgresDb) ExecStmt(stmt string) error {
	if _, err := p.db.Exec(context.Background(), stmt); err != nil {
		return err
	}
	return nil
}

func generatePropLookupQuery(table string) string {
	return fmt.Sprintf(selectByPropKeyQuery, table)
}

func generateLookupQuery(table string) string {
	return fmt.Sprintf(selectQuery, table)
}

func generateUpsertQuery(table string) string {
	return fmt.Sprintf(upsert, table, table)
}

const (
	upsert = `INSERT INTO %s (key, values, createdAt, modifiedAt) VALUES ($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT %s_pkey
DO UPDATE SET values = $2, modifiedAt = $4`
	selectQuery          = `SELECT key, values, createdAt, modifiedAt from %s where key like $1`
	selectByPropKeyQuery = `SELECT key, values, createdAt, modifiedAt from %s where values->>$1 like $2`
)
