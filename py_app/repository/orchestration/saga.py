from app.repository.ledger_repo import save_transaction
from app.clients.account_client import (
    check_balance, debit, credit, rollback
)

def process_transfer(req):
    if not check_balance(req.fromAccountId, req.amount):
        return {"status": "FAILED", "reason": "Insufficient balance"}

    debit_result = debit(req.fromAccountId, req.amount)
    if not debit_result:
        return {"status": "FAILED", "reason": "Debit failed"}

    credit_result = credit(req.toAccountId, req.amount)
    if not credit_result:
        rollback(req.fromAccountId, req.amount)
        save_transaction(req, "ROLLED_BACK")
        return {"status": "ROLLED_BACK"}

    save_transaction(req, "COMPLETED")
    return {"status": "COMPLETED"}