import { combineReducers } from 'redux';
import { reducer as announcementReducer } from './announcement/reducer';
import { reducer as authReducer } from './auth/reducer';
import { IAnnouncement } from './announcement/state';
import { IAuth } from './auth/state';

/**
 * IReducer is the entire state of the App
 */
export interface IReducer {
  announcement: IAnnouncement;
  auth: IAuth;
}

export default combineReducers<IReducer>({
  announcement: announcementReducer,
  auth: authReducer
});
