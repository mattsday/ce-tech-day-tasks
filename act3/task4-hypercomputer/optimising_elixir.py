import random
import time

def simulate_youthfulness(pepper_variety, fermentation_time):
    """Simulates the youthfulness score based on pepper variety and fermentation time."""

    # Simplified simulation logic (replace with your actual data/model)
    base_score = 50  # Base youthfulness score
    pepper_effect = random.uniform(-10, 20) if pepper_variety == "habanero" else random.uniform(-5, 10)
    fermentation_effect = random.uniform(-15, 25) if 24 <= fermentation_time <= 48 else random.uniform(-10, 10)

    # Adding in some random noise to make it more realistic.
    noise = random.uniform(-5, 5)

    return base_score + pepper_effect + fermentation_effect + noise

def find_optimal_elixir():
    """Finds the optimal pepper variety and fermentation time."""

    pepper_varieties = ["habanero", "jalapeno", "ghost_pepper"]
    fermentation_times = range(12, 72, 6)  # Fermentation times in hours

    best_score = -float('inf')
    best_variety = None
    best_time = None

    for variety in pepper_varieties:
        for time_ in fermentation_times:
            score = simulate_youthfulness(variety, time_)
            print(f"Variety: {variety}, Time: {time_} hours, Score: {score}") #Print the results to the Slurm output.
            if score > best_score:
                best_score = score
                best_variety = variety
                best_time = time_
    return best_variety, best_time, best_score

if __name__ == "__main__":
    start_time = time.time()
    optimal_variety, optimal_time, optimal_score = find_optimal_elixir()
    end_time = time.time()

    print("\n--- Optimal Elixir Formula ---")
    print(f"Pepper Variety: {optimal_variety}")
    print(f"Fermentation Time: {optimal_time} hours")
    print(f"Youthfulness Score: {optimal_score}")
    print(f"Time Taken: {end_time - start_time:.2f} seconds")