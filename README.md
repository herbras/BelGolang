# Salat CLI

Sebuah aplikasi command-line untuk menampilkan jadwal waktu sholat dengan berbagai fitur interaktif dan **geocoding otomatis**.

## ✨ Fitur Utama

- 🌍 **Smart Location Input**: Setup dengan alamat atau koordinat
- 🔍 **Geocoding Otomatis**: Konversi alamat ke koordinat secara otomatis  
- 📍 **Location Names**: Tampilan nama lokasi yang mudah dibaca
- ⚡ **Setup Cepat**: Satu perintah untuk konfigurasi lokasi
- 🕰️ **Countdown Real-time**: Countdown ke waktu sholat berikutnya
- 📺 **Live Update**: Mode watch dengan update otomatis setiap menit
- 🔔 **Notifikasi**: Opsi notifikasi saat masuk waktu sholat
- 🎨 **UI Cantik**: Interface terminal berwarna dengan emoji
- ⚙️ **Konfigurasi Fleksibel**: 8 metode perhitungan (MWL, ISNA, Egypt, Makkah, Karachi, Tehran, Kemenag, JAKIM)

## 🚀 Instalasi

1. Clone repository ini:
```bash
git clone https://github.com/herbras/BelGolang.git
cd BelGolang
```

2. Build aplikasi:
```bash
go build -o salat
```

3. (Opsional) Install ke sistem:
```bash
go install
```

## 📖 Penggunaan

### 🌍 Setup Lokasi (Baru!)

#### Setup dengan Alamat
```bash
# Setup dengan nama kota/alamat
salat setup "Jakarta, Indonesia"
salat setup "Bandung, Indonesia"
salat setup "Monas Jakarta"

# Pilih API geocoding
salat setup "Surabaya" --api photon
salat setup "Medan" --api nominatim  # default
```

#### Setup dengan Koordinat
```bash
# Format: "latitude,longitude"
salat setup -- "-6.2,106.8"        # Jakarta
salat setup -- "-6.9,107.6"        # Bandung
salat setup -- "-7.25,112.75"      # Surabaya
```

#### Setup Interaktif (Klasik)
```bash
salat setup
```

### 📅 Perintah Utama

#### Tampilkan Jadwal Hari Ini
```bash
salat show
```

Contoh output:
```
🕌 Jadwal Sholat - Friday, 30 May 2025
📍 Kota Bandung, Jawa Barat, Jawa, Indonesia (-6.921846, 107.607083) • Kemenag

WAKTU           JAM      STATUS
-------------------------------
🌙  Imsak        04:27    ✓
🌅  Subuh        04:37    ✓
☀️  Dzuhur      11:52    ✓
🌤️  Ashar       15:14    
🌇  Maghrib      17:45
✨  Isya         18:59
-------------------------------

⏰ Sholat berikutnya: 🌤️ Ashar dalam 1h 22m
```

#### Countdown Sholat Berikutnya
```bash
salat next    # atau salat n
```

#### Status Sholat Saat Ini
```bash
salat now
```

#### Mode Live Update
```bash
salat watch   # atau salat w

# Dengan notifikasi
salat watch --notify
```

### ⚙️ Konfigurasi

#### Lihat Konfigurasi
```bash
salat config show
```

#### Ubah Lokasi
```bash
# Dengan alamat
salat config set location "Yogyakarta, Indonesia"

# Dengan koordinat
salat config set location -- "-7.8,110.4"
```

#### Ubah Pengaturan Lain
```bash
salat config set timezone Asia/Jakarta
salat config set method Kemenag
salat config set geocoding_api photon
```

#### Kunci Konfigurasi yang Tersedia
- `timezone` - Zona waktu (contoh: Asia/Jakarta)
- `location` - Lokasi (alamat atau koordinat)
- `latitude` - Garis lintang
- `longitude` - Garis bujur  
- `method` - Metode perhitungan (MWL, ISNA, Egypt, Makkah, Karachi, Tehran, Kemenag, JAKIM)
- `geocoding_api` - API geocoding (nominatim, photon)

## 🌐 Geocoding APIs

Aplikasi ini mendukung dua penyedia geocoding:

### Nominatim (Default)
- API gratis dari OpenStreetMap
- Rate limit: 1 request/second
- Coverage: Global

### Photon  
- API gratis dari Komoot
- Lebih cepat dari Nominatim
- Coverage: Global
- Usage: `--api photon`

## 🎨 Opsi Tampilan

### Mode Compact
```bash
salat show --compact
```

### Tema
```bash
salat show --theme dark
salat show --theme light  # default
```

## 🔧 Development

### Build untuk Development
```bash
go run main.go setup "Jakarta"
go run main.go show
```

### Build untuk Production
```bash
go build -ldflags="-s -w" -o salat
```

### Build untuk Semua Platform
```bash
goreleaser build --snapshot --clean
```

## 📝 API Geocoding

Aplikasi menggunakan geocoding API eksternal untuk konversi alamat:

- **Forward Geocoding**: Alamat → Koordinat + Nama Lokasi
- **Reverse Geocoding**: Koordinat → Nama Lokasi  
- **One-time Only**: API dipanggil sekali saat setup, hasil disimpan di config
- **Offline-first**: Setelah setup, tidak perlu internet untuk perhitungan sholat

## 🤝 Contributing

Lihat [CONTRIBUTING.md](CONTRIBUTING.md) untuk panduan kontribusi.

## 📄 License

MIT License - lihat file [LICENSE](LICENSE) untuk detail.

## 🔗 Links

- **Repository**: https://github.com/herbras/BelGolang
- **Issues**: https://github.com/herbras/BelGolang/issues
- **NPM Package**: [salat-cli](https://www.npmjs.com/package/salat-cli)

---

**Version**: 1.6.0 • **Last Updated**: January 2025