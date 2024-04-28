# Measurements

This files contains the measurements of the execution times of the functions in this package.
It shows the differences between straight-forward and constant-time implementations.
As one can see the straight-forward implementations allow an attacker to infer the padding length from the execution times.

#### Measurements

Thr measurements were taken on the development machine.
It features an i7-10510U running at 1.80GHz.

All benchmarks were run 5 times in `Best performance` power mode.
An entry contains the minimum value of these 5 runs.
All measurements used a block size of 16 bytes, as it is the one most frequently used.
All times are nanoseconds per call.
There is always one entry for a padding of 1 byte and a padding of 15 bytes.

##### Straight-forward implementations

For comparison here are the execution times for the straight-forward implementations.

First the times for the `Pad` function with `PKCS#7` padding:

|   Type   | Time   |
|:--------:|:------:|
|  1 byte  |  101   |
| 15 bytes |  138   |

The times for 1 byte padding and 15 bytes padding differ by 37%.
They are clearly distinguishable.

Next there are the times for the `Unpad` function:

| Type                        | 1 byte  | 15 bytes  |
|-----------------------------|:-------:|:---------:|
| PKCS#7                      |   13    |    20     |
| X.923                       |   14    |    19     |
| ISO 10126                   |   13    |    13     |
| RFC 4303                    |   14    |    23     |
| ISO 7816-4                  |   13    |    24     |
| Arbitrary tail byte padding |   13    |    22     |

Here also the execution times for 1 and 15 bytes padding are clearly distinguishable, except for `ISO 10126` padding.

##### Constant-time implementations

Now follow the measurements for the constant-time implementations.
They feature nearly-constant execution times, so an attacker is unable to infer the padding length from the execution time.

First, there is a comparison of the constant-time `Pad` functions:

| Type                        | 1 byte | 15 bytes |
|-----------------------------|:------:|:--------:|
| Zero                        |  104   |   103    |
| PKCS#7                      |  108   |   107    |
| X.923                       |  102   |   102    |
| ISO 10126                   |  324   |   318    |
| RFC 4303                    |  142   |   140    |
| ISO 7816-4                  |  131   |   127    |
| Arbitrary tail byte padding |  150   |   153    |

As one can see the execution times for 1 byte and 15 bytes of padding are nearly identical.

Next, there is a comparison of the constant-time `PadLastBlock` functions:

| Type                        |1 byte | 15 bytes |
|-----------------------------|:-----:|:--------:|
| Zero                        |  48   |    48    |
| PKCS#7                      |  55   |    56    |
| X.923                       |  47   |    47    |
| ISO 10126                   |  229  |   227    |
| RFC 4303                    |  55   |    58    |
| ISO 7816-4                  |  47   |    46    |
| Arbitrary tail byte padding |  66   |    67    |

Again, the execution times for 1 byte and 15 bytes of padding are nearly identical.

At last, there are the execution times for the constant-time `Unpad` functions:

| Type                        | 1 byte | 15 bytes |
|-----------------------------|:------:|:--------:|
| Zero                        |   19   |    20    |
| PKCS#7                      |   22   |    23    |
| X.923                       |   22   |    23    |
| ISO 10126                   |   11   |    11    |
| RFC 4303                    |   22   |    23    |
| ISO 7816-4                  |   29   |    29    |
| Arbitrary tail byte padding |   18   |    20    |

Again, the execution times for 1 byte and 15 bytes of unpadding are nearly identical.
