# CoinEX
this package is not stable.

All-in-One Cryptocurrency Trading Bot

coinex is an implementation of the Cryptocurrency Exchange(Poloniex/Bithumb/Coinone/Bittrex and more) API (public and private) in Golang. 
The current implementation is very early and can not be used yet.

The following exchanges will be supported.
1. Poloniex
2. Bithumb
3. Coinone
4. Bittrex

![img](http://i.imgur.com/lNP9hgc.png)

## Dependencies  
 - https://github.com/influxdata/influxdb/client/v2
 - https://github.com/IronPark/go-poloniex
 - https://github.com/toorop/go-bittrex
## Donation
My Ethereum Wallet :
0x7EA84eFF0f9D3bd2EaD6Db190A4387B71ac42b44

## Roadmap
Goals for **CoinEX**
  
- [ ] TravisCI for this package.

- [ ] Trade simulation for algorithmic trading
  - [x] Store trade history in time series database (influxDB)
  - [ ] Alpha Model
  - [ ] Output of trade simulation report
  - [ ] Support script lang
  
- [ ] Web-base visualization
## References
- apis
    - [poloniex](https://poloniex.com/support/api/)
    - [bittrex](https://bittrex.com/Home/Api)
    - [bithumb](https://www.bithumb.com/u1/US127)
    - [poloniex](https://poloniex.com/support/api/)
    - [coinone](http://doc.coinone.co.kr)
- codes
    - [wss gist code](https://gist.github.com/ismasan/3fb75381cd2deb6bfa9c)
    - [go-poloniex](https://github.com/jyap808/go-poloniex)
    