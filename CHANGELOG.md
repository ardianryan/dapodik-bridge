# Catatan Rilis & Perubahan (Changelog)

Semua perubahan penting pada proyek **`dapodik-bridge`** akan didokumentasikan dalam berkas ini.
Format penulisan mengacu pada [Keep a Changelog](https://keepachangelog.com/id/1.0.0/) dan mengikuti [Semantic Versioning (SemVer)](https://semver.org/lang/id/).

---

## [1.0.0] - 2026-09-14

### Ditambahkan
- **Inisialisasi Daemon Dapodik Bridge**: Implementasi bridge daemon berkinerja tinggi berbasis Go (Golang) berukuran ~12 MB dan hemat memori (< 15 MB RAM).
- **Port Default Kustom (4712)**: Menghindari bentrok port dengan Dapodik (5774), Postgres (5432), atau web server lokal.
- **Jaminan 3 Lapis Strict Read-Only (`SELECT` only)**:
  - Eksekusi otomatis `SET default_transaction_read_only = on; SET transaction_read_only = on;` pada setiap koneksi.
  - Query sanitizer untuk memblokir seluruh kata kunci DDL/DML (`INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `TRUNCATE`).
  - Prepared statement dengan parameter `$1, $2, ...` untuk keamanan penuh terhadap SQL Injection.
- **Koneksi Cerdas & Fallback Password**: Otomatis mencoba koneksi ke `127.0.0.1:5432` dengan fallback password bawaan Dapodik (`""`, `"password"`, `"dapodik"`, `"dapodik123"`, `"123456"`, `"admin"`).
- **Endpoint Kesejahteraan Siswa Lengkap (`GET /api/v1/kesejahteraan`)**: Mengambil data seluruh bantuan siswa (PIP, KIP, PKH, KKS, KPS, PBI, Beasiswa) dengan filter jenis dan paginasi.
- **Endpoint Khusus PIP/KIP (`GET /api/v1/pip`)**: Shortcut cepat untuk filter data PIP dan KIP.
- **Endpoint Nilai Rapor Multi-Semester (`GET /api/v1/rapor`)**: Mengambil nilai rapor per semester (Semester 1–5 untuk PPDB Prestasi / e-Rapor).
- **Endpoint Siswa Komprehensif (`GET /api/v1/siswa/komprehensif`)**: Data siswa lengkap beserta data periodik (TB/BB, jarak sekolah, titik koordinat lintang bujur zonasi) dan profil penghasilan orang tua.
- **Endpoint GTK Lengkap (`GET /api/v1/gtk/lengkap`)**: Profil guru dan tenaga kependidikan lengkap dengan NUPTK, NIP, status kepegawaian, dan bidang studi.
- **Endpoint Rombongan Belajar (`GET /api/v1/rombel`)**: Daftar rombel beserta jumlah total anggota siswa aktif.
- **Dynamic Schema Inspector (`GET /api/v1/schema/tables` & `GET /api/v1/schema/columns`)**: Inspeksi skema tabel dan kolom database Dapodik secara langsung melalui API.
- **Push Webhook ke Server Cloud (`POST /api/v1/sync/push`)**: Pengiriman data lokal otomatis ke webhook server web sekolah pusat.
- **Middleware**: Dukungan CORS, Bearer Token Auth (`BRIDGE_API_KEY`), dan reverse proxy IP headers (`X-Forwarded-For`, `X-Real-IP`).
- **Otomasi CI/CD GitHub Actions**:
  - Validasi unit test dan linting otomatis (`test.yml`).
  - Rilis otomatis multiplatform untuk Windows 64-bit/32-bit (`.exe`), Linux, macOS Apple Silicon/Intel beserta SHA256 checksums (`release.yml`).
