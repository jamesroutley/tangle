# Array iteration in JavaScript

Let's look at three ways to iterate over an array in JavaScript. In each
example, we'll iterate over the array and print out each item.

Let's define the array to iterate over:

```javascript 1
let numbers = [1, 2, 3, 4, 5];
```

The first way is to keep track of an incrementing number `i`, which we use to
index into the array:

```javascript 2
for (let i = 0; i < numbers.length; i++) {
  console.log(numbers[i]);
}
```

Next, JavaScript has a `for...of` statement, which lets us iterate without
having to manage an index:

```javascript 3
for (let number of numbers) {
  console.log(number);
}
```

Finally, JavaScript arrays have a method `forEach`, which takes a function and
calls it for each item in the array:

```javascript 4
numbers.forEach(number => {
  console.log(number);
});
```
