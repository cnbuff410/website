=======================
TITLE
=======================
:Author: cnbuff410
:Contact: likunarmstrong@gmail.com

这个想法并不是心血来潮，而是经过了长期的思考，以及个人性格导致的结果。我在7年前
做这个Blog的时候，就曾经
`说过 <http://www.kunli.info/2007/04/26/about-this-blog/>`_:

    从开始写blog到现在，先后经历了bokee 网，sina以及现在同步更新的blogspot等阶
    段，还用上过google提供的googlepages作为主页，最后还是决心自己用wordpress建
    一个小站。真的是下了很大的决心。虽然很舍不得，但是不得不换，只有这样，网站
    的程序和数据库才都是属于自己的，页面风格才能由自己决定。只要这个世界还有空
    间在，我的小站就不会倒，而且一直不变

换句话说，我从来都对“控制”自己的站点有强烈的欲望。我10年前开始写Blog时，用的是
Bokee网，什么都是别人的系统，只有内容是自己的。后来换成Wordpress，可以简单控制
一些外观，系统特性等一些东西了。但自己毕竟没学过，而且不想学PHP，导致我能对
Wordpress的订制度很低，况且WP已经是一个非常大的系统了，去自定制WP也不是一个明智
的事情。

这次转换系统和站点，就是为了从根本上解决这个问题。新的Blog站点系统是一个我用
Go_ 写的简单系统，CSS是基于模版但是进行了大量更改，所有的Post都是用
reStructuredText_ 写的文本文件，然后转换为HTML被Blog引擎调用。也就是说，整个系
统从头到尾都是由我控制，没有任何第三方框架的限制。

也许这是大部分Geek最终不可避免的命运吧。


.. _Go: http://www.golang.org/
.. _reStructuredText: http://docutils.sourceforge.net/rst.html
