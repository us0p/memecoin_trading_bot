## 3rd party providers
### Helius
- Rate limiting: 10 requests/s for RPC endpoints.
- 1 milion credits/month.
- each RPC call costs 1 credit.

## TODO
- Create job scheduler.
- Determine how trades are going to be executed in the current workflow.
    - Job schedules is going to run wokflows based on the registered schedule.
        - Pull Tokens (1m)
            - If token is classified as a trade opportunity, put a CALL order on trade processing queue.
        - Market Data (5s)
            - If market data determines a trade should be closed, put a BID order on trade processing queue.
        - Token Largest Holders (1m)
- Must add thorough tests.
- Errors not attached to a specific token don't use the 'mint' field.
