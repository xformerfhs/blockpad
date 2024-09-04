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
| ISO 7816–4                  |   13    |    24     |
| Arbitrary tail byte padding |   13    |    22     |

Here also the execution times for 1 and 15 bytes padding are clearly distinguishable, except for `ISO 10126` padding.

##### Constant-time implementations

Now follow the measurements for the constant-time implementations.
They feature nearly-constant execution times, so an attacker is unable to infer the padding length from the execution time.

First, there is a comparison of the constant-time `Pad` functions:

| Type                        | 1 byte | 15 bytes |
|-----------------------------|:------:|:--------:|
| Zero                        |   66   |    66    |
| PKCS#7                      |   71   |    71    |
| X.923                       |   66   |    65    |
| ISO 10126                   |  222   |   214    |
| RFC 4303                    |   76   |    74    |
| ISO 7816–4                  |   70   |    70    |
| Arbitrary tail byte padding |   89   |   104    |
| Not byte padding            |   90   |   106    |

As one can see the execution times for 1 byte and 15 bytes of padding are nearly identical.

Next, there is a comparison of the constant-time `PadLastBlock` functions:

| Type                        | 1 byte | 15 bytes |
|-----------------------------|:------:|:--------:|
| Zero                        |   45   |    45    |
| PKCS#7                      |   51   |    51    |
| X.923                       |   44   |    44    |
| ISO 10126                   |  226   |   228    |
| RFC 4303                    |   51   |    51    |
| ISO 7816–4                  |   44   |    44    |
| Arbitrary tail byte padding |   65   |    66    |
| Not byte padding            |   49   |    50    |

Again, the execution times for 1 byte and 15 bytes of padding are nearly identical.

At last, there are the execution times for the constant-time `Unpad` functions:

| Type                        | 1 byte | 15 bytes |
|-----------------------------|:------:|:--------:|
| Zero                        |   21   |    23    |
| PKCS#7                      |   25   |    26    |
| X.923                       |   26   |    26    |
| ISO 10126                   |   13   |    13    |
| RFC 4303                    |   25   |    26    |
| ISO 7816-4                  |   34   |    44    |
| Arbitrary tail byte padding |   22   |    23    |
| Not byte padding            |   22   |    23    |

Again, the execution times for 1 byte and 15 bytes of unpadding are nearly identical.
