package servent

import (
	"Snapshot/app"
	"Snapshot/app/snapshotbitcake"
	"Snapshot/helper"
	"Snapshot/servent/handler"
	handlersnapshot "Snapshot/servent/handler/snapshot"
	"Snapshot/servent/message"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type SimpleServentListener struct {
	working           bool
	snapshotCollector snapshotbitcake.SnapshotCollector
	redMessages       []message.Message
	threadPool        *helper.Pool
}

func ConstructSimpleServentListener(snapshotCollector snapshotbitcake.SnapshotCollector) *SimpleServentListener {
	threadPool := helper.InitiatePool()
	return &SimpleServentListener{true, snapshotCollector, nil, threadPool}
}

func (simpleServentListener SimpleServentListener) Run(wgListener *sync.WaitGroup) {
	listener, err := net.Listen("tcp", app.MyServentInfo.Ip+":"+strconv.Itoa(app.MyServentInfo.Port))
	app.TimeStampPrint("Now waiting for connection on: " + strconv.Itoa(app.MyServentInfo.Port))
	wgListener.Done()

	if err != nil {
		log.Println(err)
		app.ExitWithErrorMsg("ERROR ESTABLISHING SERVENT LISTENER ON PORT: " + strconv.Itoa(app.MyServentInfo.Port))
	}
	defer listener.Close()

	for {
		// If working is set to false, drop out.
		if !simpleServentListener.working {
			break
		}

		var clientMessage message.Message = nil
		var connection net.Conn

		if !app.IsWhite.Load() && len(simpleServentListener.redMessages) > 0 {
			// Remove first element.
			clientMessage = simpleServentListener.redMessages[0]
			simpleServentListener.redMessages = simpleServentListener.redMessages[1:]
		} else {
			// Wait max for 1s for listener to accept.
			listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))
			// Listener will wait max 1s, and then repeat the for loop.
			connection, err = listener.Accept()

			if err != nil {
				log.Println(err)
				// The connection didn't receive any calls, so
				// we go back to the top to check if working == true.
				continue
			}
			app.TimeStampPrint("Received connection on : " + strconv.Itoa(app.MyServentInfo.Port))

		}

		go func() {
			if clientMessage != nil {
				clientMessage = snapshotbitcake.ReadMessage(connection)
			}

			if clientMsg, ok := clientMessage.(*message.BasicMessage); ok {
				app.ColorLock.Lock()

				if !clientMsg.White && app.IsWhite.Load() {
					/*
					* If the message is red, we are white, and the message isn't a marker,
					* then store it. We will get the marker soon, and then we will process
					* this message. The point is, we need the marker to know who to send
					* our info to, so this is the simplest way to work around that.
					 */
					if clientMsg.MessageType != "LY_MARKER" {
						simpleServentListener.redMessages = append(simpleServentListener.redMessages, clientMsg)
						return
					} else {
						if lybm, ok := simpleServentListener.snapshotCollector.GetBitcakeManager().(*snapshotbitcake.LaiYangBitcakeManager); ok {
							tmp, err := strconv.Atoi(clientMsg.MessageText)

							if err != nil {
								app.TimeStampErrorPrint("Error parsing integer: " + clientMsg.MessageText)
							}

							lybm.MarkerEvent(tmp, simpleServentListener.snapshotCollector)
						} else {
							app.TimeStampErrorPrint("Couldnt get bitcake manager to LaiYangBitcakeManager.")
						}
					}

					app.ColorLock.Unlock()

					var messageHandler handler.MessageHandler = handler.NullHandler{ClientMessage: clientMsg}

					switch clientMsg.MessageType {
					case "TRANSACTION":
						messageHandler = snapshotbitcake.TransactionalHandler{ClientMessage: clientMsg, BitcakeManager: simpleServentListener.snapshotCollector.GetBitcakeManager()}
					case "LY_MARKER":
						messageHandler = handlersnapshot.LYMarkerHandler{}
					case "LY_TELL":
						messageHandler = handlersnapshot.LYTellHandler{ClientMessage: clientMsg, SnapshotCollector: simpleServentListener.snapshotCollector}
					}

					// Add the message handler to the thread pool.
					simpleServentListener.threadPool.JobQueue <- helper.Job{Handler: messageHandler}
				}
			} else {
				app.TimeStampErrorPrint("Message wasn't a BasicMessage!")

			}
		}()
	}
}

func (simpleServentListener *SimpleServentListener) Stop() {
	(*simpleServentListener).working = false
}
