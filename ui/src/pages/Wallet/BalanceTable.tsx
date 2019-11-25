import React, { Component, ReactNode } from 'react';
import {connect} from "react-redux";
import {REDUX_STATE} from "../../ducks";
import {Dispatch} from "redux";
import {Table, TableCell, TableHeader, TableHeaderRow, TableRow} from "../../components/ui/Table";
import {BalanceType} from "../../ducks/user";
import "./style/balance-table.scss";

type StateProps = {
  balances: {
    [denom: string]: BalanceType
  }
}

type DispatchProps = {

}

type Props = StateProps & DispatchProps

class BalanceTable extends Component<Props> {
  render () {
    return (
      <Table className="wallet__balance-table">
        { this.renderHeaderRow() }
        { this.renderTableBody() }
      </Table>
    );
  }

  renderHeaderRow (): ReactNode {
    return (
      <TableHeaderRow>
        <TableHeader>Denom</TableHeader>
        <TableHeader>Balance</TableHeader>
        <TableHeader>Available</TableHeader>
        <TableHeader>On Orders</TableHeader>
      </TableHeaderRow>
    );
  }

  renderTableBody (): ReactNode {
    return (
      <div className="wallet__content__table__body">
        {
          Object.entries(this.props.balances)
            .map(([_, balance]) => this.renderTableRow(balance))
        }
      </div>
    );
  }

  renderTableRow (balance: BalanceType): ReactNode {
    const { denom, locked, unlocked } = balance;

    return (
      <TableRow key={denom}>
        <TableCell>{denom}</TableCell>
        <TableCell>{unlocked.toFixed(4)}</TableCell>
        <TableCell>{unlocked.toFixed(4)}</TableCell>
        <TableCell>{locked.toFixed(4)}</TableCell>
      </TableRow>
    )
  }
}

function mapStateToProps(state: REDUX_STATE): StateProps {
  return {
    balances: state.user.balances
  };
}

function mapDispatchToProps(dispatch: Dispatch): DispatchProps {
  return {

  };
}

export default connect(mapStateToProps, mapDispatchToProps)(BalanceTable);
