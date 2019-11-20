# grpc php client 
**服务端取的是https://www.jianshu.com/p/7fe7a8507745 案例**  


#### 需要执行
> export GO111MODULE=on  
go mod tidy
#### 安装依赖
>pecl install grpc  
pecl install protobuf  
同时并.so 文件放进ini中，如果是在docker中使用docker-php-ext-enable grpc,docker-php-ext-enable protobuf  
wget https://github.com/google/protobuf/releases/download/v3.5.1/protobuf-all-3.5.1.zip  
unzip protobuf-all-3.5.1.zip  
cd protobuf-3.5.1/  
./configure  
make  
make install  
protoc --version  
出现make 错误，安装apt-get install automake  
ldconfig  

1、protoc --proto_path=../../ --php_out=. helloServer.proto  
2、在client/phpclient/Rpc_proto 文件夹新建HelloClient.php  

    <?php
		namespace Rpc_proto;
		class HelloClient extends \Grpc\BaseStub{
			public function __construct($hostname, $opts, $channel = null) {
				parent::__construct($hostname, $opts, $channel);
			}
			public function SayHello(\Rpc_proto\HelloRequest $argument,$metadata=[],$options=[]){
				return $this->_simpleRequest('/rpc_proto.HelloServer/SayHello',
					$argument,
					['\Rpc_proto\HelloReply', 'decode'],
					$metadata, $options);
			}
			public function GetHelloMsg(\Rpc_proto\HelloRequest $argument,$metadata=[],$options=[]){
				return $this->_simpleRequest('/rpc_proto.HelloServer/GetHelloMsg',
					$argument,
					['\Rpc_proto\HelloMessage', 'decode'],
					$metadata, $options);
			}
		}
	?>
3、cd client/phpclient && composer install  
4、cd server && go run server.go  
5、新建hello.php,并执行php hello.php
结果：
> helloGreenHat  
this is from server HAHA!