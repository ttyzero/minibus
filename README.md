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
are `nc` & `socat`, also see `tzpipe` documented below.

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

## See also

[tzpipe](https://github.com/ttyzero/tzpipe) is a companion commandline utility for minibus which
allows shell scripts and commandline access to minibus messages (send / recieve)

[minibus-go](https://github.com/ttyzero/minibus-go) is a Golang client library for communicating
with a minibus service. It is designed for use in Golang TUI programs and serves as a reference
implementation for minibus clients in other languages.


<br/><br/>
<table>
<tr><td>
<img src='https://raw.githubusercontent.com/ttyzero/logo/master/assets/ttyzero_animated.png' alt='ttyZero Logo' title='ttyZero Logo'/>
</td>
<td style='padding-left: 10em'>
<h2>Minibus is part of the <a href='http://github.com/ttyzero'>ttyZero Project</a></h2>
<b>Minibus</b> is <i>(c) 2018 ttyZero authors</i> <br/>
 and is available under the <b>MIT license</b>. 
</td></tr>
</table>
