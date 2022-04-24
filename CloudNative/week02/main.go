package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"httpserver/metrics"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func rootFunc(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("<h1>Test Page</h1>"))

	// 接收客户端 request，并将 request 中带的 header 写入 response header
	for key, values := range r.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}

	// 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	os.Setenv("VERSION", "v0.1")
	version := os.Getenv("VERSION")
	w.Header().Set("Version", version)

	// Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	clientIp := getClientIp(r)
	fmt.Printf("Client IP: %s, Status Code: %d\n", clientIp, 200)

	w.Write([]byte("<h1>Test Page</h1>"))
}

func getClientIp(r *http.Request) (ip string) {
	addr := r.Header.Get("X-Forwarded-For")
	if addr != "" {
		ip = strings.Split(addr, ",")[0]
		return
	}

	addr = r.Header.Get("X-Real-IP")
	if addr != "" {
		return addr
	}

	addr = r.RemoteAddr
	ip = strings.Split(addr, ":")[0]
	return
}

func delay(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	randTime := rand.Intn(3000)
	time.Sleep(time.Millisecond * time.Duration(randTime))
	w.Write([]byte(fmt.Sprintf("I have delay %d", randTime)))
}

func main() {
	metrics.Register()
	mux := http.NewServeMux()
	mux.HandleFunc("/delay", delay)
	mux.HandleFunc("/", rootFunc)

	// 当访问 localhost/healthz 时，应返回 200
	mux.HandleFunc("/healthz", healthz)
	mux.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe("0.0.0.0:80", mux)
	if err != nil {
		fmt.Printf("Start Server Error: %v\n", err)
	}
}
