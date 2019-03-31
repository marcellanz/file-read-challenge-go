[![Build Status](https://dev.azure.com/marcellanz/file-read-challenge-go/_apis/build/status/marcellanz.file-read-challenge-go?branchName=master)](https://dev.azure.com/marcellanz/file-read-challenge-go/_build/latest?definitionId=1&branchName=master)

https://marcellanz.com/post/file-read-challenge/

>In January Stuart Marks published a blog post named ["Processing Large Files in Java"][marks_blog_post] as a response to a post by Paige Niedringhaus about ["Using Java to Read Really, Really Large Files"][niedringhaus_blog_post].
Niedringhaus there reports her experience with JavaScript to solve a "coding challenge" where a "very large file" has to be processed and four specific questions where asked about the processed file. After solving the challenge in JavaScript, Niedringhaus then moved forward and successfully implemented a solution to the challenge in Java as she was curious about Java and how to do it in that language.

>This article starts where Marks left and tries to improve on the performance aspect of the code further; until we [_hit the wall_][hit_the_wall].

[hit_the_wall]:https://marcellanz.com/post/file-read-challenge/#conclusion-hitting-the-wall

[marks_blog_post]: https://stuartmarks.wordpress.com/2019/01/11/processing-large-files-in-java/
[niedringhaus_blog_post]: https://itnext.io/using-java-to-read-really-really-large-files-a6f8a3f44649