import { TAnnouncement } from './reducer';
import { ThunkAction } from 'redux-thunk';
import { IStore } from 'states/type';
import { Action } from 'redux';

export const UPDATE_ANNOUNCEMENT = 'ANNOUNCEMENT/update_announcement';

export const updateAnnouncement = (
  message: string
): ThunkAction<void, IStore, unknown, Action<string>> => dispatch => {
  dispatch({ type: UPDATE_ANNOUNCEMENT, message } as TAnnouncement);
};
