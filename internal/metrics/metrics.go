package metrics

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Metric представляет метрику в памяти и в файле
type Metric struct {
	Name   string `json:"name"`
	Value  int64  `json:"value"`
	Status string `json:"status"`
}

var (
	metricsFile = filepath.Join(".", "metrics.json")
	metricsData = make(map[string]*Metric)

	// Определение счетчиков для Prometheus
	issuedOrdersCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "issued_orders_total",
			Help: "Total number of issued orders",
		},
		[]string{"status"},
	)
	createdReturnsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "created_returns_total",
			Help: "Total number of created returns",
		},
		[]string{"status"},
	)
)

func init() {
	// Регистрируем счетчики в Prometheus
	prometheus.MustRegister(issuedOrdersCounter)
	prometheus.MustRegister(createdReturnsCounter)

	// Загружаем метрики из файла при старте приложения
	if err := LoadMetrics(); err != nil {
		log.Printf("Ошибка при загрузке метрик: %v", err)
	}
}

// LoadMetrics загружает метрики из JSON-файла или создает файл с начальными значениями, если его нет
func LoadMetrics() error {
	if _, err := os.Stat(metricsFile); os.IsNotExist(err) {
		log.Println("Файл метрик не найден. Создание нового файла с начальными значениями.")

		metricsData["issued_orders_total_created"] = &Metric{Name: "issued_orders_total", Value: 0, Status: "created"}
		metricsData["created_returns_total_created"] = &Metric{Name: "created_returns_total", Value: 0, Status: "created"}

		if err := SaveMetrics(); err != nil {
			log.Printf("Ошибка при сохранении метрик: %v", err)
			return err
		}
		log.Println("Файл метрик успешно создан и заполнен начальными значениями.")
		return nil
	}

	file, err := os.Open(metricsFile)
	if err != nil {
		log.Printf("Ошибка при открытии файла метрик: %v", err)
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&metricsData)
	if err != nil {
		log.Printf("Ошибка при декодировании метрик из файла: %v", err)
		return err
	}

	log.Println("Файл метрик успешно загружен.")

	for _, metric := range metricsData {
		switch metric.Name {
		case "issued_orders_total":
			issuedOrdersCounter.WithLabelValues(metric.Status).Add(float64(metric.Value))
		case "created_returns_total":
			createdReturnsCounter.WithLabelValues(metric.Status).Add(float64(metric.Value))
		}
	}

	log.Println("Метрики успешно восстановлены в Prometheus из файла.")
	return nil
}

// SaveMetrics сохраняет метрики в JSON-файл
func SaveMetrics() error {
	log.Println("Попытка создания и записи в файл метрик...")

	file, err := os.Create(metricsFile)
	if err != nil {
		log.Printf("Ошибка при создании файла метрик: %v", err)
		return err
	}
	defer file.Close()

	log.Println("Файл метрик успешно создан, начинаем запись данных...")

	err = json.NewEncoder(file).Encode(metricsData)
	if err != nil {
		log.Printf("Ошибка при кодировании метрик в JSON: %v", err)
		return err
	}
	log.Println("Метрики успешно сохранены в файл.")
	return nil
}

// IncrementIssuedOrders увеличивает счетчик выданных заказов для Prometheus и сохраняет в JSON
func IncrementIssuedOrders(status string) {
	metricKey := "issued_orders_total_" + status
	if metric, exists := metricsData[metricKey]; exists {
		metric.Value++
	} else {
		metricsData[metricKey] = &Metric{Name: "issued_orders_total", Value: 1, Status: status}
	}
	issuedOrdersCounter.WithLabelValues(status).Inc()
	if err := SaveMetrics(); err != nil {
		log.Printf("Ошибка при сохранении метрик после IncrementIssuedOrders: %v", err)
	}
}

// IncrementCreatedReturns увеличивает счетчик созданных возвратов для Prometheus и сохраняет в JSON
func IncrementCreatedReturns(status string) {
	metricKey := "created_returns_total_" + status
	if metric, exists := metricsData[metricKey]; exists {
		metric.Value++
	} else {
		metricsData[metricKey] = &Metric{Name: "created_returns_total", Value: 1, Status: status}
	}
	createdReturnsCounter.WithLabelValues(status).Inc()
	if err := SaveMetrics(); err != nil {
		log.Printf("Ошибка при сохранении метрик после IncrementCreatedReturns: %v", err)
	}
}

// StartMetricsServer запускает HTTP-сервер для экспорта метрик Prometheus
func StartMetricsServer(addr string) {
	http.Handle("/metrics", promhttp.Handler()) // Экспорт метрик по пути /metrics
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("Ошибка при запуске сервера метрик: %v", err)
		}
	}()
}
