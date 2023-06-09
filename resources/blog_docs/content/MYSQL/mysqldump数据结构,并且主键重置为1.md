---
categories:
  - MYSQL
---
- mysqldump
```shell
mysqldump -u root -pxxx123456 -d testa  
```
- 这个是导出数据结构  

```sql
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `customerpath` char(20) DEFAULT NULL,
  `username` varchar(50) DEFAULT NULL,
  `contactid` int(10) unsigned DEFAULT NULL,
  `unique_openid` varchar(50) DEFAULT NULL,
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `unique_openid` (`unique_openid`),
  KEY `contactid` (`contactid`),
  KEY `customerpath` (`customerpath`)
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=latin1;  
```

- 发现虽然导出了，但是主键是从74条记录开始算起所以需要用sed
```shell
sed -i 's/\(AUTO_INCREMENT\)=[0-9]*/\1=1/g'
```

- 在mysqldump的时候不需要-i所以最终命令是:   
```shell
mysqldump -u root -pxxx123456 -d testa | sed 's/\(AUTO_INCREMENT\)=[0-9]*/\1=1/g' > a1.sql
```


- 参考 
1. https://blog.csdn.net/weixin_34096182/article/details/92083589
2. https://www.cnblogs.com/pengmengnan/p/9041310.html loverable 替换
3. https://segmentfault.com/q/1010000000408080
