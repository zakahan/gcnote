# gcnote

仿照语雀做了一个go的项目

> go get -u github.com/swaggo/swag/cmd/swag

## 代替

>  ioutil.ReadAll -> io.ReadAll
>  ioutil.ReadFile -> os.ReadFile
>  ioutil.ReadDir -> os.ReadDir
>  // others
>  ioutil.NopCloser -> io.NopCloser
>  ioutil.ReadDir -> os.ReadDir
>  ioutil.TempDir -> os.MkdirTemp
>  ioutil.TempFile -> os.CreateTemp
>  ioutil.WriteFile -> os.WriteFile



### go-swagger配置

```shell
go install github.com/swaggo/swag/cmd/swag@latest

# 查看
swag -v

```



其他的安装命令

```shell
go get -u github.com/swaggo/gin-swagger   
go get -u github.com/swaggo/files
# 模版
go get -u github.com/alecthomas/template

```



然后是swagger init

```shell
swag init
```





## 参考资料

1. [gin-demo](https://github.com/ngyhd/gin-demo)
2. [JWT介绍](https://blog.csdn.net/weixin_42030357/article/details/95629924)