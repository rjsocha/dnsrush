package main

// quick & dirty dns benchmark tool 
// for internal use only
// one day I will learn go... one day

// v 0.1
// Public Domain

import (
	"fmt"
	"time"
	"net"
	"os"
	"flag"
	"strings"
	"bufio"
	"math/rand"
	"github.com/miekg/dns"
)
var (
	a_ns         = flag.String("ns", "", "nameserver IPv4 address (required)")
	a_query      = flag.String("q", "", "query")
	a_qt         = flag.String("qt", "", "query type")
	a_count	     = flag.Int("c",1,"number of queries")
	a_playlist   = flag.String("playlist","","playlist input file")
	a_playmode   = flag.String("mode","random","playlist mode: random/sequential")
	a_timeout    = flag.Int("t",25,"connection timeout (in ms)")
	a_verify     = flag.Bool("verify",false,"verify playlist")
)

type Playlist struct {
  query string
  qtype uint16
}

var (
     timeout time.Duration
     pl []Playlist
)

func str2type(s string) uint16 {
 t:=uint16(0)
 switch s=strings.ToLower(s); s {
	case "a":
		t=1
	case "ns":
		t=2
	case "cname":
		t=5
	case "soa":
		t=6
	case "mx":
		t=15
	case "txt":
		t=16
	case "aaaa":
		t=28
	case "any":
		t=255
 }
 return t
}

func readPlaylist(f string) (bool) {
	fh,err:=os.Open(f)
	if(err!=nil) {
		fmt.Fprintf(os.Stderr,"%v\n",err)
		os.Exit(1)
	}
	defer fh.Close()
	scanner:=bufio.NewScanner(fh)
	scanner.Split(bufio.ScanWords)
	il:=0
	q:=""
	t:=""
	r:=""
	origin:=""
	line:=0
	for scanner.Scan() {
		r=scanner.Text()
		if(il==0) {
			q=r
		} else {
			line++
			t=r
			qt:=str2type(t)
			if(qt==0) {
				fmt.Fprintf(os.Stderr,"Unable to parse line: %d\n",line)
				os.Exit(1)
			}
			if(string(q[len(q)-1:]) != ".") {
				if(len(origin)>0) {
					if(q=="@") {
						q=origin
					} else {
						q=q+"."+origin
					}
				} else {
					fmt.Fprintf(os.Stderr,"ORIGIN not set at line: %d\n",line)
					os.Exit(1)
				}
			} else {
			   origin=q
			}
			item:=Playlist{query: q, qtype: qt}
			pl = append(pl,item)
		}
		il=(il+1)&1
	}
	if err:=scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr,"Error reading input file: %v\n",err)
		os.Exit(1)
	}
	return true
}

func exeQ(q string, t uint16,ns string) (code int, rcode int, start int64, rtt int64)  {
	m := new(dns.Msg)
	m.Id = dns.Id()
	m.RecursionDesired = false
	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{q, t, dns.ClassINET}
	c := new(dns.Client)
        c.SingleInflight = true
        c.Dialer = &net.Dialer{Timeout: timeout * time.Millisecond}
        start = time.Now().UnixNano()
	res, qrtt, err := c.Exchange(m,ns)
	rtt = int64(qrtt)
	if err != nil {
		return 1,0,start,rtt
	} else if res == nil {
		return 2,0,start,rtt
	}
	return 0,res.Rcode,start,rtt
}

// random
func play_mode1(ns string) {
	s1 := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(s1)
	l:=len(pl)
	s:=0
	for i:=0; i<*a_count;i++ {
		s=rnd.Intn(l)
		q:=pl[s].query
		qt:=pl[s].qtype
		code,rcode,start,rtt:=exeQ(q,qt,ns)
		fmt.Printf("%d %d %d %d %s %d\n",code,rcode,int64(start),int64(rtt),q,qt)
	}
}

// sequential
func play_mode2(ns string) {
	l:=len(pl)
	s:=0
	for i:=0; i<*a_count;i++ {
		q:=pl[s].query
		qt:=pl[s].qtype
		code,rcode,start,rtt:=exeQ(q,qt,ns)
		fmt.Printf("%d %d %d %d %s %d\n",code,rcode,int64(start),int64(rtt),q,qt)
		s=(s+1)%l
	}
}
func play(mode int,ns string) {
	if(mode==1) {
		play_mode1(ns)
	} else {
		play_mode2(ns)
	}
}
func usage() {
	fmt.Fprintf(os.Stderr,"Usage: %s \n",os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	timeout=time.Duration(*a_timeout)
	ns_addr := ""
	if(*a_verify==false) {
		if(*a_ns == "") {
			usage()
			os.Exit(1)
		}
		ns_ip := net.ParseIP(*a_ns)
		if ns_ip.To4() == nil {
			fmt.Fprintf(os.Stderr,"nameserver is not an IPv4 address\n")
			os.Exit(1)
		}
		ns_addr=*a_ns + ":53"
	}
	q:=""
	qt:=uint16(0)
	if(*a_playlist != "") {
		smode:=strings.ToLower(*a_playmode)
		mode:=0
		if(smode=="random") {
			mode=1
		} else if(smode=="sequential") {
			mode=2
		} else {
			fmt.Fprintf(os.Stderr,"unknow mode: %s\n",mode)
			os.Exit(1)
		}
		readPlaylist(*a_playlist)
		if(*a_verify) {
			os.Exit(0)
		}
		if(len(pl)<1) {
			fmt.Fprintf(os.Stderr,"error: empty playlist\n")
			os.Exit(1)
		}
		play(mode,ns_addr)
	} else {
		if(*a_qt == "") {
			fmt.Fprintf(os.Stderr,"query type not specified\n")
			os.Exit(1)
		}
		qt=str2type(*a_qt)
		if(qt==0) {
			fmt.Fprintf(os.Stderr,"Unknow query type: %s\n",*a_qt)
			os.Exit(1)
		}
		if(*a_query == "") {
			fmt.Fprintf(os.Stderr,"query not specified\n")
			os.Exit(1)
		}
		q=*a_query
		if(string(q[len(q)-1:]) != ".") {
			q=q+"."
		}
		for i:=0; i<*a_count;i++ {
			code,rcode,start,rtt:=exeQ(q,qt,ns_addr)
			fmt.Printf("%d %d %d %d %s %d\n",code,rcode,int64(start),int64(rtt),q,qt)
		}
	}
	os.Exit(0)
}
