## 3rd party providers
### Helius
- Rate limiting: 10 requests/s for RPC endpoints.
- 1 milion credits/month.
- each RPC call costs 1 credit.

## TODO
- Create SOL USD price web socket consumer.
- Wil determine the execution price of trades by performing the following math:
    - (Order creation) TokenPerSOL = OutputAmount / InputAmount (represents the amount of tokens per unit of SOL)
    - SOL price (retrieved from global state updated by WebSocket) / TokenPerSOL
