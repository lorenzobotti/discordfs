package discordfs

import (
	"io"

	dg "github.com/bwmarrin/discordgo"
)

// bad implementation of an iterator pattern. can you tell i've been using Rust?
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
