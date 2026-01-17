# 🏦 IMPS Success Flow — Happy Path

> **Architecture Style:** Saga Orchestration + Event Driven Ledger
>
> **Purpose:** Explain complete IMPS success transaction from API call to ledger balance.

---

## 🎯 Customer Action

Customer initiates IMPS transfer:

> **Send ₹1000 from ACC1 → ACC2**

---

## 🔹 STEP 1 — API Gateway → Transaction Manager (TM)

### API Call

```
POST /api/v1/transactions
```

### Request Body

| Field            | Value     |
| ---------------- | --------- |
| userRefId        | USER1     |
| sourceRefId      | ACC1      |
| destinationRefId | ACC2      |
| amount           | 1000      |
| paymentType      | IMMEDIATE |
| paymentMode      | IMPS      |

---

## 🔹 STEP 2 — TM Creates Initial State

### transactions table

| id    | status    | saga_status |
| ----- | --------- | ----------- |
| TX100 | INITIATED | STARTED     |

### saga_steps table

| step  | status |
| ----- | ------ |
| START | DONE   |

### Ledger

👉 No ledger entry yet (no money movement).

---

## 🔹 STEP 3 — TM Emits Debit Command

### Kafka Topic

```
tm.commands
```

### Message

```json
{
  "command": "DEBIT_ACCOUNT",
  "transactionId": "TX100",
  "accountRefId": "ACC1",
  "amount": 1000
}
```

### Saga State

| step            | status    |
| --------------- | --------- |
| DEBIT_REQUESTED | REQUESTED |

---

## 🔹 STEP 4 — Account Service Consumes Debit Command

Account Service:

• Validates balance
• Deducts ₹1000 from ACC1
• Updates its database
• Emits success event

---

## 🔹 STEP 5 — Account Service Emits Debit Event

### Kafka Topic

```
tm.events
```

### Message

```json
{
  "event": "ACCOUNT_DEBITED",
  "transactionId": "TX100",
  "accountRefId": "ACC1",
  "amount": 1000
}
```

---

## 🔹 STEP 6 — TM Consumes Debit Event

TM updates state:

### saga_steps

| step            | status |
| --------------- | ------ |
| DEBIT_REQUESTED | DONE   |

### ledger_entries

| account | D/C | amount |
| ------- | --- | ------ |
| ACC1    | D   | 1000   |

### transactions

| status      |
| ----------- |
| IN_PROGRESS |

---

## 🔹 STEP 7 — TM Emits Credit Command

### Kafka Message

```json
{
  "command": "CREDIT_ACCOUNT",
  "transactionId": "TX100",
  "accountRefId": "ACC2",
  "amount": 1000
}
```

### Saga State

| step             | status    |
| ---------------- | --------- |
| CREDIT_REQUESTED | REQUESTED |

---

## 🔹 STEP 8 — Account Service Credits Destination Account

Account Service:

• Credits ₹1000 into ACC2
• Updates its database
• Emits success event

### Event

```json
{
  "event": "ACCOUNT_CREDITED",
  "transactionId": "TX100",
  "accountRefId": "ACC2",
  "amount": 1000
}
```

---

## 🔹 STEP 9 — TM Consumes Credit Event

TM finalizes transaction:

### saga_steps

| step             | status |
| ---------------- | ------ |
| CREDIT_REQUESTED | DONE   |

### ledger_entries

| account | D/C | amount |
| ------- | --- | ------ |
| ACC2    | C   | 1000   |

### transactions

| status    | saga_status |
| --------- | ----------- |
| COMPLETED | COMPLETED   |

---

## 🎉 IMPS SUCCESS COMPLETE

### Final Ledger View

| Account | D/C | Amount |
| ------- | --- | ------ |
| ACC1    | D   | 1000   |
| ACC2    | C   | 1000   |

✔ Double entry balanced
✔ Saga completed
✔ Ledger immutable
✔ Customer experience consistent
✔ Bank books correct

---

## 🧠 Key Architecture Guarantees

| Layer       | Guarantee             |
| ----------- | --------------------- |
| API         | Idempotent request    |
| Saga        | Orchestration safety  |
| Kafka       | Event durability      |
| Ledger      | Audit correctness     |
| Transaction | Business traceability |

---


# 🏦 IMPS Ledger Documentation — Transaction Centric View

> **Purpose:** Explain how ledger entries are created, balanced, and compensated in a Same-Bank IMPS transaction using Saga pattern.

---

## 🎯 Scenario

| Field          | Value            |
| -------------- | ---------------- |
| Sender         | User A (ACC_A)   |
| Receiver       | User B (ACC_B)   |
| Amount         | ₹1000            |
| Transaction ID | **TX100**        |
| Transfer Type  | IMPS (Same Bank) |

---

## 🧾 Core Banking Principle

> **Ledger is immutable.**
>
> • Ledger never updates or deletes rows
> • Every financial event creates a **new row**
> • Compensation = new ledger entry
> • Final balance is derived by summation

---

## 🔹 Initial State

Ledger already has historical data.

No entry exists for **TX100** yet.

---

## 🔹 Step 1 — Transaction Created

### Transactions Table

| id    | status    | saga_status |
| ----- | --------- | ----------- |
| TX100 | INITIATED | STARTED     |

### Ledger

👉 No ledger entry yet (no financial movement).

---

## 🔹 Step 2 — Debit Success

TM receives: **ACCOUNT_DEBITED**

### Ledger Entry

| id  | transaction_id | account_ref_id | D/C | entry_type | amount |
| --- | -------------- | -------------- | --- | ---------- | ------ |
| L1  | TX100          | ACC_A          | D   | DEBIT      | 1000   |

### Transactions Table

| status      | saga_status |
| ----------- | ----------- |
| IN_PROGRESS | IN_PROGRESS |

---

## 🔹 Step 3 — Credit Success

TM receives: **ACCOUNT_CREDITED**

### Ledger Entry

| id  | transaction_id | account_ref_id | D/C | entry_type | amount |
| --- | -------------- | -------------- | --- | ---------- | ------ |
| L2  | TX100          | ACC_B          | C   | CREDIT     | 1000   |

### Transactions Table

| status    | saga_status |
| --------- | ----------- |
| COMPLETED | COMPLETED   |

---

## ✅ Final Ledger View (Success)

| transaction_id | account | D/C | amount |
| -------------- | ------- | --- | ------ |
| TX100          | ACC_A   | D   | 1000   |
| TX100          | ACC_B   | C   | 1000   |

### ✔ Banking Validation

| Rule                | Status |
| ------------------- | ------ |
| Debit = Credit      | ✅      |
| Same amount         | ✅      |
| Same currency       | ✅      |
| Same transaction id | ✅      |
| Balanced books      | ✅      |

---

## 👤 Customer View

| User   | Statement      |
| ------ | -------------- |
| User A | ₹1000 Debited  |
| User B | ₹1000 Credited |

---

## 🏦 Bank View

System observes a perfectly balanced double-entry record under TX100.

---

# ❌ Failure Scenario — Credit Fails

Debit already posted:

| TX100 | ACC_A | D | 1000 |

Credit fails → Saga triggers compensation.

---

## 🔹 Step 4 — Debit Reversal

TM receives: **DEBIT_REVERSED**

### Compensation Ledger Entry

| id  | transaction_id | account_ref_id | D/C | entry_type | amount |
| --- | -------------- | -------------- | --- | ---------- | ------ |
| L2  | TX100          | ACC_A          | C   | REVERSAL   | 1000   |

### Transactions Table

| status | saga_status |
| ------ | ----------- |
| FAILED | COMPENSATED |

---

## 🔁 Final Ledger View (Failure)

| transaction_id | account | D/C | amount |
| -------------- | ------- | --- | ------ |
| TX100          | ACC_A   | D   | 1000   |
| TX100          | ACC_A   | C   | 1000   |

### Net Effect

```
Debit 1000
Credit 1000
------------
Net = 0
```

👉 Customer money fully restored.
👉 Ledger remains immutable.

---

## 🧠 Why Banks Trust This Model

| Feature                | Benefit                |
| ---------------------- | ---------------------- |
| Immutable ledger       | Full audit trail       |
| Transaction ID linkage | Easy reconciliation    |
| Saga compensation      | Safe failure handling  |
| Append-only            | Regulatory compliant   |
| Double entry           | Accounting correctness |

---

## 🏁 Summary

| Concept       | Rule                |
| ------------- | ------------------- |
| Ledger writes | Append only         |
| Transaction   | Business state      |
| Saga          | Orchestration       |
| Compensation  | New ledger row      |
| Balance       | Derived, not stored |


