# 珠玑

珠玑者算机语也。

## 获取

二进制包：[Linux](https://dl.elvish.io/%e7%8f%a0%e7%8e%91-linux)、[OS X](https://dl.elvish.io/%e7%8f%a0%e7%8e%91-osx)。

从源码构建：`go get -u github.com/xiaq/zhuji/zhuji`。

## 例程

基本算术：

```
珠玑> 九加九乘十。
一百八十。
```

定义函数、操作堆栈：

```
珠玑> 倍者自加也。
珠玑> 方者自乘也。
珠玑> 二百三十三、倍。
四百六十六。
珠玑> 二百三十三、方。
四百六十六、五万四千二百八十九。
珠玑> 乘。
二千五百二十九万八千六百七十四。
珠玑> 复。
二千五百二十九万八千六百七十四、二千五百二十九万八千六百七十四。
珠玑> 弃、弃。
珠玑> 弃。
无元。
```

递归函数、条件语句：

```
珠玑> 阶乘者，复、若等于零则弃、一；非，复、减一、阶乘，乘，毕。
珠玑> 十、阶乘。
三百六十二万八千八百。
珠玑> 弃。
珠玑> 斐波那契者，复、若小于三则弃、一；非、复、减一、斐波那契、易、减二、斐波那契、和；毕。
珠玑> 十、斐波那契。
五十五。
```
