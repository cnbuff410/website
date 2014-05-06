#!/usr/bin/env python

from docutils.core import publish_parts
import sys
import os
import codecs

filename = sys.argv[1]
filebase = os.path.basename(filename).split(".")[0]
with open (filename, "r") as myfile:
    data=myfile.read()

wfile = codecs.open("../posts/" + filebase + ".html", "w", "utf-8")
txt = publish_parts(data, writer_name='html')['html_body']
print txt
wfile.write(txt)
