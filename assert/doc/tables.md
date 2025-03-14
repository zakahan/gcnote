#  MySQL表



### 用户登录表 User

| 表项         | 数据类型       | 作用              |
|------------|------------|-----------------|
| **UserId** | **String** | **用户ID，唯一标识用户** |
| UserName   | String     | 用户名（不可重复）       |
| Password   | String     | 密码              |
| Email      | String     | 邮箱，注册的时候用       |



### 用户-知识库关联表 Index

| 表项        | 数据类型   | 作用                                |
| ----------- | ---------- | ----------------------------------- |
| UserId      | String     | 用户ID，唯一标识用户                |
| **IndexId** | **String** | **知识库ID，唯一标识知识库**        |
| IndexName   | String     | 知识库名称// 重命名的话只需要改这个 |



### 知识库-文件关联表 KBFile

| 表项           | 数据类型       | 作用            |
|--------------|------------|---------------|
| IndexId      | String     | 知识库ID，唯一标识知识库 |
| **KBFileId** | **String** | **唯一标识文件**    |
| KBFileName   | String     | 文件名           |

知识库文档的真实路径会是/\$DocumentDir/knowledge_base/\$IndexId/\$KBFileId/KBFileName.md

// 因为存在图片资源，所以必须单独建文件夹



### 回收站-文件关联表 Recycle

| 表项          | 数据类型   | 作用                                       |
| ------------- | ---------- | ------------------------------------------ |
| UserId        | String     | 标识用户                                   |
| SourceIndexId | String     | 原所属知识库ID，唯一标识知识库（恢复需要） |
| **KBFileId**  | **String** | **唯一标识文件**                           |
| KBFileName    | String     | 文件名                                     |

知识库文档的真实路径会是/\$DocumentDir/recycle_bin/\$IndexId/\$KBFileId/KBFileName.md


# ElasticSearch表

就一个表把

### 知识库切片表

| 表项                | 数据类型     | 作用                       |
| ------------------- | ------------ | -------------------------- |
| page_content        | text         | 文本内容                   |
| vector              | dense_vector | 向量表                     |
| **metadata.doc_id** | **keyword**  | **文件id**                 |
| metadata.kb_file_id | keyword      | index名称                  |
| metadata.index_id   | text         | 文件名称，用于支持模糊搜索 |



### Mysql设置

```shell
docker exec -it c2851afc9957 /bin/bash

mysql -uroot -p

Password:mysql-zxcvbnm

mysql> CREATE DATABASE gcnote CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
Query OK, 1 row affected (0.06 sec)

mysql> SHOW DATABASES;
+--------------------+
| Database           |
+--------------------+
| gcnote             |
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
5 rows in set (0.17 sec)

# --------------------------------

mysql> USE gcnote
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql>
mysql> SHOW TABLES;
+------------------+
| Tables_in_gcnote |
+------------------+
| index            |
| kb_file          |
| user             |
+------------------+
3 rows in set (0.01 sec)

mysql>
```