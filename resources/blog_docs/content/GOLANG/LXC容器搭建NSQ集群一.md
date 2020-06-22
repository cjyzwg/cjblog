### LXC 容器搭建nsq集群
- 阿里云服务器cpu飙到100，登不上去，好不容易登上去之后打开lxc容器内发现连不上网
注意lxc的几个目录：/etc/default/lxc 
时间有限，暂时先贴这些
##### ecs_restart_lxc_config.sh
``` shell
	
	//检查容器内网络是否通畅
network()
{
    
    //超时时间
    local timeout=1
    lxcname=$(lxc list|awk ' NR>2 && $2!="" && $2!="|" {print $2}'|awk '{print $0}' | sed -n '$p')
    //目标网站
    local target=www.baidu.com

    //获取响应状态码
    local ret_code=`lxc exec $lxcname -- bash -c "curl -I -s --connect-timeout ${timeout} ${target} -w %{http_code} | tail -n1"`

    if [ "x$ret_code" = "x200" ]; then
        //网络畅通
        return 1
    else
        //网络不畅通
        return 0
    fi

    return 0
}

//awk 'BEGIN { FS=":";print "统计销售金额";total=0} {print NR;total=total+NR;} END {printf "销售金额总计：%.2f",total}' sx
lastnr=$(lxc list|awk ' NR>2 && $2!="" && $2!="|" {print $2}'|awk '{print NR}' | sed -n '$p')
if [ x"$lastnr" != x ];then
    state=$1
    case "$state" in
        start)
            echo "start action is right"
        ;;

        stop)
            echo "stop action is right"
        ;;

        restart|reload|force-reload)
            echo "reload action is right"
        ;;

        *)
            echo "Usage: $0 $lxcname {start|stop|restart|reload|force-reload}"
            exit 2
    esac
    // 检测网络是否正确
    network
    if [ $? -eq 0 ];then
        echo
        echo "容器内网络不畅通，请检查网络设置！"
        echo
        read -p "是否启用 ? [y/N]: " NET
        if [ "$NET" = 'y' -o "$NET" = 'Y' ]; then
            echo "先关闭容器内部网络，可能存在dnsmasq正在运行"
            ./lxc_reconnect_internet.sh  stop
            echo "开启容器内部网络"
            ./lxc_reconnect_internet.sh  start
        else
            exit 1;
        fi
    else
        echo
        echo "容器内网络非常畅通！"
        echo
        read -p "是否关闭 ? [y/N]: " NETSTOP
        if [ "$NETSTOP" = 'y' -o "$NETSTOP" = 'Y' ]; then
            ./lxc_reconnect_internet.sh  stop
        fi
    fi
    echo
    echo "检测过后容器内网络通畅"
    echo
    exitnr=$((lastnr+1))
    echo
    echo "哪一个容器你需要选择开启或关闭 openvpn和nsq?"
    echo
    text=$(lxc list|awk ' NR>2 && $2!="" && $2!="|" {print $2}'|awk '{print NR ")" $0}')
    echo "0) 所有容器"
    # echo $text #只能在同一行
    echo "$text"
    echo "$exitnr) Exit"
    read -p "Select an option [0-$exitnr]: " nr 
    # echo $nr
    if [ $nr = $exitnr ];then
        echo "退出了"
        exit;
    elif [ $nr = "0" ];then
        for container in $(lxc list|awk ' NR>2 && $2!="" && $2!="|" {print $2}'|awk '{print $0}'); do
            lxcname=$container
            # echo $container
            case "$state" in
                start)
                    ./lxc_restart.sh $lxcname
                ;;

                stop)
                    ./lxc_stop_vpn_nsq.sh $lxcname
                ;;

                restart|reload|force-reload)
                    ./lxc_stop_vpn_nsq.sh $lxcname
                    ./lxc_restart.sh $lxcname
                ;;

                *)
                    echo "Usage: $0 $lxcname {start|stop|restart|reload|force-reload}"
                    exit 2
            esac
        done
    else
        # sed -n "2p" 查找第2行内
        # echo $(lxc list|awk -F '|' ' NR%3==1 && $2!="" && NR>2 {print $2}'|sed -n "${nr}p")
        lxcname=$(lxc list|awk ' NR>2 && $2!="" && $2!="|" {print $2}'|awk '{print $0}'|sed -n "${nr}p")
        echo $lxcname
        case "$state" in
            start)
                ./lxc_restart.sh $lxcname
            ;;

            stop)
                ./lxc_stop_vpn_nsq.sh $lxcname
            ;;

            restart|reload|force-reload)
                ./lxc_stop_vpn_nsq.sh $lxcname
                ./lxc_restart.sh $lxcname
            ;;

            *)
                echo "Usage: $0 {start|stop|restart|reload|force-reload} "
                exit 2
        esac
    fi
fi
exit $?

```

##### lxc_reconnect_internet.sh
```shell
#!/bin/sh
read -p "请输入对应的lxc网卡:[默认(lxdbr0)] " NEW_LXC_BRIDGE
if [ x"$NEW_LXC_BRIDGE" = x ];then
    NEW_LXC_BRIDGE="lxdbr0"
fi
# echo $NEW_LXC_BRIDGE
read -p "请输入对应的lxc网卡地址:[默认(10.78.158.1)] " NEW_LXC_ADDR
if [ x"$NEW_LXC_ADDR" = x ];then
    NEW_LXC_ADDR="10.78.158.1"
fi
read -p "请输入对应的lxc网卡网段:[默认(10.78.158.0/24)] " NEW_LXC_NETWORK
if [ x"$NEW_LXC_NETWORK" = x ];then
    NEW_LXC_NETWORK="10.78.158.0/24"
fi
read -p "请输入对应的lxc网卡DHCP初始段:[默认(10.78.158.2)] " NEW_LXC_DHCP_FIRST
if [ x"$NEW_LXC_DHCP_FIRST" = x ];then
    NEW_LXC_DHCP_FIRST="10.78.158.2"
fi
read -p "请输入对应的lxc网卡DHCP最后段:[默认(10.78.158.254)] " NEW_LXC_DHCP_LAST
if [ x"$NEW_LXC_DHCP_LAST" = x ];then
    NEW_LXC_DHCP_LAST="10.78.158.254"
fi
NEW_LXC_DHCP_RANGE="${NEW_LXC_DHCP_FIRST},${NEW_LXC_DHCP_LAST}"

distrosysconfdir="/etc/default"
varrun="/run/lxc"
varlib="/var/lib"

# USE_LXC_BRIDGE="true"
# LXC_BRIDGE="lxdbr0"
# LXC_ADDR="10.78.158.1"
# LXC_NETMASK="255.255.255.0"
# LXC_NETWORK="10.78.158.0/24"
# LXC_DHCP_RANGE="10.78.158.2,10.78.158.254"
# LXC_DHCP_MAX="252"


USE_LXC_BRIDGE="true"
LXC_BRIDGE="lxdbr0"
LXC_BRIDGE_MAC="00:16:3e:00:00:00"
LXC_ADDR="10.0.3.1"
LXC_NETMASK="255.255.255.0"
LXC_NETWORK="10.0.3.0/24"
LXC_DHCP_RANGE="10.0.3.2,10.0.3.254"
LXC_DHCP_MAX="253"
LXC_DHCP_CONFILE=""
LXC_DOMAIN=""

LXC_IPV6_ADDR=""
LXC_IPV6_MASK=""
LXC_IPV6_NETWORK=""
LXC_IPV6_NAT="false"

write_lxc_net()
{
    local i=$1
    cat >>  $distrosysconfdir/lxc-net << EOF
# Leave USE_LXC_BRIDGE as "true" if you want to use lxcbr0 for your
# containers.  Set to "false" if you'll use virbr0 or another existing
# bridge, or mavlan to your host's NIC.
USE_LXC_BRIDGE="true"

# If you change the LXC_BRIDGE to something other than lxcbr0, then
# you will also need to update your /etc/lxc/default.conf as well as the
# configuration (/var/lib/lxc/<container>/config) for any containers
# already created using the default config to reflect the new bridge
# name.
# If you have the dnsmasq daemon installed, you'll also have to update
# and restart the system wide dnsmasq daemon.
LXC_BRIDGE="lxdbr0"
LXC_ADDR="10.78.158.1"
LXC_NETMASK="255.255.255.0"
LXC_NETWORK="10.78.158.0/24"
LXC_DHCP_RANGE="10.78.158.2,10.78.158.254"
LXC_DHCP_MAX="253"
# Uncomment the next line if you'd like to use a conf-file for the lxcbr0
# dnsmasq.  For instance, you can use 'dhcp-host=mail1,10.0.3.100' to have
# container 'mail1' always get ip address 10.0.3.100.
#LXC_DHCP_CONFILE=/etc/lxc/dnsmasq.conf

# Uncomment the next line if you want lxcbr0's dnsmasq to resolve the .lxc
# domain.  You can then add "server=/lxc/10.0.$i.1' (or your actual \$LXC_ADDR)
# to your system dnsmasq configuration file (normally /etc/dnsmasq.conf,
# or /etc/NetworkManager/dnsmasq.d/lxc.conf on systems that use NetworkManager).
# Once these changes are made, restart the lxc-net and network-manager services.
# 'container1.lxc' will then resolve on your host.
#LXC_DOMAIN="lxc"
EOF
}

configure_lxcbr0()
{
    local i=3
    cat >  $distrosysconfdir/lxc-net << EOF
# This file is auto-generated by lxc.postinst if it does not
# exist.  Customizations will not be overridden.
EOF
    # if lxcbr0 exists, keep using the same network
    if  ip addr show lxcbr0 > /dev/null 2>&1 ; then
        i=`ip addr show lxcbr0 | grep "inet\>" | awk '{ print $2 }' | awk -F. '{ print $3 }'`
        write_lxc_net $i
        return
    fi
    # if no lxcbr0, find an open 10.0.a.0 network
    for l in `ip addr show | grep "inet\>" |awk '{ print $2 }' | grep '^10\.0\.' | sort -n`; do
            j=`echo $l | awk -F. '{ print $3 }'`
            if [ $j -gt $i ]; then
                write_lxc_net $i
                return
            fi
            i=$((j+1))
    done
    if [ $i -ne 254 ]; then
        write_lxc_net $i
    fi
}

update_lxcnet_config()
{
    local i=3
    # if lxcbr0 exists, keep using the same network
    if  ip addr show lxcbr0 > /dev/null 2>&1 ; then
        return
    fi
    # our LXC_NET conflicts with an existing interface.  Probably first
    # run after system install with package pre-install.  Find a new subnet
    configure_lxcbr0

    # and re-load the newly created config
    [ ! -f $distrosysconfdir/lxc-net ] || . $distrosysconfdir/lxc-net
}

[ ! -f $distrosysconfdir/lxc ] || . $distrosysconfdir/lxc

use_iptables_lock="-w"
iptables -w -L -n > /dev/null 2>&1 || use_iptables_lock=""

_netmask2cidr ()
{
    # Assumes there's no "255." after a non-255 byte in the mask
    local x=${1##*255.}
    set -- 0^^^128^192^224^240^248^252^254^ $(( (${#1} - ${#x})*2 )) ${x%%.*}
    x=${1%%$3*}
    echo $(( $2 + (${#x}/4) ))
}

_ifdown() {
    ip addr flush dev ${LXC_BRIDGE}
    ip link set dev ${LXC_BRIDGE} down
}

_ifup() {
    MASK=`_netmask2cidr ${LXC_NETMASK}`
    CIDR_ADDR="${LXC_ADDR}/${MASK}"
    ip addr add ${CIDR_ADDR} dev ${LXC_BRIDGE}
    ip link set dev ${LXC_BRIDGE} address $LXC_BRIDGE_MAC
    ip link set dev ${LXC_BRIDGE} up
}

cleanup() {
    set +e
    if [ "$FAILED" = "1" ]; then
        echo "Failed to setup lxc-net." >&2
        stop force
        exit 1
    fi
}

start() {
    
    [ ! -f $distrosysconfdir/lxc-net ] && update_lxcnet_config

    LXC_BRIDGE=$NEW_LXC_BRIDGE
    LXC_ADDR=$NEW_LXC_ADDR
    LXC_NETWORK=$NEW_LXC_NETWORK
    LXC_DHCP_RANGE=$NEW_LXC_DHCP_RANGE

    # echo $NEW_LXC_NETWORK
    # exit 1;

    DNSMASQSRV=$(pgrep dnsmasq)
    if [ -z "$DNSMASQSRV" ]
    then
        echo "no dnsmasq found"
    else
        kill -9 $DNSMASQSRV
    fi
    [ "x$USE_LXC_BRIDGE" = "xtrue" ] || { exit 0; }

    [ ! -f "${varrun}/network_up" ] || { echo "lxc-net is already running"; exit 1; }

    if [ -d /sys/class/net/${LXC_BRIDGE} ]; then
        stop force || true
    fi

    FAILED=1

    trap cleanup EXIT HUP INT TERM
    set -e

    # set up the lxc network
    [ ! -d /sys/class/net/${LXC_BRIDGE} ] && ip link add dev ${LXC_BRIDGE} type bridge
    echo 1 > /proc/sys/net/ipv4/ip_forward
    echo 0 > /proc/sys/net/ipv6/conf/${LXC_BRIDGE}/accept_dad || true

    # if we are run from systemd on a system with selinux enabled,
    # the mkdir will create /run/lxc as init_var_run_t which dnsmasq
    # can't write its pid into, so we restorecon it (to var_run_t)
    if [ ! -d "${varrun}" ]; then
        mkdir -p "${varrun}"
        if which restorecon >/dev/null 2>&1; then
            restorecon "${varrun}"
        fi
    fi

    _ifup

    LXC_IPV6_ARG=""
    if [ -n "$LXC_IPV6_ADDR" ] && [ -n "$LXC_IPV6_MASK" ] && [ -n "$LXC_IPV6_NETWORK" ]; then
        echo 1 > /proc/sys/net/ipv6/conf/all/forwarding
        echo 0 > /proc/sys/net/ipv6/conf/${LXC_BRIDGE}/autoconf
        ip -6 addr add dev ${LXC_BRIDGE} ${LXC_IPV6_ADDR}/${LXC_IPV6_MASK}
        if [ "$LXC_IPV6_NAT" = "true" ]; then
            ip6tables $use_iptables_lock -t nat -A POSTROUTING -s ${LXC_IPV6_NETWORK} ! -d ${LXC_IPV6_NETWORK} -j MASQUERADE
        fi
        LXC_IPV6_ARG="--dhcp-range=${LXC_IPV6_ADDR},ra-only --listen-address ${LXC_IPV6_ADDR}"
    fi
    iptables $use_iptables_lock -I INPUT -i ${LXC_BRIDGE} -p udp --dport 67 -j ACCEPT
    iptables $use_iptables_lock -I INPUT -i ${LXC_BRIDGE} -p tcp --dport 67 -j ACCEPT
    iptables $use_iptables_lock -I INPUT -i ${LXC_BRIDGE} -p udp --dport 53 -j ACCEPT
    iptables $use_iptables_lock -I INPUT -i ${LXC_BRIDGE} -p tcp --dport 53 -j ACCEPT
    iptables $use_iptables_lock -I FORWARD -i ${LXC_BRIDGE} -j ACCEPT
    iptables $use_iptables_lock -I FORWARD -o ${LXC_BRIDGE} -j ACCEPT
    iptables $use_iptables_lock -t nat -A POSTROUTING -s ${LXC_NETWORK} ! -d ${LXC_NETWORK} -j MASQUERADE
    iptables $use_iptables_lock -t mangle -A POSTROUTING -o ${LXC_BRIDGE} -p udp -m udp --dport 68 -j CHECKSUM --checksum-fill

    LXC_DOMAIN_ARG=""
    if [ -n "$LXC_DOMAIN" ]; then
        LXC_DOMAIN_ARG="-s $LXC_DOMAIN -S /$LXC_DOMAIN/"
    fi

    LXC_DHCP_CONFILE_ARG=""
    if [ -n "$LXC_DHCP_CONFILE" ]; then
        LXC_DHCP_CONFILE_ARG="--conf-file=${LXC_DHCP_CONFILE}"
    fi

    # https://lists.linuxcontainers.org/pipermail/lxc-devel/2014-October/010561.html
    for DNSMASQ_USER in lxc-dnsmasq dnsmasq nobody
    do
        if getent passwd ${DNSMASQ_USER} >/dev/null; then
            break
        fi
    done

    dnsmasq $LXC_DHCP_CONFILE_ARG $LXC_DOMAIN_ARG -u ${DNSMASQ_USER} \
            --strict-order --bind-interfaces --pid-file="${varrun}"/dnsmasq.pid \
            --listen-address ${LXC_ADDR} --dhcp-range ${LXC_DHCP_RANGE} \
            --dhcp-lease-max=${LXC_DHCP_MAX} --dhcp-no-override \
            --except-interface=lo --interface=${LXC_BRIDGE} \
            --dhcp-leasefile="${varlib}"/misc/dnsmasq.${LXC_BRIDGE}.leases \
            --dhcp-authoritative $LXC_IPV6_ARG || cleanup

    touch "${varrun}"/network_up
    FAILED=0
}

stop() {
    
    [ "x$USE_LXC_BRIDGE" = "xtrue" ] || { exit 0; }

    [ -f "${varrun}/network_up" ] || [ "$1" = "force" ] || { echo "lxc-net isn't running"; exit 1; }

    LXC_BRIDGE=$NEW_LXC_BRIDGE
    LXC_ADDR=$NEW_LXC_ADDR
    LXC_NETWORK=$NEW_LXC_NETWORK
    LXC_DHCP_RANGE=$NEW_LXC_DHCP_RANGE

    if [ -d /sys/class/net/${LXC_BRIDGE} ]; then
        _ifdown 
        iptables $use_iptables_lock -D INPUT -i ${LXC_BRIDGE} -p udp --dport 67 -j ACCEPT
        iptables $use_iptables_lock -D INPUT -i ${LXC_BRIDGE} -p tcp --dport 67 -j ACCEPT
        iptables $use_iptables_lock -D INPUT -i ${LXC_BRIDGE} -p udp --dport 53 -j ACCEPT
        iptables $use_iptables_lock -D INPUT -i ${LXC_BRIDGE} -p tcp --dport 53 -j ACCEPT
        iptables $use_iptables_lock -D FORWARD -i ${LXC_BRIDGE} -j ACCEPT
        iptables $use_iptables_lock -D FORWARD -o ${LXC_BRIDGE} -j ACCEPT
        iptables $use_iptables_lock -t nat -D POSTROUTING -s ${LXC_NETWORK} ! -d ${LXC_NETWORK} -j MASQUERADE
        iptables $use_iptables_lock -t mangle -D POSTROUTING -o ${LXC_BRIDGE} -p udp -m udp --dport 68 -j CHECKSUM --checksum-fill

        if [ "$LXC_IPV6_NAT" = "true" ]; then
            ip6tables $use_iptables_lock -t nat -D POSTROUTING -s ${LXC_IPV6_NETWORK} ! -d ${LXC_IPV6_NETWORK} -j MASQUERADE
        fi

        pid=`cat "${varrun}"/dnsmasq.pid 2>/dev/null` && kill -9 $pid
        rm -f "${varrun}"/dnsmasq.pid
        # if $LXC_BRIDGE has attached interfaces, don't destroy the bridge
        ls /sys/class/net/${LXC_BRIDGE}/brif/* > /dev/null 2>&1 || ip link delete ${LXC_BRIDGE}
    fi

    rm -f "${varrun}"/network_up
}

# See how we were called.
case "$1" in
    start)
        start
    ;;

    stop)
        stop
    ;;

    restart|reload|force-reload)
        $0 stop
        $0 start
    ;;

    *)
        echo "Usage: $0 {start|stop|restart|reload|force-reload}"
        exit 2
esac

exit $?


```

##### lxc_restart.sh
``` shell
#!/bin/sh
#./opevpn/lxc_reconnect_internet.sh
CONTAINERNAME="s1"
NSQFOLDER="/root/nsq1.2"
# read -p "请输入对应的lxc容器名称: " lxcname
lxcname=$1
# echo $lxcname
if [ -n "$NSQ" ]; then 
    echo "容器$lxcname 不存在"
    exit 1;
fi
if [ ! -f "/root/$lxcname.ovpn" ]; then 
    echo "$lxcname.ovpn 不存在,请用openinstall.sh 添加新用户"
    exit 1;
fi
OPENVPN=`lxc exec $lxcname pgrep openvpn`       
if [ -n "$OPENVPN" ];then  
    echo "openvpn service is running"
else
    echo "openvpn service is not running,ok then can do this"
    lxc file push /root/$lxcname.ovpn $lxcname/root/
    lxc exec $lxcname  -- sh -c "apt-get install -y openvpn"
    # lxc exec test -- nohup bash -c "/root/test.sh &"
    # lxc exec $lxcname -- sh -c "nohup openvpn --config /root/$lxcname.ovpn &" #不生效
    lxc exec $lxcname -- nohup bash -c " openvpn --config /root/$lxcname.ovpn &"
fi
echo "等待5s openvpn启动"
sleep 5

lookupdip=`lxc exec ${CONTAINERNAME} ifconfig tun0 | grep "inet addr:" | awk '{print $2}' | cut -c 6-` 
echo "lookupdip openvpn ip is : $lookupdip"
NSQ=`lxc exec $lxcname pgrep nsq`       
if [ -n "$NSQ" ]; then 
        echo "nsq service is running"
        exit 1;
    else
        echo "nsq service is not running,we can do this;"
        lxc file push /root/nsq-1.2.0.linux-amd64.go1.12.9.tar.gz $lxcname/root/
        lxc exec $lxcname -- sh -c "tar -zxvf /root/nsq-1.2.0.linux-amd64.go1.12.9.tar.gz" 
        lxc exec $lxcname -- sh -c "rm -rf ${NSQFOLDER}"
        lxc exec $lxcname -- sh -c "mv /root/nsq-1.2.0.linux-amd64.go1.12.9 ${NSQFOLDER}"
        if [ $lxcname = ${CONTAINERNAME} ];then 
            echo "xxxx"
            #lxc exec s1 -- nohup bash -c "/root/nsq1.2/bin/nsqlookupd &"
            lxc exec $lxcname -- nohup bash -c "/${NSQFOLDER}/bin/nsqlookupd &"
            lxc exec $lxcname -- nohup bash -c "/${NSQFOLDER}/bin/nsqd --lookupd-tcp-address=127.0.0.1:4160 &"
            lxc exec $lxcname -- nohup bash -c "/${NSQFOLDER}/bin/nsqadmin --lookupd-http-address=127.0.0.1:4161 &"
        else
            lxc exec $lxcname -- nohup bash -c "/${NSQFOLDER}/bin/nsqd --lookupd-tcp-address=$lookupdip:4160 &"
        fi
fi
```


##### lxc_stop_vpn_nsq.sh
``` shell
#!/bin/sh
stopfile="stop_vpn_nsq.sh"
lxcname=$1
# read -p "请输入对应的lxc容器名称: " lxcname
echo $lxcname
lxc file push /root/openvpn/$stopfile $lxcname/root/
lxc exec $lxcname -- bash -c "/root/$stopfile "


```