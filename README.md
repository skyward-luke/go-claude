# Golang CLI for Claude Chatbot API

A playful Golang CLI for sending user input to Claude, with optional colorized output
and structured output of chatbot responses for stdin pipes (|)

Serializes the conversation history to file, in protocol buffer format.


## Run

```bash
MSG="is golang better than C?"
go run ./cmd ${MSG} -memory-id 1
MSG="is C better than Rust?"
go run ./cmd ${MSG} -memory-id 1
echo ${MSG} | go run ./cmd

# count input tokens in case history is long
go run ./cmd ${MSG} -memory-id 1 -count
```

## Install

```bash
go install ./cmd
export PATH="$PATH:$HOME/go/bin"
```

## Simple Usage

```bash
echo ${MSG} | goclaude
goclaude ${MSG}
```

## Structured output and piping

```bash
echo ${MSG} | goclaude -o json | jq
```

## User guide

```bash
goclaude -h
```

## Memory

Chat history is saved to file in protocol buffer format.

### Working with Protocol Buffers

Generate pb code in pb/chat:

`protoc -I=proto/ --go_out=. chat.proto`