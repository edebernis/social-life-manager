package sqlrepository

import (
	"context"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/ngrok/sqlmw"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type sqlInterceptor struct {
	sqlmw.NullInterceptor

	metrics *sqlMetrics
}

func newSQLInterceptor(registry prometheus.Registerer) *sqlInterceptor {
	return &sqlInterceptor{
		metrics: newSQLMetrics(registry),
	}
}

type sqlMetrics struct {
	namespace string
	subsystem string

	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func newSQLMetrics(registry prometheus.Registerer) *sqlMetrics {
	metrics := &sqlMetrics{
		namespace: "repository",
		subsystem: "sql",
	}

	metrics.requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metrics.namespace,
			Subsystem: metrics.subsystem,
			Name:      "requests_total",
			Help:      "How many SQL queries processed, partitioned by SQL verbs.",
		},
		[]string{"verb", "query"},
	)
	registry.MustRegister(metrics.requestCount)

	metrics.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metrics.namespace,
			Subsystem: metrics.subsystem,
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			Name:      "request_duration_seconds",
			Help:      "The SQL query latency bucket.",
		}, []string{"verb", "query"},
	)
	registry.MustRegister(metrics.requestDuration)

	return metrics
}

func (i *sqlInterceptor) ConnectorConnect(ctx context.Context, connector driver.Connector) (driver.Conn, error) {
	startedAt := time.Now()

	conn, err := connector.Connect(ctx)

	logger.WithContext(ctx).WithFields(logrus.Fields{
		"duration": time.Since(startedAt).Milliseconds(),
		"err":      err,
	}).Info("connected to sql db")

	return conn, err
}

func (i *sqlInterceptor) ConnPing(ctx context.Context, pinger driver.Pinger) error {
	startedAt := time.Now()

	err := pinger.Ping(ctx)

	logger.WithContext(ctx).WithFields(logrus.Fields{
		"duration": time.Since(startedAt).Milliseconds(),
		"err":      err,
	}).Info("ping sql db")

	return err
}

func (i *sqlInterceptor) ConnPrepareContext(ctx context.Context, conn driver.ConnPrepareContext, query string) (driver.Stmt, error) {
	startedAt := time.Now()

	stmt, err := conn.PrepareContext(ctx, query)

	logger.WithContext(ctx).WithFields(logrus.Fields{
		"duration": time.Since(startedAt).Milliseconds(),
		"query":    query,
		"err":      err,
	}).Info("prepared sql request")

	return stmt, err
}

func (i *sqlInterceptor) StmtExecContext(ctx context.Context, stmt driver.StmtExecContext, query string, args []driver.NamedValue) (driver.Result, error) {
	startedAt := time.Now()

	res, err := stmt.ExecContext(ctx, args)

	logger.WithContext(ctx).WithFields(logrus.Fields{
		"duration": time.Since(startedAt).Milliseconds(),
		"query":    query,
		"args":     args,
		"err":      err,
	}).Info("executed sql request")

	verb := strings.ToUpper(strings.Fields(query)[0])
	i.metrics.requestCount.WithLabelValues(verb, query).Inc()

	elapsed := float64(time.Since(startedAt)) / float64(time.Second)
	i.metrics.requestDuration.WithLabelValues(verb, query).Observe(elapsed)

	return res, err
}

func (i *sqlInterceptor) StmtQueryContext(ctx context.Context, stmt driver.StmtQueryContext, query string, args []driver.NamedValue) (driver.Rows, error) {
	startedAt := time.Now()

	rows, err := stmt.QueryContext(ctx, args)

	logger.WithContext(ctx).WithFields(logrus.Fields{
		"duration": time.Since(startedAt).Milliseconds(),
		"query":    query,
		"args":     args,
		"err":      err,
	}).Info("executed sql query")

	verb := strings.ToUpper(strings.Fields(query)[0])
	i.metrics.requestCount.WithLabelValues(verb, query).Inc()

	elapsed := float64(time.Since(startedAt)) / float64(time.Second)
	i.metrics.requestDuration.WithLabelValues(verb, query).Observe(elapsed)

	return rows, err
}
