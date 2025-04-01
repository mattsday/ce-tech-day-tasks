<<<<<<< HEAD
## The Situation

Cymbal Supplements' incredible age-reversing hot sauce is flying off the shelves! To meet this unprecedented demand, we need to build a brand-spanking-new factory. But before we break ground, we need to be smart about where we build. Nobody wants a hot sauce factory that turns into a *watery* sauce factory!

Your mission, should you choose to accept it, is to dive into the world of geospatial data and identify areas with the *highest* flood risk. We've already loaded some juicy Earth Engine flood risk data into BigQuery, specifically for the Midlands region. Let's put your data detective skills to the test!

## Your Task

1.  **Query the Flood Data:** Use BigQuery to analyze the `midlands_factory_location.midlands_earthengine_grid_data` table. Write a SQL query to identify grid cells with a high flood risk based on the following criteria:
    * Land use (`dw_classnumber`) is classified as "built", "water", or "flooded_vegetation". We want to avoid building on existing structures, bodies of water, or areas already prone to flooding. *(Hint: You can use a `SELECT DISTINCT dw_class_name` query to find the corresponding `dw_classnumber` values)*
    * The average slope (`lidar_mean_slope`) of the grid cell is below 5% (relatively flat terrain). We don't want our factory sliding down a hill!
    * The Height Above Nearest Drainage (`merit_hand`) is below 1.5 meters. Areas close to drainage are more susceptible to flooding.

2.  **Visualize the Risk:** Run your BigQuery query using the Geo Viz tool. Visualize the results on a map, clearly showing the grid cells that meet your high flood risk criteria.

3.  **Capture Your Findings:** Take a clear screenshot of the Geo Viz map. Make sure the geospatial data and the grid are visible, highlighting the areas of high flood risk. Also, copy and save your SQL query.

## Upload

Upload **one file** containing:

* A clear screenshot of the Geo Viz map illustrating the geospatial data and grid with insights into flood risk factors.
* Your SQL query used to identify the high-risk areas.

Good luck, data navigators! Let's keep Cymbal Supplements high and dry!
=======
# Part 1: Flood Risk Assessment: Keeping Dry 🌊

## Instructions

1.  Use BigQuery to query the Earth Engine flood risk data.
2.  Write a query to identify flood risk levels.
3.  Create a geospatial visualization of the flood risk data.
4.  Take a screenshot showing the geospatial visualization and BigQuery query results.
5.  Upload the screenshot.
>>>>>>> 4f9020188bb64225de548fa7529a3fcfde557687
