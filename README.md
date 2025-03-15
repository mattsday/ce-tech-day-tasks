# CE Tech Day Tasks

Here is where you define the tasks. The default ones are a great starting point - they allow you to explore Gemini across Google Cloud!

However, if you wish to customise them, then this is the right place! Customisation ranges from simple tasks (e.g. upload a screenshot) to much more complex.

## Contributing

Ideally create a new branch and commit your changes. Then issue a Merge Request in Gitlab.

### Creating New Tasks

Use the [CE Tech Day Gemini Gem](https://gemini.google.com/corp/gem/80226b51667c) to create tasks. Try it out!

### Getting Help

Don't hesitate to ping `mattsday@`!

## Deployment

Tasks go live immediately after deployment - **including updating all pages in real-time**.

For example, if you add a new task it will instantly show up on the leaderboard and be visible in the score-app right away. You can use this functionality as a Hack Director to add a level of excitement into the event.

### To Deploy

The easiest way to deploy is to use Cloud Build. Run `deploy.sh` or `fast-deploy.sh` in this folder.

For faster deployment, use the taks-tool in [utils/task-tool](../utils/task-tool/). This requires Go to be installed.

## Structure

### Folder Layout

At the root is an optional `tasks.yaml` which describes your event. If this is not present then a default configuration is used.

Each task must be configured in a file named `task.yaml` placed inside any sub-folder under this directory structure. For example, your tree could look like this:

- `tasks.yaml`
- `intro/task.yaml`
- `intro/bonus/task.yaml`
- `task1/task.yaml`
- `task2/task.yaml`

## tasks.yaml Configuration

```yaml
---
metadata:
  version: 1
event:
  name: hacksday
  scoring_enabled: true
```

### Event Section

| Key             | Description                                                                                                                 |
| --------------- | --------------------------------------------------------------------------------------------------------------------------- |
| name            | **Required**: Name for your event - e.g. hacksday                                                                           |
| scoring_enabled | **Required**: Whether scoring is enabled - set to `false` at the end. Only affects the score app - judging is not impacted. |

## task.yaml Configuration

Each `task.yaml` file has a number of options

```yaml
---
metadata:
  version: 1
task:
  id: 5-mytask
  alias:
    - mytask
  name: "My Amazing Task 🤗"
  description: "A short description of your task"
  overview:
    - This is line one of my amazin task's overview - sometimes 1 line is enough
    - In other cases though it might need a few more.
  enabled: true
  hidden: false
  lb_hidden: false
parts:
  - id: part1
    name: Image Upload Task
    open: true
    type: image
    upload_text: "Upload screenshot"
    max_points: 500
    instructions_link: part1.md
    good_examples:
      - examples/part1_a.png
    llm_instructions: Here are some extra instructions for the LLM to score this task
  - id: logo
    name: Custom Part 🧐
    challenge: true
    type: custom
    component: CustomComponent
    max_points: 1000
    instructions_link: part2.md
```

### Task Section

| Key         | Description                                                                                                                                                                                                                                                                                                                                               |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| id          | **Required**: Globally unique ID for this task. This will be sorted alphabetically in the UI (e.g. leaderboard) so keying this is important (e.g. `0-registration` will appear before `1-task1` etc).                                                                                                                                                     |
| name        | **Required**: Human-readable name for the task                                                                                                                                                                                                                                                                                                            |
| description | **Required**: Short description used in tooltips and other areas                                                                                                                                                                                                                                                                                          |
| overview    | **Required**: Array describing the task. Used as an introduction in the score app                                                                                                                                                                                                                                                                         |
| alias       | **Optional**: List of aliases for this task (unimplemented)                                                                                                                                                                                                                                                                                               |
| enabled     | **Optional**: Whether this task is active in the score app. It will be visible regardless - this controls whether it can be interacted with (default: `false` )                                                                                                                                                                                           |
| hidden      | **Optional**: Whether this task can be seen anywhere. For example, to (default: `false` )                                                                                                                                                                                                                                                                 |
| lb_hidden   | **Optional**: Whether this task is hidden from the leaderboard. Value overrides `hidden` above for the leaderboard only (i.e. if `hidden === false && lb_hidden === true` then the task will show on the leaderboard). Used, for example, to temporarily hide a bonus task from the leaderboard that the Hack Director may show later (default: `false` ) |

### Parts Section

| Key               | Description                                                                                                                                                                                                    |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| id                | **Required**: Task-unique ID for this part. Parts are kept in order (i.e. the first part goes first) and IDs are used for scoring. If in doubt, `part1` , `part2` , `part3` etc should suffice.                |
| name              | **Required**: Human-readable name for the part                                                                                                                                                                 |
| max_points        | **Required**: The maximum score for this part that judges may award                                                                                                                                            |
| instructions_link | **(Typically) Required:** The link to the markdown file including this part's instructions (e.g. `part1.md`). See Instructions Format for more information.                                                    |
| type              | **Required**: The type of task being judged. For screenshots it should be `image`, otherwise set to `custom`. See [Custom Components](#custom-components) below for more information                           |
| component         | For `custom` types set the component to be used. See [Custom Components](#custom-components) below.                                                                                                            |
| open              | **Optional**: If the task view should be open by default (default: `false`)                                                                                                                                      |
| upload_text       | The text to show in the "Upload Image" button. Set when task type is `image`.                                                                                                                                  |
| good_examples     | **Optional**: An array list of up to three good examples for the LLM to use to judge submissions. Useful for tasks where there's a clear right/wrong answer. Typically only applicable for type: `image` tasks |
| llm_instructions  | **Optional**: Any extra instructions to pass to the LLM for judging. Typically only applicable for type: `image` tasks                                                                                         |
| hidden            | **Optional**: Whether a part is hidden or not. Default value: `false`.                                                                                                                                         |

### Instructions Format

Instructions should be written in Markdown format. There are some special instructions that can be sent in addition to all standard Markdown tags.

#### Info

To call out an information block, use an info code block:

````md
::alert[It likely will not work! This is OK! Stay with us.]{severity=info}

#### Important

To mark something as important use an alert code block:

````md
::alert[Important: Ensure you create a new Notebook before continuing]{severity=warning}
````

#### Challenges

To call out a challenge, use a challenge block:

```md
:::challenge
Use Code assist to improve the app!
:::
````

#### Custom Code Blocks

To create custom code-blocks, edit [MarkdownRender.tsx](../apps/score-app/components/MarkdownRender.tsx) in the score app to configure your custom type.

### Custom Components

In situations where you wish to provide a custom component rather than an image upload task you will need to write code for the Score App to render it and also display it.

In the score-app, Look inside [components/custom/Custom.tsx](../apps/score-app/components/custom/Custom.tsx) for an example of how to register a component, and look at [CloudRunDeploy.tsx](../apps/score-app/components/custom/CloudRunDeploy.tsx) for an example of an implementation.

Most likely you will also need to create the scoring logic. Look at [libs/scoring/tasks/task3.ts](../apps/score-app/libs/scoring/tasks/task3.ts) for inspiration.
