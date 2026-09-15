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
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
  <script src="https://cdn.tailwindcss.com"></script>
  <script>
    tailwind.config = {
      darkMode: 'class',
      theme: {
        extend: {
          fontFamily: {
            sans: ['"Plus Jakarta Sans"', 'sans-serif'],
            mono: ['"JetBrains Mono"', 'monospace'],
          },
          colors: {
            brand: {
              50: '#f0f9ff',
              400: '#38bdf8',
              500: '#0ea5e9',
              600: '#0284c7',
            }
          }
        }
      }
    }
  </script>
  <style>
    body {
      background-color: #090d16;
      color: #f1f5f9;
    }
    .glass-panel {
      background: rgba(17, 24, 39, 0.75);
      backdrop-filter: blur(12px);
      border: 1px solid rgba(255, 255, 255, 0.08);
    }
    .glass-card {
      background: rgba(30, 41, 59, 0.5);
      border: 1px solid rgba(255, 255, 255, 0.06);
    }
    .glow-cyan {
      box-shadow: 0 0 35px -5px rgba(14, 165, 233, 0.25);
    }
    .glow-green {
      box-shadow: 0 0 35px -5px rgba(16, 185, 129, 0.25);
    }
    .warp-toggle {
      transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    }
  </style>
</head>
<body class="min-h-screen flex flex-col antialiased selection:bg-sky-500 selection:text-white">

  <!-- Ambient Light Orbs -->
  <div class="fixed inset-0 pointer-events-none overflow-hidden -z-10">
    <div class="absolute -top-40 left-1/2 -translate-x-1/2 w-[600px] h-[350px] bg-sky-500/10 blur-[130px] rounded-full"></div>
    <div class="absolute bottom-0 right-10 w-[400px] h-[300px] bg-emerald-500/5 blur-[120px] rounded-full"></div>
  </div>

  <!-- Top Navigation Header -->
  <header class="border-b border-slate-800/80 glass-panel sticky top-0 z-40">
    <div class="max-w-5xl mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-sky-400 to-blue-600 flex items-center justify-center shadow-lg shadow-sky-500/20 text-white font-bold">
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 19L4 12C4 7.58172 7.58172 4 12 4C16.4183 4 20 7.58172 20 12L20 19"/>
            <path d="M4 12L20 12"/>
            <circle cx="12" cy="12" r="3"/>
          </svg>
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h1 class="font-bold text-sm text-white tracking-tight">Dapodik Bridge</h1>
            <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-sky-500/10 text-sky-400 border border-sky-500/20">v` + version + `</span>
          </div>
          <p class="text-[11px] text-slate-400">Local Daemon & Push Webhook</p>
        </div>
      </div>

      <!-- Quick Status Badge -->
      <div class="flex items-center gap-3">
        <div id="status-pill" class="flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          <span id="status-pill-text">Connected</span>
        </div>
        <a href="https://github.com/ardianryan/dapodik-bridge" target="_blank" class="text-xs text-slate-400 hover:text-white px-2.5 py-1 rounded-lg hover:bg-slate-800 transition-colors">
          GitHub
        </a>
      </div>
    </div>
  </header>

  <!-- Main Container -->
  <main class="flex-1 max-w-5xl mx-auto w-full px-4 sm:px-6 py-8 space-y-6">

    <!-- HERO SECTION: Cloudflare WARP / Tailscale Style Big Toggle Card -->
    <div class="glass-panel rounded-3xl p-6 sm:p-8 glow-cyan relative overflow-hidden">
      <div class="flex flex-col md:flex-row items-center justify-between gap-6 relative z-10">
        
        <div class="space-y-2 text-center md:text-left">
          <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium bg-slate-800/80 text-slate-300 border border-slate-700/60">
            <svg class="w-3.5 h-3.5 text-sky-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/></svg>
            <span>Strict Read-Only Mode: <strong>ENFORCED</strong></span>
          </div>
          <h2 class="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
            Dapodik Local PostgreSQL Bridge
          </h2>
          <p class="text-xs sm:text-sm text-slate-400 max-w-lg leading-relaxed">
            Menghubungkan aplikasi web sekolah / CBT ke database lokal Dapodik di port <strong class="text-sky-300">` + strconv.Itoa(cfg.Port) + `</strong> secara aman tanpa risiko kerusakan data validasi.
          </p>
          
          <div class="flex flex-wrap items-center justify-center md:justify-start gap-4 pt-2 text-xs text-slate-400">
            <div class="flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-sky-400"></span>
              <span>Port: <strong class="text-white">` + strconv.Itoa(cfg.Port) + `</strong></span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
              <span>Target DB: <strong class="text-white">` + cfg.DBHost + `:` + strconv.Itoa(cfg.DBPort) + `</strong></span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-purple-400"></span>
              <span>Database: <strong class="text-white">` + cfg.DBName + `</strong></span>
            </div>
          </div>
        </div>

        <!-- Big Interactive WARP-Style Button Switch -->
        <div class="flex flex-col items-center gap-3 flex-shrink-0">
          <button 
            id="toggle-bridge-btn"
            onclick="triggerTestHealth()" 
            class="warp-toggle w-24 h-24 sm:w-28 sm:h-28 rounded-3xl bg-gradient-to-br from-sky-500 to-blue-600 hover:from-sky-400 hover:to-blue-500 shadow-xl shadow-sky-500/25 flex flex-col items-center justify-center gap-1 text-white active:scale-95 transition-all group"
            title="Klik untuk tes koneksi"
          >
            <svg class="w-8 h-8 sm:w-9 sm:h-9 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
            <span class="text-[11px] font-bold tracking-wide uppercase">PING TEST</span>
          </button>
          <span id="ping-label" class="text-[11px] font-mono text-emerald-400 font-semibold">● Active (0 ms)</span>
        </div>

      </div>
    </div>

    <!-- Quick Action Bar -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <button onclick="triggerSyncNow()" class="glass-card hover:bg-slate-800/60 p-4 rounded-2xl flex flex-col items-start gap-2 transition-all hover:-translate-y-0.5">
        <div class="w-8 h-8 rounded-xl bg-sky-500/10 text-sky-400 flex items-center justify-center">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
        </div>
        <div class="text-left">
          <h4 class="text-xs font-bold text-white">Push Webhook</h4>
          <p class="text-[11px] text-slate-400">Sinkronisasi instan</p>
        </div>
      </button>

      <a href="/api/v1/health" target="_blank" class="glass-card hover:bg-slate-800/60 p-4 rounded-2xl flex flex-col items-start gap-2 transition-all hover:-translate-y-0.5">
        <div class="w-8 h-8 rounded-xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
        </div>
        <div class="text-left">
          <h4 class="text-xs font-bold text-white">Healthcheck API</h4>
          <p class="text-[11px] text-slate-400">Status daemon & DB</p>
        </div>
      </a>

      <a href="/api/v1/schema/tables" target="_blank" class="glass-card hover:bg-slate-800/60 p-4 rounded-2xl flex flex-col items-start gap-2 transition-all hover:-translate-y-0.5">
        <div class="w-8 h-8 rounded-xl bg-purple-500/10 text-purple-400 flex items-center justify-center">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2 1 3 3 3h10c2 0 3-1 3-3V7c0-2-1-3-3-3H7C5 4 4 5 4 7z"/><path d="M9 12h6M9 16h6M9 8h6"/></svg>
        </div>
        <div class="text-left">
          <h4 class="text-xs font-bold text-white">Schema Tables</h4>
          <p class="text-[11px] text-slate-400">Inspeksi tabel Dapodik</p>
        </div>
      </a>

      <a href="/api/v1/pip" target="_blank" class="glass-card hover:bg-slate-800/60 p-4 rounded-2xl flex flex-col items-start gap-2 transition-all hover:-translate-y-0.5">
        <div class="w-8 h-8 rounded-xl bg-amber-500/10 text-amber-400 flex items-center justify-center">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/></svg>
        </div>
        <div class="text-left">
          <h4 class="text-xs font-bold text-white">Data Bansos PIP</h4>
          <p class="text-[11px] text-slate-400">Cek siswa penerima</p>
        </div>
      </a>
    </div>

    <!-- Telemetry & Live Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">

      <!-- Left: Interactive API Explorer -->
      <div class="glass-panel rounded-3xl p-6 space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800">
          <h3 class="font-bold text-sm text-white flex items-center gap-2">
            <svg class="w-4 h-4 text-sky-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
            <span>Daftar Endpoint Siap Pakai</span>
          </h3>
          <span class="text-[10px] font-mono text-slate-400">REST v1</span>
        </div>

        <div class="space-y-2 text-xs font-mono">
          <div class="p-2.5 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between hover:border-slate-700 transition-colors">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/health</span></span>
            <span class="text-[10px] text-slate-500 font-sans">Status daemon</span>
          </div>
          <div class="p-2.5 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between hover:border-slate-700 transition-colors">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/kesejahteraan</span></span>
            <span class="text-[10px] text-slate-500 font-sans">PIP, KIP, PKH, KKS</span>
          </div>
          <div class="p-2.5 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between hover:border-slate-700 transition-colors">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/rapor</span></span>
            <span class="text-[10px] text-slate-500 font-sans">Nilai per semester</span>
          </div>
          <div class="p-2.5 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between hover:border-slate-700 transition-colors">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/siswa/komprehensif</span></span>
            <span class="text-[10px] text-slate-500 font-sans">Data periodik & ortu</span>
          </div>
          <div class="p-2.5 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between hover:border-slate-700 transition-colors">
            <span class="text-emerald-400 font-semibold">GET <span class="text-slate-200">/api/v1/gtk/lengkap</span></span>
            <span class="text-[10px] text-slate-500 font-sans">Guru & Tenaga Kependidikan</span>
          </div>
          <div class="p-2.5 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between hover:border-slate-700 transition-colors">
            <span class="text-sky-400 font-semibold">POST <span class="text-slate-200">/api/v1/sync/push</span></span>
            <span class="text-[10px] text-slate-500 font-sans">Trigger Webhook Push</span>
          </div>
        </div>
      </div>

      <!-- Right: Realtime Event Terminal (Tailscale Style) -->
      <div class="glass-panel rounded-3xl p-6 space-y-4 flex flex-col">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800">
          <div class="flex items-center gap-2">
            <div class="flex gap-1.5">
              <span class="w-2.5 h-2.5 rounded-full bg-rose-500/80"></span>
              <span class="w-2.5 h-2.5 rounded-full bg-amber-500/80"></span>
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-500/80"></span>
            </div>
            <h3 class="font-bold text-sm text-white pl-2">Live Console Logs</h3>
          </div>
          <button onclick="clearLogs()" class="text-[10px] text-slate-400 hover:text-white px-2 py-0.5 rounded bg-slate-800">Clear</button>
        </div>

        <div id="terminal-logs" class="flex-1 bg-black/60 rounded-2xl p-4 font-mono text-[11px] leading-relaxed text-slate-300 space-y-1.5 overflow-y-auto max-h-[260px] border border-white/5">
          <div class="text-slate-500">// Dapodik Bridge v` + version + ` daemon online.</div>
          <div class="text-sky-400">[INFO] Listening on http://localhost:` + strconv.Itoa(cfg.Port) + `</div>
          <div class="text-emerald-400">[INFO] Strict read-only transaction mode active.</div>
        </div>
      </div>

    </div>

  </main>

  <!-- Footer -->
  <footer class="border-t border-slate-800/80 py-6 text-center text-xs text-slate-500 mt-auto">
    <p>Dapodik Read-Only Bridge Daemon • Lisensi MIT-NC • SMAN 1 Gedeg & Ryan Ardian</p>
  </footer>

  <!-- Live Interactivity Script -->
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
          pingLabel.textContent = '● Connected (' + duration + ' ms)';
          pingLabel.className = 'text-[11px] font-mono text-emerald-400 font-semibold';
          statusPill.className = 'flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20';
          statusPillText.textContent = 'Connected';
          addLog('Healthcheck OK (' + duration + 'ms): DB connected at ' + data.data.database_host, 'success');
        } else {
          pingLabel.textContent = '● Bridge UP, DB Standby';
          pingLabel.className = 'text-[11px] font-mono text-amber-400 font-semibold';
          statusPill.className = 'flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-400 border border-amber-500/20';
          statusPillText.textContent = 'Degraded';
          addLog('Bridge is running, awaiting PostgreSQL connection at ' + (data.data?.database_host || '127.0.0.1:5432'), 'warn');
        }
      } catch (err) {
        pingLabel.textContent = '● Offline';
        pingLabel.className = 'text-[11px] font-mono text-rose-400 font-semibold';
        addLog('Healthcheck request failed: ' + err.message, 'error');
      }
    }

    async function triggerSyncNow() {
      addLog('Triggering manual push webhook...', 'info');
      try {
        const res = await fetch('/api/v1/sync/push?type=welfare', { method: 'POST' });
        const data = await res.json();
        addLog('Webhook push completed: ' + JSON.stringify(data.data || data.message || 'Done'), 'success');
      } catch (err) {
        addLog('Sync webhook failed: ' + err.message, 'error');
      }
    }

    // Auto-ping once on load
    document.addEventListener('DOMContentLoaded', () => {
      triggerTestHealth();
    });
  </script>
</body>
</html>`
}
