<p align="center">
  <img src="https://dapo.kemendikdasmen.go.id/assets/logo-dapodik-BZDG7c6h.png" alt="Dapodik Logo" width="140" />
</p>

<h1 align="center">dapodik-bridge</h1>

<p align="center">
  <a href="https://github.com/ardianryan/dapodik-bridge/releases"><img src="https://img.shields.io/github/v/release/ardianryan/dapodik-bridge?style=flat-square" alt="GitHub release" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT--NC-blue.svg?style=flat-square" alt="License: MIT-NC" /></a>
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-%3E%3D1.22-00ADD8.svg?style=flat-square&logo=go" alt="Go version" /></a>
  <a href="#-konfigurasi-port-kustom-4712"><img src="https://img.shields.io/badge/Port-4712-orange.svg?style=flat-square" alt="Default Port" /></a>
  <a href="https://github.com/ardianryan/dapodik-bridge/actions/workflows/test.yml"><img src="https://img.shields.io/github/actions/workflow/status/ardianryan/dapodik-bridge/test.yml?branch=main&style=flat-square&label=Tests" alt="Test Status" /></a>
  <a href="https://www.instagram.com/smansagewithai/"><img src="https://img.shields.io/badge/Instagram-@smansagewithai-E4405F.svg?style=flat-square&logo=instagram&logoColor=white" alt="Instagram" /></a>
</p>

<p align="center">
  Daemon jembatan lokal (<i>local bridge daemon</i>) berbasis <b>Go (Golang)</b> berkinerja tinggi, berukuran ringkas (~12 MB <i>single binary</i>), dan hemat memori (&lt; 15 MB RAM) yang membaca database internal PostgreSQL aplikasi <b>Dapodik Desktop</b> secara <b>Strict Read-Only (<code>SELECT</code> only)</b>.
</p>

<p align="center">
  Dipublikasikan dan dikelola oleh <b>SMA Negeri 1 Gedeg (<a href="https://www.instagram.com/smansagewithai/">@smansagewithai</a>)</b><br />
  Dikembangkan oleh <b>Ryan Ardian</b> (<a href="mailto:inisaya@ardianryan.com">inisaya@ardianryan.com</a>)
</p>

---

> [!IMPORTANT]
> ### 📢 Pernyataan Penyangkalan (Disclaimer) & Misi Terbuka
> **`dapodik-bridge` adalah perangkat lunak *Unofficial* (tidak resmi) dan independen.** Daemon ini dikembangkan sebagai inisiatif komunitas sumber terbuka (*open-source*) oleh **SMA Negeri 1 Gedeg** dan **Ryan Ardian**, tanpa afiliasi langsung secara struktural dengan Kementerian Pendidikan Dasar dan Menengah (Kemendikdasmen).
>
> **Tujuan & Misi Pengembangan**:
> Daemon ini lahir atas semangat memajukan transformasi digital dan interoperabilitas sistem informasi sekolah di Indonesia. Tujuan utamanya adalah **memberdayakan pengembang web sekolah, portal PPDB afirmasi/zonasi/prestasi, e-Rapor, serta sistem monitoring bansos kesejahteraan siswa** agar dapat mengakses data yang tidak terekspos oleh antarmuka bawaan secara aman, instan, hemat memori, dan tanpa membebani komputer kerja operator Dapodik.
>
> Seluruh hak cipta nama, logo, dan merek dagang **Dapodik (Data Pokok Pendidikan)** adalah milik sah **Kementerian Pendidikan Dasar dan Menengah Republik Indonesia**.

---

## 📑 Daftar Isi
- [Arsitektur Sistem](#-arsitektur-sistem)
- [Mengapa Go & Bukan Framework Berat?](#-mengapa-go--bukan-framework-berat)
- [Keamanan & Jaminan Strict Read-Only](#-keamanan--jaminan-strict-read-only)
- [Kepatuhan UU Perlindungan Data Pribadi](#-kepatuhan-uu-perlindungan-data-pribadi-uu-pdp-no-272022)
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
- [Lisensi & Ketentuan Penggunaan Non-Komersial](#-lisensi--ketentuan-penggunaan-non-komersial)
- [Dokumen Pendukung Repositori](#-dokumen-pendukung-repositori)

---

## 🏛 Arsitektur Sistem

```mermaid
flowchart LR
    subgraph PC_Operator["Laptop / Komputer Operator Sekolah"]
        DAPODIK["Dapodik Desktop<br/>(PostgreSQL 127.0.0.1:5432)"]
        BRIDGE["dapodik-bridge.exe<br/>(Listen Port :4712)<br/>[Strict Read-Only]"]
        TUNNEL["Cloudflare Tunnel /<br/>Reverse Proxy Client"]
        
        DAPODIK -->|"Local IPC / SELECT only"| BRIDGE
        BRIDGE -->|"HTTP :4712"| TUNNEL
    end

    subgraph Cloud_Server["Cloud / Web Server Sekolah"]
        WEB["Website Sekolah / PPDB<br/>(Laravel, Node.js, Next.js, dsb.)"]
        WEBHOOK["Endpoint /api/dapodik/webhook"]
    end

    TUNNEL -->|"Secure TLS Tunnel"| WEB
    BRIDGE -.->|"Auto-Push Data (POST)"| WEBHOOK
```

---

## 🚀 Mengapa Go & Bukan Framework Berat?

1. **Zero Runtime Dependency**: Tidak membutuhkan instalasi PHP, Composer, Laragon, XAMPP, Node.js, atau Python di laptop operator sekolah. Cukup unduh satu file `.exe` dan langsung klik dua kali untuk menjalankan.
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

## ⚠️ Kepatuhan UU Perlindungan Data Pribadi (UU PDP No. 27/2022)

Database internal Dapodik Desktop memuat **Data Pribadi Spesifik dan Umum** (seperti NIK, NISN, nama lengkap, riwayat bansos/PIP/KIP, koordinat tempat tinggal, data orang tua/wali siswa, dan profil guru/GTK).

> [!CAUTION]
> **Sanksi Pidana & Administratif**: Setiap operator dan pengembang wajib mematuhi **UU Perlindungan Data Pribadi No. 27 Tahun 2022 Pasal 67**. Gunakan data Dapodik semata-mata untuk kepentingan sah institusi pendidikan, lindungi API Key daemon, dan dilarang keras memperjualbelikan atau mengekspos data siswa/guru ke publik tanpa enkripsi dan hak akses yang sah.

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

## ⚖️ Lisensi & Ketentuan Penggunaan Non-Komersial

Proyek ini dirilis di bawah lisensi **[MIT License with Non-Commercial Restriction (MIT-NC)](LICENSE)**.

### 📌 Ketentuan Penggunaan:
1. **100% Gratis untuk Pendidikan**: Daemon ini sepenuhnya **gratis** digunakan oleh seluruh sekolah, guru, operator, siswa, akademisi, dan lembaga pendidikan di Indonesia.
2. **Dilarang untuk Tujuan Komersial (Non-Commercial Only)**:
   - Dilarang keras memperjualbelikan, memonetisasi, menjual kembali (*reselling*), atau mengemas daemon ini ke dalam produk perangkat lunak berbayar / layanan berbayar pihak ketiga tanpa izin tertulis dari pemegang hak cipta (**Ryan Ardian & SMA Negeri 1 Gedeg**).
3. **Atribusi Hak Cipta**:
   - Hak Cipta &copy; 2026 **Ryan Ardian** ([inisaya@ardianryan.com](mailto:inisaya@ardianryan.com)) & **SMA Negeri 1 Gedeg** ([@smansagewithai](https://www.instagram.com/smansagewithai/)).

---

## 📑 Dokumen Pendukung Repositori

- 📜 [Changelog](CHANGELOG.md) - Catatan riwayat versi dan perubahan rilis.
- 🛡️ [Security Policy & UU PDP](SECURITY.md) - Kebijakan keamanan, penegakan read-only & kepatuhan perlindungan data pribadi.
- 🤝 [Contributing Guidelines](CONTRIBUTING.md) - Panduan kontribusi kode, standar Go, dan Pull Request.
- 📜 [Code of Conduct](CODE_OF_CONDUCT.md) - Kode etik komunitas kontributor.
