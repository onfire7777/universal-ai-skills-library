# Demo Recording Checklist

Quality gates for each phase of demo video production.

Purpose: Read this when Director must validate demo readiness, playback quality, security hygiene, or final delivery quality.

Contents:
- `Pre-Recording Checklist`: scenario, environment, and accessibility setup
- `Post-Recording Checklist`: playback, content, security, and output validation
- `Pre-Delivery Checklist`: packaging and reproducibility checks
- `Quick Check (3 minutes)`: rapid go/no-go review
- `Quality Scorecard`: `/65` scoring and release thresholds

---

## Pre-Recording Checklist

### Scenario
- [ ] Audience identified
- [ ] Purpose defined (1-2 sentences)
- [ ] One demo, one feature
- [ ] Story arc defined (opening / action / result)
- [ ] Emphasis points identified

### Environment
- [ ] Using `playwright.config.demo.ts`
- [ ] slowMo appropriate for audience
- [ ] `video: { mode: 'on' }` set
- [ ] Resolution matches target device
- [ ] Test data seeded
- [ ] Auth state prepared

### Accessibility
- [ ] Key screens validated with ARIA snapshots (if applicable)
- [ ] Overlay text meets contrast requirements
- [ ] Caption/subtitle plan defined for external distribution

---

## Post-Recording Checklist

### Playback
- [ ] Video plays to completion
- [ ] No choppy frames or flickering
- [ ] Text is readable at target resolution

### Content
- [ ] All operations visible at demo speed
- [ ] Pacing feels natural (not rushed, not dragging)
- [ ] Overlays readable with sufficient display time
- [ ] Final result displayed long enough

### Security
- [ ] No passwords, tokens, or production data visible
- [ ] No console errors or network errors visible
- [ ] Demo data used (realistic but fictional)

### Output
- [ ] Naming convention followed: `[feature]_[action]_[date].webm`
- [ ] MP4 conversion generated (if needed for distribution)
- [ ] File size reasonable (~5-10MB for 30s at 720p)
- [ ] Saved in `demos/output/` directory

---

## Pre-Delivery Checklist

- [ ] Demo meets original objectives
- [ ] UI matches current product state
- [ ] Scenario document saved for reproducibility
- [ ] Metadata recorded (date, version, environment)
- [ ] Output formats match distribution channels (WebM/MP4/GIF)

---

## Quick Check (3 minutes)

```
[] Plays to the end
[] Followable pace
[] No confidential information
[] Clear result
[] Correct filename and format
```

---

## Quality Scorecard

Rate each item 1-5.

| Category | Item | Score (1-5) |
|----------|------|-------------|
| **Story** | Clear beginning/ending | |
| | Appropriate pacing | |
| | Effective emphasis points | |
| **Visual** | Resolution/quality | |
| | Overlay readability | |
| | Brand consistency | |
| **Technical** | Stability/reproducibility | |
| | Selector robustness | |
| | Wait strategy (locator-based) | |
| **Accessibility** | Captions available (if external) | |
| | ARIA validation at key screens | |
| **Data** | Realism | |
| | Confidential info excluded | |
| **Total** | | /65 |

### Score Assessment

- **52-65**: Excellent - Ready for delivery
- **40-51**: Good - Deliverable with minor fixes
- **30-39**: Needs improvement - Major issues require fixing
- **< 30**: Reshoot - Review from scenario
