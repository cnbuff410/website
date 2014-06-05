=======================
Matplot-python用法速查
=======================
:Author: cnbuff410
:Contact: likunarmstrong@gmail.com

.. contents::

.. role:: python(code)
   :language: python

典型Matplot例子
===================

.. code-block:: python

    import numpy as np
    import matplotlib.pyplot as plt

    def f(t):
        return np.exp(-t) * np.cos(2 * np.pi * t)

    t1 = np.arange(0.0, 5.0, 0.1)
    t2 = np.arange(0.0, 5.0, 0.02)

    plt.figure(1)
    plt.subplot(211)
    plt.plot(t1, f(t1), 'bo', t2, f(t2), 'k')
    plt.subplot(212)
    plt.plot(t2, np.cos(2 * np.pi * t2), ‘r–’)


题目，坐标轴设置
===================

最简单的就是通过title，xlabel和ylabel设置，比如

:python:`xlabel("xxxx", fontsize=18, color="red")`

这里title是某个axis的title，如果希望整个figure有title，那么需要用suptitle() 函数，用法和title是一样的。
x和y轴的范围设定是通过

:python:`plt.axis([0, 6, 0, 20])`

这就是表示x轴是0到6，y轴是0到20。也可以用xlim和ylim单独控制，比如

:python:`ylim( ymin, ymax )`

当然还可以用yscale和xscale指定坐标轴的scale类型，分别是‘linear’ | ‘log’ | ‘symlog’，比如:

:python:`xscale("log")`

还可以通过yaxis.grid和xaxis.grid控制怎么显示网格，是否显示网格。grid的函数是这样的:

:python:`grid(self, b=None, which='major', **kwargs)`

如果用yaxis.grid (false)，就表示不用显示y轴方向的网格。

有时候我们需要在坐标轴上放文字，比如在x轴上放置月份名称之类，这就需要xticks函数。其用法是比如:

:python:`xticks( arange(5), ('Tom', 'Dick', 'Harry', 'Sally', 'Sue') )`

这里，第一个参数表示放置位置，第二个是放置内容。

设置曲线特征:

.. code-block:: python

    lines = plt.plot(x1, y1, x2, y2)
    # use keyword args
    plt.setp(lines, color='r', linewidth=2.0)
    # or matlab style string value pairs
    plt.setp(lines, 'color', 'r', 'linewidth', 2.0)

如果希望知道有哪些属性可以指定，可以传递对象给setp()函数，比如:

.. code-block:: python

    lines = plt.plot([1,2,3])
    plt.setp(lines)

具体能控制的属性参见
`pyploy tutorial <http://matplotlib.sourceforge.net/users/pyplot_tutorial.html>`_.

文字处理text()命令可以被用来放置文字在绝对坐标轴上，用法是:

:python:`text(60, .025, r'$\mu=100,\ \sigma=15$', horizontalalignment='left', verticalalignment='top')`

放置的东西前面加个r表示包括数学公式。具体能控制的属性可以参见
`属性列表 <http://matplotlib.sourceforge.net/users/text_props.html#text-properties>`_. 具体的数
学公式和符号的表达方式参见
`mathtext <http://matplotlib.sourceforge.net/users/mathtext.html>`_.

这只是放置文字，还有一种需求是要求放置文字标注, 这时候就要用到anotate()函数了，用法大概是这样的:

:python:`plt.annotate('local max', xy=(2, 1), xytext=(3, 1.5), arrowprops=dict(facecolor='black', shrink=0.05))`

第一个参数是要放置的文字，第二个参数是要指向的坐标，第三个是放置文字的坐标，最后一个是箭头属性。相关属性设置参见
`这里 <http://matplotlib.sourceforge.net/users/annotations.html#annotations-tutorial>`_.

对于经常写公式的人来讲，可能latex的格式会感觉更加常用一点。这时候，可以在放置text之前先用:

:python:`rc('text', usetex=True)`

然后与文字有关的地方都用r”"来引用，比如:

:python:`title(r"\TeX\ is Number $\displaystyle\sum_{n=1}^\infty\frac{-e^{i\pi}}{2^n}$!", fontsize=16, color='r')`

画函数
============

plt.plot最简单的形式就是take一组列表数据:

:python:`plt.plot([1,2,3])`

此时plot自动将其认为是y轴的数据，x轴数据自动分配，相当于是:

:python:`plt.plot([0,1,2],[1,2,3])`

也就是说，对于plot，标准形式是:

:python:`plt.plot([0,2,3,4], [1,4,9,16])`

当然，后面可以跟控制表示形式的参数，比如:

:python:`plt.plot([1,2,3,4], [1,4,9,16], 'ro')`

就表示用红色点表示，而不是直线。参数默认是’b-’，表示蓝色直线。

plot还支持同时在一张图里面绘制多个图形，比如:

.. code-block:: python

    import matplotlib.pyplot as plt
    plt.plot(x1, y1, 'r--', x2, y2, 'bs', x3, y3, 'g^')

就可以用不同的形式绘制三条曲线。

画子图
=======

用plt.subplot(xxx)，第一个数字表示numrows，第二个表示numcols，第三个表示fignum。fignum的数量最高也就是numrows乘以numcols。
如果是想产生多个图，每个图包含几个子图，那用figure()来控制图的数量。

画直方图
=========

用函数:

:python:`hist(x, bins=10, range=None, normed=False, cumulative=False, bottom=None, histtype='bar', align='mid', orientation='vertical', rwidth=None, log=False, **kwargs)`

具体用法参考
`hist
<http://matplotlib.sourceforge.net/api/pyplot_api.html#matplotlib.pyplot.hist>`_.

画分布图
=========

例子:
:python:`scatter(x, y, s=20, c='b', marker='o', cmap=None, norm=None, vmin=None, vmax=None, alpha=1.0, linewidths=None, verts=None,**kwargs )`

具体参数看
`scatter
<http://matplotlib.sourceforge.net/api/pyplot_api.html#matplotlib.pyplot.scatter>`_ ，唯一一个不是很确定的是s的取值，我曾经估计是s值要和x和y的列表长
度一样，也就是大概是点的个数？但具体画的时候发现不是这样，点的个数取决于x和y的
列表长度。有待研究。

画Box Plot图
=============

用matplot画Box Plot非常方便，直接用函数:

:python:`boxplot(x, notch=0, sym='+', vert=1, whis=1.5, positions=None, widths=None)`

你直接把一组数用x送进去，出来就是Box图。如果你想一个图表示多个Box，也很简单，让xappend多个list，出来的就是这些list的box图。一个例子如下:

.. code-block:: python

    for attribute in iris[type].keys():
        list = iris[type][attribute]
        data.append(list)
        subplot(2,2,count)
    title("Boxplot of %s" % type, fontsize=18, color="red")
    xticks(arange(5),("","Sepal len","Sepal wid","Petal len","Petal wid"))
    ylabel("Value")
    boxplot(data,0,'gd')

画幂函数图
============

例子:

.. code-block:: python

    from matplotlib.matlab import *

    x = linspace(-4, 4, 200)
    f1 = power(10, x)
    f2 = power(e, x)
    f3 = power(2, x)

    plot(x, f1, 'r', x, f2, 'b', x, f3, 'g', linewidth=2)
    axis([-4, 4, -0.5, 8])
    text(1, 7.5, r'$10^x$', fontsize=16)
    text(2.2, 7.5, r'$e^x$', fontsize=16)
    text(3.2, 7.5, r'$2^x$', fonsize=16)
    title('A simple example', fontsize=16)

    savefig('power.png', dpi=75)
    show()

显示图形中的数学公式
=======================

Matplotlib 可以支持一部分 TeX 的排版指令，因此用户在绘制含有数学公式的图形时会感到很方便并且可以得到比较满意的显示效果，所需要的仅仅是一些 TeX 的排版知识。下面的这个例子显示了如何在图形的不同位置上, 如坐标轴标签，图形的标题以及图形中适当的位置处，显示数学公式。相应的 Python 程序如下:

.. code-block:: python

    from matplotlib.matlab import *

    def f(x, c):
        m1 = sin(2*pi*x)
        m2 = exp(-c*x)
        return multiply(m1, m2)

    x = linspace(0, 4, 100)
    sigma = 0.5
    plot(x, f(x, sigma), 'r', linewidth=2)
    xlabel(r'$\rm{time} \ t$', fontsize=16)
    ylabel(r'$\rm{Amplitude} \ f(x)$', fontsize=16)
    title(r'$f(x) \ \rm{is \ damping \ with} \ x$', fontsize=16)
    text(2.0, 0.5, r'$f(x) = \rm{sin}(2 \pi x^2) e^{\sigma x}$', fontsize=20)
    savefig('latex.png', dpi=75)
    show()
