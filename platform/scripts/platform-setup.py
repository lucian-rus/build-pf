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


def setup_project(json_data):
    # if there is no resource dir, the type of project is not yet supported
    if "resource_dir" not in json_data:
        print("     !!! NOT YET SUPPORTED")
        return
    
    # this shall always happen
    file_move_function(json_data["resource_dir"], "./test_dir")

    # looks dumb but works
    for key, value in list(json_data.items())[1:]:
        if value["required"] == False:
            choice = input(f"{value["question"]} - yes/no: ").strip()
            if choice == "no":
                continue
        file_move_function(value["source"], value["dest"])

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
        setup_project(options[selected_key])
    except (ValueError, IndexError):
        print("Invalid choice. Please enter a valid number.")


if __name__ == "__main__":
    main()
