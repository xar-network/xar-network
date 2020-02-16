package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Based in part on github.com/compound-finance/compound-protocol/contracts/Exponential

// Exponential module for storing fixed-precision decimals
//  Exp is a struct which stores decimals with a fixed precision of 18 decimal places.
//    Thus, if we wanted to store the 5.1, mantissa would store 5.1e18. That is:
//    `NewExp(5100000000000000000)`.
type Exp struct {
	mantissa sdk.Uint
}

func NewExp(exp sdk.Uint) Exp {
	return Exp{exp}
}

func ExpScale() sdk.Uint {
	return sdk.NewUint(1e18)
}

func HalfExpScale() sdk.Uint {
	return ExpScale().QuoUint64(2)
}

func MantissaOne() sdk.Uint { return ExpScale() }

// MultiplyScalarUint64 multiplies an Exp by a scalar, returning a new Exp
func (e Exp) MultiplyScalar(scalar sdk.Uint) Exp {
	scaledMantissa := e.mantissa.Mul(scalar)
	return Exp{scaledMantissa}
}

// MultiplyScalarUint64 multiplies an Exp by a scalar, returning a new Exp
func (e Exp) MultiplyScalarUint64(scalar uint64) Exp {
	scaledMantissa := e.mantissa.MulUint64(scalar)
	return Exp{scaledMantissa}
}

// MultiplyScalarTruncate multiplies an Exp by a scalar, then truncates to return an unsigned integer
func (e Exp) MultiplyScalarTruncate(scalar sdk.Uint) sdk.Uint {
	product := e.MultiplyScalar(scalar)
	return product.Truncate()
}

// Multiply an Exp by a scalar, truncate, then add to an unsigned integer, returning an unsigned integer
func (e Exp) MultiplyScalarTruncateAddUInt(scalar sdk.Uint, addition sdk.Uint) sdk.Uint {
	product := e.MultiplyScalar(scalar)
	return product.Truncate().Add(addition)
}

// Truncates the given exp to a whole number value.
// eg NewExp(sdk.NewUint(15).Mul(ExpScale())).Truncate() = 15
func (e Exp) Truncate() sdk.Uint {
	return e.mantissa.Quo(ExpScale())
}
