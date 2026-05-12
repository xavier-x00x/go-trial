#!/usr/bin/env python3
"""Import products from MASTER_PRODUCT.txt to database using the proposal API."""

import csv
import json
import sys
import uuid

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

def get_existing_skus():
    """Get existing product SKUs from database."""
    import requests
    
    # Try to get from proposals first
    token = login()
    headers = {'Authorization': f'Bearer {token}'}
    
    # Get all PRODUCT proposals
    resp = requests.get('http://localhost:3050/api/master-data?entity_type=PRODUCT', headers=headers)
    if resp.status_code == 200:
        data = resp.json()
        if 'data' in data:
            for p in data['data']:
                payload = p.get('payload_json', '{}')
                try:
                    payload_data = json.loads(payload)
                    if 'sku' in payload_data:
                        EXISTING_SKUS.add(payload_data['sku'])
                except:
                    pass
    
    return EXISTING_SKUS

def parse_csv(filename):
    """Parse the MASTER_PRODUCT.txt CSV file."""
    products = []
    with open(filename, 'r', encoding='utf-8') as f:
        reader = csv.reader(f, delimiter=';')
        next(reader)  # Skip header
        
        for row in reader:
            if len(row) < 12:
                continue
            
            # Remove quotes from each field
            product = {
                'sku': row[0].strip('"'),
                'name': row[1].strip('"'),
                'ket_uk': row[2].strip('"'),
                'barcode': row[3].strip('"'),
                'markup': row[4].strip('"'),
                'harga_beli': row[5].strip('"'),
                'supplier': row[6].strip('"'),
                'max_stock': row[7].strip('"'),
                'rata2_penjualan': row[8].strip('"'),
                'length': row[9].strip('"'),
                'width': row[10].strip('"'),
                'height': row[11].strip('"'),
            }
            products.append(product)
    
    return products

def create_product_proposal(token, product_data, dry_run=False):
    """Create a product proposal via the API."""
    import requests
    
    # Build the product payload
    payload = {
        'sku': product_data['sku'],
        'name': product_data['name'],
        'base_uom_id': BASE_UOM_ID,
        'is_stockable': True,
        'is_stackable': True,
    }
    
    # Add optional fields
    if product_data['barcode']:
        payload['barcode'] = product_data['barcode']
    
    # Parse dimensions
    try:
        length = float(product_data['length']) if product_data['length'] else 0
        width = float(product_data['width']) if product_data['width'] else 0
        height = float(product_data['height']) if product_data['height'] else 0
    except ValueError:
        length = width = height = 0
    
    if length > 0:
        payload['length'] = str(length)
    if width > 0:
        payload['width'] = str(width)
    if height > 0:
        payload['height'] = str(height)
    
    # Create proposal request
    proposal_payload = {
        'entity_type': 'PRODUCT',
        'action_type': 'CREATE',
        'payload_json': json.dumps(payload),
        'reason': f"Bulk import: {product_data['sku']} - {product_data['name']}"
    }
    
    if dry_run:
        print(f"[DRY RUN] Would create: {product_data['sku']} - {product_data['name']}")
        return True
    
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
        print(f"Error creating {product_data['sku']}: {resp.status_code} - {resp.text[:200]}")
        return False

def main():
    import requests
    
    # Check arguments
    dry_run = '--dry-run' in sys.argv
    
    print("Logging in...")
    token = login()
    print(f"Logged in: {token[:30]}...")
    
    print("Loading existing SKUs...")
    existing_skus = get_existing_skus()
    print(f"Found {len(existing_skus)} existing SKUs")
    
    print("Parsing CSV...")
    products = parse_csv('MASTER_PRODUCT.txt')
    print(f"Found {len(products)} products in CSV")
    
    # Optionally limit for testing
    limit = 10 if dry_run else len(products)
    
    success = 0
    failed = 0
    skipped = 0
    
    print(f"Creating {limit} product proposals...")
    for i, p in enumerate(products[:limit]):
        # Skip if already exists
        if p['sku'] in existing_skus:
            skipped += 1
            continue
        
        if create_product_proposal(token, p, dry_run):
            success += 1
            print(f"[{i+1}] OK: {p['sku']} - {p['name']}")
            existing_skus.add(p['sku'])  # Add to avoid duplicates
        else:
            failed += 1
        
        # Progress indicator
        if (i + 1) % 50 == 0:
            print(f"Progress: {i+1}/{limit} (success: {success}, failed: {failed}, skipped: {skipped})")
    
    print(f"\nDone! Success: {success}, Failed: {failed}, Skipped: {skipped}")

if __name__ == '__main__':
    main()