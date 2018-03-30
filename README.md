# selpg

简介</br>
使用go语言开发Linux命令行实用程序中的selpg
selpg代表SELect PaGes。selpg 允许用户指定从输入文本抽取的页的范围，这些输入文本可以来自文件或另一个进程。

使用方法</br>
selpg -sstartPage -eendPage [ -f | -llinesPerPage ] [ -ddest ] [ inFilename ]

测试结果（这里测试数据就用lines文件，共500line）</br>
测试数据文件为lines:</br>
 ![image](https://github.com/wuchg6/selpg/raw/master/lines.png)</br>
1.测试错误输入的情况</br>
 ![image](https://github.com/wuchg6/selpg/raw/master/errorTest.png)</br>
2.测试程序的输入方式</br>
 ![image](https://github.com/wuchg6/selpg/raw/master/inputTest.png)</br>
3.测试程序结果的输出方式</br>
 ![image](https://github.com/wuchg6/selpg/raw/master/outputTest.png)</br>
4.测试lp（没有打印机，所以无法测试正确的情况）</br>
 ![image](https://github.com/wuchg6/selpg/raw/master/lpTest.png)</br>

