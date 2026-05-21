# Panduan Integrasi Otorisasi Hak Akses (Role-Permission) pada Frontend (FE)

Dokumen ini memandu tim Frontend (FE) untuk menerapkan sistem pengamanan halaman (*Route Guard*) dan elemen UI berbasis *Permission* setelah backend selesai memperbarui respons endpoint `/api/auth/me`.

---

## 1. Struktur Respons API Terbaru (`GET /api/auth/me`)

Endpoint `/api/auth/me` sekarang mengembalikan daftar string **`permissions`** yang dimiliki oleh pengguna yang sedang login.

### Contoh Response JSON:
```json
{
  "code": 200,
  "status": "OK",
  "message": "User retrieved successfully",
  "data": {
    "id": "018f89ac-3b2d-76af-8239-16a75fcd6ab2",
    "name": "Jane Doe",
    "username": "janedoe",
    "email": "jane@example.com",
    "role": "staff",
    "permissions": [
      "products:view",
      "products:create",
      "purchase-orders:view"
    ]
  }
}
```

> [!NOTE]
> Untuk pengguna dengan role **`programmer`** atau **`administrator`**, array `permissions` akan secara otomatis berisi **seluruh** daftar permission yang ada di sistem (di-bypass penuh dari sisi database backend).

---

## 2. Implementasi State Management (Pinia - Nuxt 3 / Vue 3)

Simpan daftar permissions ke dalam Global Store (misalnya **Pinia**) saat inisialisasi aplikasi atau setelah login berhasil.

```typescript
// stores/auth.ts
import { defineStore } from 'pinia'

interface User {
  id: string
  name: string
  username: string
  email: string
  role: string
  permissions: string[]
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: null as string | null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.user,
    // Getter untuk mempermudah pengecekan permission di luar store
    userPermissions: (state) => state.user?.permissions || [],
  },
  actions: {
    async fetchUser() {
      try {
        const { data } = await useApiFetch('/api/auth/me')
        this.user = data
      } catch (error) {
        this.logout()
      }
    },
    logout() {
      this.user = null
      this.token = null
      // Hapus token dari cookie/localStorage
    }
  }
})
```

---

## 3. Membuat Helper Fungsi Global `hasPermission`

Buat fungsi utilitas global untuk memeriksa kepemilikan hak akses agar pengecekan di komponen UI menjadi sangat ringkas.

```typescript
// composables/usePermission.ts
import { useAuthStore } from '~/stores/auth'

export function usePermission() {
  const authStore = useAuthStore()

  /**
   * Memeriksa apakah user memiliki permission tertentu.
   * @param permission Nama permission yang dibutuhkan (contoh: 'products:create')
   */
  const hasPermission = (permission: string): boolean => {
    if (!authStore.user) return false
    
    // Bypass otomatis di frontend untuk role programmer & administrator
    if (authStore.user.role === 'programmer' || authStore.user.role === 'administrator') {
      return true
    }

    return authStore.userPermissions.includes(permission)
  }

  return {
    hasPermission
  }
}
```

---

## 4. Proteksi Halaman (Route Middleware / Navigation Guard)

Gunakan middleware routing untuk mencegah navigasi ke halaman sensitif sebelum halaman tersebut dirender.

### Contoh Middleware Global di Nuxt 3:
```typescript
// middleware/auth-guard.global.ts
import { useAuthStore } from '~/stores/auth'
import { usePermission } from '~/composables/usePermission'

export default defineNuxtRouteMiddleware(async (to, from) => {
  const authStore = useAuthStore()
  const { hasPermission } = usePermission()

  // 1. Cek apakah route membutuhkan autentikasi
  if (to.meta.middleware?.includes('auth') || to.meta.requiredPermission) {
    if (!authStore.isLoggedIn) {
      // Pastikan data user dimuat ulang jika halaman direfresh
      await authStore.fetchUser()
    }

    if (!authStore.isLoggedIn) {
      return navigateTo('/login')
    }

    // 2. Cek kebutuhan permission spesifik halaman
    const requiredPermission = to.meta.requiredPermission as string | undefined
    if (requiredPermission && !hasPermission(requiredPermission)) {
      // Tampilkan notifikasi toast/alert (opsional)
      // triggerToast("Anda tidak memiliki hak akses untuk membuka halaman ini!", "error")
      
      // Redirect ke halaman dashboard atau 403 Forbidden
      return navigateTo('/403')
    }
  }
})
```

### Cara Menentukan Permission di Tiap Halaman (Vue/Nuxt Page):
Cukup definisikan properti `requiredPermission` pada blok `definePageMeta` di bagian atas file halaman Anda.

```vue
<!-- pages/products/create.vue -->
<script setup lang="ts">
definePageMeta({
  middleware: ['auth'],
  requiredPermission: 'products:create' // Halaman ini hanya boleh diakses jika memiliki permission ini
})
</script>

<template>
  <div class="container mx-auto py-6">
    <h1 class="text-2xl font-bold">Tambah Produk Baru</h1>
    <!-- Form Tambah Produk -->
  </div>
</template>
```

---

## 5. Proteksi Elemen UI di Template Halaman

Selain mengamankan link/route URL, tombol aksi sensitif (seperti tombol edit, hapus, tambah) harus disembunyikan menggunakan arahan `v-if` agar antarmuka bersih dan intuitif.

```vue
<!-- components/ProductTable.vue -->
<script setup lang="ts">
import { usePermission } from '~/composables/usePermission'

const { hasPermission } = usePermission()
</script>

<template>
  <div class="flex justify-between items-center mb-4">
    <h2 class="text-xl font-semibold">Daftar Produk</h2>
    
    <!-- Tombol "Tambah Produk" hanya muncul jika user punya hak akses 'products:create' -->
    <button 
      v-if="hasPermission('products:create')" 
      @click="$router.push('/products/create')"
      class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
    >
      Tambah Produk
    </button>
  </div>

  <table>
    <thead>
      <tr>
        <th>Nama Produk</th>
        <th>Harga</th>
        <th>Aksi</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="product in products" :key="product.id">
        <td>{{ product.name }}</td>
        <td>{{ product.price }}</td>
        <td class="space-x-2">
          <!-- Tombol Aksi Edit & Hapus juga diproteksi -->
          <button v-if="hasPermission('products:edit')" class="text-yellow-600">Edit</button>
          <button v-if="hasPermission('products:delete')" class="text-red-600">Hapus</button>
        </td>
      </tr>
    </tbody>
  </table>
</template>
```

---

## 6. Ringkasan Konvensi Penamaan Permission
Pastikan nama permission di frontend sinkron dengan yang didefinisikan di backend (lihat berkas `route.go` atau tabel `permissions` di database):
- `products:view` - Melihat produk / daftar produk
- `products:create` - Menambahkan produk baru
- `products:edit` - Mengubah data produk
- `products:delete` - Menghapus produk
- `purchase-orders:view` - Melihat purchase orders
- `purchase-orders:create` - Membuat purchase orders baru
- dan seterusnya.
