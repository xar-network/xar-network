import React, {Component, ReactNode} from "react";
import {RouteComponentProps, withRouter} from "react-router";
import {connect} from "react-redux";
import Button from "../../components/ui/Button";
import "./authentication.scss";
import Input from "../../components/ui/Input";
import Radio from "../../components/ui/Radio";
import {ThunkDispatch} from "redux-thunk";
import {REDUX_STATE} from "../../ducks";
import {ActionType} from "../../ducks/types";
import {login, loginZAR, loginKeystore, createZAR, createKeystore} from "../../ducks/user";
import {CREATE_WALLET__SOFTWARE, WALLET} from "../../constants/routes";

type DispatchProps = {
  login: (pw: string) => Promise<Response>

  loginZARAccount: (email: string, pw: string) => Promise<Response>,
  createZARaccount: (email: string, pw: string, firstname: string, lastname: string) => Promise<Response>,

  loginKeystore: (keystore: string, pw: string) => Promise<Response>,
  createKeystore: (pw: string) => Promise<Response>,
}

type Props = DispatchProps & RouteComponentProps;

type State = {
  password: string,
  email: string,
  firstname: string,
  lastname: string,
  method: string,
  keystore: string,

  errorMessage: string,
  isLoggingIn: boolean,
}

class Authentication extends Component<Props, State> {
  state = {
    password: '',
    email: '',
    firstname: '',
    lastname: '',
    method: 'Keystore',
    keystore: '',
    errorMessage: '',
    isLoggingIn: false,
  };

  loginKeystore = async () => {
    const {
      password,
      keystore
    } = this.state;

    if(!keystore || keystore === "") {
      this.setState({errorMessage: 'Keystore file is required'})
      return false
    }

    if(!password || password === "") {
      this.setState({errorMessage: 'Password is required'})
      return false
    }

    this.setState({
      isLoggingIn: true,
      errorMessage: '',
    });

    const resp = await this.props.loginKeystore(keystore, password);

    if (resp.status !== 204) {
      this.setState({
        isLoggingIn: false,
        errorMessage: await resp.text(),
      });
    } else {
      this.setState({
        isLoggingIn: false,
      });
      this.props.history.push(WALLET);
    }
  }

  loginZAR = async () => {
    const {
      password,
      email
    } = this.state;

    if(!email || email === "") {
      this.setState({errorMessage: 'Email address is required'})
      return false
    }

    if(!password || password === "") {
      this.setState({errorMessage: 'Password is required'})
      return false
    }

    this.setState({
      isLoggingIn: true,
      errorMessage: '',
    });

    const resp = await this.props.loginZARAccount(email, password);

    if (resp.status !== 204) {
      this.setState({
        isLoggingIn: false,
        errorMessage: await resp.text(),
      });
    } else {
      this.setState({
        isLoggingIn: false,
      });
      this.props.history.push(WALLET);
    }
  }

  loginLocal = async () => {
    const { password } = this.state;

    this.setState({
      isLoggingIn: true,
      errorMessage: '',
    });

    const resp = await this.props.login(password);

    if (resp.status !== 204) {
      this.setState({
        isLoggingIn: false,
        errorMessage: await resp.text(),
      });
    } else {
      this.setState({
        isLoggingIn: false,
      });
      this.props.history.push(WALLET);
    }
  };

  render (): ReactNode {
    const { method } = this.state;

    return (
      <div className="connect-wallet">
        <div className="authenticate">
          <div className="authenticate__title">
            Unlock your account
          </div>
          <div className="authenticate__method">
            <Radio
              label="Keystore"
              onChange={e => this.setState({
                method: e.target.value,
                errorMessage: '',
              })}
              name="method"
              value="Keystore"
              checked={method==="Keystore"}
            />
            <Radio
              label="XAR Account"
              onChange={e => this.setState({
                method: e.target.value,
                errorMessage: '',
              })}
              name="method"
              value="ZAR"
              checked={method==="ZAR"}
            />
            <Radio
              label="Local Node"
              onChange={e => this.setState({
                method: e.target.value,
                errorMessage: '',
              })}
              name="method"
              value="Local"
              checked={method==="Local"}
            />
          </div>
          <div className="authenticate__action">
            { method === "Keystore" && this.renderKeystore() }
            { method === "ZAR" && this.renderZAR() }
            { method === "Local" && this.renderLocalNode() }
          </div>
        </div>
      </div>
    )
  }

  renderKeystore() {
    const { isLoggingIn } = this.state;

    return(
      <div>
        <div className="authenticate__text">
          Connect an encrypted wallet file and input your password
        </div>
        <Input
          label="Keystore File"
          type="file"
          onChange={e => this.setState({
            keystore: e.target.value,
            errorMessage: '',
          })}
          value={this.state.keystore}
          autoFocus
        />
        <Input
          label="Password"
          type="password"
          onChange={e => this.setState({
            password: e.target.value,
            errorMessage: '',
          })}
          onKeyPress={e => {
            if (e.key === 'Enter') {
              e.stopPropagation();
              this.loginKeystore();
            }
          }}
          value={this.state.password}
          autoFocus
        />
        <div className="authenticate__button">
          <div className="authenticate__error-message">
            { this.state.errorMessage }
          </div>
          <Button
            type="primary"
            onClick={this.loginKeystore}
            disabled={isLoggingIn}
            loading={isLoggingIn}
          >
            { isLoggingIn ? 'Logging In' : 'Login' }
          </Button>
        </div>
      </div>
    );
  }

  renderZAR() {
    const { isLoggingIn } = this.state;

    return (
      <div>
        <div className="authenticate__text">
          Connect an encrypted wallet using your zar.cloud account
        </div>
        <Input
          label="Email Address"
          type="text"
          onChange={e => this.setState({
            email: e.target.value,
            errorMessage: '',
          })}
          onKeyPress={e => {
            if (e.key === 'Enter') {
              e.stopPropagation();
              this.loginZAR();
            }
          }}
          value={this.state.email}
          autoFocus
        />
        <Input
          label="Password"
          type="password"
          onChange={e => this.setState({
            password: e.target.value,
            errorMessage: '',
          })}
          onKeyPress={e => {
            if (e.key === 'Enter') {
              e.stopPropagation();
              this.loginZAR();
            }
          }}
          value={this.state.password}
        />
        <div className="authenticate__button">
          <div className="authenticate__error-message">
            { this.state.errorMessage }
          </div>
          <Button
            type="primary"
            onClick={this.loginZAR}
            disabled={isLoggingIn}
            loading={isLoggingIn}
          >
            { isLoggingIn ? 'Logging In' : 'Login' }
          </Button>
        </div>
      </div>
    );
  }

  renderLocalNode() {
    const { isLoggingIn } = this.state;

    return (
      <div>
        <div className="authenticate__text">
          Connect an encrypted wallet using your local node
        </div>
        <Input
          label="Password"
          type="password"
          onChange={e => this.setState({
            password: e.target.value,
            errorMessage: '',
          })}
          onKeyPress={e => {
            if (e.key === 'Enter') {
              e.stopPropagation();
              this.loginLocal();
            }
          }}
          value={this.state.password}
        />
        <div className="authenticate__button">
          <div className="authenticate__error-message">
            { this.state.errorMessage }
          </div>
          <Button
            type="primary"
            onClick={this.loginLocal}
            disabled={isLoggingIn}
            loading={isLoggingIn}
          >
            { isLoggingIn ? 'Logging In' : 'Login' }
          </Button>
        </div>
      </div>
    )
  }
}

function mapDispatchToProps (dispatch: ThunkDispatch<REDUX_STATE, any, ActionType<any>>): DispatchProps {
  return {
    login: (pw: string) => dispatch(login(pw)),
    loginZARAccount: (email: string, pw: string) => dispatch(loginZAR(email, pw)),
    createZARaccount: (email: string, pw: string, firstname: string, lastname: string) => dispatch(createZAR(email, pw, firstname, lastname)),
    loginKeystore: (keystore: string, pw: string) => dispatch(loginKeystore(keystore, pw)),
    createKeystore: (pw: string) => dispatch(createKeystore(pw)),
  }
}

export default withRouter(
  connect(null, mapDispatchToProps)(Authentication)
);
