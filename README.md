# Tasks Tracker

Based on the idea from roadmap.sh: https://roadmap.sh/projects/task-tracker

## Usage
```shell
# Adding a new task
task-cli add "Buy groceries"
# Output: Task added successfully (ID: 1)

# Updating and deleting tasks
task-cli update 1 "Buy groceries and cook dinner"
task-cli delete 1

# Marking a task as in progress or done
task-cli mark-in-progress 1
task-cli mark-done 1

# Listing all tasks
task-cli list

# Listing tasks by status
task-cli list done
task-cli list todo
task-cli list in-progress
```


## Project tracking features

| Feature | Description | Status
| ------ | ------ | ------
| command line | create a binary that is accessible in PATH to interface with | DONE
| storage/user data | design and implement fastest storage option for json files | TODO
| pipeline | github actions to build project and tag with version number | TODO
| document | update readme with full doc requirements | TODO