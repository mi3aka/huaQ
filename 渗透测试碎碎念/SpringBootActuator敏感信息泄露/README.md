[Spring Boot Actuator敏感信息泄露](http://r3start.net/index.php/2019/01/20/377)

[Spring Boot Actuator未授权内存分析方法](https://www.freebuf.com/vuls/339511.html)

[Springboot未授权+shrio组件getshell](https://www.freebuf.com/articles/web/288423.html)

---

```
apt install maven
git clone https://github.com/veracode-research/actuator-testbed
cd actuator-testbed
```

修改pom.xml

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>org.springframework</groupId>
    <artifactId>actuator-testbed</artifactId>
    <version>0.1.0</version>

    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <!--<version>2.0.5.RELEASE</version>-->
        <!--<version>1.5.16.RELEASE</version>-->
        <version>1.4.7.RELEASE</version>
        <!--<version>1.3.8.RELEASE</version>-->
        <!--<version>1.2.8.RELEASE</version>-->
    </parent>

    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-actuator</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.cloud</groupId>
            <artifactId>spring-cloud-starter-netflix-eureka-client</artifactId>
            <version>1.4.0.RELEASE</version>
        </dependency>
        <dependency>
            <groupId>org.jolokia</groupId>
            <artifactId>jolokia-core</artifactId>
            <version>1.6.0</version>
        </dependency>
        <dependency>
            <groupId>javax.xml.bind</groupId>
            <artifactId>jaxb-api</artifactId>
            <version>2.3.0</version>
        </dependency>
        <dependency>
            <groupId>com.sun.xml.bind</groupId>
            <artifactId>jaxb-impl</artifactId>
            <version>2.3.0</version>
        </dependency>
        <dependency>
            <groupId>com.sun.xml.bind</groupId>
            <artifactId>jaxb-core</artifactId>
            <version>2.3.0</version>
        </dependency>
        <dependency>
            <groupId>javax.activation</groupId>
            <artifactId>activation</artifactId>
            <version>1.1.1</version>
        </dependency>
    </dependencies>



    <dependencyManagement>
        <dependencies>
            <dependency>
                <groupId>org.springframework.cloud</groupId>
                <artifactId>spring-cloud-dependencies</artifactId>
                <version>Camden.RELEASE</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>

    <properties>
        <java.version>1.8</java.version>
    </properties>

    <build>
        <plugins>
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
            </plugin>
        </plugins>
    </build>
</project>
```

修改`actuator-testbed/src/main/resources/application.properties`

```conf
server.port=8090
server.address=127.0.0.1

# vulnerable configuration set 0: spring boot 1.0 - 1.4
# all spring boot versions 1.0 - 1.4 expose actuators by default without any parameters
# no configuration required to expose them

# safe configuration set 0: spring boot 1.0 - 1.4
#management.security.enabled=true

# vulnerable configuration set 1: spring boot 1.5+
# spring boot 1.5+ requires management.security.enabled=false to expose sensitive actuators
#management.security.enabled=false

# safe configuration set 1: spring boot 1.5+
# when 'management.security.enabled=false' but all sensitive actuators explicitly disabled
management.security.enabled=false

# vulnerable configuration set 2: spring boot 2+
#management.endpoints.web.exposure.include=*
spring.datasource.url=jdbc:mysql://127.0.0.1:3306/root
spring.datasource.username=root
spring.datasource.password=123456
spring.datasource.driver-class-name=com.mysql.jdbc.Driver
```

`mvn spring-boot:run`并将8090端口映射

访问[http://az.jscdn.eu.org:31433](http://az.jscdn.eu.org:31433)

![](https://img.mi3aka.eu.org/2022/09/26a8cdef715ac1ab461b9da31478aafd.png)

访问[http://az.jscdn.eu.org:31433/env](http://az.jscdn.eu.org:31433/env)

![](https://img.mi3aka.eu.org/2022/09/2413a5229a78e4f7dc58b2b5d9bbdcf1.png)

![](https://img.mi3aka.eu.org/2022/09/1770de14ceb58ec93eec9655d7d3a488.png)

访问[http://az.jscdn.eu.org:31433/heapdump](http://az.jscdn.eu.org:31433/heapdump)将数据导出

![](https://img.mi3aka.eu.org/2022/09/1a6fbccf45d2175f7db3db85f13666c8.png)

---

利用[Eclipse Memory Analyzer](https://www.eclipse.org/mat/downloads.php)对headdump进行分析

![](https://img.mi3aka.eu.org/2022/09/53f10c4ec320ce7b66e4c0bf7ea19ccc.png)

>webkit错误

![](https://img.mi3aka.eu.org/2022/09/cf3793be37c4818464b03d1b9ef8bec9.png)

```
sudo apt install libwebkit2gtk-4.0-37
sudo apt install webkit2gtk-driver
```

![](https://img.mi3aka.eu.org/2022/09/0ff9e24ec8227006d8b86f788b28a4e6.png)

![](https://img.mi3aka.eu.org/2022/09/2409bdc2857007c378eab0e8bcc2ebb5.png)

---

