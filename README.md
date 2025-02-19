# netflix-sim
Simulating Netflix's CDN to learn about Distributed Systems

# The Netflix CDN
Netflix is a popular video streaming application. It has a vast collection of titles for streaming. In 2020, Netflix spent USD 12.5 Billion on content. Video streaming represents 56% of high speed internet traffic, and Netflix accounts for 11% of this traffic. [^1]

It would be inefficient and extremely expensive to have a centralized system (even a collection of servers, like AWS Cloud) deliver terabytes of content to customers all over the world.
Streaming is best achieved with a Content Distribution Network (CDN), where the content is brought closer to the customer.

>90% of internet traffic destined for consumers is carried by CDNs, and over USD 75 Billion is invested annually into infrastructure by online service providers to bring content closer to consumers [^1]
 
Netflix owns and operates its own private CDN. This is how it is able to reach customers all over the world. It places

# A Brief Description of Content Delivery

Netflix partners with ISPs around the globe through its Open Connect program. It installs its server racks at IXPs or Within ISPs.

## Uploading a Video

1. Before Netflix distributes content, it is first ingested and processed by the origin servers, which are hosts in the Amazon Cloud.
2. The content is processed and copies in many different formats are created. This is to support a diverse array of devices. For each of these formats, versions with different bitrates are created, to support adaptive streaming using *Dynamic Adaptive Streaming over HTTP (DASH)*
3. Once all versions are created, they are rolled out to the CDN.

### How data is pushed to the CDN
- Netflix CDN servers do not use pull-caching
- Content is pushed to the CDN during off-peak hours
- If the location cannot hold the entire Netflix movie library, only popular content is pushed(determined on a day to day basis)
## Streaming a Video

The Web Pages for browsing Netflix are stored in the Amazon Cloud.
1. A user selects a video on the Netflix website
2. Netflix software determines which of its CDN servers have copies of the movie
3. The best CDN server is selected, and this IP address is passed to the client, along with a manifest file containing the URLs for different quality versions.
4. The client and server interact directly with DASH (a proprietary version of DASH)
5. Client measures throughput, runs a rate determination algorithm and determines the quality of the next chunk to request.

%% # Distributed System Concepts %%




[^1]: Netflix Open Connect Briefing Paper, 2021
