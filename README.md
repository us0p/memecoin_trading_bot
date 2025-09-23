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
    - This determines the issued usd price for the token
    - Might also take the USD price for the last MK Data entry for the token in the database.

    - (Order execution) If OutputAmount (creation) != OutputAmount (exec) then there's a price diff.
    - (TokenPerSOL = OutputAmount / InputAmount) / SOL price
    - Determines the executed price for the token
