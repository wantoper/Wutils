package WHttp

import (
	"WUtils/WTls"
	"WUtils/WTls/Util"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Server 表示HTTP服务器
type Server struct {
	Addr   string
	Router *Router
}

// NewServer 创建一个新的HTTP服务器
func NewServer(addr string) *Server {
	return &Server{
		Addr:   addr,
		Router: NewRouter(),
	}
}

func (s *Server) ListenAndServe() error {
	return StartServer(s.Addr, s.Router)
}

func (s *Server) GET(path string, handler HandlerFunc) {
	s.Router.GET(path, handler)
}

func (s *Server) POST(path string, handler HandlerFunc) {
	s.Router.POST(path, handler)
}

func create_request(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	schema, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	mehtod, rest, _ := strings.Cut(schema, " ")
	requestURI, proto, _ := strings.Cut(rest, " ")
	rawURI := "http://" + requestURI
	parseRequestURI, _ := url.ParseRequestURI(rawURI)

	headers := make(Header)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}

		if strings.EqualFold(line, "\r\n") {
			break
		}

		headerK, headerV, _ := strings.Cut(line, ":")
		headers[headerK] = strings.TrimSpace(headerV)
	}

	if parseRequestURI.Host == "" {
		parseRequestURI.Host = headers.Get("Host")
	}
	contenlength := headers.Get("Content-Length")
	realLength, err := strconv.ParseInt(contenlength, 10, 64)
	req := &Request{
		Method: mehtod,
		Proto:  proto,
		Url:    parseRequestURI,
		Header: headers,
		Body:   io.LimitReader(reader, realLength),
	}

	return req, err
}

func handleConnection(conn net.Conn, router *Router) {
	request, _ := create_request(conn)
	fmt.Println("New connection from", request.Url)

	writer := ResponseWriterImpl{
		w: *bufio.NewWriter(conn),
	}

	if router != nil {
		router.ServeHTTP(&writer, request)
	}

	writer.FinishRequest()
	conn.Close()
}

func StartServer(address string, router *Router) error {
	//listen, err := net.Listen("tcp", address)
	publickey, _ := Util.GetPublicKey("server.crt")
	privateKey, _ := Util.GetPrivateKey("server.key")
	listen, err := WTls.NewTlsServer(address, publickey, privateKey)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}

	//defer listen.Close()
	fmt.Printf("Server started on %s\n", address)

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		fmt.Println("Accepting connection from", conn.RemoteAddr())
		go handleConnection(conn, router)
	}
}
