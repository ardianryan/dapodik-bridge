package models

import "time"

// Response is the standard JSON response envelope
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Meta      *MetaInfo   `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Error     string      `json:"error,omitempty"`
}

// MetaInfo contains pagination and context metadata
type MetaInfo struct {
	TotalCount  int    `json:"total_count,omitempty"`
	Count       int    `json:"count"`
	Page        int    `json:"page,omitempty"`
	PerPage     int    `json:"per_page,omitempty"`
	Database    string `json:"database,omitempty"`
	BridgeVer   string `json:"bridge_version"`
	ExecutionMs int64  `json:"execution_ms"`
}

// HealthStatus represents bridge and database health info
type HealthStatus struct {
	Status       string    `json:"status"` // "ok" or "degraded"
	BridgeUp     bool      `json:"bridge_up"`
	DatabaseUp   bool      `json:"database_up"`
	DatabaseHost string    `json:"database_host"`
	DatabaseName string    `json:"database_name"`
	DatabaseUser string    `json:"database_user"`
	ReadOnly     bool      `json:"read_only_enforced"`
	Uptime       string    `json:"uptime"`
	Version      string    `json:"version"`
	CurrentTime  time.Time `json:"current_time"`
}

// StudentWelfare represents PIP, KIP, PKH, KKS, and other welfare assistance
type StudentWelfare struct {
	PesertaDidikID string `json:"peserta_didik_id"`
	NISN           string `json:"nisn"`
	Nama           string `json:"nama"`
	JenisKelamin   string `json:"jenis_kelamin"`
	TingkatKelas   string `json:"tingkat_kelas,omitempty"`
	Rombel         string `json:"rombel,omitempty"`
	JenisBantuan   string `json:"jenis_bantuan"` // PIP, KIP, PKH, KKS, KPS, dsb
	NomorKartu     string `json:"nomor_kartu"`
	NamaDiKartu    string `json:"nama_di_kartu"`
	TahunMulai     int    `json:"tahun_mulai,omitempty"`
	TahunSelesai   int    `json:"tahun_selesai,omitempty"`
	StatusAktif    bool   `json:"status_aktif"`
	LayakPIP       bool   `json:"layak_pip"`
	AlasanLayakPIP string `json:"alasan_layak_pip,omitempty"`
	NamaIbuKandung string `json:"nama_ibu_kandung,omitempty"`
}

// RaporGrade represents multi-semester report card grades
type RaporGrade struct {
	PesertaDidikID string   `json:"peserta_didik_id"`
	NISN           string   `json:"nisn"`
	NamaSiswa      string   `json:"nama_siswa"`
	Rombel         string   `json:"rombel"`
	SemesterID     string   `json:"semester_id"` // e.g. 20231, 20232, 20241
	MataPelajaran  string   `json:"mata_pelajaran"`
	NilaiAngka     float64  `json:"nilai_angka"`
	NilaiHuruf     string   `json:"nilai_huruf,omitempty"`
	Predikat       string   `json:"predikat,omitempty"`
	KKM            *float64 `json:"kkm,omitempty"`
}

// StudentComprehensive represents student personal, periodic, and parent details
type StudentComprehensive struct {
	PesertaDidikID   string   `json:"peserta_didik_id"`
	Nama             string   `json:"nama"`
	NISN             string   `json:"nisn"`
	NIK              string   `json:"nik,omitempty"`
	JenisKelamin     string   `json:"jenis_kelamin"`
	TempatLahir      string   `json:"tempat_lahir"`
	TanggalLahir     string   `json:"tanggal_lahir"`
	Agama            string   `json:"agama,omitempty"`
	AlamatJalan      string   `json:"alamat_jalan,omitempty"`
	RT               string   `json:"rt,omitempty"`
	RW               string   `json:"rw,omitempty"`
	KelurahanDesa    string   `json:"kelurahan_desa,omitempty"`
	Kecamatan        string   `json:"kecamatan,omitempty"`
	KodePos          string   `json:"kode_pos,omitempty"`
	Lintang          *float64 `json:"lintang,omitempty"`
	Bujur            *float64 `json:"bujur,omitempty"`
	NamaAyah         string   `json:"nama_ayah,omitempty"`
	PekerjaanAyah    string   `json:"pekerjaan_ayah,omitempty"`
	PenghasilanAyah  string   `json:"penghasilan_ayah,omitempty"`
	NamaIbuKandung   string   `json:"nama_ibu_kandung,omitempty"`
	PekerjaanIbu     string   `json:"pekerjaan_ibu,omitempty"`
	PenghasilanIbu   string   `json:"penghasilan_ibu,omitempty"`
	TinggiBadanCm    int      `json:"tinggi_badan_cm,omitempty"`
	BeratBadanKg     int      `json:"berat_badan_kg,omitempty"`
	JarakKeSekolahKm float64  `json:"jarak_ke_sekolah_km,omitempty"`
	WaktuTempuhMenit int      `json:"waktu_tempuh_menit,omitempty"`
	RombelSaatIni    string   `json:"rombel_saat_ini,omitempty"`
	TingkatKelas     string   `json:"tingkat_kelas,omitempty"`
}

// GTKLengkap represents Teacher and Educational Personnel comprehensive data
type GTKLengkap struct {
	PTKID              string `json:"ptk_id"`
	Nama               string `json:"nama"`
	GelarDepan         string `json:"gelar_depan,omitempty"`
	GelarBelakang      string `json:"gelar_belakang,omitempty"`
	NIK                string `json:"nik,omitempty"`
	JenisKelamin       string `json:"jenis_kelamin"`
	TempatLahir        string `json:"tempat_lahir"`
	TanggalLahir       string `json:"tanggal_lahir"`
	NIP                string `json:"nip,omitempty"`
	NUPTK              string `json:"nuptk,omitempty"`
	StatusKepegawaian  string `json:"status_kepegawaian,omitempty"`
	JenisPTK           string `json:"jenis_ptk,omitempty"`
	Agama              string `json:"agama,omitempty"`
	Alamat             string `json:"alamat,omitempty"`
	NoHP               string `json:"no_hp,omitempty"`
	Email              string `json:"email,omitempty"`
	PendidikanTerakhir string `json:"pendidikan_terakhir,omitempty"`
	BidangStudi        string `json:"bidang_studi,omitempty"`
	StatusTugas        string `json:"status_tugas,omitempty"`
}

// RombonganBelajar represents study group / class
type RombonganBelajar struct {
	RombelID          string `json:"rombel_id"`
	Nama              string `json:"nama"`
	TingkatPendidikan string `json:"tingkat_pendidikan"`
	Jurusan           string `json:"jurusan,omitempty"`
	SemesterID        string `json:"semester_id"`
	WaliKelasNama     string `json:"wali_kelas_nama,omitempty"`
	WaliKelasNUPTK    string `json:"wali_kelas_nuptk,omitempty"`
	JumlahSiswa       int    `json:"jumlah_siswa"`
}

// TableInfo represents schema table inspection
type TableInfo struct {
	Schema    string `json:"schema"`
	TableName string `json:"table_name"`
	TableType string `json:"table_type"`
}

// ColumnInfo represents schema column inspection
type ColumnInfo struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
	IsNullable string `json:"is_nullable"`
}

// PushPayload is the payload sent to central school cloud server
type PushPayload struct {
	SourceHost  string      `json:"source_host"`
	SchoolNPSN  string      `json:"school_npsn,omitempty"`
	PushType    string      `json:"push_type"` // welfare, rapor, siswa, gtk, all
	RecordCount int         `json:"record_count"`
	PushedAt    time.Time   `json:"pushed_at"`
	Data        interface{} `json:"data"`
}
