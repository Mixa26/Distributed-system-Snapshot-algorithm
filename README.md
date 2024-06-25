# Distributed-system-Snapshot-algorithm

## What is a distributed system?

Distributed systems are composed of multiple servers or hardware units connected over the internet.<br>
They use internet to communicate and exchange messages. As you can imagine a lot of synchronization<br>
problems can arise in communication of such servers. This the task of this project.<br>

## What is a snapshot?

A snapshot is a moment where a distributed system "records" it's current state.<br>

## Implementation specifics

This project features a Lai-Yang algorithm implementation for multiple snapshot initializations,<br>
as well as Spezialleti-Kearns implementation for combining the results of multiple simultaneous snapshots.<br>
The servents (servent and client) in our system exchange bitcake, which is an imagined currency.<br>
The point of this project is to preserve this currency and that none of the bitcake is lost during a snapshot.<br>

## Output, input and errors

To change the configuration of the system _in.txt files can be changed in the ly_snapshot/input folder. Output can be<br>
inspected in the output folder and errors in the error folder. The system configuration is read from servent_list.properties<br>
and it can be changed. Here is what it looks like:<br>

```
servent_count=7
clique=false
snapshot=ly
servent0.port=1100
servent1.port=1200
servent2.port=1300
servent3.port=1400
servent4.port=1500
servent5.port=1600
servent6.port=1700
servent0.neighbors=1,2,6
servent1.neighbors=0,3,2,5
servent2.neighbors=0,3,1,6
servent3.neighbors=1,2,4,6
servent4.neighbors=3,5
servent5.neighbors=1,4,6
servent6.neighbors=0,2,3,5
initiators=4,6,1
```

It is self explanatory.<br>
Here is an example of commands for the _in.txt file:<br>

```
pause 300
bitcake_info
transaction_burst
pause 15000
transaction_burst
pause 15000
transaction_burst
pause 15000
transaction_burst
pause 15000
stop
```
