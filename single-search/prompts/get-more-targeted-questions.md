```xml
<prompt_template>
    <!-- This is a structured prompt for an AI to create a list of targeted, narrow questions related to a given main question. The goal is to provide additional insights that can help in finding more accurate and comprehensive answers to the original query. -->
    
    <instructions>
        You are tasked with generating a series of follow-up questions based on the provided main question. Your aim is to craft questions that:
        1. Do not directly reference or paraphrase the original question.
        2. Include sufficient context so each generated question can stand alone.
        3. Are aimed at uncovering more detailed information about entities, persons, events related to the topic of the main question.
        Your output should be a JSON list containing these follow-up questions. Ensure that your response does not contain any XML tags and is purely in JSON format.
    </instructions>
    
    <!-- Below are examples demonstrating how an input question would lead to several pertinent sub-questions. -->
    <examples>
        <example>
            <input>What led to the fall of the Roman Empire?</input>
            <output>["How did economic factors contribute to the decline of the Roman Empire?", "Which external invasions played a significant role in the fall of Rome?", "In what ways did political instability within the Roman Empire accelerate its collapse?"]</output>
        </example>
        
        <example>
            <input>What is the significance of DNA in genetics?</input>
            <output>["How do mutations in DNA affect genetic traits and diseases?", "Can you explain the role of DNA replication in cell division?", "What are the implications of DNA fingerprinting in forensic science?"]</output>
        </example>
        
        <example>
            <input>Who was Leonardo da Vinci?</input>
            <output>["What were some of the major contributions made by Leonardo da Vinci to the field of art during the Renaissance?", "Can you describe any significant scientific or technological inventions attributed to Leonardo da Vinci?", "In what ways did Leonardo da Vinci's work in anatomy influence his artistic endeavors?"]</output>
        </example>
    </examples>
    
</prompt_template>
```
