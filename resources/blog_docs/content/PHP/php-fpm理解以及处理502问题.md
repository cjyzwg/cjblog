---
title: php-fpm理解以及处理502问题
categories:
  - PHP
---
### 网站出现502
- 杀死进程
``` shell
    #!/bin/sh
    pids=$(ps -e -o 'pid,comm,args,pcpu,rsz,vsz,stime,user,uid' | sort -nrk5 | grep www | grep -v grep|awk ' $2!="nginx" {print $1}')
    #echo $pids
    if [ "$pids" != "" ];then
    for  pid  in   $pids;
    do
    echo $pid
    kill -9 $pid
    done
    fi
```
- 调高进程数  
使用 netstat -napo |grep "php-fpm" | wc -l 查看一下当前 fastcgi 进程个数，如果个数接近 conf 里配置的上限，就需要调高进程数。
但也不能无休止调高，可以根据服务器内存情况，可以把 php-fpm 子进程数调到 100 或以上，在 4G 内存的服务器上 200 就可以。
- 调高 linux 内核打开文件数量  
可以使用这些命令 ( 必须是 root 帐号 )  
echo 'ulimit -HSn 65536'>> /etc/profile  
echo 'ulimit -HSn 65536'>> /etc/rc.local  
source /etc/profile  
- 调整脚本执行时间  
如果脚本因为某种原因长时间等待不返回 ，导致新来的请求不能得到处理，可以适当调小如下配置。
nginx.conf 里面主要是如下：  
fastcgi_connect_timeout 300;  
fastcgi_send_timeout 300;  
fastcgi_read_timeout 300;  
php-fpm.conf 里如要是如下  
request_terminate_timeout =10s  
- 增加缓存  
修改或增加配置到 nginx.conf  
proxy_buffer_size 64k;  
proxy_buffers  512k;  
proxy_busy_buffers_size 128k;  

### php-fpm理解
webserver每收到一个请求，都会去fork一个cgi进程，请求结束再kill掉这个进程。这样有10000个请求，就需要fork、kill php-cgi进程10000次,很浪费资源？于是，出现了cgi的改良版本，fast-cgi。fast-cgi每次处理完请求后，不会kill掉这个进程，而是保留这个进程，使这个进程可以一次处理多个请求。这样每次就不用重新fork一个进程了，大大提高了效率。    

php-fpm是什么
php-fpm即php-Fastcgi Process Manager. php-fpm是 FastCGI 的实现，并提供了进程管理的功能。 进程包含 master 进程和 worker 进程两种进程。 master 进程只有一个，负责监听端口，接收来自 Web Server 的请求，而 worker 进程则一般有多个(具体数量根据实际需要配置)，每个进程内部都嵌入了一个 PHP 解释器，是 PHP 代码真正执行的地方。

#### php-fpm.conf
``` conf
* pid = /usr/local/var/run/php-fpm.pid
* #pid设置，一定要开启,上面是Mac平台的。默认在php安装目录中的var/run/php-fpm.pid。比如centos的在: /usr/local/php/var/run/php-fpm.pid
* 
* error_log = /usr/local/var/log/php-fpm.log
* #错误日志，上面是Mac平台的，默认在php安装目录中的var/log/php-fpm.log，比如centos的在: /usr/local/php/var/log/php-fpm.log
* 
* log_level = notice
* #错误级别. 上面的php-fpm.log纪录的登记。可用级别为: alert（必须立即处理）, error（错误情况）, warning（警告情况）, notice（一般重要信息）, debug（调试信息）. 默认: notice.
* 
* emergency_restart_threshold = 60
* emergency_restart_interval = 60s
* #表示在emergency_restart_interval所设值内出现SIGSEGV或者SIGBUS错误的php-cgi进程数如果超过 emergency_restart_threshold个，php-fpm就会优雅重启。这两个选项一般保持默认值。0 表示 '关闭该功能'. 默认值: 0 (关闭).
* 
* process_control_timeout = 0
* #设置子进程接受主进程复用信号的超时时间. 可用单位: s(秒), m(分), h(小时), 或者 d(天) 默认单位: s(秒). 默认值: 0.
* 
* daemonize = yes
* #后台执行fpm,默认值为yes，如果为了调试可以改为no。在FPM中，可以使用不同的设置来运行多个进程池。 这些设置可以针对每个进程池单独设置。
* 
* listen = 127.0.0.1:9000
* #fpm监听端口，即nginx中php处理的地址，一般默认值即可。可用格式为: 'ip:port', 'port', '/path/to/unix/socket'. 每个进程池都需要设置。如果nginx和php在不同的机器上，分布式处理，就设置ip这里就可以了。
* 
* listen.backlog = -1
* #backlog数，设置 listen 的半连接队列长度，-1表示无限制，由操作系统决定，此行注释掉就行。backlog含义参考：http://www.3gyou.cc/?p=41
* 
* listen.allowed_clients = 127.0.0.1
* #允许访问FastCGI进程的IP白名单，设置any为不限制IP，如果要设置其他主机的nginx也能访问这台FPM进程，listen处要设置成本地可被访问的IP。默认值是any。每个地址是用逗号分隔. 如果没有设置或者为空，则允许任何服务器请求连接。
* 
* listen.owner = www
* listen.group = www
* listen.mode = 0666
* #unix socket设置选项，如果使用tcp方式访问，这里注释即可。
* 
* user = www
* group = www
* #启动进程的用户和用户组，FPM 进程运行的Unix用户, 必须要设置。用户组，如果没有设置，则默认用户的组被使用。
* 
* pm = dynamic
* #php-fpm进程启动模式，pm可以设置为static和dynamic和ondemand
* #如果选择static，则进程数就数固定的，由pm.max_children指定固定的子进程数。
* 
* #如果选择dynamic，则进程数是动态变化的,由以下参数决定：
* pm.max_children = 50 #子进程最大数
* pm.start_servers = 2 #启动时的进程数，默认值为: min_spare_servers + (max_spare_servers - min_spare_servers) / 2
* pm.min_spare_servers = 1 #保证空闲进程数最小值，如果空闲进程小于此值，则创建新的子进程
* pm.max_spare_servers = 3 #，保证空闲进程数最大值，如果空闲进程大于此值，此进行清理
* 
* pm.max_requests = 500
* #设置每个子进程重生之前服务的请求数. 对于可能存在内存泄漏的第三方模块来说是非常有用的. 如果设置为 '0' 则一直接受请求. 等同于 PHP_FCGI_MAX_REQUESTS 环境变量. 默认值: 0.
* 
* pm.status_path = /status
* #FPM状态页面的网址. 如果没有设置, 则无法访问状态页面. 默认值: none. munin监控会使用到
* 
* ping.path = /ping
* #FPM监控页面的ping网址. 如果没有设置, 则无法访问ping页面. 该页面用于外部检测FPM是否存活并且可以响应请求. 请注意必须以斜线开头 (/)。
* 
* ping.response = pong
* #用于定义ping请求的返回相应. 返回为 HTTP 200 的 text/plain 格式文本. 默认值: pong.
* 
* access.log = log/$pool.access.log
* #每一个请求的访问日志，默认是关闭的。
* 
* access.format = "%R - %u %t \"%m %r%Q%q\" %s %f %{mili}d %{kilo}M %C%%"
* #设定访问日志的格式。
* 
* slowlog = log/$pool.log.slow
* #慢请求的记录日志,配合request_slowlog_timeout使用，默认关闭
* 
* request_slowlog_timeout = 10s
* #当一个请求该设置的超时时间后，就会将对应的PHP调用堆栈信息完整写入到慢日志中. 设置为 '0' 表示 'Off'
* 
* request_terminate_timeout = 0
* #设置单个请求的超时中止时间. 该选项可能会对php.ini设置中的'max_execution_time'因为某些特殊原因没有中止运行的脚本有用. 设置为 '0' 表示 'Off'.当经常出现502错误时可以尝试更改此选项。
* 
* rlimit_files = 1024
* #设置文件打开描述符的rlimit限制. 默认值: 系统定义值默认可打开句柄是1024，可使用 ulimit -n查看，ulimit -n 2048修改。
* 
* rlimit_core = 0
* #设置核心rlimit最大限制值. 可用值: 'unlimited' 、0或者正整数. 默认值: 系统定义值.
* 
* chroot =
* #启动时的Chroot目录. 所定义的目录需要是绝对路径. 如果没有设置, 则chroot不被使用.
* 
* chdir =
* #设置启动目录，启动时会自动Chdir到该目录. 所定义的目录需要是绝对路径. 默认值: 当前目录，或者/目录（chroot时）
* 
* catch_workers_output = yes
* #重定向运行过程中的stdout和stderr到主要的错误日志文件中. 如果没有设置, stdout 和 stderr 将会根据FastCGI的规则被重定向到 /dev/null . 默认值: 空.
```

#### php-fpm进程分配
在之前的文章中就说过了。在fasgcgi模式下，php会启动多个php-fpm进程，来接收nginx发来的请求，那是不是进程越多，速度就越快呢？这可不一定！得根据我们的机器配置和业务量来决定。  
我们先来看下，设定进程的配置在哪里？  
pm = static | dynamic | ondemand  
pm可以设置成这样3种，我们用的最多的就上前面2种。  
pm = static 模式  
pm = static 表示我们创建的php-fpm子进程数量是固定的，那么就只有pm.max_children = 50这个参数生效。你启动php-fpm的时候就会一起全部启动51(1个主＋50个子)个进程，颇为壮观。  
pm = dynamic 模式  
pm = dynamic模式，表示启动进程是动态分配的，随着请求量动态变化的。他由 pm.max_children，pm.start_servers，pm.min_spare_servers，pm.max_spare_servers 这几个参数共同决定。  
上面已经讲过，这里再重申一下吧：  
pm.max_children ＝ 50 是最大可创建的子进程的数量。必须设置。这里表示最多只能50个子进程。  
pm.start_servers = 20 随着php-fpm一起启动时创建的子进程数目。默认值：min_spare_servers + (max_spare_servers - min_spare_servers) / 2。这里表示，一起启动会有20个子进程。
pm.min_spare_servers = 10 设置服务器空闲时最小php-fpm进程数量。必须设置。如果空闲的时候，会检查如果少于10个，就会启动几个来补上。  
pm.max_spare_servers = 30 设置服务器空闲时最大php-fpm进程数量。必须设置。如果空闲时，会检查进程数，多于30个了，就会关闭几个，达到30个的状态。  
到底选择static还数dynamic?  
很多人恐惧症来袭，不知道选什么好？  
一般原则是：动态适合小内存机器，灵活分配进程，省内存。静态适用于大内存机器，动态创建回收进程对服务器资源也是一种消耗。  
如果你的内存很大，有8-20G，按照一个php-fpm进程20M算，100个就2G内存了，那就可以开启static模式。如果你的内存很小，比如才256M，那就要小心设置了，因为你的机器里面的其他的进程也算需要占用内存的，所以设置成dynamic是最好的，比如：pm.max_chindren = 8, 占用内存160M左右，而且可以随时变化，对于一半访问量的网站足够了。  
慢日志查询  
我们有时候会经常饱受500,502问题困扰。当nginx收到如上错误码时，可以确定后端php-fpm解析php出了某种问题，比如，执行错误，执行超时。  
这个时候，我们是可以开启慢日志功能的。  

* slowlog = /usr/local/var/log/php-fpm.log.slow  
* request_slowlog_timeout = 15s  
当一个请求该设置的超时时间15秒后，就会将对应的PHP调用堆栈信息完整写入到慢日志中。  
php-fpm慢日志会记录下进程号，脚本名称，具体哪个文件哪行代码的哪个函数执行时间过长：  

* [21-Nov-2013 14:30:38] [pool www] pid 11877  
* script_filename = /usr/local/lnmp/nginx/html/www.quancha.cn/www/fyzb.php  
* [0xb70fb88c] file_get_contents() /usr/local/lnmp/nginx/html/www.quancha.cn/www/fyzb.php:2  
通过日志，我们就可以知道第2行的file_get_contents 函数有点问题，这样我们就能追踪问题了。  

转载：https://blog.csdn.net/ahaotata/article/details/83825977

