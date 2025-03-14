package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/mxdc/dns-raft/dns"
	"github.com/mxdc/dns-raft/store"
)

var (
	dnsAddr  string
	raftAddr string
	raftJoin string
	raftId   string
	zoneFile string
)

func init() {
	flag.StringVar(&dnsAddr, "dns.addr", ":5350", "DNS listen address")
	flag.StringVar(&raftAddr, "raft.addr", ":15370", "Raft bus transport bind address")
	flag.StringVar(&raftJoin, "raft.join", "", "Join to already exist cluster")
	flag.StringVar(&raftId, "id", "", "node id")
	flag.StringVar(&zoneFile, "zone.file", "", "Zone file containing resource records")
}

func main() {
	flag.Parse()

	kvs := store.InitStore(raftAddr, raftJoin, raftId)
	dns := dns.NewDNS(kvs, dnsAddr)
	go handleSignals(kvs, dns)
	dns.LoadZone(zoneFile)
	dns.Start()
}

func handleSignals(kvs *store.Store, dns *dns.DNS) {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sighupChan:
			dns.LoadZone(zoneFile)
		case <-signalChan:
			kvs.Leave(kvs.RaftID)
			dns.Shutdown()
			kvs.Shutdown()
		}
	}
}
