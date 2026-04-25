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

## Internal RBAC Template

Template role dan permission untuk provisioning client baru disimpan di:

- `internal_template_roles.csv`
- `internal_template_permissions.csv`
- `internal_template_role_permissions.csv`

File `roles.csv`, `permissions.csv`, dan `role_permissions.csv` hanya untuk data tenant/client yang sudah konkret. Jangan menambahkan baris template dengan `helper__tenant_code` kosong di tiga file tersebut.

Seeder sekarang mengharuskan helper column memakai prefix `helper__`.
Jika ada helper column tanpa prefix, proses seed akan gagal lebih awal agar format CSV tetap konsisten dan tidak membingungkan.
