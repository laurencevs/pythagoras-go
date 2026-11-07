## Pythagorean triangle finder

This program finds rational Pythagorean triangles of given area, using elliptic curves.

```
$ ./pythagoras-go -verbose 6
Finding rational Pythagorean triangles with area 6 (square-free part 6)
Tunnell's criterion passed, expect some triangles
Brute force search found initial point (12, 36)
(3, 4, 5)
(7/10, 120/7, 1201/70)
(3404/1551, 4653/851, 7776485/1319901)
(2017680/1437599, 1437599/168140, 2094350404801/241717895860)
(3122541453/2129555051, 8518220204/1040847151, 18428872963986767525/2216541307731009701)
```

If [Tunnell's theorem](https://en.wikipedia.org/wiki/Tunnell%27s_theorem) shows that no such triangle exists, the program exits. Assuming the Birch & Swinnerton-Dyer conjecture, this means the program will terminate for all areas.

```
$ ./pythagoras-go 30001
Tunnell's criterion failed (2 * 244 != 496), ending search
```

### How it works

The program works in two stages:

1. Search to find at least one initial Pythagorean triangle of area $n$, via two concurrent methods:
   - Brute force search over [Euclid's parametrisation](https://en.wikipedia.org/wiki/Pythagorean_triple#Generating_a_triple) of Pythagorean triples.
   - Descent via 2-isogeny.
2. Elliptic curve arithmetic to find arbitrarily many other Pythagorean triangles of area $n$.

Descent via 2-isogeny allows us to find triangles which would be hard to find by Euclid's parametrisation alone:

```
$ ./pythagoras-go -verbose --max-coeff=1 199
Finding rational Pythagorean triangles with area 199 (square-free part 199)
Tunnell's criterion passed, expect some triangles
Descent found initial point (14207630028601/68644524001, 14721540674529048360/17984933932786001)
(3905642100360/987560631301, 196524565628899/1952821050180, 194229730921552615921952401/1928529189133642599684180)
```

By default, stage 1 will exit as soon as one initial triangle is found. For some values of $n$ it is possible to find multiple initial triangles; to do this, use the `timeout` option.

```
$ ./pythagoras-go --timeout=1s --max-digits=5 210
(20, 21, 29)
(41/58, 24360/41, 1412881/2378)
(12, 35, 37)
(15/2, 56, 113/2)
(224/33, 495/8, 16433/264)
(1925/171, 2052/55, 366517/9405)
(1081/74, 31080/1081, 2579761/79994)
(819/187, 3740/39, 700109/7293)
```

### Options

| Name | Default | Description |
| - | - | - |
| `max-coeff` | `5` | The maximum coefficient by which to multiply each elliptic curve point in stage 2. |
| `max-digits` | `0` | The maximum number of digits in the denominator of the hypotenuse. Triangles with larger hypotenuse denominators will not be outputted. |
| `timeout` | `0` | The timeout length for the initial point search. If set to 0, the program will keep searching until it finds one point. |
| `output` | `""` | The filename for output. By default, the program will print to stdout. |
| `verbose` | `false` | Whether to include verbose output. |

### Background

A [congruent number](https://en.wikipedia.org/wiki/Congruent_number) is a number that is the area of a Pythagorean triangle with rational side lengths. For example, 6 is the area of the triangle with side lengths (3, 4, 5), so 6 is a congruent number.

1, 2, 3, and 4 are not congruent numbers, but 5, 6, and 7 are. As of November 2026, the problem of which numbers are congruent numbers is unsolved.

It turns out that finding rational Pythagorean triangles of area $n$ is equivalent to finding rational points with $y \neq 0$ on the [elliptic curve](https://en.wikipedia.org/wiki/Elliptic_curve) $$E_n: y^2 = x^3 - n^2 x.$$ Thus, given a rational Pythagorean triangle of area $n$, we can convert it to a point on $E_n$ and then use the [elliptic curve group law](https://en.wikipedia.org/wiki/Elliptic_curve#Group_law) to find arbitrarily many other Pythagorean triangles of area $n$.
