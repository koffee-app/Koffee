import { StateCreatorDev } from '../lib/stateCreatorDev';

interface Store {
  announcement: {
    message: String;
  };
}

export interface ReducerMsg {
  setAnnouncement: (msg: String) => void;
}

export type MessagingState = Store & ReducerMsg;

// type ReducerFn = (
//   set: SetState<MessagingState>,
//   get: GetState<MessagingState>
// ) => ReducerMsg;

const initialStateMessaging: Store = {
  announcement: {
    message: 'No state yet...'
  }
};

const reducerMessaging: StateCreatorDev<Store, ReducerMsg> = (set, get) => ({
  setAnnouncement: msg => {
    const state = get();
    set(
      { announcement: { ...state.announcement, message: msg } },
      'messaging.setAnnouncement'
    );
  }
});

export const messagingStore = (set, get) => ({
  ...initialStateMessaging,
  ...reducerMessaging(set, get)
});
