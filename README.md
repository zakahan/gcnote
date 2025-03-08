# gcnote

仿照语雀做了一个go的项目

### 后续计划
- 尝试加入共享文档功能
- 尝试加入RAG部分



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
| ✅ | 知识库是否存在 | /index/search_index | GET      |
|  ✅  | 展示知识库清单 | /index/show_indexs  | GET      |
|  ✅  | 重命名知识库   | /index/rename_index | POST     |

### 文档操作

| 进展 | 接口说明           | 接口地址             | 访问方式 |
| :--: | ------------------ | -------------------- | -------- |
|  ✅  | 新建文档           | /index/create_file   | POST     |
| ✅ | 读取文档           | /index/read_file     | POST     |
| ✅ | 更新文档           | /index/update_file   | POST     |
|  ✅  | 导入文档           | /index/add_file      | POST     |
|  ✅  | 重命名文档         | /index/rename_file   | POST     |
|  ✅  | 删除文档到回收站   | /index/recycle_files | POST     |
|  ✅  | 搜索文档(按文件名) | /index/search_file   | POST     |
|  ✅  | 展示知识库文档列表 | /index/show_files    | GET      |
| ✅ | 文档存在           | /index/file_exist    | GET      |


### 回收站操作

回收站就是一个单独的知识库，这里每次删除操作都是真的彻底的删除

| 进展 | 接口说明       | 接口地址               | 访问方式 |
| :--: | -------------- | ---------------------- | -------- |
|  ✅  | 展示回收站内容 | /recycle/show_files    | GET      |
|  ✅  | 彻底删除文档   | /recycle/delete_files  | POST     |
|  ✅  | 恢复文档       | /recycle/restore_files | POST     |
|  ✅  | 定期清理       | /recycle/clearup       | POST     |
|  ✅  | 清空回收站     | /recycle/clear         | POST     |


5. 搜索 - 需要ElasticSearch支持
   2. 切片搜索，搜索切片里的内容
6. AI功能(暂未完成)
   2. RAG
   3. 生成日程表或者是待办事宜表

### 涉及组件

1. Redis缓存
2. MySQL - 存储一些表，我暂时没想好存什么，大概就是知识库与其ID，知识库内的文件与其切片索引
3. ElasticSearch - 每个知识库一个index，里面存切片，然后还要有文件更新的功能。

## 参考资料

1. [gin-demo](https://github.com/ngyhd/gin-demo)
   参考了gin-demo的基本框架以及redis和自定义状态码设计。
   
2. [JWT介绍](https://blog.csdn.net/weixin_42030357/article/details/95629924)
   JWT介绍部分参考了此处
   
3. [Dify](https://github.com/langgenius/dify)

   参考了其中RAG的部分
