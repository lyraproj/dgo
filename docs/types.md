# Dgo Type syntax
The Dgo Type syntax is designed to be close to the syntax used by Go itself.

#### Primitive types  
  
|Type expression|Meaning|
|---------------|-------|
|`nil`|nil|
|`bool`|true or false|
|`true`|true|
|`false`|false|
|`string`|any string|
|`int`|any integer of any size|
|`float`|any float of any size|

#### Constrained strings

|Type expression|References|
|---------------|----------|
|`string[10,12]`|a string with a required length between 10 and 12 characters|
|`/.*abc.*/`|any string matching the regular expression|
|`"abc"`|the string "abc" verbatim|
  
#### Constrained numbers

|Type expression|References|
|---------------|----------|
|`3..28`|integer in the range 3 to 28 inclusively|
|`3...28`|integer in the range 3 to 28 with exclusive endpoint|
|`0..`|a positive integer|
|`-1.2..3.8`|a float ranging from -1.2 to 3.8|
|`-1.2...3.8`|a float ranging from -1.2 to 3.8 with exclusive endpoint|

### Arrays
#### Syntax:
`[]<element type>` or `{ <element type at position 0> [,<element type at position 1> ... ] }`

|Type expression|References an array with|
|---------------|------------------------|
|`[]int`|integers|
|`[]0..15`|integers ranging from 0 to 15 inclusively|
|`[1,10]any`|1 to 10 elements of any type|
|`[1,10]string[1]`|1 to 10 non empty strings|
|`{0..3,string,float}`|an int between 0 and 3, a string, and a float, in that order|

### Maps
#### Syntax:
`map[<key type>]<value type>`

|Sample type expression|Describes a map with|
|----------------------|--------------------|
|`map[string]int`|string keys and integer values|
|`map[string](string\|int)`|string keys and string or integer values (see anyOf below)|
|`map[string](string\|nil)`|string keys and optional string values|
|`map[string\|int]any`|string or integer keys and any type of values|
|`map[/\A[A-Z]+\z/,1,10]string[1]`|upper case string keys, non empty string values, and between 1 to 10 entries|

A map with predefined keys, where all keys are strings, is very common. Such maps are described as lists of `<key>:<value>`
associations. The `<key>` is a bit special in that it will allow identifiers that don't map to a type and treat them as
literal strings. I.e, just using `name` instead `"name"` is allowed here. 

Quoting is of course still allowed but only needed when the key conflicts with a type, i.e. `"int"` must be used to get the literal string "int" since just `int` will result in the integer type.

The characters allowed in a `<key>` identifier are $, _, 0-9, A-Z, and a-z

|Sample type expression|Describes a map with|
|----------------------|--------------------|
|`{name:string,co?:string,address:string,zip:/\d{5,5}/,city:string}`|map with named and typed entries where "co" is optional|
|`{"name":string,"co"?:string,"address":string,"zip":/\d{5,5}/,"city":string}`|same as above|

### Combinations
#### allOf syntax:
`<type>&<type>[&<type>...]`

|Sample type expression|References|
|----------------------|----------|
|`/^Ap/&string[20]`|any string that is 20 characters long and starts with the letter 'p'|

#### anyOf syntax:
`<type>|<type>[|<type>...]`

|Sample type expression|References|
|----------------------|----------|
|`"a"\|"b"\|"c"`|the string "a", "b", or "c"|
|`int\|float`|an integer or a float|
|`1\|8\|10\|16`|the integer 1, 8, 10, or 16|

### Negation
A negation matches all values that doesn't match the given type.
#### syntax:
`!<type>`

### Type Alias
New type names can be created using the assignment operator '=' which allow users to define their own
types.
```
{
  types: {
    ascii=1..127,
    slug=/^[a-z0-9-]+$/
  },
  x: map[slug]{token:ascii,value:string}
}
```
Another usage for aliases is when a type references itself. This might sound esoteric but it is in fact not that
uncommon. One obvious example is a type that describes a file system tree. Let's assume a type where each name in
a directory maps to either the mode flags of a file, or a new directory:
```
files=map[string](int|files)
```

### Type Extension
TBD, how one type can be made to extend another type, a.k.a. type inheritance.
