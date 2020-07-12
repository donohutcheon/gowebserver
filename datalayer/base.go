package datalayer

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xo/dburl"
	"os"
)

type Model struct {
	ID        int64        `json:"id" db:"id"`
	CreatedAt JsonNullTime `json:"createdAt" db:"created_at"`
	UpdatedAt JsonNullTime `json:"updatedAt" db:"updated_at"`
	DeletedAt JsonNullTime `json:"deletedAt" db:"deleted_at"`
}

type PersistenceDataLayer struct {
	conn *sqlx.DB
}

var (
	ErrNoData = sql.ErrNoRows
)

func New() (*PersistenceDataLayer, error){
	conn, err, ok := tryConnectHerokuJawsDB()
	if err != nil {
		// TODO: Use proper logger
		fmt.Printf("Could not connect to JawsDB. %s", err.Error())
		return nil, err
	} else if ok {
		return &PersistenceDataLayer{
			conn: conn,
		}, nil
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	conn, err = createCon("mysql", username, password, dbHost, dbPort, dbName)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return &PersistenceDataLayer{
		conn: conn,
	}, nil
}

func (p *PersistenceDataLayer) GetConn() *sqlx.DB {
	return p.conn
}

/*Create mysql connection*/
func createCon(driverName string, username string, password string, dbHost string, dbPort string, dbName string) (db *sqlx.DB, err error) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, dbHost, dbPort, dbName)
	db, err = sqlx.Open(driverName, dbURI)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("database is connected")
	}

	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Printf("MySQL datalayer is not connected %s", err.Error())
	}
	return db, err
}

func tryConnectHerokuJawsDB() (*sqlx.DB, error, bool){
	dbURI := os.Getenv("JAWSDB_MARIA_URL") + "?parseTime=true"
	if len(dbURI) == 0 {
		return nil, nil, false
	}

	db, err := dburl.Open( dbURI)
	if err != nil {
		return nil, err, false
	} else {
		fmt.Println("database is connected")
	}

	dbx := sqlx.NewDb(db, "mysql")

	err = dbx.Ping()
	if err != nil {
		return nil, err, false
	}

	return dbx, nil, true
}