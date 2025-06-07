import os

### this function parses and returns a list of blocks
def get_template_block_list():
    return ""

def create_ldd_project():
    print("#=====================================")
    print("creating device driver project")
    os.makedirs("device-driver")
    os.makedirs("device-driver/src")
    os.makedirs("device-driver/inc")

def create_library_project():
    print("#=====================================")
    print("creating library project")
    os.makedirs("library")
    os.makedirs("library/src")
    os.makedirs("library/inc")
    
def create_generic_project():
    print("#=====================================")
    print("creating generic project")
    os.makedirs("components-tt")
    os.makedirs("components-tt/TestComponent")
    os.makedirs("components-tt/TestComponent/src")
    os.makedirs("components-tt/TestComponent/inc")
    os.makedirs("app")

def get_input():
    print("select type of project to generate:")
    print("* ldd project        -> 1")
    print("* library            -> 2")
    print("* generic platform   -> 3")
    print("#=====================================")
    in_data = input()
    if(in_data == "1"):
        create_ldd_project()
    if(in_data == "2"):
        create_library_project()
    if(in_data == "3"):
        create_generic_project()

    print("selected: ", in_data)

def check_init():
    if os.path.exists("/settings"):
        return
    
    get_input()

check_init()
