# Notes for distributed system beginners

## What is this repo about?

This repo is a collection of notes and implementations of some classic distributed system algorithms.

The goal (for now) is to share my basic understanding of them (notice I'm also a beginner of distributed system). Since designing a distributed system algorithm that is correct, efficient and easy-to-implement at the same time is difficult, I will focus on correct and easy-to-implement.

## How to use this repo?

Most of topics are independent. So you could randomly view any topic of interest.

I'm also working on writing some topic collections, for example how to build distributed hash table by combining clock service, consistent hashing, membership protocol, failure detector etc. together. This would be some high-level study, and could form a small project that benefits beginner students.

## But, Pufan, there are already TONS of material over the Internet

Yes, of course, and many of them have better explanation than my repo. But sometimes I have a different perspective of understanding these algorithms. Also, a combination of both theory and implementation for BEGINNER is not easy to find. These are why I believe my repo holds its value among all those blogs. 

## An Informal Short Introduction

Distributed system is where machines communicate to each other through some method that does not provide bounded time delivery. For example, the most typical communication method is by network. However, nobody could guarantee how long would a message sent from US West Coast to China. The immediate result is, before a message is known to be sent (e.g. by an acknowledgement reply), it is *impossible* to know if the message sending is failed, or it is still on the way. Due to fail-unknown and asynchronous nature, programming on a single machine and on a cluster are hugely different.

The high-level goals of distributed system algorithms are the same as other computer algorithms, i.e. **correct**, **efficient** and **easy to implement**.

Correctness is the basic of any algorithm. In distributed system, correctness guarantee is much harder because there's no *global overview* for the whole cluster. In other word, machines could not know the state (running or faulty) of others in a deterministic way, such as using system bus on motherboard to let CPU knows the state of other components. Another reason is when distributed system contains hundreds of or even thousands of machines, the possibility of machine faulty is rising significantly. In all, one has to bear in mind that machine could fail at any time, as well as machine failure is normal, to design good distributed system algorithms.

Efficience guarantee is hard in distributed system. Since lacking of global overview, in order to reach some extend of consensus, typical solutions are to let machines communicate with each other multiple times. But this would result in network overhead, and number of times of communication dominates latency in distributed system. Another reason is distributed system is largely about scalability. Intuitively, increasing machines add up to more computing power. However, it is easy to have an $O(N^2)$ algorithm (obviously not scalable) and the more machines are involved, the less efficient the algorithm is.

Finaly, implementing a distributed system algorithm is hard. One reason is there are much more components involved in distributed system than single machine, so that the algorithm has to be more complicated to be fault-tolerant and robust. Also, the asynchronous nature of distributed system is counter-intuitive to human thought, which often makes the algorithm flow less clearer in our mind.

Despite the challenge of distributed system design, numerous incredible solutions are proposed and it is already the core part in machine world. Notice the Internet, without something you would go mad, is one of the largest, most well-designed distributed system. So you have already benefited from it, why not contribute to it as well?

## How to contribute?

I'm always welcome to any contribution. The best way would be post an issue or pull request. I'm still new to open source world, so let's see what is best pratice for contributing to this repo.

