import os
import re
import subprocess

BASE_DIRECTORY = 'data'
RESULTS_DIRECTORY = 'benches'
MAX_TOTAL_TIME_MS = 30000  # 30 seconds
MIN_REPETITIONS = 5
MAX_REPETITIONS = 1000

def parse_times(input_string):
    pattern = (
        r'time-decomp:\s*([\d.]+)\s*ms\s*,\s*'
        r'time-preprocess:\s*([\d.]+)\s*ms\s*,\s*'
        r'time-scheme:\s*([\d.]+)\s*ms\s*,\s*'
        r'time-reading:\s*([\d.]+)\s*ms\s*,\s*'
        r'time-comp:\s*([\d.]+)\s*ms\s*,\s*'
        r'time-total:\s*([\d.]+)\s*ms'
    )
    match = re.search(pattern, input_string)
    
    if match:
        return [float(match.group(i)) for i in range(1, 7)]
    else:
        return []


def replace_times(input_string, new_times):
    pattern = (
        r'(time-decomp:\s*)([\d.]+)(\s*ms\s*,\s*)'
        r'(time-preprocess:\s*)([\d.]+)(\s*ms\s*,\s*)'
        r'(time-scheme:\s*)([\d.]+)(\s*ms\s*,\s*)'
        r'(time-reading:\s*)([\d.]+)(\s*ms\s*,\s*)'
        r'(time-comp:\s*)([\d.]+)(\s*ms\s*,\s*)'
        r'(time-total:\s*)([\d.]+)(\s*ms)'
    )
    replacement = (
        f'time-decomp: {new_times[0]:.4f} ms, '
        f'time-preprocess: {new_times[1]:.4f} ms, '
        f'time-scheme: {new_times[2]:.4f} ms, '
        f'time-reading: {new_times[3]:.4f} ms, '
        f'time-comp: {new_times[4]:.4f} ms, '
        f'time-total: {new_times[5]:.4f} ms'
    )
    updated_string = re.sub(pattern, replacement, input_string)
    return updated_string


def run_bench(option, file_path):
    try:
        result = subprocess.run(['/usr/bin/time', '-v', './fruit'] + option.split() + [file_path],
                                stderr=subprocess.PIPE, stdout=subprocess.PIPE, text=True, timeout=300)
        if result.returncode == 0:
            return result
        else:
            print(f'Error: Call failed with return code {result.returncode}')
            return None
    except subprocess.TimeoutExpired:
        print(f'Error: Call timed out for {file_path}')
        return None
    except Exception as e:
        print(f'Error: An exception occurred: {e}')
        return None


def get_mem_usage(result):
    memory_usage = 0
    try:
        if result.stderr:
            for line in result.stderr.splitlines():
                if 'Maximum resident set size' in line:
                    memory_usage = int(line.split()[-1])
                    break
    except Exception as e:
        print(f'Error: Reading memory usage: {e}')
    return memory_usage


def estimate_num_runs(elapsed_time_ms):
    if elapsed_time_ms > 0:
        num_reps = max(int(MAX_TOTAL_TIME_MS / elapsed_time_ms), MIN_REPETITIONS)
        num_reps = min(num_reps, MAX_REPETITIONS)
    else:
        num_reps = MIN_REPETITIONS
    print(f'Estimated number of runs: {num_reps} based on the first run time of {elapsed_time_ms:.2f} ms')
    return num_reps


def run_with_options(option, sub_directory, log_file):
    for file_name in os.listdir(sub_directory):
        file_path = os.path.join(sub_directory, file_name)
        if os.path.isfile(file_path):
            print(f'Processing {file_path} {option}...')

            cumulative_memory = 0.0
            cumulative_times = [0.0, 0.0, 0.0, 0.0, 0.0, 0.0]
            output = None
            i = 0

            # Initial run
            result = run_bench(option, file_path)
            memory_usage = get_mem_usage(result) if result else 0.0
            if result is None or not result.stdout.strip():
                print('Result is not valid in initial run!')
                with open(log_file, 'a') as log:
                    log.write(
                        f'{file_name}: time-total: > 5min , memory: {memory_usage:.2f} KB , runs: 0\n'
                    )
                continue
          
            output = result.stdout.strip()

            time_values = parse_times(output)
            if time_values:
                for j in range(6):
                    cumulative_times[j] += time_values[j]

            num_reps = estimate_num_runs(cumulative_times[5])
            cumulative_memory += memory_usage

            # Average runs
            for i in range(1, num_reps):
                if cumulative_times[5] >= MAX_TOTAL_TIME_MS:
                    print(
                        f'Stopping further runs for {file_path} after {i} repetitions '
                        f'(total time exceeded).'
                    )
                    break

                result = run_bench(option, file_path)
                memory_usage = get_mem_usage(result) if result else 0.0
                if result is None or not result.stdout.strip():
                    print('Result is not valid in average runs!')
                    i -= 1
                    break

                cumulative_memory += memory_usage

                output = result.stdout.strip()

                # Updating cumulative times
                time_values = parse_times(output)
                if time_values:
                    for j in range(6):
                        cumulative_times[j] += time_values[j]

            # Calculate averages
            total_reps = i + 1
            average_memory = cumulative_memory / total_reps
            average_times = [time / total_reps for time in cumulative_times]

            updated_output = replace_times(output, average_times)

            with open(log_file, 'a') as log:
                log.write(
                    f'{file_name}: {updated_output} , '
                    f'memory: {average_memory:.2f} KB , runs: {total_reps}\n'
                )


def process_sub_directory(sub_directory):
    sub_folder_name = os.path.basename(sub_directory)
    result_dir = os.path.join(RESULTS_DIRECTORY, sub_folder_name)
    os.makedirs(result_dir, exist_ok=True)

    log_files = {
        '-b': os.path.join(result_dir, 'h3.log'),
        '-b -no': os.path.join(result_dir, 'no.log'),
        '-b -co': os.path.join(result_dir, 'co.log'),
        '-b -coc': os.path.join(result_dir, 'coc.log'),
        '-b -noc': os.path.join(result_dir, 'noc.log'),
    }

    for option, log_file in log_files.items():
        open(log_file, 'w').close()
        run_with_options(option, sub_directory, log_file)


def main():
    for sub_directory in os.listdir(BASE_DIRECTORY):
        full_sub_directory_path = os.path.join(BASE_DIRECTORY, sub_directory)
        if os.path.isdir(full_sub_directory_path):
            process_sub_directory(full_sub_directory_path)


if __name__ == '__main__':
    main()
