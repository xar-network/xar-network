package pool

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//var _ := Marshaler

var keyHandlers = map[string]func(r *ReservePool, val interface{}) error{
	"native_coins":     jsonAddNative,
	"non_native_coins": jsonAddNonNative,
	"liquidity_coins":  jsonAddLiquidity,
}

func jsonAddNative(r *ReservePool, val interface{}) error {
	coinMap, ok := val.(map[string]interface{})
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}

	amt, ok := coinMap["amount"].(string)
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}

	amount, ok := sdk.NewIntFromString(amt)
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}
	coin := sdk.NewCoin(coinMap["denom"].(string), amount)

	r.nativeCoins = coin
	return nil
}

func jsonAddNonNative(r *ReservePool, val interface{}) error {
	coinMap, ok := val.(map[string]interface{})
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}

	amt, ok := coinMap["amount"].(string)
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}

	amount, ok := sdk.NewIntFromString(amt)
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}
	coin := sdk.NewCoin(coinMap["denom"].(string), amount)

	r.nonNativeCoins = coin
	return nil
}

func jsonAddLiquidity(r *ReservePool, val interface{}) error {
	coinMap, ok := val.(map[string]interface{})
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}

	amt, ok := coinMap["amount"].(string)
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}

	amount, ok := sdk.NewIntFromString(amt)
	if !ok {
		return errors.New("cannot unmarshal native_coins")
	}
	coin := sdk.NewCoin(coinMap["denom"].(string), amount)

	r.liquidityCoins = coin
	return nil
}

// a marshaller logic is implemented below

func (r ReservePool) MarshalJSON() ([]byte, error) {
	nc, err := json.Marshal(r.nativeCoins)
	if err != nil {
		return nil, err
	}
	//nc = bytes.ReplaceAll(nc, []byte("\""), []byte("\\\""))

	nnc, err := json.Marshal(r.nonNativeCoins)
	if err != nil {
		return nil, err
	}
	//nnc = bytes.ReplaceAll(nnc, []byte("\""), []byte("\\\""))

	lc, err := json.Marshal(r.liquidityCoins)
	if err != nil {
		return nil, err
	}
	//lc = bytes.ReplaceAll(lc, []byte("\""), []byte("\\\""))

	poolJson := fmt.Sprintf(`{
		"native_coins": %s,
		"non_native_coins": %s,
		"liquidity_coins": %s
	}`, string(nc), string(nnc), string(lc))

	return []byte(poolJson), nil
}

// this is a slow solution.
func (r *ReservePool) UnmarshalJSON(bz []byte) error {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(bz, &jsonMap)
	if err != nil {
		return err
	}

	for key, handler := range keyHandlers {
		val, ok := jsonMap[key]
		if ok {
			err = handler(r, val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (r ReservePool) MarshalAmino() (string, error) {
	b, err := r.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// UnmarshalAmino defines custom decoding scheme
func (r *ReservePool) UnmarshalAmino(text string) error {
	b := []byte(text)

	return r.UnmarshalJSON(b)
}