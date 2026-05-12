#!/bin/bash

# Login to get token
TOKEN=$(curl -s -X POST http://localhost:3050/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identity":"dwi@example.com","password":"12345678"}' | jq -r '.data.token')

echo "Token: ${TOKEN:0:50}..."

# Create products from MASTER_PRODUCT.txt
# Skip header, parse CSV with semicolon delimiter

tail -n +2 MASTER_PRODUCT.txt | while IFS=';' read -r sku name ket_uk barcode markup harga_beli supp max_stock rata2_penjualan panjang lebar tinggi; do
  # Remove quotes
  sku="${sku//\"/}"
  name="${name//\"/}"
  barcode="${barcode//\"/}"
  panjang="${panjang//\"/}"
  lebar="${lebar//\"/}"
  tinggi="${tinggi//\"/}"
  
  # Create product JSON
  PRODUCT_JSON=$(cat <<EOF
{
  "sku": "$sku",
  "name": "$name",
  "base_uom_id": "019dd000-c945-7669-a461-1250bb5fc786",
  "is_stockable": true,
  "length": "$panjang",
  "width": "$lebar",
  "height": "$tinggi",
  "barcode": "$barcode"
}
EOF
)
  
  # Create proposal request
  PROPOSAL_JSON=$(cat <<EOF
{
  "entity_type": "PRODUCT",
  "action_type": "CREATE",
  "payload_json": $(echo $PRODUCT_JSON | jq -c .),
  "reason": "Bulk import: $sku - $name"
}
EOF
)
  
  # Call API
  RESPONSE=$(curl -s -X POST http://localhost:3050/api/master-data/proposals \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$PROPOSAL_JSON")
  
  if echo "$RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo "OK: $sku - $name"
  else
    echo "FAIL: $sku - $(echo $RESPONSE | jq -r '.message // .error')"
  fi
  
  # Limit to 10 for testing
  count=$((count + 1))
  if [ $count -ge 10 ]; then
    break
  fi
  
done

echo "Done! Created $count products"