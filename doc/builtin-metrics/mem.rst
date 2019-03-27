.. _builtin-mem:

内存占用
========

mem.memtotal
    :意义: 内存总量
    :取值: 0 - 无上限，整数，单位：字节
    :Tags: 无

    .. note::
        ``memtotal`` 真的会在一些极端情况下变的，并不是无意义的值

mem.memused
    :意义: 内存用量
    :取值: 0 - 无上限，整数，单位：字节
    :Tags: 无

    .. note::
        这个值不包括 page cache 和 buffer
        是通过 ``memtotal`` - ``memfree`` 算出来的

mem.memfree
    :意义: 内存剩余
    :取值: 0 - 无上限，整数，单位：字节
    :Tags: 无

    .. note::
        这个值包括 page cache 和 buffer。

mem.swaptotal
    :意义: 交换内存总量
    :取值: 0 - 无上限，整数，单位：字节
    :Tags: 无

mem.swapused
    :意义: 交换内存用量
    :取值: 0 - 无上限，整数，单位：字节
    :Tags: 无

mem.swapfree
    :意义: 交换内存剩余
    :取值: 0 - 无上限，整数，单位：字节
    :Tags: 无

mem.memfree.percent
    :意义: 内存剩余百分比
    :取值: 0 - 100.0
    :Tags: 无

mem.memused.percent
    :意义: 内存用量百分比
    :取值: 0 - 100.0
    :Tags: 无

mem.swapfree.percent
    :意义: 交换内存剩余百分比
    :取值: 0 - 100.0
    :Tags: 无

mem.swapused.percent
    :意义: 交换内存用量百分比
    :取值: 0 - 100.0
    :Tags: 无
