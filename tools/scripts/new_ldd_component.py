import os
import xml.etree.ElementTree as ET

includes = ""
macros = ""
prototypes = ""
functions = ""
ldd_start_macros = ""
ldd_end_macros = ""

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

    print("#=====================================")
    print("creating device driver project")
    ldd_name = input("> set the name of the linux device driver: ")
    if os.path.exists(ldd_name + "-device-driver"):
        print(">> the directory already exists!")
        return
    
    os.makedirs(ldd_name + "-device-driver")
    os.makedirs(ldd_name + "-device-driver/src")
    os.makedirs(ldd_name + "-device-driver/inc")
    os.makedirs(ldd_name + "-device-driver/doc")

    output = ""

    tree = ET.parse(os.path.dirname(os.path.realpath(__file__)) + "/../resources/ldd_template.xml")
    root = tree.getroot()
    
    for i in range(0, len(root)):
        if root[i][2].text == "include":
            includes += root[i][3].text + "\n"
        
        if root[i][2].text == "ldd_start_macro":
            ldd_start_macros += root[i][3].text + "\n"

        if root[i][2].text == "ldd_end_macro":
            ldd_end_macros += root[i][3].text + "\n"

        if root[i][2].text == "function":
            # functions types shall be handled
            functions += root[i][3].text + " {\n}\n\n"
            
        if root[i][2].text == "macro":
            macros += root[i][3].text + "\n"

        if root[i][2].text == "prototype":
            prototypes += root[i][3].text + "\n"


    output += includes + "\n"
    output += ldd_start_macros + "\n"
    # output += macros + "\n"
    # output += prototypes + "\n"
    output += functions + "\n"
    output += ldd_end_macros + "\n"

    # special characters to be replaced
    output = output.replace("@1@", "<")
    output = output.replace("@2@", ">")
    output = output.replace("``template_module_license``", "GPL")
    output = output.replace("``template_module_author``", "test-author")
    output = output.replace("``template_module_description``", "test-description")
    output = output.replace("``template_ldd_name``", ldd_name)

    print(output)
    file = open(os.path.dirname(os.path.realpath(__file__)) + "/../../" + ldd_name + "-device-driver/src/" + ldd_name + ".c", "w")
    file.write(output)
    file.close()

create_ldd_project()
