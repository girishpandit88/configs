package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type PGError string

// implement error interface
func (o PGError) Error() string {
	return fmt.Sprintf("Error from postgres %s", o.Error())
}

type PostgresDb struct {
	db            *pgxpool.Pool
	table         string
	preparedStmts *map[string]string
}

func NewPostgresDb(url string, table string) (*PostgresDb, error) {
	config, err := pgxpool.ParseConfig(url)
	dbConn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &PostgresDb{
		db:            dbConn,
		table:         table,
		preparedStmts: generatedPreparedStmtsMap(),
	}, nil
}

func generatedPreparedStmtsMap() *map[string]string {
	return &map[string]string{
		"UPSERT":               upsert,
		"SELECTBYKEY":          selectQuery,
		"SELECTBYPROPKEYVALUE": selectByPropKeyQuery,
	}
}

func (p *PostgresDb) Save(o *Object) error {
	if _, err := p.db.Exec(
		context.Background(),
		p.generateQuery("UPSERT", p.table),
		o.Key(),
		o.Value(),
		o.CreatedAt(),
		o.ModifiedAt()); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDb) GetByKey(key string) (*Object, error) {
	row := p.db.QueryRow(
		context.Background(),
		p.generateQuery("SELECTBYKEY", p.table),
		key+"%")
	var keyFromDB string
	var value map[string]interface{}
	var cAt, mAt time.Time
	err := row.Scan(&keyFromDB, &value, &cAt, &mAt)
	if err != nil {
		return nil, err
	}
	return &Object{createdAt: cAt, modifiedAt: mAt, key: keyFromDB, value: value}, nil
}

func (p *PostgresDb) GetByProperty(propKey string, propValueMatch string) (*[]Object, error) {
	rows, err := p.db.Query(
		context.Background(),
		p.generateQuery("SELECTBYPROPKEYVALUE", p.table),
		propKey, "%"+propValueMatch+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var objects []Object
	for rows.Next() {
		var keyFromDB string
		var value map[string]interface{}
		var cAt, mAt time.Time
		err := rows.Scan(&keyFromDB, &value, &cAt, &mAt)
		if err != nil {
			return nil, err
		}
		o := Object{createdAt: cAt, modifiedAt: mAt, key: keyFromDB, value: value}
		objects = append(objects, o)
	}
	return &objects, nil
}

func (p *PostgresDb) ExecStmt(stmt string) error {
	if _, err := p.db.Exec(context.Background(), stmt); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDb) generateQuery(query string, table string) string {
	m := *p.preparedStmts
	s := m[query]
	return fmt.Sprintf(s, table)
}

const (
	upsert               = `INSERT INTO %[1]s (key, values, createdAt, modifiedAt) VALUES ($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT %[1]s_pkey DO UPDATE SET values = $2, modifiedAt = $4`
	selectQuery          = `SELECT key, values, createdAt, modifiedAt from %s where key like $1`
	selectByPropKeyQuery = `SELECT key, values, createdAt, modifiedAt from %s where values->>$1 like $2`
)
