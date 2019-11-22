import React, { Component, ReactNode } from 'react';
import { connect } from 'react-redux';
import { REDUX_STATE } from '../../ducks';
import './exchange-header.scss';
import Numeral from "../Numeral";
import {AssetType} from "../../ducks/assets";
import {DayStatsType, MarketType} from "../../ducks/exchange";
import Dropdown, { ItemType } from '../ui/Dropdown'

type StatePropTypes = {
  baseAsset?: AssetType
  quoteAsset?: AssetType
  dayStats?: DayStatsType
  markets?: any
};

type DispatchPropTypes = {}

type PropTypes = StatePropTypes & DispatchPropTypes;

class ExchangeHeader extends Component<PropTypes> {
  render() {
    const {
      quoteAsset,
      baseAsset,
      dayStats,
      markets,
    } = this.props;

    if (!quoteAsset || !baseAsset || !dayStats || !markets) return <div />;

    const quoteSymbol = quoteAsset.symbol;
    const baseSymbol = baseAsset.symbol;

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
        { this.renderItem('Trading Pair', this.renderMarketsSelect(markets) /*`${baseSymbol}/${quoteSymbol}`*/) }
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
        label: markets[index].baseSymbol + '/' + markets[index].quoteSymbol,
        value: markets[index].baseSymbol + '/' + markets[index].quoteSymbol
      })
      index++
    }

    marketsArray.push({
      label: 'UUSD/UBTC',
      value: 'UUSD/UBTC'
    })
    marketsArray.push({
      label: 'UUSD/UFTM',
      value: 'UUSD/UFTM'
    })
    marketsArray.push({
      label: 'UUSD/UCSDT',
      value: 'UUSD/UCSDT'
    })
    marketsArray.push({
      label: 'UUSD/UETH',
      value: 'UUSD/UETH'
    })

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
            <div onClick={item.toggleDropdown} >
              <div className="exchange-header__item__text">{item.label}</div>
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
    assets: { assets, symbolToAssetId },
  } = state;

  const market = markets[selectedMarket] || { dayStats: {} };
  const baseAsset = assets[symbolToAssetId[market.baseSymbol]];
  const quoteAsset = assets[symbolToAssetId[market.quoteSymbol]];

  return {
    dayStats: market.dayStats,
    baseAsset,
    quoteAsset,
    markets,
  }
}

export default connect(mapStateToProps, null)(ExchangeHeader)
