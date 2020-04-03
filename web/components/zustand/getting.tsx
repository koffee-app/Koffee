import React from 'react';
import { useGlobalState } from '../../store';

const trying = () => {
  const lastTime = useGlobalState(state => state.lastUpdate);
  return <div>Times: {lastTime}</div>;
};

export default trying;
