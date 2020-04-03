import { create, SetState, UseStore, StoreApi } from 'zustand';
import { MessagingState, messagingStore, ReducerMsg } from './messaging';
import { WS } from '../ws/websocket';
import { devtools } from 'zustand/middleware';

// state of this store
interface State {
  lastUpdate: number;
  light: boolean;
}

// actions of this store
interface Reducer {
  tick: ({ lastUpdate: number, light: boolean }) => void;
}

type Actions = Reducer & ReducerMsg;

// The entire store and actions
type ZustandState = State & Reducer & MessagingState;

// A Zustand reducer
export type ReducerFn = (set: SetState<ZustandState>) => Reducer;

// Entire context of the app
export interface ZustandContext {
  zustandStore: StoreApi<ZustandState>;
  // More stores
  messaging: StoreApi<MessagingState>;
}

// initialState of this store
const initialState: State = {
  lastUpdate: 0,
  light: false
};

// Reducer of one store.
const reducer: ReducerFn = set => ({
  tick: ({ lastUpdate, light }) => {
    set({ lastUpdate, light });
  }
});

// Initializer
export const initializeStore = (preloadedState = initialState) => {
  return create<ZustandState>(
    devtools(
      (set, get) => ({
        ...preloadedState,
        ...reducer(set),
        ...messagingStore(set, get)
      }),
      'STATE'
    )
  );
};

// store globally the zustand state
let zustandStore: [UseStore<ZustandState>, StoreApi<ZustandState>];

export const getState = (initialState?: ZustandState) => {
  // Always make a new store if server, otherwise state is shared between requests
  if (typeof window === 'undefined') {
    return initializeStore(initialState);
  }

  // Create store if unavailable on the client and set it on the window object
  if (!zustandStore) {
    new WS();
    zustandStore = initializeStore(initialState);
  }

  return zustandStore;
};

const [_useGlobalState, storeAPI] = getState();

export const useGlobalState = _useGlobalState;
export const actions = storeAPI.getState() as Actions;
