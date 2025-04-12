package monitoring

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests",
	})
	errorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "grpc_errors_total",
		Help: "Total number of gRPC errors",
	})
	responseTimeSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "grpc_response_time_seconds",
		Help: "Summary of gRPC response times",
	})
	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: []float64{0.1, 0.5, 1, 2, 5},
	})
	ordersReturned = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "returns_rate_total",
		Help: "Rate of product returns",
	})
	orderTotalPrice = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "order_total_value",
		Help:       "Summary of order total values",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01},
	})
	ordersCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "Total number of created orders",
	})
	packagingUsage = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "packaging_usage_total",
		Help: "Total number of times a packaging type has been used",
	}, []string{"packaging_name"})
	grpcRequestByStatusCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_request_count_total",
		Help: "Number of gRPC requests, labeled by status",
	}, []string{"status"})
	grpcRequestByMethod = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_by_method_total",
			Help: "Total number of gRPC requests by method",
		},
		[]string{"method"},
	)
)

// SetResponseTimeSummary updates response time metric
func SetResponseTimeSummary(responseTime float64) {
	responseTimeSummary.Observe(responseTime)
}

// SetRequestCounter updates request count metric
func SetRequestCounter() {
	requestCounter.Inc()
}

// SetErrorCounter updates errors count metric
func SetErrorCounter() {
	errorCounter.Inc()
}

// RequestDuration updates request duration metric
func RequestDuration(duration float64) {
	requestDuration.Observe(duration)
}

// SetOrdersReturned updates returned orders metric
func SetOrdersReturned() {
	ordersReturned.Inc()
}

// SetOrderTotalPrice updates total price metric
func SetOrderTotalPrice(price int64) {
	orderTotalPrice.Observe(float64(price))
}

// SetOrdersCreated updates created orders metric
func SetOrdersCreated() {
	ordersCreated.Inc()
}

// SetPackagingUsage updates packaging usage metric
func SetPackagingUsage(packagingName string) {
	packagingUsage.WithLabelValues(packagingName).Inc()
}

// SetGrpcRequestCountWithStatus updates request count by status metric
func SetGrpcRequestCountWithStatus(status string) {
	grpcRequestByStatusCount.WithLabelValues(status).Inc()
}

// SetGrpcRequestByMethodCount updates grpc requests by method metric
func SetGrpcRequestByMethodCount(method string) {
	grpcRequestByMethod.WithLabelValues(method).Inc()
}

func init() {
	prometheus.MustRegister(
		requestCounter,
		errorCounter,
		responseTimeSummary,
		requestDuration,
		ordersReturned,
		orderTotalPrice,
		ordersCreated,
		packagingUsage,
		grpcRequestByStatusCount,
		grpcRequestByMethod,
	)
}

// StartMetricsServer starts prometheus service to collect metrics
func StartMetricsServer(errChan chan<- error) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Starting metrics server on :9001")
		if err := http.ListenAndServe(":9001", nil); err != nil {
			errChan <- err
		}
	}()
}
