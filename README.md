# CODE CHALLENGE: Computing hash Function

## The Use Case

> Build a program hat takes the given three arguments:
> 
> 1. a URL
> 2. a file name
> 3. a bandwidth for throttling (optional)
>
> And fetches the file given the URL,
> then computes its hash using the IMT hash function
> and stores the value in file under file name
> encoded as hex string.
>
> The IMT hash function is as follows:
>
> Length: 8 bytes
>
> ```
> coefficients := [8]int{ 2, 3, 5, 7, 11, 13, 17, 19 }
> for each incoming byte, ib:
>     for each byte of the hash, h
>         h[i] = ((h[i-1] + ib) * coefficient[i]) % 255
>         // in the case where i-1 == -1, h[i-1] should be 0.
> ```
> 
> For example, hashing the data:
>
> `data := []byte{12}`
>
> Should result in a hash of:
>
> `[]byte{24, 108, 90, 204, 81, 189, 102, 126}`
>
> Then, when converted to hexadecimal for writing to the output file, it should be:
> `186c5acc51bd667e`

## How to test

```sh
$ ./scripts/test.sh [-v]
```

## How to run the application

```sh
$ ./script/run.sh <url> <file> [bandwidth]
```
