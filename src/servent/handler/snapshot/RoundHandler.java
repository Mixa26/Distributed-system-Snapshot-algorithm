package servent.handler.snapshot;

import app.snapshot_bitcake.SnapshotCollector;
import app.snapshot_bitcake.SnapshotCollectorWorker;
import servent.handler.MessageHandler;
import servent.message.Message;
import servent.message.MessageType;
import servent.message.snapshot.RoundMessage;

public class RoundHandler implements MessageHandler {

    private Message clientMessage;

    private SnapshotCollector snapshotCollector;

    public RoundHandler(Message clientMessage, SnapshotCollector snapshotCollector) {
        this.clientMessage = clientMessage;
        this.snapshotCollector = snapshotCollector;
    }

    @Override
    public void run() {
        if (clientMessage.getMessageType().equals(MessageType.ROUND)){
            RoundMessage roundMessage = (RoundMessage)clientMessage;
            ((SnapshotCollectorWorker)snapshotCollector).addRoundResult(roundMessage.getLySnapshotResult());
        }
    }

    public Message getClientMessage() {
        return clientMessage;
    }
}
