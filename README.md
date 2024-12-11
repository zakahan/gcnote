# gcnote

仿照语雀做了一个go的项目

符号 



## 目前进展

### 用户接口

| 进展 | 接口说明     | 接口地址       | 访问方式 |
| :------: | ------------ | -------------- | -------- |
| :white_check_mark: | 注册用户     | /user/register | POST     |
| :white_check_mark: | 用户登录     | /user/login    | POST     |
| :white_check_mark: | 更新用户密码 | /user/update   | POST     |
| :white_check_mark: | 删除用户     | /user/delete   | POST     |
| :white_check_mark: | 展示用户信息 | /user/info     | GET      |



### 知识库操作处理



|         进展         | 接口说明       | 接口地址                  | 访问方式 |
| :------------------: | -------------- | ------------------------- | -------- |
| :white_check_mark: | 创建知识库     | /knowledge_base/create_kb | POST     |
| :white_large_square: | 删除知识库     | /knowledge_base/delete_kb | POST     |
| :white_large_square: | 知识库是否存在 | /knowledge_base/search_kb | GET      |
| :white_large_square: | 展示知识库清单 | /knowledge_base/show_kbs  | GET      |
| :white_large_square: | 重命名知识库   | /knowledge_base/rename_kb | POST     |



### 文档操作



|         进展         | 接口说明           | 接口地址                     | 访问方式 |
| :------------------: | ------------------ | ---------------------------- | -------- |
| :white_large_square: | 新建文档           | /knowledge_base/new_file     | POST     |
| :white_large_square: | 读取文档           | /knowledge_base/read_file    | POST     |
| :white_large_square: | 更新文档           | /knowledge_base/update_file  | POST     |
| :white_large_square: | 导入文档           | /knowledge_base/import_files | POST     |
| :white_large_square: | 重命名文档         | /knowledge_base/rename_file  | POST     |
| :white_large_square: | 删除文档到回收站   | /knowledge_base/delete_files | POST     |
| :white_large_square: | 搜索文档(按文件名) | /knowledge_base/search_file  | POST     |
| :white_large_square: | 展示知识库文档列表 | /knowledge_base/show_files   | GET      |
| :white_large_square: | 文档存在           | /knowledge_base/file_exist   | GET      |
| :white_large_square: | 文档共享           | /knowledge_base/share_files  | POST     |

- 新建文档，每次都是新建一个空的存进去先，这样后续所有操作都变成了更新，而不是新建 + 更新

- 导入文档，这个就是新建 + 更新的复合，但是这个功能常用，我直接分出来吧
- 文档共享：这个是咋实现的？好像得加个标记，但文档共享不行，还得有知识库共享，难道要先新建一个知识库，然后共享吗？
- 还是说对某个文档，建立一个表格，里面是允许操作的人员名单？



### 回收站操作

回收站就是一个单独的知识库，这里每次删除操作都是真的彻底的删除，我应该新建一个表的

|         进展         | 接口说明       | 接口地址                   | 访问方式 |
| :------------------: | -------------- | -------------------------- | -------- |
| :white_large_square: | 展示回收站内容 | /recycle_bin/show_files    | GET      |
| :white_large_square: | 彻底删除文档   | /recycle_bin/delete_files  | POST     |
| :white_large_square: | 恢复文档       | /recycle_bin/recycle_files | POST     |
| :white_large_square: | 定期清理       | 这功能咋实现？             |          |

回收站的话，好像只需要给kb_file加个标记-“是否位于回收站”，每次查询回收站的时候，呃，但是这样的话好像需要级联查询，有点费时间？我得研究一下。



### 搜索与RAG




5. 搜索 - 需要ElasticSearch支持
    1. 大类搜索，按照文件名做搜索
    2. 切片搜索，搜索切片里的内容
6. AI功能
    1. 单文件问答
    2. RAG
    3. 生成日程表或者是待办事宜表
11. 共享，应该是MySQL里面对文件有个标记，是私有的还是都能编辑，就是一个编辑权限表

### 需要的东西

1. Redis缓存
2. MySQL - 存储一些表，我暂时没想好存什么，大概就是知识库与其ID，知识库内的文件与其切片索引
3. ElasticSearch - 每个知识库一个index，里面存切片，然后还要有文件更新的功能。



## 参考资料

1. [gin-demo](https://github.com/ngyhd/gin-demo)
2. [JWT介绍](https://blog.csdn.net/weixin_42030357/article/details/95629924)