import React, {Component, ReactNode} from 'react';
import './style/side-panel.scss';
import Button from "../../components/ui/Button";
import {RouteComponentProps, withRouter} from "react-router";
import {DEPOSIT, WITHDRAWAL} from "../../constants/routes";
import {REDUX_STATE} from "../../ducks";
import {connect} from "react-redux";
import {bn} from "../../utils/bn";
import BigNumber from "bignumber.js";

type StateProps = {
  locked: BigNumber
  unlocked: BigNumber
  denom: String
}

type DispatchProps = {

}

type Props = StateProps & DispatchProps & RouteComponentProps;

class SidePanel extends Component<Props> {
  render() {
    const { history, locked, unlocked, denom } = this.props;
    return (
      <div className="wallet__side-panel">
        <div className="wallet__side-panel__header">
          <div className="wallet__side-panel__header__text">Overview</div>
          <div className="wallet__side-panel__header__conversion-dropdown">{denom}</div>
        </div>
        { this.renderBalance('Total Holdings', locked.plus(unlocked).toFixed(4))}
        { this.renderBalance('On Orders', locked.toFixed(4))}
        { this.renderBalance('Available Balance', unlocked.toFixed(4))}
      </div>
    )
  }

  renderBalance(label: string, value: string): ReactNode {
    return (
      <div className="wallet__side-panel__balance-group">
        <div className="wallet__side-panel__balance-group__label">{label}</div>
        <div className="wallet__side-panel__balance-group__value">{value}</div>
      </div>
    )
  }
}

function mapStateToProps (state: REDUX_STATE): StateProps {
  const { balances } = state.user;
  const balance = balances['uzar'] || {};
  return {
    locked: balance.locked || bn(0),
    unlocked: balance.unlocked || bn(0),
    denom: balance.denom || '',
  }
}

export default withRouter(
  connect(mapStateToProps)(SidePanel)
);
