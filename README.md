IMPS SUCCESS FLOW (Happy Path)
Customer sends IMPS transfer:
 - Customer sends IMPS transfer:
  
STEP 1 — API Gateway → TM
Gateway calls TM API:
 - POST /transactions

STEP 2 — TM creates initial state
TM does:
transactions table
id	status	saga_status
TX100	INITIATED	STARTED

saga_steps
step	status
START	DONE

no ledger yet
-----------------
STEP 3 — TM emits first Saga Command

TM publishes Kafka command:
{
  "command": "DEBIT_ACCOUNT",
  "transactionId": "TX100",
  "accountRefId": "ACC1",
  "amount": 1000
}

TM also saves:
saga_steps
step	status
DEBIT_REQUESTED	REQUESTED

--------------
STEP 4 — Account Service consumes command

Kafka consumer in Account Service picks:

DEBIT_ACCOUNT


Account Service:

• checks balance
• deducts money
• updates its DB

Then emits event.
-------------------
STEP 5 — Account Service emits event
Topic: tm.events
{
  "event": "ACCOUNT_DEBITED",
  "transactionId": "TX100",
  "accountRefId": "ACC1",
  "amount": 1000
}

----------------------
STEP 6 — TM consumes debit event

TM Kafka consumer picks this event.

TM now:

saga_steps
step	status
DEBIT_REQUESTED	DONE

ledger_entries
account	D/C	amount
ACC1	D	1000

transactions
status
IN_PROGRESS

-----------
STEP 7 — TM emits CREDIT command

TM publishes:
{
  "command": "CREDIT_ACCOUNT",
  "transactionId": "TX100",
  "accountRefId": "ACC2",
  "amount": 1000
}

Saga step:
step	status
CREDIT_REQUESTED	REQUESTED
----------------
STEP 8 — Account Service credits

Account Service processes and emits:
{
  "event": "ACCOUNT_CREDITED",
  "transactionId": "TX100",
  "accountRefId": "ACC2",
  "amount": 1000
}


---------------------------
STEP 9 — TM consumes credit event

TM does:

saga_steps
step	status
CREDIT_REQUESTED	DONE

ledger_entries
account	D/C	amount
ACC2	C	1000


transactions
status	saga_status
COMPLETED	COMPLETED



curl -X POST http://localhost:8080/api/v1/transactions \
-H "Content-Type: application/json" \
-d '{
  "paymentType":"IMMEDIATE",
  "paymentMode":"IMPS",
  "amount":2000,
  "currency":"INR"
}'


IMMEDATE IMPS request:
{
  "userRefId": "USR2001",
  "sourceRefId": "ACC-USR2001-01",
  "destinationRefId": "BANK",
  "paymentType": "IMMEDATE",
  "paymentMode": "IMPS",
  "dcFlag": "D",
  "amount": 2500,
  "currency": "INR",
  "networkTxnId": "",
  "gatewayTxnId": "BNK-GTWY-TX01"
}