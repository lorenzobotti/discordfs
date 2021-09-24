package discordfs

import (
	"io"

	dg "github.com/bwmarrin/discordgo"
)

// praticamente io posso chiedere all'API
// 100 messaggi alla volta. metti che il mio si
// trova al novantesimo, mi basta una richiesta
// sola no? se è al centounesimo me ne servono due
// per minimizzare il numero di richieste che faccio
// creo sta struct che tiene traccia per me
// io gli chiedo un messaggio alla volta e se
// ce l'ha me lo da, se no chiama l'API per farsene
// dare altri. questo approccio credo si chiami
// "iterator pattern" e in Rust è proprio parte del
// linguaggio ed è molto più comodo, qua devo
// un pò improvvisare
type messageIterator struct {
	s         *dg.Session
	channelId string

	mess []*dg.Message
	at   int

	lastMessageId string
	finished      bool
}

// creo un nuovo MessageIterator
func newMessageIterator(s *dg.Session, id string) messageIterator {
	return messageIterator{
		s:         s,
		channelId: id,

		mess:     []*dg.Message{},
		at:       0,
		finished: false,
	}
}

func (mi *messageIterator) next() (*dg.Message, error) {
	if mi.at >= len(mi.mess) {
		err := mi.loadMessages()
		if err != nil {
			return nil, err
		}
	}

	if mi.at >= len(mi.mess)-1 {
		return nil, io.EOF
	}

	mess := mi.mess[mi.at]
	mi.at += 1
	return mess, nil
}

const messageLimit = 100

func (mi *messageIterator) loadMessages() error {
	mess, err := mi.s.ChannelMessages(mi.channelId, messageLimit, mi.lastMessageId, "", "")
	if err != nil {
		return err
	}

	mi.mess = append(mi.mess, mess...)

	last := len(mess) - 1
	if last >= 0 {
		mi.lastMessageId = mess[last].ID
	} else {
		mi.finished = true
	}

	return nil
}
