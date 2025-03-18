Our Firestore security rules are currently about as effective as a screen door on a submarine – they need some serious reinforcement.

Your mission, should you choose to accept it (and you better, because the fate of Cymbal Supplements' reputation hangs in the balance!), is to fortify those rules and make our data impenetrable.

::challenge[Ensure only authorised, logged-in users can read from the **users collection** - at the moment anyone can!]

### Task

Update your [Firestore Security Rules](https://console.firebase.google.com/project/%%CLIENT_PROJECT_ID%%/firestore/databases/-default-/rules) with the following constraints:

- **Collection:** `users`
  - **Read**: Users with the custom claim "admin" set to true
  - **Write**: Users with the custom claim "admin" set to true
- **Collection:** `reviews`
  - **Read**: Anyone (public to the internet)
  - **Write**: No one (written by the server)
