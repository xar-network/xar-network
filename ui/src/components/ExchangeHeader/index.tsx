import React, { Component, ReactNode } from 'react';
import { connect } from 'react-redux';
import { REDUX_STATE } from '../../ducks';
import './exchange-header.scss';
import Numeral from "../Numeral";
import {DayStatsType, MarketType} from "../../ducks/exchange";
import Dropdown, { ItemType } from '../ui/Dropdown'

type StatePropTypes = {
  baseDenom?: string
  quoteDenom?: string
  dayStats?: DayStatsType
  markets?: any
};

type DispatchPropTypes = {}

type PropTypes = StatePropTypes & DispatchPropTypes;

class ExchangeHeader extends Component<PropTypes> {
  render() {
    const {
      quoteDenom,
      baseDenom,
      dayStats,
      markets,
    } = this.props;

    if (!quoteDenom || !baseDenom || !dayStats || !markets) return <div />;

    const {
      lastPrice,
      prevPrice,
      dayChange,
      dayLow,
      dayHigh,
      dayChangePercentage,
    } = dayStats;

    const changePercentage = lastPrice.minus(prevPrice).div(prevPrice).toNumber();
    let lastPriceModifier = '';

    if (changePercentage < 0) {
      lastPriceModifier = 'negative';
    } else if (changePercentage > 0) {
      lastPriceModifier = 'positive';
    }

    return (
      <div className="exchange-header">
        { this.renderItem('Trading Pair', this.renderMarketsSelect(markets)) }
        {
          this.renderItem(
            'Last Price',
            <Numeral
              value={lastPrice}
              decimals={8}
              displayDecimals={8}
              formatAsCurrency
            />,
            lastPriceModifier,
            lastPrice.isZero(),
          )
        }
        {
          this.renderItem(
            '24H Change',
            <div>
              <Numeral
                value={dayChange}
                decimals={8}
                displayDecimals={8}
                formatAsCurrency
              />
              <span>(</span>
              <span>{(dayChangePercentage * 100).toFixed(2)}</span>
              <span>%)</span>
            </div>,
            dayChange.isPositive() ? 'positive' : 'negative',
            dayChange.isZero(),
          )
        }
        {
          this.renderItem(
            '24H High',
            <Numeral
              value={dayHigh}
              decimals={8}
              displayDecimals={8}
              formatAsCurrency
            />,
            '',
            dayHigh.isZero(),
          )
        }
        {
          this.renderItem(
            '24H Low',
            <Numeral
              value={dayLow}
              decimals={8}
              displayDecimals={8}
              formatAsCurrency
            />,
            '',
            dayLow.isZero(),
          )
        }
      </div>
    )
  }

  renderMarketsSelect(markets: any) {
    if(!markets) {
      return null
    }

    const marketsArray = []
    let index = 1

    while(markets[index] != null) {
      marketsArray.push({
        label: markets[index].baseDenom + '/' + markets[index].quoteDenom,
        value: markets[index].baseDenom + '/' + markets[index].quoteDenom
      })
      index++
    }

    return (
      <Dropdown
        items={marketsArray}
        currentIndex={0}
        CurrentItem={(item: ItemType): ReactNode => (
          <div onClick={item.toggleDropdown}>
            <div className="exchange-header__item__caret" />
            <div className="exchange-header__item__text">{item.label}</div>
          </div>
        )}
        Item={(item: ItemType): ReactNode => {
          return (
            <div key={item.label} onClick={item.toggleDropdown} >
              <div  className="exchange-header__item__text">{item.label}</div>
            </div>
          );
        }}
      />
    )
  }

  renderItem(label: string, value: ReactNode, modifier: string = '', isLoading: boolean = false): ReactNode {
    if (isLoading) {
     return (
       <div className={`exchange-header__item exchange-header__item--loading`}>
         <div className="exchange-header__item__label">{label}</div>
         <div className="exchange-header__item__value" />
       </div>
     );
    }

    return (
      <div className={`exchange-header__item exchange-header__item--${modifier}`}>
        <div className="exchange-header__item__label">{label}</div>
        <div className="exchange-header__item__value">{value}</div>
      </div>
    )
  }
}

function mapStateToProps(state: REDUX_STATE) {
  const {
    exchange: {
      selectedMarket,
      markets,
    },
  } = state;

  const market = markets[selectedMarket] || { dayStats: {} };
  const baseDenom = market.baseDenom;
  const quoteDenom = market.quoteDenom

  return {
    dayStats: market.dayStats,
    baseDenom,
    quoteDenom,
    markets,
  }
}

export default connect(mapStateToProps, null)(ExchangeHeader)
