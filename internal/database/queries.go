package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ardianryan/dapodik-bridge/internal/models"
)

// GetKesejahteraan fetches all student welfare records (PIP, KIP, PKH, KKS, KPS, etc.)
func (m *DBManager) GetKesejahteraan(ctx context.Context, filterJenis string, limit, offset int) ([]models.StudentWelfare, int, error) {
	if m.Pool() == nil {
		return nil, 0, fmt.Errorf("database connection is not available")
	}

	hasKpdTable, _ := m.TableExists(ctx, "kesejahteraan_peserta_didik")

	var query string
	var countQuery string
	var args []interface{}
	argIdx := 1

	if hasKpdTable {
		// Rich join between peserta_didik and kesejahteraan_peserta_didik
		baseWhere := "WHERE (pd.soft_delete = 0 OR pd.soft_delete IS NULL)"
		if filterJenis != "" {
			baseWhere += fmt.Sprintf(" AND (LOWER(COALESCE(jk.nama, '')) LIKE $%d OR LOWER(COALESCE(kpd.nomor_kartu, '')) LIKE $%d)", argIdx, argIdx)
			args = append(args, "%"+strings.ToLower(filterJenis)+"%")
			argIdx++
		}

		countQuery = fmt.Sprintf(`
			SELECT COUNT(DISTINCT pd.peserta_didik_id)
			FROM peserta_didik pd
			LEFT JOIN kesejahteraan_peserta_didik kpd ON pd.peserta_didik_id = kpd.peserta_didik_id AND (kpd.soft_delete = 0 OR kpd.soft_delete IS NULL)
			LEFT JOIN ref.jenis_kesejahteraan jk ON kpd.jenis_kesejahteraan_id = jk.jenis_kesejahteraan_id
			%s
			AND (kpd.kesejahteraan_id IS NOT NULL OR pd.layak_pip = 1 OR pd.penerima_kip = 1 OR pd.penerima_kps = 1 OR pd.no_kks IS NOT NULL)
		`, baseWhere)

		query = fmt.Sprintf(`
			SELECT 
				pd.peserta_didik_id::text,
				COALESCE(pd.nisn, ''),
				COALESCE(pd.nama, ''),
				COALESCE(pd.jenis_kelamin, ''),
				COALESCE(CAST(rb.tingkat_pendidikan_id AS text), ''),
				COALESCE(rb.nama, ''),
				COALESCE(jk.nama, CASE 
					WHEN pd.penerima_kip = 1 THEN 'Program Indonesia Pintar (KIP)'
					WHEN pd.layak_pip = 1 THEN 'Usulan Layak PIP'
					WHEN pd.penerima_kps = 1 THEN 'Program Perlindungan Sosial (KPS)'
					WHEN pd.no_kks IS NOT NULL AND pd.no_kks != '' THEN 'Kartu Keluarga Sejahtera (KKS)'
					ELSE 'Bantuan Kesejahteraan'
				END) AS jenis_bantuan,
				COALESCE(kpd.nomor_kartu, pd.nomor_kip, pd.no_kks, pd.no_kps, ''),
				COALESCE(kpd.nama_di_kartu, pd.nama_kip, pd.nama, ''),
				COALESCE(kpd.tahun_mulai, 0),
				COALESCE(kpd.tahun_selesai, 0),
				CASE WHEN COALESCE(kpd.status, 1) = 1 THEN true ELSE false END,
				CASE WHEN pd.layak_pip = 1 OR pd.penerima_kip = 1 THEN true ELSE false END,
				COALESCE(pd.alasan_layak_pip, ''),
				COALESCE(pd.nama_ibu_kandung, '')
			FROM peserta_didik pd
			LEFT JOIN kesejahteraan_peserta_didik kpd ON pd.peserta_didik_id = kpd.peserta_didik_id AND (kpd.soft_delete = 0 OR kpd.soft_delete IS NULL)
			LEFT JOIN ref.jenis_kesejahteraan jk ON kpd.jenis_kesejahteraan_id = jk.jenis_kesejahteraan_id
			LEFT JOIN anggota_rombel ar ON pd.peserta_didik_id = ar.peserta_didik_id AND (ar.soft_delete = 0 OR ar.soft_delete IS NULL)
			LEFT JOIN rombongan_belajar rb ON ar.rombongan_belajar_id = rb.rombongan_belajar_id AND (rb.soft_delete = 0 OR rb.soft_delete IS NULL)
			%s
			AND (kpd.kesejahteraan_id IS NOT NULL OR pd.layak_pip = 1 OR pd.penerima_kip = 1 OR pd.penerima_kps = 1 OR (pd.no_kks IS NOT NULL AND pd.no_kks != ''))
			ORDER BY pd.nama ASC
			LIMIT $%d OFFSET $%d
		`, baseWhere, argIdx, argIdx+1)
	} else {
		// Fallback directly on peserta_didik columns
		baseWhere := "WHERE (pd.soft_delete = 0 OR pd.soft_delete IS NULL) AND (pd.layak_pip = 1 OR pd.penerima_kip = 1 OR pd.penerima_kps = 1 OR (pd.no_kks IS NOT NULL AND pd.no_kks != ''))"
		if filterJenis != "" {
			baseWhere += fmt.Sprintf(" AND (LOWER(COALESCE(pd.alasan_layak_pip, '')) LIKE $%d OR LOWER(COALESCE(pd.nomor_kip, '')) LIKE $%d)", argIdx, argIdx)
			args = append(args, "%"+strings.ToLower(filterJenis)+"%")
			argIdx++
		}

		countQuery = fmt.Sprintf("SELECT COUNT(*) FROM peserta_didik pd %s", baseWhere)
		query = fmt.Sprintf(`
			SELECT 
				pd.peserta_didik_id::text,
				COALESCE(pd.nisn, ''),
				COALESCE(pd.nama, ''),
				COALESCE(pd.jenis_kelamin, ''),
				COALESCE(CAST(rb.tingkat_pendidikan_id AS text), ''),
				COALESCE(rb.nama, ''),
				CASE 
					WHEN pd.penerima_kip = 1 THEN 'Program Indonesia Pintar (KIP)'
					WHEN pd.layak_pip = 1 THEN 'Usulan Layak PIP'
					WHEN pd.penerima_kps = 1 THEN 'Program Perlindungan Sosial (KPS)'
					WHEN pd.no_kks IS NOT NULL AND pd.no_kks != '' THEN 'Kartu Keluarga Sejahtera (KKS)'
					ELSE 'Bantuan Siswa'
				END AS jenis_bantuan,
				COALESCE(pd.nomor_kip, pd.no_kks, pd.no_kps, ''),
				COALESCE(pd.nama_kip, pd.nama, ''),
				0 AS tahun_mulai,
				0 AS tahun_selesai,
				true AS status_aktif,
				CASE WHEN pd.layak_pip = 1 OR pd.penerima_kip = 1 THEN true ELSE false END,
				COALESCE(pd.alasan_layak_pip, ''),
				COALESCE(pd.nama_ibu_kandung, '')
			FROM peserta_didik pd
			LEFT JOIN anggota_rombel ar ON pd.peserta_didik_id = ar.peserta_didik_id AND (ar.soft_delete = 0 OR ar.soft_delete IS NULL)
			LEFT JOIN rombongan_belajar rb ON ar.rombongan_belajar_id = rb.rombongan_belajar_id AND (rb.soft_delete = 0 OR rb.soft_delete IS NULL)
			%s
			ORDER BY pd.nama ASC
			LIMIT $%d OFFSET $%d
		`, baseWhere, argIdx, argIdx+1)
	}

	var total int
	if err := m.Pool().QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	argsWithPaging := append(args, limit, offset)
	rows, err := m.Pool().Query(ctx, query, argsWithPaging...)
	if err != nil {
		return nil, 0, fmt.Errorf("query kesejahteraan failed: %w", err)
	}
	defer rows.Close()

	var results []models.StudentWelfare
	for rows.Next() {
		var w models.StudentWelfare
		err := rows.Scan(
			&w.PesertaDidikID,
			&w.NISN,
			&w.Nama,
			&w.JenisKelamin,
			&w.TingkatKelas,
			&w.Rombel,
			&w.JenisBantuan,
			&w.NomorKartu,
			&w.NamaDiKartu,
			&w.TahunMulai,
			&w.TahunSelesai,
			&w.StatusAktif,
			&w.LayakPIP,
			&w.AlasanLayakPIP,
			&w.NamaIbuKandung,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, w)
	}

	return results, total, rows.Err()
}

// GetRaporGrades retrieves multi-semester grades for students
func (m *DBManager) GetRaporGrades(ctx context.Context, semesterID, rombelID, nisn string, limit, offset int) ([]models.RaporGrade, int, error) {
	if m.Pool() == nil {
		return nil, 0, fmt.Errorf("database connection is not available")
	}

	// Check table existence for grades: nilai_rapor or nilai_akhir
	hasNilaiRapor, _ := m.TableExists(ctx, "nilai_rapor")
	tableName := "nilai_rapor"
	if !hasNilaiRapor {
		hasNilaiAkhir, _ := m.TableExists(ctx, "nilai_akhir")
		if hasNilaiAkhir {
			tableName = "nilai_akhir"
		} else {
			return nil, 0, fmt.Errorf("tabel nilai rapor tidak ditemukan di skema Dapodik lokal")
		}
	}

	whereClauses := []string{"(nr.soft_delete = 0 OR nr.soft_delete IS NULL)"}
	var args []interface{}
	argIdx := 1

	if semesterID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("rb.semester_id = $%d", argIdx))
		args = append(args, semesterID)
		argIdx++
	}
	if rombelID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("rb.rombongan_belajar_id = $%d", argIdx))
		args = append(args, rombelID)
		argIdx++
	}
	if nisn != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("pd.nisn = $%d", argIdx))
		args = append(args, nisn)
		argIdx++
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s nr
		JOIN anggota_rombel ar ON nr.anggota_rombel_id = ar.anggota_rombel_id
		JOIN peserta_didik pd ON ar.peserta_didik_id = pd.peserta_didik_id
		JOIN rombongan_belajar rb ON ar.rombongan_belajar_id = rb.rombongan_belajar_id
		%s
	`, tableName, whereSQL)

	var total int
	if err := m.Pool().QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count rapor failed: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT 
			pd.peserta_didik_id::text,
			COALESCE(pd.nisn, ''),
			COALESCE(pd.nama, ''),
			COALESCE(rb.nama, ''),
			COALESCE(rb.semester_id, ''),
			COALESCE(mp.nama, 'Mata Pelajaran'),
			COALESCE(nr.nilai_angka, nr.nilai_kognitif, 0),
			COALESCE(nr.nilai_huruf, ''),
			COALESCE(nr.predikat, ''),
			nr.kkm
		FROM %s nr
		JOIN anggota_rombel ar ON nr.anggota_rombel_id = ar.anggota_rombel_id
		JOIN peserta_didik pd ON ar.peserta_didik_id = pd.peserta_didik_id
		JOIN rombongan_belajar rb ON ar.rombongan_belajar_id = rb.rombongan_belajar_id
		LEFT JOIN pembelajaran p ON nr.pembelajaran_id = p.pembelajaran_id
		LEFT JOIN ref.mata_pelajaran mp ON p.mata_pelajaran_id = mp.mata_pelajaran_id
		%s
		ORDER BY pd.nama ASC, mp.nama ASC
		LIMIT $%d OFFSET $%d
	`, tableName, whereSQL, argIdx, argIdx+1)

	argsWithPaging := append(args, limit, offset)
	rows, err := m.Pool().Query(ctx, query, argsWithPaging...)
	if err != nil {
		return nil, 0, fmt.Errorf("query rapor failed: %w", err)
	}
	defer rows.Close()

	var results []models.RaporGrade
	for rows.Next() {
		var g models.RaporGrade
		var kkm sql.NullFloat64
		err := rows.Scan(
			&g.PesertaDidikID,
			&g.NISN,
			&g.NamaSiswa,
			&g.Rombel,
			&g.SemesterID,
			&g.MataPelajaran,
			&g.NilaiAngka,
			&g.NilaiHuruf,
			&g.Predikat,
			&kkm,
		)
		if err != nil {
			return nil, 0, err
		}
		if kkm.Valid {
			g.KKM = &kkm.Float64
		}
		results = append(results, g)
	}

	return results, total, rows.Err()
}

// GetSiswaKomprehensif fetches comprehensive student data with periodic and family fields
func (m *DBManager) GetSiswaKomprehensif(ctx context.Context, q, rombelID string, limit, offset int) ([]models.StudentComprehensive, int, error) {
	if m.Pool() == nil {
		return nil, 0, fmt.Errorf("database connection is not available")
	}

	whereClauses := []string{"(pd.soft_delete = 0 OR pd.soft_delete IS NULL)"}
	var args []interface{}
	argIdx := 1

	if q != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(LOWER(pd.nama) LIKE $%d OR pd.nisn LIKE $%d OR pd.nik LIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+strings.ToLower(q)+"%")
		argIdx++
	}
	if rombelID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("rb.rombongan_belajar_id = $%d", argIdx))
		args = append(args, rombelID)
		argIdx++
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT pd.peserta_didik_id)
		FROM peserta_didik pd
		LEFT JOIN anggota_rombel ar ON pd.peserta_didik_id = ar.peserta_didik_id AND (ar.soft_delete = 0 OR ar.soft_delete IS NULL)
		LEFT JOIN rombongan_belajar rb ON ar.rombongan_belajar_id = rb.rombongan_belajar_id AND (rb.soft_delete = 0 OR rb.soft_delete IS NULL)
		%s
	`, whereSQL)

	var total int
	if err := m.Pool().QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count siswa failed: %w", err)
	}

	hasLongitudinal, _ := m.TableExists(ctx, "peserta_didik_longitudinal")
	longitudinalJoin := ""
	tbField := "COALESCE(pdl.tinggi_badan, 0)"
	bbField := "COALESCE(pdl.berat_badan, 0)"
	jarakField := "COALESCE(pdl.jarak_rumah_ke_sekolah_km, 0)"
	waktuField := "COALESCE(pdl.waktu_tempuh_ke_sekolah, 0)"

	if hasLongitudinal {
		longitudinalJoin = "LEFT JOIN peserta_didik_longitudinal pdl ON pd.peserta_didik_id = pdl.peserta_didik_id AND (pdl.soft_delete = 0 OR pdl.soft_delete IS NULL)"
	} else {
		tbField = "0"
		bbField = "0"
		jarakField = "0"
		waktuField = "0"
	}

	query := fmt.Sprintf(`
		SELECT 
			pd.peserta_didik_id::text,
			COALESCE(pd.nama, ''),
			COALESCE(pd.nisn, ''),
			COALESCE(pd.nik, ''),
			COALESCE(pd.jenis_kelamin, ''),
			COALESCE(pd.tempat_lahir, ''),
			COALESCE(pd.tanggal_lahir::text, ''),
			COALESCE(ref_a.nama, ''),
			COALESCE(pd.alamat_jalan, ''),
			COALESCE(pd.rt, ''),
			COALESCE(pd.rw, ''),
			COALESCE(pd.desa_kelurahan, ''),
			COALESCE(pd.kecamatan, ''),
			COALESCE(pd.kode_pos, ''),
			pd.lintang,
			pd.bujur,
			COALESCE(pd.nama_ayah, ''),
			COALESCE(ref_pa.nama, ''),
			COALESCE(ref_ha.nama, ''),
			COALESCE(pd.nama_ibu_kandung, ''),
			COALESCE(ref_pi.nama, ''),
			COALESCE(ref_hi.nama, ''),
			%s AS tinggi_badan,
			%s AS berat_badan,
			%s AS jarak_ke_sekolah,
			%s AS waktu_tempuh,
			COALESCE(rb.nama, ''),
			COALESCE(CAST(rb.tingkat_pendidikan_id AS text), '')
		FROM peserta_didik pd
		LEFT JOIN anggota_rombel ar ON pd.peserta_didik_id = ar.peserta_didik_id AND (ar.soft_delete = 0 OR ar.soft_delete IS NULL)
		LEFT JOIN rombongan_belajar rb ON ar.rombongan_belajar_id = rb.rombongan_belajar_id AND (rb.soft_delete = 0 OR rb.soft_delete IS NULL)
		LEFT JOIN ref.agama ref_a ON pd.agama_id = ref_a.agama_id
		LEFT JOIN ref.pekerjaan ref_pa ON pd.pekerjaan_id_ayah = ref_pa.pekerjaan_id
		LEFT JOIN ref.penghasilan ref_ha ON pd.penghasilan_id_ayah = ref_ha.penghasilan_id
		LEFT JOIN ref.pekerjaan ref_pi ON pd.pekerjaan_id_ibu = ref_pi.pekerjaan_id
		LEFT JOIN ref.penghasilan ref_hi ON pd.penghasilan_id_ibu = ref_hi.penghasilan_id
		%s
		%s
		ORDER BY pd.nama ASC
		LIMIT $%d OFFSET $%d
	`, tbField, bbField, jarakField, waktuField, longitudinalJoin, whereSQL, argIdx, argIdx+1)

	argsWithPaging := append(args, limit, offset)
	rows, err := m.Pool().Query(ctx, query, argsWithPaging...)
	if err != nil {
		return nil, 0, fmt.Errorf("query siswa failed: %w", err)
	}
	defer rows.Close()

	var results []models.StudentComprehensive
	for rows.Next() {
		var s models.StudentComprehensive
		var lintang, bujur sql.NullFloat64
		err := rows.Scan(
			&s.PesertaDidikID,
			&s.Nama,
			&s.NISN,
			&s.NIK,
			&s.JenisKelamin,
			&s.TempatLahir,
			&s.TanggalLahir,
			&s.Agama,
			&s.AlamatJalan,
			&s.RT,
			&s.RW,
			&s.KelurahanDesa,
			&s.Kecamatan,
			&s.KodePos,
			&lintang,
			&bujur,
			&s.NamaAyah,
			&s.PekerjaanAyah,
			&s.PenghasilanAyah,
			&s.NamaIbuKandung,
			&s.PekerjaanIbu,
			&s.PenghasilanIbu,
			&s.TinggiBadanCm,
			&s.BeratBadanKg,
			&s.JarakKeSekolahKm,
			&s.WaktuTempuhMenit,
			&s.RombelSaatIni,
			&s.TingkatKelas,
		)
		if err != nil {
			return nil, 0, err
		}
		if lintang.Valid {
			s.Lintang = &lintang.Float64
		}
		if bujur.Valid {
			s.Bujur = &bujur.Float64
		}
		results = append(results, s)
	}

	return results, total, rows.Err()
}

// GetGTKLengkap fetches complete GTK (Teachers & Staff) details
func (m *DBManager) GetGTKLengkap(ctx context.Context, q string, limit, offset int) ([]models.GTKLengkap, int, error) {
	if m.Pool() == nil {
		return nil, 0, fmt.Errorf("database connection is not available")
	}

	whereClauses := []string{"(p.soft_delete = 0 OR p.soft_delete IS NULL)"}
	var args []interface{}
	argIdx := 1

	if q != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(LOWER(p.nama) LIKE $%d OR p.nuptk LIKE $%d OR p.nip LIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+strings.ToLower(q)+"%")
		argIdx++
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ptk p %s", whereSQL)
	var total int
	if err := m.Pool().QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count gtk failed: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT 
			p.ptk_id::text,
			COALESCE(p.nama, ''),
			COALESCE(p.gelar_depan, ''),
			COALESCE(p.gelar_belakang, ''),
			COALESCE(p.nik, ''),
			COALESCE(p.jenis_kelamin, ''),
			COALESCE(p.tempat_lahir, ''),
			COALESCE(p.tanggal_lahir::text, ''),
			COALESCE(p.nip, ''),
			COALESCE(p.nuptk, ''),
			COALESCE(sk.nama, ''),
			COALESCE(jp.nama, ''),
			COALESCE(ref_a.nama, ''),
			COALESCE(p.alamat_jalan, ''),
			COALESCE(p.no_hp, ''),
			COALESCE(p.email, ''),
			COALESCE(p.status_tugas, 'Induk')
		FROM ptk p
		LEFT JOIN ref.status_kepegawaian sk ON p.status_kepegawaian_id = sk.status_kepegawaian_id
		LEFT JOIN ref.jenis_ptk jp ON p.jenis_ptk_id = jp.jenis_ptk_id
		LEFT JOIN ref.agama ref_a ON p.agama_id = ref_a.agama_id
		%s
		ORDER BY p.nama ASC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argIdx, argIdx+1)

	argsWithPaging := append(args, limit, offset)
	rows, err := m.Pool().Query(ctx, query, argsWithPaging...)
	if err != nil {
		return nil, 0, fmt.Errorf("query gtk failed: %w", err)
	}
	defer rows.Close()

	var results []models.GTKLengkap
	for rows.Next() {
		var g models.GTKLengkap
		err := rows.Scan(
			&g.PTKID,
			&g.Nama,
			&g.GelarDepan,
			&g.GelarBelakang,
			&g.NIK,
			&g.JenisKelamin,
			&g.TempatLahir,
			&g.TanggalLahir,
			&g.NIP,
			&g.NUPTK,
			&g.StatusKepegawaian,
			&g.JenisPTK,
			&g.Agama,
			&g.Alamat,
			&g.NoHP,
			&g.Email,
			&g.StatusTugas,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, g)
	}

	return results, total, rows.Err()
}

// GetRombel fetches study groups with member counts
func (m *DBManager) GetRombel(ctx context.Context, semesterID string, limit, offset int) ([]models.RombonganBelajar, int, error) {
	if m.Pool() == nil {
		return nil, 0, fmt.Errorf("database connection is not available")
	}

	whereClauses := []string{"(rb.soft_delete = 0 OR rb.soft_delete IS NULL)"}
	var args []interface{}
	argIdx := 1

	if semesterID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("rb.semester_id = $%d", argIdx))
		args = append(args, semesterID)
		argIdx++
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM rombongan_belajar rb %s", whereSQL)
	var total int
	if err := m.Pool().QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count rombel failed: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT 
			rb.rombongan_belajar_id::text,
			COALESCE(rb.nama, ''),
			COALESCE(CAST(rb.tingkat_pendidikan_id AS text), ''),
			COALESCE(j.nama, ''),
			COALESCE(rb.semester_id, ''),
			COALESCE(p.nama, ''),
			COALESCE(p.nuptk, ''),
			COUNT(ar.anggota_rombel_id) AS jumlah_siswa
		FROM rombongan_belajar rb
		LEFT JOIN jurusan_sp j ON rb.jurusan_sp_id = j.jurusan_sp_id
		LEFT JOIN ptk p ON rb.ptk_id = p.ptk_id
		LEFT JOIN anggota_rombel ar ON rb.rombongan_belajar_id = ar.rombongan_belajar_id AND (ar.soft_delete = 0 OR ar.soft_delete IS NULL)
		%s
		GROUP BY rb.rombongan_belajar_id, rb.nama, rb.tingkat_pendidikan_id, j.nama, rb.semester_id, p.nama, p.nuptk
		ORDER BY rb.tingkat_pendidikan_id ASC, rb.nama ASC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argIdx, argIdx+1)

	argsWithPaging := append(args, limit, offset)
	rows, err := m.Pool().Query(ctx, query, argsWithPaging...)
	if err != nil {
		return nil, 0, fmt.Errorf("query rombel failed: %w", err)
	}
	defer rows.Close()

	var results []models.RombonganBelajar
	for rows.Next() {
		var r models.RombonganBelajar
		err := rows.Scan(
			&r.RombelID,
			&r.Nama,
			&r.TingkatPendidikan,
			&r.Jurusan,
			&r.SemesterID,
			&r.WaliKelasNama,
			&r.WaliKelasNUPTK,
			&r.JumlahSiswa,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, r)
	}

	return results, total, rows.Err()
}
