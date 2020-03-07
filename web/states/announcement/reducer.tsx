import { IAnnouncement } from './state';
import * as announcementActions from './actions';

const initialState: IAnnouncement = {
  message: 'No announcement...'
};

export type TAnnouncement = {
  message: string;
  type: string;
};

export const reducer = (
  state = initialState,
  action: TAnnouncement
): IAnnouncement => {
  switch (action.type) {
    case announcementActions.UPDATE_ANNOUNCEMENT:
      return { ...state, message: action.message };
    default:
      return state;
  }
};
