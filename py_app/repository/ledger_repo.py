# Mock Ledger Repository

ledger = []

def save_transaction(req, status):
    ledger.append({
        "from": req.fromAccountId,
        "to": req.toAccountId,
        "amount": req.amount,
        "status": status
    })