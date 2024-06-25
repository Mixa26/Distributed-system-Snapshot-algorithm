package servent.message.snapshot;

import app.ServentInfo;
import app.snapshot_bitcake.LYSnapshotResult;
import servent.message.BasicMessage;
import servent.message.MessageType;

import java.util.List;
import java.util.Map;
import java.util.Set;

public class RoundMessage extends BasicMessage {

    Map<Integer, LYSnapshotResult> lySnapshotResult;

    public RoundMessage(ServentInfo sender, ServentInfo receiver, Map<Integer, LYSnapshotResult> lySnapshotResult) {
        super(MessageType.ROUND, sender, receiver);

        this.lySnapshotResult = lySnapshotResult;
    }

    public Map<Integer, LYSnapshotResult> getLySnapshotResult() {
        return lySnapshotResult;
    }
}
