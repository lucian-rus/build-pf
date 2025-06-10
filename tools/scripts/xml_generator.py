import os

# looks dumb but works
file = open(os.path.dirname(os.path.realpath(__file__)) + "/../resources/ldd_template.csv")
content = file.read()
file.close()

content = content.split("\n")

# get header then remove it
header = content[0]
header = header.split(",")
content = content[1:]
# remove last element (useless newline)
content = content[:-1]

output = """<?xml version="1.0" encoding="UTF-8"?>
<body>\n"""

for line in content:
    data = line.split(',')
    output += "    <item>\n"
    output += "        <{start}>{content}</{stop}>\n".format(start=header[0], content=data[0], stop=header[0])
    output += "        <{start}>{content}</{stop}>\n".format(start=header[1], content=data[1], stop=header[1])
    output += "        <{start}>{content}</{stop}>\n".format(start=header[2], content=data[2], stop=header[2])
    output += "        <{start}>{content}</{stop}>\n".format(start=header[3], content=data[3], stop=header[3])
    output += "        <{start}>{content}</{stop}>\n".format(start=header[4], content=data[4], stop=header[4])
    output += "    </item>\n\n"

output += "</body>\n"
file = open(os.path.dirname(os.path.realpath(__file__)) + "/../resources/ldd_template.xml", "w")
file.write(output)
file.close()
