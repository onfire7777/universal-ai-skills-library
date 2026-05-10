# Prompt Template

```markdown
Create a high-density infographic about: <topic>.

Target audience: <target audience>.
Layout: <layout>.
Style direction: <style>.
Aspect ratio: <aspect>.
Language: <zh|en|ja|ko>.

Main title: <title>.
Subtitle: <subtitle or leave empty>.

Sections:
- Section 1: <heading> | key points: <point 1>; <point 2>; <point 3>
- Section 2: <heading> | key points: <point 1>; <point 2>; <point 3>
- Section 3: <heading> | key points: <point 1>; <point 2>; <point 3>

Visual requirements:
- clear hierarchy
- readable labels
- consistent icon style
- enough whitespace for clarity
- strong section separation

Avoid: cluttered background, unreadable tiny text, broken charts, random decorative noise, watermark.
```

Additional rules:

- Place data, years, and terminology directly into the relevant sections whenever possible.
- If the user wants a very dense one-page information graphic, prefer `dense-modules`.
- If the user wants a technical blueprint look, add a phrase such as `blueprint-like annotations and schematic labeling` to the prompt body.
