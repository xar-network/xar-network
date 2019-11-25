import React, { Component } from 'react';
import { connect } from 'react-redux';
import {REDUX_STATE} from "../../ducks";
import {estimateBatch, reduceDepthFromOrders, sortOrders} from "../../utils/exchange-utils";
import * as am4charts from "@amcharts/amcharts4/charts";
import * as am4core from "@amcharts/amcharts4/core";
import { BigNumber as BN } from "bignumber.js";
import {bn} from "../../utils/bn";
import {Order} from "../../ducks/exchange";

type DepthChartData = {
  price: string
  value: string
}

type Props = {
  bids: Order[]
  asks: Order[]
  quoteDenom: string
  baseDenom: string
  clearingPrice: BN
  bidRation: BN
  askRation: BN
}

type TooltipType = {
  fontSize: Number
  background: {
    fill: am4core.Color,
    fillOpacity: Number,
  }
}

const COLOR_WHITE = am4core.color('#fff');
// const COLOR_BLUE = am4core.color('#3084c3');
const COLOR_GREEN = am4core.color('#53b987');
const COLOR_RED = am4core.color('#eb4d5c');

class DepthChart extends Component<Props> {
  chart?: am4charts.XYChart;
  buySeries?: am4charts.StepLineSeries;
  sellSeries?: am4charts.StepLineSeries;
  // clearingSeries?: am4charts.LineSeries;
  priceAxis?: am4charts.ValueAxis;
  amountAxis?: am4charts.ValueAxis;

  readyChart() {
    if (!this.chart) {
      const chart = am4core.create('depthchart-wrapper', am4charts.XYChart);
      const priceAxis = chart.xAxes.push(new am4charts.ValueAxis());
      const amountAxis = chart.yAxes.push(new am4charts.ValueAxis());
      const priceTooltip = priceAxis.tooltip;
      const amountTooltip = amountAxis.tooltip;
      const buySeries = chart.series.push(new am4charts.StepLineSeries());
      const sellSeries = chart.series.push(new am4charts.StepLineSeries());
      // const clearingSeries = chart.series.push(new am4charts.LineSeries());
      const cursor = new am4charts.XYCursor();

      chart.cursor = cursor;
      // Configure Tooltip on X Axis
      (priceTooltip as TooltipType).background.fill = COLOR_WHITE;
      (priceTooltip as TooltipType).background.fillOpacity = .2;
      (priceTooltip as TooltipType).fontSize = 12;

      // Configure Tooltip on Y Axis
      (amountTooltip as TooltipType).background.fill = COLOR_WHITE;
      (amountTooltip as TooltipType).background.fillOpacity = .2;
      (amountTooltip as TooltipType).fontSize = 12;

      // Configure Cursor
      cursor.lineX.stroke = COLOR_WHITE;
      cursor.lineX.strokeWidth = 1;
      cursor.lineX.strokeOpacity = 0.5;
      cursor.lineY.stroke = COLOR_WHITE;
      cursor.lineY.strokeWidth = 1;
      cursor.lineY.strokeOpacity = 0.5;

      // clearingSeries.strokeWidth = 1;
      // clearingSeries.stroke = COLOR_BLUE;
      // clearingSeries.strokeDasharray = 'dashed';

      // Configure Y Axis
      amountAxis.renderer.grid.template.strokeOpacity = .05;
      amountAxis.renderer.grid.template.stroke = COLOR_WHITE;
      amountAxis.renderer.grid.template.strokeWidth = 1;
      amountAxis.renderer.labels.template.fill = COLOR_WHITE;
      amountAxis.renderer.labels.template.opacity = .5;
      amountAxis.renderer.labels.template.fontSize = 10;
      amountAxis.renderer.labels.template.fontFamily = 'Roboto';
      amountAxis.renderer.labels.template.fontWeight = '300';
      amountAxis.min = 0;
      priceAxis.strictMinMax = true;

      // Configure X Axis
      priceAxis.renderer.minGridDistance = 50;
      priceAxis.renderer.grid.template.strokeOpacity = .05;
      priceAxis.renderer.grid.template.stroke = COLOR_WHITE;
      priceAxis.renderer.grid.template.strokeWidth = 1;
      // priceAxis.renderer.grid.template.location = 0;
      priceAxis.renderer.labels.template.fill = COLOR_WHITE;
      priceAxis.renderer.labels.template.opacity = .5;
      priceAxis.renderer.labels.template.fontSize = 10;
      priceAxis.renderer.labels.template.fontFamily = 'Roboto';
      priceAxis.renderer.labels.template.fontWeight = '300';
      priceAxis.strictMinMax = true;

      // chart.padding(0, 0, 0, 16);

      this.chart = chart;
      this.buySeries = buySeries;
      this.sellSeries = sellSeries;
      this.priceAxis = priceAxis;
      this.amountAxis = amountAxis;
      // this.clearingSeries = clearingSeries;
    }
  }

  hydrateChartWithData() {
    this.readyChart();
    if (!this.chart || !this.buySeries || !this.sellSeries || !this.priceAxis || !this.amountAxis) {
      return;
    }

    const bids: DepthChartData[] = [];
    const asks: DepthChartData[] = [];

    let accumBids = bn(0);
    let accumAsks = bn(0);

    const sortedBids = sortOrders(this.props.bids, true);

    sortedBids.forEach(b => {
      const valueN =  b.price
        // .div(10 ** 8)
        .multipliedBy(b.quantity.div(10 ** 8));

      bids.push({
        price: b.price
          // .div(10 ** 8)
          .toFixed(8),
        value: valueN
          .plus(accumBids)
          .toFixed(8),
      });
      accumBids = valueN.plus(accumBids);
    });

    this.buySeries.data = bids;
    this.buySeries.strokeWidth = 2;
    this.buySeries.stroke = COLOR_GREEN;
    this.buySeries.fill = this.buySeries.stroke;
    this.buySeries.fillOpacity = 0.1;
    this.buySeries.dataFields.valueX = 'price';
    this.buySeries.dataFields.valueY= 'value';


    const sortedAsks = sortOrders(this.props.asks);

    sortedAsks.forEach(b => {
      const valueN =  b.price
        // .div(10 ** 8)
        .multipliedBy(b.quantity.div(10 ** 8));

      asks.push({
        price: b.price
          // .div(10 ** 8)
          .toFixed(8),
        value: valueN.plus(accumAsks).toFixed(8),
      });
      accumAsks = valueN.plus(accumAsks);
    });

    this.sellSeries.data = asks;
    this.sellSeries.strokeWidth = 2;
    this.sellSeries.stroke = COLOR_RED;
    this.sellSeries.fill = this.sellSeries.stroke;
    this.sellSeries.fillOpacity = 0.1;
    this.sellSeries.dataFields.valueX = 'price';
    this.sellSeries.dataFields.valueY= 'value';

    this.chart.invalidateRawData();
  }

  componentDidMount() {
    this.hydrateChartWithData();
  }

  componentDidUpdate() {
    this.hydrateChartWithData();
  }

  componentWillUnmount() {
    if (this.chart) {
      this.chart.dispose();
    }
  }

  render() {
    return (
      <div style={{ height: '100%', width: '100%' }} id="depthchart-wrapper" />
    );
  }
}

function mapStateToProps(state: REDUX_STATE): Props {
  const {
    exchange: { selectedMarket, markets },
  } = state;
  const market = markets[selectedMarket] || {};
  const { bids = [], asks = [], quoteDenom, baseDenom } = market;
  const quoteDecimals = 8 || 0;
  const { depths: bidDepth } = reduceDepthFromOrders(bids, 8, 8);
  const { depths: askDepth } = reduceDepthFromOrders(asks, 8, 8);
  const { clearingPrice, bidRation, askRation } = estimateBatch(bids, asks, 8, 8);

  return {
    bids: bidDepth,
    asks: askDepth,
    quoteDenom,
    baseDenom,
    clearingPrice: clearingPrice.div(10 ** quoteDecimals),
    bidRation,
    askRation,
  }
}

export default connect(mapStateToProps)(DepthChart);
