# Hello World CLI

A simple command-line application written in Go using only the standard library.

## Features

- Greet a person by name
- Convert greeting to uppercase
- Repeat the greeting multiple times

## Usage

Build the application:

```bash
go build -o hello
```

Run with default options:

```bash
./hello
# Output: Hello, World!
```

Greet a specific person:

```bash
./hello -name Alice
# Output: Hello, Alice!
```

Use uppercase:

```bash
./hello -name Bob -uppercase
# Output: HELLO, BOB!
```

Repeat the greeting:

```bash
./hello -name Charlie -repeat 3
# Output:
# Hello, Charlie!
# Hello, Charlie!
# Hello, Charlie!
```

Combine options:

```bash
./hello -name Go -uppercase -repeat 2
# Output:
# HELLO, GO!
# HELLO, GO!
```

## Flags

- `-name` - Name to greet (default: "World")
- `-uppercase` - Print greeting in uppercase (default: false)
- `-repeat` - Number of times to repeat the greeting (default: 1)

## Help

View all available flags:

```bash
./hello -h
```
