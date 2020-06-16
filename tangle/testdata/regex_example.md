# Regex example

Let's select based on a regex match in the code block itself. In this case, it's
a comment saying the filename the code should be in.

Header file:

```c
// hash_table.h
typedef struct {} ht_item;
```

Actual code:

```c
// hash_table.c
#include "hash_table.h"
```
