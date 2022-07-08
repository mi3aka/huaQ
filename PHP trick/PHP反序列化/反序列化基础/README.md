## 序列化

序列化`public`,`private`,`protected`变量产生不同结果

```php
<?php
class test{
    private $test1="hello";
    public $test2="hello";
    protected $test3="hello";
}
$test = new test();
echo serialize($test);
echo "\n";
echo urlencode(serialize($test));
?>
```

```
O:4:"test":3:{s:11:" test test1";s:5:"hello";s:5:"test2";s:5:"hello";s:8:" * test3";s:5:"hello";}
O%3A4%3A%22test%22%3A3%3A%7Bs%3A11%3A%22%00test%00test1%22%3Bs%3A5%3A%22hello%22%3Bs%3A5%3A%22test2%22%3Bs%3A5%3A%22hello%22%3Bs%3A8%3A%22%00%2A%00test3%22%3Bs%3A5%3A%22hello%22%3B%7D
```

private变量`s:11:"testtest1"`的实际表示为`s:11:"\00test\00test1"`

protected变量`s:8:"*test3"`的实际表示为`s:8:"\00*\00test3"`

## 反序列化

```php
<?php
$str='O%3A4%3A%22test%22%3A3%3A%7Bs%3A11%3A%22%00test%00test1%22%3Bs%3A5%3A%22hello%22%3Bs%3A5%3A%22test2%22%3Bs%3A5%3A%22hello%22%3Bs%3A8%3A%22%00%2A%00test3%22%3Bs%3A5%3A%22hello%22%3B%7D';
$data=urldecode($str);
$obj=unserialize($data);
var_dump($obj);
?>
```

```php
object(__PHP_Incomplete_Class)[1]
  public '__PHP_Incomplete_Class_Name' => string 'test' (length=4)
  private 'test1' (test) => string 'hello' (length=5)
  public 'test2' => string 'hello' (length=5)
  protected 'test3' => string 'hello' (length=5)
```

**序列化一个对象将会保存对象的所有变量,但是不会保存对象的方法,只会保存类的名字**

**PHP在反序列化时,对类中不存在的属性也会进行反序列化**

```php
<?php
class A{
    public $var1="asdf";
    public $var2="asdf";
}
$a=new A();
var_dump(serialize($a));
$b='O:1:"A":2:{s:4:"var1";s:4:"asdf";s:4:"var3";s:4:"asdf";}';
var_dump(unserialize($b));
```

```php
string 'O:1:"A":2:{s:4:"var1";s:4:"asdf";s:4:"var2";s:4:"asdf";}' (length=56)
object(A)[2]
  public 'var1' => string 'asdf' (length=4)
  public 'var2' => string 'asdf' (length=4)
  public 'var3' => string 'asdf' (length=4)
```

## 序列化字符串中的字母含义

|字母|解释|备注|
|:---:|:---:|:---:|
|a:array|数组||
|b:boolean|布尔值||
|d:double|浮点数||
|i:integer|整数||
|r:reference|对象引用||
|s:string|字符串|大写的S(使用转义字符)替换小写的s即可用16进制表示|
|O:class|普通类||
|N:null|NULL||
|R:pointer reference|指针引用||
|U:unicode string|Unicode字符串||

**在反序列化时用大写的S替换小写的s即可用16进制表示**

```php
<?php
$str1='O:4:"test":1:{s:4:"asdf";s:4:"qwer";}';
var_dump(unserialize($str1));
$str2='O:4:"test":1:{s:4:"asdf";S:4:"\71wer";}';
var_dump(unserialize($str2));
```

```php
object(__PHP_Incomplete_Class)[1]
  public '__PHP_Incomplete_Class_Name' => string 'test' (length=4)
  public 'asdf' => string 'qwer' (length=4)
object(__PHP_Incomplete_Class)[1]
  public '__PHP_Incomplete_Class_Name' => string 'test' (length=4)
  public 'asdf' => string 'qwer' (length=4)
```

```php
<?php

class A
{
    var $value;
}

$a = new A();
$a->value = $a;
$b = new A();
$b->value = &$b;
var_dump(serialize($a));
var_dump(serialize($b));
```

```php
string(28) "O:1:"A":1:{s:5:"value";r:1;}"
string(28) "O:1:"A":1:{s:5:"value";R:1;}"
```

[关于对象引用r和指针引用R](https://hujiekang.top/2020/09/25/PHP-unserialize-advanced/)

这两者在引用方式上是有区别的,可以理解为对象引用是一个单边的引用,被赋值的那个变量可以任意修改值,而不会影响到被引用的那个对象

而指针引用则是一个双边的引用,被赋值的那个变量若做了改动,被引用的那个对象也会被修改,也就是说指针引用其实就是两个对象指针指向了同一块内存区域,所以任一指针的数值修改其实都是在对这块内存做修改,也就会影响到另一个指针的值

而对象引用的被赋值对象就像一个临时的指针,指向了被引用对象的内存区域,而当被赋值对象的值修改之后,这个临时指针就指向了另一块内存

- 对象引用r和指针引用R的引用顺序

在序列化时,r或R的后面还会有个数字,这个数字代表的就是引用的顺序

`O:4:"test":3:{s:4:"left";N;s:6:"middle";R:2;s:5:"right";R:2;}`

可以看见上面这个序列化对象中,left成员为NULL,middle和right都是指针引用,引用后跟着的数字都是2

这个数字就是所引用的对象在序列化串中第一次出现的位置,但是这个位置不是指字符的位置,而是指对象**这里的对象是泛指所有类型的量,而不仅限于对象类型**的位置

上面那个对象,首先被序列化的肯定是整个对象,也就是`test`,所以`O:4:"test"`的引用序号就是1

随后按顺序序列化它的成员变量,第一个就是`left`,所以它的引用序号就是`2`

以此类推,若后面两个成员不是引用的话,对应的引用序号也就接着递增,但是当序列化到middle时,发现它指向的对象已经被序列化了,也就是left,所以给它标上引用序号2,同理right也是2

所以实际上,这个对象的三个成员指向的都是同一块内存区域,代表的都是同一个对象

例子可见[2019年新生赛题目unserialize](https://github.com/mi3aka/xp0int-2019-ctf-web/tree/master/unserialize)

- 使用`+`绕过正则

```php
preg_match('/[oc]:\d+:/i', $var)
```

对`O:1`进行匹配,用`O:+1`进行绕过