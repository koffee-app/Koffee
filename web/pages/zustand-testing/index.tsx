import React from 'react';
import Link from 'next/link';

import { withZustand } from '../../lib/zustand';
import { NextPage } from 'next';
import { useGlobalState } from '../../store';
import { useMessaging } from '../../store/messaging';

const App: NextPage = () => {
  const announcementMessage = useMessaging(
    messaging => messaging.announcement.message
  );
  const tick = useGlobalState(({ tick }) => tick);
  return (
    <div>
      <h1>{announcementMessage}</h1>
      <button
        onClick={() => {
          tick({ lastUpdate: 1, light: false });
        }}
      ></button>
      <Link href="/zustand-testing/zustand">
        <a>Go</a>
      </Link>
    </div>
  );
};

export default withZustand(App);
