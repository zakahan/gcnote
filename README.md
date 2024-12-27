# gcnote

仿照语雀做了一个go的项目

符号

## 目前进展 -

2024 / 12 / 20

最新指导思想，放弃太多的广度，尽可能深耕某几个地方。

参加完字节讲座之后我意识到，不需要把自己的项目做的太完善，而是要关注在某几个点，自己深耕于其中。你会一堆东西，但这也只是掉包侠，AI不比你会用吗？

与其这样，不如自己深耕某几个领域，目前我选定的就是文档解析和ElasticSearch，重点关注在ES的特性以及如何导入文件这两个上，PDF2MD能力的支撑，以及DOCX2MD能力的提升，这才是能体现出我的不同的地方，而且工作量不会太高，对于系统复杂程度的要求不会太高，不想这种go-web项目一样需要消耗大量的脑力处理这种多重依赖、边界条件判断的问题。



1. 目前的宗旨：尽快完结，尽快推进到前端部分。晚上把检查存在写完，验证一下接口的bug，就可以暂时结束后端部分的开发了。读取文件之类的和前端关系比较大的后续做。
2. 还有个问题就是怎么渲染那些图片，这个事情太麻烦了，涉及到图片路径的问题。
3. 把PDF2MD提上日程，把DOCX2MD的标题-> #Heading更新提上日程
4. 现在可以说，我完成了差不多1.5/4的工作量了。剩下的里面，0.5/4是ES搜索+RAG功能，1/4是前端基础操作，1/4 共享文档操作

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
| :white_check_mark: | 知识库是否存在 | /index/search_index | GET      |
|  ✅  | 展示知识库清单 | /index/show_indexs  | GET      |
|  ✅  | 重命名知识库   | /index/rename_index | POST     |

### 文档操作

| 进展 | 接口说明           | 接口地址             | 访问方式 |
| :--: | ------------------ | -------------------- | -------- |
|  ✅  | 新建文档           | /index/create_file   | POST     |
| :white_check_mark: | 读取文档           | /index/read_file     | POST     |
| ⬜️ | 更新文档           | /index/update_file   | POST     |
|  ✅  | 导入文档           | /index/add_file      | POST     |
|  ✅  | 重命名文档         | /index/rename_file   | POST     |
|  ✅  | 删除文档到回收站   | /index/recycle_files | POST     |
|  ✅  | 搜索文档(按文件名) | /index/search_file   | POST     |
|  ✅  | 展示知识库文档列表 | /index/show_files    | GET      |
| :white_check_mark: | 文档存在           | /index/file_exist    | GET      |
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

|         进展         | 接口说明       | 接口地址         | 访问方式 |
| :------------------: | -------------- | ---------------- | -------- |
|          ✅           | 在知识库中搜索 | /index/retrieval | POST     |
| :white_large_square: | 检索生成       | /index/naive_rag | POST     |
| :white_large_square: | 直接AI问答     | /index/chat_llm  | POST     |

- 检索和AI部分，额，后续做。做不完就算了

将搜索设置为基于知识库这一层次的搜索问答，而不是单个文件，我都找到文件了？我还用你搜？

5. 搜索 - 需要ElasticSearch支持
   2. 切片搜索，搜索切片里的内容
6. AI功能
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
   
3. [Dify](https://github.com/langgenius/dify)

   参考了其中RAG的部分
