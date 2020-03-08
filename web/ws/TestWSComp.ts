import {
  TWebSocket,
  HandlerOnMessage,
  WebSocketComponent,
  HandlerOnOpen,
  IMessage,
  WS
} from './websocket';
import { MessageEvent, OpenEvent } from 'ws';
import { updateAnnouncement } from '../states/announcement/actions';

class TestWSComp extends WebSocketComponent {
  @HandlerOnOpen()
  onConnection(_: TWebSocket, __: OpenEvent) {
    WS.store.dispatch<any>(updateAnnouncement('HEMO conectao'));
  }

  @HandlerOnMessage(`chat.message`)
  onChatMsg(_: TWebSocket, __: MessageEvent, messageParsed: IMessage) {
    console.log(messageParsed);
    WS.store.dispatch<any>(updateAnnouncement('HEMOS RECIBIDO UN MENSAJE'));
  }

  public static sendChatMessage() {}
}

export default TestWSComp;
