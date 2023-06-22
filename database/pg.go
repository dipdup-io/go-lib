package database

import (
	"context"
	"fmt"
	"github.com/dipdup-net/go-lib/config"
	pg "github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

type PgGoConnection interface {
	DB() PgDB
}

type PgDB interface {
	/*
		Begin() (*pg.Tx, error)
		BeginContext(ctx context.Context) (*pg.Tx, error)
		RunInTransaction(ctx context.Context, fn func(*pg.Tx) error) error
		AddQueryHook(hook pg.QueryHook)
		beforeQuery(ctx context.Context, ormDB orm.DB, model interface{}, query interface{}, params []interface{}, fmtedQuery []byte) (context.Context, *QueryEvent, error)
		afterQuery(ctx context.Context, event *pg.QueryEvent, res pg.Result, err error) error
		afterQueryFromIndex(ctx context.Context, event *pg.QueryEvent, hookIndex int) error
		//startup(c context.Context, cn *pool.Conn, user string, password string, database string, appName string) error
		//enableSSL(c context.Context, cn *pool.Conn, tlsConf *tls.Config) error
		//auth(c context.Context, cn *pool.Conn, rd *pool.ReaderContext, user string, password string) error
		//logStartupNotice(rd *pool.ReaderContext) error
		//authCleartext(c context.Context, cn *pool.Conn, rd *pool.ReaderContext, password string) error
		//authMD5(c context.Context, cn *pool.Conn, rd *pool.ReaderContext, user string, password string) error
		//authSASL(c context.Context, cn *pool.Conn, rd *pool.ReaderContext, user string, password string) error
		PoolStats() *pg.PoolStats
		//clone() *baseDB
		//withPool(p pool.Pooler) *baseDB
		//WithTimeout(d time.Duration) *baseDB
		//WithParam(param string, value interface{}) *baseDB
		Param(param string) interface{}
		retryBackoff(retry int) time.Duration
		//getConn(ctx context.Context) (*pool.Conn, error)
		//initConn(ctx context.Context, cn *pool.Conn) error
		//releaseConn(ctx context.Context, cn *pool.Conn, err error)
		//withConn(ctx context.Context, fn func(context.Context, *pool.Conn) error) error
		shouldRetry(err error) bool
		Close() error
		Exec(query interface{}, params ...interface{}) (res pg.Result, err error)
	*/
	ExecContext(c context.Context, query interface{}, params ...interface{}) (pg.Result, error)
	/*
		//exec(ctx context.Context, query interface{}, params ...interface{}) (pg.Result, error)
		ExecOne(query interface{}, params ...interface{}) (pg.Result, error)
		ExecOneContext(ctx context.Context, query interface{}, params ...interface{}) (pg.Result, error)
		execOne(c context.Context, query interface{}, params ...interface{}) (pg.Result, error)
		Query(model interface{}, query interface{}, params ...interface{}) (res pg.Result, err error)
		QueryContext(c context.Context, model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
		query(ctx context.Context, model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
		QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
		QueryOneContext(ctx context.Context, model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
		queryOne(ctx context.Context, model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
		CopyFrom(r io.Reader, query interface{}, params ...interface{}) (res pg.Result, err error)
		//copyFrom(ctx context.Context, cn *pool.Conn, r io.Reader, query interface{}, params ...interface{}) (res Result, err error)
		CopyTo(w io.Writer, query interface{}, params ...interface{}) (res pg.Result, err error)
		//copyTo(ctx context.Context, cn *pool.Conn, w io.Writer, query interface{}, params ...interface{}) (res Result, err error)
		Ping(ctx context.Context) error
		Model(model ...interface{}) *pg.Query
		ModelContext(c context.Context, model ...interface{}) *pg.Query
		Formatter() orm.QueryFormatter
		//cancelRequest(processID int32, secretKey int32) error
		//simpleQuery(c context.Context, cn *pool.Conn, wb *pool.WriteBuffer) (*result, error)
		//simpleQueryData(c context.Context, cn *pool.Conn, model interface{}, wb *pool.WriteBuffer) (*result, error)
		Prepare(q string) (*pg.Stmt, error)
		//prepare(c context.Context, cn *pool.Conn, q string) (string, []types.ColumnInfo, error)
		//closeStmt(c context.Context, cn *pool.Conn, name string) error
	*/
}

// PgGo -
type PgGo struct {
	conn *pg.DB
}

// NewPgGo -
func NewPgGo() *PgGo {
	return new(PgGo)
}

// DB -
func (db *PgGo) DB() PgDB {
	return db.conn
}

// Connect -
func (db *PgGo) Connect(ctx context.Context, cfg config.Database) error {
	if cfg.Kind != config.DBKindPostgres {
		return errors.Wrap(ErrUnsupportedDatabaseType, cfg.Kind)
	}
	var conn *pg.DB
	if cfg.Path != "" {
		opt, err := pg.ParseURL(cfg.Path)
		if err != nil {
			return err
		}
		conn = pg.Connect(opt)
	} else {
		conn = pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			User:     cfg.User,
			Password: cfg.Password,
			Database: cfg.Database,
		})
	}
	db.conn = conn
	return nil
}

// Close -
func (db *PgGo) Close() error {
	return db.conn.Close()
}

// Ping -
func (db *PgGo) Ping(ctx context.Context) error {
	if db.conn == nil {
		return ErrConnectionIsNotInitialized
	}
	return db.conn.Ping(ctx)
}

// State -
func (db *PgGo) State(indexName string) (*State, error) {
	var s State
	err := db.conn.Model(&s).Where("index_name = ?", indexName).Limit(1).Select()
	return &s, err
}

// CreateState -
func (db *PgGo) CreateState(s *State) error {
	_, err := db.conn.Model(s).Insert()
	return err
}

// UpdateState -
func (db *PgGo) UpdateState(s *State) error {
	_, err := db.conn.Model(s).Where("index_name = ?", s.IndexName).Update()
	return err
}

// DeleteState -
func (db *PgGo) DeleteState(s *State) error {
	_, err := db.conn.Model(s).Where("index_name = ?", s.IndexName).Delete()
	return err
}
