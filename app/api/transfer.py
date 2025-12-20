from fastapi import APIRouter, HTTPException
from app.domain.models import TransferRequest
from app.orchestration.saga import process_transfer

router = APIRouter()

@router.post("/transfer")
def transfer_funds(request: TransferRequest):
    return process_transfer(request)