# Kebijakan Keamanan & Kepatuhan Perlindungan Data (Security & Data Protection Policy)

Kami memprioritaskan integritas perangkat lunak serta perlindungan data pribadi dalam ekosistem pendidikan Indonesia. Dokumen ini menjelaskan kebijakan penanganan kerentanan keamanan, jaminan *Strict Read-Only*, dan kepatuhan hukum atas penggunaan `dapodik-bridge`.

---

## ⚖️ Kepatuhan UU Perlindungan Data Pribadi (UU PDP No. 27 Tahun 2022)

Database internal Dapodik Desktop memuat **Data Pribadi Spesifik dan Umum** (seperti NIK, NISN, nama lengkap, riwayat bansos/PIP/KIP, koordinat tempat tinggal, data orang tua/wali siswa, dan profil guru/GTK).

Setiap pengguna, operator, atau pengembang yang memanfaatkan daemon **`dapodik-bridge`** diwajibkan:
1. **Menghormati Dasar Pemrosesan Data**: Mengakses dan memproses data hanya untuk kepentingan sah institusi pendidikan/sekolah terkait dengan mandat kedinasan yang sah.
2. **Menjaga Kerahasiaan Port & API Key**:
   - Pasang API Key (`BRIDGE_API_KEY`) jika bridge dapat diakses di luar `localhost`.
   - Gunakan Reverse Proxy terenkripsi (seperti Cloudflare Tunnel atau HTTPS) saat menghubungkan ke cloud web sekolah.
3. **Mencegah Pengungkapan Tanpa Izin**: Dilarang mengekspos data pribadi siswa atau guru ke domain publik tanpa mekanisme autentikasi dan otorisasi yang ketat.

> [!CAUTION]
> **Sanksi Hukum**: Segala bentuk penyalahgunaan, pembocoran, atau pemrosesan data pribadi secara melawan hukum dapat dikenakan sanksi pidana dan denda administratif sesuai ketentuan **UU Perlindungan Data Pribadi No. 27 Tahun 2022 Pasal 67**.

---

## 🛡️ Jaminan Keamanan Sistem Database (Strict Read-Only)

`dapodik-bridge` dirancang dengan arsitektur **Zero Modification** terhadap database Dapodik:
1. Setiap sesi koneksi PostgreSQL secara otomatis di-lock dengan:
   ```sql
   SET default_transaction_read_only = on;
   SET transaction_read_only = on;
   ```
2. Seluruh permintaan query DDL/DML (`INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `TRUNCATE`) ditolak secara native oleh database server PostgreSQL dan query sanitizer internal bridge.
3. Penggunaan daemon ini tidak merubah struktur tabel, trigger, indeks, maupun catatan riwayat sinkronisasi Dapodik pusat.

---

## 🛡️ Versi yang Didukung (Supported Versions)

| Versi Daemon | Status Pemeliharaan Keamanan |
| :--- | :--- |
| **`1.x.x`** | ✅ Didukung Secara Penuh (*Active Support*) |
| `< 1.0.0` | ❌ Tidak Didukung |

---

## 🚨 Melaporkan Kerentanan Keamanan (Reporting a Vulnerability)

Jika Anda menemukan potensi celah keamanan (*vulnerability*) atau bug keamanan pada daemon ini, mohon **JANGAN** membuat *public issue* di GitHub.

Silakan laporkan secara privat melalui kanal berikut:
1. **Email Pengembang**: Kontak langsung ke [inisaya@ardianryan.com](mailto:inisaya@ardianryan.com) atau DM resmi via Instagram [@smansagewithai](https://www.instagram.com/smansagewithai/).
2. **GitHub Security Advisory**: Gunakan fitur *Report a vulnerability* di tab **Security** repositori GitHub kami.

Tim pemelihara akan merespons laporan dalam waktu maksimal **2 x 24 jam** dan segera merilis perbaikan (*patch*) pada versi rilis berikutnya.
