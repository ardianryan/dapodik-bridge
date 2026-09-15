# Catatan Rilis & Perubahan (Changelog)

Semua perubahan penting pada proyek **`dapodik-bridge`** akan didokumentasikan dalam berkas ini.
Format penulisan mengacu pada [Keep a Changelog](https://keepachangelog.com/id/1.0.0/) dan mengikuti [Semantic Versioning (SemVer)](https://semver.org/lang/id/).

---

## [1.1.0] - 2026-09-15

### Ditambahkan
- **Mode Desktop System Tray (`DapodikBridge-gui-windows-amd64.exe`)**:
  - Kompilasi Windows GUI tanpa konsol (`-H=windowsgui`), menghilangkan jendela hitam CMD secara total.
  - Ikon status dinamis di pojok kanan bawah taskbar (System Tray) dengan menu interaktif.
  - Kebal penutupan tidak sengaja (*anti-close*): tombol close tidak mematikan service, melainkan sembunyi di System Tray.
- **Autostart Otomatis Bebas UAC**:
  - Pendaftaran startup otomatis saat booting pada level user (`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`).
  - Tidak memicu popup UAC (*Run as Admin*) sehingga dijamin berjalan otomatis dan *silent* ke Tray.
  - Opsi toggle aktif/nonaktif autostart langsung dari klik kanan menu System Tray.
- **Dukungan Native Windows Service (`services.msc`)**:
  - Integrasi `kardianos/service` untuk komputer server sekolah 24/7.
  - Subperintah CLI: `service install`, `service uninstall`, `service start`, `service stop`, dan `service status`.
  - Dukungan restart otomatis jika terjadi kegagalan sistem (*auto-recovery*).
- **Otomatisasi Keamanan Dependabot**:
  - Konfigurasi `.github/dependabot.yml` untuk memantau pembaruan berkala dependensi Go (`gomod`) dan workflow (`github-actions`).
  - Penambahan audit kerentanan otomatis Go (`govulncheck`) pada alur CI `test.yml`.

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
