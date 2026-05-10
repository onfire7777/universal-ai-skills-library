---
name: data-analyst-engine
description: Comprehensive data analysis expertise covering statistical analysis, data visualization, exploratory data analysis, hypothesis testing, machine learning, data cleaning, and data storytelling.
license: Unspecified
metadata:
  version: 1.0.0
  author: Custom Meta-Skill
  tags:
  - data-analysis
  - statistics
  - visualization
  - pandas
  - numpy
  - matplotlib
  - EDA
  - hypothesis-testing
  - data-storytelling
---
# Data Analyst Engine

## Purpose
Provide world-class data analysis capabilities — from raw data ingestion through cleaning, exploration, statistical analysis, visualization, and insight communication.

## The Data Analysis Lifecycle

### Phase 1: Data Understanding
1. **What is the data?** Schema, types, dimensions, source, collection method
2. **What's the question?** Define the specific analytical question
3. **What's the context?** Business domain, stakeholders, decisions to be informed
4. **What are the constraints?** Time, tools, data quality, privacy

### Phase 2: Data Cleaning (80% of the work)
**Systematic cleaning checklist:**
- Check for missing values: `df.isnull().sum()`, decide: drop, impute (mean/median/mode/model), or flag
- Check for duplicates: `df.duplicated().sum()`
- Check data types: `df.dtypes` — are dates parsed? Are numbers stored as strings?
- Check for outliers: IQR method, Z-score, domain knowledge
- Check for inconsistencies: Standardize categories, fix typos, normalize formats
- Check for impossible values: Negative ages, future dates, values outside valid range
- Document every cleaning decision and its rationale

### Phase 3: Exploratory Data Analysis (EDA)
**Univariate Analysis** (one variable at a time):
- Numeric: histogram, box plot, descriptive stats (mean, median, std, skew, kurtosis)
- Categorical: bar chart, frequency table, mode, cardinality

**Bivariate Analysis** (two variables):
- Numeric vs Numeric: scatter plot, correlation coefficient, regression line
- Numeric vs Categorical: box plot by group, violin plot, t-test/ANOVA
- Categorical vs Categorical: contingency table, chi-square test, stacked bar chart

**Multivariate Analysis**:
- Correlation matrix heatmap
- Pair plots
- PCA for dimensionality reduction
- Clustering (K-means, DBSCAN) for pattern discovery

### Phase 4: Statistical Analysis

**Descriptive Statistics**:
| Measure | What It Tells You | When to Use |
|---------|-------------------|-------------|
| Mean | Central tendency | Symmetric distributions |
| Median | Central tendency | Skewed distributions |
| Mode | Most common value | Categorical data |
| Std Dev | Spread | Normal-ish distributions |
| IQR | Spread | Robust to outliers |
| Skewness | Asymmetry | Distribution shape |
| Kurtosis | Tail heaviness | Outlier tendency |

**Hypothesis Testing Framework**:
1. State null hypothesis (H₀) and alternative (H₁)
2. Choose significance level (α = 0.05 typically)
3. Select appropriate test
4. Calculate test statistic and p-value
5. Decide: reject or fail to reject H₀
6. Report effect size and confidence interval (not just p-value!)

**Test Selection Guide**:
| Question | Data Type | Test |
|----------|-----------|------|
| Compare 2 group means | Continuous, normal | Independent t-test |
| Compare 2 group means | Continuous, non-normal | Mann-Whitney U |
| Compare 3+ group means | Continuous, normal | One-way ANOVA |
| Compare 3+ group means | Continuous, non-normal | Kruskal-Wallis |
| Association between categories | Categorical | Chi-square |
| Relationship between 2 variables | Continuous | Pearson/Spearman correlation |
| Predict continuous outcome | Mixed | Linear regression |
| Predict binary outcome | Mixed | Logistic regression |
| Before/after comparison | Paired continuous | Paired t-test / Wilcoxon |

### Phase 5: Data Visualization

**Chart Selection Guide**:
| Purpose | Chart Type | Tool |
|---------|-----------|------|
| Distribution | Histogram, KDE, box plot | matplotlib, seaborn |
| Comparison | Bar chart, grouped bar, dot plot | matplotlib, plotly |
| Relationship | Scatter plot, bubble chart | seaborn, plotly |
| Composition | Pie chart (sparingly!), stacked bar, treemap | plotly |
| Trend over time | Line chart, area chart | matplotlib, plotly |
| Correlation matrix | Heatmap | seaborn |
| Geographic | Choropleth, scatter map | plotly, folium |
| Part-to-whole | Waterfall, Sankey | plotly |

**Visualization Best Practices**:
- Title every chart with the insight, not just the data description
- Label axes with units
- Use colorblind-friendly palettes (viridis, cividis)
- Remove chart junk (unnecessary gridlines, 3D effects, decorations)
- Start bar chart y-axis at zero
- Use consistent color encoding across related charts
- Annotate key data points
- Choose the right aspect ratio for the data

### Phase 6: Data Storytelling

**The Insight Communication Framework**:
1. **Context**: Why does this analysis matter? What question are we answering?
2. **Key Finding**: State the main insight in one sentence
3. **Evidence**: Show the data that supports the finding (visualization + statistics)
4. **Implications**: What does this mean for decisions?
5. **Recommendations**: What action should be taken?
6. **Caveats**: What are the limitations? What could be wrong?

**Storytelling Anti-Patterns**:
- Showing all the data instead of the insight
- Leading with methodology instead of findings
- Using jargon the audience doesn't understand
- Presenting correlation as causation
- Cherry-picking data to support a narrative
- Ignoring contradictory evidence

## Python Data Analysis Toolkit

### Core Libraries
```python
import pandas as pd          # DataFrames and data manipulation
import numpy as np            # Numerical computing
import matplotlib.pyplot as plt  # Static visualizations
import seaborn as sns         # Statistical visualizations
import plotly.express as px   # Interactive visualizations
from scipy import stats       # Statistical tests
from sklearn import *         # Machine learning
```

### Pandas Power Patterns
- **Chaining**: `df.query('age > 30').groupby('city').agg({'salary': ['mean', 'median']}).sort_values(('salary', 'mean'), ascending=False)`
- **Pivot Tables**: `pd.pivot_table(df, values='sales', index='region', columns='quarter', aggfunc='sum')`
- **Window Functions**: `df['rolling_avg'] = df['value'].rolling(window=7).mean()`
- **Merge/Join**: `pd.merge(df1, df2, on='key', how='left')`
- **Apply**: `df['category'] = df['value'].apply(lambda x: 'high' if x > threshold else 'low')`

### Performance Tips
- Use `pd.read_csv(dtype={...})` to specify types upfront
- Use `.query()` instead of boolean indexing for readability
- Use `pd.Categorical` for low-cardinality string columns
- Use `numpy` vectorized operations instead of loops
- For large datasets: consider `polars`, `dask`, or `vaex`
- Profile with `%timeit` and `df.memory_usage(deep=True)`

## Machine Learning Quick Reference

### When to Use ML vs. Statistics
- **Statistics**: When you need to understand WHY (inference, hypothesis testing)
- **ML**: When you need to predict WHAT (prediction, classification, clustering)

### ML Workflow
1. Define the problem (classification, regression, clustering, ranking)
2. Prepare features (encoding, scaling, feature engineering)
3. Split data (train/validation/test — typically 70/15/15)
4. Train baseline model (simple model first)
5. Evaluate (appropriate metrics for the problem type)
6. Iterate (feature engineering, hyperparameter tuning, model selection)
7. Validate on held-out test set (only once!)
8. Deploy and monitor

### Evaluation Metrics
| Problem Type | Metrics | When |
|-------------|---------|------|
| Classification | Accuracy, Precision, Recall, F1, AUC-ROC | Balanced classes |
| Classification (imbalanced) | Precision-Recall AUC, F1, Balanced Accuracy | Rare events |
| Regression | MAE, RMSE, R², MAPE | Continuous prediction |
| Clustering | Silhouette, Calinski-Harabasz, Domain evaluation | Unsupervised |
| Ranking | NDCG, MAP, MRR | Search/recommendation |
