<<<<<<< HEAD
## The Grand Finale: Finding the Perfect Spot!

You've braved the flood risks and navigated the urban landscape to identify potential brownfield sites. Now comes the ultimate challenge: combining your findings to suggest the *best* locations for the new Cymbal Supplements factory!

We need to find brownfield land that is **not** in a high flood risk area. To do this efficiently, we'll use a SQL query that cleverly joins the results of your previous analyses.

We've provided a starting query structure that uses the `WITH` clause to create two temporary tables:

* `brownfield`: This table (which you helped define in Part 2) lists all the brownfield OSM polygons within our area of interest.
* `high_flood_risk_areas`: This table will contain a single `MULTIPOLYGON` geometry representing the union of all the high flood risk grid cells you identified in Part 1.

## Here's the query structure:

```sql
WITH brownfield AS
(
    SELECT
        osm_way_id,
        ST_AREA(geometry) AS building_area_sqm,
        geometry,
    FROM `bigquery-public-data.geo_openstreetmap.planet_features_multipolygons` AS osm_table,
        (SELECT full_outer_polygon FROM `midlands_factory_location.midlands_areas_of_interest`
        WHERE sample_area_desc = "Midlands sample area of interest 1") AS area_of_interest
    WHERE EXISTS
        (
            UNNEST(osm_table.all_tags) AS tags
            WHERE tags.key = 'landuse' AND tags.value = 'brownfield'
        )
    AND ST_WITHIN(geometry,area_of_interest.full_outer_polygon)
), high_flood_risk_areas AS
(
    SELECT
        ST_UNION(gridcell_outline_polygon) AS all_high_risk_areas
    FROM `midlands_factory_location.midlands_earthengine_grid_data`
    WHERE (
        -- WHERE CONDITION FROM PART 1 HERE
    )
)

SELECT
    brownfield.osm_way_id,
    brownfield.geometry,
    brownfield.building_area_sqm
FROM brownfield
INNER JOIN high_flood_risk_areas
ON ST_DISJOINT(brownfield.geometry, high_flood_risk_areas.all_high_risk_areas)
ORDER BY brownfield.building_area_sqm DESC
LIMIT 5
```


## Your Task

1.  **Complete the Aggregate Function:** In the `high_flood_risk_areas` CTE (Common Table Expression), you need to complete the `AGGREGATE FUNCTION` to group together all the individual high flood risk grid cell polygons into a single `MULTIPOLYGON` geometry. This is crucial for the subsequent spatial join. The function you need is `ST_UNION()`.

2.  **Complete the WHERE Condition:** In the `high_flood_risk_areas` CTE, you need to insert the `WHERE` condition you used in **Part 1** to identify the high flood risk grid cells.

3.  **Run the Final Query:** Execute the completed SQL query in BigQuery. This query will identify the top 5 largest brownfield areas that are spatially disjoint (i.e., do not overlap) with the high flood risk areas.
4.  **Visualize and Capture:** Use the Geo Viz tool to visualize the results of this final query. Focus on the top 5 suggested brownfield locations. Take a screenshot of the Geo Viz map, ideally in a satellite view, clearly highlighting these potential factory sites.

## Upload
Upload one file containing:
A satellite view of the Geo Viz map highlighting the top 5 suggested factory locations.








Congratulations, you've successfully navigated the complexities of geospatial data to help Cymbal Supplements find the perfect location for their expansion! Now, let the hot sauce empire grow!
=======
# Part 3: Location Suggestions: Pinpointing Potential 📍

## Instructions

1.  Based on the flood risk and land use analyses, identify suitable locations for the new factory.
2.  Circle the top 3 suggested locations on the location map.
3.  Take a screenshot of the map with the circled locations.
4.  Upload the screenshot.
>>>>>>> 4f9020188bb64225de548fa7529a3fcfde557687
