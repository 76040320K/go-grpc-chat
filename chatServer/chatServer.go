package chatServer

import (
	"context"
	"encoding/json"
	"go-grpc-chat/protoDir"
	"go-grpc-chat/redisDb"
	"log"
	"math/rand"
	"strings"
	// "sync"
	"time"
)

type messageUnit struct {
	ClientName string
	MessageBody string
	MessageUniqueCode int
	ClientUniqueCode int
}

// type messageQue struct {
// 	MQue []messageUnit
// 	// mu sync.Mutex
// }

var (
	// messageQueIntf = messageQue{}
	userCount = 0
)

type ChatServer struct {
	protoDir.ServicesServer
}

func (csstrt *ChatServer) ChatService(cs protoDir.Services_ChatServiceServer) error {
	clientUniqueCode := rand.Intn(1e3)
	userCount += 1
	go receiveFromStream(cs, clientUniqueCode)
	errch := make(chan error)

	go subscribeTopic(cs, clientUniqueCode, errch)
	return <- errch
}

var (
	ctx = context.Background()
)

func receiveFromStream(cs protoDir.Services_ChatServiceServer, clientUniqueCode int) {
	rdsCli := redisDb.RdsCli
	for {
		req, err := cs.Recv()
		if err != nil {
			// user out
			if strings.Contains(err.Error(), "Canceled desc = context canceled") {
				userCount -= 1
			} else {
				log.Panic(err)
			}
		} else {
			log.Print(req)
			messageQue := messageUnit{
				ClientName: req.Name,
				MessageBody: req.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode: clientUniqueCode,
			}
			payload, err := json.Marshal(messageQue)
			if err != nil {
				log.Panic(err)
			}
			if err := rdsCli.Publish(ctx, "user-chat-data", payload).Err(); err != nil {
				log.Panic(err)
			}
		}
	}
}

func sendToStream(cs protoDir.Services_ChatServiceServer, messageQue messageUnit, clientUniqueCode int, errch chan error) {
	time.Sleep(200 * time.Millisecond)

	err := cs.Send(&protoDir.FromServer{Name: messageQue.ClientName, Body: messageQue.MessageBody})

	if err != nil {
		errch <- err
	}
}

func subscribeTopic(cs protoDir.Services_ChatServiceServer, clientUniqueCode int, errch chan error) {
	for {
		rdsCli := redisDb.RdsCli
		subscriber := rdsCli.Subscribe(ctx, "user-chat-data")

		var messageQue messageUnit

		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				errch <- err
			}

			if err := json.Unmarshal([]byte(msg.Payload), &messageQue); err != nil {
				errch <- err
			}

			sendToStream(cs, messageQue, clientUniqueCode, errch)
		}
	}
}