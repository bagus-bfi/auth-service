# README

Contains ONLY data store interface (repository interface) along with table entity struct or per use-case entity struct.

As we are trying not to limit how a team access the data store, we provide the following example of the implementation of data store access:

1. Using direct SQL using plain SQL query and vanilla Go. See [`pg` folder](./pg)
2. Using query builder and struct auto populate. See [`pgwithqb` folder](./pgwithqb)
3. Using ORM. See [`pgwithorm` folder](./pgwithorm)

But please keep in mind that if you choose to add an external library for data access,
make sure it does not have many dependencies or make the team member learn lots of libraries or frameworks to achieve their goal.

We recommend #1 or at most #2. Fewer dependencies on external libraries are better.