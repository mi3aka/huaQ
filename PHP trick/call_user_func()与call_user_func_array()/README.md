## call_user_func

函数定义`call_user_func(callable $callback, mixed $parameter = ?, mixed $... = ?): mixed`

第一个参数`callback`是被调用的回调函数,其余参数是回调函数的参数

>传入call_user_func()的参数不能为引用传递

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202232126748.png)

```php
<?php
function out($a, $b)
{
    echo "$a $b\n";
}
call_user_func('out', "1", "2");
call_user_func('out', "3", "4");
```

命名空间的使用

```php
<?php
namespace Foobar;
class Foo {
    static public function test() {
        print "Hello world!\n";
    }
}
call_user_func(__NAMESPACE__ .'\Foo::test');
call_user_func(array(__NAMESPACE__ .'\Foo', 'test'));
```

调用类中的方法

```php
<?php

class myclass
{
    static function out($str)
    {
        echo "$str\n";
    }
}

$classname = "myclass";
call_user_func(array($classname, 'out'), "asdf");
call_user_func($classname . '::out', "asdf");
$myobject = new myclass();
call_user_func(array($myobject, 'out'), "asdf");
```

匿名函数

```php
<?php
call_user_func(function($arg) { echo $arg; }, 'test');
```

## call_user_func_array

函数定义`call_user_func_array(callable $callback, array $param_arr): mixed`

第一个参数作为回调函数`callback`调用,把参数数组作`param_arr`为回调函数的的参数传入

>参数可以为引用传递

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202232127926.png)

```php
<?php
function out($a, $b)
{
    echo "$a $b\n";
}
call_user_func_array('out', array("1", "2"));
call_user_func_array('out', array("3", "4"));
```

命名空间的使用

```php
<?php
namespace Foobar;
class Foo {
    static public function test() {
        print "Hello world!\n";
    }
}
call_user_func_array(__NAMESPACE__ .'\Foo::test',array());
call_user_func_array(array(__NAMESPACE__ .'\Foo', 'test'),array());
```

调用类中的方法

```php
<?php

class myclass
{
    static function out($str)
    {
        echo "$str\n";
    }
}

$classname = "myclass";
call_user_func_array(array($classname, 'out'), array("asdf"));
call_user_func_array($classname . '::out', array("asdf"));
$myobject = new myclass();
call_user_func_array(array($myobject, 'out'), array("asdf"));
```

匿名函数

```php
<?php
call_user_func_array(function($arg) { echo $arg; }, array('test'));
```