┌──────┐
│ USER │
└───┬──┘
    │ 1. IMPS transfer request (HTTP)
    │
    │ T=INITIATED
    │ S=INIT
    │ L=NONE
    ▼
┌────────────────────────────┐
│ Transaction Manager (TM)   │
└───┬────────────────────────┘
    │
    │ PRODUCE
    │ account.commands.balance-block
    │
    │ T=BLOCK_REQUESTED
    │ S=BALANCE_BLOCK(IN_PROGRESS)
    │ L=NONE
    ▼
┌────────────────────────────┐
│ Account Service            │
└───┬────────────────────────┘
    │
    │ PRODUCE
    │ account.events.balance-blocked
    │
    │ T=BLOCKED
    │ S=BALANCE_BLOCK(COMPLETED)
    │ L=AUTH(HOLD)
    ▼
┌────────────────────────────┐
│ Transaction Manager (TM)   │
└───┬────────────────────────┘
    │
    │ PRODUCE
    │ payment.commands.debit (channel=IMPS)
    │
    │ T=NETWORK_REQUESTED
    │ S=IMPS_DEBIT(IN_PROGRESS)
    │ L=AUTH(HOLD)
    ▼
┌──────────────────────────────────────────────┐
│ Internal Payment Network + NPCI              │
│ (IMPS Adapter + External IMPS Switch)        │
└───┬──────────────────────────────────────────┘
    │
    │ NPCI processing
    │
    │ PRODUCE
    │ payment.events.debit-success
    │ OR payment.events.debit-failed / timeout
    │
    │ T=NETWORK_RESPONSE
    │ S=IMPS_DEBIT(COMPLETED / FAILED)
    │ L=AUTH(HOLD)
    ▼
┌────────────────────────────┐
│ Transaction Manager (TM)   │
└───┬────────────────────────┘
    │
    ├──────── SUCCESS ─────────────────────────────────┐
    │                                                   │
    │ PRODUCE                                          │ PRODUCE
    │ account.commands.final-debit                     │ account.commands.release-hold
    │                                                   │
    │ T=DEBIT_REQUESTED                                │ T=RELEASE_REQUESTED
    │ S=FINAL_DEBIT(IN_PROGRESS)                       │ S=RELEASE(IN_PROGRESS)
    │ L=AUTH(HOLD)                                     │ L=AUTH(HOLD)
    ▼                                                   ▼
┌────────────────────────────┐         ┌────────────────────────────┐
│ Account Service            │         │ Account Service            │
│ FINAL DEBIT                │         │ RELEASE HOLD               │
└───┬────────────────────────┘         └───┬────────────────────────┘
    │                                           │
    │ PRODUCE                                  │ PRODUCE
    │ account.events.balance-debited           │ account.events.balance-released
    │                                           │
    │ T=COMPLETED                               │ T=FAILED
    │ S=FINAL_DEBIT(COMPLETED)                  │ S=RELEASE(COMPLETED)
    │ L=DEBIT + SETTLEMENT                      │ L=RELEASE
    ▼                                           ▼
┌────────────────────────────┐         ┌────────────────────────────┐
│ Ledger Service / TM Ledger │         │ Ledger Service / TM Ledger │
│ DEBIT + SETTLEMENT ENTRY   │         │ RELEASE ENTRY              │
└──────────────┬─────────────┘         └──────────────┬─────────────┘
               │                                       │
               ▼                                       ▼
        ┌──────────────────────────────────────────────┐
        │ TM PRODUCE FINAL EVENT                        │
        │ imps.events.transaction-final                 │
        │                                               │
        │ T=COMPLETED / FAILED                          │
        │ S=COMPLETED / FAILED                          │
        │ L=FINALIZED                                   │
        └──────────────────────────────────────────────┘
