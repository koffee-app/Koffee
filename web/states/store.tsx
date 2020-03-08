import { createStore, applyMiddleware } from 'redux';
import { composeWithDevTools } from 'redux-devtools-extension';
import thunkMiddleware from 'redux-thunk';
import reducer from './reducer';
import { WS } from '../ws/websocket';
import all from '../ws/all';

export const initStore = initialState => {
  const store = createStore(
    reducer,
    initialState,
    composeWithDevTools(applyMiddleware(thunkMiddleware))
  );
  WS.initializeStore(store);
  all(); //lol
  return store;
};
