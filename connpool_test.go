package connpool

import(
	"net"
	"sync"
	"time"
	"testing"
	"fmt"
	"os"
	"log"
	"bytes"
)

func factory() (net.Conn,error) {
	return net.Dial("tcp","127.0.0.1:7060")
}

func TestPool(t *testing.T) {
	pool1 := NewPool(10, 30 ,factory)
	defer pool1.Close()
	var err1 error
	var wg sync.WaitGroup
	for i:=0;i<50;i++ {
		wg.Add(1)
		go func(n int){
			conn ,err := pool1.Get()
			if err!=nil {
				err1 = err
			}
			//t.Log(n,err, time.Now().Format("15:04:05"))
			time.Sleep( 5*time.Second )
			msg := []byte( fmt.Sprintf("msg%d\n",n) )
			buf := make([]byte,len(msg))
			_,err = conn.Write( msg )
			if err!=nil{
				t.Log(err)
			}
			n1,err := conn.Read( buf )
			if err!=nil{
				t.Log(err)
			}else if bytes.Compare( msg[:n1] , buf[:n1] )!=0 {
				t.Log("Fail",msg[:n1],buf[:n1])
			}
			pool1.Put(conn)
			wg.Done()
		}(i)
	}
	wg.Wait()
	if err1!=nil {
		t.Fatal(err1)
	}
}

func TestMain(m *testing.M) {
	l, e := net.Listen("tcp", "127.0.0.1:7060")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer l.Close()
	go echoServer(l)
	os.Exit(m.Run())
}

func echoServer(l net.Listener) {
	for i:=0;i<50;i++ {
		conn,err := l.Accept()
		if err!=nil {
			break
		}
		go func(conn1 net.Conn) {
			var buf [100]byte
			for {
				n,err:=conn1.Read(buf[:])
				if err==nil{
					conn1.Write(buf[:n])
				}else{
					break
				}
			}
			conn1.Close()
		}(conn)
	}
}
