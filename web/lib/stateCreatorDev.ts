import { SetState, State, StateCreator } from 'zustand';

/** To allow named functions when using devtools */
export type StateCreatorDev<T extends State, Returns> = (
  set: (partial: Parameters<SetState<T>>[0], name?: string) => void,
  get: Parameters<StateCreator<T>>[1]
) => Returns;
