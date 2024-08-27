```xml
<prompt_template>
    <instructions>Upon receiving a question and context, your role is to evaluate whether the provided context sufficiently answers the given question. If the answer can be derived from the context, respond with "PROBLEM-SOLVED". Conversely, if the context does not contain enough information to address the question, you should indicate this by responding with "FAILED". Ensure that no XML tags are included in your response.</instructions>
    <examples>
        <example>
            Question: What is the capital of France? Context: Paris is recognized as the capital and most populous city of France. Output: PROBLEM-SOLVED
        </example>
        <example>
            Question: Who discovered penicillin? Context: Antibiotics are medicines that combat bacterial infections; they can be natural, semi-synthetic, or synthetic. Output: FAILED
        </example>
        <example>
            Question: What are the primary colors of light? Context: In additive color systems like those used in displays and television screens, red, green, and blue lights combine to produce a broad spectrum of colors. Output: PROBLEM-SOLVED
        </example>
    </examples>
</prompt_template>
```
