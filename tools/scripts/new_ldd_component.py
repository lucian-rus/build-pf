import os
import xml.etree.ElementTree as ET

### ensure that these are generated automatically -> shall use dictionaries and everything shall be generic
includes = ""
macros = ""
prototypes = ""
functions = ""
ldd_start_macros = ""
ldd_end_macros = ""
declarations = ""

XML_NAME_POS = 0
XML_CATEGORY_POS = 2
XML_PARENT_POS = 3
XML_CODE_POS = 4

function_obj = {}

### this function parses and returns a list of blocks
def get_template_block_list():
    return ""

def create_ldd_project():
    global includes
    global macros
    global prototypes
    global functions
    global ldd_start_macros
    global ldd_end_macros
    global declarations

    print("#=====================================")
    print("creating device driver project")
    ldd_name = input("> set the name of the linux device driver: ")
    # if os.path.exists(ldd_name + "-device-driver"):
    #     print(">> the directory already exists!")
    #     return
    
    # os.makedirs(ldd_name + "-device-driver")
    # os.makedirs(ldd_name + "-device-driver/src")
    # os.makedirs(ldd_name + "-device-driver/inc")
    # os.makedirs(ldd_name + "-device-driver/doc")

    output = ""

    tree = ET.parse(os.path.dirname(os.path.realpath(__file__)) + "/../resources/ldd_template.xml")
    root = tree.getroot()
    
    for i in range(0, len(root)):
        if root[i][XML_CATEGORY_POS].text == "include":
            includes += root[i][XML_CODE_POS].text + "\n"
        
        if root[i][XML_CATEGORY_POS].text == "ldd_start_macro":
            ldd_start_macros += root[i][XML_CODE_POS].text + "\n"

        if root[i][XML_CATEGORY_POS].text == "ldd_end_macro":
            ldd_end_macros += root[i][XML_CODE_POS].text + "\n"

        if root[i][XML_CATEGORY_POS].text == "function":
            # functions return types shall be handled
            # functions += root[i][XML_CODE_POS].text + " {\n}\n\n"
            function_obj[root[i][XML_NAME_POS].text] = {"name": root[i][XML_CODE_POS].text, "children": []}
            
        if root[i][XML_CATEGORY_POS].text == "macro":
            macros += root[i][XML_CODE_POS].text + "\n"

        if root[i][XML_CATEGORY_POS].text == "prototype":
            prototypes += root[i][XML_CODE_POS].text + "\n"

        if root[i][XML_CATEGORY_POS].text == "function_call":
            function_obj[root[i][XML_PARENT_POS].text]["children"].append(root[i][XML_CODE_POS].text)

        if root[i][XML_CATEGORY_POS].text == "declaration":
            declarations += root[i][XML_CODE_POS].text + "\n"

    output += includes + "\n"
    output += ldd_start_macros + "\n"
    # output += macros + "\n"
    # output += prototypes + "\n"
    output += declarations + "\n"

    print(function_obj)
    for key, values in function_obj.items():
        func = values["name"] + " {\n"
        for item in values["children"]:
            func += "   " + item + "\n"

        func += "}\n\n"
        output += func

    # output += functions + "\n"
    output += ldd_end_macros + "\n"

    # special characters to be replaced
    output = output.replace("@1@", "<")
    output = output.replace("@2@", ">")
    output = output.replace("@3@", ",")
    output = output.replace("@4@", "&")
    output = output.replace("``template_module_license``", "GPL")
    output = output.replace("``template_module_author``", "test-author")
    output = output.replace("``template_module_description``", "test-description")
    output = output.replace("``template_ldd_name``", ldd_name)

    print(output)
    file = open(os.path.dirname(os.path.realpath(__file__)) + "/../../" + ldd_name + "-device-driver/src/" + ldd_name + ".c", "w")
    file.write(output)
    file.close()

create_ldd_project()
