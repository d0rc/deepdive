```xml
<instructions>
    To complete this task, you need to summarize the provided search engine output while addressing a specific question. The input will consist of three components:
    1. A question that the user is seeking an answer for.
    2. Keywords used for web searching, which can help understand the query's context better.
    3. An output list from a search engine in YAML format containing triplets (url, title, content), where each triplet represents a search result.

    Your task involves analyzing these triplets and summarizing them with respect to the user’s question. It is important to note that some results may be misleading or irrelevant; therefore, ensure your summary focuses on relevant information only and includes URL references.

    The output should not contain any XML tags, but instead present findings in plain text format.
</instructions>

<example>
If the question was "What are the health benefits of green tea?", keywords included 'healthy', 'green tea' and search engine results contained articles discussing various aspects such as weight loss, antioxidants properties etc., your summary should focus on how these points answer or contribute to answering user's initial query about health benefits. Ensure each relevant point is supported by a citation through url.

For instance:
Input: Question - "What are the health benefits of green tea?"
Keywords - 'healthy' , 'green tea'
Search Engine Output -
[
    (url1, title1, content1),
    (url2, title2, content2),
    (url3, title3, content3)
]
Output: Green tea has several health benefits including its role in weight loss and high antioxidant content. According to [url1], green tea aids metabolism contributing positively towards healthy weight management. Another source at [url2] highlights that regular consumption of green tea can boost the body's antioxidant levels which protect cells from damage caused by free radicals.
</example>

<example>
Another example could involve searching for "Best practices in cybersecurity" where keywords might be 'cybersecurity', 'best' and results range across different domains like government recommendations, industry standards etc. Your summary should identify key points from these sources while linking them back to original articles.

Input: Question - "What are best practices in cybersecurity?"
Keywords - 'cybersecurity', 'best'
Search Engine Output -
[
    (url1, title1, content1),
    (url2, title2, content2),
    (url3, title3, content3)
]
Output: Best practices in cybersecurity include implementing strong password policies and training employees on phishing awareness. A government report at [url1] suggests updating software regularly to prevent vulnerabilities. Industry guidelines from [url2] emphasize the importance of multi-factor authentication for accessing sensitive information.
</example>

<example>
In an inquiry about "Recent developments in AI research", keywords could be 'AI', 'research' and results may vary widely covering topics such as new algorithms, ethics discussions etc. Your summary should highlight recent trends or breakthroughs while referencing specific sources.

Input: Question - "What are the recent developments in AI research?"
Keywords - 'AI', 'research'
Search Engine Output -
[
    (url1, title1, content1),
    (url2, title2, content2),
    (url3, title3, content3)
]
Output: Recent developments in AI research include advancements in natural language processing and ethical considerations for algorithmic biases. A study published at [url1] explores improvements in machine translation accuracy using neural networks. Another article from [url2] discusses the need for transparent AI models to ensure fairness.
</example>
```