import React, {Component, ReactNode} from "react";
import "./transfer.scss";
import Icon from "../../components/ui/Icon";
import ArrowRight from "../../assets/icons/arrow-right.svg";
import Button from "../../components/ui/Button";
import {WALLET} from "../../constants/routes";
import {RouteComponentProps, withRouter} from "react-router";
import {connect} from "react-redux";
import {REDUX_STATE} from "../../ducks";
import Dropdown from "../../components/ui/Dropdown";
import Input from "../../components/ui/Input";
import CheckIcon from "../../assets/icons/check-green.svg";

enum FeeType {
  Slow = 'slow',
  Normal = 'normal',
  Fast = 'fast',
}

const FeeOptions = [
  { label: 'Slow - 0.002 DEMO', value: FeeType.Slow},
  { label: 'Normal - 0.004 DEMO', value: FeeType.Normal},
  { label: 'Fast - 0.013 DEMO', value: FeeType.Fast},
];

type StateProps = {
  denoms: {
    [k: string]: string
  }
  address: string
}

type Props = StateProps & RouteComponentProps

type State = {
  selectedDenom: string | null
  recipientAddress: string
  selectedFee: FeeType
  amount: string
  isReviewing: boolean
  isCompleted: boolean
  isSending: boolean
}

class Transfer extends Component<Props, State> {
  state = {
    selectedDenom: '',
    recipientAddress: '',
    selectedFee: FeeType.Normal,
    amount: '',
    isReviewing: false,
    isCompleted: false,
    isSending: false,
  };

  isValid(): boolean {
    const {
      selectedFee,
      selectedDenom,
      amount,
      recipientAddress,
      isSending,
    } = this.state;

    return !!selectedFee && !!amount && !!recipientAddress && !isSending;
  }

  send = () => {
    this.setState({ isSending: true });
    setTimeout(() => this.setState({
      isSending: false,
      isReviewing: false,
      isCompleted: true,
    }), 2000);
  };

  render(): ReactNode {
    return (
      <div className="wallet__content transfer">
        { this.renderHeader() }
        { this.renderContent() }
      </div>
    );
  }

  renderHeader(): ReactNode {
    const { history } = this.props;
    const { isReviewing, isCompleted } = this.state;

    if (isCompleted) {
      return (
        <div className="wallet__content__header">
          <span>My Wallet</span>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <span>Transfer</span>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <span>Review</span>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <span>Completed</span>
        </div>
      );
    }

    if (isReviewing) {
      return (
        <div className="wallet__content__header">
          <Button type="link" onClick={() => history.push(WALLET)}>
            My Wallet
          </Button>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <Button type="link" onClick={() => this.setState({ isReviewing: false })}>
            Transfer
          </Button>
          <Icon className="wallet__content__header__arrow" url={ArrowRight} />
          <span>Review</span>
        </div>
      );
    }

    return (
      <div className="wallet__content__header">
        <Button type="link" onClick={() => history.push(WALLET)}>
          My Wallet
        </Button>
        <Icon className="wallet__content__header__arrow" url={ArrowRight} />
        <span>Transfer</span>
      </div>
    );
  }

  renderContent(): ReactNode {
    const { isReviewing, isCompleted } = this.state;

    if (isCompleted) {
      return this.renderCompleted();
    }

    if (isReviewing) {
      return this.renderReview();
    }

    return this.renderSend();
  }

  renderCompleted(): ReactNode {
    const { denoms, history } = this.props;
    const {
      selectedDenom,
      recipientAddress,
      amount,
    } = this.state;

    return (
      <div className="wallet__content__body deposit deposit--done">
        <div className="deposit__body">
          <div className="deposit__body__hero">
            <div className="deposit__body__icon-wrapper">
              <Icon url={CheckIcon} width={60} height={60} />
            </div>
          </div>
          <div className="deposit__body__text">
            {`Your have initiated a transfer of ${amount} ${denoms[selectedDenom]} to ${recipientAddress}. The transaction will take about 5 minutes to complete.`}
          </div>
          <div className="deposit__footer">
            <Button
              type="primary"
              onClick={() => history.push(WALLET)}
            >
              Close
            </Button>
          </div>
        </div>
      </div>
    );
  }


  renderReview(): ReactNode {
    const { denoms } = this.props;
    const {
      selectedDenom,
      selectedFee,
      recipientAddress,
      amount,
      isSending,
    } = this.state;

    const feedIndex = FeeOptions.findIndex(({ value }) => value === selectedFee);

    return (
      <div className="transfer__form">
        { this.renderReviewRow('Asset Type', `${denoms[selectedDenom]}`)}
        { this.renderReviewRow('Recipient Address', recipientAddress)}
        { this.renderReviewRow('Fee', FeeOptions[feedIndex].label)}
        { this.renderReviewRow('Amount', `${amount} ${denoms[selectedDenom]}`)}
        <div className="transfer__form__actions">
          <Button
            type="primary"
            disabled={!this.isValid()}
            onClick={this.send}
            loading={isSending}
          >
            Transfer
          </Button>
        </div>
      </div>
    );
  }

  renderSend(): ReactNode {
    const { denoms } = this.props;
    const {
      selectedDenom,
      selectedFee,
      recipientAddress,
      amount,
    } = this.state;
    const denomItems = Object.entries(denoms)
      .map(([ denom ]) => ({
        label: `${denom}`,
        value: denom,
      }));

    const currentIndex = denomItems.findIndex(({ value }) => value === selectedDenom);
    const feedIndex = FeeOptions.findIndex(({ value }) => value === selectedFee);

    return (
      <div className="transfer__form">
        <div className="transfer__form__row">
          <div className="transfer__form__label">Denom</div>
          <Dropdown
            className="transfer__form__content transfer__asset-dropdown"
            items={denomItems}
            currentIndex={currentIndex}
            onSelect={id => this.setState({ selectedDenom: id })}
          />
        </div>
        <div className="transfer__form__row">
          <div className="transfer__form__label">Recipient Address</div>
          <Input
            type="text"
            className="transfer__form__content transfer__input"
            onChange={e => this.setState({ recipientAddress: e.target.value })}
            value={recipientAddress}
            placeholder="cosmos1j689jv..."
          />
        </div>
        <div className="transfer__form__row">
          <div className="transfer__form__label">Amount</div>
          <Input
            type="number"
            className="transfer__form__content transfer__input"
            step="0.01"
            onChange={e => this.setState({ amount: e.target.value })}
            value={amount}
            placeholder="0.01"
          />
        </div>
        <div className="transfer__form__actions">
          <Button
            type="primary"
            disabled={!this.isValid()}
            onClick={() => this.setState({ isReviewing: true })}
          >
            Review
          </Button>
        </div>
      </div>
    );
  }

  renderReviewRow(label: string, value: string) {
    return (
      <div className="transfer__review__row">
        <div className="transfer__review__row__label">{label}</div>
        <div className="transfer__review__row__value">{value}</div>
      </div>
    )
  }
}

function mapStateToProps(state: REDUX_STATE): StateProps {
  const {
    user: { denoms },
    user: { address },
  } = state;

  return {
    denoms,
    address,
  };
}

export default withRouter(
  connect(mapStateToProps)(Transfer)
);
