package main

import(
	"fmt"
	"flag"
	"os"
	"os/exec"
	"net"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

const config_file string = "config.cfg"
const log_file string = "system.log"


type Config struct {
  Server string
  Port uint
}


func readConfig() Config{
  file, err := ioutil.ReadFile(config_file)
  if err != nil {
	  log.Fatal(err)
  }
  var configs Config
  err = json.Unmarshal(file, &configs)
  if err != nil {
	  log.Fatal(err)
  }
  return configs
}


func start_server(){
//	configs := readConfig()
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":9090")
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil{
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}

}


func handleConnection(conn net.Conn){
	log.Println("Connection received")
	var command string
	var buf [5]byte
	for{
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Fatal(err)
		}
		command = command+string(buf[0:n])
		if strings.HasSuffix(command,"\n"){
			exe_command(strings.TrimRight(command,"\r"))
			command = ""
		}
	}
}

func exe_command(comm string) string{
	out, err := exec.Command("ls","").Output()
	if err != nil{ log.Fatal(err) }
	return string(out)
}



func start_client(){
	fmt.Println("Starting client")
}

func main(){
	log_file, err := os.OpenFile(log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(log_file)

	var server_mode bool
	flag.BoolVar(&server_mode, "server", false,"Start in server mode")
	flag.Parse()
	if server_mode {
		log.Println("Starting in server mode")
		start_server()
	}else{
		log.Println("Starting in client mode")
		start_client()
	}
}
