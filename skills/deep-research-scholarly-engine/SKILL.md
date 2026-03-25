---
name: Deep Research & Scholarly Analysis Engine
description: Comprehensive research methodology combining scholarly research, scientific data analysis, evidence synthesis, source evaluation, and systematic review practices for producing research-grade outputs.
version: 1.0.0
author: Custom Meta-Skill
tags: [research, scholarly, scientific, evidence-based, systematic-review, data-analysis, google-scholar, consensus]
---

# Deep Research & Scholarly Analysis Engine

## Purpose
Enable rigorous, evidence-based research that meets academic standards. Combine multiple research methodologies, source evaluation frameworks, and data synthesis techniques to produce outputs grounded in the best available evidence.

## Research Methodology Hierarchy

### Level 1: Systematic Review (Highest Rigor)
Use when: Critical decisions, health/safety, policy recommendations
1. Define a precise research question (PICO format for clinical: Population, Intervention, Comparison, Outcome)
2. Develop comprehensive search strategy across multiple databases
3. Apply inclusion/exclusion criteria systematically
4. Extract data using standardized forms
5. Assess quality of evidence (GRADE framework)
6. Synthesize findings with appropriate methods
7. Report following PRISMA guidelines

### Level 2: Structured Literature Review
Use when: Technical decisions, architecture choices, best practices
1. Define scope and research questions
2. Search 3+ independent sources
3. Evaluate source quality (CRAAP test)
4. Identify themes and patterns
5. Synthesize with explicit methodology
6. Acknowledge limitations

### Level 3: Rapid Evidence Assessment
Use when: Time-constrained decisions, initial exploration
1. Focused search on 2-3 key sources
2. Quick quality assessment
3. Extract key findings
4. Provide confidence-weighted conclusions

## Source Evaluation: The CRAAP Test
For every source, evaluate:
- **Currency**: When was it published/updated? Is it current enough for the topic?
- **Relevance**: Does it directly address the research question?
- **Authority**: Who is the author? What are their credentials? Is the publisher reputable?
- **Accuracy**: Is the information supported by evidence? Can it be verified? Is it peer-reviewed?
- **Purpose**: Why does this information exist? Is there bias? Is it trying to sell something?

Score each dimension 1-5. Sources scoring below 15/25 should be used cautiously or discarded.

## Evidence Quality Pyramid (Strongest to Weakest)
1. **Systematic Reviews & Meta-Analyses** — Gold standard
2. **Randomized Controlled Trials** — Strong causal evidence
3. **Cohort Studies** — Good observational evidence
4. **Case-Control Studies** — Moderate evidence
5. **Case Series / Case Reports** — Weak evidence
6. **Expert Opinion / Editorials** — Lowest evidence level
7. **Anecdotal / Blog Posts** — Not evidence (but may suggest hypotheses)

## Search Strategy Best Practices

### Academic/Scholarly Search
- **Google Scholar**: Use `site:`, `intitle:`, date ranges, cited-by chains
- **Consensus.app**: For AI-synthesized scientific consensus on specific claims
- **Semantic Scholar**: For citation graph exploration and related papers
- **PubMed**: For biomedical and life sciences
- **arXiv**: For preprints in CS, physics, math, AI/ML
- **SSRN**: For social sciences and economics

### Search Technique: Citation Chaining
1. Find one highly relevant paper
2. **Forward chain**: Who cited this paper? (Google Scholar "Cited by")
3. **Backward chain**: What did this paper cite? (Reference list)
4. **Lateral chain**: What other papers do the same authors write?
5. Repeat until saturation (no new relevant papers found)

### Search Technique: Boolean Strategy
- Use AND to narrow: `"machine learning" AND "healthcare" AND "diagnosis"`
- Use OR to broaden: `"deep learning" OR "neural network"`
- Use NOT to exclude: `"AI" NOT "artificial insemination"`
- Use quotes for exact phrases: `"transformer architecture"`
- Use wildcards: `optim*` matches optimize, optimization, optimal

## Data Analysis Framework

### Quantitative Analysis Protocol
1. **Data Cleaning**: Check for missing values, outliers, inconsistencies
2. **Descriptive Statistics**: Mean, median, mode, SD, range, distribution shape
3. **Exploratory Analysis**: Visualizations, correlations, patterns
4. **Inferential Statistics**: Hypothesis testing, confidence intervals, effect sizes
5. **Interpretation**: What do the numbers actually mean in context?

### Statistical Reasoning Checklist
- Is the sample size adequate?
- Is the sample representative?
- Are the statistical tests appropriate for the data type?
- Is statistical significance confused with practical significance?
- Are confidence intervals reported (not just p-values)?
- Is the effect size meaningful?
- Could there be confounding variables?
- Is correlation being confused with causation?

### Qualitative Analysis Protocol
1. **Thematic Analysis**: Identify recurring themes across sources
2. **Content Analysis**: Systematic categorization of textual data
3. **Comparative Analysis**: How do different sources agree/disagree?
4. **Gap Analysis**: What questions remain unanswered?

## Evidence Synthesis Methods

### Narrative Synthesis
- Organize findings by theme, not by source
- Identify areas of agreement and disagreement
- Weight findings by evidence quality
- Explicitly state the strength of evidence for each conclusion

### Vote Counting
- How many studies support conclusion A vs B?
- Weight by study quality and sample size
- Report the ratio with confidence assessment

### Triangulation
- Do multiple independent sources/methods converge on the same conclusion?
- If yes: High confidence
- If mixed: Moderate confidence, report the disagreement
- If contradictory: Low confidence, investigate why

## Research Output Standards

### Every Research Output Must Include:
1. **Clear Research Question**: What exactly are we investigating?
2. **Methodology Statement**: How did we search and what criteria did we use?
3. **Source Documentation**: All sources cited with full references
4. **Evidence Quality Assessment**: How strong is the evidence?
5. **Confidence Level**: How confident are we in each conclusion?
6. **Limitations**: What are the gaps and weaknesses?
7. **Recommendations**: What actions does the evidence support?

### Citation Standards
- Always cite specific sources for factual claims
- Prefer primary sources over secondary
- Include publication date for currency assessment
- Note if a source is pre-print, peer-reviewed, or grey literature
- Use inline numeric citations with reference list

## Internet Parsing & Search Mastery

### Web Content Extraction Hierarchy
1. **Structured APIs** (best): Use official APIs when available
2. **Structured Data**: Look for JSON-LD, schema.org markup, RSS feeds
3. **Clean HTML Parsing**: Extract from semantic HTML elements
4. **Full Page Rendering**: For JavaScript-heavy sites
5. **Screenshot + OCR**: Last resort for complex layouts

### Source Triangulation Protocol
For any factual claim from the internet:
1. Find the **primary source** (original study, official announcement, raw data)
2. Find **2+ independent confirmations** from reputable sources
3. Check for **contradicting evidence** actively
4. Assess **recency** — is this still current?
5. Check for **corrections or retractions**

### OSINT Best Practices
- Start broad, narrow progressively
- Use multiple search engines (results differ)
- Check the Wayback Machine for historical context
- Verify images with reverse image search
- Cross-reference social media claims with official sources
- Be aware of information warfare and deliberate misinformation

## Anti-Patterns in Research
- **Cherry-picking**: Selecting only evidence that supports a predetermined conclusion
- **Appeal to Authority**: Accepting claims because of who said them, not the evidence
- **Recency Bias**: Assuming newer = better without evaluation
- **Survivorship Bias**: Only looking at successful cases
- **Publication Bias**: Published studies skew positive; negative results are underreported
- **P-hacking**: Statistical manipulation to achieve significance
- **HARKing**: Hypothesizing After Results are Known
- **Ecological Fallacy**: Applying group-level findings to individuals
