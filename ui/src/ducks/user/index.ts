import { Dispatch } from "redux";
import { ActionType } from '../types';
import {post, get, postZar, GetUserOrderResponse, BalanceResponse} from "../../utils/fetch";
import { unlockAccount } from '../../utils/zar-util';
import {addOrders, OrderType} from "../exchange";
import BigNumber from "bignumber.js";
import {bn} from "../../utils/bn";

const crypto = require('crypto');
const sha256 = require('sha256');
const bip39 = require('bip39');

export const ADD_USER_ORDERS = 'app/user/addUserOrders';
export const ADD_USER_ADDRESS = 'app/user/addUserAddress';
export const ADD_USER_KEYSTORE = 'app/user/addUserKeystore';
export const SET_ORDER_HISTORY_FILTER = 'app/user/setOrderHistoryFilter';
export const SET_BALANCE = 'app/user/setBalance';
export const SET_LOGIN = 'app/user/setLogin';

export enum ORDER_HISTORY_FILTERS {
  'ALL',
  'OPEN',
}

export type BalanceType = {
  denom: string,
  locked: BigNumber,
  unlocked: BigNumber,
}

export type UserStateType = {
  orderHistoryFilter: ORDER_HISTORY_FILTERS
  orders: string[]
  transactions: string[]
  balances: {
    [denom: string]: BalanceType
  }
  address: string
  isLoggedIn?: boolean
}

const initialState = {
  orderHistoryFilter: ORDER_HISTORY_FILTERS.OPEN,
  orders: [],
  transactions: [],
  balances: {},
  address: '',
  isLoggedIn: undefined,
};

export const setOrderHistoryFilter = (filter: ORDER_HISTORY_FILTERS): ActionType<ORDER_HISTORY_FILTERS> => ({
  type: SET_ORDER_HISTORY_FILTER,
  payload: filter,
});

export const addUserOrders = (orders: OrderType[]): ActionType<OrderType[]> => ({
  type: ADD_USER_ORDERS,
  payload: orders,
});

export const login = (password: string) => async (dispatch: Dispatch): Promise<Response> => {
  const resp = await post('/auth/login', { username: 'validator', password });

  if (resp.status === 204) {
    const addrRes = await get('/auth/me');
    const addrJSON: { address: string} = await addrRes.json();
    dispatch(setAddress(addrJSON.address));
    dispatch({ type: '%INIT' });
  }

  return resp;
};

export const loginZAR = (email: string, password: string) => async (dispatch: Dispatch): Promise<Response> => {
  const resp = await postZar('/api/v1/login', _encrypt({ email_address: email, password }, '/api/v1/login'), '', '');


  // ok, so the call is successful. Now we need to figure out session management using this.
  // dispatch(setAddress(addrJSON.address));
  // dispatch({ type: '%INIT' });


  return resp;
};

export const loginKeystore = (keystore: string, password: string) => async (dispatch: Dispatch) => {
  const resp = await unlockAccount({ keystore, password });

  dispatch(setAddress(resp.address));
  dispatch(setKeystore(keystore));
  dispatch({ type: '%INIT' });

  return resp;
};

export const createZAR = (email: string, password: string, firstname: string, lastname: string) => async (dispatch: Dispatch): Promise<Response> => {
  const resp = await post('/auth/login', { username: 'validator', password });
  return resp;
};

export const createKeystore = (password: string) => async (dispatch: Dispatch): Promise<Response> => {
  const resp = await post('/auth/login', { username: 'validator', password });
  return resp;
};

function _encrypt(postData:Object, url:string) {
  const signJson = JSON.stringify(postData);
  const signMnemonic = bip39.generateMnemonic();
  const cipher = crypto.createCipher('aes-256-cbc', signMnemonic);
  const signEncrypted =
    cipher.update(signJson, 'utf8', 'base64') + cipher.final('base64');
  var signData = {
    e: hexEncode(signEncrypted),
    m: hexEncode(signMnemonic),
    u: sha256(url.toLowerCase()),
    p: sha256(sha256(url.toLowerCase())),
    t: new Date().getTime(),
    s: undefined
  };
  const signSeed = JSON.stringify(signData);
  const signSignature = sha256(signSeed);
  signData.s = signSignature;
  postData = JSON.stringify(signData);

  return postData;
}

function hexEncode (str:String) {
  var hex, i;
  var result = '';
  for (i = 0; i < str.length; i++) {
    hex = str.charCodeAt(i).toString(16);
    result += ('000' + hex).slice(-4);
  }
  return result;
};




export const logout = () => (dispatch: Dispatch) => {
  dispatch(setAddress(''));
  dispatch(setLogin(false));

  return true
}

export const setLogin = (payload: boolean): ActionType<boolean> => ({
  type: SET_LOGIN,
  payload,
});

export const setBalance  = (payload: BalanceType): ActionType<BalanceType> => ({
  type: SET_BALANCE,
  payload,
});

export const setAddress = (payload: string): ActionType<string> => ({
  type: ADD_USER_ADDRESS,
  payload,
});

export const setKeystore = (payload: string): ActionType<string> => ({
  type: ADD_USER_KEYSTORE,
  payload,
});

export const checkLogin = () => async (dispatch: Dispatch) => {
  try {
    const resp = await get('/user/balances');
    switch (resp.status) {
      case 401:
        return dispatch(setLogin(false));
      case 200:
        const addrRes = await get('/auth/me');
        const addrJSON: { address: string} = await addrRes.json();
        dispatch({ type: '%INIT' });
        dispatch(setAddress(addrJSON.address));
        dispatch(setLogin(true));
    }
  } catch (e) {
    console.log(e)
  }
};

export const fetchUserOrders = () => async (dispatch: Dispatch<ActionType<OrderType[]>>) => {
  try {
    const resp = await get('/user/orders');
    const json: GetUserOrderResponse = await resp.json();
    dispatch(addOrders(json.orders || []));
    dispatch(addUserOrders(json.orders || []));
  } catch (e) {
    console.log(e);
  }
};

export const addUserAddress = (address: string) => ({
  type: ADD_USER_ADDRESS,
  payload: address,
});

export const addUserKeystore = (keystore: string) => ({
  type: ADD_USER_KEYSTORE,
  payload: keystore,
});

export const fetchBalance = () => async (dispatch: Dispatch<ActionType<BalanceType>>) => {
  try {
    const resp = await get('/user/balances');
    const json: BalanceResponse = await resp.json();

    json.balances.forEach(balance => {
      console.log(balance)
      dispatch(setBalance({
        denom: balance.denom,
        locked: bn(balance.at_risk),
        unlocked: bn(balance.amount),
      }))
    })
  } catch (e) {
    console.log(e);
  }
};

export default function userReducer(state: UserStateType = initialState, action: ActionType<any>): UserStateType {
  switch (action.type) {
    case SET_ORDER_HISTORY_FILTER:
      return handleSetOrderHistoryFilter(state, action);
    case ADD_USER_ORDERS:
      return handleAddOrders(state, action);
    case SET_BALANCE:
      return handleSetBalance(state, action);
    case ADD_USER_ADDRESS:
      return {
        ...state,
        address: action.payload,
      };
    case SET_LOGIN:
      return {
        ...state,
        isLoggedIn: action.payload,
      };
    case '%INIT':
      return {
        ...state,
        isLoggedIn: true,
      };
    default:
      return state;
  }
}

function handleAddOrders(state: UserStateType, action: ActionType<OrderType[]>): UserStateType {
  return {
    ...state,
    orders: action.payload.map(({ id }) => id),
  }
}

function handleSetOrderHistoryFilter(state: UserStateType, action: ActionType<ORDER_HISTORY_FILTERS>): UserStateType {
  return {
    ...state,
    orderHistoryFilter: action.payload,
  };
}

function handleSetBalance(state: UserStateType, action: ActionType<BalanceType>): UserStateType {
  const { denom } = action.payload;
  return {
    ...state,
    balances: {
      ...state.balances,
      [denom]: action.payload,
    },
  };
}
