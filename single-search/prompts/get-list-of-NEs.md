```xml
<prompt-template>
    <!-- This is a template to instruct an AI system to identify named entities from provided text, then return them as a parsable JSON list of strings without any additional information or XML tags -->

    <instructions>
        When you are presented with a piece of context in the form of textual data,
        your task is to identify all the named entities present within this context.
        A named entity refers to an object that has a specific name, typically including (but not limited to)
        names of people, places, organizations, products, and dates.

        Here's how you should proceed:
        1. Analyze the provided context for any mentions of named entities.
        2. Ensure you are identifying all types of named entities: Persons, Organizations, Locations,
        Products, Dates or other specific names that might be relevant to the given text.
        3. Format your output by creating a JSON list that only contains strings (the identified named entities).
        - The JSON list should omit any additional information such as explanations or justifications for inclusions.
        - Do not include any XML tags, comments, or extraneous data in your response; the output must be strictly
        a parsable list of strings.

        Remember, it's crucial that you only return named entities and nothing else to maintain the integrity of the task.
    </instructions>

    <examples>
        <!-- Example 1 -->
        <example>
            Provided context: "Apple Inc., headquartered in Cupertino, California, was founded by Steve Jobs."

            Your output should be:
            [
            "Apple Inc.",
            "Cupertino",
            "California",
            "Steve Jobs"
            ]
        </example>

        <!-- Example 2 -->
        <example>
            Given context: "The Eiffel Tower in Paris is a world-famous landmark, designed by Gustave Eiffel."

            Your output should look like:
            [
            "Eiffel Tower",
            "Paris",
            "Gustave Eiffel"
            ]
        </example>

        <!-- Example 3 -->
        <example>
            Provided text: "After visiting the Vatican City, Pope Francis gave a speech at St. Peter's Square."

            Expected output:
            [
            "Vatican City",
            "Pope Francis",
            "St. Peter's Square"
            ]
        </example>
    </examples>

</prompt-template>
```
