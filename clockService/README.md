# Clock System

Clock system is a fundamental feature of a distributed system. It is straightforward that a system would like to know if one event happens before another. 

For example, assuming three nodes (A, B, C) in a news system. Node A send a piece of news to nodes B and C. Node B receives the news and send a comment to both nodes A and C. Now node C gets two events, i.e. a new piece of news and a new comment. How should it determine which one should come first? 

It seems extremely easy at first glance. But it turns out to be quite complicated in a real-world distributed system involving network delay, network partition etc.

In this article, we will use the former example (news system) throughout the passage, introducing several types of clock services in practice, as well as discussing about their pros and cons.

## TL;DR

We will talk about these clock system: logical clock, vector clock, bounded vector clock and dependency clock. They have tradeoff in network overhead, ability to determine occurrence sequence of events, implementation complexity, and memory overhead.

## Organization

This article will be organized as following:

First, we will briefly talk about why a traditional clock service doesn't work in distributed system.

Second, we will introduce a conceptually simple clock service called Logical Clock/Lamport Clock, which is useful in most of systems despite of its simplicity.

Then, another type of clock service called Vector Clock is introduced. This is perhaps the most famous and widely used clock service in practice.

Finally, two more advanced clock service, namely Bounded Vector Clock and Dependency Clock is introduced. These are designed to overcome the shortcoming comes with Vector Clock, but with tradeoff of either ability to determine sequence of occurrence events, or implementation complexity.


## Conclusion



> Written with [StackEdit](https://stackedit.io/).
