import "strings"

func plugin() string {
	var ret []string
	ret = append(ret, "Scanning subnet...")
	ret = append(ret, "# nmap -v 127.0.0.1")
	ret = append(ret, "Starting Nmap 7.93 ( https://nmap.org ) at 2022-09-31 24:01 CST")
	ret = append(ret, "Initiating Ping Scan at 24:01")
	ret = append(ret, "Scanning 127.0.0.1 [2 ports]")
	ret = append(ret, "Completed Ping Scan at 24:01, 0.00s elapsed (1 total hosts)")
	ret = append(ret, "Initiating Connect Scan at 24:01")
	ret = append(ret, "Scanning localhost (127.0.0.1) [1000 ports]")
	ret = append(ret, "Discovered open port 22/tcp on 127.0.0.1")
	ret = append(ret, "Discovered open port 80/tcp on 127.0.0.1")
	ret = append(ret, "Discovered open port 3306/tcp on 127.0.0.1")
	ret = append(ret, "Completed Connect Scan at 24:01, 0.02s elapsed (1000 total ports)")
	ret = append(ret, "Nmap scan report for localhost (127.0.0.1)")
	ret = append(ret, "Host is up (0.000052s latency).")
	ret = append(ret, "Not shown: 997 closed tcp ports (conn-refused)")
	ret = append(ret, "PORT     STATE SERVICE")
	ret = append(ret, "22/tcp   open  ssh")
	ret = append(ret, "80/tcp   open  http")
	ret = append(ret, "3306/tcp open  mysql")
	ret = append(ret, "")
	ret = append(ret, "Read data files from: /opt/homebrew/bin/../share/nmap")
	ret = append(ret, "Nmap done: 1 IP address (1 host up) scanned in 0.06 seconds")
	return strings.Join(ret, "\n")
}
