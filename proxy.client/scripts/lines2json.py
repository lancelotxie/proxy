import json
import os
import sys


def convert(file):
    f = open(file)
    lines = f.readlines()
    f.close()

    ss = []
    for line in lines:
        line = line.replace("\n", "")
        ss.append(line)

    print("一共有", len(ss), "行被转换")
    j = json.dumps(ss)

    new_file_name = newFileName(file)
    f = open(new_file_name, mode='w')
    f.write(j)
    f.close()


def realFileName(file):
    base = os.path.splitext(file)[0]
    return base


def newFileName(old):
    old = realFileName(old)
    new = old+".json"
    return new


if __name__ == "__main__":
    filename = sys.argv[1]
    print("转换中：", filename)
    convert(filename)
    print("转换完成")
