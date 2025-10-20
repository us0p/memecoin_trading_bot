## 3rd party providers
### Helius
- Rate limiting: 10 requests/s for RPC endpoints.
- 1 milion credits/month.
- each RPC call costs 1 credit.

## TODO
- Add dead token identification.
- Must add thorough tests.
- Errors not attached to a specific token don't use the 'mint' field.
- Run notification system in a separate thread.
- should think of a way to clean intermitent errors when they fix themselves, to avoid error clustering.

- FIX ERROR REPORT, AFTER ADDING MUTEX, NOTIFICATIONS AREN'T BEING RECEIVED.
- IMPROVE SELL ORDER LOCKING, USING DATABASE MIGHT INCREASE LATENCY.
- ADD SELL SIMULATION FOR WALLET HOLDING, AFTER SIMULATION BUY, WHEN TRYING TO SELL THE AMMOUNT IN THE WALLET IS ALWAYS GOING TO BE 0.
