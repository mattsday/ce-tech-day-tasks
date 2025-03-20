# Task Tool

This tool deploys tasks to the live service. Changes are immediate, including in open sessions.

For more information on task configuration, see [the Task documentation](../../tasks/README.md).

## Usage

In most cases you can use [deploy.sh](../../tasks/deploy.sh) in the Task folder to deploy without invoking this utility directly. This will also fallback to using Cloud Build should Go fail for any reason.

### Updating all Tasks

```sh
. ../../setup/config.sh
export HOST_BUCKET="${HOST_PROJECT}-score-assets"
go run main.go --base-folder "../../tasks" --bucket "${HOST_BUCKET}" --host-pid "${HOST_PROJECT}"
```

### Fast Updating

You can pass `--upload=false` to skip file uploads which makes the deployment time significantly faster. Only do this if the markdown and image files have not changed:

```sh
. ../../setup/config.sh
export HOST_BUCKET="${HOST_PROJECT}-score-assets"
go run main.go --base-folder "../../tasks" --bucket "${HOST_BUCKET}" --host-pid "${HOST_PROJECT}" --upload=false
```
