```xml
<prompt_template>
    <!-- The following instructions describe how to create a list of questions that explore different angles of an original question, without referencing it directly. -->
    <instructions>
        When you receive a question from the user, your goal is not to answer this question but to interpret its essence and generate new questions that could provide deeper insights or help in finding a perfect answer. Ensure each created question is self-contained; no reference to the original question should be made. Instead, include all necessary context within the generated questions themselves.
        Structure your output as a JSON list of strings, where each string represents an individual question. Remember, the aim here is to look at the topic from various perspectives and offer questions that could shed more light on it or provide hints leading to a better understanding.
        Always use original language of the Question.
        
        Avoid using any XML tags in your response and ensure all generated questions are clear, concise, and relevant to the original topic.
    </instructions>

    <!-- Below are examples demonstrating how to apply the instructions -->
    <examples>
        <!-- Example 1 -->
        <example>
            The user's question might be about "The history of pizza". A proper output could include:
            <quote>
                [
                "What were the earliest forms of pizza and their origins?",
                "How has the preparation and ingredients of pizza evolved over time?",
                "What are some significant cultural impacts that pizza has had globally?"
                ]
            </quote>
        </example>

        <!-- Example 2 -->
        <example>
            If asked about "Best practices for remote work", you could create questions such as:
            <quote>
                [
                "How do successful companies manage team communication in a fully remote setting?",
                "What are the key differences between managing remote and onsite teams?",
                "Which tools have proven most effective in enhancing productivity during remote work?"
                ]
            </quote>
        </example>

        <!-- Example 3 -->
        <example>
            For a question on "Impact of social media on mental health", potential outputs include:
            <quote>
                [
                "How do frequent users of social media report their emotional state compared to non-users?",
                "What role does the frequency and type of content consumed play in influencing mental health outcomes on social media platforms?",
                "Are there specific age groups or demographics more susceptible to negative impacts from social media exposure?"
                ]
            </quote>
        </example>
    </examples>

</prompt_template>
```
