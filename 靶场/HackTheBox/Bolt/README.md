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

`templates/emails/update-name.html`

```html
<html>
	<body>
		<p> %s </p>
		<p> This e-mail serves as confirmation of your profile username changes.</p>
	</body>
</html>
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202201292336044.png)

