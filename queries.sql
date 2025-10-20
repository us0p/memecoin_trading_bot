-- Trending data analysis
with hourly_movimentation as (
	select (row_number() over() / 12) / 60 hour,
	num_buys_1h + num_sells_1h + num_traders_1h trading_activity,
	buy_volume_1h + sell_volume_1h volume,
	buy_organic_volume_1h / sell_organic_volume_1h organic_volume,
	num_net_buyers_1h,
	buy_volume_1h,
	sell_volume_1h
	from market_data 
	order by priced_at
) select 
	hour, 
	cast(avg(trading_activity) as int) avg_trading_activity,
	cast(avg(volume) as int) avg_volume,
	cast(avg(organic_volume) as int) avg_organic_volume,
	cast(avg(num_net_buyers_1h) as int) avg_net_buyers,
	cast(avg(buy_volume_1h) as int) avg_buy_volume,
	cast(avg(sell_volume_1h) as int) avg_sell_volume
from hourly_movimentation 
group by hour;
