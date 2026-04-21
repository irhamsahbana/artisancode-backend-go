# User Invitation Flow

## Goal

Menyediakan flow undangan yang jelas supaya produk tidak bergantung pada menu generik `Users`, tetapi tetap punya lifecycle akses yang rapi untuk:

- owner
- admin
- employee

Dokumen ini adalah planning awal produk + backend. Implementasi FE dan BE bisa mengikuti bertahap.

## Product Direction

Prinsip utama:

- `users` adalah identity/auth layer, bukan menu utama produk.
- `employees` adalah domain object untuk data tenaga kerja.
- akses diberikan lewat flow sesuai aktor bisnis, bukan lewat CRUD user generik.

Implikasi UX:

- owner masuk lewat login/register owner.
- admin diundang atau dibuat oleh owner.
- employee dibuat lewat modul employee, lalu akses login diaktifkan saat perlu.

## Target Flows

### 1. Owner

- Owner mendaftar tenant/company saat register awal.
- Setelah login, owner bisa mengelola access untuk admin.
- Owner tidak dibuat dari invitation flow biasa kecuali nanti ada kebutuhan multi-owner.

### 2. Admin

- Owner membuka halaman `Akses` atau `Admin`.
- Owner mengirim undangan ke email admin.
- Admin menerima link, membuat password, lalu status invitation berubah menjadi accepted.

### 3. Employee

- Owner/admin membuat employee record lebih dulu.
- Employee bisa tetap aktif sebagai data HR walau belum punya akun login.
- Saat akses dibutuhkan, owner/admin mengirim undangan dari detail/list employee.
- Employee menerima link aktivasi, membuat password, lalu employee terhubung ke akun login.

## Proposed Domain Model

Model minimum yang disarankan:

- `users`
  - identity login
  - email, password hash, status verifikasi, tenant scope
- `employees`
  - domain profile pegawai
  - optional `user_id` sampai employee mengaktifkan akun
- `user_roles` atau membership setara
  - relasi role seperti owner, admin, employee
- `user_invitations`
  - source of truth undangan

## Proposed `user_invitations` Fields

Kolom awal yang disarankan:

- `id`
- `tenant_id`
- `employee_id` nullable
- `email`
- `role_code`
- `status`
- `token_hash`
- `expires_at`
- `accepted_at` nullable
- `revoked_at` nullable
- `invited_by`
- `last_sent_at`
- `created_at`
- `updated_at`

Nilai `status` minimum:

- `pending`
- `accepted`
- `expired`
- `revoked`

Catatan:

- Simpan hash token, bukan raw token.
- Email invitation harus immutable per record agar audit jelas.
- Untuk employee invitation, `employee_id` mengikat invitation ke domain record yang benar.

## Backend API Plan

Tahap awal cukup fokus ke API berikut:

### Protected APIs

- `POST /user-invitations`
  - membuat invitation baru untuk admin atau employee
- `GET /user-invitations`
  - list/filter invitation
- `POST /user-invitations/:id/resend`
  - kirim ulang invitation yang masih valid atau regenerate token
- `POST /user-invitations/:id/revoke`
  - batalkan invitation

### Public APIs

- `GET /user-invitations/accept`
  - validasi token dan tampilkan summary invitation
- `POST /user-invitations/accept`
  - set password, verifikasi email, bind role, bind employee bila ada

## Acceptance Flow

Saat invitation diterima:

1. validasi token, status, tenant, dan expiry
2. cek apakah email invitation masih sesuai
3. buat user baru bila belum ada
4. bila user sudah ada dalam tenant yang sama, reuse user itu dengan aturan ketat
5. assign role yang diundang
6. link ke `employees.user_id` bila invitation employee
7. set `email_verified_at`
8. tandai invitation accepted

## Guardrails

Aturan yang disarankan sejak awal:

- Tidak boleh ada dua invitation aktif untuk email + role + tenant yang sama.
- Employee yang sudah terhubung ke user aktif tidak boleh diundang ulang tanpa flow reset/reinvite yang eksplisit.
- Acceptance harus idempotent terhadap token yang sama setelah sukses.
- Admin tidak boleh mengundang role di atas dirinya.
- Owner-only action untuk undangan admin level tinggi.

## FE Plan

Tahap FE yang disarankan:

### Phase 1

- sembunyikan menu `Users` dari sidebar/dashboard
- tetap pertahankan route internal bila masih dipakai tim
- tambahkan titik masuk invitation di:
  - halaman employee
  - halaman roles/access untuk admin

### Phase 2

- buat halaman `Invitations` atau panel invitation di `Roles`
- tampilkan status:
  - pending
  - accepted
  - expired
  - revoked
- sediakan action:
  - invite
  - resend
  - revoke

### Phase 3

- buat public accept page
- form set password
- tampilkan identitas tenant, role, dan employee summary

## Delivery Plan

Urutan implementasi yang paling aman:

1. Tambah tabel `user_invitations` + repository + core dasar.
2. Tambah protected API create/list/resend/revoke.
3. Tambah email template invitation.
4. Tambah public accept API.
5. Tambah FE invitation entry point dari employee dan access management.
6. Tambah FE accept invitation page.
7. Matikan dependensi bisnis pada menu `Users`.

## Open Decisions

Beberapa keputusan produk yang masih perlu disepakati:

- apakah admin boleh membuat invitation admin lain
- apakah invitation employee otomatis mengaktifkan app access atau perlu approval tambahan
- apakah satu email boleh terhubung ke lebih dari satu employee di tenant berbeda
- apakah owner tambahan nanti pakai invitation flow yang sama atau flow khusus

## Recommended MVP

MVP paling efektif:

- invitation khusus `admin` dan `employee`
- create, resend, revoke, accept
- employee record wajib sudah ada sebelum invite employee
- owner register tetap lewat flow existing

Dengan pendekatan ini, menu `Users` bisa tetap tersembunyi tanpa menghilangkan kontrol akses yang dibutuhkan operasional.
