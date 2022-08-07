package main

import (
	"bufio"
	// "fmt"
	// "fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type hostTarget struct {
    name string
    lookupIP []string
    cname string
    lookupNS []string
    // lookupTXT []string
}
//String slice of valid URLs
var urls = make([]string, 0)
var hostTargets = make([]hostTarget, 0)

//Open a file and make a slice from the contents of the file.
func readFile(file string) []string {
    urls = make([]string, 0)
    f, err := os.Open(file)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        urls = append(urls, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return urls
}





func main() {
    /* urls = append(urls, "http://www.google.com", "http://www.yahoo.com", "http://www.bing.com")
    urls = append(urls, "http://www.yahoo.com")
    urls = append(urls, "http://www.amazon.com")
    urls = append(urls, "http://www.facebook.com")
    urls = append(urls, "http://www.twitter.com") */
    urls = readFile("./urls.txt")
    //Concurrently validateUrl each URL in the slice
    for _, u := range urls {
        ch :=make (chan bool)
        go validateUrl(u, ch)
        result := <-ch
        log.Printf("%s is %v\n", u, result)
    }
/*
    for _, u := range urls {
        tar := dnsLookUp(u)
        hostTargets = append(hostTargets, tar)
    }

    for _, ht := range hostTargets {
        fmt.Printf("Name: %s\n", ht.name)
        fmt.Printf("IP: %s\n", ht.lookupIP)
        fmt.Printf("CNAME: %s\n", ht.cname)
        fmt.Printf("NS: %s\n", ht.lookupNS)
        // fmt.Printf("TXT: %s\n", ht.lookupTXT)
        fmt.Println("")
    } */

}

//Function to perform DNS lookup and return struct with results
func dnsLookUp (target string) hostTarget {
    var ht hostTarget
    ht.name = target
    lookupIPs, err := net.LookupIP(target)
    if err != nil {
        // log.Fatal(err)
        ht.lookupIP = append(ht.lookupIP, "Error")
    } else {
        for _, ip := range lookupIPs {
            ht.lookupIP = append(ht.lookupIP, ip.String())
        }
    }
    cnames, err := net.LookupCNAME(target)
    if err != nil {
        // log.Fatal(err)
        ht.cname = "Error"
    } else {
          ht.cname = cnames
    }
    lookupNS, err := net.LookupNS(target)
    if err != nil {
        // log.Fatal(err)
        ht.lookupNS = append(ht.lookupNS, "Error")
    } else {
        for _, ns := range lookupNS {
            ht.lookupNS = append(ht.lookupNS, ns.Host)
        }
    }
    return ht
}

//Create a function that validates a URL and returns true if it is valid and false if it is invalid.
func validateUrl(u string, ch chan bool) {
    var result = false
    _, err := http.Get(u)
    if err == nil {
        result = true
    }
    ch <- result
}
