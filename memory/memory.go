package memory

import (
	"chat"
	"log"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate stringer -type=Role

type Role byte

const (
	User Role = iota
	Assistant
)

type SavedMemories struct {
	FilePath string
}

func (m *SavedMemories) Get(memoryid int32) (*chat.Memory, error) {
	memories, err := m.read()
	if err != nil {
		return nil, err
	}

	memory := memories.Batch[memoryid]
	if memory == nil {
		memory = &chat.Memory{}
	}
	return memory, nil
}

func (m *SavedMemories) Save(content string, role Role, memoryid int32) {
	memories, err := m.read()
	if err != nil {
		memories = &chat.Memories{}
		log.Println(err)
	}

	memory := memories.Batch[memoryid]
	if memory == nil {
		memory = &chat.Memory{}
	}
	if len(memories.Batch) == 0 {
		memories.Batch = make(map[int32]*chat.Memory)
	}

	now := timestamppb.Now()
	msg := &chat.ChatMessage{Role: strings.ToLower(role.String()), Content: content, Ts: now}

	memory.ChatMessages = append(memory.ChatMessages, msg)
	memory.LastUsed = now
	memories.Batch[memoryid] = memory

	log.Println(memory)
	log.Println(memories)

	if err := m.write(memories); err != nil {
		log.Fatalln("failed to save memories:", err)
	}
}

func (m *SavedMemories) fileExists() bool {
	_, err := os.Stat(m.FilePath)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

func (m *SavedMemories) read() (*chat.Memories, error) {
	if !m.fileExists() {
		return &chat.Memories{}, nil
	}

	in, err := os.ReadFile(m.FilePath)
	if err != nil {
		return nil, err
	}

	memories := &chat.Memories{}
	if err := proto.Unmarshal(in, memories); err != nil {
		return nil, err
	}
	return memories, nil
}

func (m *SavedMemories) write(memories *chat.Memories) error {
	out, err := proto.Marshal(memories)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(m.FilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(out))
	if err != nil {
		return err
	}
	return nil
}
