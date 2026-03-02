# Skills Repository

Software Mansion's official collection of agent skills for Claude Code.

## Repository Structure

```
skills/
├── <skill-name>/
│   ├── SKILL.md          # Required — metadata + core instructions
│   ├── references/       # Optional — detailed docs loaded on demand
│   ├── examples/         # Optional — runnable code samples
│   ├── scripts/          # Optional — utility scripts
│   └── templates/        # Optional — template files for output
```

## Skill Conventions

### SKILL.md format

Every skill must have a `SKILL.md` with YAML frontmatter:

```yaml
---
name: skill-name            # lowercase, hyphen-separated
description: "..."          # trigger phrases and description (third-person)
version: 1.0.0
metadata:
  author: Software Mansion
  license: MIT
allowed-tools: Bash(...)    # optional tool whitelist
user-invocable: true        # whether users can invoke directly (default: true)
---
```

### Writing guidelines

- Keep `SKILL.md` body under 2,000 words. Move detailed content to `references/`.
- Use imperative, instructional language: "Run the validator" over "You should run the validator."
- Write `description` in third person with specific trigger phrases in quotes.
- Include concrete examples in the body or in `examples/`.

### Naming

- Skill directories use lowercase with hyphens: `react-native-debugging`, `expo-router`.
- Reference files use lowercase with hyphens: `common-patterns.md`, `api-reference.md`.

### Testing

Before submitting a skill, verify:
1. `SKILL.md` has valid YAML frontmatter with `name` and `description`.
2. The skill triggers correctly on the phrases listed in `description`.
3. All referenced files in `references/`, `examples/`, `scripts/` exist.
4. Instructions are clear and produce correct results.
