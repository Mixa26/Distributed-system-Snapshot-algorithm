package servent.message.snapshot;

import app.ServentInfo;
import servent.message.BasicMessage;
import servent.message.MessageType;

public class LYMarkerMessage extends BasicMessage {

	private static final long serialVersionUID = 388942509576636228L;

	private Integer snapshotNumber;

	private LYTellMessage lyTellMessage;

	private Integer collectorId;

	public LYMarkerMessage(int snapshotNumber, ServentInfo sender, ServentInfo receiver, int collectorId, LYTellMessage lyTellMessage) {
		super(MessageType.LY_MARKER, sender, receiver, String.valueOf(collectorId));

		this.snapshotNumber = snapshotNumber;
		this.lyTellMessage = lyTellMessage;
		this.collectorId = collectorId;
	}

	public Integer getSnapshotNumber() {
		return snapshotNumber;
	}

	public LYTellMessage getLyTellMessage() {
		return lyTellMessage;
	}

	public Integer getCollectorId() {
		return collectorId;
	}
}
