package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/textproto"
	"os"
	"strconv"
	"strings"
)

type maild struct {
	from    string
	to      string
	subject string
	body    []byte
}

type client struct {
	serverName string
	localName  string
	conn       net.Conn
	txt        *textproto.Conn
	maild      maild
}

var (
	Info  *log.Logger
	Error *log.Logger
)

func (c *client) Command(code int, cmd string) error {
	id, err := c.txt.Cmd(cmd)
	if err != nil {
		return err
	}

	c.txt.StartResponse(id)
	defer c.txt.EndResponse(id)

	if _, _, err := c.txt.ReadCodeLine(code); err != nil {
		return err
	}

	return nil
}

func Dial(addr string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
	}

	return conn, nil
}

func (c *client) Ehlo() error {

	if err := c.Command(220, "HELO "+c.localName); err != nil {
		return err
	}

	return nil
}

func (c *client) Mail() error {
	_, err := c.txt.Cmd("Mail from:xxx@gmail.com")
	if err != nil {
		return err
	}

	if _, err := c.txt.ReadLineBytes(); err != nil {
		return err
	}

	return nil
}

func (c *client) Rcpt() error {
	_, err := c.txt.Cmd("Rcpt to:jheyu@xxx")
	if err != nil {
		return err
	}

	if _, err := c.txt.ReadLineBytes(); err != nil {
		return err
	}

	return nil
}

func (c *client) Data() error {
	_, err := c.txt.Cmd("DATA")
	if err != nil {
		return err
	}

	if _, err = c.txt.ReadLineBytes(); err != nil {
		return err
	}

	//_, err = c.txt.Cmd("From:" + c.maild.from)
	//if err != nil {
	//      Error.Println("Data From error")
	//      return err
	//}

	//_, err = c.txt.Cmd("To:" + c.maild.to)
	//if err != nil {
	//      Error.Println("Data To error")
	//      return err
	//}

	//_, err = c.txt.Cmd("Subject:" + c.maild.subject)
	//if err != nil {
	//      Error.Println("Data subject error")
	//      return err
	//}

	//_, err = c.txt.Cmd("\n")
	//if err != nil {
	//      fmt.Println("Data \n error")
	//      return err
	//}

	_, err = c.txt.Cmd(string(c.maild.body))
	if err != nil {
		return err
	}

	_, err = c.txt.Cmd(".")
	if err != nil {
		return err
	}

	if _, err := c.txt.ReadLineBytes(); err != nil {
		return err
	}

	return nil
}

func (c *client) Rset() error {
	if _, err := c.txt.Cmd("RSET"); err != nil {
		return err
	}

	//str, _ := c.txt.ReadLineBytes()
	//Info.Println(string(str))

	return nil
}

func (c *client) Quit() error {
	if _, err := c.txt.Cmd("Quit"); err != nil {
		return err
	}

	//str, _ := c.txt.ReadLineBytes()

	return nil
}

func SendMail(addr string, mailt maild) error {
	conn, err := Dial(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	text := textproto.NewConn(conn)
	c := &client{txt: text, conn: conn, serverName: addr, localName: "jheyu.localhost", maild: mailt}

	if err = c.Ehlo(); err != nil {
		return err
	}

	if err = c.Mail(); err != nil {
		return err
	}

	if err = c.Rcpt(); err != nil {
		return err
	}

	if err = c.Data(); err != nil {
		return err
	}

	if err = c.Quit(); err != nil {
		return err
	}

	return nil
}

func Stat(user string, pass string) error {
	conn, err := Dial("xxx.xxx.xxx.xxx:110")
	if err != nil {
		return err
	}
	defer conn.Close()

	text := textproto.NewConn(conn)

	res, _ := text.ReadLineBytes()
	fmt.Println("pop3 connect sucessed response: ", string(res))

	_, err = text.Cmd("user jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("user command response: ", string(res))

	_, err = text.Cmd("pass jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("pass command response: ", string(res))

	_, err = text.Cmd("STAT")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("STAT command response: ", string(res))

	_, err = text.Cmd("QUIT")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("QUIT command response: ", string(res))

	return nil
}

func Retr(user string, pass string, index string, mail_file string) error {
	conn, err := Dial("xxx.xxx.xxx.xxx:110")
	if err != nil {
		return err
	}
	defer conn.Close()

	text := textproto.NewConn(conn)
	res, _ := text.ReadLineBytes()
	fmt.Println("pop3 connect sucessed response: ", string(res))

	_, err = text.Cmd("user jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("user command response: ", string(res))

	_, err = text.Cmd("pass jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("pass command response: ", string(res))

	cmd := "RETR " + index
	_, err = text.Cmd(cmd)
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("RETR command response: ", string(res))

	var mail_contents string
	res, _ = text.ReadLineBytes()
	for string(res) != "." {
		mail_contents = mail_contents + string(res) + "\n"
		res, _ = text.ReadLineBytes()
	}

	file_contents := []byte(mail_contents)
	err = ioutil.WriteFile(mail_file, file_contents, 0644)

	_, err = text.Cmd("QUIT")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("QUIT command response: ", string(res))

	return nil
}

func Del(user string, pass string, index string) error {
	conn, err := Dial("xxx.xxx.xxx.xxx:110")
	if err != nil {
		return err
	}
	defer conn.Close()

	text := textproto.NewConn(conn)

	res, _ := text.ReadLineBytes()
	fmt.Println("pop3 connect sucessed response: ", string(res))

	_, err = text.Cmd("user jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("user command response: ", string(res))

	_, err = text.Cmd("pass jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("pass command response: ", string(res))

	cmd := "DELE " + index
	_, err = text.Cmd(cmd)
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("DELE", index, "command response: ", string(res))

	_, err = text.Cmd("QUIT")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("QUIT command response: ", string(res))

	return nil
}

func GetAll(user string, pass string, mail_dir string) error {
	conn, err := Dial("xxx.xxx.xxx.xxx:110")
	if err != nil {
		return err
	}
	defer conn.Close()

	text := textproto.NewConn(conn)
	res, _ := text.ReadLineBytes()
	fmt.Println("pop3 connect sucessed response: ", string(res))

	_, err = text.Cmd("user jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("user command response: ", string(res))

	_, err = text.Cmd("pass jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("pass command response: ", string(res))

	// PARSE START
	_, err = text.Cmd("STAT")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("STAT command response: ", string(res))
	fmt.Printf("Fields are: %q\n", strings.Fields(string(res)))

	var fields []string
	fields = strings.Fields(string(res))

	var count int
	count, _ = strconv.Atoi(fields[1])
	fmt.Println("count: ", count)

	for i := 1; i <= count; i++ {
		cmd := "RETR " + strconv.Itoa(i)
		_, err = text.Cmd(cmd)
		if err != nil {
			return err
		}
		res, _ = text.ReadLineBytes()
		fmt.Println("RETR command response: ", string(res))

		var mail_contents string
		res, _ = text.ReadLineBytes()
		for string(res) != "." {
			mail_contents = mail_contents + string(res) + "\n"
			res, _ = text.ReadLineBytes()
		}

		file_contents := []byte(mail_contents)
		err = ioutil.WriteFile(mail_dir+"/"+strconv.Itoa(i), file_contents, 0644)
	}

	_, err = text.Cmd("QUIT")
	if err != nil {
		return err
	}

	res, _ = text.ReadLineBytes()
	fmt.Println("QUIT command response: ", string(res))

	return nil
}

func DelAll(user string, pass string) error {
	conn, err := Dial("xxx.xxx.xxx.xxx:110")
	if err != nil {
		return err
	}
	defer conn.Close()

	text := textproto.NewConn(conn)

	res, _ := text.ReadLineBytes()
	fmt.Println("pop3 connect sucessed response: ", string(res))

	_, err = text.Cmd("user jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("user command response: ", string(res))

	_, err = text.Cmd("pass jheyu")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("pass command response: ", string(res))

	_, err = text.Cmd("STAT")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("STAT command response: ", string(res))
	fmt.Printf("Fields are: %q\n", strings.Fields(string(res)))

	var fields []string
	fields = strings.Fields(string(res))

	var count int
	count, _ = strconv.Atoi(fields[1])
	fmt.Println("count: ", count)

	for i := 1; i <= count; i++ {
		cmd := "DELE " + strconv.Itoa(i)
		_, err = text.Cmd(cmd)
		if err != nil {
			return err
		}
		res, _ = text.ReadLineBytes()
		fmt.Println("DELE", i, "command response: ", string(res))

	}

	_, err = text.Cmd("QUIT")
	if err != nil {
		return err
	}
	res, _ = text.ReadLineBytes()
	fmt.Println("QUIT command response: ", string(res))

	return nil
}

func main() {

	// Initial log file
	info, _ := os.OpenFile("stress_info_log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0700)
	Info = log.New(info, "[INFO] ", log.Ltime)
	log.SetOutput(info)
	defer info.Close()

	//_, _, user, pass, commandmail_dir := os.Args[1], os.Args[2], os.Args[3], os.Args[4]
	flag.Parse()
	fmt.Printf("args=%s, num=%d\n", flag.Args(), flag.NArg())
	for i := 0; i != flag.NArg(); i++ {
		fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))
	}

	switch flag.Arg(4) {
	case "count":
		fmt.Println("<<<<< command count >>>>>")
		Stat(flag.Arg(2), flag.Arg(3))
	case "get":
		fmt.Println("<<<<< command get index mail_file >>>>>")
		Retr(flag.Arg(2), flag.Arg(3), flag.Arg(5), flag.Arg(6))
	case "delete":
		fmt.Println("<<<<< command delete index >>>>>")
		Del(flag.Arg(2), flag.Arg(3), flag.Arg(5))
	case "getall":
		fmt.Println("<<<<< command getall >>>>>")
		GetAll(flag.Arg(2), flag.Arg(3), flag.Arg(5))
	case "deleteall":
		fmt.Println("<<<<< command deleteall >>>>>")
		DelAll(flag.Arg(2), flag.Arg(3))
	}
}
