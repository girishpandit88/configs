A Simple Config service in Go
============================================


- You need to start docker container for postgres included.
```bash
docker-compose up -d
```
- Create tables and indexes
```go
go run init_table.go
```

- To initialize a db store: 
```go
dbstore := postgres.NewDBA(dburl)
```

- To save a value:
```go
dbstore.Save(dbObject)
```

- To retrieve by key:
```go
dbstore.GetConfig(key)
```

- To retrieve by matching config name and value:
```go
//returns config objects which have a url config and contains value "goo"
dbstore.GetConfigByProperty("url", "goo")
```
