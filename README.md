# selpg

简介</br>
使用go语言开发Linux命令行实用程序中的selpg
selpg代表SELect PaGes。selpg 允许用户指定从输入文本抽取的页的范围，这些输入文本可以来自文件或另一个进程。

使用方法</br>
selpg -sstartPage -eendPage [ -f | -llinesPerPage ] [ -ddest ] [ inFilename ]

selpg程序逻辑</br>
如前面所说的那样，selpg是从文本输入选择页范围的实用程序。该输入可以来自作为最后一个命令行参数指定的文件，在没有给出文件名参数时也可以来自标准输入。selpg 首先处理所有的命令行参数。在扫描了所有的选项参数（也就是那些以连字符为前缀的参数）后，如果 selpg 发现还有一个参数，则它会接受该参数为输入文件的名称并尝试打开它以进行读取。如果没有其它参数，则 selpg 假定输入来自标准输入。

参数处理</br>
“-sstartPage”和“--eendPage”强制选项：</br>
selpg 要求用户用两个命令行参数“-sstartPage”（例如，“-s10”表示从第 10 页开始）和“--eendPage”（例如，“-e20”表示在第 20 页结束）指定要抽取的页面范围的起始页和结束页。selpg 对所给的页号进行合理性检查；换句话说，它会检查两个数字是否为有效的正整数以及结束页是否不小于起始页。这两个选项，“-sstartPage”和“--eendPage”是强制性的，而且必须是命令行上在命令名 selpg 之后的头两个参数：
$ selpg -s10 -e20 ...
（... 是命令的余下部分，下面对它们做了描述）。
</br>“-llinesPerPage”和“-f”可选选项：</br>
selpg 可以处理两种输入文本：
类型 1：该类文本的页行数固定。这是缺省类型，因此不必给出选项进行说明。也就是说，如果既没有给出“-llinesPerPage”也没有给出“-f”选项，则 selpg 会理解为页有固定的长度（每页 72 行）。
选择 72 作为缺省值是因为在行打印机上这是很常见的页长度。这样做的意图是将最常见的命令用法作为缺省值，这样用户就不必输入多余的选项。该缺省值可以用“-llinesPerPage”选项覆盖，如下所示：
$ selpg -s10 -e20 -l66 ...
这表明页有固定长度，每页为 66 行。
类型 2：该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C 中用“\f”表示）定界。该格式与“每页行数固定”格式相比的好处在于，当每页的行数有很大不同而且文件有很多页时，该格式可以节省磁盘空间。在含有文本的行后面，类型 2 的页只需要一个字符 ― 换页 ― 就可以表示该页的结束。打印机会识别换页符并自动根据在新的页开始新行所需的行数移动打印头。
将这一点与类型 1 比较：在类型 1 中，文件必须包含 PAGELEN - CURRENTPAGELEN 个新的行以将文本移至下一页，在这里 PAGELEN 是固定的页大小而 CURRENTPAGELEN 是当前页上实际文本行的数目。在此情况下，为了使打印头移至下一页的页首，打印机实际上必须打印许多新行。这在磁盘空间利用和打印机速度方面效率都很低（尽管实际的区别可能不太大）。
类型 2 格式由“-f”选项表示，如下所示：
$ selpg -s10 -e20 -f ...
该命令告诉 selpg 在输入中寻找换页符，并将其作为页定界符处理。
注：“-llinesPerPage”和“-f”选项是互斥的。
</br>“-ddest”可选选项：</br>
selpg 还允许用户使用“-ddest”选项将选定的页直接发送至打印机。这里，“dest”应该是 lp 命令“-d”选项（请参阅“man lp”）可接受的打印目的地名称。该目的地应该存在 ― selpg 不检查这一点。在运行了带“-d”选项的 selpg 命令后，若要验证该选项是否已生效，请运行命令“lpstat -t”。该命令应该显示添加到“dest”打印队列的一项打印作业。如果当前有打印机连接至该目的地并且是启用的，则打印机应打印该输出。在下面的示例中，我们打开到命令
$ lp -ddest
的管道以便输出，并写至该管道而不是标准输出：
selpg -s10 -e20 -dlp1
该命令将选定的页作为打印作业发送至 lp1 打印目的地。您应该可以看到类似“request id is lp1-6”的消息。该消息来自 lp 命令；它显示打印作业标识。如果在运行 selpg 命令之后立即运行命令 lpstat -t | grep lp1 ，您应该看见 lp1 队列中的作业。如果在运行 lpstat 命令前耽搁了一些时间，那么您可能看不到该作业，因为它一旦被打印就从队列中消失了。

输入处理</br>
一旦处理了所有的命令行参数，就使用这些指定的选项以及输入、输出源和目标来开始输入的实际处理。
selpg 通过以下方法记住当前页号：如果输入是每页行数固定的，则 selpg 统计新行数，直到达到页长度后增加页计数器。如果输入是换页定界的，则 selpg 改为统计换页符。这两种情况下，只要页计数器的值在起始页和结束页之间这一条件保持为真，selpg 就会输出文本（逐行或逐字）。当那个条件为假（也就是说，页计数器的值小于起始页或大于结束页）时，则 selpg 不再写任何输出。瞧！您得到了想输出的那些页。

测试结果（这里测试数据就用selpg.go文件）</br>
1.</br>
[cg@localhost selpg]$ ./selgp</br>
./selgp: not enough arguments</br>
</br>
USAGE: ./selgp -sstartPage -eendPage [ -f | -llinesPerPage ] [ -ddest ] [ inFilename ]</br>
2. （startPage > endPage）</br>
[cg@localhost selpg]$ ./selgp -s10 -e5</br>
./selgp: invalid end page 5</br>
</br>
USAGE: ./selgp -sstartPage -eendPage [ -f | -llinesPerPage ] [ -ddest ] [ inFilename ]</br>
3.</br>
[cg@localhost selpg]$ cat selgp.go | ./selgp -s2 -e3 -l3</br>
import (</br>
	"fmt"</br>
	"os"</br>
	"os/exec"</br>
	"strconv"</br>
	"bufio"v
4.</br>
[cg@localhost selpg]$ ./selgp -s2 -e2 -l2</br>
a</br>
b</br>
c</br>
d</br>
end</br>
b</br>
c</br>
5.</br>
[cg@localhost selpg]$ ./selgp -s3 -e7 -l2 selgp.go </br>
	"fmt"</br>
	"os"</br>
	"os/exec"</br>
	"strconv"</br>
	"bufio"</br>
	"io"</br>
	"io/ioutil"</br>
)</br>
</br>
var (</br>
6.</br>
[cg@localhost selpg]$ ./selgp -s3 -e7 -l2 selgp.go | grep os</br>
	"os"</br>
	"os/exec"</br>
7.</br>
[cg@localhost selpg]$ ./selgp -s3 -e7 -l2 selgp.go >result 2>debug </br>
[cg@localhost selpg]$ cat result </br>
	"fmt"</br>
	"os"</br>
	"os/exec"</br>
	"strconv"</br>
	"bufio"</br>
	"io"</br>
	"io/ioutil"</br>
)</br>
</br>
var (</br>
[cg@localhost selpg]$ cat debug</br> 
[cg@localhost selpg]$ </br>
8.</br>
[cg@localhost selpg]$ ./selgp -s-1 -e7 -l2 selgp.go >result 2>debug </br>
[cg@localhost selpg]$ cat result </br>
[cg@localhost selpg]$ cat debug </br>
./selgp: invalid start page -1</br>
</br>
USAGE: ./selgp -sstartPage -eendPage [ -f | -llinesPerPage ] [ -ddest ] [ inFilename ]</br>
9.</br>
[cg@localhost selpg]$ ls</br>
debug  result  selgp  selgp.go</br>
[cg@localhost selpg]$ ./selgp -s2 -e2 -l2 test >output 2>&1</br>
[cg@localhost selpg]$ cat output </br>
./selgp: input file "test" does not exist</br>
</br>
10. （没有打印机，无法进一步测试）</br>
[cg@localhost selpg]$ ./selgp -s3 -e7 -l2 -dp selgp.go </br>
lp: The printer or class does not exist.</br>
