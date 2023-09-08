package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/miekg/dns"
)

func randomIP() string {
	return fmt.Sprintf("172.20.%d.%d", rand.Intn(20), 1+rand.Intn(253))
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			for i := 0; i < 8; i++ {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, randomIP()))
				rr.Header().Ttl = 5
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	// attach request handler func
	dns.HandleFunc("test.org", handleDnsRequest)

	// start server
	port := 53
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
