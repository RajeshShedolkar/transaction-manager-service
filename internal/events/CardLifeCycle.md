Card Lifecycle in TM
CARD_AUTH  →  CARD_SETTLEMENT  →  COMPLETED
        ↘
         CARD_AUTH_RELEASE → RELEASED

Card Payment flow:
Merchant → Network → Bank → Transaction Manager → Ledger


Characters in Card Payment

1. Customer
2. Merchant (Amazon, Flipkart, etc.)
3. Payment Gateway (Razorpay, Stripe, PayU)
4. Card Network (Visa / Mastercard / RuPay)
5. Issuing Bank (Customer’s bank)
6. Acquiring Bank (Merchant’s bank)
7. Transaction Manager (your service)

Big Picture

Customer never talks directly to bank.
Everything flows through merchant → gateway → network → banks.

Your Transaction Manager is inside the banking system side.


STEP 1 — Customer enters card details on merchant site
Customer enters:
 - Card number
 - Expiry
 - CVV
 - OTP

Merchant NEVER talks to bank directly.
Merchant sends details to Payment Gateway.

STEP 2 — Payment Gateway contacts Card Network

Gateway sends:
 - Card details + amount + merchant id
To:
 - Visa / Mastercard / RuPay

STEP 3 — Card Network contacts Issuing Bank
Network routes request to:
 - Customer’s bank (issuer)

STEP 4 — Issuing Bank AUTHORISES
Issuing bank checks:
 - Card valid?
 - OTP correct?
 - Balance available?

If YES:
 - Bank blocks the amount (not debited yet)
 - This is called AUTHORIZATION
Bank responds:
 - AUTH_SUCCESS
To:
 - to Network → Gateway → Merchant.
-------------------
Who blocks the money?
 - Issuing Bank blocks the money
 - Not merchant.
 - Not gateway.
 - Not network.
At this point:
Money is reserved, not transferred.
But Customer sees:
 - “Payment successful”
But technically:
Only AUTH is done, not settlement.

STEP 5 — How Transaction Manager knows?
Issuing bank publishes->CARD_AUTH event -> to internal systems.

Transaction Manager consumes this event.
Event contains:
 - Network transaction ID
 - Amount
 - Card type
 - Merchant info

TM does based on event:
transactions table - row entry with AUTHORIZED
ledger table - row entry with AUTH


STEP 6:
Settlement happens
Gateway → Network → Issuing bank
Issuing bank now:

 - Actually transfers money
 - Removes block
 - Credits acquiring bank

This is SETTLEMENT

STEP 7 — Transaction Manager gets settlement
Issuing bank emits:
 - CARD_SETTLEMENT event

TM does based on event:
transactions table - row entry with COMPLETED
ledger table - row entry with SETTLEMENT

OTHER CASE:
If merchant cancels
Issuing bank:
 - Releases blocked amount.
TM gets:
 - CARD_AUTH_RELEASE event

TM updates:
| transactions | RELEASED |
| ledger | RELEASE |

