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

// Definindo as métricas para o Prometheus
var (
	cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "homelab_cpu_usage_threads",
		Help: "Número de Goroutines (threads leves) em uso no MacBook 2011",
	})
	// Métrica de memória RAM alocada
	memUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "homelab_memory_heap_alloc_bytes",
		Help: "Memória RAM (Heap) alocada pelo processo em bytes",
	})
)

// Estrutura de log compatível com o Grafana Loki (JSON)
type LogEntry struct {
	TS    string `json:"ts"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
	App   string `json:"app"`
	Host  string `json:"host"`
}

func main() {
	// Goroutine para o endpoint de métricas (Porta 8080)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Printf("📊 Endpoint de métricas ativo na porta :8080/metrics\n")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("Erro ao iniciar servidor HTTP: %s\n", err)
		}
	}()

	// Loop principal de coleta e logging
	for {
		// Captura estatísticas de memória do runtime do Go
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// Coleta dados atuais
		currentThreads := runtime.NumGoroutine()
		heapAllocBytes := m.HeapAlloc

		// Atualiza as métricas que o Prometheus irá ler
		cpuUsage.Set(float64(currentThreads))
		memUsage.Set(float64(heapAllocBytes))

		// Converte bytes para MB para facilitar a leitura no log
		allocMB := float64(heapAllocBytes) / 1024 / 1024

		// Cria a entrada de log estruturada
		entry := LogEntry{
			TS:    time.Now().Format(time.RFC3339),
			Level: "info",
			// MENSAGEM DE TESTE PARA O GITOPS V2:
			Msg:   fmt.Sprintf("Homelab V2 - GitOps funcional! Metrics: Threads=%d, RAM=%.2fMB", currentThreads, allocMB),
			App:   "go-homelab",
			Host:  "debian-vm",
		}

		// Serializa para JSON e imprime no stdout (onde o Promtail/Loki captura)
		payload, _ := json.Marshal(entry)
		fmt.Println(string(payload))

		// Intervalo de 10 segundos para não sobrecarregar o MacBook 2011
		time.Sleep(10 * time.Second)
	}
}
