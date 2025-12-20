# Transaction Manager Service

## Overview
This service orchestrates fund transfers (IMPS / UPI / NEFT) using Saga Orchestration.

## How it works
1. Validate request
2. Check balance (mocked)
3. Debit source account
4. Credit destination account
5. Rollback on failure
6. Record transaction in ledger

## Run
pip install fastapi uvicorn
uvicorn app.main:app --reload

## API
POST /transactions/transfer