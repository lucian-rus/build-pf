import os
import json
import shutil

# these paths shall be relative to the parent dir
CONFIG_FILE_PATH = "../../config.json"
TEMP_DIR_PATH = "../../tools/.temp"

SCRIPT_DIR_PATH = os.path.dirname(os.path.abspath(__file__))


########################################################################################
#                               helper functions
########################################################################################
def get_parsed_json_data():
    """
    Reads and parses the JSON configuration file.

    Returns
    -------
    dict
        Parsed JSON data from the configuration file.
    """
    file = open(os.path.join(SCRIPT_DIR_PATH, CONFIG_FILE_PATH))
    json_data = json.load(file)

    file.close()
    return json_data


########################################################################################
#                               parser functions
########################################################################################
def get_data_from_json(key):
    """
    Retrieves the value from the JSON File based on the given key.

    Returns
    -------
    generic data
        Value of the JSON key.
    """
    json_data = get_parsed_json_data()
    return json_data[key]


########################################################################################
#                               runner functions
########################################################################################
def run_post_build_cleaner(directory):
    """
    Cleans the output directory by removing all files except from the ones from the exclusion list.
    """
    print("> running post build cleaner")
    cleaner_exception_list = get_data_from_json("output_cleaner_exceptions")
    for root, dirs, files in os.walk(directory):
        for name in files:
            if name == "Makefile":
                continue

            # dumb but works -> ensure that `mod.c` files are deleted as well if they are not in the exception list
            if (name.endswith(".c") or name.endswith(".h")) and 1 == name.count("."):
                continue

            # create extension here
            if any(name.endswith("." + ext) for ext in cleaner_exception_list):
                continue

            os.remove(os.path.join(root, name))


# @todo: maybe make this more generic?
def run_post_clean_mover(directory):
    """
    Moves files to the output directory.
    """
    print("> running post build mover")
    mover_exception_list = get_data_from_json("output_mover_exceptions")
    output_directory_name = get_data_from_json("output_mover_target_dir")
    for root, dirs, files in os.walk(directory):
        for name in files:
            if name == "Makefile":
                continue

            # dumb but works -> ensure that `mod.c` files are deleted as well if they are not in the exception list
            if (name.endswith(".c") or name.endswith(".h")) and 1 == name.count("."):
                continue

            # create extension here
            if any(name.endswith("." + ext) for ext in mover_exception_list):
                continue

            shutil.move(
                os.path.join(root, name),
                os.path.join(os.getcwd(), output_directory_name, name),
            )


def restore_temporary_files():
    """
    Removes temporary directory and restores temporary files.
    """
    print("> removing temporary directory")
    temp_dir = os.path.join(SCRIPT_DIR_PATH, TEMP_DIR_PATH)
    temp_makefile = os.path.join(temp_dir, "Makefile")
    makefile_dir_path = get_data_from_json("makefile_directory_path")
    makefile_path = os.path.join(SCRIPT_DIR_PATH, makefile_dir_path)
    shutil.copy(temp_makefile, makefile_path)
    shutil.rmtree(temp_dir)


def post_build_runner(directory):
    """
    Runs all post-build steps, such as removing directories and moving files.
    """
    run_post_build_cleaner(directory)
    run_post_clean_mover(directory)
    restore_temporary_files()


if __name__ == "__main__":
    print("------------------------------------------")
    print("         running post-build script")
    print("------------------------------------------")
    script_dir = os.path.dirname(os.path.abspath(__file__))
    makefile_dir_path = get_data_from_json("makefile_directory_path")
    dir_path = os.path.join(script_dir, makefile_dir_path)
    post_build_runner(dir_path)
    print("> done")
