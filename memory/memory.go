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

func Get(memoryid int32) (*chat.Memory, error) {
	memories, err := read()
	if err != nil {
		return nil, err
	}

	memory := memories.Batch[memoryid]
	if memory == nil {
		memory = &chat.Memory{}
	}
	return memory, nil
}

func Save(content string, role Role, memoryid int32) {
	memories, err := read()
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

	if err := write(memories); err != nil {
		log.Fatalln("failed to save memories:", err)
	}
}

func read() (*chat.Memories, error) {
	in, err := os.ReadFile("memories.bin")
	if err != nil {
		return nil, err
	}

	memories := &chat.Memories{}
	if err := proto.Unmarshal(in, memories); err != nil {
		return nil, err
	}
	return memories, nil
}

func write(memories *chat.Memories) error {
	out, err := proto.Marshal(memories)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("memories.bin", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
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
