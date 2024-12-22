import sys
import networkx as nx

def write_graph(graph, file):
    with open(file, 'w') as f:
        f.write(f'n: {graph.number_of_nodes()}\n')
        for edge in graph.edges():
            f.write(f'{edge[0]} {edge[1]}\n')


def read_input():
    if len(sys.argv) < 4:
        print('Usage: python3 generate_graph.py <flag> <n> [additional-params] <target-file>')
        sys.exit(1)

    flag = sys.argv[1]
    
    try:
        n = int(sys.argv[2])
    except (ValueError, IndexError):
            print('Error: <n> must be an integer.')
            sys.exit(1)

    try:
        target_file_name = sys.argv[-1]
    except (ValueError, IndexError):
            print('Error reading target file.')
            sys.exit(1)

    params = sys.argv[3:-1]
    
    return flag, n, target_file_name, params


def generate_graph(flag, n, params):
    if flag == 'gnm':
        try:
            m = int(params[0])
        except (ValueError, IndexError):
            print('Error: <m> must be an integer for gnm_random_graph.')
            sys.exit(1)
        digraph = nx.gnm_random_graph(n, m, directed=True)
    
    elif flag == 'gn':
        digraph = nx.gn_graph(n)

    elif flag == 'gnc':
        digraph = nx.gnc_graph(n)

    elif flag == 'sf':
        digraph = nx.scale_free_graph(n)
    else:
        print('Invalid flag.')
        sys.exit(1)
    
    return digraph


def main():
    flag, n, target_file_name, params = read_input()
    digraph = generate_graph(flag, n, params)
    
    write_graph(digraph, target_file_name)


if __name__ == '__main__':
    main()
