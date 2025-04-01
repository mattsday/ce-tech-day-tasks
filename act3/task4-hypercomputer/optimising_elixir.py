import random
import time
import math
import argparse  # To accept command-line arguments for intensity

# --- Configuration for the Absurd Elixir ---

# Youthful pepper varieties
PEPPER_VARIETIES = [
    "carolina_reaper",  # Ironically good at reversing age
    "bell_pepper",  # For the cautious youth-seeker
    "dragons_breath", "pepper_x",  # Real contenders for hottest
    "shishito",  # Mostly mild, sometimes surprising!
    "quantum_quandary_pepper",  # Fictional, tastes like uncertainty
    "crybaby_serrano",  # Extra potent?
    "tickle_me_pink_pepper"  # Fictional, sounds harmless...
]

# What vessel holds the fermenting magic?
FERMENTATION_VESSELS = [
    "oak_barrel",  # Steady and solid
    "hollowed_out_moon_rock",  # Premium choice
    "discarded_welly_boot",  # Questionable hygiene, adds 'character'
    "the_holy_grail"  # If you can find it
]

# Does background music influence the elixir's potency? Of course!
AMBIENT_MUSICS = [
    "MF Doom"  # Mm Food
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

# --- Simulation Function ---


def simulate_youthfulness_intensive(
    pepper_variety,
    fermentation_time,
    vessel,
    music,
    phase,
    base_iterations=100
):
    """
    Simulates the youthfulness score with CPU load, now considering
    even more ludicrously contrived factors.
    """
    # --- Determine calculation intensity based on inputs ---

    pepper_factors = {
        "habanero": 1.2, "jalapeno": 0.8, "ghost_pepper": 1.5, "carolina_reaper": 1.8,
        "bell_pepper": 0.5, "dragons_breath": 1.9, "pepper_x": 2.0, "shishito": 0.9,
        "quantum_quandary_pepper": 1.3, "crybaby_serrano": 1.1, "tickle_me_pink_pepper": 0.7,
        "default": 1.0  # Fallback
    }
    pepper_factor = pepper_factors.get(
        pepper_variety, pepper_factors["default"])

    vessel_factors = {
        "oak_barrel": 1.1, "stainless_steel": 1.0, "clay_pot": 0.95, "glass_carboy": 1.05,
        "hollowed_out_moon_rock": 1.6, "discarded_welly_boot": 0.7, "the_holy_grail": 1.4,
        "default": 1.0
    }
    vessel_factor = vessel_factors.get(vessel, vessel_factors["default"])

    music_factors = {
        "classical_bach": 1.0, "heavy_metal_slayer": 1.3, "whale_songs": 0.85,
        # Reduces potency?
        "uplifting_polka": 1.1, "smooth_jazz": 0.9, "barry_manilow_collection": 0.6,
        "absolute_silence": 1.05, "ambient_static_from_jupiter": 1.2,
        "default": 1.0
    }
    music_factor = music_factors.get(music, music_factors["default"])

    phase_factors = {
        "new_moon": 1.0, "waxing_crescent": 1.05, "first_quarter": 1.1, "waxing_gibbous": 1.15,
        "full_moon": 1.25, "waning_gibbous": 0.95, "last_quarter": 0.9, "waning_crescent": 0.85,
        "tuesday": 1.11,  # Tuesdays are special
        "default": 1.0
    }
    phase_factor = phase_factors.get(phase, phase_factors["default"])

    # Adjust iterations based on fermentation time (non-linear scaling)
    time_factor = 1 + \
        math.sqrt(max(0, fermentation_time - 12)) / math.sqrt(72 - 12)

    # Combine all factors to get total iterations
    total_iterations = int(
        base_iterations * pepper_factor * time_factor *
        vessel_factor * music_factor * phase_factor
    )
    # Ensure a minimum number of iterations to prevent zero division later
    total_iterations = max(100, total_iterations)

    # --- Perform CPU-intensive work ---
    result = 0.0
    for i in range(1, total_iterations + 1):
        # Mix of operations using factors for slight variation per combo
        val = math.sin(i * 0.01 * pepper_factor) * \
            math.cos(i * 0.005 * vessel_factor)
        val += math.sqrt(abs(math.log(i) * time_factor *
                         music_factor if i > 1 else 1.0))
        val = (val / (i % 100 + 1)) * phase_factor  # Incorporate phase factor
        result += val

    # --- Derive a 'score' from the result ---
    # Use modulo for pseudo-randomness based on calculation
    pseudo_random_base = (result % 60.0)  # Value between -60 and 60

    # Add small score adjustments based on well-tested factors (before clamping)
    score_offset = 0
    if vessel == "hollowed_out_moon_rock":
        score_offset += 5
    if vessel == "discarded_welly_boot":
        score_offset -= 10  # Penalty!
    if music == "barry_manilow_collection":
        score_offset -= 5
    if pepper_variety == "tickle_me_pink_pepper":
        score_offset += 3
    if phase == "tuesday":
        score_offset += 2  # Tuesday bonus

    # Add minor random noise
    noise = random.uniform(-2, 2)

    # Combine to get final score
    final_score = 40 + pseudo_random_base + score_offset + noise
    final_score = max(0, min(100, final_score))  # Clamp score to 0-100 range

    return final_score

# --- Optimization Function ---


def find_optimal_elixir(iterations_per_simulation):
    """
    Finds the optimal combination across all contrived dimensions using the
    CPU-intensive simulation function.
    """
    fermentation_times = range(
        12, 72, 12)  # Reduced time steps slightly to compensate for other dimensions (12, 24, 36, 48, 60)

    best_score = -float('inf')
    best_params = {}

    # Calculate total simulations - WARNING: This can get VERY large!
    total_simulations = (
        len(PEPPER_VARIETIES) * len(fermentation_times) *
        len(FERMENTATION_VESSELS) * len(AMBIENT_MUSICS) * len(LUNAR_PHASES)
    )
    print("Starting The Grand Elixir Optimization Quest!")
    print(
        f"Searching across {len(PEPPER_VARIETIES)} peppers, {len(fermentation_times)} times,")
    print(f"{len(FERMENTATION_VESSELS)} vessels, {len(AMBIENT_MUSICS)} musics, {len(LUNAR_PHASES)} phases.")
    print(f"Total simulations to run: {total_simulations:,}")
    print(
        f"Base iterations per simulation step: {iterations_per_simulation:,}")
    if total_simulations > 100000:
        print("\nWARNING: This is a *large* number of simulations and may take a very long time!")
        print("Perfect for stress-testing Slurm, perhaps less so for quick results.\n")
    print("-" * 60)

    count = 0
    # The MEGA-LOOP - iterating through all dimensions
    for variety in PEPPER_VARIETIES:
        for time_ in fermentation_times:
            for vessel in FERMENTATION_VESSELS:
                for music in AMBIENT_MUSICS:
                    for phase in LUNAR_PHASES:
                        count += 1
                        sim_start_time = time.time()

                        # Call the intensive simulation
                        score = simulate_youthfulness_intensive(
                            variety, time_, vessel, music, phase, iterations_per_simulation
                        )

                        sim_end_time = time.time()
                        duration = sim_end_time - sim_start_time

                        # Print progress (limit frequency if too many simulations?)
                        # Only print every N steps if it's overwhelming
                        # Print every 100 steps, or if slow, or last step
                        if count % 100 == 0 or duration > 1.0 or count == total_simulations:
                            print(f"[{count:>{len(str(total_simulations))}}/{total_simulations:,}] "
                                  f"Time: {duration:5.2f}s | Score: {score:6.2f} | "
                                  f"{variety[:10]:<10}, {time_}h, {vessel[:10]:<10}, {music[:10]:<10}, {phase[:10]:<10}")

                        if score > best_score:
                            best_score = score
                            best_params = {
                                "variety": variety,
                                "time": time_,
                                "vessel": vessel,
                                "music": music,
                                "phase": phase,
                            }
                            print(
                                f"*** New Best Score Found: {best_score:.2f} with params: {best_params} ***")

    print("-" * 60)
    return best_params, best_score

# --- Main Execution ---


if __name__ == "__main__":
    # Set up argument parsing
    parser = argparse.ArgumentParser(
        description="Find the optimal youth elixir via a ludicrously complex and CPU-intensive simulation.",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter  # Show defaults in help
    )
    parser.add_argument(
        '-i', '--iterations',
        type=int,
        # Lowered default due to vastly increased simulation count! Adjust as needed.
        default=10000,
        help='Base number of iterations per *single* simulation step. Higher values increase CPU load.'
    )
    args = parser.parse_args()

    print(
        f"Initializing simulation with base iterations = {args.iterations:,}")
    overall_start_time = time.time()

    # Run the optimization
    optimal_params, optimal_score = find_optimal_elixir(args.iterations)

    overall_end_time = time.time()

    # Print the final results
    print("\n--- The Ultimate Elixir Formula Has Been Revealed! ---")
    if optimal_params:
        print(f"Best Pepper Variety : {optimal_params['variety']}")
        print(f"Best Fermentation Time: {optimal_params['time']} hours")
        print(f"Best Vessel         : {optimal_params['vessel']}")
        print(f"Best Ambient Music  : {optimal_params['music']}")
        print(f"Best Lunar Phase    : {optimal_params['phase']}")
        print(f"Achieved Youth Score: {optimal_score:.2f}")
    else:
        print(
            "No optimal solution found (this shouldn't happen unless no simulations ran).")

    print(
        f"\nTotal Optimization Quest Duration: {overall_end_time - overall_start_time:.2f} seconds")

    # Context for fun
    print("\nSimulation completed near: Solar System, Milky Way, Sector ZZ9 Plural Z Alpha, Western Spiral Arm of the Galaxy")
    # Uses system time/zone
    print(
        f"Approximate real time was: {time.strftime('%Y-%m-%d %H:%M:%S %Z')}")
