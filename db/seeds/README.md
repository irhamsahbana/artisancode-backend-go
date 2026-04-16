# CSV Seed Conventions

Untuk mengurangi kebingungan antara kolom tabel database dan kolom bantu seed:

- Kolom tanpa prefix adalah kolom data utama yang mewakili field tabel atau field seed utama.
- Kolom dengan prefix `helper__` adalah helper column.

Contoh:

- `tenant_code` pada proses seed sekarang ditulis sebagai `helper__tenant_code`
- `company_code` ditulis sebagai `helper__company_code`
- `password_plain` ditulis sebagai `helper__password_plain`

Helper column dipakai untuk:

- lookup relasi ke record lain
- backfill ID antar file CSV
- transformasi data sebelum disimpan ke DB

Seeder sekarang mengharuskan helper column memakai prefix `helper__`.
Jika ada helper column tanpa prefix, proses seed akan gagal lebih awal agar format CSV tetap konsisten dan tidak membingungkan.
