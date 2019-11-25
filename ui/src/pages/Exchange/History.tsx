import React, {Component} from 'react';
import {connect} from 'react-redux';
import {ThunkDispatch} from 'redux-thunk';
import {Table, TableCell, TableHeader, TableHeaderRow, TableRow} from '../../components/ui/Table';
import {Module, ModuleContent, ModuleHeader, ModuleHeaderButton} from '../../components/Module';
import {REDUX_STATE} from '../../ducks';
import {OrderType} from '../../ducks/exchange';
import {ActionType} from '../../ducks/types';
import {fetchUserOrders, ORDER_HISTORY_FILTERS, setOrderHistoryFilter} from '../../ducks/user';

type StateProps = {
  orders: OrderType[]
  orderHistoryFilter: ORDER_HISTORY_FILTERS
  baseDenom: string
  quoteDenom: string
}

type DispatchProps = {
  fetchUserOrders: () => void
  setOrderHistoryFilter: (f: ORDER_HISTORY_FILTERS) => ActionType<ORDER_HISTORY_FILTERS>
}

type Props = StateProps & DispatchProps

class History extends Component<Props> {
  componentWillMount() {
    this.props.fetchUserOrders();
  }

  render() {
    const {
      setOrderHistoryFilter,
      orderHistoryFilter,
      orders,
      baseDenom,
      quoteDenom,
    } = this.props;
    const filteredOrders = orders/*.filter(({status}) => {
      switch (orderHistoryFilter) {
        case ORDER_HISTORY_FILTERS.ALL:
          return true;
        case ORDER_HISTORY_FILTERS.OPEN:
          return status === 'OPEN';
        default:
          return false;
      }
    });*/

    if (!baseDenom || !quoteDenom) return <noscript />;

    return (
      <Module className="exchange__history">
        <ModuleHeader>
          <ModuleHeaderButton
            onClick={() => setOrderHistoryFilter(ORDER_HISTORY_FILTERS.OPEN)}
            active={orderHistoryFilter === ORDER_HISTORY_FILTERS.OPEN}
          >
            Open
          </ModuleHeaderButton>
          <ModuleHeaderButton
            onClick={() => setOrderHistoryFilter(ORDER_HISTORY_FILTERS.ALL)}
            active={orderHistoryFilter === ORDER_HISTORY_FILTERS.ALL}
          >
            All
          </ModuleHeaderButton>
        </ModuleHeader>
        <ModuleContent>
          <Table className="exchange__history__table">
            <TableHeaderRow>
              <TableHeader>Block</TableHeader>
              <TableHeader>Pair</TableHeader>
              <TableHeader>Type</TableHeader>
              <TableHeader>Side</TableHeader>
              <TableHeader>{`Price (${quoteDenom})`}</TableHeader>
              <TableHeader>{`Amount (${baseDenom})`}</TableHeader>
              <TableHeader>Filled</TableHeader>
              <TableHeader>Status</TableHeader>
            </TableHeaderRow>
            <div className="exchange__history__table-content">
              {
                filteredOrders.length
                  ?
                  filteredOrders.map(this.renderRow)
                  :
                  <TableRow>
                    <TableCell>No orders to display.</TableCell>
                  </TableRow>
              }
            </div>
          </Table>
        </ModuleContent>
      </Module>
    )
  }

  renderRow = (order: OrderType): React.ReactNode => {
    const { baseDenom, quoteDenom } = this.props;

    if (!baseDenom || !quoteDenom || !order) return <noscript />;

    return (
      <TableRow key={order.id}>
        <TableCell>{ order.created_block }</TableCell>
        <TableCell>{`${baseDenom}/${quoteDenom}`}</TableCell>
        <TableCell>LIMIT</TableCell>
        <TableCell>{ order.direction }</TableCell>
        <TableCell>
          {
            order.price
              .div(10 ** 8)
              .toFixed(Math.min(6, 8))
          }
        </TableCell>
        <TableCell>
          {
            order.quantity
              .div(10 ** 8)
              .toFixed(Math.min(6, 8))
          }
        </TableCell>
        <TableCell>
          {
            order.quantity_filled
              .div(10 ** 8)
              .toFixed(Math.min(6, 8))
          }
        </TableCell>
        <TableCell>{ order.status }</TableCell>
      </TableRow>
    );
  }
}

function mapStateToProps(state: REDUX_STATE): StateProps {
  const {
    user: {
      orders: history,
      orderHistoryFilter,
    },
    exchange: { orders, selectedMarket, markets },
  } = state;
  const market = markets[selectedMarket] || {};
  const { baseDenom, quoteDenom } = market;
  return {
    orders: history.map(id => orders[id]),
    orderHistoryFilter,
    quoteDenom,
    baseDenom,
  }
}

function mapDispatchToProps(dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>): DispatchProps {
  return {
    fetchUserOrders: () => dispatch(fetchUserOrders()),
    setOrderHistoryFilter: (filter: ORDER_HISTORY_FILTERS) => dispatch(setOrderHistoryFilter(filter)),
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(History);
