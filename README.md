# gcnote

仿照语雀做了一个go的项目

符号

## 目前进展 -

2024 / 12 / 17

1. 后端简单curd的基本上完成了，判断文件存在性之类的好像也不是很必要，搁置争议
2. read_file之类的和前端技术选型有关，我暂时不好判断怎么写
3. 删除用户应该级联删除Index表和KBFile表，但是这个就先这样了，就当（保留用户数据）
4. 接下来要推进ES方面的处理了。（晚上做or明天做）
5. 现在可以说，我完成了差不多1/4的工作量了。剩下的里面，1/4是ES搜索+RAG功能，1/4是前端基础操作，1/4 共享文档操作

### 用户接口

| 进展 | 接口说明     | 接口地址       | 访问方式 |
| :--: | ------------ | -------------- | -------- |
|  ✅  | 注册用户     | /user/register | POST     |
|  ✅  | 用户登录     | /user/login    | POST     |
|  ✅  | 更新用户密码 | /user/update   | POST     |
|  ✅  | 删除用户     | /user/delete   | POST     |
|  ✅  | 展示用户信息 | /user/info     | GET      |

### 知识库操作处理

| 进展 | 接口说明       | 接口地址            | 访问方式 |
| :--: | -------------- | ------------------- | -------- |
|  ✅  | 创建知识库     | /index/create_index | POS      |
|  ✅  | 删除知识库     | /index/delete_index | POST     |
| ⬜️ | 知识库是否存在 | /index/search_index | GET      |
|  ✅  | 展示知识库清单 | /index/show_indexs  | GET      |
|  ✅  | 重命名知识库   | /index/rename_index | POST     |

### 文档操作

| 进展 | 接口说明           | 接口地址             | 访问方式 |
| :--: | ------------------ | -------------------- | -------- |
|  ✅  | 新建文档           | /index/create_file   | POST     |
| ⬜️ | 读取文档           | /index/read_file     | POST     |
| ⬜️ | 更新文档           | /index/update_file   | POST     |
|  ✅  | 导入文档           | /index/add_file      | POST     |
|  ✅  | 重命名文档         | /index/rename_file   | POST     |
|  ✅  | 删除文档到回收站   | /index/recycle_files | POST     |
|  ✅  | 搜索文档(按文件名) | /index/search_file   | POST     |
|  ✅  | 展示知识库文档列表 | /index/show_files    | GET      |
| ⬜️ | 文档存在           | /index/file_exist    | GET      |
| ⬜️ | 文档共享           | /index/share_files   | POST     |

- 新建文档，每次都是新建一个空的存进去先，这样后续所有操作都变成了更新，而不是新建 + 更新
- 导入文档，这个就是新建 + 更新的复合，但是这个功能常用，我直接分出来吧
- 文档共享：这个是咋实现的？好像得加个标记，但文档共享不行，还得有知识库共享，难道要先新建一个知识库，然后共享吗？
- 还是说对某个文档，建立一个表格，里面是允许操作的人员名单？

### 回收站操作

回收站就是一个单独的知识库，这里每次删除操作都是真的彻底的删除，我应该新建一个表的

| 进展 | 接口说明       | 接口地址               | 访问方式 |
| :--: | -------------- | ---------------------- | -------- |
|  ✅  | 展示回收站内容 | /recycle/show_files    | GET      |
|  ✅  | 彻底删除文档   | /recycle/delete_files  | POST     |
|  ✅  | 恢复文档       | /recycle/restore_files | POST     |
|  ✅  | 定期清理       | /recycle/clearup       | POST     |
|  ✅  | 清空回收站     | /recycle/clear         | POST     |

回收站单独一张表。设置了一个接口clearup，清理更新时间大于30的文件，这个接口倒是每天24点被前端调用一次。

### 搜索与RAG

将搜索设置为基于知识库这一层次的搜索问答，而不是单个文件，我都找到文件了？我还用你搜？

5. 搜索 - 需要ElasticSearch支持
   1. 大类搜索，按照文件名做搜索
   2. 切片搜索，搜索切片里的内容
6. AI功能
   1. 单文件问答
   2. RAG
   3. 生成日程表或者是待办事宜表
7. 共享，应该是MySQL里面对文件有个标记，是私有的还是都能编辑，就是一个编辑权限表

### 需要的东西

1. Redis缓存
2. MySQL - 存储一些表，我暂时没想好存什么，大概就是知识库与其ID，知识库内的文件与其切片索引
3. ElasticSearch - 每个知识库一个index，里面存切片，然后还要有文件更新的功能。

## 参考资料

1. [gin-demo](https://github.com/ngyhd/gin-demo)
   参考了gin-demo的基本框架以及redis和自定义状态码设计。
2. [JWT介绍](https://blog.csdn.net/weixin_42030357/article/details/95629924)
   JWT介绍部分参考了此处
