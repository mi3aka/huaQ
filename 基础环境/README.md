大部分基础环境均运行在Docker上

本机为黑苹果,Docker运行于VMware上的一个Debian Server中,Debian设置了静态ip为172.16.172.202

VMware设置如下

在`/Library/Preferences/VMware Fusion/vmnet8/dhcpd.conf`中追加如下内容

```
host Server{
truehardware ethernet 00:50:56:35:FB:00;
truefixed-address 172.16.172.202;
}
```

`host xxx`为vmware虚拟机的名字,`truehardware ethernet`是mac地址,`truefixed-address`是要为其设定的静态ip(注意ip的有效范围)

![](./截屏2021-08-30%2020.17.14.png)

---

Debian设置如下

![](./截屏2021-08-30%2020.20.03.png)

修改`/lib/systemd/system/docker.service`中的`ExecStart=/usr/bin/dockerd -H fd:`为以下内容

```
ExecStart=/usr/bin/dockerd -H unix:///var/run/docker.sock -H tcp://0.0.0.0:2375
```

修改完成后重新加载配置并重启docker

```
systemctl daemon-reload
service docker restart
```

1. php74

利用[whistle](https://wproxy.org/whistle/)将流量转发到Debian中从而进行xdebug调试

![](截屏2021-08-30%2020.27.02.png)