package sqlwrapper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

var (
	ip string
	Log *zap.Logger
)

type Stmt struct {
	stmt    *sql.Stmt
	prepare string
	debug   bool
	slow    time.Duration
}

type Tx struct {
	tx    *sql.Tx
	debug bool
	slow  time.Duration
}

type DB struct {
	db    *sql.DB
	slow  time.Duration
	debug bool
}

func init() {
	Log, _ = zap.NewProduction()
	ip, _ = getExternalIP()
}

// connect returns SQL database connection from the pool
func (db *DB) Conn(ctx context.Context) (*sql.Conn, error) {
	c, err := db.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

func getExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}






func (t *Tx) Commit() error {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if t.debug || total >= t.slow {
			Log.Info("tx commit",
				zap.String("use-time", total.String()),
				zap.String("ip", ip), )
		}
	}()
	return t.tx.Commit()
}

func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if t.debug || total >= t.slow {
			logInfo("tx exec", total.String(), query, args)
		}
	}()
	return t.tx.Exec(query, args...)
}
func (t *Tx) Prepare(query string) (*Stmt, error) {
	s, err := t.tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	stmt := &Stmt{
		stmt:    s,
		debug:   t.debug,
		prepare: query,
		slow:    t.slow,
	}
	return stmt, nil
}
func (t *Tx) Rollback() error {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if t.debug || total >= t.slow {
			Log.Info("tx rollback",
				zap.String("use-time", total.String()),
				zap.String("ip", ip),
			)
		}
	}()
	return t.tx.Rollback()
}
func (t *Tx) Stmt(stmt *Stmt) *Stmt {
	s := t.tx.Stmt(stmt.stmt)
	stmt.stmt = s
	return stmt
}
func (t *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if t.debug || total >= t.slow {
			logInfo("tx query", total.String(), query, args)
		}
	}()
	return t.tx.Query(query, args...)
}
func (t *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if t.debug || total >= t.slow {
			logInfo("tx query row", total.String(), query, args)
		}
	}()
	return t.tx.QueryRow(query, args...)
}


func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if s.debug || total >= s.slow {
			logInfo("stmt query row", total.String(), s.prepare, args)
		}
	}()
	return s.stmt.Exec(args...)
}
func (s *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if s.debug || total >= s.slow {
			logInfo("stmt query", total.String(), s.prepare, args)
		}
	}()
	return s.stmt.Query(args...)
}
func (s *Stmt) QueryRow(args ...interface{}) *sql.Row {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if s.debug || total >= s.slow {
			logInfo("stmt query row", total.String(), s.prepare, args)
		}
	}()
	return s.stmt.QueryRow(args...)
}
func (s *Stmt) Close() error {
	return s.stmt.Close()
}


func WrapperDB(db *sql.DB, debug bool, slow time.Duration) (d *DB) {
	return &DB{
		db:    db,
		slow:  slow,
		debug: debug,
	}
}
func (d *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if d.debug || total >= d.slow {
			logInfo("db ExecContext", total.String(), query, args)
		}
	}()
	return d.db.ExecContext(ctx, query, args...)

}
func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if d.debug || total >= d.slow {
			logInfo("db exec", total.String(), query, args)
		}
	}()
	return d.db.Exec(query, args...)

}

func (d *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if d.debug || total >= d.slow {
			logInfo("db QueryContext", total.String(), query, args)
		}
	}()
	return d.db.QueryContext(ctx, query, args...)
}
func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if d.debug || total >= d.slow {
			logInfo("db query", total.String(), query, args)
		}
	}()
	return d.db.Query(query, args...)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if d.debug || total >= d.slow {
			logInfo("db QueryRowContext", total.String(), query, args)
		}
	}()
	return d.db.QueryRowContext(ctx, query, args...)
}
func (d *DB) Ping() error {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if d.debug || total >= d.slow {
			logInfo("db ping", total.String(), "", "")
		}
	}()
	return d.db.Ping()
}
func (d *DB) PingContext(ctx context.Context) error {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if d.debug || total >= d.slow {
			logInfo("db PingContext", total.String(), "", "")
		}
	}()
	return d.db.PingContext(ctx)
}
func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	st := time.Now()
	defer func() {
		et := time.Now()
		total := et.Sub(st)
		if d.debug || total >= d.slow {
			logInfo("db query row", total.String(), query, args)
		}
	}()
	return d.db.QueryRow(query, args...)
}
func (d *DB) Close() error {
	return d.db.Close()
}
func (d *DB) BeginTX(ctx context.Context, opts *sql.TxOptions) (t *Tx, err error) {
	tx, err := d.db.BeginTx(ctx, opts)
	if err != nil {
		return
	}
	t = &Tx{
		tx:    tx,
		debug: d.debug,
		slow:  d.slow,
	}
	return
}

func (d *DB) Begin() (t *Tx, err error) {
	tx, err := d.db.Begin()
	if err != nil {
		return
	}
	t = &Tx{
		tx:    tx,
		debug: d.debug,
		slow:  d.slow,
	}
	return
}
func (d *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	s, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{
		stmt:    s,
		prepare: query,
		debug:   d.debug,
		slow:    d.slow,
	}, nil
}
func (d *DB) Prepare(query string) (*Stmt, error) {
	s, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{
		stmt:    s,
		prepare: query,
		debug:   d.debug,
		slow:    d.slow,
	}, nil
}

func logInfo(tag string, time string, query string, args ...interface{}) {
	argStr := ""
	for _, v:= range args {
		argStr += fmt.Sprintf("%v", v)
	}
	Log.Info(tag,
		zap.String("use-time", time),
		zap.String("sql", query),
		zap.String("args", argStr),
		zap.String("ip", ip),
	)
}