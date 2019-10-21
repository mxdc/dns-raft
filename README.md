# dns-raft

DNS cluster using Raft protocol for resource records replication.

The purpose of this case study is to implement the [Raft](https://raft.github.io/) library from [Hashicorp](https://github.com/hashicorp/raft) in order to maintain consistent DNS records across multiple machines.

```
                             ┌────────┐
                             │zone.txt│
                             └────▲───┘
                                  │
                             read │
                                  │
   ┌──────────┐             ┌─────┴────┐             ┌──────────┐
   │node 02   │             │node 01   │             │node 03   │
   │          │             │          │             │          │
   │follower  ◀─────────────┤leader    ├─────────────▶follower  │
   │          │ replication │          │ replication │          │
   │          │             │          │             │          │
   └──────────┘             └──────────┘             └──────────┘
```

## Build

Compile source code:
```
$ go build -o bin/dns-raft cmd/main.go
```

## Run

Start three nodes:
```
$ bin/dns-raft -id id0 -raft.addr ":8300" -dns.addr ":8600" -zone.file "./zones/zone.txt"
$ bin/dns-raft -id id1 -raft.addr ":8301" -dns.addr ":8601" -raft.join ":8300"
$ bin/dns-raft -id id2 -raft.addr ":8302" -dns.addr ":8602" -raft.join ":8300"
```

## DNS

Resources records are loaded from [zone file](zones/zone.txt) at execution.

Resolve address from first node:
```
$ dig @127.0.0.1 -p 8600 example.com
```

Resolve address from second node:
```
$ dig @127.0.0.1 -p 8601 example.com
```

Resolve address from third node:
```
$ dig @127.0.0.1 -p 8602 example.com
```

Add a DNS record:
```
echo 'database                 60 A     1.2.3.7' >> zones/zone.txt
```

Reload zone file by sending SIGHUP to leader node:
```
$ pkill -SIGHUP dns-raft
```

Resolve new address from follower node:
```
$ dig @127.0.0.1 -p 8602 database.example.com
```

## Play with KV Store

Ping the first node:
```
$ echo "kv ping" | nc localhost 8300
PONG
```

Add a key:
```
$ echo "kv set toto titi" | nc localhost 8300
SUCCESS
```

Get a key from any node:
```
$ echo "kv get toto" | nc localhost 8300
titi
$ echo "kv get toto" | nc localhost 8301
titi
$ echo "kv get toto" | nc localhost 8302
titi
```

Remove the key:
```
$ echo "kv del toto" | nc localhost 8500
```

## Inspirations

* http://www.scs.stanford.edu/17au-cs244b/labs/projects/orbay_fisher.pdf
* https://github.com/hashicorp/consul
* https://github.com/otoolep/hraftd
* https://github.com/yongman/leto

## todo

* forward to leader set command
* load zone file on every node
