Great! The search widget is working! Byan's followers will be very happy.

Now it's time to embed the widget onto our reviews site!

## Task

1. **Open the Code**:
   1. Open Cloud Workstations and start your instance: [Cloud Workstations](https://console.cloud.google.com/workstations/overview?project=%%CLIENT_PROJECT_ID%%).
   2. Inside workstations, select **Open Folder** and select **retail-site**
   3. Select the menu in the top-left and select **Terminal --> New Terminal**
2. Start a mini http server on port 8000:

```bash
cd public
python -m http.server 8000
```

3. Open the local preview of the site and ensure it is working
4. Use the integration steps for your search widget in the previous steps to add it to the retail site - replacing the existing search box

::info[**Hint**: Don't forget to allowlist your workstation preview domain in your AI Application Config before testing!]

5. Take a screenshot of the embedded search app and upload it below for points!
