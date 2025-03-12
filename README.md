# Golang CLI for Claude Sonnet

A playful Golang CLI for sending user input to Claude, with optional colorized output
and structured output of chatbot responses for stdin pipes (|)

## Run

```bash
go run ./cmd is golang better than C?
```

## Install

```bash
go install ./cmd
export PATH="$PATH:$HOME/go/bin"
```

## Simple Usage

```bash
goclaude is golang better than C?
```

## Structured output and piping

```bash
goclaude is golang better than C? -o json | jq
```

## User guide

```bash
goclaude -h
```
