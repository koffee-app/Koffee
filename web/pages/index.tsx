import React from 'react';
import Link from 'next/link';

import { IAnnouncement } from 'states/announcement/state';
import { connect } from 'react-redux';
import { updateAnnouncement } from '../states/announcement/actions';
import { IStore } from 'states/type';

interface IProps {
  announcement: IAnnouncement;
  updateAnnouncement: (msg: string) => void;
}

const App: React.FC<IProps> = props => {
  return (
    <div>
      <h1>{props.announcement.message}</h1>
      <button onClick={() => props.updateAnnouncement('yes')}></button>
      <Link href="/test">
        <a>Go</a>
      </Link>
    </div>
  );
};

const mapStateToProps = (state: IStore) => ({
  announcement: state.announcement
});

const actions = { updateAnnouncement };

export default connect(mapStateToProps, actions)(App);
