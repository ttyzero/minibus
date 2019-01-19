# minibus

🚐 Minibus is a simple UNIX domain socket-based message bus, 
specifically for plain text communication between terminal 
UIs / shell scripts and terminal multiplexers.


## Messages

Messages are delivered to minibus as datagram packets containing utf8 strings
prefixed with a channel name, followed by a `:`. Messages are delivered to any
listeners on the channel specified. If a message sent to minibus does not
match the `chan: msg` format it will be dropped.

Messages are delivered by sending a datagram packet to the `minibus` UNIX 
datagram socket file in the minibus working directory, usually `~/.cache/minibus` 
on linux, or `~/Library/Caches/minibus` on macOS. 

Some shell commands capable of making this request 
are `nc` & `socat`, also see `tzmsg` documented below.

## Channels

Channels are specified in two locations, as the prefix to a Message, and as the
channel component of a socket connection from a listening process. Messages sent
to Minibus with a $CHANNEL prefix will be delivered to any Client Connection 
bearing the channel name.


## Client connections 

To listen to a channel, a process should establish a datagram socket in the minubus working 
directory with the name `$PID-$CHANNEL` where PID is the process ID of the 
connecting process, and CHANNEL is the name of the channel to recieve Messages
for.

## tzmsg

`tzmsg` is a companion command to Minibus which will stream anything on stdin to
a given Minibus channel. It can also be used in single message mode to send one 
message to a channel.

#### Using tzmsg

tzmsg with only one argument assumes the argument is a channel and accepts
stdin and sends each 'line' of input to the channel via the `minibus`
datagram socket.

tzmsg with multiple arguments assumes the first argument is the channel and 
all subsequent arguments constitute the message.

Any trailing colon on the end of the first argument is ignored, thus the channel
delimiter is optional.

#### Examples

Redirecting the output of a command to a minibus channel 'foo':

```bash
./your-server | tzmsg foo
```

Sending a single message to a minibus channel 'bar':

```bash
tzmsg bar 'this is a message'
```

```bash
tzmsg bar: this is a message
```

<br/><br/>
<table>
<tr><td>
<img src='https://raw.githubusercontent.com/ttyzero/logo/master/assets/ttyzero_animated.png' alt='ttyZero Logo' title='ttyZero Logo'/>
</td>
<td style='padding-left: 10em'>
<h2>Minibus is part of the <a href='http://github.com/ttyzero'>ttyZero Project</a></h2>
<b>Minibus</b> is <i>(c) 2019 ttyZero authors</i> <br/>
 and is available under the <b>MIT license</b>. 
</td></tr>
</table>
