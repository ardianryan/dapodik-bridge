# Panduan Kontribusi (Contributing Guidelines)

Terima kasih atas minat Anda untuk berkontribusi pada pengembangan **`dapodik-bridge`**! Proyek ini bersifat *open-source* dan kami menyambut baik kontribusi dalam bentuk pelaporan bug, usulan fitur, pembaruan dokumentasi, maupun pengiriman kode (*Pull Request*).

---

## 🛠️ Alur Pengembangan Lokal (Local Setup)

### Prasyarat
- [Go (Golang)](https://go.dev/) (versi >= 1.22)
- Git
- Make (opsional)

### Langkah Menjalankan Proyek
1. **Fork** repositori ini ke akun GitHub Anda.
2. **Clone** hasil fork ke komputer lokal:
   ```bash
   git clone https://github.com/ardianryan/dapodik-bridge.git
   cd dapodik-bridge
   ```
3. **Unduh dependensi**:
   ```bash
   go mod download
   ```
4. **Jalankan Unit Test**:
   ```bash
   go test -v ./...
   ```
5. **Jalankan Daemon Lokal**:
   ```bash
   go run ./cmd/bridge -port=4712
   ```

---

## 🌿 Standar Branch & Commit

### 1. Penamaan Branch
Gunakan format nama branch yang deskriptif:
- `feat/<nama-fitur>`: Penambahan fitur baru
- `fix/<nama-bug>`: Perbaikan bug
- `docs/<nama-dokumen>`: Perbaikan atau penambahan dokumentasi
- `refactor/<nama-komponen>`: Refaktor kode tanpa mengubah fungsionalitas

### 2. Konvensi Pesan Commit (Conventional Commits)
Kami mengikuti standar [Conventional Commits](https://www.conventionalcommits.org/):
- `feat: add schema column inspector endpoint`
- `fix: correct port parsing when env var is empty`
- `docs: update reverse proxy examples in README`
- `test: add unit test for strict read-only query sanitizer`

---

## 🧪 Aturan Pengujian & Kualitas Kode

Sebelum mengajukan *Pull Request* (PR), pastikan:
1. Menjalankan `go fmt ./...` dan `go vet ./...`.
2. Seluruh unit test lolos 100%:
   ```bash
   go test -v -race ./...
   ```
3. Memastikan binary cross-compile berjalan lancar:
   ```bash
   make build-all
   ```
