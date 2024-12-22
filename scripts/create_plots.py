import os
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.ticker import LogLocator, FuncFormatter

from utils import LABELS, TYPES, BASE_DIRECTORY, parse_log_file

PLOTS_DIRECTORY = 'plots'

def format_ticks(value, _):
    if value >= 1 and value < 1000:
        return f'{value:.0f}' # No decimals for values between 1 and 1000
    elif value >= 1000:
        return f'$10^{{{int(np.log10(value))}}}$' # Powers of 10 for larger values
    elif value > 0:
        return f'{value:.2f}' # Two decimals for values between 0 and 1
    else:
        return '0' # Display 0 for zero values


def set_axes_scaling(plt, log_scale_x, log_scale_y):
    if log_scale_x:
        plt.xscale('log')
        plt.gca().xaxis.set_major_locator(LogLocator(base=10.0))
        plt.gca().xaxis.set_minor_locator(LogLocator(base=10.0, subs='auto'))
    
    plt.gca().xaxis.set_major_formatter(FuncFormatter(format_ticks))

    if log_scale_y:
        plt.yscale('log')
        plt.gca().yaxis.set_major_locator(LogLocator(base=10.0))
        plt.gca().yaxis.set_minor_locator(LogLocator(base=10.0, subs='auto'))
    
    plt.gca().yaxis.set_major_formatter(FuncFormatter(format_ticks))


def plot_it(log_scale_x, log_scale_y, metric, output_dir, label):
    set_axes_scaling(plt, log_scale_x, log_scale_y)

    plt.grid(True, which='both', ls='--')
    plt.xlabel(LABELS[label])
    plt.ylabel(LABELS[metric])
    plt.legend()

    output_path = os.path.join(output_dir, f'{metric}.pdf')
    plt.savefig(output_path)
    plt.close()


def plot_with_distinct_styles(x, y, label, index):
    line_styles = ['-', '--', '-.', ':']
    markers = ['s', 'D', '^', 'o', 'x']
    
    line_style = line_styles[index % len(line_styles)]
    marker = markers[index % len(markers)]
    
    plt.plot(x, y, label=label, linestyle=line_style, marker=marker, markersize=8, linewidth=2)


def sort_by_vertices(data):
    return sorted(data, key=lambda x: x['nodes'])


def all_positive(values):
    return all(v is not None and v > 0 for v in values)


def sort_by_edges(data):
    return sorted(data, key=lambda x: x['edges'])


def get_unique_p_values(data_dict):
    unique_p_values = set()
    for data_entries in data_dict.values():
        unique_p_values.update(entry['P'] for entry in data_entries)
    return sorted(unique_p_values)


def group_by_nodes(data):
    groups = {}
    for entry in data:
        v = entry['nodes']
        if v not in groups:
            groups[v] = []
        groups[v].append(entry)
    return groups


def create_gnm_cross_file_comparison(data_dict, subdir):
    decomp_dir = os.path.join(PLOTS_DIRECTORY, subdir, 'decomp')
    os.makedirs(decomp_dir, exist_ok=True)

    all_node_counts = set()
    for data in data_dict.values():
        all_node_counts.update(entry['nodes'] for entry in data)
    all_node_counts = sorted(all_node_counts)

    for v in all_node_counts:
        nodes_data = {}

        for file_name, data in data_dict.items():
            nodes_data[file_name] = [entry for entry in data if entry['nodes'] == v]

        if nodes_data:
            for metric in nodes_data[list(nodes_data.keys())[0]][0].keys():
                if metric not in ['edges', 'nodes', 'runs', 'P']:

                    plt.figure()
                    log_scale_x, log_scale_y = True, True

                    for id, (file_name, group) in enumerate(nodes_data.items()):
                        sorted_group = sort_by_edges(group)
                        x = [entry['edges'] for entry in sorted_group]
                        y = [entry[metric] for entry in sorted_group]

                        if not all_positive(x):
                            log_scale_x = False
                        if not all_positive(y):
                            log_scale_y = False

                        legend_label = LABELS[file_name]
                        plot_with_distinct_styles(x, y, legend_label, id)

                    plt.title(f'|V| = {v}')
                    set_axes_scaling(plt, log_scale_x, log_scale_y)

                    plt.grid(True, which='both', ls='--')
                    plt.xlabel(LABELS['edges'])
                    plt.ylabel(LABELS[metric])
                    plt.legend()

                    output_path = os.path.join(decomp_dir, f'{metric}_vs_edges_{v}.pdf')
                    plt.savefig(output_path)
                    plt.close()
                

def create_base_cross_file_comparison(data_dict, subdir):
    os.makedirs(os.path.join(PLOTS_DIRECTORY, subdir), exist_ok=True)

    if data_dict:
        for metric in data_dict[list(data_dict.keys())[0]][0].keys():
            if metric not in ['edges', 'nodes', 'runs', 'P']:

                plt.figure()
                log_scale_x, log_scale_y = True, True

                for id, (file_name, data) in enumerate(data_dict.items()):
                    sorted_data = sort_by_vertices(data)
                    x = [entry['nodes'] for entry in sorted_data if entry[metric] is not None]
                    y = [entry[metric] for entry in sorted_data if entry[metric] is not None]

                    if not all_positive(x):
                        log_scale_x = False
                    if not all_positive(y):
                        log_scale_y = False

                    legend_label = LABELS[file_name]

                    set_axes_scaling(plt, log_scale_x, log_scale_y)
                    plot_with_distinct_styles(x, y, legend_label, id)

                output_path = os.path.join(PLOTS_DIRECTORY, subdir)
                plot_it(log_scale_x, log_scale_y, metric, output_path, 'nodes')


def create_gnm_plots(data_dict, output_dir):
    grouped_data = group_by_nodes(data_dict)
    
    for metric in data_dict[0].keys():
        if metric not in ['edges', 'nodes', 'runs', 'P']:
            plt.figure()
            log_scale_x, log_scale_y = True, True

            for id, (v, group) in enumerate(grouped_data.items()):
                sorted_group = sort_by_edges(group)
                x = [entry['edges'] for entry in sorted_group if entry['edges'] > 0]
                y = [entry[metric] for entry in sorted_group if entry['edges'] > 0]

                if not all_positive(x):
                    log_scale_x = False
                if not all_positive(y):
                    log_scale_y = False

                plot_with_distinct_styles(x, y, f'|V|={v}', id)

            plot_it(log_scale_x, log_scale_y, metric, output_dir, 'edges')


def create_time_vs_edges_plot(data_dict, output_dir, v=10000):
    filtered_data = [entry for entry in data_dict['h3.log'] if entry['nodes'] == v]

    if not filtered_data:
        print(f"No data found for |V| = {v}")
        return

    output_dir = os.path.join('plots', output_dir, 'preprocessing')
    os.makedirs(output_dir, exist_ok=True)
    print(output_dir)
    
    sorted_data = sort_by_edges(filtered_data)

    x = [entry['edges'] for entry in sorted_data if entry['edges'] > 0]
    time_metrics = ['time_collapse', 'time_topo', 'time_remove_edges', 'time_topo_edges_time']
    y_data = {metric: [entry[metric] for entry in sorted_data if entry['edges'] > 0] for metric in time_metrics}

    plt.figure()

    log_scale_x, log_scale_y = True, True

    for id, metric in enumerate(time_metrics):
        y = y_data[metric]

        if not all_positive(x):
            log_scale_x = False
        if not all_positive(y):
            log_scale_y = False

        plot_with_distinct_styles(x, y, LABELS[metric], id)
    
    plot_it(log_scale_x, log_scale_y, 'time', output_dir, 'edges')


def process_subdir(subdir):
    subsubdir = os.path.join(BASE_DIRECTORY, subdir)
    files = [f for f in os.listdir(subsubdir) if f in TYPES]

    output_subsubdir = None
    data_dict = {}
    for file_name in files:
        file_path = os.path.join(subsubdir, file_name)
        print(f'Processing {file_path}')

        if subdir == 'gnm':
            file_prefix = file_name.split('.')[0]
            output_subsubdir = os.path.join(PLOTS_DIRECTORY, subdir, file_prefix)
            os.makedirs(output_subsubdir, exist_ok=True)

        data = parse_log_file(file_path)
        data_dict[file_name] = data
        
        if subdir == 'gnm':
            create_gnm_plots(data, output_subsubdir)

    return data_dict


def process_gnm_subdir(subdir):
    data_dict = process_subdir(subdir)

    if data_dict:
        create_time_vs_edges_plot(data_dict, subdir)
        create_gnm_cross_file_comparison(data_dict, subdir)


def process_base_subdir(subdir):
    directory = os.path.join(BASE_DIRECTORY, subdir)
    files = [f for f in os.listdir(directory) if os.path.isfile(os.path.join(directory, f))]

    df_dict = {}
    for file_name in files:
        file_path = os.path.join(directory, file_name)

        data = parse_log_file(file_path)
        if data:
            df_dict[file_name] = data

    if df_dict:
        create_base_cross_file_comparison(df_dict, subdir)


def main():
    subdirectories = [d for d in os.listdir(BASE_DIRECTORY) if os.path.isdir(os.path.join(BASE_DIRECTORY, d))]

    for subdir in subdirectories:
        if subdir == 'gnm':
            process_gnm_subdir(subdir)
        else:
            process_base_subdir(subdir)


if __name__ == '__main__':
    main()
