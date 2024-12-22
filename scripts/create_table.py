import os
import pandas as pd
import sys

from utils import parse_log_file

def create_table(base_directory, attribute):
    data_dict = {}
    unique_V_values = set()
    has_edges = False

    for filename in os.listdir(base_directory):
        filepath = os.path.join(base_directory, filename)
        
        if os.path.isfile(filepath):
            try:
                parsed_data_array = parse_log_file(filepath)
                
                for parsed_data in parsed_data_array:
                    if 'nodes' in parsed_data and attribute in parsed_data:
                        unique_V_values.add(parsed_data['nodes'])
                        
                        if filename not in data_dict:
                            data_dict[filename] = {}

                        data_dict[filename][parsed_data['nodes']] = parsed_data[attribute]

                        if 'edges' in parsed_data:
                            has_edges = True
                            if 'edges' not in data_dict:
                                data_dict['edges'] = {}
                            data_dict['edges'][parsed_data['nodes']] = parsed_data['edges']
            
            except Exception as e:
                print(f'Error processing {filename}: {e}')

    unique_V_values = sorted(unique_V_values)
    df_data = {'V': unique_V_values}

    if has_edges:
        df_data['edges'] = [data_dict['edges'].get(v, None) for v in unique_V_values]
    else:
        df_data['edges'] = [None for _ in unique_V_values]

    for filename in data_dict:
        if filename not in ['edges']:
            df_data[filename] = [data_dict[filename].get(v, None) for v in unique_V_values]

    df = pd.DataFrame(df_data)
    df = df[['V', 'edges'] + [col for col in df.columns if col not in ['V', 'edges']]]

    latex_table = df.to_latex(index=False)
    print(latex_table)

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print('Usage: create_table.py <BASE_DIRECTORY> <ATTRIBUTE_NAME>')
        sys.exit(1)

    base_directory = sys.argv[1]
    attribute = sys.argv[2]
    create_table(base_directory, attribute)
