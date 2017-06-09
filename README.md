# CoinEX
All-in-One Cryptocurrency Exchange API

coinex is an implementation of the Cryptocurrency Exchange(Poloniex/Bithumb/Coinone/Bittrex and more) API (public and private) in Golang. 
The current implementation is very early and can not be used yet.

The following exchanges will be supported.
1. Poloniex
2. Bithumb
3. Coinone
4. Bittrex

![img](http://i.imgur.com/lNP9hgc.png)

## Dependencies  
 - github.com/influxdata/influxdb/client/v2
 
## Donation
My Ethereum Wallet :
0x7EA84eFF0f9D3bd2EaD6Db190A4387B71ac42b44

## Roadmap
Goals for **CoinEX**

- [ ] Public API implements All exchanges
  - [ ] bithumb
  - [ ] poloniex
  - [ ] bittrex
  - [ ] coinone
  
- [ ] Trade API implements All exchanges
  - [ ] bithumb
  - [ ] poloniex
  - [ ] bittrex
  - [ ] coinone
  
- [ ] TravisCI for this package.

- [ ] Trade simulation for algorithmic trading
  - [x] Store trade history in time series database (influxDB)
  - [ ] Alpha Model
  - [ ] Output of trade simulation report
  - [ ] Support script lang
  
- [ ] Web-base visualization
## References
- [bithumb](https://www.bithumb.com/u1/US127)
- [poloniex](https://poloniex.com/support/api/)
- [coinone](http://doc.coinone.co.kr)