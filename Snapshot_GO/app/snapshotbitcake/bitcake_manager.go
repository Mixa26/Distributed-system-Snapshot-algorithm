package snapshotbitcake

type BitcakeManager interface {
	TakeSomeBitcakes(amount int)
	AddSomeBitcakes(amount int)
	GetCurrentBitcakeAmount() int
}
