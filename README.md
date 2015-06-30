# holdingpattern

Holdingpattern talks to [skyapi](https://github.com/mediocregopher/skyapi/) and advertises a "service"
based on the command-line options/arguments. Makes testing remote endpoints extremely simple.

## Usage

```
holdingpattern [--api=127.0.0.1:8053] [--weight=100] [--priority=1] [hostname] [addr]
```
`hostname` and `addr` can be expressed using flags as well with `--hostname` and `--addr` respectively

## Example

Let's saying you're running skydns with a root as `example`.
```
$ holdingpattern test 127.0.0.1:8000
2015/06/29 22:59:33 Advertising [127.0.0.1:8053]: test on 127.0.0.1:8000 with priority 1 weight 100
```
Will add an SRV record as `test.services.example` pointing to `127.0.0.1` and port `8000`
```
$ dig @localhost SRV test.services.example

;; QUESTION SECTION:
;test.services.example.         IN      SRV

;; ANSWER SECTION:
test.services.example.  27      IN      SRV     100 100 8000 e4d29d5fcb01c98adb72.test.services.example.

;; ADDITIONAL SECTION:
e4d29d5fcb01c98adb72.test.services.example. 27 IN A 127.0.0.1
```
