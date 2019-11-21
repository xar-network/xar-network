import React, { Component, ReactNode } from 'react';
import {RouteComponentProps, withRouter} from "react-router";
import { connect } from 'react-redux';
import Button from '../ui/Button';
import {ThunkDispatch} from "redux-thunk";
import {REDUX_STATE} from "../../ducks";
import {ActionType} from "../../ducks/types";
import {logout} from "../../ducks/user";
import './wallet-header.scss';
import {HOME} from "../../constants/routes";

type DispatchProps = {
  logout: () => boolean
}

type PropTypes = DispatchProps & RouteComponentProps;

class WalletHeader extends Component<PropTypes> {
  state = {
    isLoggingOut: false,
  };

  login = async () => {
    const resp = await this.props.logout();

    this.setState({
      isLoggingOut: false,
    });
    this.props.history.push(HOME);
  };

  render() {

    const { isLoggingOut } = this.state

    return (
      <div className="wallet-header">
        <Button
          type="primary"
          onClick={this.login}
          disabled={isLoggingOut}
          loading={isLoggingOut}
        >
          { isLoggingOut ? 'Logging Out' : 'Logout' }
        </Button>
      </div>
    )
  }
}

function mapDispatchToProps (dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>): DispatchProps {
  return {
    logout: () => dispatch(logout()),
  }
}


export default withRouter(
  connect(null, mapDispatchToProps)(WalletHeader)
);
