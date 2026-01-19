from kafka import KafkaProducer, KafkaConsumer
from flask import Flask
import threading
import json
import time
import random

BROKER = "localhost:9092"

# ---------------- Kafka setup ----------------

producer = KafkaProducer(
    bootstrap_servers=BROKER,
    value_serializer=lambda v: json.dumps(v).encode("utf-8")
)

def consumer(topic, group):
    return KafkaConsumer(
        topic,
        bootstrap_servers=BROKER,
        group_id=group,
        auto_offset_reset="earliest",
        value_deserializer=lambda v: json.loads(v.decode("utf-8")),
        enable_auto_commit=True,
    )

# ---------------- Flask ----------------

app = Flask(__name__)

@app.route("/health")
def health():
    return {"status": "UP"}

# ===================================================
# IMPS ACCOUNT SERVICE
# ===================================================

def account_block():
    c = consumer("account.commands.balance-block", "account-group")
    for msg in c:
        data = msg.value
        if data.get("source") == "ACCOUNT":
            continue

        print("[ACCOUNT] BLOCK CMD:", data)

        data["transaction_status"] = "BLOCKED"
        data["saga_status"] = "BALANCE_BLOCK_COMPLETED"
        data["ledger_status"] = "AUTH_HOLD"
        data["source"] = "ACCOUNT"

        producer.send("account.events.balance-blocked", data)

def account_final_debit():
    c = consumer("account.commands.final-debit", "account-group")
    for msg in c:
        data = msg.value
        if data.get("source") == "ACCOUNT":
            continue

        print("[ACCOUNT] FINAL DEBIT CMD:", data)

        data["transaction_status"] = "COMPLETED"
        data["saga_status"] = "FINAL_DEBIT_COMPLETED"
        data["ledger_status"] = "DEBIT_SETTLEMENT"
        data["source"] = "ACCOUNT"

        producer.send("account.events.balance-debited", data)

def account_release_hold():
    c = consumer("account.commands.release-hold", "account-group")
    for msg in c:
        data = msg.value
        if data.get("source") == "ACCOUNT":
            continue

        print("[ACCOUNT] RELEASE HOLD CMD:", data)

        data["transaction_status"] = "FAILED"
        data["saga_status"] = "RELEASE_COMPLETED"
        data["ledger_status"] = "RELEASE"
        data["source"] = "ACCOUNT"

        producer.send("account.events.balance-released", data)

# ===================================================
# IMPS PAYMENT NETWORK
# ===================================================

def payment_network():
    c = consumer("payment.commands.debit", "network-group")
    for msg in c:
        data = msg.value
        if data.get("source") == "NETWORK":
            continue

        print("[NETWORK-IMPS] DEBIT REQ:", data)
        time.sleep(1)

        outcome = "SUCCESS" if data.get("channel") == "IMPS" else random.choice(["SUCCESS","FAIL"])

        data["source"] = "NETWORK"

        if outcome == "SUCCESS":
            data["transaction_status"] = "NETWORK_SUCCESS"
            data["saga_status"] = "IMPS_DEBIT_COMPLETED"
            data["ledger_status"] = "AUTH_HOLD"
            producer.send("payment.events.debit-success", data)
            print("[NETWORK-IMPS] SUCCESS")
        else:
            data["transaction_status"] = "NETWORK_FAILED"
            data["saga_status"] = "IMPS_DEBIT_FAILED"
            data["ledger_status"] = "AUTH_HOLD"
            producer.send("payment.events.debit-failed", data)
            print("[NETWORK-IMPS] FAILED")

# ===================================================
# NEFT ACCOUNT SERVICE
# ===================================================

def neft_account_block():
    c = consumer("neft.payment.commands.debit", "neft-account-group")
    for msg in c:
        data = msg.value
        if data.get("source") == "ACCOUNT":
            continue

        print("[NEFT ACCOUNT] BLOCK CMD:", data)

        data["transaction_status"] = "BLOCKED"
        data["saga_status"] = "BALANCE_BLOCK_COMPLETED"
        data["ledger_status"] = "AUTH_HOLD"
        data["source"] = "ACCOUNT"

        producer.send("neft.account.events.balance-blocked", data)

def neft_account_final_debit():
    c = consumer("account.commands.final-debit", "neft-account-group")
    for msg in c:
        data = msg.value
        if data.get("source") == "ACCOUNT":
            continue

        print("[NEFT ACCOUNT] FINAL DEBIT CMD:", data)

        data["transaction_status"] = "COMPLETED"
        data["saga_status"] = "FINAL_DEBIT_COMPLETED"
        data["ledger_status"] = "DEBIT_SETTLEMENT"
        data["source"] = "ACCOUNT"

        producer.send("neft.account.events.balance-debited", data)

def neft_account_release_hold():
    c = consumer("account.commands.release-hold", "neft-account-group")
    for msg in c:
        data = msg.value
        if data.get("source") == "ACCOUNT":
            continue

        print("[NEFT ACCOUNT] RELEASE HOLD CMD:", data)

        data["transaction_status"] = "FAILED"
        data["saga_status"] = "RELEASE_COMPLETED"
        data["ledger_status"] = "RELEASE"
        data["source"] = "ACCOUNT"

        producer.send("neft.account.events.balance-released", data)

# ===================================================
# NEFT PAYMENT NETWORK
# ===================================================

def neft_payment_network():
    c = consumer("neft.payment.commands.debit", "neft-network-group")
    for msg in c:
        data = msg.value
        if data.get("source") == "NETWORK":
            continue

        print("[NETWORK-NEFT] DEBIT REQ:", data)
        time.sleep(2)

        outcome = random.choice(["SUCCESS","FAIL"])
        data["source"] = "NETWORK"

        if outcome == "SUCCESS":
            data["transaction_status"] = "NETWORK_SUCCESS"
            data["saga_status"] = "NEFT_DEBIT_COMPLETED"
            data["ledger_status"] = "AUTH_HOLD"
            producer.send("neft.payment.events.debit-success", data)
            print("[NETWORK-NEFT] SUCCESS")
        else:
            data["transaction_status"] = "NETWORK_FAILED"
            data["saga_status"] = "NEFT_DEBIT_FAILED"
            data["ledger_status"] = "AUTH_HOLD"
            producer.send("neft.payment.events.debit-failed", data)
            print("[NETWORK-NEFT] FAILED")

# ---------------- Thread bootstrap ----------------

def start():
    # IMPS
    threading.Thread(target=account_block, daemon=True).start()
    threading.Thread(target=account_final_debit, daemon=True).start()
    threading.Thread(target=account_release_hold, daemon=True).start()
    threading.Thread(target=payment_network, daemon=True).start()

    # NEFT
    threading.Thread(target=neft_account_block, daemon=True).start()
    threading.Thread(target=neft_account_final_debit, daemon=True).start()
    threading.Thread(target=neft_account_release_hold, daemon=True).start()
    threading.Thread(target=neft_payment_network, daemon=True).start()

# ---------------- Main ----------------

start()

if __name__ == "__main__":
    app.run(port=7000)
