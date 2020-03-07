import { ThunkAction } from 'redux-thunk';
import { Action } from 'redux';
import { IStore } from 'states/type';
import { TAuth } from './reducer';
import AuthAPI from '../../api/Auth';

export const REGISTER_USER = 'AUTH/register_user';

export interface IUserForm {
  password: string;
  email: string;
}

export const registerAction = (
  userData: IUserForm
): ThunkAction<void, IStore, unknown, Action<string>> => dispatch => {
  AuthAPI.register(userData)
    .success(user => dispatch({ type: REGISTER_USER, user: user } as TAuth))
    .error(e => {
      console.log('ERRORS HAVE BEEN MADE');
      console.log(e);
    });
};
