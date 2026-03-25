---
name: AI Expertise Engine
description: Comprehensive AI/ML expertise covering prompt engineering, LLM architecture, AI agent design, RAG systems, fine-tuning, AI safety, and cutting-edge AI research for building and leveraging AI systems.
version: 1.0.0
author: Custom Meta-Skill
tags: [AI, machine-learning, LLM, prompt-engineering, agents, RAG, fine-tuning, AI-safety, deep-learning]
---

# AI Expertise Engine

## Purpose
Provide world-class AI expertise across the full spectrum — from prompt engineering and LLM usage to AI system architecture, agent design, RAG pipelines, fine-tuning, safety, and cutting-edge research.

## Prompt Engineering Mastery

### The Prompt Engineering Hierarchy
1. **System Prompt Design**: Define persona, constraints, output format, and behavioral rules
2. **Few-Shot Examples**: Provide 2-5 high-quality input/output examples
3. **Chain-of-Thought**: "Think step by step" / "Let's work through this systematically"
4. **Structured Output**: Specify exact JSON schema, Markdown format, or template
5. **Meta-Prompting**: Prompt the AI to generate better prompts

### Advanced Prompt Techniques
- **Role Assignment**: "You are a senior security researcher with 20 years of experience..."
- **Constraint Setting**: "You MUST cite sources. You MUST NOT speculate without evidence."
- **Output Templating**: "Respond in this exact format: [template]"
- **Self-Consistency**: Generate multiple responses and pick the most common answer
- **Tree of Thought**: Explore multiple reasoning paths, evaluate each, select the best
- **ReAct Pattern**: Reason → Act → Observe → Reason → Act (for tool-using agents)
- **Reflection Prompting**: "Review your answer. What might be wrong? Revise if needed."
- **Decomposition**: Break complex tasks into subtasks with separate prompts for each

### Prompt Anti-Patterns
- Vague instructions without specific output format
- Too many instructions at once (cognitive overload)
- Contradictory constraints
- Assuming the model knows your context
- Not providing examples when the task is ambiguous
- Over-constraining creativity when exploration is needed

## LLM Architecture Understanding

### Transformer Architecture
- **Self-Attention**: Allows each token to attend to all other tokens (O(n²) complexity)
- **Multi-Head Attention**: Multiple attention patterns in parallel
- **Feed-Forward Networks**: Position-wise transformations
- **Layer Normalization**: Stabilizes training
- **Positional Encoding**: Injects sequence order information
- **KV Cache**: Stores key-value pairs for efficient autoregressive generation

### Key Model Parameters
- **Temperature**: 0.0 (deterministic) → 1.0 (creative) → 2.0 (chaotic)
- **Top-p (nucleus sampling)**: Cumulative probability threshold (0.9 = top 90% probability mass)
- **Top-k**: Consider only top k tokens
- **Max tokens**: Output length limit
- **Frequency/Presence penalty**: Reduce repetition
- **Stop sequences**: Define where generation should stop

### Model Selection Guide
| Use Case | Best Model Type | Why |
|----------|----------------|-----|
| Complex reasoning | Large frontier models (GPT-4.1, Claude 3.5, Gemini 2.5) | Maximum capability |
| Fast simple tasks | Small models (GPT-4.1-mini, Haiku, Flash) | Speed + cost efficiency |
| Code generation | Code-specialized models | Domain optimization |
| Embedding/search | Embedding models (text-embedding-3, voyage) | Vector representation |
| Image understanding | Multimodal models | Vision capability |
| Real-time/streaming | Models with streaming support | Low latency |

## AI Agent Architecture

### Agent Design Patterns
1. **ReAct Agent**: Reason → Act → Observe loop with tool access
2. **Plan-and-Execute**: Create full plan first, then execute steps
3. **Reflexion Agent**: Execute → Reflect → Improve → Re-execute
4. **Multi-Agent Systems**: Specialized agents collaborating (researcher, coder, reviewer)
5. **Hierarchical Agents**: Manager agent delegates to worker agents
6. **Agentic Workflows**: DAG-based task orchestration with conditional branching

### Tool Use Best Practices
- Define tools with clear names, descriptions, and parameter schemas
- Provide examples of when to use each tool
- Handle tool errors gracefully with retry logic
- Implement rate limiting and cost controls
- Log all tool calls for debugging and auditing
- Use structured output (JSON) for tool parameters

### Agent Memory Systems
- **Short-term**: Conversation context window
- **Working Memory**: Scratchpad for current task state
- **Long-term**: Vector database for retrieval (RAG)
- **Episodic**: Specific past interactions and outcomes
- **Semantic**: General knowledge and facts
- **Procedural**: How to perform specific tasks (skills)

## RAG (Retrieval-Augmented Generation)

### RAG Pipeline Architecture
1. **Ingestion**: Document loading → Chunking → Embedding → Vector store
2. **Retrieval**: Query embedding → Similarity search → Re-ranking → Context assembly
3. **Generation**: Retrieved context + Query → LLM → Response with citations

### Chunking Strategies
- **Fixed-size**: Simple but may split semantic units
- **Semantic**: Split on paragraph/section boundaries
- **Recursive**: Try large chunks first, split smaller if needed
- **Agentic**: Use LLM to determine optimal chunk boundaries
- **Overlap**: Include 10-20% overlap between chunks for context continuity

### Retrieval Optimization
- **Hybrid Search**: Combine vector similarity + keyword (BM25) search
- **Re-ranking**: Use cross-encoder models to re-rank initial results
- **Query Expansion**: Generate multiple query variants for broader recall
- **Metadata Filtering**: Pre-filter by date, source, category before vector search
- **Contextual Compression**: Compress retrieved chunks to only relevant parts

## Fine-Tuning & Training

### When to Fine-Tune vs. Prompt Engineer
- **Prompt Engineering**: Try this first. Works for most use cases.
- **Few-Shot + RAG**: When you need domain knowledge but not style changes.
- **Fine-Tuning**: When you need consistent style, format, or domain-specific behavior that prompting can't achieve.
- **Pre-Training**: Almost never needed. Only for entirely new domains or languages.

### Fine-Tuning Best Practices
- Start with high-quality training data (quality > quantity)
- Use at least 50-100 high-quality examples
- Include diverse examples covering edge cases
- Evaluate on a held-out test set
- Monitor for overfitting (training loss vs. validation loss)
- Use LoRA/QLoRA for parameter-efficient fine-tuning
- Version control your training data and model checkpoints

## AI Safety & Alignment

### Key Safety Principles
- **Harmlessness**: Don't generate harmful, illegal, or dangerous content
- **Honesty**: Don't fabricate information; acknowledge uncertainty
- **Helpfulness**: Actually solve the user's problem
- **Transparency**: Be clear about capabilities and limitations
- **Privacy**: Don't leak training data or user information

### Hallucination Mitigation
1. Ground responses in retrieved context (RAG)
2. Ask the model to cite specific sources
3. Use lower temperature for factual tasks
4. Implement fact-checking pipelines
5. Use structured output to constrain responses
6. Chain-of-thought to make reasoning explicit and verifiable

### Evaluation Metrics
- **Accuracy**: Factual correctness of outputs
- **Relevance**: How well the output addresses the query
- **Coherence**: Logical consistency and readability
- **Groundedness**: Are claims supported by provided context?
- **Toxicity**: Presence of harmful or biased content
- **Latency**: Response time for real-time applications
- **Cost**: Token usage and API costs per query

## Cutting-Edge AI Research Areas (2025-2026)
- **Reasoning Models**: o1/o3-style chain-of-thought reasoning at inference time
- **Multimodal Agents**: Vision + Language + Action in unified models
- **Long Context**: 1M+ token context windows with efficient attention
- **Mixture of Experts**: Sparse activation for efficiency (Mixtral, Switch Transformer)
- **Constitutional AI**: Self-supervised alignment without human labels
- **Agentic AI**: Autonomous agents that plan, use tools, and self-correct
- **Synthetic Data**: Using AI to generate training data for AI
- **Test-Time Compute**: Spending more compute at inference for better reasoning
