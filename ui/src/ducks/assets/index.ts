import {ActionType} from "../types";
import {FundingSources} from "../../constants/clearinghouse";

export type AssetType = {
  symbol: string
  name: string
  decimals: number
  nativeDecimals: number
  sources: FundingSources[]
  chainId: string
}

export type ChainType = {
  id: string
  name: string
  depositFinality: number
}

export type AssetStateType = {
  symbolToAssetId: {
    [k: string]: string
  }
  assets: {
    [k: string]: AssetType
  }
  chains: {
    [k: string]: ChainType
  }
}

const initialState: AssetStateType = {
  symbolToAssetId: {
    UCSDT: '1',
    UFTM: '2',
  },
  assets: {
    '1': {
      symbol: 'UCSDT',
      name: 'Collateralized Stable Debt Token',
      decimals: 18,
      nativeDecimals: 4,
      sources: [],
      chainId: '',
    },
    '2': {
      symbol: 'UFTM',
      name: 'Fantom Token',
      decimals: 18,
      nativeDecimals: 4,
      sources: [],
      chainId: '',
    },
  },
  chains: {
    ETH: {
      id: 'ETH',
      name: 'Rinkeby - Ethereum Testnet',
      depositFinality: 1,
    },
  },
};

export default function assetReducer(state: AssetStateType = initialState, action: ActionType<any>) {
  switch (action.type) {
    default:
      return state;
  }
}
