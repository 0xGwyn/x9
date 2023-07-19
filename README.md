# Installation
```bash
go install github.com/0xgwyn/x9@latest
```
or
```bash
git clone https://github.com/0xGwyn/x9.git 
cd x9
go build -o $GOPATH/bin/x8 main.go
```
# Usage 
```bash
x9 -h
```
This will display help for the tool. Here are all the switches it supports.
```
Usage:
  x9 [flags]

Flags:
  -l, -list string                  List of URLS to edit (stdin could be used alternatively)
  -p, -parameters string            Parameter wordlist
  -c, -chunk int                    Number of parameters in each URL (default 15)
  -v, -value string[]               Value for the parameters
  -gs, -generate-strategy string[]
                                    Select the mode strategy from the available choices:
                                        normal:  Remove all parameters and put the wordlist
                                        combine: Pitchfork combine on the existing parameters
                                        ignore:  Don't touch the URL and append the parameters to the URL

  -vs, -value-strategy string       Select the strategy from the available choices:
                                        replace: Replace the current URL values with the given values
                                        suffix:  Append the value to the end of the parameters
                                        (default "suffix")
  -o, -output string                File to write output results
  -de, -double-encode               Double encode the values
```

# Examples
paramFile contents:
```
param1
param2
...
param10
```

### Normal mode 
Replaces all the values with the new provided values
```bash
echo "https://domain.tld/test/rot?pa1=val1&pa2=val2" | x9 -p paramsFile -gs normal -v newVal1 -v newVal2 -c 7

https://domain.tld/test/rot?pa1=newVal1&pa2=newVal1&paramP10=newVal1&paramP6=newVal1&paramP7=newVal1&paramP8=newVal1&paramP9=newVal1
https://domain.tld/test/rot?pa1=newVal1&pa2=newVal1&paramP1=newVal1&paramP2=newVal1&paramP3=newVal1&paramP4=newVal1&paramP5=newVal1
https://domain.tld/test/rot?pa1=newVal2&pa2=newVal2&paramP10=newVal2&paramP6=newVal2&paramP7=newVal2&paramP8=newVal2&paramP9=newVal2
https://domain.tld/test/rot?pa1=newVal2&pa2=newVal2&paramP1=newVal2&paramP2=newVal2&paramP3=newVal2&paramP4=newVal2&paramP5=newVal2
```


### Ignore mode 
Leaves default params and their values unchanged and adds new params with the new values
```bash
echo "https://domain.tld/test/rot?pa1=val1&pa2=val2" | x9 -p paramsFile -gs ignore -v newVal1 -v newVal2 -c 7

https://domain.tld/test/rot?pa1=val1&pa2=val2&paramP10=newVal1&paramP6=newVal1&paramP7=newVal1&paramP8=newVal1&paramP9=newVal1
https://domain.tld/test/rot?pa1=val1&pa2=val2&paramP1=newVal1&paramP2=newVal1&paramP3=newVal1&paramP4=newVal1&paramP5=newVal1
https://domain.tld/test/rot?pa1=val1&pa2=val2&paramP10=newVal2&paramP6=newVal2&paramP7=newVal2&paramP8=newVal2&paramP9=newVal2
https://domain.tld/test/rot?pa1=val1&pa2=val2&paramP1=newVal2&paramP2=newVal2&paramP3=newVal2&paramP4=newVal2&paramP5=newVal2
```


### Combine mode 
For each parameter in each url, either replaces the value completely or appends the new provided value as a suffix

with replace flag:
```bash
echo "https://domain.tld/test/rot?pa1=val1&pa2=val2" | x9 -p paramsFile -gs combine -v newVal1 -v newVal2 -c 7 -vs replace

https://domain.tld/test/rot?pa1=val1&pa2=newVal1
https://domain.tld/test/rot?pa1=newVal1&pa2=val2
https://domain.tld/test/rot?pa1=val1&pa2=newVal2
https://domain.tld/test/rot?pa1=newVal2&pa2=val2
```
with suffix flag:
```bash
echo "https://domain.tld/test/rot?pa1=val1&pa2=val2" | x9 -p paramsFile -gs combine -v newVal1 -v newVal2 -c 7 -vs suffix

https://domain.tld/test/rot?pa1=val1newVal1&pa2=val2
https://domain.tld/test/rot?pa1=val1&pa2=val2newVal1
https://domain.tld/test/rot?pa1=val1newVal2&pa2=val2
https://domain.tld/test/rot?pa1=val1&pa2=val2newVal2
```

# Notes
1. Multiple modes could be provided as comma-separated values (-gs combine,ignore,normal)
2. Each values should be provided separately (-v value1 -v value2)

