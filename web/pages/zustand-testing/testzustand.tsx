import React from 'react';
import { NextPage } from 'next';
// import { ZustandContext } from '../store';
import Getting from '../../components/zustand/getting';
import Trying from '../../components/zustand/trying';
import { withZustand } from '../../lib/zustand';
import Link from 'next/link';

const TestZustand: NextPage = () => {
  return (
    <>
      <Getting></Getting>
      <Trying></Trying>
      <Link href="/zustand-testing/zustand">
        <a>Go</a>
      </Link>
    </>
  );
};

// Zustand.getInitialProps = ({
//   zustandStore
// }: ZustandContext & NextPageContext) => {
//   zustandStore.setState({ light: true, lastUpdate: 3 });
//   return {};
// };

export default withZustand(TestZustand);
