package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

//Checar eventuais erros
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %s. \n", err.Error())
		os.Exit(1)
	}
}

func handleURL(a string) (string, string, []string) {
	label := strings.Split(a, "\r\n")
	url := label[0]
	label1 := strings.Split(url, " ")
	url1 := label1[1]
	label3 := strings.Split(url1, "/")
	url2 := label3[1]
	label4 := strings.Split(url2, " ")
	url3 := label4[0]
	url4 := label3[2:]
	temp := strings.Join(url4, "/")

	return url3, temp, label3
}

//Manipular os clientes e fecha a conexão
func handleClient(conn net.Conn) {
	var buf ([500000]byte)  //buffer
	var buf1 ([500000]byte) //buffer
	num, err := conn.Read(buf[0:])
	checkError(err)

	//Transforma o buf em string
	a := string(buf[0:num])

	url3, temp, label3 := handleURL(a)

	//Da o nome dos arquivos
	url5 := label3[1:]
	subs := strings.Join(url5, "%")
	leitor, err := ioutil.ReadFile(subs)
	text := string(leitor)
	data := time.Now()
	//Se o cache não tiver sido criado
	//Post it
	no_cache := ("\n<p style='" + "z-index:9999; position:fixed; top:20px; left:20px;width:200px;height:100px; background-color:yellow;padding:10px; font-weight:bold;'>" + "Cache: " + (data.Format("02/01/2006 15:04:05 ")) + "</p>")
	guia1 := strings.Replace(text, "<body>", "<body>"+no_cache, 1)
	conn.Write([]byte(guia1))

	if err != nil {
		if strings.Contains(a, "favicon.ico") {
			response := "HTTP/1.1 304\r\n\r\n"
			conn.Write([]byte(response))
			conn.Close()
			return
		}

		conn1, err := net.Dial("tcp", url3+":80")
		checkError(err)
		_, err = conn1.Write([]byte("GET /" + temp + " HTTP/1.1\r\nHost: " + url3 + "\r\n\r\n"))
		checkError(err)
		n, err := (conn1.Read(buf1[0:]))
		checkError(err)
		_, err = conn.Write(buf1[0:n])
		checkError(err)
		// Criando o arquivo em cache
		file := ioutil.WriteFile(subs, buf1[0:n], 0644)
		checkError(file)

	} else {
		no_cache := ("\n<p style='" + "z-index:9999; position:fixed; top:20px; left:20px;width:200px;height:100px; background-color:yellow;padding:10px; font-weight:bold;'>" + "Novo em: " + (data.Format("02/01/2006 15:04:05 ")) + "</p>")
		guia1 := strings.Replace(text, "<body>", "<body>"+no_cache, 1)
		conn.Write([]byte(guia1))
	}
}

func main() {
	fmt.Println("Servidor aguardando conexões...")

	listener, err := net.Listen("tcp", ":8888")
	checkError(err)

	//conn socket que está entre o programa e o navegador
	for {
		conn, err := listener.Accept()
		fmt.Println("Conexão aceita:")
		if err != nil {
			checkError(err)
			continue
		}

		go handleClient(conn)

	}

}
