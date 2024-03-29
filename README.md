# connpool

#### 介绍
一个TCP连接池，用法参考connpool_test.go

#### 软件架构
```
type Conn
    func (s *Conn) Close() error
    func (s *Conn) Read(p []byte) (int, error)
    func (s *Conn) Timeout() bool
    func (s *Conn) Write(p []byte) (int, error)
type Pool
    func NewPool(max, timeout int, factory func() (net.Conn, error)) *Pool
    func (s *Pool) Close()
    func (s *Pool) Get() (*Conn, error)
    func (s *Pool) Put(conn1 *Conn)

type Conn
	Conn - Wrap of net.Conn

	type Conn struct {
	    // contains filtered or unexported fields
	}

func (s *Conn) Close() error
	Close - Close the connection


func (s *Conn) Read(p []byte) (int, error)
	Read - Compatible io.Reader

func (s *Conn) Timeout() bool
	Timeout - Test the connection is timeout. return true/false.

func (s *Conn) Write(p []byte) (int, error)
	Write - Compatible io.Write

type Pool
	Pool - Please create Pool with function NewPool

	type Pool struct {
	    // contains filtered or unexported fields
	}

func NewPool(max, timeout int, factory func() (net.Conn, error)) *Pool
	NewPool - Create a new Pool struct,and initialize it.

	max - the max connection number.
	timeout - the program will use it to set deadline value.
	factory - wrap of net.Dial(...).

func (s *Pool) Close()
	Close - Close all connections in this Pool

func (s *Pool) Get() (*Conn, error)
	Get - Get a new connection from the Pool

func (s *Pool) Put(conn1 *Conn)
	Put - Put a connection back to the Pool
```

#### 安装教程

go get -v -u github.com/rocket049/connpool

或：

go get -v -u gitee.com/rocket049/connpool

#### 使用说明

```
import "github.com/rocket049/connpool"
//import "gitee.com/rocket049/connpool"

func factory() (net.Conn,error) {
	return net.Dial("tcp","127.0.0.1:7060")
}

func UsePool() {
	pool1 := connpool.NewPool(10, 30 ,factory)
	defer pool1.Close()
	var wg sync.WaitGroup
	for i:=0;i<50;i++ {
		wg.Add(1)
		go func(n int){
		    // connect
			conn ,err := pool1.Get()
			if err!=nil {
				...
			}
			//send
			_,err = conn.Write( msg )
			if err!=nil{
				...
			}
			//recv
			n1,err := conn.Read( buf )
			if err!=nil{
				...
			}
			//timeout
			if conn.Timeout() {
				pool1.Put(conn)
				conn ,err := pool1.Get()
				...
			}
			//close
			pool1.Put(conn)
			wg.Done()
		}(i)
	}
	wg.Wait()

}
```

#### 参与贡献

1. Fork 本仓库
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request


#### 码云特技

1. 使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2. 码云官方博客 [blog.gitee.com](https://blog.gitee.com)
3. 你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解码云上的优秀开源项目
4. [GVP](https://gitee.com/gvp) 全称是码云最有价值开源项目，是码云综合评定出的优秀开源项目
5. 码云官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6. 码云封面人物是一档用来展示码云会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
