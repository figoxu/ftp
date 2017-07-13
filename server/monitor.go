package server

import "fmt"
import "time"
import "net/http"

func countdown(upsince int64) string {
	secs := time.Now().Unix() - upsince
	us := time.Unix(secs, 0)
	str := us.UTC().String()
	return str[11:19]
}

func trimGuid(guid string) string {
	return guid[0:6]
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "%d client(s), Up for %s\n",
		len(ConnectionMap), countdown(UpSince))

	for k, v := range ConnectionMap {
		fmt.Fprintf(w, "   %s %s, %s\n", trimGuid(k), countdown(v.connectedAt), v.user)
		pv := v.passive
		pk := v.passive.cid
		fmt.Fprintf(w, "     %s %s, %d %s %s\n", trimGuid(pk), countdown(pv.listenAt), pv.port, pv.command, pv.param)
	}
}

func Monitor() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5010", nil)
}

func Monitor2() {
	for {
		fmt.Println("9 clients, 4 passive, Up for 12:33:12")
		fmt.Println("   PFC106 00:17:34")
		fmt.Println("   PFC101 00:09:04")
		fmt.Println("   PFC109 00:07:04")
		fmt.Println("     PFC109p 03:33")
		fmt.Println("     PFC109p 01:33")
		fmt.Println("     PFC109p 00:00, 00:45")
		fmt.Println("   PFC116 00:06:04")
		fmt.Println("   PFC126 00:05:30")
		fmt.Println("   PFC104 00:04:30")
		fmt.Println("   PFC206 00:03:30")
		fmt.Println("     PFC206p 00:02:30")
		fmt.Println("   PFC306 00:02:30")
		fmt.Println("   PFC301 00:01:30")
		fmt.Println("Last 5 STORs:")
		fmt.Println("   PFC4FF /home/dir/path/filename.dat 22MB, 00:31:33")
		fmt.Println("   AFC4FF /home/dir/path/other.dat 292MB, 01:31:33")
		fmt.Println("Last 5 APPEs:")
		fmt.Println("   EFC4FF /home/dir/path/filename.dat 2MB, 05:31:33")
		fmt.Println("   AFC4FF /home/dir/path/other.dat 92MB, 02:31:33")
		time.Sleep(5 * time.Second)
	}
}
