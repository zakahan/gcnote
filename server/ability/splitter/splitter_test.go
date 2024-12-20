// -------------------------------------------------
// Package splitter
// Author: hanzhi
// Date: 2024/12/20
// -------------------------------------------------

package splitter

import "testing"

func TestSplitter(t *testing.T) {
	text := `
# Attention

### 整体网络架构图

<img src="E:\02-BUPT\02-EByte\2023夏\008-语言情感分析（马祖耀）20230804\U3EwOx.png" alt="U3EwOx" style="zoom:67%;" />

### 分析厘清

从输入开始的故事！！！

从中间开始，实在是太傻比的行为了，为什么不从输入开始，一步步的像数学公式那样导出来

### Input Embedding

首先，输入的是句子，而直接的句子是无法表示信息的，我们必须将其表示为向量，也就是**词汇表里的每一个单词，都要有一个对应的向量来表示**，这种方式同时还涉及了一个降维提高信息密度，以及提高联系的过程。
$$
V_{input}=InputEmbedding(Seq_v)
$$

#### Positional Embedding

这是个什么呢，简单的来说就是把位置信息加到向量$V_{input}$当中。这里他是直接采用的加法，就是两个向量相加，很奇怪，不过无所谓了
$$
V = V_{input} + Pos
$$

------------------------


## Encoder

### Self-Attention

举例的时候就以机器翻译为例子，这样的话有具体的代入感比较方便

首先token的概念，token在机器翻译中是“单词” - **文本的最小单位**

每个token都可以被表示为三个向量，qkv。

假设我当前的句子S，这里面有n个token，$S=[s^1,s^2...,s^n]$、

每个$s^i$对应一组 qkv，这里右上角并不是幂次方的意思，就是编号而已
$$
s^i ->  [q^i,k^i,v^i]
$$

但是！

这特么是个怎么样的对应环节？对应这个词，一点也不数学

实际上，是通过矩阵变换得到的

首先说明一下，小写的qkv是向量，大写的QKV代表矩阵，也就是好几个向量
$$
\begin{cases}
	q^i=s^iW^q\\
	k^i=s^iW^k\\
	v^i=s^iW^v
\\
\end{cases}
$$

这样就得到了这样一组向量。

其中，s的维度是dl，q和k都是dk维的，v是dv维的，（行向量） (1,维度)

所以说也就能推理出三个W权重矩阵的大小分别是因此也就推理出分别为(dl,dk)(dl,dk)(dl,dv)三种格式

先不要去管他们的意义，意义之类的事情，按下不表。

------

之前得到了QKV，假设

Q由n个行向量组成，$Q=[q^1, q^2, ...,q^n]^T$。$(n,dk)$

K由m个行向量组成，$K=[k^1,k^2,...,k^m]^T$。$(m,dk)$

V由m个行向量组成，$V=[v^1,v^2,...,v^m]^T$。$(m,dv)$

K和V的数量一定是相同的，但是Q的数量可以是不同的!!!行向量，有几行，就是有几个，

首先是Attention的公式，这是涉及到了具体的内容，就是一个softmax函数后乘以V向量
$$
Attention\left( Q,K,V \right) =soft\max \left( \frac{QK^T}{\sqrt{d_k}} \right) V
$$
这样来计算

$QK^T$应该是$(n,dk) *  (m,dk)^T$的样式，得到的就是一个$(n,m)$大小的新矩阵，然后经过除以根号 dk以及softmax，仍然是(n,m).

$(n,m) * (m, dv) = (n, dv)$，得到了n个dv维度的向量。

这就是，自注意力，值得注意的是，在 transformer中，用的·并**不是**self Attention，而是multihead Attention，讲self Attention只是相当于打了个基础！！！！！！！

### 稍微说一下意义

QK的内积，表征了q向量和k向量二者之间的相似度，然后根号dk是把这个方便归一化的，之后在进行softmax操作就可以了

----

### MultiHead Attention

注意，这个才是真正用到的东西

而multiHead，多头注意力机制，则是对许多个head进行concat。
$$
MultiHead(Q,K,V)=Concat(head_1,head_2,...,head_h)W^O
\\
where\,\,head_i=Attention(QW_{i}^{Q},KW_{i}^{K},VW_{i}^{V})
$$
公式里出现的新东西比较多，一个一个来吧

首先是$W^Q_i、W^K_i和W^V_i$这三个东西，他们分别表示的是三种投影。head一共有h个，因此这类W也就一共有h组。

1、先对QKV进行投影，2、用Self-Attention机选，3、然后Concat连接，4、乘以一个Wo矩阵，（FIXME：这个矩阵为什么有，以及怎么来的，大小格式，我都一概不知

还是来，细细的推理

$QW^Q_i-(n,d_{model})*(d_{model}，dk)$，注意，这里和之前的不一样了，这里和self-attention不一样。

作者说

> 我们发现，与其使用dmodel维度的键、值和查询执行单一的注意力函数，不如将查询、键和值分别以不同的、学习过的线性投影h次线性投影到dk、dk和dv维度，这是有益的。

所以，self-Attention里说的是直接就确定是dk个，但是如今呢，我们让他是dmodel个，然后再经过一次映射到dmodel个，原来是
$$
\begin{cases}
	Q=sW^Q\\
	K=sW^K\\
	V=sW^V
\\
\end{cases}
\\
where \,\, head=Attention(Q,K,V)
$$
如今的多头注意力机制是这样的。
$$
\begin{cases}
	A_1=QW^Q_i=s^iW^QW^Q_i\\
	A_2=KW^K_i=s^iW^KW^K_i\\
	A_3=VW^K_i=s^iW^VW^K_i
\\
\end{cases}
\\
where \,\, head_i=Attention(A_1,A_2,A_3)
$$
这样的对比就很清晰明了，之所以多头注意力机制，一方面是加入了head这个中间的过程，让数量增加，增加了W个数，然后又增加了一次新的映射W^{QKV}的这些个新的映射，总之就是增加了权重矩阵个数，也就增加了可学习的东西。

关于dkdvdmodel之间的关系，作者用的是
$$
d_k=d_v=d_{model}/h = 64
$$

这就是个参数设置，无所谓了。



最后一步
$$
MultiHead(Q,K,V)=Concat(head_1,head_2,...,head_h)W^O
$$
就形成了多头注意力

已知$W^O\in \mathbb{R}^{hd_v\times d_{model}}$，也就意味着con_head应该是$\mathbb{R}^{n\times hd_v}$的，这里的concat就是直接的同维度的连接啊。。。。所以从n,dv 变成了n,hdv个。（懒得用tex写了）

n还是Query的个数。



<img src="C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20230818111238645.png" alt="image-20230818111238645" style="zoom:67%;" />

注意图上面的Linear，其实就是乘了一个W矩阵，WO也好WQ也好，无非就是增加线性变化度罢了。

‘假设m和n是一致的，（Selef-Attention的情况下，这二者肯定是一致的）

那么输出的矩阵大小就是  $(n,d_{model})$，这里的$d_{model}$和 输入是否一致？是一致的！$d_{model}就是d_l$，就是输入的s的词向量长度！！！！

续图

<img src="E:\02-BUPT\02-EByte\2023夏\008-语言情感分析（马祖耀）20230804\U3EwOx.png" alt="U3EwOx" style="zoom:60%;" />

### Add & Norm

经过上一部的Multi-Head Attention层，输出数据的格式就变成了 $(n,d_{model})$，

已知输入s也是 $(n,d_{model})$的尺寸，所以这两者可以通过残差连接在一起

Add & LayerNorm是有两个部分的，那么久都写一下吧
$$
result_m=LayerNorm(s+MultiHeadAttention(s))\\
result_f=LayerNorm(s+FeedForward(s))
$$


这是它的结构，其实就是一个残差连接，再接上一次Norm。Layer Normalization 会将每一层神经元的输入都转成均值方差都一样的，这样可以加快收敛。

### LayerNorm

就是说他和batchNorm，批量标准化的区别吧

一个词向量是一维，一个[m,n]的矩阵，m是样本数量，n是特征数量，（行向量）

对每一个样本，计算均值和方差，计算得
$$
\hat{x}_{ij}=\frac{x_{ij}-u_i}{\sqrt{\sigma _{i}^{2}+\varepsilon}}
$$
之后再进行标准化后的特征进行线性变换和平移
$$
y_{ij}=\gamma_j\hat{x}_{ij}+\beta_j
$$
其中gamma和 beta分别是可训练的参数，用于对标准化后的特征进行缩放和平移。

与BN不同的是，LN是在特征维度上进行归一化处理。



### Feed Forward

Feed Forward 层比较简单，是一个**两层的全连接层**，第一层的激活函数为 Relu，第二层不使用激活函数，对应的公式如下。
$$
FF = max(0,XW_1+b_1)W_2+b_2
$$


内层max是第一层带relu的，第二层是纯FC

其实就是

FC-> relu -> FC

然后再来一次Add&Norm

如此就构成了Encoder



算一下输出

<img src="E:\02-BUPT\02-EByte\2023夏\008-语言情感分析（马祖耀）20230804\U3EwOx.png" alt="U3EwOx" style="zoom:67%;" />![c7w7rD](E:\02-BUPT\02-EByte\2023夏\008-语言情感分析（马祖耀）20230804\c7w7rD.png)

在第一次Add&Norm之后输出的矩阵大小是 $[n,d_{model}]$

然后经过FeedForward应该是也么有变化的，因为这个也进行了一次残差链接

所以，输出仍然是不变的。

// 所以我大概到现在是明白EE链接DD链接以及ED链接是怎么做的了，唯一的骚操作就是乘了一个矩阵罢了，

### 总结

$$
input\,\, s\in\mathbb{R}^{n\times d_{model}}\\
\begin{cases}
	A_1=QW^Q_i=s^iW^QW^Q_i\\
	A_2=KW^K_i=s^iW^KW^K_i\\
	A_3=VW^K_i=s^iW^VW^K_i
\\
\end{cases}
\\
head_i=Attention(A_1,A_2,A_3) \\
MultiHead(Q,K,V)=Concat(head_1,head_2,...,head_h)W^O \\
x=LayerNorm(s+MultiHeadAttention(s))\\
x=LayerNorm(x+max(0,XW_1+b_1)W_2+b_2)\\

\\
output\,\, x\in\mathbb{R}^{n\times d_{model}}
$$



## Decoder

我觉得我也没啥可以写的了

却别就在于第二层MultiHeadAttention那里，无非就是QKV的来源不同而已。

还有一个点就是Masked Attention这里

Masked的配图我始终不太理解，这他妈不是矩阵吗？怎么  什么时候变成了好像每个点对应一个词了？

好了，反正我理解了。

![c7w7rD](E:\02-BUPT\02-EByte\2023夏\008-语言情感分析（马祖耀）20230804\c7w7rD.png)





## 总结思考



Attention，还有FeedForward里的Relu，是负责获得非线性特征的

各种W矩阵变换以及FC，都是用来获取线性特征的。

所以说，其实好像和神经网络的数学本质，区别不是很大？

当然这么说，实在是太过于概括了。`
	splitMarkdown(text, 250)
}
