# 🕌 Salat CLI

Aplikasi CLI untuk menampilkan jadwal sholat dengan **smart location input** dan **geocoding otomatis** - powered by Go, distributed via NPM.

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
- 🌐 **Cross-platform**: Windows, macOS, Linux

## 📦 Instalasi

```bash
# Install global
npm install -g salat-cli

# Atau jalankan langsung tanpa install
npx salat-cli
```

## 📖 Penggunaan

### 🌍 Setup Lokasi (Fitur Baru!)

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

#### Ubah Lokasi (Baru!)
```bash
# Dengan alamat - auto geocoding
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
- `location` - Lokasi (alamat atau koordinat) **[BARU!]**
- `latitude` - Garis lintang
- `longitude` - Garis bujur  
- `method` - Metode perhitungan (MWL, ISNA, Egypt, Makkah, Karachi, Tehran, Kemenag, JAKIM)
- `geocoding_api` - API geocoding (nominatim, photon) **[BARU!]**

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

## 🏗️ Arsitektur

Package NPM ini adalah wrapper untuk binary Go yang dikompilasi untuk berbagai platform:

```
salat-cli/
├── bin/
│   ├── darwin-arm64/    # macOS Apple Silicon
│   ├── darwin-x64/      # macOS Intel
│   ├── linux-arm64/     # Linux ARM64
│   ├── linux-x64/       # Linux x64
│   ├── win32-arm64/     # Windows ARM64
│   └── win32-x64/       # Windows x64
├── index.js             # NPM wrapper script
└── package.json
```

Binary yang sesuai akan dipilih otomatis berdasarkan platform dan arsitektur sistem Anda.

## 📝 API Geocoding

Aplikasi menggunakan geocoding API eksternal untuk konversi alamat:

- **Forward Geocoding**: Alamat → Koordinat + Nama Lokasi
- **Reverse Geocoding**: Koordinat → Nama Lokasi  
- **One-time Only**: API dipanggil sekali saat setup, hasil disimpan di config
- **Offline-first**: Setelah setup, tidak perlu internet untuk perhitungan sholat

## 🔧 Requirements

- **Node.js**: >= 14.0.0
- **OS**: Windows, macOS, atau Linux
- **Arch**: x64 atau ARM64

## 📂 Konfigurasi

File konfigurasi disimpan di `~/.config/salat/config.yaml`:

```yaml
timezone: Asia/Jakarta
latitude: -6.921846
longitude: 107.607083
location_name: "Kota Bandung, Jawa Barat, Jawa, Indonesia"
method: Kemenag
geocoding_api: nominatim
```

## 🆕 What's New in v1.6.0

- ✨ **Smart Location Input**: Setup lokasi dengan alamat atau koordinat
- 🌍 **Geocoding Support**: Auto-convert alamat ke koordinat
- 📍 **Location Names**: Tampilan nama lokasi yang user-friendly
- ⚡ **Enhanced Config**: Manajemen lokasi yang lebih mudah
- 🔄 **Cached Results**: Performa optimal dengan geocoding sekali saja

## 🤝 Contributing

1. Fork repository: [github.com/herbras/BelGolang](https://github.com/herbras/BelGolang)
2. Buat branch fitur (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

## 📄 License

MIT License - lihat file [LICENSE](https://github.com/herbras/BelGolang/blob/main/LICENSE) untuk detail.

## 🔗 Links

- **Repository**: https://github.com/herbras/BelGolang
- **Issues**: https://github.com/herbras/BelGolang/issues
- **NPM Package**: https://www.npmjs.com/package/salat-cli

---

**Version**: 1.6.0 • **Powered by Go** • **Distributed via NPM**