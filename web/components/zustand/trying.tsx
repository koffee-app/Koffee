import React from 'react';
import { useGlobalState } from '../../store';

const trying = () => {
  const tick = useGlobalState(state => state.tick);
  const lastTime = useGlobalState(state => state.lastUpdate);
  return (
    <div>
      <button
        onClick={() => {
          tick({ lastUpdate: lastTime + 1, light: true });
        }}
      >
        Click click
      </button>
    </div>
  );
};

export default trying;
