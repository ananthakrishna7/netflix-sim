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

![[Pasted image 20250220144202.png]]

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

# Distributed System Concepts

## Non Blocking Synchronous Communication
The Amazon Server needs to communicate with the CDN and expects a reply for every message. But it would be wasteful for it to block when waiting for replies.

## Middleware
There does not seem to be an explicit middleware in the architecture, but content is created in various different formats to support multiple devices with different configurations. This can be thought of as converting 
## Global Snapshot
The system needs to know where each resource is located to generate a list of servers and select the best server for s specific client. Since the resources on the servers change dynamically, and hosts may fail, we need a Global Snapshot from time to time to keep track of resources. 
## Mutual Exclusion
Since Netflix has a distributed server on Amazon cloud, when one of the hosts is working with user data, the other hosts must not use it. Additionally changes here must be propagated to all hosts in case a distributed database is used. 
## Leader Election
Netflix has its own private CDN, so it has been able to optimize and simplify the software that runs on the CDN. All the compute intensive tasks run on the Amazon cloud, and the CDN handles only streaming. Hosts on the Amazon cloud are similar machines, hence there may be a need for leader election.

# Our Implementation

Since it would not be practical to simulate Netflix on a large scale, we plan to treat all hosts homogeneously, and run a leader election algorithm to elect a leader that coordinates our mini-CDN, and handles administrative tasks. 

We plan to implement this project in Go for the following reasons:
1. Support for concurrency
2. Support for HTTP servers
3. Support for MPEG-DASH protocol (open source library on GitHub)

We will implement the core features of the Netflix CDN:
1. Leader selection - Admin server selected through Ring algorithm leader election (Other servers must redirect requests here). We select the Ring algorithm because we do not want the Admin server to change if a failed process with higher rank comes back online. We want to minimize leader switching.
2. Global Snapshot Algorithm - We will implement the Chandy-Lamport Algorithm, because TCP guarantees reliable, in-order, error-checked delivery (TCP is the underlying Transport Layer Protocol for HTTP) and hence the network channel is a FIFO channel.
3. File processing and distribution - Video is uploaded to Admin server, and Admin server converts it into different formats and quality. These are distributed throughout the CDN.

# Sources
1. Netflix Open Connect Briefing Paper, 2021
2. Computer Networks, A Top Down Approach, Eighth Edition, James F. Kurose and Keith W. Ross
