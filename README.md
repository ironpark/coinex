# CoinEx
[Korean](https://github.com/ironpark/coinex/blob/master/README-KR.md) | [English](https://github.com/ironpark/coinex/blob/master/README.md)

**CoinEX (temporary name)** is open source software that can **develop, test, share, execute and manage** trading strategies using various cryptographic exchange data.

The names and repositories will change as alpha test versions are developed.

**Current exchange** is as follows.
1. Bittrex / UPbit
2. Bitfinex
3. Poloniex
4. Coinone


## dependency
 - https://github.com/influxdata/influxdb/client/v2
 - https://github.com/ironpark/go-poloniex
 - https://github.com/toorop/go-bittrex
 - https://github.com/gin-gonic/gin
 - https://github.com/asaskevich/EventBus
 - https://github.com/sirupsen/logrus


## donate
Ethereum Address:

![ethereum address](https://chart.googleapis.com/chart?cht=qr&chl=%200x7EA84eFF0f9D3bd2EaD6Db190A4387B71ac42b44&chs=300x300&choe=UTF-8&chld=L|2)

0x7EA84eFF0f9D3bd2EaD6Db190A4387B71ac42b44

**Attention!** Do not send Ethereum Classic!

## Roadmap
Goals for **CoinEX**

- [ ] TravisCI support

- [ ] Simulation support for algorithmic trading
  - [x] Storing and managing real-time price data (OHCL) using a time series database (influxDB)
  - [ ] Plug-in system for custom trading strategy
  - [ ] Plug-in system for trading indicators ex) SMA, EMA, ...
  - [ ] Trading simulation using historical data and report support for performance evaluation of trading strategy
  - [ ] Scripting language support for plug-in ecosystem (python / js)

- Support visualization based on web technology

## References
- apis
    - [poloniex](https://poloniex.com/support/api/)
    - [bittrex](https://bittrex.com/Home/Api)
    - [coinone](http://doc.coinone.co.kr)
- codes
    - [go-poloniex](https://github.com/jyap808/go-poloniex)

## License
[MPL 2.0 (Mozilla Public License Version 2.0)](https://www.mozilla.org/en-US/MPL/2.0/)