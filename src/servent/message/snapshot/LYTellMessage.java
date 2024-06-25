package servent.message.snapshot;

import java.util.List;
import java.util.Set;

import app.ServentInfo;
import app.snapshot_bitcake.LYSnapshotResult;
import servent.message.BasicMessage;
import servent.message.Message;
import servent.message.MessageType;

public class LYTellMessage extends BasicMessage {

	private static final long serialVersionUID = 3116394054726162318L;

	private Integer snapshotNumber;
	private List<LYSnapshotResult> lySnapshotResult;

	private Set<Integer> idBorderSet;

	public LYTellMessage(Integer snapshotNumber , ServentInfo sender, ServentInfo receiver, List<LYSnapshotResult> lySnapshotResult, Set<Integer> idBorderSet) {
		super(MessageType.LY_TELL, sender, receiver);

		this.snapshotNumber = snapshotNumber;
		this.lySnapshotResult = lySnapshotResult;
		this.idBorderSet = idBorderSet;
	}

	public List<LYSnapshotResult> getLYSnapshotResult() {
		return lySnapshotResult;
	}

	public Set<Integer> getIdBorderSet() {
		return idBorderSet;
	}

    /*@Override
	public Message setRedColor() {
		Message toReturn = new LYTellMessage(snapshotNumber ,getMessageType(), getOriginalSenderInfo(), getReceiverInfo(),
				false, getRoute(), getMessageText(), getMessageId(), getLYSnapshotResult());
		return toReturn;
	}*/

	public int getSnapshotNumber() {
		return snapshotNumber;
	}
}
