#!/usr/bin/env python3
"""Approve and execute all pending product proposals."""

import json
import sys

def login():
    """Login and get access token."""
    import requests
    resp = requests.post('http://localhost:3050/api/auth/login', json={
        'identity': 'dwi@example.com',
        'password': '12345678'
    })
    data = resp.json()
    return data['data']['access_token']

def approve_and_execute(token):
    """Approve and execute all pending product proposals."""
    import requests
    
    headers = {
        'Content-Type': 'application/json',
        'Authorization': f'Bearer {token}'
    }
    
    # Get all pending proposals for PRODUCT
    resp = requests.get('http://localhost:3050/api/master-data/entity/PRODUCT?status=PENDING', headers=headers)
    if resp.status_code != 200:
        print(f"Error getting proposals: {resp.status_code}")
        return 0, 0
    
    data = resp.json()
    proposals = data.get('data', [])
    print(f"Found {len(proposals)} pending PRODUCT proposals")
    
    success = 0
    failed = 0
    
    for p in proposals:
        proposal_id = p['id']
        
        # Approve the proposal
        review_resp = requests.post(
            f'http://localhost:3050/api/master-data/{proposal_id}/review',
            json={'status': 'APPROVED'},
            headers=headers
        )
        
        if review_resp.status_code >= 400:
            print(f"[{proposal_id[:8]}] Failed to approve: {review_resp.status_code}")
            failed += 1
            continue
        
        # Execute the proposal (create the product)
        exec_resp = requests.post(
            f'http://localhost:3050/api/master-data/{proposal_id}/execute',
            headers=headers
        )
        
        if exec_resp.status_code >= 400:
            print(f"[{proposal_id[:8]}] Failed to execute: {exec_resp.status_code} - {exec_resp.text[:100]}")
            failed += 1
        else:
            success += 1
            if success % 100 == 0:
                print(f"Progress: {success} approved & executed")
    
    return success, failed

def main():
    print("Logging in...")
    token = login()
    print(f"Logged in")
    
    print("\nApproving and executing product proposals...")
    success, failed = approve_and_execute(token)
    
    print(f"\nDone! Success: {success}, Failed: {failed}")

if __name__ == '__main__':
    main()