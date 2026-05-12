import json

base_data = {
    "no": "PO/2024/001",
    "tanggal": "15 Januari 2024",
    "supplier": "PT Sumber Rejeki",
    "alamat": "Jl. Merdeka No. 123, Jakarta Pusat",
    "telp": "(021) 1234567",
    "email": "info@sumberrekaji.com",
    "store": "Toko Utama Jakarta",
    "namaItem": "Item Detail Ke-",
    "qty": 1,
    "harga": 10000,
    "total": 10000,
    "sat": "pcs",
    "subtotal": 500000,
    "diskon": 0,
    "ppn": 55000,
    "grandtotal": 555000,
    "paymentTerms": 30,
    "paymentMode": "Tempo",
    "approvedBy": "John Doe",
    "notes": "Pengiriman dilakukan hari Senin-Jumat"
}

data_list = []
for i in range(1, 51):
    item = base_data.copy()
    item["namaItem"] = f"Item Detail Ke-{i} yang memiliki nama sangat panjang sekali untuk mengetest fitur stretch height dan page break agar terlihat apakah berfungsi dengan baik atau tidak pada saat mencapai batas halaman"
    item["qty"] = i
    item["total"] = i * 10000
    data_list.append(item)

with open('test_output/po_test_data.json', 'w') as f:
    json.dump(data_list, f, indent=2)

print("Generated 50 items in test_output/po_test_data.json")
