**Result of nmap'ing our floating IP (188.166.193.132)**

How I did it:

`docker run --rm -it instrumentisto/nmap -A -T4 scanme.nmap.org -v -A -sV 188.166.193.132`

Result:

```
Starting Nmap 7.91 ( https://nmap.org ) at 2021-04-19 10:04 UTC
NSE: Loaded 153 scripts for scanning.
NSE: Script Pre-scanning.
Initiating NSE at 10:04
Completed NSE at 10:04, 0.00s elapsed
Initiating NSE at 10:04
Completed NSE at 10:04, 0.00s elapsed
Initiating NSE at 10:04
Completed NSE at 10:04, 0.00s elapsed
Initiating Ping Scan at 10:04
Scanning 2 hosts [4 ports/host]
Completed Ping Scan at 10:04, 0.01s elapsed (2 total hosts)
Initiating Parallel DNS resolution of 2 hosts. at 10:04
Completed Parallel DNS resolution of 2 hosts. at 10:04, 1.13s elapsed
Initiating SYN Stealth Scan at 10:04
Scanning 2 hosts [1000 ports/host]
Discovered open port 22/tcp on 188.166.193.132
Discovered open port 8080/tcp on 188.166.193.132
Discovered open port 22/tcp on 45.33.32.156
Discovered open port 9200/tcp on 188.166.193.132
Discovered open port 8081/tcp on 188.166.193.132
Completed SYN Stealth Scan against 188.166.193.132 in 0.50s (1 host left)
Discovered open port 80/tcp on 45.33.32.156
Discovered open port 31337/tcp on 45.33.32.156
Discovered open port 9929/tcp on 45.33.32.156
Completed SYN Stealth Scan at 10:04, 1.55s elapsed (2000 total ports)
Initiating Service scan at 10:04
Scanning 8 services on 2 hosts
Completed Service scan at 10:04, 11.08s elapsed (8 services on 2 hosts)
Initiating OS detection (try #1) against 2 hosts
Retrying OS detection (try #2) against 2 hosts
Initiating Traceroute at 10:05
Completed Traceroute at 10:05, 0.03s elapsed
Initiating Parallel DNS resolution of 2 hosts. at 10:05
Completed Parallel DNS resolution of 2 hosts. at 10:05, 0.03s elapsed
NSE: Script scanning 2 hosts.
Initiating NSE at 10:05
Completed NSE at 10:05, 6.62s elapsed
Initiating NSE at 10:05
Completed NSE at 10:05, 2.33s elapsed
Initiating NSE at 10:05
Completed NSE at 10:05, 0.00s elapsed
Nmap scan report for scanme.nmap.org (45.33.32.156)
Host is up (0.034s latency).
Other addresses for scanme.nmap.org (not scanned): 2600:3c01::f03c:91ff:fe18:bb2f
Not shown: 996 closed ports
PORT      STATE SERVICE    VERSION
22/tcp    open  ssh        OpenSSH 6.6.1p1 Ubuntu 2ubuntu2.13 (Ubuntu Linux; protocol 2.0)
| ssh-hostkey: 
|   1024 ac:00:a0:1a:82:ff:cc:55:99:dc:67:2b:34:97:6b:75 (DSA)
|   2048 20:3d:2d:44:62:2a:b0:5a:9d:b5:b3:05:14:c2:a6:b2 (RSA)
|   256 96:02:bb:5e:57:54:1c:4e:45:2f:56:4c:4a:24:b2:57 (ECDSA)
|_  256 33:fa:91:0f:e0:e1:7b:1f:6d:05:a2:b0:f1:54:41:56 (ED25519)
80/tcp    open  http       Apache httpd 2.4.7 ((Ubuntu))
|_http-favicon: Nmap Project
| http-methods: 
|_  Supported Methods: OPTIONS GET HEAD POST
|_http-server-header: Apache/2.4.7 (Ubuntu)
|_http-title: Go ahead and ScanMe!
9929/tcp  open  nping-echo Nping echo
31337/tcp open  tcpwrapped
OS fingerprint not ideal because: Host distance (11 network hops) is greater than five
No OS matches for host
Network Distance: 2 hops
Service Info: OS: Linux; CPE: cpe:/o:linux:linux_kernel

TRACEROUTE (using port 80/tcp)
HOP RTT     ADDRESS
-   Hop 1 is the same as for 188.166.193.132
2   0.54 ms scanme.nmap.org (45.33.32.156)

Nmap scan report for 188.166.193.132
Host is up (0.0047s latency).
Not shown: 996 closed ports
PORT     STATE SERVICE VERSION
22/tcp   open  ssh     OpenSSH 8.2p1 Ubuntu 4ubuntu0.2 (Ubuntu Linux; protocol 2.0)
| ssh-hostkey: 
|   3072 c0:03:75:6d:ac:9c:4c:88:41:c6:a6:3b:e4:67:6e:9e (RSA)
|   256 c7:d6:ce:66:11:3d:2d:01:33:aa:b0:de:e1:91:b5:56 (ECDSA)
|_  256 9f:0d:0e:1a:22:b6:53:32:e0:67:d2:32:a6:b5:5a:c7 (ED25519)
8080/tcp open  http    Golang net/http server (Go-IPFS json-rpc or InfluxDB API)
| http-methods: 
|_  Supported Methods: GET
|_http-open-proxy: Proxy might be redirecting requests
|_http-title:  Public Timeline | MiniTwit
8081/tcp open  http    Golang net/http server (Go-IPFS json-rpc or InfluxDB API)
|_http-title: Site doesn't have a title (text/plain; charset=utf-8).
9200/tcp open  http    Elasticsearch REST API 7.12.0 (name: elasticsearch; cluster: docker-cluster; Lucene 8.8.0)
|_http-favicon: Unknown favicon MD5: 6177BFB75B498E0BB356223ED76FFE43
| http-methods: 
|   Supported Methods: HEAD GET DELETE OPTIONS
|_  Potentially risky methods: DELETE
|_http-title: Site doesn't have a title (application/json; charset=UTF-8).
Device type: printer|print server|switch
Running (JUST GUESSING): HP embedded (86%), Dell embedded (85%)
OS CPE: cpe:/h:hp:designjet_650c cpe:/h:hp:jetdirect_170x cpe:/h:dell:powerconnect_5424
Aggressive OS guesses: HP DesignJet 650C printer (86%), HP 170X print server or Inkjet 3000 printer (85%), Dell PowerConnect 5424 switch (85%), HP LaserJet 4000 printer (85%)
No exact OS matches for host (test conditions non-ideal).
Network Distance: 2 hops
TCP Sequence Prediction: Difficulty=135 (Good luck!)
IP ID Sequence Generation: Randomized
Service Info: OS: Linux; CPE: cpe:/o:linux:linux_kernel

TRACEROUTE (using port 80/tcp)
HOP RTT     ADDRESS
1   0.07 ms 172.17.0.1
2   0.56 ms 188.166.193.132

NSE: Script Post-scanning.
Initiating NSE at 10:05
Completed NSE at 10:05, 0.00s elapsed
Initiating NSE at 10:05
Completed NSE at 10:05, 0.00s elapsed
Initiating NSE at 10:05
Completed NSE at 10:05, 0.00s elapsed
Read data files from: /usr/bin/../share/nmap
OS and Service detection performed. Please report any incorrect results at https://nmap.org/submit/ .
Nmap done: 2 IP addresses (2 hosts up) scanned in 30.28 seconds
           Raw packets sent: 2230 (104.504KB) | Rcvd: 2115 (86.104KB)
```

**Result of nmap'ing a worker node (swarm-worker-01)**

Command I used:

`anderskaas$ docker run --rm -it instrumentisto/nmap -A -T4 scanme.nmap.org -v -A -sV 165.227.130.157`

Result:

```
Starting Nmap 7.91 ( https://nmap.org ) at 2021-04-19 10:07 UTC
NSE: Loaded 153 scripts for scanning.
NSE: Script Pre-scanning.
Initiating NSE at 10:07
Completed NSE at 10:07, 0.00s elapsed
Initiating NSE at 10:07
Completed NSE at 10:07, 0.00s elapsed
Initiating NSE at 10:07
Completed NSE at 10:07, 0.00s elapsed
Initiating Ping Scan at 10:07
Scanning 2 hosts [4 ports/host]
Completed Ping Scan at 10:07, 0.01s elapsed (2 total hosts)
Initiating Parallel DNS resolution of 2 hosts. at 10:07
Completed Parallel DNS resolution of 2 hosts. at 10:07, 0.36s elapsed
Initiating SYN Stealth Scan at 10:07
Scanning 2 hosts [1000 ports/host]
Discovered open port 8080/tcp on 165.227.130.157
Discovered open port 22/tcp on 165.227.130.157
Discovered open port 9200/tcp on 165.227.130.157
Discovered open port 8081/tcp on 165.227.130.157
Completed SYN Stealth Scan against 165.227.130.157 in 0.50s (1 host left)
Discovered open port 80/tcp on 45.33.32.156
Discovered open port 22/tcp on 45.33.32.156
Discovered open port 9929/tcp on 45.33.32.156
Discovered open port 31337/tcp on 45.33.32.156
Completed SYN Stealth Scan at 10:07, 1.71s elapsed (2000 total ports)
Initiating Service scan at 10:07
Scanning 8 services on 2 hosts
Completed Service scan at 10:07, 11.22s elapsed (8 services on 2 hosts)
Initiating OS detection (try #1) against 2 hosts
Retrying OS detection (try #2) against 2 hosts
adjust_timeouts2: packet supposedly had rtt of -611325 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -611325 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -465729 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -465729 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -569155 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -569155 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -433639 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -433639 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -533691 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -533691 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -507445 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -507445 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -464391 microseconds.  Ignoring time.
adjust_timeouts2: packet supposedly had rtt of -464391 microseconds.  Ignoring time.
Initiating Traceroute at 10:07
Completed Traceroute at 10:07, 0.03s elapsed
Initiating Parallel DNS resolution of 2 hosts. at 10:07
Completed Parallel DNS resolution of 2 hosts. at 10:07, 0.01s elapsed
NSE: Script scanning 2 hosts.
Initiating NSE at 10:07
Completed NSE at 10:08, 5.94s elapsed
Initiating NSE at 10:08
Completed NSE at 10:08, 2.69s elapsed
Initiating NSE at 10:08
Completed NSE at 10:08, 0.00s elapsed
Nmap scan report for scanme.nmap.org (45.33.32.156)
Host is up (0.037s latency).
Other addresses for scanme.nmap.org (not scanned): 2600:3c01::f03c:91ff:fe18:bb2f
Not shown: 996 closed ports
PORT      STATE SERVICE    VERSION
22/tcp    open  ssh        OpenSSH 6.6.1p1 Ubuntu 2ubuntu2.13 (Ubuntu Linux; protocol 2.0)
| ssh-hostkey: 
|   1024 ac:00:a0:1a:82:ff:cc:55:99:dc:67:2b:34:97:6b:75 (DSA)
|   2048 20:3d:2d:44:62:2a:b0:5a:9d:b5:b3:05:14:c2:a6:b2 (RSA)
|   256 96:02:bb:5e:57:54:1c:4e:45:2f:56:4c:4a:24:b2:57 (ECDSA)
|_  256 33:fa:91:0f:e0:e1:7b:1f:6d:05:a2:b0:f1:54:41:56 (ED25519)
80/tcp    open  http       Apache httpd 2.4.7 ((Ubuntu))
|_http-favicon: Nmap Project
| http-methods: 
|_  Supported Methods: OPTIONS GET HEAD POST
|_http-server-header: Apache/2.4.7 (Ubuntu)
|_http-title: Go ahead and ScanMe!
9929/tcp  open  nping-echo Nping echo
31337/tcp open  tcpwrapped
OS fingerprint not ideal because: Host distance (11 network hops) is greater than five
No OS matches for host
Network Distance: 2 hops
TCP Sequence Prediction: Difficulty=135 (Good luck!)
IP ID Sequence Generation: Randomized
Service Info: OS: Linux; CPE: cpe:/o:linux:linux_kernel

TRACEROUTE (using port 80/tcp)
HOP RTT     ADDRESS
1   0.11 ms 172.17.0.1
2   0.87 ms scanme.nmap.org (45.33.32.156)

Nmap scan report for 165.227.130.157
Host is up (0.0087s latency).
Not shown: 996 closed ports
PORT     STATE SERVICE VERSION
22/tcp   open  ssh     OpenSSH 8.2p1 Ubuntu 4ubuntu0.2 (Ubuntu Linux; protocol 2.0)
| ssh-hostkey: 
|   3072 b6:d8:47:4e:41:a3:85:9a:c0:46:66:91:33:50:2e:15 (RSA)
|   256 bb:12:7a:8b:24:7b:8b:3e:ff:91:16:3a:d2:d0:ea:b9 (ECDSA)
|_  256 dc:d8:f1:fa:74:4e:8c:22:51:f1:f8:28:0e:00:62:20 (ED25519)
8080/tcp open  http    Golang net/http server (Go-IPFS json-rpc or InfluxDB API)
| http-methods: 
|_  Supported Methods: GET
|_http-open-proxy: Proxy might be redirecting requests
|_http-title:  Public Timeline | MiniTwit
8081/tcp open  http    Golang net/http server (Go-IPFS json-rpc or InfluxDB API)
|_http-title: Site doesn't have a title (text/plain; charset=utf-8).
9200/tcp open  http    Elasticsearch REST API 7.12.0 (name: elasticsearch; cluster: docker-cluster; Lucene 8.8.0)
|_http-favicon: Unknown favicon MD5: 6177BFB75B498E0BB356223ED76FFE43
| http-methods: 
|   Supported Methods: HEAD GET DELETE OPTIONS
|_  Potentially risky methods: DELETE
|_http-title: Site doesn't have a title (application/json; charset=UTF-8).
OS fingerprint not ideal because: Host distance (12 network hops) is greater than five
No OS matches for host
Network Distance: 2 hops
TCP Sequence Prediction: Difficulty=136 (Good luck!)
IP ID Sequence Generation: Randomized
Service Info: OS: Linux; CPE: cpe:/o:linux:linux_kernel

TRACEROUTE (using port 80/tcp)
HOP RTT     ADDRESS
-   Hop 1 is the same as for 45.33.32.156
2   0.32 ms 165.227.130.157

NSE: Script Post-scanning.
Initiating NSE at 10:08
Completed NSE at 10:08, 0.00s elapsed
Initiating NSE at 10:08
Completed NSE at 10:08, 0.00s elapsed
Initiating NSE at 10:08
Completed NSE at 10:08, 0.00s elapsed
Read data files from: /usr/bin/../share/nmap
OS and Service detection performed. Please report any incorrect results at https://nmap.org/submit/ .
Nmap done: 2 IP addresses (2 hosts up) scanned in 30.23 seconds
           Raw packets sent: 2246 (103.348KB) | Rcvd: 2115 (86.696KB)

```



