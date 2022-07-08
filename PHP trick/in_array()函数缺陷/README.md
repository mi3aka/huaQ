`in_array`用于检查数组中是否存在某个值

函数定义 `in_array(mixed $needle, array $haystack, bool $strict = false): bool`

- `needle` 待搜索的值
- `haystack` 待搜索的数组
- `strict` 如果`strict`的值为`true`则`in_array()`函数会对`needle`进行类型检查

```php
<?php
var_dump(array(0, 1, 2, 3));
var_dump(in_array(0, array(0, 1, 2, 3)));# true
var_dump(in_array(0, array(0, 1, 2, 3), true));# true
var_dump(in_array('0', array(0, 1, 2, 3)));# true
var_dump(in_array('0', array(0, 1, 2, 3), true));# false
var_dump(in_array('4', array(0, 1, 2, 3)));# false
var_dump(in_array('4', array(0, 1, 2, 3), true));# false
var_dump(in_array('hua1Q', array(0, 1, 2, 3)));# true
var_dump(in_array('hua1Q', array(0, 1, 2, 3), true));# false
?>
```

```php
array (size=4)
  0 => int 0
  1 => int 1
  2 => int 2
  3 => int 3
boolean true
boolean true
boolean true
boolean false
boolean false
boolean false
boolean true
boolean false
```

当`strict`没有被设置为`true`时,`'0','4','hua1Q'`被转换为`0,4,1`后进行`in_array`判断

例题 php-security-calendar-2017 Day 1 - Wish List

```php
<?php
highlight_file(__FILE__);
define ('SITE_ROOT', realpath(dirname(__FILE__)));

class Challenge
{
    const UPLOAD_DIRECTORY = SITE_ROOT.'/solutions/';
    private $file;
    private $whitelist;

    public function __construct($file)
    {
        if(!is_dir('solutions')){
            mkdir('solutions');
        }
        $this->file = $file;
        $this->whitelist = range(1, 24);
    }

    public function __destruct()
    {
        if (in_array($this->file['name'], $this->whitelist)) {
            move_uploaded_file($this->file['tmp_name'], self::UPLOAD_DIRECTORY . $this->file['name']);
        }
    }
}

if (isset($_POST['submit'])) {
    $challenge = new Challenge($_FILES['solution']);
}
?>
<form action="" method="post" enctype="multipart/form-data">
    <label for="file">Filename:</label>
    <input type="file" name="solution" id="file"/>
    <br/>
    <input type="submit" name="submit" value="Submit"/>
</form>
```

`in_array($this->file['name'], $this->whitelist)`限制了文件名必须是`1~24`,但是`strict`参数没有被设置为`true`,因此可以上传如`1.php`等文件来绕过限制

除了能够用于文件上传绕过,还可以用于sql注入

```php
<?php
$num = array(1, 2, 3, 4);
$input = $_POST["a"];
if (!isset($input) or !in_array($input, $num)) {
    return false;
}
$query = 'SELECT id,username FROM user WHERE id="' . $input . '"';
var_dump($query);

$db = new mysqli("172.16.172.202", "user", "password", "www");
$result = $db->query($query);
if ($result->num_rows) {
    while ($row = $result->fetch_assoc()) {
        echo "id:" . $row["id"] . " username:" . $row["username"] . "<br>";
    }
} else {
    echo "NULL";
}
$db->close();
?>
```

![](./截屏2021-08-30%2019.53.57.png)

![](./截屏2021-08-30%2019.56.49.png)

![](./截屏2021-08-30%2019.57.42.png)