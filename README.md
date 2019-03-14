In Januray of this year Stuart Marks posted a blog article named "Processing Large Files in Java". The blog post references an another post by Paige Niedringhaus about "Using Java to Read Really, Really Large Files". Niedringhaus there reports her experience with JavaScript to implement a "coding challenge" where a "very large file" has to be processed and four specific questiones where asked about the processed file. The mission statement of the challenge is the following:

- Write a program that will print out the total number of lines in the file.
- Notice that the 8th column contains a person’s name. Write a program that loads in this data and creates an array with all name strings. Print out the 432nd and 43243rd names.
- Notice that the 5th column contains a form of date. Count how many donations occurred in each month and print out the results.
- Notice that the 8th column contains a person’s name. Create an array with each first name. Identify the most common first name in the data and how many times it occurs.

The file is a freely available file about election donations provided by the U.S. Federal Elecetions Commission and is about 3.3GB in size uncompressed if taken at 4th of January 2019. 

After solving the challenge in JavaScript, Niedringhaus then moved forward and successfully implemented the same in Java as she was curious about Java and how to do it in that language.

In his article Marks takes the opportunity to analyze the Java version for the challenge to present variations and "focus on changing aspects of the computation to improve runtime performance" and also "present some interesting opportunities for optimizations and use of newer APIs". Overall, Marks establishes a baseline performance point, transforms the given Java programm and reduces the programs runtime from 108 seconds down to 32 seconds in the course of seven iterations.

After reading Marks article, I started to think how to implement that challenges solution in Go and how low can we go in terms of its running time. I tried to set a baseline for my environment by executing the existing implementations given by Niedringhaus and Marks and run them on my Laptop, a 2017 MacBookPro (Model 14,3, I7-7920HQ):

Niedringhaus, JavaScript with event-stream: 82s
Niedringhaus, Java: 84s
Marks, Java Variation 7: 23s

Then I started with Variation 7 of Marks implementation in Java and ported it straight forward to Go. 

Revision-0:
I ran it and it took about 38 seconds... or +60% of Variation 7. Thats Interesting.
 
<source of revision-0>

So what to do if a Go program is "slow"? Right! you get out the "Hammer of Go" and throw some "chunk of work" at all your CPU cores concurrently using go-routines and channels and try to beat the JAVA guy (I appologise for this little emotional outbreak).

Basically the challenge can be divided into:

a), counting lines
b), parse and collect:
- parse first names
- parse last names
- parse dates
and then 
c), processes the collected data:
- creating a frequency table of dates for donations
- find the most common first name
- find last names at three specific indexes within the data

Revision 10 => >720s, CPU ~80%

Why not take the "parse and collect" part into a go-routine and "send over" pared data over a channel where it would be collected to be processed after every lines was parsed. It took over 12 minutes this way, obiviously not a good idea.

What happens here: Revision 10 ensures that the CPU cores are highly saturated with "work". For every line we fire up a go-routine and these go routines might get distributed over to all cores, but it seems that the overhead of communication in relation to the small part of parsing a line into three data fields (firstname, lastname, date) is too high.

Revision 11.0 => 43s, CPU ~80%

So the "Hammer of Go" stupidly applied to a problem does not help as we learned it quite fast. With the next revision we try to reduce this overhead. Before we create one go-routine for every line we read, and send over every firstname, lastname and date over a dedicated channel, we could pack this data into a struct (entry) and send it over a channel to a collecting go-routine that unpacks the and appends them to their collecting lists of names, firstnames and dates.

With 110 seconds for this revision, we come back to a kind of normal runtime vs the 12 minutes with revision 10. This revision shows, that firing off a go-routine for every of these 18.2 mio lines of text is not completely unreasonable. There where still running 18.2 mio go-routines and then data was collected over one channel concurrently.

Revision 11.1 => 26s, CPU ~80%

Next, before starting a go-routine we build chunks of lines and then parse them together concurrently in a go-routine. Starting from 1k, 4k, 8k, 32k and 64k of lines, we come down from 43s down to ~26s of runtime.

Revision 12 => 11.8s, CPU ~80%

Instead of sending over every entry for every line, we collect a bunch of entries and then send them over to the concurrent go-routine to append the entries data to the global lists. We gain over 14 seconds in relation to revision 11.1, this is a lot.
One can see here clearly, that the communication overhead still is significant.

12 seconds now is nearly half the time Variation 7 in JAVA takes for the task. But it seems we can go faster.

Revision 13 => 14s, 

This revision removes the collecting go-routine and the channel and replaces the section where the entries get over the channel with a mutex to protect the appending code to the three lists for firstname, name and the dates.
The go-routines here most probably all will queue up before the mutex is available and then they can append their data to the lists.
It seems, to be proved, that this is less efficient than a dedidacted channel to collect the data.  
 
With the next revisions we peel off second by second by addressing specific aspects of Go: regex, sync.Pool and then rearrangements of code that helps to achive better performance. It seems to be impossible to gain -10s and go down to <4 seconds from here, right?
 

Go: 3.99s

[0] https://stuartmarks.wordpress.com/2019/01/11/processing-large-files-in-java/

https://github.com/paigen11/read-file-java/blob/27856f488cd61c6e13602a90293b4e752f92de1c/src/main/java/com/example/readFile/readFileJava/ReadFileJavaApplicationBufferedReader.java

https://itnext.io/using-node-js-to-read-really-really-large-files-pt-1-d2057fe76b33
https://itnext.io/streams-for-the-win-a-performance-comparison-of-nodejs-methods-for-reading-large-datasets-pt-2-bcfa732fa40e
https://itnext.io/using-java-to-read-really-really-large-files-a6f8a3f44649
https://github.com/paigen11/file-read-challenge

 
[1] ttps://www.fec.gov/files/bulk-downloads/2018/indiv18.zip