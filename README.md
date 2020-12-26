# go-switch

simple template filler.

### Source file

defines source patterns with flag `--sourceFile`
```
$$pattern: header
<header>
  this is header
</header>
$$end

$$pattern: footer
<footer>
  this is footer
</footer>
$$end
```

### Target file
defines target file with flag `--targetFile`
```
${header}
this is content
${footer}
```

command for filling is `go run main.go --sourceFile="from.txt" --targetFile="to.txt"`
end the result will be `new_${targetFile}` 

> works for the working directory for now

result for this example will be:

```
<header>
  this is header
</header>
this is content
<footer>
  this is footer
</footer>
```

-----

## other flags

- `--patternBegins` : defines a string for the beginning of the search pattern. 
> works for only target file.
- `--patternEnds`   : defines a string for the ending of the search pattern. 
> works for only target file.
