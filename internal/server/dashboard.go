package server

import (
	"net/http"
	"strconv"

	"github.com/ardianryan/dapodik-bridge/internal/config"
)

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(dashboardHTML(s.cfg, config.AppVersion)))
}

func dashboardHTML(cfg *config.Config, version string) string {
	return `<!DOCTYPE html>
<html lang="id" class="dark">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Dapodik Bridge - Dashboard</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <script>
    tailwind.config = {
      darkMode: 'class',
      theme: {
        extend: {
          colors: {
            brand: {
              500: '#0284c7',
              600: '#0369a1',
            }
          }
        }
      }
    }
  </script>
</head>
<body class="min-h-screen flex flex-col bg-slate-950 text-slate-100 antialiased selection:bg-sky-500 selection:text-white">

  <!-- Top Navigation Header -->
  <header class="border-b border-slate-800 bg-slate-900 sticky top-0 z-40">
    <div class="max-w-5xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-sky-600 flex items-center justify-center text-white font-bold">
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 19L4 12C4 7.58172 7.58172 4 12 4C16.4183 4 20 7.58172 20 12L20 19"/>
            <path d="M4 12L20 12"/>
            <circle cx="12" cy="12" r="3"/>
          </svg>
        </div>
        <div class="flex items-center gap-2">
          <h1 class="font-semibold text-sm text-white">Dapodik Bridge</h1>
          <span class="px-2 py-0.5 rounded text-[11px] font-mono bg-slate-800 text-slate-300 border border-slate-700">v` + version + `</span>
        </div>
      </div>

      <!-- Quick Status Badge -->
      <div class="flex items-center gap-3">
        <div id="status-pill" class="flex items-center gap-2 px-2.5 py-1 rounded text-xs font-medium bg-emerald-950/60 text-emerald-400 border border-emerald-800">
          <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
          <span id="status-pill-text">Memeriksa...</span>
        </div>
        <a href="https://github.com/ardianryan/dapodik-bridge" target="_blank" rel="noopener noreferrer" class="text-xs text-slate-400 hover:text-white px-2 py-1 rounded hover:bg-slate-800 transition-colors">
          GitHub
        </a>
      </div>
    </div>
  </header>

  <!-- Main Container -->
  <main class="flex-1 max-w-5xl mx-auto w-full px-4 sm:px-6 py-6 space-y-6">

    <!-- Status & Info Card -->
    <div class="bg-slate-900 border border-slate-800 rounded-xl p-6">
      <div class="flex flex-col md:flex-row items-center justify-between gap-6">
        
        <div class="space-y-2 text-center md:text-left">
          <div class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded text-xs font-medium bg-slate-800 text-slate-300 border border-slate-700">
            <svg class="w-3.5 h-3.5 text-sky-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/></svg>
            <span>Strict Read-Only Mode: <strong>Enforced</strong></span>
          </div>
          <h2 class="text-xl sm:text-2xl font-bold text-white tracking-tight">
            Dapodik Local PostgreSQL Bridge
          </h2>
          <p class="text-xs sm:text-sm text-slate-400 max-w-xl leading-relaxed">
            Menghubungkan aplikasi web sekolah ke database lokal Dapodik secara aman tanpa risiko penulisan yang merusak data validasi.
          </p>
          
          <div class="flex flex-wrap items-center justify-center md:justify-start gap-4 pt-2 text-xs text-slate-400 font-mono">
            <div>Port: <strong class="text-slate-200">` + strconv.Itoa(cfg.Port) + `</strong></div>
            <div>Target DB: <strong class="text-slate-200">` + cfg.DBHost + `:` + strconv.Itoa(cfg.DBPort) + `</strong></div>
            <div>Database: <strong class="text-slate-200">` + cfg.DBName + `</strong></div>
          </div>
        </div>

        <!-- Ping Test Button -->
        <div class="flex flex-col items-center gap-2 flex-shrink-0">
          <button 
            id="toggle-bridge-btn"
            onclick="triggerTestHealth()" 
            class="px-5 py-2.5 rounded-lg bg-sky-600 hover:bg-sky-500 text-white font-medium text-xs flex items-center gap-2 active:bg-sky-700 transition-colors"
            title="Klik untuk tes koneksi"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
            <span>Tes Koneksi</span>
          </button>
          <span id="ping-label" class="text-xs font-mono text-slate-400">Menghubungkan...</span>
        </div>

      </div>
    </div>

    <!-- Quick Action Bar -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <button onclick="triggerSyncNow()" class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-4 rounded-xl flex flex-col items-start gap-2 transition-colors">
        <div class="w-7 h-7 rounded-md bg-sky-950 text-sky-400 flex items-center justify-center">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
        </div>
        <div class="text-left">
          <h4 class="text-xs font-semibold text-white">Push Webhook</h4>
          <p class="text-[11px] text-slate-400">Kirim data ke cloud</p>
        </div>
      </button>

      <a href="/api/v1/health" target="_blank" class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-4 rounded-xl flex flex-col items-start gap-2 transition-colors">
        <div class="w-7 h-7 rounded-md bg-emerald-950 text-emerald-400 flex items-center justify-center">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
        </div>
        <div class="text-left">
          <h4 class="text-xs font-semibold text-white">Healthcheck API</h4>
          <p class="text-[11px] text-slate-400">Status daemon & DB</p>
        </div>
      </a>

      <a href="/api/v1/schema/tables" target="_blank" class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-4 rounded-xl flex flex-col items-start gap-2 transition-colors">
        <div class="w-7 h-7 rounded-md bg-purple-950 text-purple-400 flex items-center justify-center">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2 1 3 3 3h10c2 0 3-1 3-3V7c0-2-1-3-3-3H7C5 4 4 5 4 7z"/><path d="M9 12h6M9 16h6M9 8h6"/></svg>
        </div>
        <div class="text-left">
          <h4 class="text-xs font-semibold text-white">Schema Tables</h4>
          <p class="text-[11px] text-slate-400">Inspeksi tabel</p>
        </div>
      </a>

      <a href="/api/v1/pip" target="_blank" class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-4 rounded-xl flex flex-col items-start gap-2 transition-colors">
        <div class="w-7 h-7 rounded-md bg-amber-950 text-amber-400 flex items-center justify-center">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/></svg>
        </div>
        <div class="text-left">
          <h4 class="text-xs font-semibold text-white">Data Bansos PIP</h4>
          <p class="text-[11px] text-slate-400">Penerima bantuan</p>
        </div>
      </a>
    </div>

    <!-- Telemetry & Logs -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">

      <!-- Left: Endpoint Explorer -->
      <div class="bg-slate-900 border border-slate-800 rounded-xl p-5 space-y-3">
        <div class="flex items-center justify-between pb-2 border-b border-slate-800">
          <h3 class="font-semibold text-xs text-slate-200">
            Daftar Endpoint REST API
          </h3>
          <span class="text-[11px] font-mono text-slate-500">v1</span>
        </div>

        <div class="space-y-1.5 text-xs font-mono">
          <div class="p-2 rounded bg-slate-950 border border-slate-800/80 flex items-center justify-between">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/health</span></span>
            <span class="text-[11px] text-slate-500 font-sans">Status daemon</span>
          </div>
          <div class="p-2 rounded bg-slate-950 border border-slate-800/80 flex items-center justify-between">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/kesejahteraan</span></span>
            <span class="text-[11px] text-slate-500 font-sans">PIP, KIP, PKH, KKS</span>
          </div>
          <div class="p-2 rounded bg-slate-950 border border-slate-800/80 flex items-center justify-between">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/rapor</span></span>
            <span class="text-[11px] text-slate-500 font-sans">Nilai per semester</span>
          </div>
          <div class="p-2 rounded bg-slate-950 border border-slate-800/80 flex items-center justify-between">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/siswa/komprehensif</span></span>
            <span class="text-[11px] text-slate-500 font-sans">Data periodik & ortu</span>
          </div>
          <div class="p-2 rounded bg-slate-950 border border-slate-800/80 flex items-center justify-between">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/gtk/lengkap</span></span>
            <span class="text-[11px] text-slate-500 font-sans">Guru & Tenaga Kependidikan</span>
          </div>
          <div class="p-2 rounded bg-slate-950 border border-slate-800/80 flex items-center justify-between">
            <span class="text-sky-400 font-semibold">POST <span class="text-slate-200">/api/v1/sync/push</span></span>
            <span class="text-[11px] text-slate-500 font-sans">Trigger Webhook Push</span>
          </div>
        </div>
      </div>

      <!-- Right: Realtime Console Logs -->
      <div class="bg-slate-900 border border-slate-800 rounded-xl p-5 space-y-3 flex flex-col">
        <div class="flex items-center justify-between pb-2 border-b border-slate-800">
          <h3 class="font-semibold text-xs text-slate-200">Live Console Logs</h3>
          <button onclick="clearLogs()" class="text-[11px] text-slate-400 hover:text-white px-2 py-0.5 rounded bg-slate-800 transition-colors">Clear</button>
        </div>

        <div id="terminal-logs" class="flex-1 bg-slate-950 rounded-lg p-3 font-mono text-[11px] leading-relaxed text-slate-300 space-y-1 overflow-y-auto max-h-[240px] border border-slate-800">
          <div class="text-slate-500">// Dapodik Bridge v` + version + ` daemon online.</div>
          <div class="text-sky-400">[INFO] Listening on http://localhost:` + strconv.Itoa(cfg.Port) + `</div>
          <div class="text-emerald-400">[INFO] Strict read-only transaction mode active.</div>
        </div>
      </div>

    </div>

  </main>

  <!-- Footer -->
  <footer class="border-t border-slate-800 py-4 text-center text-xs text-slate-500 mt-auto">
    <p>Dapodik Read-Only Bridge Daemon • Lisensi MIT-NC • SMAN 1 Gedeg & Ryan Ardian</p>
  </footer>

  <!-- Interactivity Script -->
  <script>
    function addLog(message, type = 'info') {
      const term = document.getElementById('terminal-logs');
      const timeStr = new Date().toLocaleTimeString();
      const div = document.createElement('div');
      
      let colorClass = 'text-slate-300';
      if (type === 'success') colorClass = 'text-emerald-400';
      if (type === 'warn') colorClass = 'text-amber-400';
      if (type === 'error') colorClass = 'text-rose-400';
      
      div.className = colorClass;
      div.textContent = '[' + timeStr + '] ' + message;
      term.appendChild(div);
      term.scrollTop = term.scrollHeight;
    }

    function clearLogs() {
      const term = document.getElementById('terminal-logs');
      term.innerHTML = '<div class="text-slate-500">// Log cleared.</div>';
    }

    async function triggerTestHealth() {
      const pingLabel = document.getElementById('ping-label');
      const statusPill = document.getElementById('status-pill');
      const statusPillText = document.getElementById('status-pill-text');
      
      pingLabel.textContent = 'Pinging...';
      const start = performance.now();
      
      try {
        const res = await fetch('/api/v1/health');
        const end = performance.now();
        const duration = Math.round(end - start);
        const data = await res.json();
        
        if (data.data && data.data.database_up) {
          pingLabel.textContent = 'Active (' + duration + ' ms)';
          pingLabel.className = 'text-xs font-mono text-emerald-400 font-medium';
          statusPill.className = 'flex items-center gap-2 px-2.5 py-1 rounded text-xs font-medium bg-emerald-950/60 text-emerald-400 border border-emerald-800';
          statusPillText.textContent = 'Connected';
          addLog('Healthcheck OK (' + duration + 'ms): DB connected at ' + data.data.database_host, 'success');
        } else {
          pingLabel.textContent = 'Bridge UP, DB Standby';
          pingLabel.className = 'text-xs font-mono text-amber-400 font-medium';
          statusPill.className = 'flex items-center gap-2 px-2.5 py-1 rounded text-xs font-medium bg-amber-950/60 text-amber-400 border border-amber-800';
          statusPillText.textContent = 'Degraded';
          addLog('Bridge is running, awaiting PostgreSQL connection at ' + (data.data?.database_host || '127.0.0.1:5432'), 'warn');
        }
      } catch (err) {
        pingLabel.textContent = 'Offline';
        pingLabel.className = 'text-xs font-mono text-rose-400 font-medium';
        addLog('Healthcheck request failed: ' + err.message, 'error');
      }
    }

    async function triggerSyncNow() {
      addLog('Triggering manual push webhook...', 'info');
      try {
        const res = await fetch('/api/v1/sync/push?type=welfare', { method: 'POST' });
        const data = await res.json();
        addLog('Webhook push response: ' + JSON.stringify(data.data || data.message || 'Done'), 'success');
      } catch (err) {
        addLog('Sync webhook failed: ' + err.message, 'error');
      }
    }

    document.addEventListener('DOMContentLoaded', () => {
      triggerTestHealth();
    });
  </script>
</body>
</html>`
}
