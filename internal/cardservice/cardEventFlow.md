| Event (entry_type) | What happened | Transaction.status update |
|-------------------|--------------|----------------------------|
| AUTH              | Amount blocked | AUTHORIZED |
| SETTLEMENT        | Settlement started | PROCESSING |
| DEBIT             | Money debited | PROCESSING |
| CREDIT            | Money credited to the merchant | COMPLETED |
| RELEASE           | Block removed | RELEASED |
| REVERSAL          | AUTH undone | FAILED |
| REFUND            | Money returned | REFUNDED |


### Transaction Table (Business View)

| Scenario | Direction | Payment Type | Transaction Status Journey |
|---------|-----------|--------------|-----------------------------|
| IMPS (outgoing) | OUTGOING | IMMEDIATE | INITIATED → COMPLETED |
| NEFT (outgoing) | OUTGOING | DEFERRED | INITIATED → PENDING → AUTHORIZED → PROCESSING → COMPLETED |
| Card (outgoing) | OUTGOING | DEFERRED | INITIATED → PENDING → AUTHORIZED → PROCESSING → COMPLETED |
| Card fail after auth | OUTGOING | DEFERRED | AUTHORIZED → FAILED |
| Cancel after auth | OUTGOING | DEFERRED | AUTHORIZED → RELEASED |
| Incoming NEFT / IMPS | INCOMING | — | COMPLETED |
| Refund to user | INCOMING | — | REFUNDED |
| Chargeback | INCOMING | CARD | CHARGEBACK |


### Ledger Table (Financial Truth)

| Scenario | Ledger Entry Type(s) | dc_flag | Meaning |
|---------|---------------------|--------|--------|
| IMPS outgoing success | DEBIT + CREDIT | D + C | Instant money movement |
| NEFT / Card auth | AUTH | D | Amount blocked |
| NEFT / Card settlement start | SETTLEMENT | — | Processing marker |
| Final debit | DEBIT | D | Money debited |
| Merchant credit | CREDIT | C | Money credited |
| Cancel after auth | RELEASE | — | Block removed |
| Fail after auth | REVERSAL | — | Block undone |
| Incoming payment | CREDIT | C | User balance increased |
| Refund | REFUND | C | Money returned |
| Chargeback | CHARGEBACK | D | Money taken back |


### Event → Ledger Entry Mapping

| Event | When it occurs | Ledger Entry Type | dc_flag | Balance Impact |
|------|---------------|-------------------|--------|----------------|
| PAYMENT_INITIATED | User submits payment | — | — | No impact |
| BANK_ACCEPTED | Bank accepts request (NEFT/Card) | AUTH | D | Amount blocked |
| AUTH_SUCCESS | Authorization successful | AUTH | D | Amount blocked |
| SETTLEMENT_STARTED | Clearing / batch started | SETTLEMENT | — | No impact (marker) |
| DEBIT_CONFIRMED | Money debited from payer | DEBIT | D | Balance reduced |
| CREDIT_CONFIRMED | Money credited to receiver | CREDIT | C | Balance increased |
| CANCEL_REQUESTED | User cancels after auth | RELEASE | — | Block removed |
| AUTH_FAILED | Authorization failed | — | — | No impact |
| SETTLEMENT_FAILED | Failed after auth | REVERSAL | — | Block undone |
| REFUND_PROCESSED | Money returned to user | REFUND | C | Balance increased |
| CHARGEBACK_RAISED | Dispute / chargeback | CHARGEBACK | D | Balance reduced |


### Example Flows

| Scenario | Event | Ledger Entry | Transaction Status |
|--------|-------|--------------|--------------------|
| IMPS success | DEBIT_CONFIRMED | DEBIT | PROCESSING |
|  | CREDIT_CONFIRMED | CREDIT | COMPLETED |
| NEFT/Card auth | AUTH_SUCCESS | AUTH | AUTHORIZED |
| NEFT settlement | SETTLEMENT_STARTED | SETTLEMENT | PROCESSING |
| Final success | CREDIT_CONFIRMED | CREDIT | COMPLETED |
| Cancel after auth | CANCEL_REQUESTED | RELEASE | RELEASED |
| Fail after auth | SETTLEMENT_FAILED | REVERSAL | FAILED |
| Incoming payment | CREDIT_CONFIRMED | CREDIT | COMPLETED |
| Refund | REFUND_PROCESSED | REFUND | REFUNDED |
