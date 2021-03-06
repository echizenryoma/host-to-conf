package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var inputPath string
var dnsConfigPath string
var proxyConfigPath string
var ipv6ConfigPath string
var dns string
var ipset string

func init() {
    flag.StringVar(&inputPath, "in", "hosts", "Input hosts file path")
    flag.StringVar(&dnsConfigPath, "do", "gfw_dns.conf", "Out dns config file path")
    flag.StringVar(&proxyConfigPath, "po", "gfw_proxy.conf", "Out proxy config file path")
    flag.StringVar(&ipv6ConfigPath, "v6", "gfw_ipv6.conf", "Out ipv6 config file path")
    flag.StringVar(&dns, "dns", "127.0.0.1#65053", "DNS Server")
    flag.StringVar(&ipset, "ipset", "proxy", "IPSet")
}


func main() {
    flag.Parse()

	inFile, err := os.Open(inputPath)
	fmt.Println(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	dnsFile, err := os.Create(dnsConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dnsFile.Close()

	proxyFile, err := os.Create(proxyConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer proxyFile.Close()

	ipv6File, err := os.Create(ipv6ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer ipv6File.Close()

	scanner := bufio.NewScanner(inFile)
	hosts := false
	var domainMap map[string]int
	for scanner.Scan() {
		line := scanner.Text()

		if ok, _ := regexp.MatchString("^#.*(End)$", line); ok {
			var domains []string
			for domain := range domainMap {
				domains = append(domains, domain)
			}
			sort.Strings(domains)
			for _, domain := range domains {
				fmt.Fprintln(dnsFile, fmt.Sprintf("server=/.%s/%s", domain, dns))
				fmt.Fprintln(proxyFile, fmt.Sprintf("ipset=/.%s/%s", domain, ipset))
				fmt.Fprintln(ipv6File, fmt.Sprintf("address=/.%s/::", domain))
			}
			fmt.Println(line)
			fmt.Fprintln(dnsFile, line)
			fmt.Fprintln(proxyFile, line)
			fmt.Fprintln(ipv6File, line)

			hosts = false
			domainMap = map[string]int{}
			continue
		}

		if ok, _ := regexp.MatchString("^#.*(Start)$", line); ok {
			fmt.Println(line)
			fmt.Fprintln(dnsFile, line)
			fmt.Fprintln(proxyFile, line)
			fmt.Fprintln(ipv6File, line)

			hosts = true
			domainMap = map[string]int{}
			continue
		}

		if len(line) > 0 && hosts {
			line := strings.Replace(line, "\t", " ", -1)
			domain := strings.Split(line, " ")[1]
			if ok, _ := regexp.MatchString("^localhost$", domain); ok {
				continue
			}

			eSLD, err := publicsuffix.EffectiveTLDPlusOne(domain)
			if err != nil {
				eSLD = domain
			}
			_, exist := domainMap[eSLD]
			if !exist {
				domainMap[eSLD] = 1
			}
		} else {
			fmt.Println(line)
			fmt.Fprintln(dnsFile, line)
			fmt.Fprintln(proxyFile, line)
			fmt.Fprintln(ipv6File, line)
		}
	}
}
