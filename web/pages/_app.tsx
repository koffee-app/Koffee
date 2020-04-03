import React from 'react';
import App, { AppInitialProps } from 'next/app';
import initializeWebSocketComponents from '../ws/initializeWebSocketComponents';

export default class MyApp extends App<AppInitialProps> {
  static async getInitialProps({ Component, ctx }) {
    console.log('xxx');
    initializeWebSocketComponents();
    return {
      pageProps: Component.getInitialProps
        ? await Component.getInitialProps(ctx)
        : {}
    };
  }

  render() {
    const { Component, pageProps } = this.props;
    return <Component {...pageProps} />;
  }
}
