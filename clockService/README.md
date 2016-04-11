# Clock System

Clock system is a fundamental feature of a distributed system. The reason is a system would like to know if one event happens before another.

For example, assuming three nodes (A, B, C) in a news system. Node A send a piece of news to nodes B and C. Node B receives the news and send a comment to both nodes A and C. Now node C gets two events, i.e. a new piece of news and a new comment. How should it determine which one should come first? 

It seems extremely easy at first glance. But it turns out to be quite complicated in a real-world distributed system involving network delay, network partition etc.

In this article, we will use the former example (news system) throughout the passage, introducing several types of clock services in practice, as well as discussing about their pros and cons.

## TL;DR

We will talk about two basic clock systems: logical clock and vector clock. They have tradeoff in network overhead, ability to determine occurrence sequence of events and implementation complexity.

## Organization

This article will be organized as following:

**First**, we introduce message passing model in distributed system. **Then** two important concepts -- concurrent, happens before -- would be introduced. **Next**, we will introduce a conceptually simple clock service called Logical Clock/Lamport Clock, which is useful in most of systems despite of its simplicity. **Finally**, another type of clock service called Vector Clock is introduced. This is perhaps the most famous and widely used clock service in practice.

## Message Passing Model

Communication among distributed systems could not be synchronized. This is to say, though intra-system communication could be synchronized, all inter-system communication is generally asynchronized. 

One communication model is message passing model. So that system $A$ talks with system $B$ by sending and receiving information through network, not by sharing address space or other methods. The  model doesn't care the implementation of network transportation.

There are three major things that could change the time.
**Send**: System $A$ sends a message to system $B$.
**Receive**: System $A$ receives a message from system $B$.
**Event**: System $A$ triggers an inner job to happen.

The last thing -- event -- could be bound to many things, including user input, open a file, system reboot etc.

## Two important concepts

Two fundamental concepts in distributed system clock service is *concurrent* and *happens before*. The interpretation is really straightforward.

Let's say job $J_1$ has time $t_1$, and job $J_2$ has time $t_2$. Then $J_1$ is ***concurrent*** with $J_2$ if $t_1=t_2$. And $J_1$ ***happens before*** $J_2$ if $t_1<t_2$.

## Logical Clock / Lamport Clock

#### Idea
The key concept or breakthrough of logical clock, I believe, is that it release clock service from a single user/a single system, to a global view. Instead of each sub-system maintain clock locally, now the whole system together maintain a clock service. By doing this, all sub-systems could come to consensus when ordering events.

Each sub-system uses a counter as clock locally. The counter indicates the sequence of jobs. System increment their clock when any of send, receive or event happens. They synchronize their clocks when sending and receiving message to/from other systems. The counter could only increment, and job that happens later should have higher counter number.

#### Implementation

To be specific, each sub-system has counter initialized as zero. It acts as below:
1. Send: It attaches current time in the message, sends it, and then increment local time.
2. Receive: It receive a message, updates local time to the time attached with the message if it's larger, and then increment local time.
3. Event: It increment local time.

You could see the implementation here at GitHub: [Logical Clock](https://github.com/otnt/distributed-system-how-to/blob/master/clockService/logicalClock.go)

#### Pros and Cons

Pros:
1. It's conceptually simple.
2. It guarantees that things causally happen before, would not have larger time.
3. It saves network bandwidth.

Cons:
1. We could not determine which job happens first by comparing the time.

## Vector Clock

Coming soon.

## Conclusion

Coming soon.


> Written with [StackEdit](https://stackedit.io/).
