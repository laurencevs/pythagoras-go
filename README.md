## Pythagorean triangle finder

This program finds rational Pythagorean triangles of given area, using elliptic curves.

```
$ ./pythagoras-go 6
Finding rational Pythagorean triangles with area 6
(3, 4, 5)
(7/10, 120/7, 1201/70)
(4653/851, 3404/1551, 7776485/1319901)
(1437599/168140, 2017680/1437599, 2094350404801/241717895860)
(3122541453/2129555051, 8518220204/1040847151, 18428872963986767525/2216541307731009701)
```

### How it works

The program works in two stages:

1. Brute force search to find at least one initial Pythagorean triangle of area $n$.
   - This uses [Euclid's parametrisation](https://en.wikipedia.org/wiki/Pythagorean_triple#Generating_a_triple) of Pythagorean triples.
2. Elliptic curve arithmetic to find arbitrarily many other Pythagorean triangles of area $n$.

By default, stage 1 will exit as soon as one initial triangle is found. For some values of $n$ it is possible to find multiple initial trianges; to do this, use the `timeout` option.

```
$ ./pythagoras-go --timeout=1s --max-digits=5 210
Finding rational Pythagorean triangles with area 210
(21, 20, 29)
(41/58, 24360/41, 1412881/2378)
(35, 12, 37)
(15/2, 56, 113/2)
(495/8, 224/33, 16433/264)
(1925/171, 2052/55, 366517/9405)
(1081/74, 31080/1081, 2579761/79994)
(819/187, 3740/39, 700109/7293)
```

### Options

| Name | Default | Description |
| - | - | - |
| `max-coeff` | `5` | The maximum coefficient by which to multiply each elliptic curve point. |
| `max-digits` | `0` | The maximum number of digits in the denominator of the hypotenuse. Triangles with larger hypotenuse denominators will not be outputted. |
| `timeout` | `0` | The timeout length for the initial point search. If set to 0, the program will keep searching until it finds one point. |
| `output` | `""` | The filename for output. By default, the program will print to stdout. |

### Background

A [congruent number](https://en.wikipedia.org/wiki/Congruent_number) is a number that is the area of a Pythagorean triangle with rational side lengths. For example, 6 is the area of the triangle with side lengths (3, 4, 5), so 6 is a congruent number.

1, 2, 3, and 4 are not congruent numbers, but 5, 6, and 7 are. As of January 2025, the problem of which numbers are congruent numbers is unsolved.

It turns out that finding rational Pythagorean triangles of area $n$ is equivalent to finding rational points with $y \neq 0$ on the [elliptic curve](https://en.wikipedia.org/wiki/Elliptic_curve) $E_n$. Thus, given a rational Pythagorean triangle of area $n$, we can use the [elliptic curve group law](https://en.wikipedia.org/wiki/Elliptic_curve#The_group_law) to find arbitrarily many other Pythagorean triangles of area $n$.
