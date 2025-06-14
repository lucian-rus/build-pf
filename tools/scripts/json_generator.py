import os
import json

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

backlog = []
output = {}

for line in content:
    data = line.split(',')
    obj = {}

    for i in range(0, len(header)):
        obj[header[i]] = data[i]
    
    if obj["parent"] == "root":
        if not obj["class"] in output:
            output[obj["class"]] = []

        obj_class = obj["class"]
        to_delete = []
        for key, value in obj.items():
            if value == "":
                to_delete.append(key)

        # cleanup
        for key in to_delete:
            del obj[key]
        del obj["class"]
        del obj["parent"]

        output[obj_class].append(obj)
    else:
        # probably a better way of doing this -> move to tree 
        for item in output["scope"]:
            if obj["parent"] == item["name"]:
                if not "children" in item:
                    item["children"] = []
                
                del obj["parent"]
                item["children"].append(obj)
                break

print(json.dumps(output))

file = open(os.path.dirname(os.path.realpath(__file__)) + "/../resources/ldd_template.json", "w")
file.write(json.dumps(output))
file.close()
