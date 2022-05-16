# q

An attempt at building out the lock-free, concurrent, generic queue in 32 bits described [here](https://nullprogram.com/blog/2022/05/14/).

We might not get it in 32 bits, but let's see.

## Caveats, Notes, etc.

This is probably not ready for production. There are a lot of caveats for working with it (see the package docs for more information).

This is not meant as a "better channel" implementation, it is very much not an implementation of a channel, and it barely does anything that channels do.

This is 100% a learning project to see whether I could implement the design found in the link above, along with seeing how to squeeze out every inch of performance that I can from it.