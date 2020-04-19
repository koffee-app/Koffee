import {
  TWebSocket,
  HandlerOnMessage,
  WebSocketComponent,
  HandlerOnOpen,
  IMessage
} from './websocket';
import { MessageEvent, OpenEvent } from 'ws';
import { actions } from '../store';

class TestWSComp extends WebSocketComponent {
  @HandlerOnOpen()
  onConnection(_: TWebSocket, __: OpenEvent) {
    const methods = actions;
    console.log('worked');
    methods.setAnnouncement('worked');
  }

  @HandlerOnMessage(`chat.message`)
  onChatMsg(_: TWebSocket, __: MessageEvent, messageParsed: IMessage) {
    const methods = actions;
    console.log(messageParsed);
    methods.setAnnouncement(`The new msg is: ${messageParsed.data}`);
  }

  public static sendChatMessage() {}
}

export default TestWSComp;
