## 3rd party providers
### Helius
- Rate limiting: 10 requests/s for RPC endpoints.
- 1 milion credits/month.
- each RPC call costs 1 credit.

## TODO
- Re-think about trade storing, maybe storing separate operations instead 
  of updating the trade entry might be more clean and easier to implement.
- Must also rethinkg structure of the table to store the most relevant data
  that we receive from the provider.
    - Must add buy fee and sell fee to trade entries in the database.
- Create job scheduler.
- Determine how trades are going to be executed in the current workflow.
- Must add thorough tests.
- Errors not attached to a specific token don't use the 'mint' field.
