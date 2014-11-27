#!/usr/bin/env python
import urllib
import urllib2

url = "http://likunarmstrong.appspot.com/blog/update"
values = {"test":"bar"}
data = urllib.urlencode(values)
req = urllib2.Request(url, data)
response = urllib2.urlopen(req)
result = response.read()
print result
