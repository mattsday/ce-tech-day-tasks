import random
import time
import math
import argparse
import os
import json  # For saving partial results
import itertools  # For easily generating combinations

# Youthful pepper varieties
PEPPER_VARIETIES = [
    "habanero", "jalapeno", "ghost_pepper",  # Classics
    "carolina_reaper",  # Ironically good at reversing age
    "bell_pepper",  # For the cautious youth-seeker
    "dragons_breath", "pepper_x",  # Real contenders for hottest
    "shishito",  # Mostly mild, sometimes surprising!
    "quantum_quandary_pepper",  # Tastes like uncertainty
    "crybaby_serrano",  # Extra potent?
    "tickle_me_pink_pepper"  # Sounds harmless...
]

# What vessel holds the fermenting magic?
FERMENTATION_VESSELS = [
    "stainless_steel", "clay_pot", "glass_carboy",  # Mainstream
    "oak_barrel",  # Steady and solid
    "hollowed_out_moon_rock",  # Premium choice
    "discarded_welly_boot",  # Questionable hygiene, adds 'character'
    "the_holy_grail"  # If you can find it
]

# Does background music influence the elixir's potency? Of course!
AMBIENT_MUSICS = [
    "mf_doom",  # Mm Food
    "anthony_kiedis", # Funky vibes
    "heavy_metal_slayer",  # YOLO
    "barry_manilow_collection",  # Risky!
    "absolute_silence",  # For the purists
    "ambient_static_from_jupiter"  # Good vibes only
]

# The cosmos must surely play a part!
LUNAR_PHASES = [
    "full_moon",  # Classic spooky choice
    "the_waterboys",  # The whole of the moon
    "tuesday"  # A very specific cosmic alignment
]

FERMENTATION_TIMES = range(12, 72, 12)  # 12, 24, 36, 48, 60

# --- Simulation Function (Keep as before) ---


def simulate_youthfulness_intensive(
    pepper_variety, fermentation_time, vessel, music, phase, base_iterations=100000
):
    """
    Simulates the youthfulness score with INCREASED CPU load per iteration,
    considering all the ludicrously contrived factors.
    """
    # --- Determine calculation intensity factors (Same as before) ---
    pepper_factors = {
        "habanero": 1.2, "jalapeno": 0.8, "ghost_pepper": 1.5, "carolina_reaper": 1.8,
        "bell_pepper": 0.5, "dragons_breath": 1.9, "pepper_x": 2.0, "shishito": 0.9,
        "quantum_quandary_pepper": 1.3, "crybaby_serrano": 1.1, "tickle_me_pink_pepper": 0.7,
        "default": 1.0
    }
    pepper_factor = pepper_factors.get(pepper_variety, pepper_factors["default"])
    vessel_factors = {
        "oak_barrel": 1.1, "stainless_steel": 1.0, "clay_pot": 0.95, "glass_carboy": 1.05,
        "hollowed_out_moon_rock": 1.6, "discarded_welly_boot": 0.7, "the_holy_grail": 1.4,
        "default": 1.0
    }
    vessel_factor = vessel_factors.get(vessel, vessel_factors["default"])
    music_factors = {
        "classical_bach": 1.0, "heavy_metal_slayer": 1.3, "whale_songs": 0.85,
        "uplifting_polka": 1.1, "smooth_jazz": 0.9, "barry_manilow_collection": 0.6,
        "absolute_silence": 1.05, "ambient_static_from_jupiter": 1.2,
        "default": 1.0
    }
    music_factor = music_factors.get(music, music_factors["default"])
    phase_factors = {
        "new_moon": 1.0, "waxing_crescent": 1.05, "first_quarter": 1.1, "waxing_gibbous": 1.15,
        "full_moon": 1.25, "waning_gibbous": 0.95, "last_quarter": 0.9, "waning_crescent": 0.85,
        "tuesday": 1.11,
        "default": 1.0
    }
    phase_factor = phase_factors.get(phase, phase_factors["default"])
    time_factor = 1 + math.sqrt(max(0, fermentation_time - 12)) / math.sqrt(72 - 12)
    total_iterations = int(base_iterations * pepper_factor * time_factor * vessel_factor * music_factor * phase_factor)
    total_iterations = max(100, total_iterations) # Ensure minimum iterations

    # --- Perform HEAVIER CPU-intensive work ---
    result = 0.0
    for i in range(1, total_iterations + 1): # Main loop

        # --- Step 1: Original calculations (keep variation) ---
        trig_val = math.sin(i * 0.01 * pepper_factor) * math.cos(i * 0.005 * vessel_factor)
        log_sqrt_val = math.sqrt(abs(math.log(i * time_factor + 1) * music_factor)) # Avoid log(0) or log(negative)
        base_val = (trig_val + log_sqrt_val / (i % 100 + 1)) * phase_factor

        # --- Step 2: Add more complex math - Exponentiation ---
        # Use pow() or **, ensure base/exponent vary slightly and don't explode too fast
        exp_base = 1.001 + abs(math.sin(i * 0.0001 * pepper_factor)) * 0.005 # Keep base slightly > 1
        exp_exponent = 5.0 + (i % 5) + vessel_factor # Small exponent, varies
        try:
            # Power calculation can be costly
            exp_val = pow(exp_base, exp_exponent)
        except OverflowError:
            exp_val = float('inf') # Handle potential overflow if numbers get huge

        # --- Step 3: Add a small inner loop performing calculations ---
        # This forces repeated work within each main iteration
        inner_sum = 0
        # Scale inner loop size modestly, e.g., based on i or a factor
        inner_loop_limit = 5 + (i % 15) # Small inner loop (5 to 19 reps)
        for j in range(inner_loop_limit):
             # Use another math function, combine factors in different ways
             inner_term = math.atan(base_val * 0.1 * j * music_factor)
             inner_term += math.sinh(pepper_factor * 0.01 * j) # Hyperbolic sine
             inner_sum += inner_term / (j + 1)

        # --- Step 4: Combine results ---
        # Add the results of the new steps. Scale if necessary to prevent huge numbers.
        # The exact combination isn't important, just that work is done.
        result += base_val + (exp_val * 0.0001) + (inner_sum * 0.1) # Scale down exp/inner results

        # --- Step 5: Optional - Add yet another costly operation ---
        # Example: Factorial (computationally expensive but grows VERY fast)
        # Use sparingly and only for small numbers derived from i
        # if i % 50 == 0: # Only do this occasionally
        #    small_num = i % 7 + 2 # Calculate factorial of numbers between 2 and 8
        #    try:
        #        result += math.factorial(small_num) * 0.00001 # Scale heavily
        #    except ValueError: # Handles negative input, though shouldn't happen here
        #        pass
        # --> Decided against factorial for now as it grows too fast and might dominate.

    # --- Derive a 'score' (Same as before) ---
    pseudo_random_base = (result % 60.0)
    score_offset = 0
    if vessel == "hollowed_out_moon_rock": score_offset += 5
    if vessel == "discarded_welly_boot": score_offset -= 10
    if music == "barry_manilow_collection": score_offset -= 5
    if pepper_variety == "tickle_me_pink_pepper": score_offset += 3
    if phase == "tuesday": score_offset += 2
    noise = random.uniform(-2, 2)
    final_score = 40 + pseudo_random_base + score_offset + noise
    final_score = max(0, min(100, final_score))

    return final_score
# --- Main Logic ---


def run_task(task_id, task_count, iterations_per_simulation, output_dir):
    """Runs the simulations assigned to a specific task."""

    # 1. Generate all possible combinations
    all_combinations = list(itertools.product(
        PEPPER_VARIETIES,
        FERMENTATION_TIMES,
        FERMENTATION_VESSELS,
        AMBIENT_MUSICS,
        LUNAR_PHASES
    ))
    total_combinations = len(all_combinations)
    print(
        f"[Task {task_id}/{task_count}] Total combinations in search space: {total_combinations:,}")

    local_best_score = -float('inf')
    local_best_params = None
    combinations_this_task = 0

    start_time = time.time()

    # 2. Iterate through combinations and select subset for this task
    for index, combo in enumerate(all_combinations):
        if index % task_count == (task_id - 1):
            combinations_this_task += 1
            variety, time_, vessel, music, phase = combo

            # Run the simulation for this combination
            score = simulate_youthfulness_intensive(
                variety, time_, vessel, music, phase, iterations_per_simulation
            )

            # Update local best if this one is better
            if score > local_best_score:
                local_best_score = score
                local_best_params = {
                    "variety": variety, "time": time_, "vessel": vessel,
                    "music": music, "phase": phase
                }
                # Optional: Print progress when local best is found
                # print(f"[Task {task_id}] New local best score: {score:.2f} (Combo index {index})")

    end_time = time.time()
    print(f"[Task {task_id}/{task_count}] Processed {combinations_this_task} combinations in {end_time - start_time:.2f} seconds.")

    # 3. Save the local best result to a unique file
    output_filename = os.path.join(
        output_dir, f"partial_results_{task_id}.json")
    result_data = {
        "task_id": task_id,
        "task_count": task_count,
        # Handle case where no combos were assigned/run
        "best_score": local_best_score if local_best_params else None,
        "best_params": local_best_params,
        "combinations_processed": combinations_this_task
    }
    try:
        with open(output_filename, 'w') as f:
            json.dump(result_data, f, indent=2)
        print(
            f"[Task {task_id}/{task_count}] Local best result saved to {output_filename}")
    except Exception as e:
        print(f"[Task {task_id}/{task_count}] ERROR saving results: {e}")


def collect_results(output_dir):
    """Reads all partial result files and finds the global best."""
    print(f"\n--- Collecting Results from Directory: {output_dir} ---")
    all_results = []
    try:
        result_files = [f for f in os.listdir(output_dir) if f.startswith(
            "partial_results_") and f.endswith(".json")]
        if not result_files:
            print("ERROR: No partial result files found.")
            return

        for filename in result_files:
            filepath = os.path.join(output_dir, filename)
            try:
                with open(filepath, 'r') as f:
                    data = json.load(f)
                    # Only consider tasks that found a result
                    if data.get("best_score") is not None:
                        all_results.append(data)
            except Exception as e:
                print(f"Warning: Could not read or parse {filename}: {e}")

    except FileNotFoundError:
        print(f"ERROR: Output directory '{output_dir}' not found.")
        return
    except Exception as e:
        print(f"ERROR reading results directory: {e}")
        return

    if not all_results:
        print("No valid results found in partial files.")
        return

    # Find the best score among all tasks
    global_best_score = -float('inf')
    global_best_result = None

    for result in all_results:
        if result["best_score"] > global_best_score:
            global_best_score = result["best_score"]
            global_best_result = result

    print("\n--- The Ultimate Elixir Formula Has Been Revealed! ---")
    if global_best_result:
        params = global_best_result['best_params']
        print(f"Best Pepper Variety : {params['variety']}")
        print(f"Best Fermentation Time: {params['time']} hours")
        print(f"Best Vessel         : {params['vessel']}")
        print(f"Best Ambient Music  : {params['music']}")
        print(f"Best Lunar Phase    : {params['phase']}")
        print(f"Achieved Youth Score: {global_best_result['best_score']:.2f}")
        print(
            f"(Found by Task {global_best_result['task_id']} out of {global_best_result['task_count']})")
    else:
        print("No optimal solution found across all tasks.")

    total_processed = sum(r.get("combinations_processed", 0)
                          for r in all_results)
    print(
        f"\nTotal combinations processed across all tasks: {total_processed:,}")
    # Note: Total time requires looking at Slurm accounting or start/end of the whole array job.


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Run or collect results for a parallelized, CPU-intensive elixir quest.",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter
    )
    parser.add_argument(
        'mode', choices=['run', 'collect'],
        help="Mode: 'run' to execute a task's portion, 'collect' to gather results."
    )
    parser.add_argument(
        # VERY LOW default for 10min limit!
        '-i', '--iterations', type=int, default=1000000,
        help='Base iterations per simulation. CRITICAL for fitting in time limit.'
    )
    parser.add_argument(
        '-o', '--output-dir', type=str, default="results",
        help='Directory to store/read partial results.'
    )
    args = parser.parse_args()

    # Ensure output directory exists for 'run' mode
    if args.mode == 'run' and not os.path.exists(args.output_dir):
        try:
            os.makedirs(args.output_dir)
            print(f"Created output directory: {args.output_dir}")
        except OSError as e:
            print(f"Error creating output directory {args.output_dir}: {e}")
            exit(1)  # Exit if we can't create the dir

    if args.mode == 'run':
        # Get Slurm variables (provide defaults for local testing)
        task_id = int(os.environ.get('SLURM_ARRAY_TASK_ID', 1))
        task_count = int(os.environ.get('SLURM_ARRAY_TASK_COUNT', 1))
        print(f"Starting Elixir Quest Task {task_id}/{task_count}")
        print(f"Base iterations: {args.iterations:,}")
        print(f"Output directory: {args.output_dir}")
        run_task(task_id, task_count, args.iterations, args.output_dir)
        print(f"Finished Elixir Quest Task {task_id}/{task_count}")
    elif args.mode == 'collect':
        collect_results(args.output_dir)
