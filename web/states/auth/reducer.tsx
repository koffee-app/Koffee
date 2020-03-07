import { IAuth, IUser } from './state';
import * as authActions from './actions';

const initialState: IAuth = {
  user: null,
  authenticated: false
};

export type TAuth = {
  user: IUser;
  type: string;
};

export const reducer = (state = initialState, action: TAuth): IAuth => {
  switch (action.type) {
    case authActions.REGISTER_USER:
      return { ...state, user: action.user, authenticated: true };
    default:
      return state;
  }
};
