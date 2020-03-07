import React from 'react';
import { Provider } from 'react-redux';
import App, { AppInitialProps } from 'next/app';
import withRedux from 'next-redux-wrapper';
import { initStore } from '../states/store';
import { IStore } from '../states/type';

export default withRedux(initStore)(
  class MyApp extends App<AppInitialProps & { store: IStore }> {
    static async getInitialProps({ Component, ctx }) {
      return {
        pageProps: Component.getInitialProps
          ? await Component.getInitialProps(ctx)
          : {}
      };
    }

    render() {
      const { Component, pageProps, store } = this.props;
      return (
        <Provider store={store}>
          <Component {...pageProps} />
        </Provider>
      );
    }
  }
);
