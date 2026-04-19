package metrics

func (m *Metrics) AddInFlight(method string, route string, delta float64) {
	m.requestsInFlight.WithLabelValues(method, route).Add(delta)
}

func (m *Metrics) ObserveRequest(method string, route string, status string, durationSeconds float64, requestSizeBytes int, responseSizeBytes int) {
	m.requestsTotal.WithLabelValues(method, route, status).Inc()
	m.requestDuration.WithLabelValues(method, route, status).Observe(durationSeconds)
	m.requestSizeBytes.WithLabelValues(method, route).Observe(float64(requestSizeBytes))
	m.responseSizeBytes.WithLabelValues(method, route, status).Observe(float64(responseSizeBytes))
}
