BASE_DIRECTORY = 'benches'
MEMORY_CONVERSION_FACTOR = 1e3  # KB to MB

TYPES = ['no.log', 'co.log', 'noc.log', 'coc.log', 'h3.log']

LABELS = {
    'nodes': '|V|',
    'edges': '|E|',
    'scc': '|SCC|',
    'chains': '$k_c$',
    'scheme_size': 'Scheme-Size',
    'removed_edges': '|Removed Edges|',
    'collapse_nodes': '|Vertices Processed during Collapse|',
    'collapse_edges': '|Edges Processed during Collapse|',
    'decomp_nodes': '|Vertices Processed during Decomposition|',
    'decomp_edges': '|Edges Processed during Decomposition|',
    'scheme_nodes': '|Vertices Processed during Scheme|',
    'scheme_edges': '|Edges Processed during Scheme|',
    'time_decomp': 'Decomposition Time (ms)',
    'time_preprocess': 'Preprocessing Time (ms)',
    'time_scheme': 'Scheme Time (ms)',
    'time_reading': 'Reading Time (ms)',
    'time_comp': 'Computation Time (ms)',
    'time_total': 'Total Time (ms)',
    'time_collapse': 'Collapse to DAG Time (ms)',
    'time_topo': 'Topological Sort Time (ms)',
    'time_remove_edges': 'Remove Transitive Edges Time (ms)',
    'time_topo_edges_time': 'Topological Sort Edges Time (ms)',
    'memory': 'Memory Usage (MB)',
    'no.log': 'Node-Order',
    'co.log': 'Chain-Order',
    'noc.log': 'Node-Order-Concat',
    'coc.log': 'Chain-Order-Concat',
    'h3.log': 'H3-Concat',
    'time': 'Time (ms)',
}


def parse_log_file(filepath):
    data = []
    with open(filepath, 'r') as file:
        print('Parsing ' + filepath)
        for line in file:
            if 'time-total: > 5min' in line:
                entry_name = line.split(': ')[0]
                print('Ignoring faulty run: ' + entry_name)
                continue
            
            parts = line.split(', ')
            entry_name = parts[0].split(': ')[0]

            parsed_data = {
                'nodes': int(parts[0].split('#nodes: ')[1].split()[0]),
                'edges': int(parts[1].split(': ')[1]),
                'scc': int(parts[2].split(': ')[1]),
                'chains': int(parts[3].split(': ')[1]),
                'scheme_size': int(parts[4].split(': ')[1]),
                'removed_edges': int(parts[5].split(': ')[1]),
                'collapse_nodes': int(parts[6].split(': ')[1]),
                'collapse_edges': int(parts[7].split(': ')[1]),
                'decomp_nodes': int(parts[8].split(': ')[1]),
                'decomp_edges': int(parts[9].split(': ')[1]),
                'scheme_nodes': int(parts[10].split(': ')[1]),
                'scheme_edges': int(parts[11].split(': ')[1]),
                'time_decomp': float(parts[12].split(': ')[1].split()[0]),
                'time_preprocess': float(parts[13].split(': ')[1].split()[0]),
                'time_scheme': float(parts[14].split(': ')[1].split()[0]),
                'time_reading': float(parts[15].split(': ')[1].split()[0]),
                'time_comp': float(parts[16].split(': ')[1].split()[0]),
                'time_total': float(parts[17].split(': ')[1].split()[0]),
                'time_collapse': float(parts[18].split(': ')[1].split()[0]),
                'time_topo': float(parts[19].split(': ')[1].split()[0]),
                'time_remove_edges': float(parts[20].split(': ')[1].split()[0]),
                'time_topo_edges_time': float(parts[21].split(': ')[1].split()[0]),
                'memory': float(parts[22].split(': ')[1].split()[0]) / MEMORY_CONVERSION_FACTOR,
                'runs': int(parts[23].split(': ')[1]),
            }
            
            if parsed_data['nodes'] < 2:
                continue

            data.append(parsed_data)

    return data
