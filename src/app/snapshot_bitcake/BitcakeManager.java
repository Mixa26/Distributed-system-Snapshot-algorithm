package app.snapshot_bitcake;

import java.util.Set;

/**
 * Describes a bitcake manager. These classes will have the methods
 * for handling snapshot recording and sending info to a collector.
 * 
 * @author bmilojkovic
 *
 */
public interface BitcakeManager {

	public void takeSomeBitcakes(int amount);
	public void addSomeBitcakes(int amount);

	public void addIdBorderSet(Set<Integer> idBorderSet);
	public int getCurrentBitcakeAmount();
	
}
