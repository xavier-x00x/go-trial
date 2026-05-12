#!/usr/bin/env python3
"""Import ProductPrice from MASTER_PRODUCT.txt."""

import csv
import json

BASE_UOM_ID = "019dd000-c945-7669-a461-1250bb5fc786"  # PCS
PRICE_LIST_ID = "019dcffc-0996-74e6-b58d-6766cac4a695"  # NORMAL

def login():
    import requests
    resp = requests.post('http://localhost:3050/api/auth/login', json={
        'identity': 'dwi@example.com',
        'password': '12345678'
    })
    data = resp.json()
    return data['data']['access_token']

def get_product_mapping(token):
    import requests
    headers = {'Authorization': f'Bearer {token}'}
    resp = requests.get('http://localhost:3050/api/products', headers=headers)
    if resp.status_code != 200:
        return {}
    data = resp.json()
    return {p['sku']: p['id'] for p in data.get('data', [])}

def parse_csv(filename):
    products = []
    with open(filename, 'r', encoding='utf-8') as f:
        reader = csv.reader(f, delimiter=';')
        next(reader)
        for row in reader:
            if len(row) < 12:
                continue
            product = {
                'sku': row[0].strip('"'),
                'markup': row[4].strip('"'),
                'harga_beli': row[5].strip('"'),
            }
            products.append(product)
    return products

def create_product_price_proposal(token, product_data, product_id):
    import requests
    
    harga_beli = float(product_data['harga_beli']) if product_data['harga_beli'] else 0
    markup = float(product_data['markup']) if product_data['markup'] else 0
    
    if harga_beli <= 0 or markup <= 0:
        return False
    
    # Calculate sell price: buy_price + (buy_price * markup / 100)
    sell_price = harga_beli + (harga_beli * markup / 100)
    
    payload = {
        'price_list_id': PRICE_LIST_ID,
        'product_id': product_id,
        'uom_id': BASE_UOM_ID,
        'markup_pct': str(markup),
        'sell_price': str(round(sell_price, 2)),
    }
    
    proposal_payload = {
        'entity_type': 'PRODUCT_PRICE',
        'action_type': 'CREATE',
        'payload_json': json.dumps(payload),
        'reason': f"Set price for {product_data['sku']}"
    }
    
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {token}'
    }
    
    resp = requests.post('http://localhost:3050/api/master-data', json=proposal_payload, headers=headers)
    return resp.status_code >= 200 and resp.status_code < 300

def main():
    print("Logging in...")
    token = login()
    print(f"Logged in")
    
    print("Getting product mapping...")
    product_mapping = get_product_mapping(token)
    print(f"Found {len(product_mapping)} products")
    
    print("Parsing CSV...")
    products = parse_csv('MASTER_PRODUCT.txt')
    print(f"Found {len(products)} products")
    
    success = 0
    failed = 0
    
    for i, p in enumerate(products):
        product_id = product_mapping.get(p['sku'])
        if not product_id:
            continue
        
        if create_product_price_proposal(token, p, product_id):
            success += 1
            if success % 500 == 0:
                print(f"Progress: {success}")
        else:
            failed += 1
    
    print(f"\nDone! Success: {success}, Failed: {failed}")

if __name__ == '__main__':
    main()