#!/usr/bin/env python3
"""Approve and execute all pending product price proposals."""

import json

def login():
    import requests
    resp = requests.post('http://localhost:3050/api/auth/login', json={
        'identity': 'dwi@example.com',
        'password': '12345678'
    })
    return resp.json()['data']['access_token']

def approve_and_execute(token):
    import requests
    headers = {'Authorization': f'Bearer {token}'}
    
    resp = requests.get('http://localhost:3050/api/master-data/entity/PRODUCT_PRICE?status=PENDING', headers=headers)
    proposals = resp.json().get('data', [])
    print(f"Found {len(proposals)} pending PRODUCT_PRICE proposals")
    
    success = 0
    failed = 0
    
    for p in proposals:
        proposal_id = p['id']
        
        review_resp = requests.post(
            f'http://localhost:3050/api/master-data/{proposal_id}/review',
            json={'status': 'APPROVED'},
            headers=headers
        )
        
        if review_resp.status_code >= 400:
            failed += 1
            continue
        
        exec_resp = requests.post(
            f'http://localhost:3050/api/master-data/{proposal_id}/execute',
            headers=headers
        )
        
        if exec_resp.status_code >= 400:
            failed += 1
        else:
            success += 1
            if success % 500 == 0:
                print(f"Progress: {success}")
    
    return success, failed

def main():
    token = login()
    print("Approving and executing product price proposals...")
    success, failed = approve_and_execute(token)
    print(f"\nDone! Success: {success}, Failed: {failed}")

if __name__ == '__main__':
    main()