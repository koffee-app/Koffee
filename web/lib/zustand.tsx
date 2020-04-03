import React from 'react';
import { getState, ZustandContext } from '../store';
import { NextPage, NextPageContext } from 'next';

export const withZustand = (
  PageComponent: NextPage<any>,
  { ssr = true } = {}
) => {
  const WithZustand = ({ initialZustandState, ...props }) => {
    const [, zustandStore] = getState(initialZustandState);
    zustandStore.setState(initialZustandState);
    return <PageComponent {...props} />;
  };

  // Set the correct displayName in development
  if (process.env.NODE_ENV !== 'production') {
    const displayName =
      PageComponent.displayName || PageComponent.name || 'Component';

    WithZustand.displayName = `withZustand(${displayName})`;
  }

  // Get or Create the store with `undefined` as initialState
  // This allows you to set a custom default initialState
  const [, zustandStore] = getState();

  if (ssr || PageComponent.getInitialProps) {
    WithZustand.getInitialProps = async (
      context: NextPageContext & ZustandContext
    ) => {
      // Provide the store to getInitialProps of pages
      context.zustandStore = zustandStore;

      // Run getInitialProps from HOCed PageComponent
      const pageProps =
        typeof PageComponent.getInitialProps === 'function'
          ? await PageComponent.getInitialProps(context)
          : {};

      // Pass props to PageComponent
      return {
        ...pageProps,
        initialZustandState: zustandStore.getState()
      };
    };
  }

  return WithZustand;
};
