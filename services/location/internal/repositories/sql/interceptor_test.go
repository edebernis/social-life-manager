package sqlrepository

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ngrok/sqlmw"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func newInterceptedDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *sqlInterceptor) {
	registry := prometheus.NewRegistry()
	interceptor := newSQLInterceptor(registry)

	db, mock, err := sqlmock.NewWithDSN(
		"mockDSN",
		sqlmock.MonitorPingsOption(true),
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.FailNow()
	}

	sql.Register("mock-mw", sqlmw.Driver(db.Driver(), interceptor))
	inDB, err := sql.Open("mock-mw", "mockDSN")
	if err != nil {
		t.FailNow()
	}

	return inDB, mock, interceptor
}

func TestSQLInterceptor(t *testing.T) {
	ctx := newTestContext()
	db, mock, i := newInterceptedDB(t)
	defer db.Close()

	t.Run("TestPing", func(t *testing.T) {
		mock.ExpectPing()

		err := db.PingContext(ctx)
		assert.NoError(t, err)
	})

	t.Run("TestExecContest", func(t *testing.T) {
		query := "DELETE FROM test WHERE id = $1"
		mock.ExpectPrepare(query).ExpectExec().WithArgs("123").WillReturnResult(sqlmock.NewResult(0, 1))

		stmt, err := db.PrepareContext(ctx, query)
		assert.NoError(t, err)

		_, err = stmt.ExecContext(ctx, "123")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())

		err = testutil.CollectAndCompare(i.metrics.requestCount, strings.NewReader(fmt.Sprintf(`
		# HELP repository_sql_requests_total How many SQL queries processed, partitioned by SQL verbs.
		# TYPE repository_sql_requests_total counter
		repository_sql_requests_total{query="%s",verb="DELETE"} 1
		`, query)), "repository_sql_requests_total")
		assert.NoError(t, err)

		i.metrics.requestCount.Reset()
		i.metrics.requestDuration.Reset()
	})

	t.Run("TestQueryContext", func(t *testing.T) {
		query := "SELECT * FROM test WHERE id = $1"
		mock.ExpectPrepare(query).ExpectQuery().WithArgs("123").WillReturnRows(sqlmock.NewRows(nil))

		stmt, err := db.PrepareContext(ctx, query)
		assert.NoError(t, err)

		_, err = stmt.QueryContext(ctx, "123")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())

		err = testutil.CollectAndCompare(i.metrics.requestCount, strings.NewReader(fmt.Sprintf(`
		# HELP repository_sql_requests_total How many SQL queries processed, partitioned by SQL verbs.
		# TYPE repository_sql_requests_total counter
		repository_sql_requests_total{query="%s",verb="SELECT"} 1
		`, query)), "repository_sql_requests_total")
		assert.NoError(t, err)
	})
}
