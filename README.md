# gzog

This describes an experimental implementation of a simple protocol for
communications between a number of servers (daemons) organized in a full
mesh and one or more clients external to the ring.

For historical reasons the daemon is called `zog` and its Go
implementation `gzog`.  There is a client called `sprog`; its Go
implementation is of course `gsprog`.

Participants send messages, simple strings, to one another using
**tcp/ip**, one of the
[standard Internet protocols](https://en.wikipedia.org/wiki/internet_protocol_suite).

### Ring

<img align="right" src="img/fullMesh.png">
The system runs on a set of five or so machines, each housing a **zog daemon**
listening for messages from the other daemons on port 55551.
Servers are connected in a full mesh, seeing other machines as `+1`, `+2`,
etc on the ring.
Each server has to send a keep-alive
at least every 300s to keep a connection open.

Daemons log all messages to either `access.log` or `error.log`, with
CRLF terminated messages of reasonable length (<512 bytes) being
logged to `access.log` and all other messages reported to `error.log`.

Daemons also listen for traffic from clients on port 55552.

For the purpose of this exercise, the ring consists of five machine
with addresses in
[private address space](https:/tools.ietf.org/rfc1918), 192.168/8

* **MV**, 192.168.152.242
* **PA**, 192.168.152.252
* **G**,  192.168.152.254
* **LA**, 192.168.152.253
* **S**,  192.168.152.11

### Daemons

<TBD>

### Clients

Clients run **gsprog**, which talks to the gzog daemons on port 55552.
At least one gsprog client will run. By default the
client will send messages at configurable intervals to a randomly selected
server.  Command line switches will
be available to change this behavior, to round-robin (-r/--roundRobin) or
to target a specific host (`-s/--specificHost`).  If the testing option
is chosen, `specificHost` will be `127.0.0.1`.

The default interval should be 10s.

Gzog messages are simple strings,
[Unicode](http://unicode.org/standard/WhatIsUnicode.html)
encoded as
[UTF-8](https://en.wikipedia.org/wiki/UTF-8).

### Messages

In this implementation,
messages are standard Go strings preceded by a single character, a digit, the
**wrapping** byte.

The wrapping byte will be one of `0`, `1`, or `2`.  That is, it will be
a single ASCII character, a digit.

Content will be quasi-random with lengths (`-i/--interval`) in the range
`[32..63]` inclusive, ***not*** counting the wrapping byte.  Content will
consist of printing characters.  [Content length should be adjusted as
necessary to keep log lines within 80 characters.]

### Message Handling

On receiving a message other than a keep-alive from a client, a server will
forward it to its `+1` and `+2` peers with `1` and `2` **wrappings**
respectively.

#### Daemon-Daemon Messages

On receiving a message from a peer, a server will examine its envelope.
All messages will be logged AFTER any forwarding.

If the wrapping is `0` (the digit zero), the message will just be logged.  If the wrapping
is `1`, the message will be forwarded to the peer's `+3` server with a zero
wrapping.  If the wrapping is `2`, the message will be forwarded to the
`+4` server with a zero wrapping.

#### Client-Daemon Messages

On receiving a message from a client which is other than a keep-alive
a server will forward it to its `+1` and `+2` peers with a `2` wrapping
in each case.

If the wrapping is `0` (zero), the message will just be logged.

If the wrapping is `1`, the message will be forwarded to the receiver's
`+1` peer with a zero wrapping.  [We will not use this wrapping in the
current configuration.]

If the wrapping is `2`, the message will be forwarded to the
receiver's `+2` peer with a zero wrapping.

#### Protocol

gzog
[protocol](http://www.webopedia.com/TERM/P/protocol.html)
should be `org.xlattice.gzog`
and should at this point
define two message pairs: `keepAlive` and `ack`, and `data` and `ok`.

## Project Status

Pre-Alpha.

## On-line Documentation

More information on the **gzog** project can be found
[here](https://jddixon.github.io/gzog)
