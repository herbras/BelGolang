# Panduan Kontribusi

Terima kasih telah mempertimbangkan untuk berkontribusi pada proyek Salat CLI! Berikut adalah panduan untuk membantu Anda berkontribusi.

## Proses Kontribusi

1. Fork repositori ini
2. Clone fork Anda: `git clone https://github.com/USERNAME/salat.git`
3. Buat branch fitur: `git checkout -b fitur-baru`
4. Lakukan perubahan Anda
5. Commit perubahan: `git commit -am 'Menambahkan fitur baru'`
6. Push ke branch: `git push origin fitur-baru`
7. Kirim Pull Request

## Pengembangan Lokal

### Prasyarat

- Go 1.18 atau lebih baru
- Git

### Setup

```bash
# Clone repositori
git clone https://github.com/herbras/BelGolang/salat.git
cd salat

# Install dependensi
go mod download

# Build aplikasi
go build -o salat
```

### Testing

```bash
go test ./...
```

## Struktur Proyek

```
salat/
├── cmd/           # Command handlers (Cobra)
├── config/        # Konfigurasi aplikasi
├── salat/         # Core logic perhitungan waktu sholat
├── internal/      # Kode internal yang tidak di-export
├── main.go        # Entry point aplikasi
└── go.mod         # Dependensi Go
```

## Panduan Gaya Kode

- Ikuti [Effective Go](https://golang.org/doc/effective_go) dan [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Gunakan `gofmt` untuk memformat kode Anda
- Tambahkan komentar untuk fungsi dan struktur data yang kompleks
- Tulis unit test untuk kode baru

## Menambahkan Fitur Baru

1. Diskusikan fitur baru di issue tracker sebelum memulai pekerjaan besar
2. Tulis unit test untuk fitur baru
3. Update dokumentasi jika diperlukan
4. Pastikan semua test lulus

## Pelaporan Bug

Ketika melaporkan bug, sertakan:

- Versi Salat CLI yang Anda gunakan
- Sistem operasi dan versinya
- Langkah-langkah untuk mereproduksi bug
- Output yang diharapkan dan yang sebenarnya
- Screenshot jika memungkinkan

## Lisensi

Dengan berkontribusi pada proyek ini, Anda setuju bahwa kontribusi Anda akan dilisensikan di bawah lisensi MIT proyek.