#!/usr/bin/env python3
"""Import ProductPrice and ProductSupplier from MASTER_PRODUCT.txt."""

import csv
import json
import sys

BASE_UOM_ID = "019dd000-c945-7669-a461-1250bb5fc786"  # PCS
EXISTING_SKUS = set()

def login():
    """Login and get access token."""
    import requests
    resp = requests.post('http://localhost:3050/api/auth/login', json={
        'identity': 'dwi@example.com',
        'password': '12345678'
    })
    data = resp.json()
    return data['data']['access_token']

def get_supplier_mapping(token):
    """Get supplier code to ID mapping."""
    import requests
    headers = {'Authorization': f'Bearer {token}'}
    
    resp = requests.get('http://localhost:3050/api/suppliers', headers=headers)
    if resp.status_code != 200:
        print(f"Error getting suppliers: {resp.status_code}")
        return {}
    
    data = resp.json()
    mapping = {}
    for s in data.get('data', []):
        code = s.get('code', '')
        supplier_id = s.get('id', '')
        mapping[code] = supplier_id
    
    return mapping

def get_product_mapping(token):
    """Get product SKU to ID mapping."""
    import requests
    headers = {'Authorization': f'Bearer {token}'}
    
    resp = requests.get('http://localhost:3050/api/products', headers=headers)
    if resp.status_code != 200:
        print(f"Error getting products: {resp.status_code}")
        return {}
    
    data = resp.json()
    mapping = {}
    for p in data.get('data', []):
        sku = p.get('sku', '')
        product_id = p.get('id', '')
        mapping[sku] = product_id
    
    return mapping

def parse_csv(filename):
    """Parse the MASTER_PRODUCT.txt CSV file."""
    products = []
    with open(filename, 'r', encoding='utf-8') as f:
        reader = csv.reader(f, delimiter=';')
        next(reader)  # Skip header
        
        for row in reader:
            if len(row) < 12:
                continue
            
            product = {
                'sku': row[0].strip('"'),
                'name': row[1].strip('"'),
                'barcode': row[3].strip('"'),
                'markup': row[4].strip('"'),
                'harga_beli': row[5].strip('"'),
                'supplier': row[6].strip('"'),
                'max_stock': row[7].strip('"'),
            }
            products.append(product)
    
    return products

def create_product_supplier_proposal(token, product_data, product_id, supplier_id, product_mapping):
    """Create a product supplier proposal via the API."""
    import requests
    
    offered_price = float(product_data['harga_beli']) if product_data['harga_beli'] else 0
    max_stock = int(product_data['max_stock']) if product_data['max_stock'] else 1
    
    if offered_price <= 0:
        return False
    
    # Build the product supplier payload
    payload = {
        'product_id': product_id,
        'supplier_id': supplier_id,
        'supplier_sku': product_data['sku'],
        'is_primary': True,
        'offered_price': str(offered_price),
        'min_order_qty': str(max_stock),
    }
    
    # Create proposal request
    proposal_payload = {
        'entity_type': 'PRODUCT_SUPPLIER',
        'action_type': 'CREATE',
        'payload_json': json.dumps(payload),
        'reason': f"Link supplier for {product_data['sku']}"
    }
    
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {token}'
    }
    
    resp = requests.post(
        'http://localhost:3050/api/master-data',
        json=proposal_payload,
        headers=headers
    )
    
    if resp.status_code >= 200 and resp.status_code < 300:
        return True
    else:
        print(f"Error creating product supplier for {product_data['sku']}: {resp.status_code} - {resp.text[:100]}")
        return False

def main():
    import requests
    
    print("Logging in...")
    token = login()
    print(f"Logged in: {token[:30]}...")
    
    print("Getting supplier mapping...")
    supplier_mapping = get_supplier_mapping(token)
    print(f"Found {len(supplier_mapping)} suppliers")
    
    print("Getting product mapping...")
    product_mapping = get_product_mapping(token)
    print(f"Found {len(product_mapping)} products")
    
    print("Parsing CSV...")
    products = parse_csv('MASTER_PRODUCT.txt')
    print(f"Found {len(products)} products in CSV")
    
    # For each product
    success = 0
    failed = 0
    
    for i, p in enumerate(products):  # Process all products
        sku = p['sku']
        supplier_code = p['supplier']
        
        # Get product ID
        product_id = product_mapping.get(sku)
        if not product_id:
            continue
        
        # Get supplier ID
        supplier_id = supplier_mapping.get(supplier_code)
        if not supplier_id:
            # Try to find by partial match
            supplier_id = None
            for code, sid in supplier_mapping.items():
                if supplier_code in code or code in supplier_code:
                    supplier_id = sid
                    break
            
            if not supplier_id:
                continue
        
        if create_product_supplier_proposal(token, p, product_id, supplier_id, product_mapping):
            success += 1
            print(f"[{i+1}] OK: {sku} - {supplier_code}")
        else:
            failed += 1
        
        if (i + 1) % 20 == 0:
            print(f"Progress: {i+1} (success: {success}, failed: {failed})")
    
    print(f"\nDone! Success: {success}, Failed: {failed}")

if __name__ == '__main__':
    main()