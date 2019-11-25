import { combineReducers } from 'redux';
import general, { GeneralStateType } from './general';
import exchange, { ExchangeStateType } from './exchange';
import user, { UserStateType } from './user';

export type REDUX_STATE = {
  general: GeneralStateType
  exchange: ExchangeStateType
  user: UserStateType
}

export default combineReducers({
  general,
  exchange,
  user,
});
