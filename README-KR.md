[Korean](https://github.com/ironpark/coinex/blob/master/README-KR.md) | [English](https://github.com/ironpark/coinex/blob/master/README.md)

# CoinEx
**CoinEX(임시 가칭)** 는 다양한 암호화폐 거래소에서 거래 데이터를 실시간으로 받아와 로컬 데이터베이스에 저장 및 관리해주며 이를 기반으로 트레이딩 전략을 **개발**, **테스트**, **공유**, **관리**할 수 있는 오픈소스 플랫폼 입니다.

Alpha 테스트 버전이 개발 완료되는 시점에서 이름 및 저장소가 변경될 예정입니다.

**현재 지원 예정인 거래소**는 아래와 같습니다.
1. Bittrex / UPbit
2. Bitfinex
3. Poloniex
4. Coinone


## 의존성
 - https://github.com/influxdata/influxdb/client/v2
 - https://github.com/ironpark/go-poloniex
 - https://github.com/toorop/go-bittrex
 - https://github.com/gin-gonic/gin
 - https://github.com/asaskevich/EventBus
 - https://github.com/sirupsen/logrus

## 기부
이더리움(Ethereum) 주소 :

![ethereum address](https://chart.googleapis.com/chart?cht=qr&chl=%200x7EA84eFF0f9D3bd2EaD6Db190A4387B71ac42b44&chs=300x300&choe=UTF-8&chld=L|2')

0x7EA84eFF0f9D3bd2EaD6Db190A4387B71ac42b44

**주의!** 이더리움 클래식(Ethereum Classic)을 보내지 마세요!

## 로드맵
Goals for **CoinEX**

- [ ] TravisCI 지원

- [ ] 알고리드믹 트레이딩을 위한 시뮬레이션 지원
  - [x] 실시간 가격 데이터(OHCL)를 시계열 데이터베이스를 이용하여 보관/관리(influxDB)
  - [ ] 사용자 정의 거래 전략을 위한 플러그인 시스템
  - [ ] 거래 지표를 위한 플러그인 시스템 ex) SMA,EMA,...
  - [ ] 거래전략의 성능평가를 위한 거래 시뮬레이션과 보고서 지원
  - [ ] 플러그인 시스템 생태계를 위한 스크립트 언어지원 (python/js)

- [ ] 웹기술을 기반으로한 시각화 지원

## References
- apis
    - [poloniex](https://poloniex.com/support/api/)
    - [bittrex](https://bittrex.com/Home/Api)
    - [coinone](http://doc.coinone.co.kr)
- codes
    - [go-poloniex](https://github.com/jyap808/go-poloniex)

## License
[MPL 2.0 (Mozilla Public License Version 2.0)](https://www.mozilla.org/en-US/MPL/2.0/)