import React, {Component, ChangeEvent, ClipboardEvent, KeyboardEvent} from 'react';
import './radio.scss';

type PropTypes = {
  className?: string
  label?: string
  disabled?: boolean
  onChange?: (event: ChangeEvent<HTMLInputElement>) => void
  value?: string
  name?: string
  checked?: boolean
}

class Radio extends Component<PropTypes> {
  render() {
    const {
      className = '',
      label,
      disabled,
      onChange,
      value,
      name,
      checked
    } = this.props;

    return (
      <div className={`radio ${className}`}>
        { label && <div className="radio__label">{label}</div> }
        <div className="radio__wrapper">
          <input
            className="radio__input"
            type='radio'
            disabled={disabled}
            onChange={onChange}
            value={value}
            name={name}
            checked={checked}
          />
        </div>
      </div>
    )
  }
}

export default Radio;
