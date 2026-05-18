# ISSUE: Perencanaan Implementasi Otorisasi Akses Route (Role-Permission) di Frontend

## Deskripsi Fitur
Fitur ini bertujuan untuk mengamankan route/halaman pada aplikasi frontend agar hanya dapat diakses oleh user yang memiliki hak akses (permission) yang sesuai. Verifikasi hak akses dilakukan sebelum halaman dirender (menggunakan *Navigation Guard* / *Route Middleware*).

---

## Tahapan Implementasi

### Langkah 1: Persiapkan Data Otorisasi dari Backend (API Support & Middleware)
Agar frontend tahu apa saja hak akses dari user yang sedang login, backend harus menyediakan data permissions dan memproteksi endpoint API menggunakan Middleware Otorisasi.

**Tugas Backend:**
1. Update `UserResponse` DTO pada backend agar menyertakan array list permissions (berupa string nama/path permission).
   - Contoh output JSON yang diharapkan pada `/api/auth/me`:
     ```json
     {
       "code": 200,
       "status": "OK",
       "message": "User retrieved successfully",
       "data": {
         "id": "user-uuid",
         "name": "John Doe",
         "role": "staff",
         "permissions": [
           "products:view",
           "products:create",
           "purchase-orders:view"
         ]
       }
     }
     ```
2. Pastikan kueri database untuk memuat permissions user dilakukan menggunakan **JOIN SQL** demi efisiensi dan performa (menghindari multipel query `Preload` dan looping manual di sisi memori aplikasi).
   - *Rekomendasi Query GORM (Single Query Count):*
     ```go
     var count int64
     db.Table("users").
         Joins("JOIN roles ON users.role = roles.name").
         Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
         Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
         Where("users.id = ? AND permissions.name = ?", userID, requiredPermission).
         Where("users.deleted_at IS NULL AND roles.deleted_at IS NULL AND permissions.deleted_at IS NULL").
         Count(&count)
     ```
3. Implementasikan fungsi pembantu `check` pada route setup (misalnya `route.go`) agar penulisan route guard di backend tetap pendek dan bersih:
   ```go
   check := func(permission string) fiber.Handler {
       return middleware.RequirePermission(reg.DB, permission)
   }
   r.Post("", check("products:create"), h.Create)
   ```

---

### Langkah 2: Menyimpan State Hak Akses di Frontend
Simpan data user beserta array `permissions` ke dalam Global State Management frontend (misalnya **Pinia** jika menggunakan Nuxt/Vue, atau **Redux/Context** jika React).

**Tugas:**
1. Setelah login sukses, atau saat memuat ulang halaman pertama kali (aplikasi diinisialisasi), panggil `/api/auth/me`.
2. Simpan array `permissions` (contoh: `['products:view', 'products:create']`) ke dalam global state/store Auth.

---

### Langkah 3: Menentukan Kebutuhan Akses di Setiap Halaman (Meta Route)
Setiap halaman/route yang sensitif harus ditandai dengan permission apa saja yang dibutuhkan untuk membukanya.

**Tugas:**
1. Gunakan properti `meta` pada konfigurasi route/page di frontend.
2. Contoh implementasi (jika menggunakan Nuxt 3 dengan `definePageMeta`):
   ```typescript
   // Halaman: pages/products/create.vue
   definePageMeta({
     middleware: ['auth'],
     meta: {
       requiredPermission: 'products:create'
     }
   })
   ```
3. Jika menggunakan Vue Router biasa, definisikan di array routes:
   ```javascript
   {
     path: '/products/create',
     component: ProductCreateComponent,
     meta: {
       requiredPermission: 'products:create'
     }
   }
   ```

---

### Langkah 4: Membuat Route Guard / Middleware Proteksi Halaman
Buat middleware routing global untuk mencegat setiap navigasi user dan mencocokkan hak aksesnya sebelum masuk ke halaman tujuan.

**Tugas:**
1. Buat file middleware global baru (misalnya `auth-guard.global.ts` di Nuxt atau gunakan `router.beforeEach` di Vue Router).
2. Tulis logika pengecekan di dalam middleware:
   - Ambil target halaman yang dituju (`to`).
   - Cek apakah target halaman memerlukan permission khusus (`to.meta.requiredPermission`).
   - Jika **tidak memerlukan**, biarkan navigasi berlanjut.
   - Jika **memerlukan**:
     1. Ambil daftar `permissions` milik user dari state store auth.
     2. Cek apakah daftar tersebut mengandung permission yang dibutuhkan (`to.meta.requiredPermission`).
     3. Jika **punya akses**, biarkan navigasi berlanjut.
     4. Jika **tidak punya akses**:
        - Batalkan navigasi ke halaman tersebut.
        - Tampilkan notifikasi peringatan (*Toast Alert*) seperti: *"Anda tidak memiliki hak akses untuk membuka halaman ini!"*.
        - Redirect user ke halaman dashboard utama atau halaman error 403 (Forbidden).

---

### Langkah 5: Helper Tambahan untuk Elemen UI (Opsional tapi Sangat Direkomendasikan)
Selain mengamankan link URL/route, kita juga harus menyembunyikan tombol-tombol aksi (seperti tombol "Tambah", "Edit", atau "Hapus") dari user yang tidak memiliki wewenang.

**Tugas:**
1. Buat helper function global, misalnya `hasPermission(permissionName: string): boolean`.
   - Fungsi ini mengambil daftar `permissions` dari store, lalu return `true` jika user memilikinya, atau `false` jika tidak.
2. Gunakan helper ini di template UI menggunakan `v-if` (Vue/Nuxt):
   ```html
   <!-- Tombol Tambah Produk hanya muncul jika user punya permission products:create -->
   <button v-if="hasPermission('products:create')" @click="addProduct">
     Tambah Produk
   </button>
   ```

---

## Rencana Pengujian (Testing)
1. **Uji Coba User Admin:** Login sebagai Admin yang memiliki semua permission. Pastikan semua menu, tombol, dan halaman bisa diakses tanpa kendala.
2. **Uji Coba User Staff Biasa:** Login sebagai Staff yang hanya memiliki permission terbatas (misal: hanya `products:view`).
   - Pastikan Staff bisa masuk ke halaman daftar produk.
   - Coba akses URL `/products/create` secara langsung lewat address bar browser. Sistem harus otomatis memblokir dan mengarahkan kembali ke dashboard dengan notifikasi error 403/Forbidden.
   - Pastikan tombol "Tambah Produk" tidak kelihatan di halaman daftar produk.
