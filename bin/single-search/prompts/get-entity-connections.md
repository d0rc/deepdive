```xml
<template>
    <instructions>Given a List of Entities, Context (in form of unstructured text or an article), and The Entity Of Interest, your task is to identify connections between The Entity Of Interest and other entities from the list provided. These connections should be described using natural language and organized into a JSON map. This map will contain lists of each type of connection you find. Ensure that the output does not include any XML tags or elements; it should only consist of valid JSON data. Follow these steps to complete the task:
        <ol>
            <li>Read through the Context provided to understand the relationships between entities.</li>
            <li>Determine all connections for The Entity Of Interest with other entities listed in List of Entities based on information from the Context.</li>
            <li>Classify each connection into types, such as 'is a', 'part of', 'related to', etc., describing how The Entity Of Interest is connected to another entity.</li>
            <li>Create lists for different types of connections found in step 3 and organize these lists into a JSON map where keys are the types of connections and values are lists of entities that share those connections with The Entity Of Interest.</li>
            <li>Output the resulting JSON map without including any XML tags or elements.</li>
        </ol>
    </instructions>
    
    <examples>
        <example>
            Input:
                List of Entities: ['John Doe', 'Jane Smith', 'CityX']
                Context: John Doe is a resident of CityX and Jane Smith's brother. Jane Smith works at the local library in CityX.
                The Entity Of Interest: 'John Doe'
            Expected Output:
                {
                    "is a": ["resident"],
                    "related to": ["Jane Smith"],
                    "lives in": ["CityX"]
                }
        </example>
        
        <example>
            Input:
                List of Entities: ['Apple', 'Banana', 'Orange']
                Context: Apple and Orange are fruits found on trees, while Banana grows from a plant that is not considered a tree.
                The Entity Of Interest: 'Apple'
            Expected Output:
                {
                    "is": ["fruit"],
                    "found on": ["trees"]
                }
        </example>
        
        <example>
            Input:
                List of Entities: ['Paris', 'France', 'Eiffel Tower']
                Context: Paris is the capital city of France and home to the iconic Eiffel Tower.
                The Entity Of Interest: 'Paris'
            Expected Output:
                {
                    "is": ["capital"],
                    "part of": ["France"],
                    "home to": ["Eiffel Tower"]
                }
        </example>
    </examples>
</template>
```
