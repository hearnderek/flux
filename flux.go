package main

import (
  "fmt"
  "time"
  "math/rand"
  "sort"
)

func main() {
  rand_init()
  // loop()
  test_many()
}

func test_many () {
  days := 364
  starting_market := 14000.00
  starting_investment := 0.00
  daily := 20.00

  max_div := 0.00
  max_val := 0.00

  min_div := 1000000.00
  min_val := 1000000.00

  div := 0.00
  var divs []float64
  val := 0.00
  var vals []float64

  i := 0
  times := 10000
  for i < times {
    i++
    s := test(days, starting_market, starting_investment, daily)
    if s.dividends > max_div {
      max_div = s.dividends
    }
    if s.dividends < min_div {
      min_div = s.dividends
    }
    div += s.dividends

    if s.value > max_val {
      max_val = s.value
    }
    if s.value < min_val {
      min_val = s.value
    }
    val += s.value
    divs = append(divs, s.dividends)
    vals = append(vals, s.value)
  }

  sort.Float64s(divs)
  sort.Float64s(vals)
  div_mean := divs[len(divs)/2]
  val_mean := vals[len(vals)/2]
  total_invested := starting_investment + float64(days) * daily

  fmt.Println("Div:","Max:",max_div,"Min:",min_div,"Avg:",div/float64(times),"Mean:",div_mean)
  fmt.Println("Val:","Max:",max_val,"Min:",min_val,"Avg:",val/float64(times),"Mean:",val_mean)
  fmt.Println("total Invested:",total_invested)
  fmt.Println("%return:",div_mean / total_invested)
  fmt.Println("%return:",(val_mean + div_mean - total_invested) / total_invested)
  fmt.Println("Div per month:",div_mean / 12)
}

func loop() {
  count := 0
  index := 17000.00
  cash := 100.00
  p := new_portfolio(cash)

  for count < 365 {
    index = pull_market(index)
    count++

    // fmt.Println(count,"BTC:",index)

    p = p.act(index)

    p.pp(index)

    time.Sleep(100 * time.Millisecond)
    p = p.invest(100)
  }
  p.pp(index)
}

func test(days int, starting_index float64, starting_cash float64, invest_daily float64) Scenario {
  count := 0
  index := starting_index
  p := new_portfolio(starting_cash)

  for count < days {
    index = pull_market(index)
    count++

    p = p.act(index)

    p = p.invest(invest_daily)
  }
  return Scenario {p.dividends, p.invested, p.value(index), index}
}

//----//

type Scenario struct {
  dividends float64
  invested float64
  value float64
  market_value float64
}

//----//

type Sec struct {
  strike float64
  invested float64
  btc float64
}

func (purchase Sec) sell (market_value float64) (cash float64) {
  return purchase.btc * market_value
}

func (purchase Sec) should_sell (market_value float64) bool {
  return (market_value - purchase.strike) / purchase.strike >= 0.03
}

//----//

type Portfolio struct {
  cash float64
  purchases []Sec
  invested float64
  dividends float64
}


func new_portfolio (inital_capital float64) Portfolio {
  var purchases []Sec
  portfolio := Portfolio { inital_capital, purchases, inital_capital, 0 }
  return portfolio
}

func (p Portfolio) value (market_value float64) float64 {
  total := p.cash

  for _, s := range p.purchases {
    // fmt.Println("s:",s)
    total += value(market_value, s.btc)
  }

  return total
}

func (p Portfolio) invest (usd float64) Portfolio {
  p.invested += usd
  p.cash += usd
  return p
}

func (p Portfolio) buy (market_value float64, usd float64) Portfolio {
  p.cash -= usd
  sec := buy(market_value, usd)
  p.purchases = append(p.purchases, sec)
  return p;
}

func (p Portfolio) act (market_value float64) Portfolio {
  var updated_purchases []Sec

  for _, purchase := range p.purchases {
    if(purchase.should_sell(market_value)) {
      received_cash := purchase.sell(market_value)
      dividend := received_cash * 0.0625

      p.cash += received_cash - dividend
      p.dividends += dividend
    } else {
      updated_purchases = append(updated_purchases, purchase)
    }
  }
  p.purchases = updated_purchases

  if(p.should_buy(market_value)) {
     p = p.buy_auto(market_value)
  }

  return p
}

func (p Portfolio) should_buy (market_value float64) bool {
  opinion := p.cash > 5 && market_value > 0
  if(opinion) {
    // fmt.Println("Opinion:BUY")
  }
  return opinion
}

func (p Portfolio) buy_auto (market_value float64) Portfolio {
  return p.buy(market_value, p.cash)
}

func (p Portfolio) pp (index float64) {
  // fmt.Println("Portfolio",p.purchases)
  fmt.Println("Portfolio Val",p.value(index))
  fmt.Println("Portfolio Div",p.dividends)
  fmt.Println("Portfolio Inv",p.invested)
}

//----//

func rand_init() {
  // Set so we have a new seed everytime
  rand.Seed(time.Now().UnixNano())
}

//----//

func pull_market (last_known float64) float64 {
  rval := rand.Intn(17000*0.08)

  if rand.Intn(2) == 0 {
    rval *= -1
  }
  new_market_value := last_known + float64(rval)
  if new_market_value < 0 {
    new_market_value = 0
  }
  return new_market_value
}

//----//

func buy(market_value float64, usd float64) (sec Sec) {
  return Sec{market_value, usd, usd/market_value}
}

func value(market_value float64, btc float64) float64 {
  return btc * market_value
}
