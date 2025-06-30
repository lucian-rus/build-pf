import json
import os
import shutil


# this path shall be relative to the root of the platform project
PLATFORM_JSON_PATH = "./platform/scripts/platform-presets.json"


def file_move_function(source_path, dest_path):
    if not os.path.exists(dest_path):
        os.makedirs(dest_path)

    for entry in os.listdir(source_path):
        full_source_path = os.path.join(source_path, entry)
        if os.path.isfile(full_source_path):
            print(f"moving {full_source_path} to {dest_path}")
            shutil.copy(full_source_path, dest_path)


# don't put too much effort into this, as it'll be deleted soon
def copy_driver_files(json_data):
    if not os.path.exists("./driver"):
        os.makedirs("./driver")
        os.makedirs("./driver/template_ldd")
        os.makedirs("./driver/template_ldd/src")
        os.makedirs("./driver/template_ldd/inc")

    source_path = os.path.join(json_data["templates"]["source"], "template_ldd.c")
    shutil.copy(source_path, "driver/template_ldd/src")

    source_path = os.path.join(json_data["templates"]["source"], "template_ldd.h")
    print(source_path, "driver/inc")
    shutil.copy(source_path, "driver/template_ldd/inc")

    source_path = os.path.join(json_data["templates"]["source"], "Makefile")
    print(source_path, "driver")
    shutil.copy(source_path, "driver")


def setup_project(json_data, proj_type):
    # if there is no resource dir, the type of project is not yet supported
    if "resource_dir" not in json_data:
        print("     !!! NOT YET SUPPORTED")
        return

    # # this shall always happen
    # file_move_function(json_data["resource_dir"], "./")

    # # looks dumb but works
    # for key, value in list(json_data.items())[1:]:
    #     if not value["required"]:
    #         choice = input(f"{value['question']} - yes/no: ").strip()
    #         if choice == "no":
    #             continue
    #     file_move_function(value["source"], value["dest"])

    # @todo make this properly configurable through the json
    if proj_type == "ldd_project":
        copy_driver_files(json_data)


def main():
    # Read options from JSON file
    with open(PLATFORM_JSON_PATH, "r") as f:
        options = json.load(f)

    # Display options
    print("Select an option:")
    for idx, key in enumerate(options.keys(), 1):
        print(f"{idx}. {key}")

    choice = input(f"Enter your choice (1-{len(options)}): ").strip()

    try:
        choice_idx = int(choice) - 1
        selected_key = list(options.keys())[choice_idx]
        print(f"You selected: {selected_key}")

        # scripts, resources and templates are required
        if (
            ("scripts" not in options[selected_key])
            or ("resources" not in options[selected_key])
            or ("templates" not in options[selected_key])
        ):
            print("     !!! DOES NOT SATISFY THE CONDITIONS. CHECK PRESETS")
            return

        setup_project(options[selected_key], selected_key)
    except (ValueError, IndexError):
        print("Invalid choice. Please enter a valid number.")


if __name__ == "__main__":
    main()
