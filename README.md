# Golang CLI for Claude Sonnet

A playful Golang CLI for sending user input to Claude, with optional colorized output
and structured output of chatbot responses for stdin pipes (|)

## Run

```bash
MSG="is golang better than C?"
echo ${MSG} | go run ./cmd
```

## Install

```bash
go install ./cmd
export PATH="$PATH:$HOME/go/bin"
```

## Simple Usage

```bash
echo ${MSG} | goclaude
```

## Structured output and piping

```bash
echo ${MSG} | goclaude -o json | jq
```

## User guide

```bash
goclaude -h
```
