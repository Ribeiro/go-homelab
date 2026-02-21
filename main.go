package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Definindo as métricas
var (
	cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "homelab_cpu_usage_threads",
		Help: "Número de Goroutines (threads leves) em uso no MacBook 2011",
	})
	// Nova métrica de memória RAM
	memUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "homelab_memory_heap_alloc_bytes",
		Help: "Memória RAM (Heap) alocada pelo processo em bytes",
	})
)

type LogEntry struct {
	TS    string `json:"ts"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
	App   string `json:"app"`
	Host  string `json:"host"`
}

func main() {
	// Endpoint para o Prometheus coletar as métricas
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Printf("📊 Endpoint de métricas ativo na porta :8080/metrics\n")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Erro ao iniciar servidor HTTP: %s\n", err)
		}
	}()

	for {
		// Captura estatísticas de memória do runtime do Go
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// Atualiza as métricas no Prometheus
		currentThreads := runtime.NumGoroutine()
		heapAllocBytes := m.HeapAlloc

		cpuUsage.Set(float64(currentThreads))
		memUsage.Set(float64(heapAllocBytes))

		// Converte bytes para MB para facilitar a leitura no log
		allocMB := float64(heapAllocBytes) / 1024 / 1024

		entry := LogEntry{
			TS:    time.Now().Format(time.RFC3339),
			Level: "info",
			Msg:   fmt.Sprintf("Métricas: Threads=%d, RAM=%.2fMB", currentThreads, allocMB),
			App:   "go-homelab",
			Host:  "debian-vm",
		}

		payload, _ := json.Marshal(entry)
		fmt.Println(string(payload))

		time.Sleep(10 * time.Second)
	}
}
