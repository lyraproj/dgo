# Dgo Type syntax
The Dgo Type syntax is designed to be close to the syntax used by Go itself.

The "Corresponding Puppet type" column present in the tables below is there for clarity only for those
familiar with that syntax. Dgo has no dependency on Puppet. 

#### Primitive types  
  
|Type expression|Meaning|Corresponding Puppet type|
|---------------|-------|-------------------------|
|`nil`|nil|`Undef`|
|`bool`|true or false|`Boolean`|
|`true`|true|`Boolean[true]`|
|`false`|false|`Boolean[false]`|
|`string`|any string|`String`|
|`int`|any integer of any size|`Integer`|
|`float`|any float of any size|`Float`|

#### Constrained strings

|Type expression|References|Corresponding Puppet type|
|---------------|----------|-------------------------|
|`string[10,12]`|a string with a required length of 10 to 12 characters|`String[10,12]`
|`/.*abc.*/`|any string matching the regular expression|`Pattern[/.*abc.*/]`
|`"abc"`|the string "abc" verbatim|not applicable|
  
#### Constrained numbers

|Type expression|References|Corresponding Puppet type|
|---------------|----------|-------------------------|
|`3..28`|integer in the range 3 to 28 inclusively|`Integer[3,28]`
|`3...28`|integer in the range 3 to 28 with exclusive endpoint|`Integer[3,27]`
|`0..`|a positive integer|`Integer[0]`
|`-1.2..3.8`|a float ranging from -1.2 to 3.8|`Float[-1.2, 3.8]` 

The exclusive endpoint is not permitted for floats since the predecessor is impossible to compute.
### Arrays
#### Syntax:
`[]<element type>` or `{ <element type at position 0> [,<element type at position 1> ... ] }`

|Type expression|References an array with|Corresponding Puppet type|
|---------------|------------------------|-------------------------|
|`[]int`|integers|`Array[Integer]`|
|`[]0..15`|integers ranging from 0 to 15 inclusively|`Array[Integer[0,15]]`|
|`[1,10]any`|1 to 10 elements of any type|`Array[1,10]`|
|`[1,10]string[1]`|1 to 10 non empty strings|`Array[String[1],1,10]`|
|`{0..3,string,float}`|an int between 0 and 3, a string, and a float, in that order|`Tuple[Integer[0,3],String,Float]`|

### Maps
#### Syntax:
`map[<key type>]<value type>`

|Sample type expression|Describes a map with|Corresponding Puppet type|
|----------------------|--------------------|-------------------------|
|`map[string]int`|string keys and integer values|`Hash[String,Integer]`|
|`map[string](string\|int)`|string keys and string or integer values (see anyOf below)|`Hash[String,Variant[String,Integer]]`|
|`map[string](string\|nil)`|string keys and optional string values|`Hash[String,Optional[String]]`|
|`map[string\|int]any`|string or integer keys and any type of values|`Hash[Variant[String,Integer],Any]`|
|`map[/\A[A-Z]+\z/,1,10]string[1]`|upper case string keys, non empty string values, and between 1 to 10 entries|`Hash[Pattern[/\A[A-Z]+\z/],String[1],1,10]`|

A map with predefined keys, where all keys are strings, is very common. Such maps are described as lists of <key>:<value>
associations. The <key> is a bit special in that it will allow identifiers that doesn't map to a type and treat them as
literal strings. I.e, just using `name` instead `"name"` is allowed here. Quoting is of course still allowed but only
needed when the key conflicts with a type, i.e. `"int"` must be used to get the literal string "int"
since just `int` will result in the integer type.

The allowed in a <key> identifier are $, _, 0-9, A-Z, and a-z

|Sample type expression|Describes a map with|Corresponding Puppet type|
|----------------------|--------------------|-------------------------|
|`{name:string,co?:string,address:string,zip:/\d{5,5}/,city:string}`|map with named and typed entries where "co" is optional|`Struct[name=>String,Optional[co]=>String,address=>String,zip=>Pattern[/\d{5,5}/],city=>String]`
|`{"name":string,"co"?:string,"address":string,"zip":/\d{5,5}/,"city":string}`|same as above|same as above

### Combinations
#### allOf syntax:
`<type>&<type>[&<type>...]`

|Sample type expression|References|Corresponding Puppet type|
|----------------------|----------|-------------------------|
|`/^Ap/&string[20]`|any string that is 20 characters long and starts with the letter 'p'|not applicable|

#### anyOf syntax:
`<type>|<type>[|<type>...]`

|Sample type expression|References|Corresponding Puppet type|
|----------------------|----------|-------------------------|
|`"a"\|"b"\|"c"`|the string "a", "b", or "c"|`Enum[a,b,c]`|
|`int\|float`|an integer or a float|`Variant[Integer,Float]`|
|`1\|8\|10\|16`|the integer 1, 8, 10, or 16|`Variant[Integer[1,1],Integer[8,8],Integer[10,10],Integer[16,16]]`|

### Negation
A negation matches all values that doesn't match the given type.
#### syntax:
`!<type>` (not applicable in Puppet)

### Type Alias
New type names can be created using the assignment operator '='. Here's an example of a self
referencing type that describes a dictionary of strings which is nested to aribtrary depth:
```
dir=map[string](string|dir)
```

### Type Extension
TBD, how one type can be made to extend another type, a.k.a. type inheritance.
