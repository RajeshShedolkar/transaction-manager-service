from pydantic import BaseModel
from enum import Enum

class TransferMethod(str, Enum):
    IMPS = "IMPS"
    UPI = "UPI"
    NEFT = "NEFT"

class TransferRequest(BaseModel):
    fromAccountId: str
    toAccountId: str
    amount: float
    transferMethod: TransferMethod
    idempotencyKey: str