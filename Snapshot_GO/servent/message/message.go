package message

type Message interface {
	// Alters the message(returns copy) so we are the sender.
	// Everything else stays intact, used for rerouting
	// the middle nodes.
	MakeMeASender() Message
	// Alters the message(returns copy) to change receiver.
	// Everything else stays intact, used for rerouting
	// the middle nodes.
	ChangeReceiver(newReceiverId int) Message
	// Alters the message(returns copy) to change color.
	// White - message was sent before local snapshot.
	// Red - message was sent after local snapshot.
	SetRedColor() Message
	// Alters the message(returns copy) to change color.
	// White - message was sent before local snapshot.
	// Red - message was sent after local snapshot.
	SetWhiteColor() Message
	SendEffect()
}
