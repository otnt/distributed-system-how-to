# Clock System

Clock system is a fundamental feature of a distributed system. The reason is a system would like to know if one event happens before another.

For example, assuming three machines (A, B, C) in a news system. Machine A send a piece of news to machine B and C. Machine B receives the news and send a comment to both machines A and C. Now machine C gets two events, i.e. a new piece of news and a new comment. How should it determine which one should come first? 

It seems extremely easy at first glance. But it turns out to be quite complicated in a real-world distributed system involving network delay, network partition etc.

In this article, we will use the former example (news system) throughout the passage, introducing several types of clock services in practice, as well as discussing about their pros and cons.

## TL;DR

We will talk about two basic clock systems: logical clock and vector clock. Logical clock is very simple and synchronizing using logical clock hardly add any burden to network. Vector clock results in much more network overhead, but it is more capable to order things.

## Organization

This article will be organized as following:

**First**, we introduce message passing model in distributed system. 

**Then**, two important concepts -- concurrent, happens before -- would be introduced. 

**Next**, we will introduce a conceptually simple clock service called Logical Clock/Lamport Clock, which is useful in most of systems despite of its simplicity. 

**Finally**, another type of clock service called Vector Clock is introduced. This is perhaps the most famous and widely used clock service in practice.

## Message Passing Model

Communication among distributed systems could not be synchronized. This is to say, though intra-machine communication (within a single machine) could be synchronized, all inter-machine communication (among multiple machines) are generally asynchronized. 

One communication model is message passing model. So that machine A talks with machine B by sending and receiving information directly through network, not by sharing address space or other methods.

In message passing model, all jobs happend could be put to three categories.

<img src="https://raw.githubusercontent.com/otnt/distributed-system-notes/master/clockService/img/send_receive_event.png" align="right" width="400" alt="Message Passing Model">

1. **Send**: Machine A sends a message to machine B.
2. **Receive**: Machine A receives a message from machine B.
3. **Event**: Machine A triggers an inner job to happen.

The last job, i.e. event, could refer to many things, including user input, open a file, system reboot etc.

## Two important concepts

Two fundamental concepts in distributed system clock service is *concurrent* and *happens before*. The interpretation is really straightforward.

Let's say job j1 has time t1, and job j2 has time t2. Then j1 is ***concurrent*** with j2 if t1 = t2. And j1 ***happens before*** j2 if t1 < t2.

## Logical Clock / Lamport Clock

#### Idea
The key concept of logical clock is that it applies clock service from a single user/a single machine, to a global view. Instead of each machine maintains its own clock locally, now the whole system together maintain a clock service. By doing this, all machines could come to consensus when ordering jobs.

Each machine uses a counter as clock locally. The counter indicates the sequence of jobs. Machines increment their clock when any of send, receive or event happens. They synchronize their clocks when sending and receiving message to/from other machines. The counter could only increment, and job that causally happens later could not have lower counter number.

#### Implementation

To be specific, each machine has counter initialized as zero. It acts as below:

<img src="https://raw.githubusercontent.com/otnt/distributed-system-notes/master/clockService/img/logical_clock.png" align="right" width="400" alt="Logical Clock">

1. **Send**: It attaches current time in the message, sends it, and then increment local time.
2. **Receive**: It receive a message, updates local time to the time attached with the message if it's larger, and then increment local time.
3. **Event**: It increment local time.

For example, in right image, even network between A and C is conjusted, so that message sent from A to C is delayed, C still knows A's message should be received before B's message.

#### Pros and Cons

Pros:

1. It's conceptually simple.
2. It guarantees that things causally happen before, would not have larger time.
3. It saves network bandwidth.

Cons:

1. We could not determine which job happens first by comparing the time. In otherword, if job j1 has smaller time than job j2, it could be either j1 is concurrent with j2 or j1 happens before j2. But we are not able to determine which one is correct.

## Vector Clock

Coming soon.

## Conclusion

Coming soon.


> Written with [StackEdit](https://stackedit.io/).
