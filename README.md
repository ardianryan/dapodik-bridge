# 🌉 Dapodik Bridge

[![Test & Lint](https://github.com/ardianryan/dapodik-bridge/actions/workflows/test.yml/badge.svg)](https://github.com/ardianryan/dapodik-bridge/actions/workflows/test.yml)
[![Release](https://github.com/ardianryan/dapodik-bridge/actions/workflows/release.yml/badge.svg)](https://github.com/ardianryan/dapodik-bridge/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev)
[![Default Port](https://img.shields.io/badge/Default%20Port-4712-orange.svg)](#konfigurasi-port)

Daemon jembatan lokal (*local bridge daemon*) berbasis **Go (Golang)** berkinerja tinggi, berukuran ringkas (~12 MB single binary), dan hemat memori (< 15 MB RAM) yang membaca database internal PostgreSQL aplikasi **Dapodik Desktop** secara **Strict Read-Only (`SELECT` only)**.

`dapodik-bridge` dirancang khusus untuk kebutuhan pengembangan **Website Sekolah, Portal PPDB Afirmasi, Sistem Monitoring Bansos/Kesejahteraan Siswa (PIP, KIP, PKH, KKS, KPS), serta Sinkronisasi Nilai Rapor Multi-Semester**.

---

## 📑 Daftar Isi
- [Arsitektur Sistem](#-arsitektur-sistem)
- [Mengapa Go & Bukan Framework Berat?](#-mengapa-go--bukan-framework-berat)
- [Keamanan & Jaminan Strict Read-Only](#-keamanan--jaminan-strict-read-only)
- [Daftar Lengkap Endpoint REST API](#-daftar-lengkap-endpoint-rest-api)
  - [1. Healthcheck (`GET /api/v1/health`)](#1-healthcheck-get-apiv1health)
  - [2. Bansos & Kesejahteraan Lengkap (`GET /api/v1/kesejahteraan`)](#2-bansos--kesejahteraan-lengkap-get-apiv1kesejahteraan)
  - [3. Shortcut Khusus PIP/KIP (`GET /api/v1/pip`)](#3-shortcut-khusus-pipkip-get-apiv1pip)
  - [4. Nilai Rapor Multi-Semester (`GET /api/v1/rapor`)](#4-nilai-rapor-multi-semester-get-apiv1rapor)
  - [5. Siswa Komprehensif (`GET /api/v1/siswa/komprehensif`)](#5-siswa-komprehensif-get-apiv1siswakomprehensif)
  - [6. GTK & Tenaga Kependidikan (`GET /api/v1/gtk/lengkap`)](#6-gtk--tenaga-kependidikan-get-apiv1gtklengkap)
  - [7. Rombongan Belajar (`GET /api/v1/rombel`)](#7-rombongan-belajar-get-apiv1rombel)
  - [8. Schema Database Inspector (`GET /api/v1/schema/tables` & `columns`)](#8-schema-database-inspector)
  - [9. Sinkronisasi Push ke Cloud Web Sekolah (`POST /api/v1/sync/push`)](#9-sinkronisasi-push-ke-cloud-web-sekolah-post-apiv1syncpush)
- [Panduan Instalasi & Menjalankan](#-panduan-instalasi--menjalankan)
  - [Untuk Operator Dapodik (Windows)](#untuk-operator-dapodik-windows)
  - [Menjalankan sebagai Background Service Windows (NSSM)](#menjalankan-sebagai-background-service-windows-nssm)
- [Konfigurasi Reverse Proxy (Port 4712)](#-konfigurasi-reverse-proxy-port-4712)
  - [1. Cloudflare Tunnel (cloudflared)](#1-cloudflare-tunnel-cloudflared-direkomendasikan)
  - [2. Nginx Reverse Proxy](#2-nginx-reverse-proxy)
- [Variabel Lingkungan (.env) & CLI Flags](#-variabel-lingkungan-env--cli-flags)

---

## 🏛 Arsitektur Sistem

```mermaid
flowchart LR
    subgraph PC_Operator["Laptop / Komputer Operator Sekolah"]
        DAPODIK["Dapodik Desktop\n(PostgreSQL 127.0.0.1:5432)"]
        BRIDGE["dapodik-bridge.exe\n(Listen Port :4712)\n[Strict Read-Only]"]
        TUNNEL["Cloudflare Tunnel /\nReverse Proxy Client"]
        
        DAPODIK -->|Local IPC / SELECT only| BRIDGE
        BRIDGE -->|HTTP :4712| TUNNEL
    end

    subgraph Cloud_Server["Cloud / Web Server Sekolah"]
        WEB["Website Sekolah / PPDB\n(Laravel, Node.js, Next.js, dsb.)"]
        WEBHOOK["Endpoint /api/dapodik/webhook"]
    end

    TUNNEL -->|Secure TLS Tunnel| WEB
    BRIDGE -.->|Auto-Push Data (POST)| WEBHOOK
```

---

## 🚀 Mengapa Go & Bukan Framework Berat?

1. **Zero Runtime Dependency**: Tidak memerlukan instalasi PHP, Composer, Laragon, XAMPP, Node.js, atau Python di laptop operator sekolah. Cukup unduh satu file `.exe` dan langsung klik dua kali untuk menjalankan.
2. **Super Ringan**: Penggunaan memori hanya **~12 - 15 MB RAM**, tidak membebani laptop operator sekolah yang sedang menjalankan Dapodik Desktop.
3. **Instan Boot**: Menyala dalam hitungan milidetik.
4. **Port Kustom Aman (`4712`)**: Port default disetel ke `4712` agar tidak konflik dengan port default Dapodik (`5774`), PostgreSQL (`5432`), atau Apache/Nginx (`80/8080`).

---

## 🛡 Keamanan & Jaminan Strict Read-Only

`dapodik-bridge` menerapkan **3 lapis pengamanan read-only** untuk menjamin database Dapodik tidak pernah mengalami modifikasi, kerusakan, maupun diskualifikasi validasi:

1. **Session Enforcement**: Pada setiap koneksi PostgreSQL yang dibuka, bridge mengeksekusi:
   ```sql
   SET default_transaction_read_only = on;
   SET transaction_read_only = on;
   ```
   PostgreSQL secara native akan menolak query penulisan apa pun (`ERROR: cannot execute INSERT/UPDATE/DELETE in a read-only transaction`).
2. **Query Level Sanitization**: Seluruh query internal hanya menggunakan perintah `SELECT`. Query yang mengandung kata kunci DDL/DML seperti `INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `TRUNCATE` langsung diblokir sebelum mencapai database.
3. **Prepared Statements**: Seluruh parameter pencarian menggunakan parameterized placeholder (`$1`, `$2`, dst.) sehingga 100% bebas dari risiko SQL Injection.

---

## 📡 Daftar Lengkap Endpoint REST API

Semua respons endpoint dibungkus (*enveloped*) dalam format JSON standar berikut:

```json
{
  "success": true,
  "data": [...],
  "meta": {
    "total_count": 120,
    "count": 50,
    "page": 1,
    "per_page": 50,
    "bridge_version": "1.0.0",
    "execution_ms": 14
  },
  "timestamp": "2026-09-14T08:45:00Z"
}
```

---

### 1. Healthcheck (`GET /api/v1/health`)
Memeriksa status daemon bridge, koneksi ke PostgreSQL Dapodik lokal, uptime, dan penegakan status read-only. Endpoint ini tidak memerlukan autentikasi API Key sehingga dapat digunakan oleh probe reverse proxy.

```bash
curl http://localhost:4712/api/v1/health
```

**Contoh Response:**
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "bridge_up": true,
    "database_up": true,
    "database_host": "127.0.0.1:5432",
    "database_name": "dapodik_dasmen",
    "database_user": "postgres",
    "read_only_enforced": true,
    "uptime": "2h14m30s",
    "version": "1.0.0",
    "current_time": "2026-09-14T08:45:00Z"
  }
}
```

---

### 2. Bansos & Kesejahteraan Lengkap (`GET /api/v1/kesejahteraan`)
Mengambil seluruh data siswa penerima bantuan kesejahteraan sosial (tidak hanya PIP, tetapi juga **KIP, PKH, KKS, KPS, PBI, dan Beasiswa**).

* **Query Parameters**:
  * `jenis`: Filter jenis bantuan (contoh: `pip`, `kip`, `pkh`, `kks`, `kps`). Jika dikosongkan, menampilkan semua jenis bantuan.
  * `page`: Nomor halaman (default: `1`).
  * `per_page`: Jumlah data per halaman (default: `50`, maks `500`).

```bash
# Mengambil semua bantuan
curl http://localhost:4712/api/v1/kesejahteraan

# Filter hanya penerima PKH
curl "http://localhost:4712/api/v1/kesejahteraan?jenis=pkh&page=1&per_page=100"
```

**Contoh Data:**
```json
{
  "peserta_didik_id": "9a31b234-5678-4abc-def0-1234567890ab",
  "nisn": "0081234567",
  "nama": "ADITYA PRATAMA",
  "jenis_kelamin": "L",
  "tingkat_kelas": "10",
  "rombel": "X TJKT 1",
  "jenis_bantuan": "Program Keluarga Harapan (PKH)",
  "nomor_kartu": "PKH-357801-09876",
  "nama_di_kartu": "SITI AMINAH",
  "tahun_mulai": 2023,
  "tahun_selesai": 2026,
  "status_aktif": true,
  "layak_pip": true,
  "alasan_layak_pip": "Pemegang PKH/KPS/KKS",
  "nama_ibu_kandung": "SITI AMINAH"
}
```

---

### 3. Shortcut Khusus PIP/KIP (`GET /api/v1/pip`)
Alias cepat untuk mengambil data siswa pemegang Kartu Indonesia Pintar (KIP) dan usulan Program Indonesia Pintar (PIP).

```bash
curl "http://localhost:4712/api/v1/pip?page=1&per_page=50"
```

---

### 4. Nilai Rapor Multi-Semester (`GET /api/v1/rapor`)
Mengambil nilai rapor siswa per semester untuk keperluan PPDB jalur prestasi nilai rapor (Semester 1 s.d. 5) atau sinkronisasi nilai e-Rapor ke website sekolah.

* **Query Parameters**:
  * `semester_id`: ID semester (contoh: `20231` = 2023/2024 Ganjil, `20232` = 2023/2024 Genap).
  * `rombel_id`: UUID Rombel tertentu (opsional).
  * `nisn`: Filter per NISN siswa (opsional).
  * `page`, `per_page`: Paginasi data.

```bash
curl "http://localhost:4712/api/v1/rapor?semester_id=20231&nisn=0081234567"
```

**Contoh Data:**
```json
{
  "peserta_didik_id": "9a31b234-5678-4abc-def0-1234567890ab",
  "nisn": "0081234567",
  "nama_siswa": "ADITYA PRATAMA",
  "rombel": "X TJKT 1",
  "semester_id": "20231",
  "mata_pelajaran": "Matematika (Umum)",
  "nilai_angka": 88.5,
  "nilai_huruf": "A",
  "predikat": "Sangat Baik",
  "kkm": 75.0
}
```

---

### 5. Siswa Komprehensif (`GET /api/v1/siswa/komprehensif`)
Menyajikan data siswa lengkap dengan data periodik (tinggi/berat badan, lingkar kepala, jarak rumah ke sekolah, koordinat lintang bujur untuk PPDB Zonasi) dan pekerjaan/penghasilan orang tua.

* **Query Parameters**:
  * `q`: Pencarian berdasarkan Nama, NISN, atau NIK siswa.
  * `rombel_id`: Filter per rombongan belajar.
  * `page`, `per_page`: Paginasi data.

```bash
curl "http://localhost:4712/api/v1/siswa/komprehensif?q=Aditya"
```

---

### 6. GTK & Tenaga Kependidikan (`GET /api/v1/gtk/lengkap`)
Mengambil data guru dan tenaga kependidikan, NUPTK, NIP, status kepegawaian, bidang studi pengajaran, dan kontak resmi.

* **Query Parameters**:
  * `q`: Pencarian nama atau NUPTK/NIP.
  * `page`, `per_page`: Paginasi.

```bash
curl "http://localhost:4712/api/v1/gtk/lengkap"
```

---

### 7. Rombongan Belajar (`GET /api/v1/rombel`)
Mengambil daftar rombel, tingkat kelas, jurusan, wali kelas, beserta total jumlah anggota siswa aktif.

```bash
curl "http://localhost:4712/api/v1/rombel?semester_id=20241"
```

---

### 8. Schema Database Inspector
Memungkinkan pengembang web sekolah mengecek langsung tabel dan kolom apa saja yang tersedia di dalam database Dapodik lokal secara dinamis:

* **Daftar Tabel**:
  ```bash
  curl http://localhost:4712/api/v1/schema/tables
  ```
* **Daftar Kolom Tabel**:
  ```bash
  curl "http://localhost:4712/api/v1/schema/columns?table=peserta_didik"
  ```

---

### 9. Sinkronisasi Push ke Cloud Web Sekolah (`POST /api/v1/sync/push`)
Memicu daemon bridge untuk mengemas data lokal Dapodik dan mengirimkannya via HTTP POST (Webhook) ke server web sekolah pusat di internet.

* **Query Parameters / Body JSON**:
  * `type`: Jenis data yang dikirim (`welfare`, `rapor`, `siswa`, `gtk`).
  * `target_url`: URL webhook penerima (opsional jika sudah diatur di `.env`).
  * `secret_token`: Token autentikasi webhook.

```bash
curl -X POST "http://localhost:4712/api/v1/sync/push?type=welfare" \
  -H "Content-Type: application/json" \
  -d '{
    "target_url": "https://sekolah.sch.id/api/dapodik/webhook",
    "secret_token": "rahasia-sekolah-123"
  }'
```

**Payload yang Diterima Server Web Sekolah:**
```json
{
  "source_host": "PC-OPERATOR-DAPODIK",
  "school_npsn": "20500000",
  "push_type": "welfare",
  "record_count": 142,
  "pushed_at": "2026-09-14T08:45:00Z",
  "data": [...]
}
```

---

## 💻 Panduan Instalasi & Menjalankan

### Untuk Operator Dapodik (Windows)

1. Buka halaman **[Releases](https://github.com/ardianryan/dapodik-bridge/releases)**.
2. Unduh file binary `dapodik-bridge-windows-amd64.exe` (untuk Windows 64-bit) atau `dapodik-bridge-windows-386.exe` (untuk Windows 32-bit).
3. Letakkan di folder mana saja (misal: `C:\dapodik-bridge\`).
4. Jalankan melalui Command Prompt atau klik dua kali:
   ```cmd
   dapodik-bridge-windows-amd64.exe
   ```
5. Daemon langsung menyala di port `4712` dan otomatis mencari kredensial Dapodik lokal:
   ```text
   [INFO] Initializing Dapodik Bridge on port 4712 (host: 0.0.0.0)...
   [INFO] Connected to Dapodik PostgreSQL at 127.0.0.1:5432/dapodik_dasmen (read-only enforced)
   [INFO] Dapodik Bridge Daemon is listening on http://0.0.0.0:4712
   ```

### Menjalankan sebagai Background Service Windows (NSSM)

Agar bridge otomatis berjalan setiap kali laptop/server operator menyala tanpa perlu membuka terminal:

1. Unduh [NSSM (Non-Sucking Service Manager)](https://nssm.cc/).
2. Buka Command Prompt Administrator dan ketik:
   ```cmd
   nssm install DapodikBridge C:\dapodik-bridge\dapodik-bridge-windows-amd64.exe
   nssm start DapodikBridge
   ```

---

## 🌐 Konfigurasi Reverse Proxy (Port 4712)

Karena database PostgreSQL Dapodik terkunci di `127.0.0.1` dan IP publik sekolah umumnya berada di balik NAT/Indihome/Biznet, gunakan Reverse Proxy agar web sekolah di cloud bisa mengakses API Bridge.

### 1. Cloudflare Tunnel (`cloudflared`) (Direkomendasikan)
Sangat mudah, gratis, dan tidak membutuhkan IP Publik Statis maupun port forwarding modem.

1. Jalankan tunnel mengarah ke port `4712`:
   ```bash
   cloudflared tunnel --url http://localhost:4712
   ```
2. Atau pasang di `config.yml` cloudflared:
   ```yaml
   ingress:
     - hostname: bridge.sekolah.sch.id
       service: http://localhost:4712
     - service: http_status:404
   ```

### 2. Nginx Reverse Proxy
Jika komputer Dapodik berada dalam satu jaringan LAN dengan server lokal:

```nginx
server {
    listen 80;
    server_name bridge.sekolah.lan;

    location / {
        proxy_pass http://127.0.0.1:4712;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## ⚙️ Variabel Lingkungan (.env) & CLI Flags

Seluruh opsi dapat diatur via file `.env` maupun argumen perintah (*flags*):

| CLI Flag | Variabel `.env` | Default | Deskripsi |
| :--- | :--- | :--- | :--- |
| `-port` | `BRIDGE_PORT` | `4712` | Port HTTP yang didengarkan bridge |
| `-host` | `BRIDGE_HOST` | `0.0.0.0` | Host bind address |
| `-api-key` | `BRIDGE_API_KEY` | *(kosong)* | Kunci Bearer Token untuk melindungi endpoint |
| `-db-host` | `DB_HOST` | `127.0.0.1` | Host database PostgreSQL Dapodik |
| `-db-port` | `DB_PORT` | `5432` | Port database Dapodik |
| `-db-user` | `DB_USER` | `postgres` | Username database Dapodik |
| `-db-pass` | `DB_PASSWORD` | *(auto)* | Password DB (otomatis dicoba password bawaan jika kosong) |
| `-db-name` | `DB_NAME` | `dapodik_dasmen` | Nama database Dapodik lokal |
| `-push-url` | `PUSH_TARGET_URL` | *(kosong)* | URL server cloud target push sinkronisasi |
| `-push-secret`| `PUSH_SECRET_TOKEN`| *(kosong)* | Token autentikasi webhook push |
| `-npsn` | `SCHOOL_NPSN` | *(kosong)* | Nomor Pokok Sekolah Nasional |

Contoh menjalankan dengan flags:
```bash
./dapodik-bridge -port=4712 -api-key=supersecret123
```

---

## 📄 Lisensi
Didistribusikan di bawah lisensi [MIT](LICENSE). Dibuat untuk mendukung transparansi dan kemudahan pengelolaan data pendidikan di Indonesia.
