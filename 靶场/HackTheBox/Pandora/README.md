nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201221513433.png)

---

发现`161/udp`开放,尝试使用`snmpwalk`进行读取

`sudo nmap -sU -p 161 --script=snmp-processes 10.10.11.136`

```
Starting Nmap 7.80 ( https://nmap.org ) at 2022-01-22 15:27 CST
Nmap scan report for 10.10.11.136
Host is up (0.081s latency).

PORT    STATE SERVICE
161/udp open  snmp
| snmp-processes: 
|   1: 
|     Name: systemd
|     Path: /sbin/init
|     Params: maybe-ubiquity
...
|   838: 
|     Name: cron
|     Path: /usr/sbin/CRON
|     Params: -f
|   843: 
|     Name: sh
|     Path: /bin/sh
|     Params: -c sleep 30; /bin/bash -c '/usr/bin/host_check -u daniel -p HotelBabylon23'
|   858: 
|     Name: atd
|     Path: /usr/sbin/atd
|     Params: -f
|   860: 
|     Name: snmpd
|     Path: /usr/sbin/snmpd
|     Params: -LOw -u Debian-snmp -g Debian-snmp -I -smux mteTrigger mteTriggerConf -f -p /run/snmpd.pid
|   861: 
|     Name: sshd
|     Path: sshd: /usr/sbin/sshd -D [listener] 1 of 10-100 startups
|   872: 
|     Name: apache2
|     Path: /usr/sbin/apache2
|     Params: -k start
|   897: 
|     Name: agetty
|     Path: /sbin/agetty
|     Params: -o -p -- \u --noclear tty1 linux
|   909: 
|     Name: polkitd
|     Path: /usr/lib/policykit-1/polkitd
|     Params: --no-debug
|   949: 
|     Name: mysqld
|     Path: /usr/sbin/mysqld
|   1131: 
|     Name: host_check
|     Path: /usr/bin/host_check
|     Params: -u daniel -p HotelBabylon23
...
|   31386: 
|     Name: sshd
|     Path: sshd: daniel@pts/0
|   31387: 
|     Name: bash
|     Path: -bash
|   31589: 
|     Name: chisel_1.7.6_li
|     Path: ./chisel_1.7.6_linux_amd64
|     Params: client 10.10.14.21:1234 R:30000:127.0.0.1:80
|   36786: 
|     Name: kworker/0:0-memcg_kmem_cache
|   36823: 
|     Name: kworker/1:2-events
|   37386: 
|     Name: kworker/1:1-events
|   37522: 
|     Name: kworker/u4:0-events_power_efficient
|   42454: 
...
| 
|_  71459: 

Nmap done: 1 IP address (1 host up) scanned in 357.43 seconds
```

敏感信息`/usr/bin/host_check -u daniel -p HotelBabylon23`,`Params: -u daniel -p HotelBabylon23`,`chisel_1.7.6_linux_amd64`

`user.txt`在`/home/matt`目录下,不在`/home/daniel`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201221539616.png)

```
find / -user root -perm -4000 -print 2>/dev/null
/usr/bin/sudo
/usr/bin/pkexec
/usr/bin/chfn
/usr/bin/newgrp
/usr/bin/gpasswd
/usr/bin/umount
/usr/bin/pandora_backup
/usr/bin/passwd
/usr/bin/mount
/usr/bin/su
/usr/bin/fusermount
/usr/bin/chsh
/usr/lib/openssh/ssh-keysign
/usr/lib/dbus-1.0/dbus-daemon-launch-helper
/usr/lib/eject/dmcrypt-get-device
/usr/lib/policykit-1/polkit-agent-helper-1
```

`/usr/bin/host_check`反编译

```cpp
int __cdecl main(int argc, const char **argv, const char **envp)
{
  int result; // eax

  if ( argc == 5 )
  {
    puts("PandoraFMS host check utility");
    puts("Now attempting to check PandoraFMS registered hosts.");
    puts("Files will be saved to ~/.host_check");
    if ( system(
           "/usr/bin/curl 'http://127.0.0.1/pandora_console/include/api.php?op=get&op2=all_agents&return_type=csv&other_m"
           "ode=url_encode_separator_%7C&user=daniel&pass=HotelBabylon23' > ~/.host_check 2>/dev/null") )
    {
      printf("Host check unsuccessful!\nPlease check your credentials.\nTerminating program!");
      result = 1;
    }
    else
    {
      sleep(0x1B186Eu);
      puts("Host check successful!\nTerminating program!");
      result = 0;
    }
  }
  else
  {
    if ( argc > 4 )
      puts("Two arguments expected.");
    else
      puts("Ussage: ./host_check -u username -p password.");
    result = 0;
  }
  return result;
}
```

在`/var/www/`有`html`和`pandora`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201221547490.png)

`ssh -L 80:127.0.0.1:80 daniel@10.10.11.136`将端口转发到本地(注意权限)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201221614269.png)

[http://127.0.0.1/pandora_console/include/api.php?op=get&op2=all_agents&return_type=csv&other_m%22%22ode=url_encode_separator_%7C&user=daniel&pass=HotelBabylon23](http://127.0.0.1/pandora_console/include/api.php?op=get&op2=all_agents&return_type=csv&other_m%22%22ode=url_encode_separator_%7C&user=daniel&pass=HotelBabylon23)

```
1;localhost.localdomain;192.168.1.42;Created by localhost.localdomain;Linux;;09fbaa6fdf35afd44f8266676e4872f299c1d3cbb9846fbe944772d913fcfc69;3
2;localhost.localdomain;;Pandora FMS Server version 7.0NG.742_FIX_PERL2020;Linux;;localhost.localdomain;3
```

这个FMS有个sql注入的漏洞[https://blog.sonarsource.com/pandora-fms-742-critical-code-vulnerabilities-explained](https://blog.sonarsource.com/pandora-fms-742-critical-code-vulnerabilities-explained)

