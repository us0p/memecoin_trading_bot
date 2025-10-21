## 3rd party providers
### Helius
- Rate limiting: 10 requests/s for RPC endpoints.
- 1 milion credits/month.
- each RPC call costs 1 credit.

## TODO
- Must add thorough tests.
- Errors not attached to a specific token don't use the 'mint' field.
- Run notification system in a separate thread.
- should think of a way to clean intermitent errors when they fix themselves, to avoid error clustering.
- Add dead token identification. (need more data).
