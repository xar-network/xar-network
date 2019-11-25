import React, {Component, ReactNode} from "react";
import BigNumber from "bignumber.js";
import {SPREAD_TYPE} from "../../../ducks/exchange";
import {TableCell, TableRow} from "../../../components/ui/Table";
import c from "classnames";

type Props = {
  price: BigNumber
  quantity: BigNumber
  max: BigNumber
  side: SPREAD_TYPE
  baseDenom: string
  quoteDenom: string
}

type State = {
  shouldHighlight: boolean
}

class DepthTableRow extends Component<Props, State> {
  state = {
    shouldHighlight: false,
  };

  flashOnTimeout: any | null = null;
  flashOffTimeout: any | null = null;

  componentDidMount() {
    this.flash();
  }

  componentDidUpdate(lastProps: Props) {
    if (!lastProps.quantity.isEqualTo(this.props.quantity)) {
      this.flash();
    }
  }

  componentWillUnmount() {
    if (this.flashOnTimeout) clearTimeout(this.flashOnTimeout);
    if (this.flashOffTimeout) clearTimeout(this.flashOffTimeout);
  }

  flash () {
    this.flashOnTimeout = setTimeout(() => this.setState({ shouldHighlight: true}), 0);
    this.flashOffTimeout = setTimeout(() => this.setState({ shouldHighlight: false }), 1000);
  }

  render (): ReactNode {
    const {
      max,
      price,
      quantity,
    } = this.props;
    const { shouldHighlight } = this.state;

    const priceN = price;
    let quantityN = quantity.div(10 ** 8);
    let isPendingFills = false;

    if (quantityN.isZero()) {
      return null;
    }

    const totalN = priceN.multipliedBy(quantityN);
    const percentN = price.multipliedBy(quantity).div(max).multipliedBy(100);

    const p = priceN.toFixed(8);
    const q = quantityN.toFixed(Math.min(8, 8))

    return (
      <TableRow
        className={c({
          'exchange__depth__crossed-row': isPendingFills,
          'exchange__depth__highlight': shouldHighlight,
        })}
      >
        <TableCell>{p}</TableCell>
        <TableCell>{q}</TableCell>
        <TableCell>{totalN.toFixed(8)}</TableCell>
        <div
          className="exchange__depth__percentage"
          style={{
            width: `${percentN.toFixed(0)}%`
          }}
        />
      </TableRow>
    );
  }
}

export default DepthTableRow;
