package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"strconv"
)

type mapstrint map[string]int
type strIntPair struct{
	val string // the value string like IP, request or timestamp
	cnt int    // count of occurences we need to order by
}
const DEFAULTCOUNT = 10 // default number of most frequent strings/values to print (configurable via 2nd cmd agrument)
var COUNT = DEFAULTCOUNT

// method to find and print from a"mapstrint" map the most frequent occurences of given key/string; "COUNT" souch values to be printed
// there are ways to implement it in more elegant - but slower - way, like converting the map to an array and sorting the *whole* array
// but since we need just top few items, this is expected to be much faster - traverse the map just once and populate "max" field sequentially
func (m mapstrint) biggest_in_map(header string) {
	fmt.Println(COUNT, "most frequent", header, ":")
	max := make([]strIntPair, COUNT+1) // the "+1" to simplify a loop of shifting items
	for val, cnt := range m {
		pos := COUNT-1
		for {
			if cnt > max[pos].cnt {
				max[pos+1] = max[pos]
			}
			if (pos == 0) || (cnt <= max[pos].cnt) {
				break
			}
			pos--
		}
		if cnt > max[pos].cnt {
			max[pos].val = val
			max[pos].cnt = cnt
		} else if cnt > max[pos+1].cnt {
                        max[pos+1].val = val
                        max[pos+1].cnt = cnt
                }

	}
//	fmt.Println(m)
	for pos := COUNT-1; pos >=0; pos-- {
		if max[pos].cnt > 0 {
			fmt.Println(max[pos].cnt, "\t", max[pos].val)
		}
	}
	fmt.Println()
}

func main() {
	if len(os.Args)<2 {
		fmt.Println("Missing argument of httpd access log file.")
		os.Exit(1)
	}
	if len(os.Args)>2 {
		C, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("2nd argument must be int:", err)
			os.Exit(2)
		}
		COUNT = C
	}
	filePath := os.Args[1]
	readFile, err := os.Open(filePath)
	if err != nil {
        	fmt.Println(err)
		os.Exit(3)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	// example of parsed line:
	// 1.2.3.4 - - [16/Jun/2022:14:14:47 +0200] "GET /rhsm/status HTTP/1.1" 200 426 "-" "RHSM/1.0 (cmd=rhsmcertd-worker) subscription-manager/1.28.13-6.el8_4"
	IPs := make(mapstrint)
	TSHours := make(mapstrint)
	TSs := make(mapstrint)
	requests := make(mapstrint)
	unifiedRequests := make(mapstrint)
	retCodes := make(mapstrint)
	commands := make(mapstrint)
	// main regular expression parsing whole line
	re := regexp.MustCompile(`(?P<ip>.*) .* .* \[(?P<timestamp>.*):[0-9][0-9].*\] \"(?P<request>.*?)\" (?P<retcode>\d*) .* .* \"(?P<prog>.*?)\"`)
	// regular expression to trim leading "RHSM/1.0 " from a command and cut anything after " "
	rhsmcmdre := regexp.MustCompile(`^(RHSM/\d.\d )?(?P<cmd>\S+)`)
	// regular expression to replace RHSM UUIDs by "UUID"
	uuidre := regexp.MustCompile(`/(?P<uuid>[a-z0-9-]{36})`)
	// regular expression to trim trailing "HTTP/1.1" from a request
	httpre := regexp.MustCompile(` HTTP/\d\.\d$`)
	for fileScanner.Scan() {
		matches := re.FindStringSubmatch(fileScanner.Text())
		if len(matches) < 6 {
			continue
		}
		ip := matches[1]
		ts := matches[2]
		request := matches[3]
		retCode := matches[4]
		// from command, remove optional leading "RHSM/1.0 " and take the first string only
		cmdmatch := rhsmcmdre.FindStringSubmatch(matches[5])
		cmd := cmdmatch[2]
		tsHours := strings.Join(strings.Split(ts, ":")[:2], ":")
		unifiedRequest := httpre.ReplaceAllString(uuidre.ReplaceAllString(request, "/UUID"), "")
//		fmt.Println(ip, ts, tsHours, request, unifiedRequest, retCode, cmd)
		IPs[ip]++
		TSHours[tsHours]++
		TSs[ts]++
		requests[request]++
		unifiedRequests[unifiedRequest]++
		retCodes[retCode]++
		commands[cmd]++
	}
	readFile.Close()

	IPs.biggest_in_map("IPs")
	TSHours.biggest_in_map("timestamps (hours)")
	TSs.biggest_in_map("timestamps (minutes)")
	requests.biggest_in_map("requests")
	unifiedRequests.biggest_in_map("requests (unified)")
	retCodes.biggest_in_map("return codes")
	commands.biggest_in_map("commands")
}
