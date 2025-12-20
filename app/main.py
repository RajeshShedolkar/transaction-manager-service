from fastapi import FastAPI
from app.api.transfer import router as transfer_router

app = FastAPI(title="Transaction Manager Service")

app.include_router(transfer_router, prefix="/transactions")

@app.get("/health")
def health():
    return {"status": "UP"}