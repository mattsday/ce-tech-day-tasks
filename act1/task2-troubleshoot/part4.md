::important[This step may have an affect on future tasks and point scoring!]

Right now, your orders-service is running with the compute engine default service accounts - **which is like leaving the keys to the time machine under the doormat**.

In the world of supplements, security is paramount. You need to create dedicated service accounts for your the Cloud Run service, granting it only the necessary permissions. This ensures that even if a breach occurs, the damage is contained. We need to secure the formula, it is too important to be left to chance.

### Task

1. Configure a new service accounts and assign it to your Cloud Run instance
2. Ensure they have only the permissions required for their given tasks
    * `roles/cloudsql.client`
3. Verify by putting their names below
