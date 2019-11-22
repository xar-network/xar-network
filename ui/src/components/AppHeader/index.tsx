import React, { Component } from 'react';
import SearchBox from '../SeachBox';
import ExchangeHeader from '../ExchangeHeader';
import WalletHeader from '../WalletHeader';
import NetworkDropdown from '../NetworkDropdown';
import {Route, Switch} from 'react-router';
import "./app-header.scss";
import {EXCHANGE, WALLET, AUTHENTICATE, HOME} from "../../constants/routes";
import Button from '../ui/Button';
import { connect } from 'react-redux';
import { REDUX_STATE } from "../../ducks";
import {ActionType} from "../../ducks/types";
import { withRouter, RouteComponentProps } from 'react-router-dom';
import {ThunkDispatch} from "redux-thunk";
import {logout} from "../../ducks/user";

type StateProps = {
  isLoggedIn?: boolean
}

type DispatchProps = {
  logout: () => boolean
}

type PropTypes = StateProps & DispatchProps & RouteComponentProps;

class AppHeader extends Component<PropTypes> {
  state = {
    isLoggingOut: false,
  };

  logout = async () => {
    const resp = await this.props.logout();

    this.setState({
      isLoggingOut: false,
    });

    this.props.history.push(HOME);
  }

  render() {
    return (
      <div className="app-header">
        <div className="app-header__content">
          <Switch>
            <Route path={EXCHANGE} component={ExchangeHeader} />
            <Route path={HOME} component={ExchangeHeader} />
          </Switch>
        </div>
        {this.props.isLoggedIn !== true && <div className="app-header__login">
          <Button
            type="primary"
            onClick={() => this.props.history.push(AUTHENTICATE)}
          >
            Login
          </Button>
        </div>}
        {this.props.isLoggedIn === true && <div className="app-header__login">
          <Button
            type="primary"
            onClick={this.logout}
          >
            Logout
          </Button>
        </div>}
      </div>
    )
  }
}

function mapDispatchToProps (dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>): DispatchProps {
  return {
    logout: () => dispatch(logout()),
  }
}

function mapStateToProps(state: REDUX_STATE): StateProps {
  return {
    isLoggedIn: state.user.isLoggedIn,
  };
}



export default withRouter(
  connect(mapStateToProps, mapDispatchToProps)(AppHeader)
);
