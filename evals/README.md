# Eval Harness

Static assertion grader for skill-creator eval workspaces. After `/skill-creator` runs evals (with LLM-based grading), this tool adds deterministic checks by applying assertions from `evals.json` to the output files.

## Structure

```
evals/
├── README.md          # This file
├── evals.json         # Test case definitions with static assertions
├── go.mod             # Go module
├── main.go            # Entry point
├── types.go           # Data structures
├── loader.go          # evals.json parser
├── grader.go          # Assertion grading logic
├── workspace.go       # skill-creator workspace integration
├── skill-eval         # Compiled binary (gitignored)
└── results/           # Output from eval runs (gitignored)
```

## Quick Start

Run from the repo root using [Task](https://taskfile.dev):

```bash
# Grade a skill-creator workspace
task eval:grade -- /path/to/workspace

# Or point directly to an iteration directory
task eval:grade -- /path/to/workspace/iteration-1
```

## Workflow

1. Use `/skill-creator` to run evals on a skill. This produces a workspace with LLM-graded results.
2. Run `task eval:grade -- /path/to/workspace` to apply static assertions from `evals.json`.
3. Results appear as `static_grading.json` inside each `with_skill/` and `without_skill/` directory, plus a `static_summary.json` at the iteration level.

## Prerequisites

- **[Task](https://taskfile.dev)** for running commands
- **Go 1.25+** for building

## Adding Test Cases

Add evals to `evals.json` under the appropriate skill's `evals` array. Each eval needs:

| Field | Required | Description |
|-------|----------|-------------|
| `id` | yes | Numeric identifier matching the skill-creator `eval-N` directory name |
| `prompt` | yes | Human-readable description of the eval task |
| `expected_output` | no | Human-readable description of what a good response looks like |
| `assertions` | no | Array of machine-checkable assertions for static grading |

### Example eval

```json
{
  "id": 0,
  "prompt": "Implement a spinner loader animation that rotates continuously.",
  "expected_output": "Should use CSS Animations API, not the shared value API",
  "assertions": [
    {
      "type": "contains",
      "value": "animationName",
      "text": "Uses CSS animation API (animationName)"
    },
    {
      "type": "not_contains",
      "value": "useSharedValue",
      "text": "Does not use shared value API"
    }
  ]
}
```

### Assertion types

| Type | Value | Passes when |
|------|-------|-------------|
| `contains` | substring | Output files include the substring (case-insensitive) |
| `not_contains` | substring | Output files do not include the substring |
| `file_exists` | file path | A file exists at the given path in the run's `outputs/` directory |
| `exit_code` | number | The run's `metadata.json` shows the given exit code |

## Workspace structure

The grader expects a skill-creator workspace laid out like this:

```
iteration-1/
├── eval-0/
│   ├── with_skill/
│   │   ├── outputs/          ← output files from the with-skill run
│   │   ├── eval_metadata.json
│   │   ├── timing.json
│   │   └── static_grading.json   ← written by the grader
│   └── without_skill/
│       ├── outputs/
│       ├── eval_metadata.json
│       ├── timing.json
│       └── static_grading.json
├── eval-1/
│   └── ...
└── static_summary.json           ← written by the grader
```

## Output

The grader writes `static_grading.json` in skill-creator's expectations format:

```json
{
  "expectations": [
    {
      "text": "Uses CSS animation API (animationName)",
      "passed": true,
      "evidence": "Found 'animationName' in response"
    }
  ],
  "summary": {
    "passed": 2,
    "failed": 1,
    "total": 3,
    "pass_rate": 0.67
  }
}
```

## Tips

- Eval IDs in `evals.json` must match the `eval-N` directory names in the workspace (0-indexed)
- Write assertions for the important parts only; don't over-constrain
- Use `expected_output` as a human-readable note even when you have no machine-checkable assertions
- The grader reads all text files in `outputs/` and checks assertions against their combined content
