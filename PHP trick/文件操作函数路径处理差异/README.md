`file_put_contents`,`file_get_contents`等对文件内容进行处理的函数,在处理传入的文件路径时,**会将其转换成绝对路径**

`unlink`,`file_exists`,`rename`等文件操作函数,在处理传入的文件路径时,**不会将其转换成绝对路径**

例子

```php
<?php
$filename=md5(time())."/../test.txt";
var_dump($filename);
file_put_contents($filename,"test");
var_dump(file_get_contents($filename));
system("cat test.txt");
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202261452516.png)

---

当`unlink`函数在处理`$filename=md5(time())."/../test.txt"`这种路径时会发生报错`No such file or directory`,无法正常处理

而`file_exist`函数在处理`$filename=md5(time())."/../test.txt"`时,尽管文件确实存在,但是却会返回`false`

例子

```php
<?php
$filename=md5(time())."/../test.txt";
var_dump($filename);
file_put_contents($filename,"test");

var_dump(file_exists($filename));
var_dump(file_exists("test.txt"));
unlink($filename);

var_dump(file_get_contents($filename));
system("cat test.txt");
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202261457197.png)

---

源码分析

调用链如下图,版本为7.4.26

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202261619214.png)

全局搜索`file_put_contents`,在`file.c`中的`stream = php_stream_open_wrapper_ex(filename, mode, ((flags & PHP_FILE_USE_INCLUDE_PATH) ? USE_PATH : 0) | REPORT_ERRORS, NULL, context);`下断点并跟进

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202261621069.png)

经过一系列的`stream_open`后,进行`expand_filepath`

```cpp
	if (IS_ABSOLUTE_PATH(filepath, path_len)) {
		cwd[0] = '\0';
	} else {
		const char *iam = SG(request_info).path_translated;
		const char *result;
		if (relative_to) {
			if (relative_to_len > MAXPATHLEN-1U) {
				return NULL;
			}
			result = relative_to;
			memcpy(cwd, relative_to, relative_to_len+1U);
		} else {
			result = VCWD_GETCWD(cwd, MAXPATHLEN);
		}

		if (!result && (iam != filepath)) {
			int fdtest = -1;

			fdtest = VCWD_OPEN(filepath, O_RDONLY);
			if (fdtest != -1) {
				/* return a relative file path if for any reason
				 * we cannot cannot getcwd() and the requested,
				 * relatively referenced file is accessible */
				copy_len = path_len > MAXPATHLEN - 1 ? MAXPATHLEN - 1 : path_len;
				if (real_path) {
					memcpy(real_path, filepath, copy_len);
					real_path[copy_len] = '\0';
				} else {
					real_path = estrndup(filepath, copy_len);
				}
				close(fdtest);
				return real_path;
			} else {
				cwd[0] = '\0';
			}
		} else if (!result) {
			cwd[0] = '\0';
		}
	}
```

先判断是否为绝对路径,不是则进行路径延展

`result = VCWD_GETCWD(cwd, MAXPATHLEN);`得到当前路径

跟进到`virtual_file_ex`中

```cpp
memcpy(resolved_path, state->cwd, state_cwd_length);//将当前绝对路径复制到resolved_path中
if (resolved_path[state_cwd_length-1] == DEFAULT_SLASH) {
	memcpy(resolved_path + state_cwd_length, path, path_length + 1);
	path_length += state_cwd_length;
} else {
	resolved_path[state_cwd_length] = DEFAULT_SLASH;
	memcpy(resolved_path + state_cwd_length + 1, path, path_length + 1);//将传入的路径复制到resolved_path后面,即进行路径延展
	path_length += state_cwd_length + 1;
}
```

`path_length = tsrm_realpath_r(resolved_path, start, path_length, &ll, &t, use_realpath, 0, NULL);`进入`tsrm_realpath_r`函数,对路径递归处理

```cpp
		i = len;
		while (i > start && !IS_SLASH(path[i-1])) {//确定路径中最后一个斜杠/的位置,进行文件名与路径的切割
			i--;
		}
		assert(i < MAXPATHLEN);

		if (i == len ||
			(i + 1 == len && path[i] == '.')) {//处理index.php/.这种情况
			/* remove double slashes and '.' */
			len = EXPECTED(i > 0) ? i - 1 : 0;
			is_dir = 1;
			continue;
		} else if (i + 2 == len && path[i] == '.' && path[i+1] == '.') {//处理..这种情况
			/* remove '..' and previous directory */
			is_dir = 1;
			if (link_is_dir) {
				*link_is_dir = 1;
			}
			if (i <= start + 1) {
				return start ? start : len;
			}
			j = tsrm_realpath_r(path, start, i-1, ll, t, use_realpath, 1, NULL);
			if (j > start && j != (size_t)-1) {
				j--;
				assert(i < MAXPATHLEN);
				while (j > start && !IS_SLASH(path[j])) {
					j--;
				}
				assert(i < MAXPATHLEN);
				if (!start) {
					/* leading '..' must not be removed in case of relative path */
					if (j == 0 && path[0] == '.' && path[1] == '.' &&
							IS_SLASH(path[2])) {
						path[3] = '.';
						path[4] = '.';
						path[5] = DEFAULT_SLASH;
						j = 5;
					} else if (j > 0 &&
							path[j+1] == '.' && path[j+2] == '.' &&
							IS_SLASH(path[j+3])) {
						j += 4;
						path[j++] = '.';
						path[j++] = '.';
						path[j] = DEFAULT_SLASH;
					}
				}
			} else if (!start && !j) {
				/* leading '..' must not be removed in case of relative path */
				path[0] = '.';
				path[1] = '.';
				path[2] = DEFAULT_SLASH;
				j = 2;
			}
			return j;
		}
```

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202261657890.png)

![](https://cdn.jsdelivr.net/gh/AMDyesIntelno/PicGoImg@master/202202261723732.png)

得到`real_path`,而`unlink`等函数没有对`..`进行处理

```php
<?php
highlight_file(__FILE__);
if (isset($_POST["name"]) and isset($_POST["data"])) {
    $name = $_POST["name"];
    $data = $_POST["data"];
    file_put_contents($name, $data);
    if (file_exists($name)) {
        unlink($name);
    }
}
```

传入`name=asdf/../a.php&data=<?php phpinfo();?>`或`name=a.php/.&data=<?php phpinfo();?>`即可