package middlewares

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

const (
	secretKey        = "secret_key"
	defaultSecretKey = "replace_me_with_strong_string"
)

const (
	reqsName    = "http_requests_total"
	latencyName = "http_request_duration_seconds"
)

var (
	curSecretValue = defaultSecretKey

	// DefaultBuckets prometheus buckets in seconds.
	DefaultBuckets = []float64{0.3, 1.2, 5.0}
)

type tickProm struct {
	prom *Prometheus
}

type TickPromConfig struct {
	ServiceName string
	SecretKey   string
}

func createDefaultTickPromConfig() TickPromConfig {
	return TickPromConfig{
		ServiceName: "unknown",
		SecretKey:   defaultSecretKey,
	}
}

func newTickProm(config TickPromConfig) *tickProm {
	//统计请求在 1ms内 10ms 100ms 1s 10s 30s区间内的数量
	context.SetHandlerName("github.com/iris-contrib/middleware/prometheus.*", "iris-contrib.prometheus")
	tick := New(config.ServiceName, 0.001, 0.01, 0.1, 1, 10, 30.0)

	if len(config.SecretKey) > 0 {
		curSecretValue = config.SecretKey
	}

	return &tickProm{prom: tick}
}

func (t *tickProm) ServeHTTP(ctx iris.Context) {
	t.prom.ServeHTTP(ctx)
}

//serName : parse current project service name here
func TickProm(config ...TickPromConfig) context.Handler {
	c := createDefaultTickPromConfig()
	if len(config) > 0 {
		c = config[0]
	}

	h := newTickProm(c).ServeHTTP
	return h
}

func TickMetricHandler() context.Handler {
	handler := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
	)

	return func(ctx iris.Context) {
		key := ctx.URLParamDefault(secretKey, "")

		if len(key) > 0 && key == curSecretValue {
			handler.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		} else {
			ctx.StatusCode(403)
			_, _ = ctx.WriteString(http.StatusText(403))
			ctx.StopExecution()
		}
	}
}

// Prometheus is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
//
// Usage: pass its `ServeHTTP` to a route or globally.
type Prometheus struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// New returns a new prometheus middleware.
//
// If buckets are empty then `DefaultBuckets` are set.
func New(name string, buckets ...float64) *Prometheus {
	p := Prometheus{}
	p.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.reqs)

	if len(buckets) == 0 {
		buckets = DefaultBuckets
	}

	p.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.latency)

	return &p
}

func (p *Prometheus) ServeHTTP(ctx iris.Context) {
	start := time.Now()
	ctx.Next()
	r := ctx.Request()
	statusCode := strconv.Itoa(ctx.GetStatusCode())

	p.reqs.WithLabelValues(statusCode, r.Method, r.URL.Path).
		Inc()

	p.latency.WithLabelValues(statusCode, r.Method, r.URL.Path).
		Observe(float64(time.Since(start).Nanoseconds()) / 1000000000)
}
