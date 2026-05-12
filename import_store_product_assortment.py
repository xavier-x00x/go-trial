#!/usr/bin/env python3
"""Import StoreProductAssortment from MASTER_PRODUCT.txt - batch insert."""

import csv
import subprocess
import uuid
import sys

STORE_ID = "019dcffc-088f-7850-ba6a-88b7e8fbec11"  # Toko Test

def parse_csv(filename):
    products = []
    with open(filename, 'r', encoding='utf-8') as f:
        reader = csv.reader(f, delimiter=';')
        next(reader)
        for row in reader:
            if len(row) < 12:
                continue
            products.append({
                'sku': row[0].strip('"'),
                'max_stock': row[7].strip('"'),
                'avg_sales': row[8].strip('"'),
            })
    return products

def get_product_sku_to_id():
    """Get product SKU to ID mapping from database."""
    result = subprocess.run([
        'docker', 'exec', 'mysql-db', 'mysql',
        '-uroot', '-prootpassword',
        '-e', 'SELECT id, sku FROM noto.products'
    ], capture_output=True, text=True)
    
    mapping = {}
    lines = result.stdout.strip().split('\n')[1:]
    for line in lines:
        if line:
            parts = line.split('\t')
            if len(parts) >= 2:
                mapping[parts[1]] = parts[0]
    
    return mapping

def insert_batch(values):
    if not values:
        return True
    
    sql = f"""
        INSERT IGNORE INTO noto.store_product_assortments 
        (id, store_id, product_id, status, display_facing_qty, safety_stock_qty, max_stock_qty, average_daily_sales)
        VALUES {', '.join(values)}
    """
    
    # Use echo to pass the SQL
    process = subprocess.Popen([
        'docker', 'exec', '-i', 'mysql-db', 'mysql',
        '-uroot', '-prootpassword', 'noto'
    ], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    
    stdout, stderr = process.communicate(input=sql.encode(), timeout=30)
    
    return process.returncode == 0

def main():
    print("Parsing CSV...")
    products = parse_csv('MASTER_PRODUCT.txt')
    print(f"Found {len(products)} products in CSV")
    
    print("Getting product mapping...")
    product_mapping = get_product_sku_to_id()
    print(f"Found {len(product_mapping)} products in DB")
    
    # Build values in batches
    success = 0
    batch_size = 100
    values = []
    
    for p in products:
        product_id = product_mapping.get(p['sku'])
        if not product_id:
            continue
        
        max_stock = int(p['max_stock']) if p['max_stock'] else 1
        avg_sales = float(p['avg_sales']) if p['avg_sales'] else 0
        new_id = str(uuid.uuid4())
        values.append(f"('{new_id}', '{STORE_ID}', '{product_id}', 'ACTIVE', 1, {max_stock}, {max_stock}, {avg_sales})")
    
    total = len(values)
    print(f"Inserting {total} records in batches of {batch_size}...")
    
    for i in range(0, total, batch_size):
        batch = values[i:i + batch_size]
        if insert_batch(batch):
            success += len(batch)
            print(f"Progress: {success}/{total}")
    
    # Verify
    result = subprocess.run([
        'docker', 'exec', 'mysql-db', 'mysql',
        '-uroot', '-prootpassword',
        '-e', 'SELECT COUNT(*) FROM noto.store_product_assortments'
    ], capture_output=True, text=True)
    
    count = result.stdout.strip().split('\n')[1]
    print(f"\nDone! Total StoreProductAssortment: {count}")

if __name__ == '__main__':
    main()