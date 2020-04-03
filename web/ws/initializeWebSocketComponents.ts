import test from './TestWSComp';

let initialized = false;

export default function initializeWebSocketComponents() {
  if (initialized) return;
  initialized = true;
  new test();
}
