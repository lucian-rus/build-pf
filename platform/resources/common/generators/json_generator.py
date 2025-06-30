import os
import json

def get_csv_file_content(file_path):
    """! reads and returns the raw CSV data """
    file = open(file_path)
    raw_content = file.read()
    file.close()

    return raw_content

def write_json_data_to_file(file_path, output):
    file = open(file_path, "w")
    file.write(json.dumps(output))
    file.close()

def extract_and_remove_header(content):
    content = content.split("\n")

    # get header then remove it
    header = content[0]
    header = header.split(",")
    content = content[1:]
    # remove last element (useless newline)
    content = content[:-1]
    
    return header, content

def parse_csv_file(raw_csv_data, header):
    output = []

    for line in raw_csv_data:
        data = line.split(',')
        temporary_dict_container = {}

        # create temp dict -> skip if value should be empty. this is to clean the objects up
        for i in range(0, len(header)):
            if not data[i] == "":
                temporary_dict_container[header[i]] = data[i]
        output.append(temporary_dict_container)

    return output

def setup_output_dict(backlog):
    output = {}
    # list where to add indexes for items to be deleted
    to_delete_index_list = []

    backlog_item_counter = 0
    for item in backlog:
        if item["parent"] == "root":
            # ensure existence of top layer
            if not item["class"] in output:
                output[item["class"]] = []

            temporary_class_data = item["class"]
            # those are not required anymore, so delete them for more clarity in the JSON file
            del item["parent"]
            del item["class"]
            output[temporary_class_data].append(item)
            to_delete_index_list.append(backlog_item_counter)

        # increase this counter
        backlog_item_counter += 1
    
    # run backlog cleanup -> go backwards
    for index in reversed(to_delete_index_list):
        del backlog[index]

    return output, backlog

# add support for nested scopes
def parse_backlog(backlog, output):
    # list where to add indexes for items to be deleted
    to_delete_index_list = []

    backlog_item_counter = 0
    for item in backlog:
        # probably a better way to do this -> do this recursively
        for scope in output["scope"]:
            if item["parent"] == scope["name"]:
                if not "children" in scope:
                    scope["children"] = []
                

                scope["children"].append(item)
                to_delete_index_list.append(backlog_item_counter)
    # run backlog cleanup -> go backwards
    for index in reversed(to_delete_index_list):
        del backlog[index]

    return output, backlog

def __main__():
    raw_content = get_csv_file_content(os.path.dirname(os.path.realpath(__file__)) + "/../resources/ldd_template.csv")
    header, headerless_content = extract_and_remove_header(raw_content)

    # while could be called parsed data, is called backlog because it represents the data that shall be parsed
    backlog = parse_csv_file(headerless_content, header)
    output, backlog = setup_output_dict(backlog)

    # do this until no items are left in backlog -> does not support nesting yet
    while not len(backlog) == 0:
        print(backlog)
        output, backlog = parse_backlog(backlog, output)

    write_json_data_to_file(os.path.dirname(os.path.realpath(__file__)) + "/../resources/ldd_template.json", output)

__main__()
