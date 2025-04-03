The Plot Thickens... (Not the Sauce, Hopefully!)

Now that we know where *not* to build (thanks to your flood risk assessment!), let's focus on finding suitable land for our state-of-the-art hot sauce factory. We're not just looking for any patch of green; we're on the hunt for **brownfield land**.

Brownfield land, for those not in the know, is previously developed land that may be vacant, derelict, or underutilized. Often, old buildings have been demolished and cleared, making it potentially easier (and with fewer bureaucratic hurdles) to build new structures.

For this part of the mission, we'll be using the amazing **OpenStreetMap (OSM)** dataset, which is conveniently available as a public dataset in BigQuery. Specifically, we're interested in the `planet_features_multipolygons` table, which contains polygon features representing various land uses.

### Task

1.  **Explore OSM Data:** We're interested in land plots marked with the `landuse` key set to the value `brownfield`. These tags provide valuable information about how the land is currently (or was previously) used. You can learn more about the `landuse` key and its possible values here: [https://wiki.openstreetmap.org/wiki/Key:landuse](https://wiki.openstreetmap.org/wiki/Key:landuse)

2.  **Complete the Query:** We've started writing a BigQuery SQL query for you. Your task is to complete the `UNNEST` clause to correctly filter for brownfield land within our area of interest (the "Midlands sample area of interest 1").

    ```sql
    SELECT
        osm_way_id,
        ST_AREA(geometry) AS building_area_sqm,
        geometry,
    FROM `bigquery-public-data.geo_openstreetmap.planet_features_multipolygons` AS osm_table,
        (SELECT full_outer_polygon FROM `midlands_factory_location.midlands_areas_of_interest`
        WHERE sample_area_desc = "Midlands sample area of interest 1") AS area_of_interest
    WHERE EXISTS
        (
            -- COMPLETE THE UNNEST QUERY HERE
        )
    AND ST_WITHIN(geometry,area_of_interest.full_outer_polygon)
    ```

    **Tips for the UNNEST query:**

    * The `all_tags` column in the `osm_table` contains an array of key-value pairs. You'll need to use the `UNNEST` function to transform this array into individual rows, where each row has a `key` and a `value`.
    * Create a condition within the `EXISTS` clause to check if at least one tag in the `all_tags` array has a `key` equal to 'landuse' and a corresponding `value` equal to 'brownfield'.

3.  **Visualize Brownfield Areas:** Run your completed query in BigQuery using the Geo Viz tool. Visualize the resulting brownfield polygons on a map.

4.  **Capture Your Findings:** Take a clear screenshot of the Geo Viz map, highlighting the potential brownfield sites within the Midlands sample area. Also, make sure to save your completed SQL query.

### Scoring

Upload a screenshot containing:

* A clear screenshot of the Geo Viz map highlighting potential brownfield sites.
* Your completed SQL query used to identify brownfield land.

Let's find some prime real estate for Cymbal Supplements! No more squeezing sauce in cramped conditions!
