import subprocess
import math
import os

def main():
    os.makedirs('./data/er', exist_ok=True)
    os.makedirs('./data/gnm', exist_ok=True)
    os.makedirs('./data/gn', exist_ok=True)
    os.makedirs('./data/gnc', exist_ok=True)
    os.makedirs('./data/sf', exist_ok=True)
    command = ['python3', './scripts/generate_graph.py', 'FLAG']

    # Create GNM-Models
    print('Creating GNM-Graphs...')
    command[2] = 'gnm'

    for i in range(0, 7):
        curr_n = pow(10, i)
        j = 0
        curr_max = math.floor((curr_n * (curr_n-1)) / 2)
        while pow(10, j) < curr_max:
            curr_m = pow(10, j)
            path = f'./data/gnm/gnm_{str(curr_n)}_{str(curr_m)}.gr'
            print("creating: " + path)
            subprocess.run(command + [str(curr_n), str(curr_m), path], text=True)
            j += 1
        path = f'./data/gnm/gnm_{str(curr_n)}_{str(curr_max)}.gr'
        print("creating: " + path)
        subprocess.run(command + [str(curr_n), str(curr_max), path], text=True)

    # Create GN-Models
    print('Creating GN-Graphs...')
    command[2] = 'gn'

    for i in range(1, 6):
        curr_n = str(pow(10, i))
        path = f'./data/gn/gn_{curr_n}.gr'
        print("creating: " + path)
        subprocess.run(command + [curr_n, path], text=True)

    # Create GNC-Models
    print('Creating GNC-Graphs...')
    command[2] = 'gnc'

    for i in range(1, 8):
        curr_n = str(pow(10, i))
        path = f'./data/gnc/gnc_{curr_n}.gr'
        print("creating: " + path)
        subprocess.run(command + [curr_n, path], text=True)

    # Create SF-Models
    print('Creating SF-Graphs...')
    command[2] = 'sf'

    for i in range(1, 9):
        curr_n = str(pow(10, i))
        path = f'./data/sf/sf_{curr_n}.gr'
        print("creating: " + path)
        subprocess.run(command + [curr_n, path], text=True)


if __name__ == '__main__':
    main()
