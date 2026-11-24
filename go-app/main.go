package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"runtime"
	"time"
)

var startTime = time.Now()

type StatusInfo struct {
	Hostname  string `json:"hostname"`
	PodIP     string `json:"pod_ip"`
	Namespace string `json:"namespace"`
	UID       int    `json:"uid"`
	GID       int    `json:"gid"`
	Arch      string `json:"architecture"`
	GoVersion string `json:"go_version"`
	Uptime    string `json:"uptime"`
	Timestamp string `json:"timestamp"`
}

const pageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>KCSA / CKS Multi-Env Go Service</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #0f172a; color: #f8fafc; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; }
    .card { background: #1e293b; border: 1px solid #334155; border-radius: 14px; padding: 2rem 2.5rem; max-width: 540px; width: 100%; box-shadow: 0 20px 35px rgba(0,0,0,0.4); }
    .badge { display: inline-block; padding: 0.3rem 0.75rem; border-radius: 999px; font-size: 0.8rem; font-weight: 600; margin-bottom: 1rem; }
    .secure { background: rgba(16, 185, 129, 0.15); color: #34d399; border: 1px solid #059669; }
    .root { background: rgba(245, 158, 11, 0.15); color: #fbbf24; border: 1px solid #d97706; }
    h1 { margin: 0 0 0.5rem; font-size: 1.5rem; }
    p { color: #94a3b8; margin: 0 0 1.5rem; font-size: 0.95rem; }
    table { width: 100%; border-collapse: collapse; }
    td { padding: 0.65rem 0; border-bottom: 1px solid #334155; font-size: 0.9rem; }
    td:first-child { color: #94a3b8; font-weight: 500; }
    td:last-child { text-align: right; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; color: #38bdf8; }
    tr:last-child td { border-bottom: none; }
    .footer { margin-top: 1.5rem; font-size: 0.8rem; color: #64748b; text-align: center; }
    a { color: #38bdf8; text-decoration: none; }
  </style>
</head>
<body>
  <div class="card">
    {{if eq .UID 0}}
      <span class="badge root">⚠️ Running as Root (UID 0)</span>
    {{else}}
      <span class="badge secure">🛡️ Restricted Non-Root Mode (UID {{.UID}})</span>
    {{end}}
    <h1>Go Multi-Stage Microservice</h1>
    <p>Live Kubernetes Pod Runtime &amp; Security Telemetry</p>
    <table>
      <tr><td>Pod Hostname</td><td>{{.Hostname}}</td></tr>
      <tr><td>UID / GID</td><td>{{.UID}} / {{.GID}}</td></tr>
      <tr><td>OS / Arch</td><td>{{.Arch}}</td></tr>
      <tr><td>Go Runtime</td><td>{{.GoVersion}}</td></tr>
      <tr><td>Uptime</td><td>{{.Uptime}}</td></tr>
      <tr><td>Server Time</td><td>{{.Timestamp}}</td></tr>
    </table>
    <div class="footer">JSON Telemetry: <a href="/healthz">/healthz</a></div>
  </div>
</body>
</html>`

func getStatus() StatusInfo {
	hostname, _ := os.Hostname()
	return StatusInfo{
		Hostname:  hostname,
		UID:       os.Getuid(),
		GID:       os.Getgid(),
		Arch:      fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		GoVersion: runtime.Version(),
		Uptime:    time.Since(startTime).Round(time.Second).String(),
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
	}
}

func main() {
	tmpl := template.Must(template.New("index").Parse(pageHTML))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tmpl.Execute(w, getStatus())
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(getStatus())
	})

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
		os.Exit(1)
	}
}
