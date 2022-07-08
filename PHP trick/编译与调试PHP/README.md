[c语言调试方法（调试 PHP 底层、扩展）](https://learnku.com/articles/33565)

[PHP 源码阅读笔记：编译与调试 PHP](https://her-cat.com/2021/08/06/php-internals-compile-and-debugging.html)

[命令执行底层原理探究-PHP](https://qftm.github.io/2020/12/01/command-execution-research-php/)

从[https://www.php.net/downloads.php](https://www.php.net/downloads.php)处下载源代码

我选择的是[php-7.4.26.tar.gz](https://www.php.net/distributions/php-7.4.26.tar.gz)

下载完成后,将其解压到`~/PHP-Debug`目录下,在`~/PHP-Debug/php-7.4.26`中运行如下命令

`./configure --prefix=$HOME/PHP-Debug/php-7.4.26/build --disable-all --enable-debug --enable-cli --enable-phar`

`make -j8 && make install`

等待执行完成后,可以在`$HOME/PHP-Debug/php-7.4.26/build/bin`目录下看到相关可执行文件

![](2021-12-07%2015-02-00%20的屏幕截图.png)

![](2021-12-07%2015-27-58%20的屏幕截图.png)

在`sapi/cli/php_cli.c`的`main`函数中下断点即可

![](2021-12-07%2019-19-49%20的屏幕截图.png)

---

```php
<?php
var_dump("Hello");
?>
```

在`ext/standard/var.c`中的`PHPAPI void php_var_dump(zval *struc, int level)`下断点

```cpp
PHPAPI void php_var_dump(zval *struc, int level) /* {{{ */
{
	HashTable *myht;
	zend_string *class_name;
	int is_ref = 0;
	zend_ulong num;
	zend_string *key;
	zval *val;
	uint32_t count;

	if (level > 1) {
		php_printf("%*c", level - 1, ' ');
	}

again:
	switch (Z_TYPE_P(struc)) {
		case IS_FALSE:
			php_printf("%sbool(false)\n", COMMON);
			break;
		case IS_TRUE:
			php_printf("%sbool(true)\n", COMMON);
			break;
		case IS_NULL:
			php_printf("%sNULL\n", COMMON);
			break;
		case IS_LONG:
			php_printf("%sint(" ZEND_LONG_FMT ")\n", COMMON, Z_LVAL_P(struc));
			break;
		case IS_DOUBLE:
			php_printf("%sfloat(%.*G)\n", COMMON, (int) EG(precision), Z_DVAL_P(struc));
			break;
		case IS_STRING:
			php_printf("%sstring(%zd) \"", COMMON, Z_STRLEN_P(struc));
			PHPWRITE(Z_STRVAL_P(struc), Z_STRLEN_P(struc));
			PUTS("\"\n");
			break;
		case IS_ARRAY:
			myht = Z_ARRVAL_P(struc);
			if (!(GC_FLAGS(myht) & GC_IMMUTABLE)) {
				if (level > 1) {
					if (GC_IS_RECURSIVE(myht)) {
						PUTS("*RECURSION*\n");
						return;
					}
					GC_PROTECT_RECURSION(myht);
				}
				GC_ADDREF(myht);
			}
			count = zend_array_count(myht);
			php_printf("%sarray(%d) {\n", COMMON, count);
			ZEND_HASH_FOREACH_KEY_VAL_IND(myht, num, key, val) {
				php_array_element_dump(val, num, key, level);
			} ZEND_HASH_FOREACH_END();
			if (!(GC_FLAGS(myht) & GC_IMMUTABLE)) {
				if (level > 1) {
					GC_UNPROTECT_RECURSION(myht);
				}
				GC_DELREF(myht);
			}
			if (level > 1) {
				php_printf("%*c", level-1, ' ');
			}
			PUTS("}\n");
			break;
		case IS_OBJECT:
			if (Z_IS_RECURSIVE_P(struc)) {
				PUTS("*RECURSION*\n");
				return;
			}
			Z_PROTECT_RECURSION_P(struc);

			myht = zend_get_properties_for(struc, ZEND_PROP_PURPOSE_DEBUG);
			class_name = Z_OBJ_HANDLER_P(struc, get_class_name)(Z_OBJ_P(struc));
			php_printf("%sobject(%s)#%d (%d) {\n", COMMON, ZSTR_VAL(class_name), Z_OBJ_HANDLE_P(struc), myht ? zend_array_count(myht) : 0);
			zend_string_release_ex(class_name, 0);

			if (myht) {
				zend_ulong num;
				zend_string *key;
				zval *val;

				ZEND_HASH_FOREACH_KEY_VAL(myht, num, key, val) {
					zend_property_info *prop_info = NULL;

					if (Z_TYPE_P(val) == IS_INDIRECT) {
						val = Z_INDIRECT_P(val);
						if (key) {
							prop_info = zend_get_typed_property_info_for_slot(Z_OBJ_P(struc), val);
						}
					}

					if (!Z_ISUNDEF_P(val) || prop_info) {
						php_object_property_dump(prop_info, val, num, key, level);
					}
				} ZEND_HASH_FOREACH_END();
				zend_release_properties(myht);
			}
			if (level > 1) {
				php_printf("%*c", level-1, ' ');
			}
			PUTS("}\n");
			Z_UNPROTECT_RECURSION_P(struc);
			break;
		case IS_RESOURCE: {
			const char *type_name = zend_rsrc_list_get_rsrc_type(Z_RES_P(struc));
			php_printf("%sresource(%d) of type (%s)\n", COMMON, Z_RES_P(struc)->handle, type_name ? type_name : "Unknown");
			break;
		}
		case IS_REFERENCE:
			//??? hide references with refcount==1 (for compatibility)
			if (Z_REFCOUNT_P(struc) > 1) {
				is_ref = 1;
			}
			struc = Z_REFVAL_P(struc);
			goto again;
			break;
		default:
			php_printf("%sUNKNOWN:0\n", COMMON);
			break;
	}
}
```

```cpp
case IS_STRING:
	php_printf("%sstring(%zd) \"", COMMON, Z_STRLEN_P(struc));
	PHPWRITE(Z_STRVAL_P(struc), Z_STRLEN_P(struc));
	PUTS("\"\n");
	break;
```

![](2021-12-07%2019-42-59%20的屏幕截图.png)

在`main/output.c`中

```cpp
PHPAPI size_t php_output_write(const char *str, size_t len)
{
	if (OG(flags) & PHP_OUTPUT_ACTIVATED) {
		php_output_op(PHP_OUTPUT_HANDLER_WRITE, str, len);
		return len;
	}
	if (OG(flags) & PHP_OUTPUT_DISABLED) {
		return 0;
	}
	return php_output_direct(str, len);
}

...

static inline void php_output_op(int op, const char *str, size_t len)
{
	php_output_context context;
	php_output_handler **active;
	int obh_cnt;

	if (php_output_lock_error(op)) {
		return;
	}

	php_output_context_init(&context, op);

	/*
	 * broken up for better performance:
	 *  - apply op to the one active handler; note that OG(active) might be popped off the stack on a flush
	 *  - or apply op to the handler stack
	 */
	if (OG(active) && (obh_cnt = zend_stack_count(&OG(handlers)))) {
		context.in.data = (char *) str;
		context.in.used = len;

		if (obh_cnt > 1) {
			zend_stack_apply_with_argument(&OG(handlers), ZEND_STACK_APPLY_TOPDOWN, php_output_stack_apply_op, &context);
		} else if ((active = zend_stack_top(&OG(handlers))) && (!((*active)->flags & PHP_OUTPUT_HANDLER_DISABLED))) {
			php_output_handler_op(*active, &context);
		} else {
			php_output_context_pass(&context);
		}
	} else {
		context.out.data = (char *) str;
		context.out.used = len;
	}

	if (context.out.data && context.out.used) {
		php_output_header();

		if (!(OG(flags) & PHP_OUTPUT_DISABLED)) {
#if PHP_OUTPUT_DEBUG
			fprintf(stderr, "::: sapi_write('%s', %zu)\n", context.out.data, context.out.used);
#endif
			sapi_module.ub_write(context.out.data, context.out.used);

			if (OG(flags) & PHP_OUTPUT_IMPLICITFLUSH) {
				sapi_flush();
			}

			OG(flags) |= PHP_OUTPUT_SENT;
		}
	}
	php_output_context_dtor(&context);
}
```

在`sapi/cli/php_cli.c`中

```cpp
static size_t sapi_cli_ub_write(const char *str, size_t str_length) /* {{{ */
{
	const char *ptr = str;
	size_t remaining = str_length;
	ssize_t ret;

	if (!str_length) {
		return 0;
	}

	if (cli_shell_callbacks.cli_shell_ub_write) {
		size_t ub_wrote;
		ub_wrote = cli_shell_callbacks.cli_shell_ub_write(str, str_length);
		if (ub_wrote != (size_t) -1) {
			return ub_wrote;
		}
	}

	while (remaining > 0)
	{
		ret = sapi_cli_single_write(ptr, remaining);
		if (ret < 0) {
#ifndef PHP_CLI_WIN32_NO_CONSOLE
			EG(exit_status) = 255;
			php_handle_aborted_connection();
#endif
			break;
		}
		ptr += ret;
		remaining -= ret;
	}

	return (ptr - str);
}

...

PHP_CLI_API ssize_t sapi_cli_single_write(const char *str, size_t str_length) /* {{{ */
{
	ssize_t ret;

	if (cli_shell_callbacks.cli_shell_write) {
		cli_shell_callbacks.cli_shell_write(str, str_length);
	}

#ifdef PHP_WRITE_STDOUT
	do {
		ret = write(STDOUT_FILENO, str, str_length);
	} while (ret <= 0 && errno == EAGAIN && sapi_cli_select(STDOUT_FILENO));
#else
	ret = fwrite(str, 1, MIN(str_length, 16384), stdout);
	if (ret == 0 && ferror(stdout)) {
		return -1;
	}
#endif
	return ret;
}
```

`ret = write(STDOUT_FILENO, str, str_length);`使用write函数进行输出

![](2021-12-07%2020-17-31%20的屏幕截图.png)

---

`lex`词法分析`Zend/zend_language_scanner.l`,`yacc`语法分析`Zend/zend_language_parser.y`

对文件`index.php`进行调试

```php
<?php
echo "Hello";
?>
```

在`Zend/zend_language_scanner.l`中

```
<ST_IN_SCRIPTING>"echo" {
	RETURN_TOKEN(T_ECHO);
}
```

在`Zend/zend_language_parser.y`中

```
	|	T_ECHO echo_expr_list ';'		{ $$ = $2; }

...

echo_expr_list:
		echo_expr_list ',' echo_expr { $$ = zend_ast_list_add($1, $3); }
	|	echo_expr { $$ = zend_ast_create_list(1, ZEND_AST_STMT_LIST, $1); }
;
echo_expr:
	expr { $$ = zend_ast_create(ZEND_AST_ECHO, $1); }
;
```

![](2021-12-07%2016-10-46%20的屏幕截图.png)

```cpp
void zend_compile_echo(zend_ast *ast) /* {{{ */
{
	zend_op *opline;
	zend_ast *expr_ast = ast->child[0];

	znode expr_node;
	zend_compile_expr(&expr_node, expr_ast);

	opline = zend_emit_op(NULL, ZEND_ECHO, &expr_node, NULL);
	opline->extended_value = 0;
}
```

`opcode`为`ZEND_ECHO`

在`Zend/zend_vm_opcodes.h`中

```cpp
#define ZEND_ECHO                       136
```

在`Zend/zend_vm_def.h`中找到`ZEND_ECHO`的定义

```cpp
ZEND_VM_HANDLER(136, ZEND_ECHO, CONST|TMPVAR|CV, ANY)
{
	USE_OPLINE
	zend_free_op free_op1;
	zval *z;

	SAVE_OPLINE();
	z = GET_OP1_ZVAL_PTR_UNDEF(BP_VAR_R);

	if (Z_TYPE_P(z) == IS_STRING) {
		zend_string *str = Z_STR_P(z);

		if (ZSTR_LEN(str) != 0) {
			zend_write(ZSTR_VAL(str), ZSTR_LEN(str));
		}
	} else {
		zend_string *str = zval_get_string_func(z);

		if (ZSTR_LEN(str) != 0) {
			zend_write(ZSTR_VAL(str), ZSTR_LEN(str));
		} else if (OP1_TYPE == IS_CV && UNEXPECTED(Z_TYPE_P(z) == IS_UNDEF)) {
			ZVAL_UNDEFINED_OP1();
		}
		zend_string_release_ex(str, 0);
	}

	FREE_OP1();
	ZEND_VM_NEXT_OPCODE_CHECK_EXCEPTION();
}
```

如果是字符串的话就直接进行`zend_write`操作,如果不是就先进行`zval_get_string_func`操作

在`Zend/zend_operators.c`中

```cpp
ZEND_API zend_string* ZEND_FASTCALL zval_get_string_func(zval *op) /* {{{ */
{
	return __zval_get_string_func(op, 0);
}


static zend_always_inline zend_string* __zval_get_string_func(zval *op, zend_bool try) /* {{{ */
{
try_again:
	switch (Z_TYPE_P(op)) {
		case IS_UNDEF:
		case IS_NULL:
		case IS_FALSE:
			return ZSTR_EMPTY_ALLOC();
		case IS_TRUE:
			return ZSTR_CHAR('1');
		case IS_RESOURCE: {
			return zend_strpprintf(0, "Resource id #" ZEND_LONG_FMT, (zend_long)Z_RES_HANDLE_P(op));
		}
		case IS_LONG: {
			return zend_long_to_str(Z_LVAL_P(op));
		}
		case IS_DOUBLE: {
			return zend_strpprintf(0, "%.*G", (int) EG(precision), Z_DVAL_P(op));
		}
		case IS_ARRAY:
			zend_error(E_NOTICE, "Array to string conversion");
			return (try && UNEXPECTED(EG(exception))) ?
				NULL : ZSTR_KNOWN(ZEND_STR_ARRAY_CAPITALIZED);
		case IS_OBJECT: {
			zval tmp;
			if (Z_OBJ_HT_P(op)->cast_object) {
				if (Z_OBJ_HT_P(op)->cast_object(op, &tmp, IS_STRING) == SUCCESS) {
					return Z_STR(tmp);
				}
			} else if (Z_OBJ_HT_P(op)->get) {
				zval *z = Z_OBJ_HT_P(op)->get(op, &tmp);
				if (Z_TYPE_P(z) != IS_OBJECT) {
					zend_string *str = try ? zval_try_get_string(z) : zval_get_string(z);
					zval_ptr_dtor(z);
					return str;
				}
				zval_ptr_dtor(z);
			}
			if (!EG(exception)) {
				zend_throw_error(NULL, "Object of class %s could not be converted to string", ZSTR_VAL(Z_OBJCE_P(op)->name));
			}
			return try ? NULL : ZSTR_EMPTY_ALLOC();
		}
		case IS_REFERENCE:
			op = Z_REFVAL_P(op);
			goto try_again;
		case IS_STRING:
			return zend_string_copy(Z_STR_P(op));
		EMPTY_SWITCH_DEFAULT_CASE()
	}
	return NULL;
}
```

在`Zend/zend.h`中

```cpp
extern ZEND_API zend_write_func_t zend_write;

typedef int (*zend_write_func_t)(const char *str, size_t str_length);
```

在`Zend/zend.c`中

```cpp
zend_write = (zend_write_func_t) utility_functions->write_function;
```

在`main/main.c`中

```cpp
zuf.write_function = php_output_write;
```

在`main/output.c`中

```cpp
PHPAPI size_t php_output_write(const char *str, size_t len)
{
	if (OG(flags) & PHP_OUTPUT_ACTIVATED) {
		php_output_op(PHP_OUTPUT_HANDLER_WRITE, str, len);
		return len;
	}
	if (OG(flags) & PHP_OUTPUT_DISABLED) {
		return 0;
	}
	return php_output_direct(str, len);
}

...

static inline void php_output_op(int op, const char *str, size_t len)
{
	php_output_context context;
	php_output_handler **active;
	int obh_cnt;

	if (php_output_lock_error(op)) {
		return;
	}

	php_output_context_init(&context, op);

	/*
	 * broken up for better performance:
	 *  - apply op to the one active handler; note that OG(active) might be popped off the stack on a flush
	 *  - or apply op to the handler stack
	 */
	if (OG(active) && (obh_cnt = zend_stack_count(&OG(handlers)))) {
		context.in.data = (char *) str;
		context.in.used = len;

		if (obh_cnt > 1) {
			zend_stack_apply_with_argument(&OG(handlers), ZEND_STACK_APPLY_TOPDOWN, php_output_stack_apply_op, &context);
		} else if ((active = zend_stack_top(&OG(handlers))) && (!((*active)->flags & PHP_OUTPUT_HANDLER_DISABLED))) {
			php_output_handler_op(*active, &context);
		} else {
			php_output_context_pass(&context);
		}
	} else {
		context.out.data = (char *) str;
		context.out.used = len;
	}

	if (context.out.data && context.out.used) {
		php_output_header();

		if (!(OG(flags) & PHP_OUTPUT_DISABLED)) {
#if PHP_OUTPUT_DEBUG
			fprintf(stderr, "::: sapi_write('%s', %zu)\n", context.out.data, context.out.used);
#endif
			sapi_module.ub_write(context.out.data, context.out.used);

			if (OG(flags) & PHP_OUTPUT_IMPLICITFLUSH) {
				sapi_flush();
			}

			OG(flags) |= PHP_OUTPUT_SENT;
		}
	}
	php_output_context_dtor(&context);
}
```

在`sapi/cli/php_cli.c`中

```cpp
static size_t sapi_cli_ub_write(const char *str, size_t str_length) /* {{{ */
{
	const char *ptr = str;
	size_t remaining = str_length;
	ssize_t ret;

	if (!str_length) {
		return 0;
	}

	if (cli_shell_callbacks.cli_shell_ub_write) {
		size_t ub_wrote;
		ub_wrote = cli_shell_callbacks.cli_shell_ub_write(str, str_length);
		if (ub_wrote != (size_t) -1) {
			return ub_wrote;
		}
	}

	while (remaining > 0)
	{
		ret = sapi_cli_single_write(ptr, remaining);
		if (ret < 0) {
#ifndef PHP_CLI_WIN32_NO_CONSOLE
			EG(exit_status) = 255;
			php_handle_aborted_connection();
#endif
			break;
		}
		ptr += ret;
		remaining -= ret;
	}

	return (ptr - str);
}

...

PHP_CLI_API ssize_t sapi_cli_single_write(const char *str, size_t str_length) /* {{{ */
{
	ssize_t ret;

	if (cli_shell_callbacks.cli_shell_write) {
		cli_shell_callbacks.cli_shell_write(str, str_length);
	}

#ifdef PHP_WRITE_STDOUT
	do {
		ret = write(STDOUT_FILENO, str, str_length);
	} while (ret <= 0 && errno == EAGAIN && sapi_cli_select(STDOUT_FILENO));
#else
	ret = fwrite(str, 1, MIN(str_length, 16384), stdout);
	if (ret == 0 && ferror(stdout)) {
		return -1;
	}
#endif
	return ret;
}
```

`ret = write(STDOUT_FILENO, str, str_length);`使用write函数进行输出

![](2021-12-07%2020-10-40%20的屏幕截图.png)