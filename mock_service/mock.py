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

# ---------------- Flask (health only) ----------------

app = Flask(__name__)

@app.route("/health")
def health():
    return {"status": "UP"}

# ===================================================
# ACCOUNT SERVICE (STATUS DRIVEN)
# ===================================================

def account_block():
    c = consumer("account.commands.balance-block", "account-svc")
    for msg in c:
        data = msg.value

        # EXPECTED STATUS
        # T=BLOCK_REQUESTED, S=BALANCE_BLOCK(IN_PROGRESS)
        print("[ACCOUNT] BLOCK CMD:", data)

        data["transaction_status"] = "BLOCKED"
        data["saga_status"] = "BALANCE_BLOCK_COMPLETED"
        data["ledger_status"] = "AUTH_HOLD"

        producer.send("account.events.balance-blocked", data)

def account_final_debit():
    c = consumer("account.commands.final-debit", "account-svc")
    for msg in c:
        data = msg.value

        # EXPECTED STATUS
        # T=DEBIT_REQUESTED, S=FINAL_DEBIT(IN_PROGRESS)
        print("[ACCOUNT] FINAL DEBIT CMD:", data)

        data["transaction_status"] = "COMPLETED"
        data["saga_status"] = "FINAL_DEBIT_COMPLETED"
        data["ledger_status"] = "DEBIT_SETTLEMENT"

        producer.send("account.events.balance-debited", data)

def account_release_hold():
    c = consumer("account.commands.release-hold", "account-svc")
    for msg in c:
        data = msg.value

        # EXPECTED STATUS
        # T=RELEASE_REQUESTED, S=RELEASE(IN_PROGRESS)
        print("[ACCOUNT] RELEASE HOLD CMD:", data)

        data["transaction_status"] = "FAILED"
        data["saga_status"] = "RELEASE_COMPLETED"
        data["ledger_status"] = "RELEASE"

        producer.send("account.events.balance-released", data)

# ===================================================
# PAYMENT NETWORK (STATUS DRIVEN)
# ===================================================

def payment_network():
    c = consumer("payment.commands.debit", "payment-network")
    for msg in c:
        data = msg.value

        # EXPECTED STATUS
        # T=NETWORK_REQUESTED, S=IMPS_DEBIT(IN_PROGRESS)
        print("[NETWORK] DEBIT REQ:", data)
        time.sleep(1)
        

        channel = data.get("channel", "IMPS")

        # Demo logic
        if channel == "IMPS":
            outcome = "SUCCESS"
            # outcome = "FAIL"
        else:
            outcome = random.choice(["SUCCESS", "FAIL"])

        if outcome == "SUCCESS":
            data["transaction_status"] = "NETWORK_SUCCESS"
            data["saga_status"] = "IMPS_DEBIT_COMPLETED"
            data["ledger_status"] = "AUTH_HOLD"

            producer.send("payment.events.debit-success", data)
            print("[NETWORK] SUCCESS")

        else:
            data["transaction_status"] = "NETWORK_FAILED"
            data["saga_status"] = "IMPS_DEBIT_FAILED"
            data["ledger_status"] = "AUTH_HOLD"

            producer.send("payment.events.debit-failed", data)
            print("[NETWORK] FAILED")

# ---------------- Thread bootstrap ----------------

def start():
    threading.Thread(target=account_block, daemon=True).start()
    threading.Thread(target=account_final_debit, daemon=True).start()
    threading.Thread(target=account_release_hold, daemon=True).start()
    threading.Thread(target=payment_network, daemon=True).start()

# ---------------- Main ----------------
start()
if __name__ == "__main__":
    app.run(port=7000)
