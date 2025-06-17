import os
import json

def get_json_file_content(file_path):
    """! reads and returns the raw JSON data """
    file = open(file_path)
    raw_content = json.load(file)
    file.close()

    return raw_content

def create_ldd_project():

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

def run_question_query(item):
    if item["required"] == "false":
        should_add = input("{:64s} > ".format(item["question"]))
        if not should_add == "yes":
            return False
    
    return True

# probably a better way to do this
def add_data(json_input, output, category, subcategory, add_children):
    for item in json_input[category]:
        if not subcategory == None:
            if item["subclass"] == subcategory:
                if run_question_query(item) == True:
                    output += item["text"]
                
                if (add_children == True) and ("children" in item):
                    output += " {\n"
                    for child in item["children"]:
                        if run_question_query(child):
                            output += "    " + child["text"] + "\n"
                    output += "}\n"
                output += "\n"
        else:
            if run_question_query(item) == True:
                    output += item["text"] + "\n"
    output += "\n"

    return output

def translate_json_to_file(json_input, file_path):
    output = ""
    output = add_data(json_input, output, "include", None, False)
    output = add_data(json_input, output, "macro", "ldd_start", False)
    output = add_data(json_input, output, "prototype", None, False)
    output = add_data(json_input, output, "expression", None, False)
    output = add_data(json_input, output, "scope", "generic", True)
    output = add_data(json_input, output, "macro", "ldd_end", False)

    print("----------------------------------------\n", output)

def __main__():
    # create_ldd_project()
    raw_json_content = get_json_file_content(os.path.dirname(os.path.realpath(__file__)) + "/../resources/ldd_template.json")
    translate_json_to_file(raw_json_content, None)

__main__()
