package main

import (
	"database/sql"
	"github.com/jackc/go_db_bench/raw"
	"github.com/jackc/pgx"
	"sync"
	"testing"
	"time"
)

var (
	setupOnce     sync.Once
	pgxPool       *pgx.ConnPool
	pgxStdlib     *sql.DB
	pq            *sql.DB
	rawConn       *raw.Conn
	randPersonIDs []int32
)

var selectPersonNameSQL = `select first_name from person where id=$1`
var selectPersonSQL = `
select id, first_name, last_name, sex, birth_date, weight, height
from person
where id=$1`
var selectMultiplePeopleSQL = `
select id, first_name, last_name, sex, birth_date, weight, height
from person
where id between $1 and $1 + 25`

var rxBuf []byte

type person struct {
	id        int32
	firstName string
	lastName  string
	sex       string
	birthDate time.Time
	weight    int32
	height    int32
}

func setup(b *testing.B) {
	setupOnce.Do(func() {
		config := extractConfig()

		config.AfterConnect = func(conn *pgx.Conn) error {
			_, err := conn.Prepare("selectPersonName", selectPersonNameSQL)
			if err != nil {
				return err
			}

			_, err = conn.Prepare("selectPerson", selectPersonSQL)
			if err != nil {
				return err
			}

			_, err = conn.Prepare("selectMultiplePeople", selectMultiplePeopleSQL)
			if err != nil {
				return err
			}

			return nil
		}

		err := loadTestData(config)
		if err != nil {
			b.Fatalf("loadTestData failed: %v", err)
		}

		pgxPool, err = openPgxNative(config)
		if err != nil {
			b.Fatalf("openPgxNative failed: %v", err)
		}

		pgxStdlib, err = openPgxStdlib(config)
		if err != nil {
			b.Fatalf("openPgxNative failed: %v", err)
		}

		pq, err = openPq(config)
		if err != nil {
			b.Fatalf("openPq failed: %v", err)
		}

		rawConfig := raw.ConnConfig{
			Host:     config.Host,
			Port:     config.Port,
			User:     config.User,
			Password: config.Password,
			Database: config.Database,
		}
		rawConn, err = raw.Connect(rawConfig)
		if err != nil {
			b.Fatalf("raw.Connect failed: %v", err)
		}
		_, err = rawConn.Prepare("selectPersonName", selectPersonNameSQL)
		if err != nil {
			b.Fatalf("rawConn.Prepare failed: %v", err)
		}
		_, err = rawConn.Prepare("selectPerson", selectPersonSQL)
		if err != nil {
			b.Fatalf("rawConn.Prepare failed: %v", err)
		}
		_, err = rawConn.Prepare("selectMultiplePeople", selectMultiplePeopleSQL)
		if err != nil {
			b.Fatalf("rawConn.Prepare failed: %v", err)
		}

		rxBuf = make([]byte, 16384)

		// Get random person ids in random order outside of timing
		ids, err := pgxPool.SelectValues("select id from person order by random()")
		if err != nil {
			b.Fatalf("pgxPool.SelectValues failed: %v", err)
		}
		for _, id := range ids {
			randPersonIDs = append(randPersonIDs, id.(int32))
		}
	})
}

func BenchmarkPgxNativeSelectSingleValueUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkPgxNativeSelectSingleValue(b, selectPersonNameSQL)
}

func benchmarkPgxNativeSelectSingleValue(b *testing.B, sql string) {
	for i := 0; i < b.N; i++ {
		id := randPersonIDs[i%len(randPersonIDs)]
		firstName, err := pgxPool.SelectValue(sql, id)
		if err != nil {
			b.Fatalf("pgxPool.SelectValue failed: %v", err)
		}
		if len(firstName.(string)) == 0 {
			b.Fatal("firstName was empty")
		}
	}
}

func BenchmarkPgxStdlibSelectSingleValueUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkSelectSingleValueUnprepared(b, pgxStdlib)
}

func BenchmarkPqSelectSingleValueUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkSelectSingleValueUnprepared(b, pq)
}

func benchmarkSelectSingleValueUnprepared(b *testing.B, db *sql.DB) {
	for i := 0; i < b.N; i++ {
		id := randPersonIDs[i%len(randPersonIDs)]
		row := db.QueryRow(selectPersonNameSQL, id)
		var firstName string
		err := row.Scan(&firstName)
		if err != nil {
			b.Fatalf("row.Scan failed: %v", err)
		}
		if len(firstName) == 0 {
			b.Fatal("firstName was empty")
		}
	}
}

func BenchmarkPgxNativeSelectSingleValuePrepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkPgxNativeSelectSingleValue(b, "selectPersonName")
}

func BenchmarkPgxStdlibSelectSingleValuePrepared(b *testing.B) {
	setup(b)
	stmt, err := pgxStdlib.Prepare(selectPersonNameSQL)
	if err != nil {
		b.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()

	b.ResetTimer()
	benchmarkSelectSingleValuePrepared(b, stmt)
}

func BenchmarkPqSelectSingleValuePrepared(b *testing.B) {
	setup(b)
	stmt, err := pq.Prepare(selectPersonNameSQL)
	if err != nil {
		b.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()

	b.ResetTimer()
	benchmarkSelectSingleValuePrepared(b, stmt)
}

func benchmarkSelectSingleValuePrepared(b *testing.B, stmt *sql.Stmt) {
	for i := 0; i < b.N; i++ {
		id := randPersonIDs[i%len(randPersonIDs)]
		row := stmt.QueryRow(id)
		var firstName string
		err := row.Scan(&firstName)
		if err != nil {
			b.Fatalf("row.Scan failed: %v", err)
		}
		if len(firstName) == 0 {
			b.Fatal("firstName was empty")
		}
	}
}

func BenchmarkRawSelectSingleValuePrepared(b *testing.B) {
	setup(b)
	b.ResetTimer()

	txBuf := []byte{0x42, 0x0, 0x0, 0x0, 0x28, 0x0, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x11, 0x61, 0x0, 0x1, 0x0, 0x0, 0x45, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x53, 0x0, 0x0, 0x0, 0x4}

	for i := 0; i < b.N; i++ {
		_, err := rawConn.Conn.Write(txBuf)
		if err != nil {
			b.Fatalf("rawConn.Conn.Write failed: %v", err)
		}

		rxRawUntilReady(b)
	}
}

func BenchmarkPgxNativeSelectSingleRowUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkPgxNativeSelectSingleRow(b, selectPersonSQL)
}

func benchmarkPgxNativeSelectSingleRow(b *testing.B, sql string) {
	for i := 0; i < b.N; i++ {
		var p person
		id := randPersonIDs[i%len(randPersonIDs)]
		err := pgxPool.SelectFunc(sql, func(r *pgx.DataRowReader) error {
			p.id = r.ReadValue().(int32)
			p.firstName = r.ReadValue().(string)
			p.lastName = r.ReadValue().(string)
			p.sex = r.ReadValue().(string)
			p.birthDate = r.ReadValue().(time.Time)
			p.weight = r.ReadValue().(int32)
			p.height = r.ReadValue().(int32)
			return nil
		}, id)
		if err != nil {
			b.Fatalf("pgxPool.SelectFunc failed: %v", err)
		}

		checkPersonWasFilled(b, p)
	}
}

func checkPersonWasFilled(b *testing.B, p person) {
	if p.id == 0 {
		b.Fatal("id was 0")
	}
	if len(p.firstName) == 0 {
		b.Fatal("firstName was empty")
	}
	if len(p.lastName) == 0 {
		b.Fatal("lastName was empty")
	}
	if len(p.sex) == 0 {
		b.Fatal("sex was empty")
	}
	var zeroTime time.Time
	if p.birthDate == zeroTime {
		b.Fatal("birthDate was zero time")
	}
	if p.weight == 0 {
		b.Fatal("weight was 0")
	}
	if p.height == 0 {
		b.Fatal("height was 0")
	}
}

func BenchmarkPgxStdlibSelectSingleRowUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkSelectSingleRowUnprepared(b, pgxStdlib)
}

func BenchmarkPqSelectSingleRowUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkSelectSingleRowUnprepared(b, pq)
}

func benchmarkSelectSingleRowUnprepared(b *testing.B, db *sql.DB) {
	for i := 0; i < b.N; i++ {
		id := randPersonIDs[i%len(randPersonIDs)]
		row := db.QueryRow(selectPersonSQL, id)
		var p person
		err := row.Scan(&p.id, &p.firstName, &p.lastName, &p.sex, &p.birthDate, &p.weight, &p.height)
		if err != nil {
			b.Fatalf("row.Scan failed: %v", err)
		}

		checkPersonWasFilled(b, p)
	}
}

func BenchmarkPgxNativeSelectSingleRowPrepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkPgxNativeSelectSingleRow(b, "selectPerson")
}

func BenchmarkPgxStdlibSelectSingleRowPrepared(b *testing.B) {
	setup(b)
	stmt, err := pgxStdlib.Prepare(selectPersonSQL)
	if err != nil {
		b.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()

	b.ResetTimer()
	benchmarkSelectSingleRowPrepared(b, stmt)
}

func BenchmarkPqSelectSingleRowPrepared(b *testing.B) {
	setup(b)
	stmt, err := pq.Prepare(selectPersonSQL)
	if err != nil {
		b.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()

	b.ResetTimer()
	benchmarkSelectSingleRowPrepared(b, stmt)
}

func benchmarkSelectSingleRowPrepared(b *testing.B, stmt *sql.Stmt) {
	for i := 0; i < b.N; i++ {
		id := randPersonIDs[i%len(randPersonIDs)]
		row := stmt.QueryRow(id)
		var p person
		err := row.Scan(&p.id, &p.firstName, &p.lastName, &p.sex, &p.birthDate, &p.weight, &p.height)
		if err != nil {
			b.Fatalf("row.Scan failed: %v", err)
		}

		checkPersonWasFilled(b, p)
	}
}

func BenchmarkRawSelectSingleRowPrepared(b *testing.B) {
	setup(b)
	b.ResetTimer()

	txBuf := []byte{0x42, 0x0, 0x0, 0x0, 0x30, 0x0, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0xc, 0x7e, 0x0, 0x7, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x45, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x53, 0x0, 0x0, 0x0, 0x4}

	for i := 0; i < b.N; i++ {
		_, err := rawConn.Conn.Write(txBuf)
		if err != nil {
			b.Fatalf("rawConn.Conn.Write failed: %v", err)
		}

		rxRawUntilReady(b)
	}
}

func BenchmarkPgxNativeSelectMultipleRowsUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkPgxNativeSelectMultipleRows(b, selectMultiplePeopleSQL)
}

func benchmarkPgxNativeSelectMultipleRows(b *testing.B, sql string) {
	for i := 0; i < b.N; i++ {
		var people []person
		id := randPersonIDs[i%len(randPersonIDs)]
		err := pgxPool.SelectFunc(sql, func(r *pgx.DataRowReader) error {
			var p person
			p.id = r.ReadValue().(int32)
			p.firstName = r.ReadValue().(string)
			p.lastName = r.ReadValue().(string)
			p.sex = r.ReadValue().(string)
			p.birthDate = r.ReadValue().(time.Time)
			p.weight = r.ReadValue().(int32)
			p.height = r.ReadValue().(int32)
			people = append(people, p)
			return nil
		}, id)
		if err != nil {
			b.Fatalf("pgxPool.SelectFunc failed: %v", err)
		}

		for _, p := range people {
			checkPersonWasFilled(b, p)
		}
	}
}

func BenchmarkPgxStdlibSelectMultipleRowsUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkSelectMultipleRowsUnprepared(b, pgxStdlib)
}

func BenchmarkPqSelectMultipleRowsUnprepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkSelectMultipleRowsUnprepared(b, pq)
}

func benchmarkSelectMultipleRowsUnprepared(b *testing.B, db *sql.DB) {
	for i := 0; i < b.N; i++ {
		var people []person
		id := randPersonIDs[i%len(randPersonIDs)]
		rows, err := db.Query(selectMultiplePeopleSQL, id)
		if err != nil {
			b.Fatalf("db.Query failed: %v", err)
		}

		for rows.Next() {
			var p person
			err := rows.Scan(&p.id, &p.firstName, &p.lastName, &p.sex, &p.birthDate, &p.weight, &p.height)
			if err != nil {
				b.Fatalf("rows.Scan failed: %v", err)
			}
			people = append(people, p)
		}

		if rows.Err() != nil {
			b.Fatalf("rows.Err() returned an error: %v", err)
		}

		for _, p := range people {
			checkPersonWasFilled(b, p)
		}
	}
}

func BenchmarkPgxNativeSelectMultipleRowsPrepared(b *testing.B) {
	setup(b)
	b.ResetTimer()
	benchmarkPgxNativeSelectMultipleRows(b, "selectMultiplePeople")
}

func BenchmarkPgxStdlibSelectMultipleRowsPrepared(b *testing.B) {
	setup(b)
	stmt, err := pgxStdlib.Prepare(selectMultiplePeopleSQL)
	if err != nil {
		b.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()
	b.ResetTimer()
	benchmarkSelectMultipleRowsPrepared(b, stmt)
}

func BenchmarkPqSelectMultipleRowsPrepared(b *testing.B) {
	setup(b)
	stmt, err := pq.Prepare(selectMultiplePeopleSQL)
	if err != nil {
		b.Fatalf("Prepare failed: %v", err)
	}
	defer stmt.Close()

	b.ResetTimer()
	benchmarkSelectMultipleRowsPrepared(b, stmt)
}

func benchmarkSelectMultipleRowsPrepared(b *testing.B, stmt *sql.Stmt) {
	for i := 0; i < b.N; i++ {
		var people []person
		id := randPersonIDs[i%len(randPersonIDs)]
		rows, err := stmt.Query(id)
		if err != nil {
			b.Fatalf("db.Query failed: %v", err)
		}

		for rows.Next() {
			var p person
			err := rows.Scan(&p.id, &p.firstName, &p.lastName, &p.sex, &p.birthDate, &p.weight, &p.height)
			if err != nil {
				b.Fatalf("rows.Scan failed: %v", err)
			}
			people = append(people, p)
		}

		if rows.Err() != nil {
			b.Fatalf("rows.Err() returned an error: %v", err)
		}

		for _, p := range people {
			checkPersonWasFilled(b, p)
		}
	}
}

func BenchmarkRawSelectMultipleRowsPrepared(b *testing.B) {
	setup(b)
	b.ResetTimer()

	txBuf := []byte{0x42, 0x0, 0x0, 0x0, 0x38, 0x0, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x65, 0x50, 0x65, 0x6f, 0x70, 0x6c, 0x65, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x11, 0x83, 0x0, 0x7, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x45, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x53, 0x0, 0x0, 0x0, 0x4}

	for i := 0; i < b.N; i++ {
		_, err := rawConn.Conn.Write(txBuf)
		if err != nil {
			b.Fatalf("rawConn.Conn.Write failed: %v", err)
		}

		rxRawUntilReady(b)
	}
}

func rxRawUntilReady(b *testing.B) {
	for {
		n, err := rawConn.Conn.Read(rxBuf)
		if err != nil {
			b.Fatalf("rawConn.Conn.Read failed: %v", err)
		}
		if rxBuf[n-6] == 'Z' && rxBuf[n-2] == 5 && rxBuf[n-1] == 'I' {
			return
		}
	}
}